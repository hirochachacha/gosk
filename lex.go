// Copyright (c) 2014 Hiroshi Ioka. All rights reserved.

// this source code originally come from text/template/parse/lex.go
// http://code.google.com/p/go/source/browse/src/pkg/text/template/parse/lex.go

// Copyright (c) 2012 The Go Authors. All rights reserved.
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are
// met:
//
//    * Redistributions of source code must retain the above copyright
// notice, this list of conditions and the following disclaimer.
//    * Redistributions in binary form must reproduce the above
// copyright notice, this list of conditions and the following disclaimer
// in the documentation and/or other materials provided with the
// distribution.
//    * Neither the name of Google Inc. nor the names of its
// contributors may be used to endorse or promote products derived from
// this software without specific prior written permission.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
// "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
// LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
// A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
// OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
// SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
// LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
// DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
// THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
// (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
// OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

package gosk

import (
	"errors"
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"
)

type Pos int

// item represents a token or text string returned from the scanner.
type Item struct {
	typ itemType // The type of this item.
	pos Pos      // The starting position, in bytes, of this item in the input string.
	val string   // The value of this item.
}

func (i Item) String() string {
	switch i.typ {
	case itemError:
		return fmt.Sprintf("error: %d %s", i.pos, i.val)
	case itemBool:
		return fmt.Sprintf("bool: %s", i.val)
	case itemChar:
		return fmt.Sprintf("char: %s", i.val)
	case itemRune:
		return fmt.Sprintf("rune: %s", i.val)
	case itemImaginary:
		return fmt.Sprintf("imaginary: %s", i.val)
	case itemNewLine:
		return "NEWLINE"
	case itemEOF:
		return "EOF"
	case itemField:
		return fmt.Sprintf("field: %s", i.val)
	case itemGlobalIdentifier:
		return fmt.Sprintf("global_identifier: %s", i.val)
	case itemIdentifier:
		return fmt.Sprintf("identifier: %s", i.val)
	case itemLeftBrace:
		return fmt.Sprintf("left_brace: %s", i.val)
	case itemLeftParen:
		return fmt.Sprintf("left_paren: %s", i.val)
	case itemInt:
		return fmt.Sprintf("int: %s", i.val)
	case itemFloat:
		return fmt.Sprintf("float: %s", i.val)
	case itemRawString:
		return fmt.Sprintf("raw_string: %s", i.val)
	case itemRightBrace:
		return fmt.Sprintf("right_brace: %s", i.val)
	case itemRightParen:
		return fmt.Sprintf("right_paren: %s", i.val)
	case itemSpace:
		return "SPACE"
	case itemString:
		return fmt.Sprintf("string: %s", i.val)
	case itemIndent:
		return "INDENT"
	case itemDedent:
		return "DEDENT"
	case itemOrOr:
		return "OROR"
	case itemAndAnd:
		return "ANDAND"
	case itemEqual:
		return "EQ"
	case itemNonEqual:
		return "NE"
	case itemLessThan:
		return "LT"
	case itemLessThanOrEqual:
		return "LE"
	case itemGreaterThan:
		return "GT"
	case itemGreaterThanOrEqual:
		return "GE"
	case itemLeftShift:
		return "LSH"
	case itemRightShift:
		return "RSH"
	case itemAndNot:
		return "ANDNOT"
	}

	if i.typ > itemKeyword {
		if i.typ == itemReturn {
			return "<return>"
		}
		return fmt.Sprintf("<%s>", i.val)
	}

	// unreachable
	return fmt.Sprintf("unknown: %q", i.val)
}

// itemType identifies the type of lex items.
type itemType int

const (
	itemError itemType = iota // error occurred; value is text of error
	itemBool                  // boolean constant
	itemChar                  // printable ASCII character; grab bag for comma etc.
	itemRune                  // character constant
	itemInt
	itemFloat
	itemImaginary // imaginary constant (2i)
	itemNewLine
	itemEOF
	itemField            // alphanumeric identifier starting with '.'
	itemGlobalIdentifier // alphanumeric identifier starting with '$'
	itemIdentifier       // alphanumeric identifier not starting with '.'
	itemLeftBrace
	itemLeftParen // '(' inside action
	itemRawString // raw quoted string (includes quotes)
	itemRightBrace
	itemRightParen // ')' inside action
	itemSpace      // run of spaces separating arguments
	itemString     // quoted string (includes quotes)

	itemIndent // only internal use
	itemDedent // only internal use

	itemOrOr
	itemAndAnd
	itemEqual
	itemNonEqual
	itemLessThan
	itemLessThanOrEqual
	itemGreaterThan
	itemGreaterThanOrEqual
	itemLeftShift
	itemRightShift
	itemAndNot

	itemKeyword // used only to delimit the keywords

	itemDot    // the cursor, spelled '.'
	itemDollar // the cursor, spelled '$'
	itemReturn
	itemIf   // if keyword
	itemElse // else keyword
	itemSwitch
	itemCase
	itemDefault
	itemNil // the untyped nil constant, easiest to treat as a keyword
)

