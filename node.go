package gosk

import (
	"fmt"
	"reflect"
	"strings"
)

type nodeType int

const (
	nodeRoot nodeType = iota
	nodeLiteral
	nodeGlobalIdentifier
	nodeIdentifier
	nodeField

	nodeSelector
	nodeIndex
	nodeSlice
	nodeCall

	nodeOp
	nodeUnaryOp

	nodeBlock

	nodeReturn
	nodeAssign

	nodeIf
	nodeSwitch
	nodeCase
	nodeDefault
	nodeAssignBlock
)

func (typ nodeType) String() string {
	var s string
	switch typ {
	case nodeRoot:
		s = "nodeRoot"
	case nodeLiteral:
		s = "nodeLiteral"
	case nodeIdentifier:
		s = "nodeIdentifier"
	case nodeField:
		s = "nodeField"
	case nodeSelector:
		s = "nodeSelector"
	case nodeIndex:
		s = "nodeIndex"
	case nodeSlice:
		s = "nodeSlice"
	case nodeCall:
		s = "nodeCall"
	case nodeOp:
		s = "nodeOp"
	case nodeUnaryOp:
		s = "nodeUnaryOp"
	case nodeBlock:
		s = "nodeBlock"
	case nodeReturn:
		s = "nodeReturn"
	case nodeAssign:
		s = "nodeAssign"
	case nodeIf:
		s = "nodeIf"
	case nodeSwitch:
		s = "nodeSwitch"
	case nodeCase:
		s = "nodeCase"
	case nodeDefault:
		s = "nodeDefault"
	case nodeAssignBlock:
		s = "nodeAssignBlock"
	}
	return s
}

type opType int

const (
	opOrOr opType = iota
	opAndAnd
	opEqual
	opNonEqual
	opLessThan
	opLessThanOrEqual
	opGreaterThan
	opGreaterThanOrEqual
	opPlus
	opMinus
	opOr
	opMulti
	opDivide
	opMod
	opAnd
	opLeftShift
	opRightShift
	opAndNot

	opNegate
	opNot
	opBitNot
)

func (typ opType) String() string {
	var s string
	switch typ {
	case opOrOr:
		s = "opOrOr"
	case opAndAnd:
		s = "opAndAnd"
	case opEqual:
		s = "opEqual"
	case opNonEqual:
		s = "opNonEqual"
	case opLessThan:
		s = "opLessThan"
	case opLessThanOrEqual:
		s = "opLessThanOrEqual"
	case opGreaterThan:
		s = "opGreaterThan"
	case opGreaterThanOrEqual:
		s = "opGreaterThanOrEqual"
	case opPlus:
		s = "opPlus"
	case opMinus:
		s = "opMinus"
	case opOr:
		s = "opOr"
	case opMulti:
		s = "opMulti"
	case opDivide:
		s = "opDivide"
	case opMod:
		s = "opMod"
	case opAnd:
		s = "opAnd"
	case opLeftShift:
		s = "opLeftShift"
	case opRightShift:
		s = "opRightShift"
	case opAndNot:
		s = "opAndNot"
	case opNegate:
		s = "opNegate"
	case opNot:
		s = "opNot"
	case opBitNot:
		s = "opBitNot"
	}
	return s
}

type Node struct {
	typ nodeType

	cond *Node

	left  *Node
	right *Node

	list *NodeList

	val reflect.Value

	block *Node
}

func (node *Node) debug(depth int) {
	if node == nil {
		return
	}

	indent := strings.Repeat("  ", depth)

	fmt.Println(indent+"typ: ", node.typ)
	if node.val.IsValid() {
		if node.typ == nodeOp {
			fmt.Println(indent+"val: ", node.val.Interface().(opType))
		} else {
			fmt.Println(indent+"val: ", node.val.String())
		}
	}

	depth++

	if node.cond != nil {
		node.cond.debug(depth)
	}
	if node.left != nil {
		node.left.debug(depth)
	}
	if node.right != nil {
		node.right.debug(depth)
	}
	if node.list != nil {
		node.list.traverse(func(n *Node) { n.debug(depth) })
	}
	if node.block != nil {
		node.block.debug(depth)
	}
}

func (node *Node) Debug() {
	node.debug(0)
}

type listNode struct {
	next *listNode
	prev *listNode

	node *Node
}

type NodeList struct {
	first *listNode
	last  *listNode
	len   int
}

func NewNodeList(nodes ...*Node) *NodeList {
	nodeList := &NodeList{}

	for _, node := range nodes {
		nodeList.append(node)
	}

	return nodeList
}

func (nodeList *NodeList) traverse(f func(node *Node)) {
	iter := nodeList.first
	for iter != nil {
		f(iter.node)
		iter = iter.next
	}
}

func (nodeList *NodeList) prepend(node *Node) {
	first := nodeList.first

	newListNode := &listNode{
		next: first,
		node: node,
	}

	if first == nil {
		nodeList.last = newListNode
	} else {
		first.prev = newListNode
	}
	nodeList.first = newListNode
	nodeList.len++
}

func (nodeList *NodeList) append(node *Node) {
	last := nodeList.last

	newListNode := &listNode{
		prev: last,
		node: node,
	}

	if last == nil {
		nodeList.first = newListNode
	} else {
		last.next = newListNode
	}
	nodeList.last = newListNode
	nodeList.len++
}

func (nodeList *NodeList) applyable() {
	last := nodeList.last
	lastVal := last.node.val
	lastKind := lastVal.Kind()

	switch lastKind {
	case reflect.Array, reflect.Slice:
		for i := 0; i < lastVal.Len(); i++ {
			newListNode := &listNode{
				prev: last,
				node: lastVal.Index(i).Interface().(*Node),
			}

			last = last.next

			last.next = newListNode
			nodeList.last = newListNode
			nodeList.len++
		}
	default:
		panic("last argument should be slice or array")
	}
}
