//line gosk.y:2
package main

import __yyfmt__ "fmt"

//line gosk.y:2
import "reflect"

//line gosk.y:14
type yySymType struct {
	yys      int
	node     *Node
	nodelist *NodeList
	op       opType
}

const OROR = 57346
const ANDAND = 57347
const EQ = 57348
const NE = 57349
const LT = 57350
const LE = 57351
const GT = 57352
const GE = 57353
const LSH = 57354
const RSH = 57355
const ANDNOT = 57356
const Return = 57357
const If = 57358
const Else = 57359
const Switch = 57360
const Case = 57361
const Default = 57362
const Range = 57363
const Error = 57364
const Bool = 57365
const Rune = 57366
const Imaginary = 57367
const Field = 57368
const Identifier = 57369
const Int = 57370
const Float = 57371
const RawString = 57372
const String = 57373
const Nil = 57374
const DEDENT = 57375
const INDENT = 57376

var yyToknames = []string{
	"OROR",
	"ANDAND",
	"EQ",
	"NE",
	"LT",
	"LE",
	"GT",
	"GE",
	" +",
	" -",
	" |",
	" ^",
	" *",
	" /",
	" %",
	" &",
	"LSH",
	"RSH",
	"ANDNOT",
	"Return",
	"If",
	"Else",
	"Switch",
	"Case",
	"Default",
	"Range",
	"Error",
	"Bool",
	"Rune",
	"Imaginary",
	"Field",
	"Identifier",
	"Int",
	"Float",
	"RawString",
	"String",
	"Nil",
	"DEDENT",
	"INDENT",
}
var yyStatenames = []string{}

const yyEofCode = 1
const yyErrCode = 2
const yyMaxDepth = 200

//line gosk.y:310

//line yacctab:1
var yyExca = []int{
	-1, 1,
	1, -1,
	-2, 0,
}

const yyNprod = 86
const yyPrivate = 57344

var yyTokenNames []string
var yyStates []string

const yyLast = 516

