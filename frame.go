package main

import (
	"reflect"
	"sync"
)

type Frame struct {
	l       sync.RWMutex
	dlocals dVarMap
	parent  *Frame
}

func NewFrame(parent *Frame) *Frame {
	return &Frame{
		dlocals: make(dVarMap),
		parent:  parent,
	}
}

func (frame *Frame) getDVar(globals varMap, id string) *Dataflow {
	var d *Dataflow
	f := frame
	for f != nil {
		f.l.RLock()
		d, ok := f.dlocals[id]
		f.l.RUnlock()
		if ok {
			return d
		}
		f = f.parent
	}
	if _, ok := globals[id]; ok {
		return nil
	}
	frame.l.Lock()
	if d, ok := frame.dlocals[id]; !ok {
		d = NewDataflow()
		frame.dlocals[id] = d
	}
	frame.l.Unlock()
	return d
}

func (frame *Frame) getVar(globals varMap, id string) reflect.Value {
	d := frame.getDVar(globals, id)
	if d == nil {
		return globals[id]
	}
	return d.get()
}

func (frame *Frame) walkRoot(globals varMap, lctx reflect.Value, node *Node, c chan<- reflect.Value) {
	frame.walkStmts(globals, lctx, node, c)
}

func (frame *Frame) walkStmts(globals varMap, lctx reflect.Value, node *Node, c chan<- reflect.Value) {
	node.list.traverse(func(node *Node) {
		if node == nil {
			panic("illegal stmts")
		}
		frame.walkStmt(globals, lctx, node, c)
	})
}

func (frame *Frame) walkStmt(globals varMap, lctx reflect.Value, node *Node, c chan<- reflect.Value) {
	switch node.typ {
	case nodeReturn:
		frame.walkReturn(globals, lctx, node, c)
	case nodeAssign:
		frame.walkAssign(globals, lctx, node, c)
	case nodeIf:
		frame.walkIf(globals, lctx, node, c)
	case nodeSwitch:
		frame.walkSwitch(globals, lctx, node, c)
	case nodeAssignBlock:
		frame.walkAssignBlock(globals, lctx, node, c)
	default:
		panic("illegal stmt")
	}
}

func (frame *Frame) walkBlock(globals varMap, lctx reflect.Value, node *Node, c chan<- reflect.Value) {
	newFrame := NewFrame(frame)
	newFrame.walkStmts(globals, lctx, node, c)
}

func (frame *Frame) walkReturn(globals varMap, lctx reflect.Value, node *Node, c chan<- reflect.Value) {
	frame.walkExpr(globals, lctx, node.left, c)
	if node.block != nil {
		frame.walkBlock(globals, lctx, node.block, c)
	}
}

func (frame *Frame) walkAssign(globals varMap, lctx reflect.Value, node *Node, c chan<- reflect.Value) {
	id := node.left.val.String()
	d := frame.getDVar(globals, id)
	newC := make(chan reflect.Value)
	go frame.walkExpr(globals, lctx, node.right, newC)
	err := d.set(<-newC)
	if err != nil {
		panic(err)
	}
	if node.block != nil {
		frame.walkBlock(globals, lctx, node.block, c)
	}
}

func (frame *Frame) walkIf(globals varMap, lctx reflect.Value, node *Node, c chan<- reflect.Value) {
	newC := make(chan reflect.Value)
	go frame.walkExpr(globals, lctx, node.cond, newC)
	cond := <-newC
	if cond.Kind() != reflect.Bool {
		panic("can't use non bool as if condition")
	}
	if cond.Bool() {
		frame.walkBlock(globals, lctx, node.left, c)
	} else {
		if node.right != nil {
			frame.walkBlock(globals, lctx, node.right, c)
		}
	}
}

func (frame *Frame) walkSwitch(globals varMap, lctx reflect.Value, node *Node, c chan<- reflect.Value) {
	newC := make(chan reflect.Value)
	go frame.walkExpr(globals, lctx, node.cond, newC)
	lcond := <-newC
	node.list.traverse(func(node *Node) {
		switch node.typ {
		case nodeCase:
			go frame.walkExpr(globals, lctx, node.cond, newC)
			rcond := <-newC
			if lcond.Interface() == rcond.Interface() {
				frame.walkBlock(globals, lctx, node.block, c)
			}
		case nodeDefault:
			frame.walkBlock(globals, lctx, node.block, c)
		default:
			panic("illegal switch")
		}
	})
}

func (frame *Frame) walkAssignBlock(globals varMap, lctx reflect.Value, node *Node, c chan<- reflect.Value) {
	id := node.left.val.String()
	d := frame.getDVar(globals, id)
	newC := make(chan reflect.Value)
	go frame.walkBlock(globals, lctx, node.block, newC)
	err := d.set(<-newC)
	if err != nil {
		panic(err)
	}
}

