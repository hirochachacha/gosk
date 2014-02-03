package gosk

import (
	"reflect"
)

type Template struct {
	ast     *Node
	globals varMap
}

func NewTemplate(format string) *Template {
	lexer := NewLexer("gosk", format)
	yyParse(lexer)
	return &Template{
		ast:     rootNode,
		globals: make(varMap),
	}
}

func (t *Template) Vars(vars VarMap) *Template {
	for k, v := range vars {
		t.globals[k] = reflect.ValueOf(v)
	}
	return t
}

func (t *Template) Execute(ctx ...interface{}) string {
	newFrame := NewFrame(nil)
	c := make(chan reflect.Value)

	lctx := valueOfNil
	if len(ctx) > 0 {
		lctx = reflect.ValueOf(ctx[0])
		kind := lctx.Kind()
		for kind == reflect.Ptr {
			lctx = lctx.Elem()
			kind = lctx.Kind()
		}
		if !(kind == reflect.Struct || kind == reflect.Map) {
			panic("can't use non struct or non map as context")
		}
		if kind == reflect.Map && !typeOfString.AssignableTo(lctx.Type().Key()) {
			panic("can't use non string key map as context")
		}
	}
	go newFrame.walkRoot(t.globals, lctx, t.ast, c)
	return (<-c).String()
}