var yyAct = []int{

	19, 18, 25, 96, 27, 13, 29, 94, 108, 107,
	109, 95, 39, 75, 17, 40, 41, 42, 87, 21,
	38, 65, 30, 31, 36, 38, 16, 34, 35, 33,
	32, 37, 71, 102, 101, 70, 2, 20, 83, 28,
	73, 74, 76, 72, 44, 45, 49, 50, 51, 52,
	53, 54, 55, 56, 57, 26, 58, 59, 60, 63,
	61, 62, 64, 79, 12, 14, 97, 15, 85, 84,
	22, 81, 86, 23, 43, 27, 16, 29, 66, 78,
	48, 13, 89, 93, 6, 112, 47, 46, 3, 98,
	111, 5, 104, 30, 31, 36, 38, 16, 34, 35,
	33, 32, 37, 110, 4, 88, 99, 113, 20, 114,
	28, 115, 103, 117, 100, 11, 10, 9, 119, 44,
	45, 49, 50, 51, 52, 53, 54, 55, 56, 57,
	7, 58, 59, 60, 63, 61, 62, 64, 44, 45,
	49, 50, 51, 52, 53, 54, 55, 56, 57, 8,
	58, 59, 60, 63, 61, 62, 64, 69, 12, 14,
	106, 15, 68, 67, 24, 105, 1, 0, 0, 0,
	16, 0, 0, 0, 0, 0, 0, 0, 6, 91,
	0, 0, 0, 0, 90, 44, 45, 49, 50, 51,
	52, 53, 54, 55, 56, 57, 0, 58, 59, 60,
	63, 61, 62, 64, 44, 45, 49, 50, 51, 52,
	53, 54, 55, 56, 57, 0, 58, 59, 60, 63,
	61, 62, 64, 0, 0, 0, 0, 0, 0, 0,
	0, 120, 0, 0, 0, 0, 0, 0, 0, 44,
	45, 49, 50, 51, 52, 53, 54, 55, 56, 57,
	118, 58, 59, 60, 63, 61, 62, 64, 0, 0,
	0, 0, 0, 0, 0, 27, 0, 29, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 80, 30, 31, 36, 38, 16, 34, 35,
	33, 32, 37, 0, 0, 0, 0, 0, 20, 0,
	28, 0, 92, 44, 45, 49, 50, 51, 52, 53,
	54, 55, 56, 57, 0, 58, 59, 60, 63, 61,
	62, 64, 0, 0, 0, 0, 0, 27, 0, 29,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 116, 30, 31, 36, 38, 16,
	34, 35, 33, 32, 37, 27, 0, 29, 0, 82,
	20, 0, 28, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 30, 31, 36, 38, 16, 34, 35,
	33, 32, 37, 0, 27, 75, 29, 0, 20, 0,
	28, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 30, 31, 36, 38, 16, 34, 35, 33,
	32, 37, 0, 0, 0, 0, 0, 20, 0, 28,
	44, 45, 49, 50, 51, 52, 53, 54, 55, 56,
	57, 0, 58, 59, 60, 63, 61, 62, 64, 44,
	45, 49, 50, 51, 52, 53, 54, 55, 56, 57,
	0, 58, 59, 60, 63, 61, 62, 64, 0, 77,
	0, 27, 0, 29, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 75, 30,
	31, 36, 38, 16, 34, 35, 33, 32, 37, 0,
	0, 0, 0, 0, 0, 0, 28, 44, 45, 49,
	50, 51, 52, 53, 54, 55, 56, 57, 0, 58,
	59, 60, 63, 61, 62, 64,
}
var yyPact = []int{

	-1000, -1000, 135, -1000, -29, -1000, -1000, -1000, -1000, -1000,
	-1000, -1000, 371, -32, 371, 371, -1000, -25, -1000, 493,
	371, -1000, -14, 448, -1000, -1000, -1000, -1000, -1000, -1000,
	-1000, -1000, -1000, -1000, -1000, -1000, -1000, -1000, -1000, 342,
	435, 416, -1000, 448, -1000, -1000, -1000, -1000, -1000, -1000,
	-1000, -1000, -1000, -1000, -1000, -1000, -1000, -1000, -1000, -1000,
	-1000, -1000, -1000, -1000, -1000, 235, -1000, -1000, -1000, -1000,
	314, -9, -1000, 493, -1000, -25, -7, -1000, 41, -1000,
	-1000, 134, 252, -1000, -40, -49, 493, -30, 6, -1000,
	-1000, 62, -1000, 115, -1000, -38, -1000, -43, -1000, -1000,
	-1000, -35, 371, -1000, 40, -1000, 371, -1000, 371, -30,
	299, -1000, 371, 200, 493, -1000, -30, 181, -1000, -1000,
	-1000,
}
var yyPgo = []int{

	0, 166, 2, 164, 55, 163, 162, 157, 149, 130,
	117, 116, 115, 114, 106, 105, 104, 91, 88, 36,
	87, 86, 80, 74, 73, 70, 19, 0, 69, 68,
	66, 1,
}
var yyR1 = []int{

	0, 1, 19, 19, 18, 18, 18, 18, 16, 16,
	8, 9, 17, 17, 17, 10, 10, 11, 11, 15,
	15, 13, 14, 12, 31, 27, 27, 27, 26, 26,
	23, 23, 23, 23, 23, 20, 20, 20, 20, 20,
	20, 21, 21, 21, 22, 22, 22, 22, 22, 22,
	22, 24, 24, 24, 25, 25, 25, 25, 25, 25,
	25, 2, 4, 5, 6, 6, 6, 6, 6, 6,
	7, 7, 7, 28, 28, 29, 30, 30, 3, 3,
	3, 3, 3, 3, 3, 3,
}
var yyR2 = []int{

	0, 1, 2, 0, 2, 2, 1, 1, 1, 1,
	3, 2, 1, 1, 1, 3, 5, 4, 5, 2,
	0, 4, 3, 3, 4, 3, 1, 3, 1, 2,
	1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
	1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
	1, 1, 1, 1, 1, 1, 1, 2, 2, 2,
	2, 1, 1, 3, 3, 4, 4, 5, 6, 7,
	2, 3, 4, 1, 2, 2, 3, 0, 1, 1,
	1, 1, 1, 1, 1, 1,
}
var yyChk = []int{

	-1000, -1, -19, -18, -16, -17, 43, -9, -8, -10,
	-11, -12, 23, -2, 24, 26, 35, 43, -31, -27,
	46, -26, -25, -24, -3, -2, -4, 13, 48, 15,
	31, 32, 39, 38, 36, 37, 33, 40, 34, 44,
	-27, -27, 42, -23, 4, 5, -20, -21, -22, 6,
	7, 8, 9, 10, 11, 12, 13, 14, 16, 17,
	18, 20, 21, 19, 22, -27, -4, -5, -6, -7,
	49, 46, -26, -27, -31, 43, -31, 43, -19, -26,
	47, -27, 45, 47, -28, -29, -27, 25, -15, 41,
	50, 45, 50, -27, 47, 51, 52, -30, -31, -14,
	-13, 28, 27, 50, -27, 50, 45, 47, 51, 45,
	-27, 50, 45, -27, -27, -31, 45, -27, 50, -31,
	50,
}
var yyDef = []int{

	3, -2, 1, 2, 0, 6, 7, 8, 9, 12,
	13, 14, 0, 0, 0, 0, 61, 4, 5, 11,
	0, 26, 28, 0, 54, 55, 56, 51, 52, 53,
	78, 79, 80, 81, 82, 83, 84, 85, 62, 0,
	0, 0, 3, 0, 30, 31, 32, 33, 34, 35,
	36, 37, 38, 39, 40, 41, 42, 43, 44, 45,
	46, 47, 48, 49, 50, 0, 57, 58, 59, 60,
	0, 0, 29, 10, 23, 0, 15, 20, 0, 27,
	25, 0, 0, 70, 0, 73, 77, 0, 17, 24,
	63, 0, 64, 0, 71, 0, 74, 75, 16, 18,
	19, 0, 0, 65, 0, 66, 0, 72, 0, 0,
	0, 67, 0, 0, 76, 22, 0, 0, 68, 21,
	69,
}
var yyTok1 = []int{

	1, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	43, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 48, 3, 3, 3, 18, 19, 3,
	46, 47, 16, 12, 51, 13, 52, 17, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 45, 3,
	3, 44, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 49, 3, 50, 15, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 14,
}
var yyTok2 = []int{

	2, 3, 4, 5, 6, 7, 8, 9, 10, 11,
	20, 21, 22, 23, 24, 25, 26, 27, 28, 29,
	30, 31, 32, 33, 34, 35, 36, 37, 38, 39,
	40, 41, 42,
}
var yyTok3 = []int{
	0,
}