var (
	indentItem = Item{typ: itemIndent}
	dedentItem = Item{typ: itemDedent}
)

var key = map[string]itemType{
	".":       itemDot,
	"return":  itemReturn,
	"if":      itemIf,
	"else":    itemElse,
	"switch":  itemSwitch,
	"case":    itemCase,
	"default": itemDefault,
	"nil":     itemNil,
}

const eof = -1

// stateFn represents the state of the scanner as a function that returns the next state.
type stateFn func(*lexer) stateFn

// lexer holds the state of the scanner.
type lexer struct {
	name       string    // the name of the input; used only for error reports
	input      string    // the string being scanned
	state      stateFn   // the next lexing function to enter
	pos        Pos       // current position in the input
	start      Pos       // start position of this item
	width      Pos       // width of last rune read from input
	items      chan Item // channel of scanned items
	parenDepth int       // nesting depth of ( ) exprs
	blockDepth int
}

// next returns the next rune in the input.
func (l *lexer) next() rune {
	if int(l.pos) >= len(l.input) {
		l.width = 0
		return eof
	}
	r, w := utf8.DecodeRuneInString(l.input[l.pos:])
	l.width = Pos(w)
	l.pos += l.width
	return r
}

// peek returns but does not consume the next rune in the input.
func (l *lexer) peek() rune {
	r := l.next()
	l.backup()
	return r
}

// backup steps back one rune. Can only be called once per call of next.
func (l *lexer) backup() {
	l.pos -= l.width
}

// emit passes an item back to the client.
func (l *lexer) emit(t itemType) {
	l.items <- Item{t, l.start, l.input[l.start:l.pos]}
	l.start = l.pos
}

// ignore skips over the pending input before this point.
func (l *lexer) ignore() {
	l.start = l.pos
}

func (l *lexer) acceptUntil(pred func(rune) bool) bool {
	isSuccess := false
	var r rune
	for {
		r = l.next()
		if !pred(r) {
			break
		}
		isSuccess = true
	}
	l.backup()
	return isSuccess
}

func (l *lexer) acceptWord() (string, error) {
	isSuccess := l.acceptUntil(isAlphaNumeric)

	word := l.input[l.start:l.pos]
	if !isSuccess || !l.atTerminator() {
		return "", errors.New(fmt.Sprintf("bad character %#U", l.input[l.pos-1]))
	}
	return word, nil
}

// accept consumes the next rune if it's from the valid set.
func (l *lexer) accept(valid string) bool {
	if strings.IndexRune(valid, l.next()) >= 0 {
		return true
	}
	l.backup()
	return false
}

// acceptUntilRuneIn consumes a run of runes from the valid set.
func (l *lexer) acceptUntilRuneIn(valid string) bool {
	pred := func(r rune) bool {
		return strings.IndexRune(valid, r) >= 0
	}
	return l.acceptUntil(pred)
}

func (l *lexer) errorf(format string, args ...interface{}) stateFn {
	l.items <- Item{itemError, l.start, fmt.Sprintf(format, args...)}
	return nil
}

func (l *lexer) error(err error) stateFn {
	l.items <- Item{itemError, l.start, err.Error()}
	return nil
}

// lex creates a new scanner for the input string.
func lex(name, input string) *lexer {
	l := &lexer{
		name:  name,
		input: input,
		items: make(chan Item),
	}
	go l.run()
	return l
}

// run runs the state machine for the lexer.
func (l *lexer) run() {
	for l.state = lexMain; l.state != nil; {
		l.state = l.state(l)
	}
	close(l.items)
}

// state functions

const (
	lineComment  = "//"
	leftComment  = "/*"
	rightComment = "*/"
)

