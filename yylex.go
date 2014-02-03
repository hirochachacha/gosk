package gosk

import (
	"fmt"
	"reflect"
	"strconv"
	"unicode/utf8"
)

const MaxIndentDepth = 100

type Lexer struct {
	chunk        chan Item
	internal     *lexer
	indentLength int
	indentDepth  int
	isEndOfLine  bool
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
		var deltaSpace int
		var deltaDepth int
		var isBlankLine bool

		for {
			item = <-lexer.internal.items
			if item.typ != itemSpace {
				if item.typ == itemNewLine {
					lexer.isEndOfLine = true
					isBlankLine = true
				}
				break
			}
			nspaces++
		}

		deltaSpace = nspaces - lexer.indentLength*lexer.indentDepth

		if deltaSpace == 0 {
			return item
		} else {
			switch {
			case lexer.indentLength == 0:
				if isBlankLine {
					return item
				}
				lexer.indentLength = deltaSpace
				fallthrough
			case lexer.indentLength > 0:
				if isBlankLine {
					deltaDepth = -1
				} else {
					if deltaSpace%lexer.indentLength != 0 {
						return Item{itemError, 0, fmt.Sprintf("unaligned indentation: expected %d aligned, but %d", lexer.indentLength, deltaSpace)}
					}
					deltaDepth = deltaSpace / lexer.indentLength
				}

				switch {
				case deltaDepth > 0:
					if deltaDepth != 1 {
						return Item{itemError, 0, fmt.Sprintf("deep indentation: expected %d, but %d", lexer.indentLength, deltaSpace)}
					}
					lexer.chunk <- indentItem
				case deltaDepth < 0:
					for i := 0; i > deltaDepth; i-- {
						lexer.chunk <- dedentItem
					}
				}

				lexer.indentDepth += deltaDepth
				lexer.chunk <- item
				item = <-lexer.chunk
			case lexer.indentLength < 0:
				return Item{itemError, 0, "unknown indentation"}
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
	fmt.Printf("error: %s\n", s)
}

func (lexer *Lexer) Debug() {
	for {
		item := lexer.nextItem()
		fmt.Println(item)
		if item.typ == itemEOF || item.typ == itemError {
			break
		}
	}
}

func (lexer *Lexer) Lex(sym *yySymType) int {
	item := lexer.nextItem()
	switch item.typ {
	case itemError:
		panic(item)
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
		val := item.val[:len(item.val)-1]
		f, err := strconv.ParseFloat(val, 64)
		if err != nil {
			panic(fmt.Sprintf("can't convert %s to float", val))
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
	case itemGlobalIdentifier:
		sym.node = &Node{
			typ: nodeGlobalIdentifier,
			val: reflect.ValueOf(item.val),
		}
		return GlobalIdentifer
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
			panic(fmt.Sprintf("can't convert %s to int", item.val))
		}

		sym.node = &Node{
			typ: nodeLiteral,
			val: reflect.ValueOf(f),
		}
		return Int
	case itemFloat:
		f, err := strconv.ParseFloat(item.val, 64)
		if err != nil {
			panic(fmt.Sprintf("can't convert %s to float", item.val))
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
	}

	// unreachable
	return 0
}