//line yaccpar:1

/*	parser for yacc output	*/

var yyDebug = 0

type yyLexer interface {
	Lex(lval *yySymType) int
	Error(s string)
}

const yyFlag = -1000

func yyTokname(c int) string {
	// 4 is TOKSTART above
	if c >= 4 && c-4 < len(yyToknames) {
		if yyToknames[c-4] != "" {
			return yyToknames[c-4]
		}
	}
	return __yyfmt__.Sprintf("tok-%v", c)
}

func yyStatname(s int) string {
	if s >= 0 && s < len(yyStatenames) {
		if yyStatenames[s] != "" {
			return yyStatenames[s]
		}
	}
	return __yyfmt__.Sprintf("state-%v", s)
}

func yylex1(lex yyLexer, lval *yySymType) int {
	c := 0
	char := lex.Lex(lval)
	if char <= 0 {
		c = yyTok1[0]
		goto out
	}
	if char < len(yyTok1) {
		c = yyTok1[char]
		goto out
	}
	if char >= yyPrivate {
		if char < yyPrivate+len(yyTok2) {
			c = yyTok2[char-yyPrivate]
			goto out
		}
	}
	for i := 0; i < len(yyTok3); i += 2 {
		c = yyTok3[i+0]
		if c == char {
			c = yyTok3[i+1]
			goto out
		}
	}

out:
	if c == 0 {
		c = yyTok2[1] /* unknown char */
	}
	if yyDebug >= 3 {
		__yyfmt__.Printf("lex %s(%d)\n", yyTokname(c), uint(char))
	}
	return c
}

