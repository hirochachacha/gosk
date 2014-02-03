package gosk

import (
	"fmt"
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
	if val, ok := globals[id]; ok {
		return NewDataflowWith(val)
	}
	frame.l.Lock()
	d, ok := frame.dlocals[id]
	if !ok {
		d = NewDataflow()
		frame.dlocals[id] = d
	}
	frame.l.Unlock()
	return d
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
	expr := frame.levalExpr(globals, lctx, node.left)
	if node.block != nil {
		go frame.walkBlock(globals, lctx, node.block, c)
	}
	c <- expr.Get()
}

func (frame *Frame) walkAssign(globals varMap, lctx reflect.Value, node *Node, c chan<- reflect.Value) {
	id := node.left.val.String()
	d := frame.getDVar(globals, id)
	expr := frame.levalExpr(globals, lctx, node.right)
	if node.block != nil {
		go frame.walkBlock(globals, lctx, node.block, c)
	}
	err := d.Set(expr.Get())
	if err != nil {
		panic(err)
	}
}

func (frame *Frame) walkIf(globals varMap, lctx reflect.Value, node *Node, c chan<- reflect.Value) {
	cond := frame.levalExpr(globals, lctx, node.cond).Get()
	if cond.Kind() != reflect.Bool {
		panic(fmt.Sprintf("can't use non bool as if condition: %s", cond.Kind()))
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
	lcond := frame.levalExpr(globals, lctx, node.cond).Get()
	node.list.traverse(func(node *Node) {
		switch node.typ {
		case nodeCase:
			left := lcond.Interface()
			node.list.traverse(func(child *Node) {
				if node == nil {
					panic("illegal case")
				}
				rcond := frame.levalExpr(globals, lctx, child).Get()
				if left == rcond.Interface() {
					frame.walkBlock(globals, lctx, node.block, c)
				}
			})
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
	err := d.Set(<-newC)
	if err != nil {
		panic(err)
	}
}

func (frame *Frame) levalExpr(globals varMap, lctx reflect.Value, node *Node) *Dataflow {
	switch node.typ {
	case nodeUnaryOp:
		return frame.levalUnaryOp(globals, lctx, node)
	case nodeOp:
		return frame.levalOp(globals, lctx, node)
	case nodeLiteral:
		return frame.levalLiteral(globals, lctx, node)
	case nodeGlobalIdentifier:
		return frame.levalGlobalIdentifier(globals, lctx, node)
	case nodeIdentifier:
		return frame.levalIdentifier(globals, lctx, node)
	case nodeField:
		return frame.levalField(globals, lctx, node)
	case nodeSelector:
		return frame.levalSelector(globals, lctx, node)
	case nodeIndex:
		return frame.levalIndex(globals, lctx, node)
	case nodeSlice:
		return frame.levalSlice(globals, lctx, node)
	case nodeCall:
		return frame.levalCall(globals, lctx, node)
	default:
		panic("illegal expr")
	}
	return nil
}

func (frame *Frame) levalUnaryOp(globals varMap, lctx reflect.Value, node *Node) *Dataflow {
	d := NewDataflow()
	ld := frame.levalExpr(globals, lctx, node.left)
	go func() {
		expr := ld.Get()
		kind := expr.Kind()
		switch node.val.Interface().(opType) {
		case opNegate:
			if kind != reflect.Int64 && kind != reflect.Float64 && kind != reflect.Complex64 {
				panic("can't use non number with '-'")
			}
			if kind == reflect.Int64 {
				d.Set(reflect.ValueOf(-1 * expr.Int()))
			} else if kind == reflect.Float64 {
				d.Set(reflect.ValueOf(-1 * expr.Float()))
			} else {
				d.Set(reflect.ValueOf(-1 * expr.Complex()))
			}
		case opNot:
			if kind != reflect.Bool {
				panic("can't use non bool with '!'")
			}
			d.Set(reflect.ValueOf(!expr.Bool()))
		case opBitNot:
			if kind != reflect.Int64 {
				panic("can't use non int with '^'")
			}
			d.Set(reflect.ValueOf(int64(^expr.Int())))
		default:
			panic("illegal unary_op")
		}
	}()
	return d
}

func (frame *Frame) levalOp(globals varMap, lctx reflect.Value, node *Node) *Dataflow {
	d := NewDataflow()
	ld := frame.levalExpr(globals, lctx, node.left)
	rd := frame.levalExpr(globals, lctx, node.right)
	go func() {
		lexpr := ld.Get()
		rexpr := rd.Get()
		lkind := lexpr.Kind()
		rkind := rexpr.Kind()

		switch node.val.Interface().(opType) {
		case opOrOr:
			if !(lkind == reflect.Bool && rkind == reflect.Bool) {
				panic("can't use non bool with '||'")
			}
			d.Set(reflect.ValueOf(lexpr.Bool() || rexpr.Bool()))
		case opAndAnd:
			if !(lkind == reflect.Bool && rkind == reflect.Bool) {
				panic("can't use non bool with '&&'")
			}
			d.Set(reflect.ValueOf(lexpr.Bool() && rexpr.Bool()))
		case opEqual:
			d.Set(reflect.ValueOf(lexpr.Interface() == rexpr.Interface()))
		case opNonEqual:
			d.Set(reflect.ValueOf(lexpr.Interface() != rexpr.Interface()))
		case opLessThan:
			if lkind != reflect.Int64 && lkind != reflect.Float64 {
				panic("can't use non number with '<'")
			}
			if rkind != reflect.Int64 && rkind != reflect.Float64 {
				panic("can't use non number with '<'")
			}
			d.Set(reflect.ValueOf(lexpr.Convert(typeOfFloat64).Float() < rexpr.Convert(typeOfFloat64).Float()))
		case opLessThanOrEqual:
			if lkind != reflect.Int64 && lkind != reflect.Float64 {
				panic("can't use non number with '<='")
			}
			if rkind != reflect.Int64 && rkind != reflect.Float64 {
				panic("can't use non number with '<='")
			}
			d.Set(reflect.ValueOf(lexpr.Convert(typeOfFloat64).Float() <= rexpr.Convert(typeOfFloat64).Float()))
		case opGreaterThan:
			if lkind != reflect.Int64 && lkind != reflect.Float64 {
				panic("can't use non number with '>'")
			}
			if rkind != reflect.Int64 && rkind != reflect.Float64 {
				panic("can't use non number with '>'")
			}
			d.Set(reflect.ValueOf(lexpr.Convert(typeOfFloat64).Float() > rexpr.Convert(typeOfFloat64).Float()))
		case opGreaterThanOrEqual:
			if lkind != reflect.Int64 && lkind != reflect.Float64 {
				panic("can't use non number with '>='")
			}
			if rkind != reflect.Int64 && rkind != reflect.Float64 {
				panic("can't use non number with '>='")
			}
			d.Set(reflect.ValueOf(lexpr.Convert(typeOfFloat64).Float() >= rexpr.Convert(typeOfFloat64).Float()))
		case opPlus:
			if lkind == reflect.String && rkind == reflect.String {
				d.Set(reflect.ValueOf(lexpr.String() + rexpr.String()))
			} else {
				if lkind != reflect.Int64 && lkind != reflect.Float64 && lkind != reflect.Complex64 {
					panic(fmt.Sprintf("can't use non number with '+': %s", lkind))
				}
				if rkind != reflect.Int64 && rkind != reflect.Float64 && rkind != reflect.Complex64 {
					panic(fmt.Sprintf("can't use non number with '+': %s", rkind))
				}
				if lkind == reflect.Complex64 || rkind == reflect.Complex64 {
					d.Set(reflect.ValueOf(lexpr.Convert(typeOfComplex64).Complex() + rexpr.Convert(typeOfComplex64).Complex()))
				} else if lkind == reflect.Int64 && rkind == reflect.Int64 {
					d.Set(reflect.ValueOf(lexpr.Int() + rexpr.Int()))
				} else {
					d.Set(reflect.ValueOf(lexpr.Convert(typeOfFloat64).Float() + rexpr.Convert(typeOfFloat64).Float()))
				}
			}
		case opMinus:
			if lkind != reflect.Int64 && lkind != reflect.Float64 && lkind != reflect.Complex64 {
				panic("can't use non number with '-'")
			}
			if rkind != reflect.Int64 && rkind != reflect.Float64 && rkind != reflect.Complex64 {
				panic("can't use non number with '-'")
			}
			if lkind == reflect.Complex64 || rkind == reflect.Complex64 {
				d.Set(reflect.ValueOf(lexpr.Convert(typeOfComplex64).Complex() - rexpr.Convert(typeOfComplex64).Complex()))
			} else if lkind == reflect.Int64 && rkind == reflect.Int64 {
				d.Set(reflect.ValueOf(lexpr.Int() - rexpr.Int()))
			} else {
				d.Set(reflect.ValueOf(lexpr.Convert(typeOfFloat64).Float() - rexpr.Convert(typeOfFloat64).Float()))
			}
		case opOr:
			if !(lkind == reflect.Int64 && rkind == reflect.Int64) {
				panic("can't use non int with '|'")
			}
			d.Set(reflect.ValueOf(lexpr.Int() | rexpr.Int()))
		case opMulti:
			if lkind != reflect.Int64 && lkind != reflect.Float64 && lkind != reflect.Complex64 {
				panic("can't use non number with '*'")
			}
			if rkind != reflect.Int64 && rkind != reflect.Float64 && rkind != reflect.Complex64 {
				panic("can't use non number with '*'")
			}
			if lkind == reflect.Complex64 || rkind == reflect.Complex64 {
				d.Set(reflect.ValueOf(lexpr.Convert(typeOfComplex64).Complex() * rexpr.Convert(typeOfComplex64).Complex()))
			} else if lkind == reflect.Int64 && rkind == reflect.Int64 {
				d.Set(reflect.ValueOf(lexpr.Int() * rexpr.Int()))
			} else {
				d.Set(reflect.ValueOf(lexpr.Convert(typeOfFloat64).Float() * rexpr.Convert(typeOfFloat64).Float()))
			}
		case opDivide:
			if lkind != reflect.Int64 && lkind != reflect.Float64 && lkind != reflect.Complex64 {
				panic("can't use non number with '/'")
			}
			if rkind != reflect.Int64 && rkind != reflect.Float64 && rkind != reflect.Complex64 {
				panic("can't use non number with '/'")
			}
			if lkind == reflect.Complex64 || rkind == reflect.Complex64 {
				d.Set(reflect.ValueOf(lexpr.Convert(typeOfComplex64).Complex() / rexpr.Convert(typeOfComplex64).Complex()))
			} else if lkind == reflect.Int64 && rkind == reflect.Int64 {
				d.Set(reflect.ValueOf(lexpr.Int() / rexpr.Int()))
			} else {
				d.Set(reflect.ValueOf(lexpr.Convert(typeOfFloat64).Float() / rexpr.Convert(typeOfFloat64).Float()))
			}
		case opMod:
			if !(lkind == reflect.Int64 && rkind == reflect.Int64) {
				panic("can't use non int with '%'")
			}
			d.Set(reflect.ValueOf(lexpr.Int() % rexpr.Int()))
		case opLeftShift:
			if !(lkind == reflect.Int64 && rkind == reflect.Int64) {
				panic("can't use non int with '<<'")
			}
			l := lexpr.Int()
			r := rexpr.Int()
			if l < 0 || r < 0 {
				panic("can't use negative int with '<<'")
			}
			d.Set(reflect.ValueOf(uint(l) << uint(r)))
		case opRightShift:
			if !(lkind == reflect.Int64 && rkind == reflect.Int64) {
				panic("can't use non int with '>>'")
			}
			l := lexpr.Int()
			r := rexpr.Int()
			if l < 0 || r < 0 {
				panic("can't use negative int with '>>'")
			}
			d.Set(reflect.ValueOf(uint(l) >> uint(r)))
		case opAnd:
			if !(lkind == reflect.Int64 && rkind == reflect.Int64) {
				panic("can't use non int with '%'")
			}
			d.Set(reflect.ValueOf(lexpr.Int() & rexpr.Int()))
		case opAndNot:
			if !(lkind == reflect.Int64 && rkind == reflect.Int64) {
				panic("can't use non int with '&^'")
			}
			d.Set(reflect.ValueOf(lexpr.Int() &^ rexpr.Int()))
		default:
			panic("illegal op")
		}
	}()
	return d
}

func (frame *Frame) levalLiteral(globals varMap, lctx reflect.Value, node *Node) *Dataflow {
	return NewDataflowWith(node.val)
}

func (frame *Frame) levalGlobalIdentifier(globals varMap, lctx reflect.Value, node *Node) *Dataflow {
	id := node.val.String()[1:]
	val, ok := globals[id]
	if !ok {
		panic(fmt.Sprintf("unknown globals: %s", id))
	}
	return NewDataflowWith(val)
}

func (frame *Frame) levalIdentifier(globals varMap, lctx reflect.Value, node *Node) *Dataflow {
	id := node.val.String()
	return frame.getDVar(globals, id)
}

func (frame *Frame) levalField(globals varMap, lctx reflect.Value, node *Node) *Dataflow {
	id := node.val.String()[1:]

	switch lctx.Kind() {
	case reflect.Struct:
		return NewDataflowWith(lctx.FieldByName(id))
	case reflect.Map:
		val := lctx.MapIndex(reflect.ValueOf(id))
		if !val.IsValid() {
			panic(fmt.Sprintf("unknown field: %s", id))
		}
		return NewDataflowWith(val.Elem())
	default:
		panic(fmt.Sprintf("illegal field: %s", id))
	}
	return nil
}

func (frame *Frame) levalSelector(globals varMap, lctx reflect.Value, node *Node) *Dataflow {
	d := NewDataflow()
	containerd := frame.levalExpr(globals, lctx, node.left)
	go func() {
		container := containerd.Get()
		if container.Kind() != reflect.Struct {
			panic("can't use non struct with selector")
		}
		selector := node.val.String()
		d.Set(container.FieldByName(selector))
	}()
	return d
}

func (frame *Frame) levalIndex(globals varMap, lctx reflect.Value, node *Node) *Dataflow {
	d := NewDataflow()
	containerd := frame.levalExpr(globals, lctx, node.left)
	indexd := frame.levalExpr(globals, lctx, node.right)
	go func() {
		container := containerd.Get()
		kind := container.Kind()
		if kind != reflect.Array && kind != reflect.Slice {
			panic("can't use non array or non slice with index")
		}
		index := int(indexd.Get().Int())
		d.Set(container.Index(index))
	}()
	return d
}

func (frame *Frame) levalSlice(globals varMap, lctx reflect.Value, node *Node) *Dataflow {
	d := NewDataflow()
	containerd := frame.levalExpr(globals, lctx, node.left)
	ds := frame.levalSliceArgs(globals, lctx, node)
	go func() {
		container := containerd.Get()
		kind := container.Kind()
		if kind != reflect.Array && kind != reflect.Slice && kind != reflect.String {
			panic("can't use non array or non slice or non string with slice")
		}

		args := make([]reflect.Value, 0, node.list.len)
		for _, d := range ds {
			args = append(args, d.Get())
		}

		var i, j, k int
		i = int(args[0].Int())
		if args[1] == valueOfNil {
			j = container.Len()
		} else {
			j = int(args[1].Int())
		}
		if args[2] == valueOfNil {
			d.Set(container.Slice(i, j))
		} else {
			k = int(args[2].Int())
			d.Set(container.Slice3(i, j, k))
		}
	}()
	return d
}

func (frame *Frame) levalSliceArgs(globals varMap, lctx reflect.Value, node *Node) []*Dataflow {
	ds := make([]*Dataflow, 0, node.list.len)
	node.list.traverse(func(node *Node) {
		var d *Dataflow
		if node == nil {
			d = NewDataflowWith(valueOfNil)
		} else {
			d = frame.levalExpr(globals, lctx, node)
		}
		ds = append(ds, d)
	})
	return ds
}

func (frame *Frame) levalCall(globals varMap, lctx reflect.Value, node *Node) *Dataflow {
	d := NewDataflow()
	containerd := frame.levalExpr(globals, lctx, node.left)
	ds := frame.levalCallArgs(globals, lctx, node)
	go func() {
		container := containerd.Get()
		args := make([]reflect.Value, 0, node.list.len)
		for _, d := range ds {
			args = append(args, d.Get())
		}
		d.Set(container.Call(args)[0])
	}()
	return d
}

func (frame *Frame) levalCallArgs(globals varMap, lctx reflect.Value, node *Node) []*Dataflow {
	ds := make([]*Dataflow, 0, node.list.len)
	node.list.traverse(func(node *Node) {
		if node == nil {
			panic("illegal call")
		}
		d := frame.levalExpr(globals, lctx, node)
		ds = append(ds, d)
	})
	return ds
}
