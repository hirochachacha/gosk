package main

import (
	"reflect"
	"sync/atomic"
)

var isDone int32

type Template struct {
	globals varMap
}

func NewTemplate(format string) *Template {
	if !atomic.CompareAndSwapInt32(&isDone, 0, 1) {
		panic("currently, only single instance is supported")
	}
	lexer := NewLexer("gosk", format)
	yyParse(lexer)
	return &Template{
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
		if !(kind == reflect.Struct || kind == reflect.Map) {
			panic("can't use non struct or non map as context")
		}
		if kind == reflect.Map && !typeOfString.AssignableTo(lctx.Type().Key()) {
			panic("can't use non string key map as context")
		}
	}
	go newFrame.walkRoot(t.globals, lctx, rootNode, c)
	return (<-c).String()
}