func (frame *Frame) walkExpr(globals varMap, lctx reflect.Value, node *Node, c chan<- reflect.Value) {
	switch node.typ {
	case nodeUnaryOp:
		frame.walkUnaryOp(globals, lctx, node, c)
	case nodeOp:
		frame.walkOp(globals, lctx, node, c)
	case nodeLiteral:
		frame.walkLiteral(globals, lctx, node, c)
	case nodeIdentifier:
		frame.walkIdentifier(globals, lctx, node, c)
	case nodeField:
		frame.walkField(globals, lctx, node, c)
	case nodeSelector:
		frame.walkSelector(globals, lctx, node, c)
	case nodeIndex:
		frame.walkIndex(globals, lctx, node, c)
	case nodeSlice:
		frame.walkSlice(globals, lctx, node, c)
	case nodeCall:
		frame.walkCall(globals, lctx, node, c)
	default:
		panic("illegal expr")
	}
}

func (frame *Frame) walkUnaryOp(globals varMap, lctx reflect.Value, node *Node, c chan<- reflect.Value) {
	newC := make(chan reflect.Value)
	go frame.walkExpr(globals, lctx, node.left, newC)
	expr := <-newC
	kind := expr.Kind()
	switch node.val.Interface().(opType) {
	case opNegate:
		if kind != reflect.Int64 && kind != reflect.Float64 && kind != reflect.Complex64 {
			panic("can't use non number with '-'")
		}
		if kind == reflect.Int64 {
			c <- reflect.ValueOf(-1 * expr.Int())
		}
		if kind == reflect.Float64 {
			c <- reflect.ValueOf(-1 * expr.Float())
		}
		c <- reflect.ValueOf(-1 * expr.Complex())
	case opNot:
		if kind != reflect.Bool {
			panic("can't use non bool with '!'")
		}
		c <- reflect.ValueOf(!expr.Bool())
	case opBitNot:
		if kind != reflect.Int64 {
			panic("can't use non int with '^'")
		}
		c <- reflect.ValueOf(int64(^expr.Int()))
	default:
		panic("illegal unary_op")
	}
}