func lexMain(l *lexer) stateFn {
	switch r := l.next(); {
	case r == eof:
		l.emit(itemNewLine)
		l.emit(itemEOF)
		return nil
	case isEndOfLine(r):
		l.emit(itemNewLine)
	case isSpace(r):
		return lexSpace
	case r == '/':
		r = l.peek()
		switch r {
		case '/':
			l.backup()
			return lexComment
		case '*':
			l.backup()
			return lexMultiComment
		default:
			l.emit(itemChar)
			return lexMain
		}
	case r == '"':
		return lexQuote
	case r == '`':
		return lexRawQuote
	case r == '\'':
		return lexRune
	case r == '.':
		if l.pos < Pos(len(l.input)) {
			r := l.input[l.pos]
			if r < '0' || '9' < r {
				return lexField
			}
		}
		l.backup()
		return mkLexNumber(true)
	case r == '$':
		if l.pos < Pos(len(l.input)) {
			r := l.input[l.pos]
			if r < '0' || '9' < r {
				return lexGlobalIdentifier
			}
		}
		l.backup()
		return mkLexNumber(true)
	case r == '+' || r == '-' || ('0' <= r && r <= '9'):
		l.backup()
		return mkLexNumber(false)
	case isAlphaNumeric(r):
		l.backup()
		return lexIdentifier
	case r == '!':
		switch l.next() {
		case '=':
			l.emit(itemNonEqual)
		default:
			l.backup()
			l.emit(itemChar)
		}
		return lexMain
	case r == '=':
		switch l.next() {
		case '=':
			l.emit(itemEqual)
		default:
			l.backup()
			l.emit(itemChar)
		}
		return lexMain
	case r == '<':
		switch l.next() {
		case '-':
			l.emit(itemReturn)
		case '=':
			l.emit(itemLessThanOrEqual)
		case '<':
			l.emit(itemLeftShift)
		default:
			l.backup()
			l.emit(itemLessThan)
		}
		return lexMain
	case r == '>':
		switch l.next() {
		case '=':
			l.emit(itemGreaterThanOrEqual)
		case '>':
			l.emit(itemRightShift)
		default:
			l.backup()
			l.emit(itemGreaterThan)
		}
		return lexMain
	case r == '&':
		switch l.next() {
		case '^':
			l.emit(itemAndNot)
		default:
			l.backup()
			l.emit(itemChar)
		}
		return lexMain
	case r == '{':
		l.emit(itemLeftBrace)
		l.blockDepth++
		return lexMain
	case r == '}':
		l.emit(itemRightBrace)
		l.blockDepth--
		if l.blockDepth < 0 {
			return l.errorf("unexpected right brace %#U", r)
		}
		return lexMain
	case r == '(':
		l.emit(itemLeftParen)
		l.parenDepth++
		return lexMain
	case r == ')':
		l.emit(itemRightParen)
		l.parenDepth--
		if l.parenDepth < 0 {
			return l.errorf("unexpected right paren %#U", r)
		}
		return lexMain
	case r <= unicode.MaxASCII && unicode.IsPrint(r):
		l.emit(itemChar)
		return lexMain
	default:
		return l.errorf("unrecognized character in action: %#U", r)
	}
	return lexMain
}

func lexComment(l *lexer) stateFn {
	l.pos += Pos(len(lineComment))
	i := strings.IndexAny(l.input[l.pos:], "\r\n")
	if i < 0 {
		// maybe eof, exit
		return nil
	}
	l.pos += Pos(i)
	l.ignore()
	return lexMain
}

// lexMultiComment scans a comment. The left comment marker is known to be present.
func lexMultiComment(l *lexer) stateFn {
	l.pos += Pos(len(leftComment))
	i := strings.Index(l.input[l.pos:], rightComment)
	if i < 0 {
		return l.errorf("unclosed comment")
	}
	l.pos += Pos(i + len(rightComment))
	l.ignore()
	return lexMain
}

// lexSpace scans a run of space characters.
// One space has already been seen.
func lexSpace(l *lexer) stateFn {
	l.emit(itemSpace)
	return lexMain
}

// lexGlobalIdentifier scans a field: $Alphanumeric.
// The $ has been scanned.
func lexGlobalIdentifier(l *lexer) stateFn {
	if l.atTerminator() {
		l.emit(itemDollar)
		return lexMain
	}

	_, err := l.acceptWord()
	if err != nil {
		return l.error(err)
	}

	l.emit(itemGlobalIdentifier)
	return lexMain
}

