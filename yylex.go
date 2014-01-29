package main

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"unicode/utf8"
)

const MaxIndentDepth = 100

type Lexer struct {
	chunk       chan Item
	internal    *lexer
	indent      int
	depth       int
	isEndOfLine bool
}

func NewLexer(name, input string) *Lexer {
	lexer := new(Lexer)
	lexer.internal = lex(name, input)
	lexer.chunk = make(chan Item, MaxIndentDepth)
	lexer.isEndOfLine = true
	return lexer
}

func (lexer *Lexer) nextItem() Item {
	var item Item

	if lexer.isEndOfLine {
		lexer.isEndOfLine = false

		nspaces := 0
		var diffSpace int
		var diffDepth int
		var isBlankLine bool

		for {
			item = <-lexer.internal.items
			if item.typ != itemSpace {
				break
			}
			nspaces++
		}

		diffSpace = nspaces - lexer.indent*lexer.depth

		if diffSpace == 0 {
			return item
		} else {
			switch {
			case lexer.indent == 0:
				lexer.indent = diffSpace
				fallthrough
			case lexer.indent > 0:
				if diffSpace%lexer.indent != 0 {
					return Item{itemError, 0, "unaligned indent"}
				}

				isBlankLine = item.typ == itemNewLine

				if isBlankLine && diffSpace < 0 {
					diffDepth = -1
				} else {
					diffDepth = diffSpace / lexer.indent
				}

				switch {
				case diffDepth > 0:
					if diffDepth != 1 {
						return Item{itemError, 0, "too deep indent"}
					}
					lexer.chunk <- indentItem
				case diffDepth < 0:
					for i := 0; i > diffDepth; i-- {
						lexer.chunk <- dedentItem
					}
				}

				lexer.depth += diffDepth
				lexer.chunk <- item
				item = <-lexer.chunk
			case lexer.indent < 0:
				item = Item{itemError, 0, "unknown indent"}
			}
		}
	} else {
		for {
			if len(lexer.chunk) > 0 {
				item = <-lexer.chunk
			} else {
				item = <-lexer.internal.items
			}
			if item.typ != itemSpace {
				if item.typ == itemNewLine {
					lexer.isEndOfLine = true
				}
				break
			}
		}
	}
	return item
}

func (lexer *Lexer) Error(s string) {
	fmt.Printf("syntax error: %s\n", s)
}

func (lexer *Lexer) Debug() {
	var item Item
	for {
		item = lexer.nextItem()
		fmt.Println(item)
		if item.typ == itemEOF {
			break
		}
	}
}

func (lexer *Lexer) Lex(sym *yySymType) int {
	item := lexer.nextItem()
	switch item.typ {
	case itemError:
		sym.node = &Node{
			typ: nodeLiteral,
			val: reflect.ValueOf(errors.New(item.val)),
		}
		return Error
	case itemBool:
		sym.node = &Node{
			typ: nodeLiteral,
			val: reflect.ValueOf(item.val == "true"),
		}
		return Bool
	case itemChar:
		r, _ := utf8.DecodeRuneInString(item.val)
		return int(r)
	case itemRune:
		r, _ := utf8.DecodeRuneInString(item.val[1:])
		sym.node = &Node{
			typ: nodeLiteral,
			val: reflect.ValueOf(r),
		}
		return Rune
	case itemImaginary:
		f, err := strconv.ParseFloat(item.val[:len(item.val)-1], 64)
		if err != nil {
			sym.node = &Node{
				typ: nodeLiteral,
				val: reflect.ValueOf(err),
			}
			return Error
		}

		sym.node = &Node{
			typ: nodeLiteral,
			val: reflect.ValueOf(complex(0, f)),
		}
		return Imaginary
	case itemNewLine:
		return int('\n')
	case itemEOF:
		return 0
	case itemField:
		sym.node = &Node{
			typ: nodeField,
			val: reflect.ValueOf(item.val),
		}
		return Field
	case itemIdentifier:
		sym.node = &Node{
			typ: nodeIdentifier,
			val: reflect.ValueOf(item.val),
		}
		return Identifier
	case itemLeftBrace:
		return int('{')
	case itemLeftParen:
		return int('(')
	case itemInt:
		f, err := strconv.ParseInt(item.val, 10, 0)
		if err != nil {
			sym.node = &Node{
				typ: nodeLiteral,
				val: reflect.ValueOf(err),
			}
			return Error
		}

		sym.node = &Node{
			typ: nodeLiteral,
			val: reflect.ValueOf(f),
		}
		return Int
	case itemFloat:
		f, err := strconv.ParseFloat(item.val, 64)
		if err != nil {
			sym.node = &Node{
				typ: nodeLiteral,
				val: reflect.ValueOf(err),
			}
			return Error
		}

		sym.node = &Node{
			typ: nodeLiteral,
			val: reflect.ValueOf(f),
		}
		return Float
	case itemRawString:
		sym.node = &Node{
			typ: nodeLiteral,
			val: reflect.ValueOf(item.val[1 : len(item.val)-1]),
		}
		return RawString
	case itemRightBrace:
		return int('}')
	case itemRightParen:
		return int(')')
	case itemIndent:
		return INDENT
	case itemDedent:
		return DEDENT
	case itemString:
		sym.node = &Node{
			typ: nodeLiteral,
			val: reflect.ValueOf(item.val[1 : len(item.val)-1]),
		}
		return String
	case itemOrOr:
		return OROR
	case itemAndAnd:
		return ANDAND
	case itemEqual:
		return EQ
	case itemNonEqual:
		return NE
	case itemLessThan:
		return LT
	case itemLessThanOrEqual:
		return LE
	case itemGreaterThan:
		return GT
	case itemGreaterThanOrEqual:
		return GE
	case itemLeftShift:
		return LSH
	case itemRightShift:
		return RSH
	case itemAndNot:
		return ANDNOT
	case itemDot:
		return int('.')
	case itemReturn:
		return Return
	case itemIf:
		return If
	case itemElse:
		return Else
	case itemSwitch:
		return Switch
	case itemCase:
		return Case
	case itemDefault:
		return Default
	case itemNil:
		sym.node = &Node{
			typ: nodeLiteral,
			val: reflect.ValueOf(nil),
		}
		return Nil
	case itemRange:
		return Range
	}

	// unreachable
	return 0
}
