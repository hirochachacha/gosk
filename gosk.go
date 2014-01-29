package main

import "reflect"

type VarMap map[string]interface{}
type varMap map[string]reflect.Value
type dVarMap map[string]*Dataflow

var (
	typeOfString    = reflect.TypeOf("")
	typeOfFloat64   = reflect.TypeOf(float64(0))
	typeOfComplex64 = reflect.TypeOf(complex64(0))

	valueOfNil = reflect.ValueOf(nil)
)

var rootNode *Node