// lexIdentifier scans an alphanumeric.
func lexIdentifier(l *lexer) stateFn {
	word, err := l.acceptWord()
	if err != nil {
		return l.error(err)
	}

	switch {
	case key[word] > itemKeyword:
		l.emit(key[word])
	case word[0] == '.':
		l.emit(itemField)
	case word == "true", word == "false":
		l.emit(itemBool)
	default:
		l.emit(itemIdentifier)
	}
	return lexMain
}

// lexField scans a field: .Alphanumeric.
// The . has been scanned.
func lexField(l *lexer) stateFn {
	if l.atTerminator() {
		l.emit(itemDot)
		return lexMain
	}

	_, err := l.acceptWord()
	if err != nil {
		return l.error(err)
	}

	l.emit(itemField)
	return lexMain
}

// atTerminator reports whether the input is at valid termination character to
// appear after an identifier. Breaks .X.Y into two pieces. Also catches cases
// like "$x+2" not being acceptable without a space, in case we decide one
// day to implement arithmetic.
func (l *lexer) atTerminator() bool {
	r := l.peek()
	if isSpace(r) || isEndOfLine(r) {
		return true
	}
	switch r {
	case eof, '.', ',', '|', ':', ')', '(', '[', ']':
		return true
	}
	return false
}

// lexRune scans a character constant. The initial quote is already
// scanned. Syntax checking is done by the parser.
func lexRune(l *lexer) stateFn {
Loop:
	for {
		switch l.next() {
		case '\\':
			if r := l.next(); r != eof && r != '\n' {
				break
			}
			fallthrough
		case eof, '\n':
			return l.errorf("unterminated character constant")
		case '\'':
			break Loop
		}
	}
	l.emit(itemRune)
	return lexMain
}

func mkLexNumber(isFloat bool) func(l *lexer) stateFn {
	// lexNumber scans a number: decimal, octal, hex, float, or imaginary. This
	// isn't a perfect number scanner - for instance it accepts "." and "0x0.2"
	// and "089" - but when it's wrong the input is invalid and the parser (via
	// strconv) will notice.
	return func(l *lexer) stateFn {
		// Optional leading sign.
		onlySign := l.accept("+-")
		// Is it hex?
		digits := "0123456789"
		if l.accept("0") {
			onlySign = false
			if l.accept("xX") {
				digits = "0123456789abcdefABCDEF"
			}
		}

		if l.acceptUntilRuneIn(digits) {
			onlySign = false
		}

		if l.accept(".") {
			onlySign = false
			isFloat = true
			l.acceptUntilRuneIn(digits)
		}
		if l.accept("eE") {
			onlySign = false
			l.accept("+-")
			l.acceptUntilRuneIn("0123456789")
		}
		// Is it imaginary?
		isImaginary := l.accept("i")
		// Next thing mustn't be alphanumeric.
		if isAlphaNumeric(l.peek()) {
			l.next()
			return l.errorf("bad number syntax: %q", l.input[l.start:l.pos])
		}
		if onlySign {
			l.emit(itemChar)
			return lexMain
		}

		if isImaginary {
			l.emit(itemImaginary)
		} else {
			if isFloat {
				l.emit(itemFloat)
			} else {
				l.emit(itemInt)
			}
		}

		if sign := l.next(); sign == '+' || sign == '-' {
			l.emit(itemChar)
		} else {
			l.backup()
		}

		return lexMain
	}
}

// lexQuote scans a quoted string.
func lexQuote(l *lexer) stateFn {
Loop:
	for {
		switch l.next() {
		case '\\':
			if r := l.next(); r != eof && r != '\n' {
				break
			}
			fallthrough
		case eof, '\n':
			return l.errorf("unterminated quoted string")
		case '"':
			break Loop
		}
	}
	l.emit(itemString)
	return lexMain
}

// lexRawQuote scans a raw quoted string.
func lexRawQuote(l *lexer) stateFn {
Loop:
	for {
		switch l.next() {
		case eof:
			return l.errorf("unterminated raw quoted string")
		case '`':
			break Loop
		}
	}
	l.emit(itemRawString)
	return lexMain
}

// isSpace reports whether r is a space character.
func isSpace(r rune) bool {
	return r == ' ' || r == '\t'
}

// isEndOfLine reports whether r is an end-of-line character.
func isEndOfLine(r rune) bool {
	return r == '\r' || r == '\n'
}

// isAlphaNumeric reports whether r is an alphabetic, digit, or underscore.
func isAlphaNumeric(r rune) bool {
	return r == '_' || unicode.IsLetter(r) || unicode.IsDigit(r)
}