func (frame *Frame) walkOp(globals varMap, lctx reflect.Value, node *Node, c chan<- reflect.Value) {
	newC := make(chan reflect.Value)
	newC2 := make(chan reflect.Value)

	go frame.walkExpr(globals, lctx, node.left, newC)
	go frame.walkExpr(globals, lctx, node.right, newC2)

	lexpr := <-newC
	rexpr := <-newC2

	lkind := lexpr.Kind()
	rkind := rexpr.Kind()

	switch node.val.Interface().(opType) {
	case opOrOr:
		if !(lkind == reflect.Bool && rkind == reflect.Bool) {
			panic("can't use non bool with '||'")
		}
		c <- reflect.ValueOf(lexpr.Bool() || rexpr.Bool())
	case opAndAnd:
		if !(lkind == reflect.Bool && rkind == reflect.Bool) {
			panic("can't use non bool with '&&'")
		}
		c <- reflect.ValueOf(lexpr.Bool() && rexpr.Bool())
	case opEqual:
		c <- reflect.ValueOf(lexpr.Interface() == rexpr.Interface())
	case opNonEqual:
		c <- reflect.ValueOf(lexpr.Interface() != rexpr.Interface())
	case opLessThan:
		if lkind != reflect.Int64 && lkind != reflect.Float64 {
			panic("can't use non number with '<'")
		}
		if rkind != reflect.Int64 && rkind != reflect.Float64 {
			panic("can't use non number with '<'")
		}
		c <- reflect.ValueOf(lexpr.Convert(typeOfFloat64).Float() < rexpr.Convert(typeOfFloat64).Float())
	case opLessThanOrEqual:
		if lkind != reflect.Int64 && lkind != reflect.Float64 {
			panic("can't use non number with '<='")
		}
		if rkind != reflect.Int64 && rkind != reflect.Float64 {
			panic("can't use non number with '<='")
		}
		c <- reflect.ValueOf(lexpr.Convert(typeOfFloat64).Float() <= rexpr.Convert(typeOfFloat64).Float())
	case opGreaterThan:
		if lkind != reflect.Int64 && lkind != reflect.Float64 {
			panic("can't use non number with '>'")
		}
		if rkind != reflect.Int64 && rkind != reflect.Float64 {
			panic("can't use non number with '>'")
		}
		c <- reflect.ValueOf(lexpr.Convert(typeOfFloat64).Float() > rexpr.Convert(typeOfFloat64).Float())
	case opGreaterThanOrEqual:
		if lkind != reflect.Int64 && lkind != reflect.Float64 {
			panic("can't use non number with '>='")
		}
		if rkind != reflect.Int64 && rkind != reflect.Float64 {
			panic("can't use non number with '>='")
		}
		c <- reflect.ValueOf(lexpr.Convert(typeOfFloat64).Float() >= rexpr.Convert(typeOfFloat64).Float())
	case opPlus:
		if lkind != reflect.Int64 && lkind != reflect.Float64 && lkind != reflect.Complex64 {
			panic("can't use non number with '+'")
		}
		if rkind != reflect.Int64 && rkind != reflect.Float64 && rkind != reflect.Complex64 {
			panic("can't use non number with '+'")
		}
		if lkind == reflect.Complex64 || rkind == reflect.Complex64 {
			c <- reflect.ValueOf(lexpr.Convert(typeOfComplex64).Complex() + rexpr.Convert(typeOfComplex64).Complex())
		}
		if lkind == reflect.Int64 && rkind == reflect.Int64 {
			c <- reflect.ValueOf(lexpr.Int() + rexpr.Int())
		}
		c <- reflect.ValueOf(lexpr.Convert(typeOfFloat64).Float() + rexpr.Convert(typeOfFloat64).Float())
	case opMinus:
		if lkind != reflect.Int64 && lkind != reflect.Float64 && lkind != reflect.Complex64 {
			panic("can't use non number with '-'")
		}
		if rkind != reflect.Int64 && rkind != reflect.Float64 && rkind != reflect.Complex64 {
			panic("can't use non number with '-'")
		}
		if lkind == reflect.Complex64 || rkind == reflect.Complex64 {
			c <- reflect.ValueOf(lexpr.Convert(typeOfComplex64).Complex() - rexpr.Convert(typeOfComplex64).Complex())
		}
		if lkind == reflect.Int64 && rkind == reflect.Int64 {
			c <- reflect.ValueOf(lexpr.Int() - rexpr.Int())
		}
		c <- reflect.ValueOf(lexpr.Convert(typeOfFloat64).Float() - rexpr.Convert(typeOfFloat64).Float())
	case opOr:
		if !(lkind == reflect.Int64 && rkind == reflect.Int64) {
			panic("can't use non int with '|'")
		}
		c <- reflect.ValueOf(lexpr.Int() | rexpr.Int())
	case opMulti:
		if lkind != reflect.Int64 && lkind != reflect.Float64 && lkind != reflect.Complex64 {
			panic("can't use non number with '*'")
		}
		if rkind != reflect.Int64 && rkind != reflect.Float64 && rkind != reflect.Complex64 {
			panic("can't use non number with '*'")
		}
		if lkind == reflect.Complex64 || rkind == reflect.Complex64 {
			c <- reflect.ValueOf(lexpr.Convert(typeOfComplex64).Complex() * rexpr.Convert(typeOfComplex64).Complex())
		}
		if lkind == reflect.Int64 && rkind == reflect.Int64 {
			c <- reflect.ValueOf(lexpr.Int() * rexpr.Int())
		}
		c <- reflect.ValueOf(lexpr.Convert(typeOfFloat64).Float() * rexpr.Convert(typeOfFloat64).Float())
	case opDivide:
		if lkind != reflect.Int64 && lkind != reflect.Float64 && lkind != reflect.Complex64 {
			panic("can't use non number with '/'")
		}
		if rkind != reflect.Int64 && rkind != reflect.Float64 && rkind != reflect.Complex64 {
			panic("can't use non number with '/'")
		}
		if lkind == reflect.Complex64 || rkind == reflect.Complex64 {
			c <- reflect.ValueOf(lexpr.Convert(typeOfComplex64).Complex() / rexpr.Convert(typeOfComplex64).Complex())
		}
		if lkind == reflect.Int64 && rkind == reflect.Int64 {
			c <- reflect.ValueOf(lexpr.Int() / rexpr.Int())
		}
		c <- reflect.ValueOf(lexpr.Convert(typeOfFloat64).Float() / rexpr.Convert(typeOfFloat64).Float())
	case opMod:
		if !(lkind == reflect.Int64 && rkind == reflect.Int64) {
			panic("can't use non int with '%'")
		}
		c <- reflect.ValueOf(lexpr.Int() % rexpr.Int())
	case opAnd:
		if !(lkind == reflect.Int64 && rkind == reflect.Int64) {
			panic("can't use non int with '%'")
		}
		c <- reflect.ValueOf(lexpr.Int() & rexpr.Int())
	case opLeftShift:
		if !(lkind == reflect.Int64 && rkind == reflect.Int64) {
			panic("can't use non int with '<<'")
		}
		l := lexpr.Int()
		r := rexpr.Int()
		if l < 0 || r < 0 {
			panic("can't use negative int with '<<'")
		}
		c <- reflect.ValueOf(uint(l) << uint(r))
	case opRightShift:
		if !(lkind == reflect.Int64 && rkind == reflect.Int64) {
			panic("can't use non int with '>>'")
		}
		l := lexpr.Int()
		r := rexpr.Int()
		if l < 0 || r < 0 {
			panic("can't use negative int with '>>'")
		}
		c <- reflect.ValueOf(uint(l) >> uint(r))
	case opAndNot:
		if !(lkind == reflect.Int64 && rkind == reflect.Int64) {
			panic("can't use non int with '&^'")
		}
		c <- reflect.ValueOf(lexpr.Int() &^ rexpr.Int())
	default:
		panic("illegal op")
	}
}

