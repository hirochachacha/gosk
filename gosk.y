%{
package main

import "reflect"

%}

%left OROR                             /* '||' */
%left ANDAND                           /* '&&' */
%left EQ NE LT LE GT GE                /* '==' '!=' '<' '<=' '>' '>=' */
%left '+' '-' '|' '^'
%left '*' '/' '%' '&' LSH RSH ANDNOT   /* '<<' '>>' '&^' */

%union{
  node *Node
  nodelist *NodeList
  op opType
}

%start grammar

%token            Return If Else Switch Case Default Range
%token <node>     Error Bool Rune Imaginary Field Identifier Int Float RawString String Nil
%type  <node>     identifier literal field
%type  <node>     index
%type  <nodelist> slice call

%type  <node>     assign_stmt return_stmt
%type  <node>     if_stmt switch_stmt assign_block_stmt
%type  <node>     case_stmt default_stmt
%type  <nodelist> case_stmts

%type  <node>     simple_stmt compound_stmt
%type  <node>     stmt
%type  <nodelist> stmts

%token            DEDENT INDENT
%token            OROR ANDAND EQ NE LT LE GT GE LSH RSH ANDNOT
%type  <op>       rel_op add_op mul_op
%type  <op>       binary_op unary_op
%type  <node>     primary_expr
%type  <node>     unary_expr
%type  <node>     expr
%type  <nodelist> args exprs _exprs

%type  <node> block

%%

grammar : stmts {
        rootNode = &Node{
          typ: nodeRoot,
          list: $1,
        }
};

stmts : stmts stmt {
      if $2 != nil {
        $1.append($2)
        $$ = $1
      }
}
      | {
      $$ = NewNodeList()
};

stmt : simple_stmt '\n' {
     $$ = $1
}
     | simple_stmt block {
     $1.block = $2
     $$ = $1
}
     | compound_stmt
     | '\n' {
     $$ = nil
};

simple_stmt : return_stmt
            | assign_stmt
;

assign_stmt : identifier "=" expr {
            $$ = &Node{
              typ: nodeAssign,
              left: $1,
              right: $3,
            }
};

return_stmt : Return expr {
              $$ = &Node{
                typ: nodeReturn,
                left: $2,
              }
};

compound_stmt : if_stmt
              | switch_stmt
              | assign_block_stmt
;

if_stmt : If expr block {
        $$ = &Node{
          typ: nodeIf,
          cond: $2,
          left: $3,
        }
}
        | If expr block Else block {
        $$ = &Node{
          typ: nodeIf,
          cond: $2,
          left: $3,
          right: $5,
        }
};

switch_stmt : Switch expr '\n' case_stmts {
            $$ = &Node{
              typ: nodeSwitch,
              cond: $2,
              list: $4,
            }
}
            | Switch expr '\n' case_stmts default_stmt {
            cases := $4
            cases.append($5)

            $$ = &Node{
              typ: nodeSwitch,
              cond: $2,
              list: cases,
            }
};

case_stmts : case_stmts case_stmt {
           $1.append($2)
           $$ = $1
}
           | {
           $$ = NewNodeList()
};

case_stmt : Case expr ":" block {
          $$ = &Node{
            typ: nodeCase,
            cond: $2,
            block: $4,
          }
};

default_stmt: Default ":" block {
            $$ = &Node{
              typ: nodeDefault,
              block: $3,
            }
};

assign_block_stmt : identifier "=" block {
                  $$ = &Node{
                    typ: nodeAssignBlock,
                    left: $1,
                    block: $3,
                  }
};

block : '\n' INDENT stmts DEDENT {
      $$ = &Node{
        typ: nodeBlock,
        list: $3,
      }
};

expr : "(" expr ")" { $$ = $2 }
     | unary_expr
     | expr binary_op unary_expr {
     $$ = &Node{
       typ: nodeOp,
       left: $1,
       right: $3,
       val: reflect.ValueOf($2),
     }
};

unary_expr : primary_expr
           | unary_op unary_expr {
           $$ = &Node{
             typ: nodeUnaryOp,
             left: $2,
             val: reflect.ValueOf($1),
           }
}
;

binary_op : OROR   { $$ = opOrOr }
          | ANDAND { $$ = opAndAnd }
          | rel_op
          | add_op
          | mul_op
;

rel_op : EQ { $$ = opEqual }
       | NE { $$ = opNonEqual }
       | LT { $$ = opLessThan }
       | LE { $$ = opLessThanOrEqual }
       | GT { $$ = opGreaterThan }
       | GE { $$ = opGreaterThanOrEqual }
;

add_op : "+" { $$ = opPlus }
       | "-" { $$ = opMinus }
       | "|" { $$ = opOr }
;

mul_op : "*"    { $$ = opMulti }
       | "/"    { $$ = opDivide }
       | "%"    { $$ = opMod }
       | LSH    { $$ = opLeftShift }
       | RSH    { $$ = opRightShift }
       | "&"    { $$ = opAnd }
       | ANDNOT { $$ = opAndNot }
;

unary_op : "-" { $$ = opNegate }
         | "!" { $$ = opNot }
         | "^" { $$ = opBitNot }
;

primary_expr : literal
             | identifier
             | field
             | primary_expr field {
             $$ = &Node{
               typ: nodeSelector,
               left: $1,
               val: $2.val,
             }
}
             | primary_expr index {
             $$ = &Node{
               typ: nodeIndex,
               left: $1,
               right: $2,
             }
}
             | primary_expr slice {
             $$ = &Node{
               typ: nodeSlice,
               left: $1,
               list: $2,
             }
}
             | primary_expr call {
             $$ = &Node{
               typ: nodeCall,
               left: $1,
               list: $2,
             }
};

identifier : Identifier;

field : Field;

index : "[" expr "]" { $$ = $2 };

slice : "[" ":" "]"                    { $$ = NewNodeList(nil, nil, nil) }
      | "[" expr ":" "]"               { $$ = NewNodeList($2, nil, nil) }
      | "[" ":" expr "]"               { $$ = NewNodeList(nil, $3, nil) }
      | "[" expr ":" expr "]"          { $$ = NewNodeList($2, $4, nil) }
      | "[" ":" expr ":" expr "]"      { $$ = NewNodeList(nil, $3, $5) }
      | "[" expr ":" expr ":" expr "]" { $$ = NewNodeList($2, $4, $6) }
;

call : "(" ")"           { $$ = nil }
     | "(" args ")"      { $$ = $2 }
     | "(" args ","  ")" { $$ = $2 }
;

args : exprs
     | exprs "..." {
     $1.applyable()
     $$ = $1
};

exprs : expr _exprs {
      $2.prepend($1)
      $$ = $2
};

_exprs : _exprs "," expr {
       $1.append($3)
       $$ = $1
}
       | {
       $$ = NewNodeList()
};

literal : Bool
        | Rune
        | String
        | RawString
        | Int
        | Float
        | Imaginary
        | Nil
;

%%