func yyParse(yylex yyLexer) int {
	var yyn int
	var yylval yySymType
	var yyVAL yySymType
	yyS := make([]yySymType, yyMaxDepth)

	Nerrs := 0   /* number of errors */
	Errflag := 0 /* error recovery flag */
	yystate := 0
	yychar := -1
	yyp := -1
	goto yystack

ret0:
	return 0

ret1:
	return 1

yystack:
	/* put a state and value onto the stack */
	if yyDebug >= 4 {
		__yyfmt__.Printf("char %v in %v\n", yyTokname(yychar), yyStatname(yystate))
	}

	yyp++
	if yyp >= len(yyS) {
		nyys := make([]yySymType, len(yyS)*2)
		copy(nyys, yyS)
		yyS = nyys
	}
	yyS[yyp] = yyVAL
	yyS[yyp].yys = yystate

yynewstate:
	yyn = yyPact[yystate]
	if yyn <= yyFlag {
		goto yydefault /* simple state */
	}
	if yychar < 0 {
		yychar = yylex1(yylex, &yylval)
	}
	yyn += yychar
	if yyn < 0 || yyn >= yyLast {
		goto yydefault
	}
	yyn = yyAct[yyn]
	if yyChk[yyn] == yychar { /* valid shift */
		yychar = -1
		yyVAL = yylval
		yystate = yyn
		if Errflag > 0 {
			Errflag--
		}
		goto yystack
	}

yydefault:
	/* default state action */
	yyn = yyDef[yystate]
	if yyn == -2 {
		if yychar < 0 {
			yychar = yylex1(yylex, &yylval)
		}

		/* look through exception table */
		xi := 0
		for {
			if yyExca[xi+0] == -1 && yyExca[xi+1] == yystate {
				break
			}
			xi += 2
		}
		for xi += 2; ; xi += 2 {
			yyn = yyExca[xi+0]
			if yyn < 0 || yyn == yychar {
				break
			}
		}
		yyn = yyExca[xi+1]
		if yyn < 0 {
			goto ret0
		}
	}
	if yyn == 0 {
		/* error ... attempt to resume parsing */
		switch Errflag {
		case 0: /* brand new error */
			yylex.Error("syntax error")
			Nerrs++
			if yyDebug >= 1 {
				__yyfmt__.Printf("%s", yyStatname(yystate))
				__yyfmt__.Printf(" saw %s\n", yyTokname(yychar))
			}
			fallthrough

		case 1, 2: /* incompletely recovered error ... try again */
			Errflag = 3

			/* find a state where "error" is a legal shift action */
			for yyp >= 0 {
				yyn = yyPact[yyS[yyp].yys] + yyErrCode
				if yyn >= 0 && yyn < yyLast {
					yystate = yyAct[yyn] /* simulate a shift of "error" */
					if yyChk[yystate] == yyErrCode {
						goto yystack
					}
				}

				/* the current p has no shift on "error", pop stack */
				if yyDebug >= 2 {
					__yyfmt__.Printf("error recovery pops state %d\n", yyS[yyp].yys)
				}
				yyp--
			}
			/* there is no state on the stack with an error shift ... abort */
			goto ret1

		case 3: /* no shift yet; clobber input char */
			if yyDebug >= 2 {
				__yyfmt__.Printf("error recovery discards %s\n", yyTokname(yychar))
			}
			if yychar == yyEofCode {
				goto ret1
			}
			yychar = -1
			goto yynewstate /* try again in the same state */
		}
	}

	/* reduction by production yyn */
	if yyDebug >= 2 {
		__yyfmt__.Printf("reduce %v in:\n\t%v\n", yyn, yyStatname(yystate))
	}

	yynt := yyn
	yypt := yyp
	_ = yypt // guard against "declared and not used"

	yyp -= yyR2[yyn]
	yyVAL = yyS[yyp+1]

	/* consult goto table to find next state */
	yyn = yyR1[yyn]
	yyg := yyPgo[yyn]
	yyj := yyg + yyS[yyp].yys + 1

	if yyj >= yyLast {
		yystate = yyAct[yyg]
	} else {
		yystate = yyAct[yyj]
		if yyChk[yystate] != -yyn {
			yystate = yyAct[yyg]
		}
	}
	// dummy call; replaced with literal code
	switch yynt {

	case 1:
		//line gosk.y:50
		{
			rootNode = &Node{
				typ:  nodeRoot,
				list: yyS[yypt-0].nodelist,
			}
		}
	case 2:
		//line gosk.y:57
		{
			if yyS[yypt-0].node != nil {
				yyS[yypt-1].nodelist.append(yyS[yypt-0].node)
				yyVAL.nodelist = yyS[yypt-1].nodelist
			}
		}
	case 3:
		//line gosk.y:63
		{
			yyVAL.nodelist = NewNodeList()
		}
	case 4:
		//line gosk.y:67
		{
			yyVAL.node = yyS[yypt-1].node
		}
	case 5:
		//line gosk.y:70
		{
			yyS[yypt-1].node.block = yyS[yypt-0].node
			yyVAL.node = yyS[yypt-1].node
		}
	case 6:
		yyVAL.node = yyS[yypt-0].node
	case 7:
		//line gosk.y:75
		{
			yyVAL.node = nil
		}
	case 8:
		yyVAL.node = yyS[yypt-0].node
	case 9:
		yyVAL.node = yyS[yypt-0].node
	case 10:
		//line gosk.y:83
		{
			yyVAL.node = &Node{
				typ:   nodeAssign,
				left:  yyS[yypt-2].node,
				right: yyS[yypt-0].node,
			}
		}
	case 11:
		//line gosk.y:91
		{
			yyVAL.node = &Node{
				typ:  nodeReturn,
				left: yyS[yypt-0].node,
			}
		}
	case 12:
		yyVAL.node = yyS[yypt-0].node
	case 13:
		yyVAL.node = yyS[yypt-0].node
	case 14:
		yyVAL.node = yyS[yypt-0].node
	case 15:
		//line gosk.y:103
		{
			yyVAL.node = &Node{
				typ:  nodeIf,
				cond: yyS[yypt-1].node,
				left: yyS[yypt-0].node,
			}
		}
	case 16:
		//line gosk.y:110
		{
			yyVAL.node = &Node{
				typ:   nodeIf,
				cond:  yyS[yypt-3].node,
				left:  yyS[yypt-2].node,
				right: yyS[yypt-0].node,
			}
		}
	case 17:
		//line gosk.y:119
		{
			yyVAL.node = &Node{
				typ:  nodeSwitch,
				cond: yyS[yypt-2].node,
				list: yyS[yypt-0].nodelist,
			}
		}
	case 18:
		//line gosk.y:126
		{
			cases := yyS[yypt-1].nodelist
			cases.append(yyS[yypt-0].node)

			yyVAL.node = &Node{
				typ:  nodeSwitch,
				cond: yyS[yypt-3].node,
				list: cases,
			}
		}
	case 19:
		//line gosk.y:137
		{
			yyS[yypt-1].nodelist.append(yyS[yypt-0].node)
			yyVAL.nodelist = yyS[yypt-1].nodelist
		}
	case 20:
		//line gosk.y:141
		{
			yyVAL.nodelist = NewNodeList()
		}
	case 21:
		//line gosk.y:145
		{
			yyVAL.node = &Node{
				typ:   nodeCase,
				cond:  yyS[yypt-2].node,
				block: yyS[yypt-0].node,
			}
		}
	case 22:
		//line gosk.y:153
		{
			yyVAL.node = &Node{
				typ:   nodeDefault,
				block: yyS[yypt-0].node,
			}
		}
	case 23:
		//line gosk.y:160
		{
			yyVAL.node = &Node{
				typ:   nodeAssignBlock,
				left:  yyS[yypt-2].node,
				block: yyS[yypt-0].node,
			}
		}
	case 24:
		//line gosk.y:168
		{
			yyVAL.node = &Node{
				typ:  nodeBlock,
				list: yyS[yypt-1].nodelist,
			}
		}
	case 25:
		//line gosk.y:175
		{
			yyVAL.node = yyS[yypt-1].node
		}
	case 26:
		yyVAL.node = yyS[yypt-0].node
	case 27:
		//line gosk.y:177
		{
			yyVAL.node = &Node{
				typ:   nodeOp,
				left:  yyS[yypt-2].node,
				right: yyS[yypt-0].node,
				val:   reflect.ValueOf(yyS[yypt-1].op),
			}
		}
	case 28:
		yyVAL.node = yyS[yypt-0].node
	case 29:
		//line gosk.y:187
		{
			yyVAL.node = &Node{
				typ:  nodeUnaryOp,
				left: yyS[yypt-0].node,
				val:  reflect.ValueOf(yyS[yypt-1].op),
			}
		}
	case 30:
		//line gosk.y:196
		{
			yyVAL.op = opOrOr
		}
	case 31:
		//line gosk.y:197
		{
			yyVAL.op = opAndAnd
		}
	case 32:
		yyVAL.op = yyS[yypt-0].op
	case 33:
		yyVAL.op = yyS[yypt-0].op
	case 34:
		yyVAL.op = yyS[yypt-0].op
	case 35:
		//line gosk.y:203
		{
			yyVAL.op = opEqual
		}
	case 36:
		//line gosk.y:204
		{
			yyVAL.op = opNonEqual
		}
	case 37:
		//line gosk.y:205
		{
			yyVAL.op = opLessThan
		}
	case 38:
		//line gosk.y:206
		{
			yyVAL.op = opLessThanOrEqual
		}
	case 39:
		//line gosk.y:207
		{
			yyVAL.op = opGreaterThan
		}
	case 40:
		//line gosk.y:208
		{
			yyVAL.op = opGreaterThanOrEqual
		}
	case 41:
		//line gosk.y:211
		{
			yyVAL.op = opPlus
		}
	case 42:
		//line gosk.y:212
		{
			yyVAL.op = opMinus
		}
	case 43:
		//line gosk.y:213
		{
			yyVAL.op = opOr
		}
	case 44:
		//line gosk.y:216
		{
			yyVAL.op = opMulti
		}
	case 45:
		//line gosk.y:217
		{
			yyVAL.op = opDivide
		}
	case 46:
		//line gosk.y:218
		{
			yyVAL.op = opMod
		}
	case 47:
		//line gosk.y:219
		{
			yyVAL.op = opLeftShift
		}
	case 48:
		//line gosk.y:220
		{
			yyVAL.op = opRightShift
		}
	case 49:
		//line gosk.y:221
		{
			yyVAL.op = opAnd
		}
	case 50:
		//line gosk.y:222
		{
			yyVAL.op = opAndNot
		}
	case 51:
		//line gosk.y:225
		{
			yyVAL.op = opNegate
		}
	case 52:
		//line gosk.y:226
		{
			yyVAL.op = opNot
		}
	case 53:
		//line gosk.y:227
		{
			yyVAL.op = opBitNot
		}
	case 54:
		yyVAL.node = yyS[yypt-0].node
	case 55:
		yyVAL.node = yyS[yypt-0].node
	case 56:
		yyVAL.node = yyS[yypt-0].node
	case 57:
		//line gosk.y:233
		{
			yyVAL.node = &Node{
				typ:  nodeSelector,
				left: yyS[yypt-1].node,
				val:  yyS[yypt-0].node.val,
			}
		}
	case 58:
		//line gosk.y:240
		{
			yyVAL.node = &Node{
				typ:   nodeIndex,
				left:  yyS[yypt-1].node,
				right: yyS[yypt-0].node,
			}
		}
	case 59:
		//line gosk.y:247
		{
			yyVAL.node = &Node{
				typ:  nodeSlice,
				left: yyS[yypt-1].node,
				list: yyS[yypt-0].nodelist,
			}
		}
	case 60:
		//line gosk.y:254
		{
			yyVAL.node = &Node{
				typ:  nodeCall,
				left: yyS[yypt-1].node,
				list: yyS[yypt-0].nodelist,
			}
		}
	case 61:
		yyVAL.node = yyS[yypt-0].node
	case 62:
		yyVAL.node = yyS[yypt-0].node
	case 63:
		//line gosk.y:266
		{
			yyVAL.node = yyS[yypt-1].node
		}
	case 64:
		//line gosk.y:268
		{
			yyVAL.nodelist = NewNodeList(nil, nil, nil)
		}
	case 65:
		//line gosk.y:269
		{
			yyVAL.nodelist = NewNodeList(yyS[yypt-2].node, nil, nil)
		}
	case 66:
		//line gosk.y:270
		{
			yyVAL.nodelist = NewNodeList(nil, yyS[yypt-1].node, nil)
		}
	case 67:
		//line gosk.y:271
		{
			yyVAL.nodelist = NewNodeList(yyS[yypt-3].node, yyS[yypt-1].node, nil)
		}
	case 68:
		//line gosk.y:272
		{
			yyVAL.nodelist = NewNodeList(nil, yyS[yypt-3].node, yyS[yypt-1].node)
		}
	case 69:
		//line gosk.y:273
		{
			yyVAL.nodelist = NewNodeList(yyS[yypt-5].node, yyS[yypt-3].node, yyS[yypt-1].node)
		}
	case 70:
		//line gosk.y:276
		{
			yyVAL.nodelist = nil
		}
	case 71:
		//line gosk.y:277
		{
			yyVAL.nodelist = yyS[yypt-1].nodelist
		}
	case 72:
		//line gosk.y:278
		{
			yyVAL.nodelist = yyS[yypt-2].nodelist
		}
	case 73:
		yyVAL.nodelist = yyS[yypt-0].nodelist
	case 74:
		//line gosk.y:282
		{
			yyS[yypt-1].nodelist.applyable()
			yyVAL.nodelist = yyS[yypt-1].nodelist
		}
	case 75:
		//line gosk.y:287
		{
			yyS[yypt-0].nodelist.prepend(yyS[yypt-1].node)
			yyVAL.nodelist = yyS[yypt-0].nodelist
		}
	case 76:
		//line gosk.y:292
		{
			yyS[yypt-2].nodelist.append(yyS[yypt-0].node)
			yyVAL.nodelist = yyS[yypt-2].nodelist
		}
	case 77:
		//line gosk.y:296
		{
			yyVAL.nodelist = NewNodeList()
		}
	case 78:
		yyVAL.node = yyS[yypt-0].node
	case 79:
		yyVAL.node = yyS[yypt-0].node
	case 80:
		yyVAL.node = yyS[yypt-0].node
	case 81:
		yyVAL.node = yyS[yypt-0].node
	case 82:
		yyVAL.node = yyS[yypt-0].node
	case 83:
		yyVAL.node = yyS[yypt-0].node
	case 84:
		yyVAL.node = yyS[yypt-0].node
	case 85:
		yyVAL.node = yyS[yypt-0].node
	}
	goto yystack /* stack new state and value */
}