func (frame *Frame) walkLiteral(globals varMap, lctx reflect.Value, node *Node, c chan<- reflect.Value) {
	c <- node.val
}

func (frame *Frame) walkIdentifier(globals varMap, lctx reflect.Value, node *Node, c chan<- reflect.Value) {
	id := node.val.String()
	c <- frame.getVar(globals, id)
}

func (frame *Frame) walkField(globals varMap, lctx reflect.Value, node *Node, c chan<- reflect.Value) {
	id := node.val.String()[1:]

	switch lctx.Kind() {
	case reflect.Struct:
		c <- lctx.FieldByName(id)
	case reflect.Map:
		c <- lctx.MapIndex(reflect.ValueOf(id))
	default:
		panic("illegal field")
	}
}

func (frame *Frame) walkSelector(globals varMap, lctx reflect.Value, node *Node, c chan<- reflect.Value) {
	newC := make(chan reflect.Value)
	go frame.walkExpr(globals, lctx, node.left, newC)
	container := <-newC
	if container.Kind() != reflect.Struct {
		panic("can't use non struct with selector")
	}
	selector := node.val.String()
	c <- container.FieldByName(selector)
}

func (frame *Frame) walkIndex(globals varMap, lctx reflect.Value, node *Node, c chan<- reflect.Value) {
	newC := make(chan reflect.Value)
	go frame.walkExpr(globals, lctx, node.left, newC)
	container := <-newC
	kind := container.Kind()
	if kind != reflect.Array && kind != reflect.Slice {
		panic("can't use non array or non slice with index")
	}
	go frame.walkExpr(globals, lctx, node.right, newC)
	index := int((<-newC).Int())
	c <- container.Index(index)
}

func (frame *Frame) walkSlice(globals varMap, lctx reflect.Value, node *Node, c chan<- reflect.Value) {
	newC := make(chan reflect.Value)
	go frame.walkExpr(globals, lctx, node.left, newC)
	container := <-newC
	kind := container.Kind()
	if kind != reflect.Array && kind != reflect.Slice {
		panic("can't use non array or non slice with slice")
	}

	args := frame.walkSliceArgs(globals, lctx, node, nil)

	var i, j, k int
	i = int(args[0].Int())
	if args[1] == valueOfNil {
		j = container.Len()
	} else {
		j = int(args[1].Int())
	}
	if args[2] == valueOfNil {
		c <- container.Slice(i, j)
	}
	k = int(args[2].Int())
	c <- container.Slice3(i, j, k)
}

func (frame *Frame) walkSliceArgs(globals varMap, lctx reflect.Value, node *Node, _ chan<- reflect.Value) []reflect.Value {
	cs := make([]chan reflect.Value, 0, node.list.length)
	node.list.traverse(func(node *Node) {
		newC := make(chan reflect.Value)
		if node == nil {
			newC <- valueOfNil
		} else {
			go frame.walkExpr(globals, lctx, node, newC)
		}
		cs = append(cs, newC)
	})

	args := make([]reflect.Value, 0, node.list.length)
	for _, c := range cs {
		args = append(args, <-c)
	}
	return args

}

func (frame *Frame) walkCall(globals varMap, lctx reflect.Value, node *Node, c chan<- reflect.Value) {
	newC := make(chan reflect.Value)
	go frame.walkExpr(globals, lctx, node.left, newC)
	container := <-newC

	args := frame.walkCallArgs(globals, lctx, node, nil)

	c <- container.Call(args)[0]
}

func (frame *Frame) walkCallArgs(globals varMap, lctx reflect.Value, node *Node, _ chan<- reflect.Value) []reflect.Value {
	cs := make([]chan reflect.Value, 0, node.list.length)
	node.list.traverse(func(node *Node) {
		newC := make(chan reflect.Value)
		if node == nil {
			panic("illegal call")
		}
		go frame.walkExpr(globals, lctx, node, newC)
		cs = append(cs, newC)
	})

	args := make([]reflect.Value, 0, node.list.length)
	for _, c := range cs {
		args = append(args, <-c)
	}
	return args
}
