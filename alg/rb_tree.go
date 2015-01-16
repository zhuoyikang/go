package alg

import (
	"fmt"
)

const (
	RB_NODE_COLOR_R = 0
	RB_NODE_COLOR_B = 1
)

var (
	_nil_node *rbNode
)

func init() {
	_nil_node = &rbNode{}
	_nil_node.setColor(RB_NODE_COLOR_B)
}

type rbNode struct {
	left   *rbNode
	right  *rbNode
	parent *rbNode
	key    BinKeyI
	value  interface{}
	status uint8
}

// 最后一位放置颜色类型
func (node *rbNode) getColor() uint8 {
	return node.status & 0x01
}

func (node *rbNode) setColor(color uint8) {
	node.status = (node.status & 0xFE) | color
}

//查找一棵树中最小的节点
func (node *rbNode) minimum() *rbNode {
	for node.left != _nil_node {
		node = node.left
	}
	return node
}

//查找一棵树中最小的节点
func (node *rbNode) maximum() *rbNode {
	for node.right != _nil_node {
		node = node.right
	}
	return node
}

//查找后继
func (node *rbNode) successor() *rbNode {
	if node.right != _nil_node {
		return node.right.minimum()
	}

	left := node
	parent := node.parent

	for parent != _nil_node && node == parent.right {
		left = node
		node = parent
		parent = node.parent
	}

	switch {
	case parent == _nil_node:
		if left == node.right {
			return _nil_node
		}
	case node == parent.left:
		return parent
	}

	return _nil_node
}

//简单删除
func (z *rbNode) delete(rbTree *RBTree) *rbNode {
	var y *rbNode
	var x *rbNode
	if z.left == _nil_node || z.right == _nil_node {
		y = z
	} else {
		y = z.successor()
	}

	if y.left != _nil_node {
		x = y.left
	} else {
		x = y.right
	}

	if x != _nil_node {
		x.parent = y.parent
	}

	if y.parent == _nil_node {
		rbTree.root = x
	} else if y == y.parent.left {
		y.parent.left = x
	} else {
		y.parent.right = x
	}

	if y != z {
		z.key, y.key = y.key, z.key
		z.value, y.value = y.value, z.value
	}

	if y.getColor() == RB_NODE_COLOR_B {
		rbTree.deleteFix(x)
	}
	return y
}

// 处理叶子nil结点为黑色
func getColor(node *rbNode) uint8 {
	if node == nil {
		return RB_NODE_COLOR_B
	} else {
		return node.getColor()
	}
}

func (x *rbNode) leftRotate(T *RBTree) {
	y := x.right
	x.right = y.left
	if y.left != _nil_node {
		y.left.parent = x
	}
	y.parent = x.parent
	if x.parent == _nil_node {
		T.root = y
	} else if x == x.parent.left {
		x.parent.left = y
	} else {
		x.parent.right = y
	}
	y.left = x
	x.parent = y
}

func (y *rbNode) rightRotate(T *RBTree) {
	x := y.left
	y.left = x.right
	if x.right != _nil_node {
		x.right.parent = y
	}
	x.parent = y.parent
	if y.parent == _nil_node {
		T.root = x
	} else if y == y.parent.left {
		y.parent.left = x
	} else {
		y.parent.right = x
	}
	x.right = y
	y.parent = x
}

//------------------------------------------------------------------------------
// RBTree
//------------------------------------------------------------------------------

type RBTree struct {
	root *rbNode
}

func (rb *RBTree) deleteFix(x *rbNode) {
	var w *rbNode
	for (x != rb.root) && (x.getColor() == RB_NODE_COLOR_B) {
		if x == x.parent.left {
			w = x.parent.right
			if w.getColor() == RB_NODE_COLOR_R {
				w.setColor(RB_NODE_COLOR_B)
				x.parent.setColor(RB_NODE_COLOR_R)
				x.parent.leftRotate(rb)
				w = x.parent.right
			}
			if w.left.getColor() == RB_NODE_COLOR_B &&
				w.right.getColor() == RB_NODE_COLOR_B {
				w.setColor(RB_NODE_COLOR_R)
				x = x.parent
			} else if w.right.getColor() == RB_NODE_COLOR_B {
				w.left.setColor(RB_NODE_COLOR_B)
				w.setColor(RB_NODE_COLOR_R)
				w.rightRotate(rb)
				w = x.parent.right
			} else {
				w.setColor(x.parent.getColor())
				x.parent.setColor(RB_NODE_COLOR_B)
				w.right.setColor(RB_NODE_COLOR_B)
				x.parent.leftRotate(rb)
				x = rb.root
			}

		} else if x == x.parent.right {
			w = x.parent.left
			if w.getColor() == RB_NODE_COLOR_R {
				w.setColor(RB_NODE_COLOR_B)
				x.parent.setColor(RB_NODE_COLOR_R)
				x.parent.leftRotate(rb)
				w = x.parent.left
			}
			if w.right.getColor() == RB_NODE_COLOR_B &&
				w.left.getColor() == RB_NODE_COLOR_B {
				w.setColor(RB_NODE_COLOR_R)
				x = x.parent
			} else if w.left.getColor() == RB_NODE_COLOR_B {
				w.right.setColor(RB_NODE_COLOR_B)
				w.setColor(RB_NODE_COLOR_R)
				w.leftRotate(rb)
				w = x.parent.left
			} else {
				w.setColor(x.parent.getColor())
				x.parent.setColor(RB_NODE_COLOR_B)
				w.left.setColor(RB_NODE_COLOR_B)
				x.parent.leftRotate(rb)
				x = rb.root
			}
		} else {
			fmt.Printf("%s\n", "ignore")
		}
	}

}

func (rb *RBTree) insertFix(z *rbNode) {
	var y *rbNode
	for getColor(z.parent) == RB_NODE_COLOR_R {
		//fmt.Printf("z %d\n", z.key)
		if z.parent == z.parent.parent.left {
			y = z.parent.parent.right
			if getColor(y) == RB_NODE_COLOR_R {
				//fmt.Printf("%s\n", "case 1")
				z.parent.setColor(RB_NODE_COLOR_B)
				y.setColor(RB_NODE_COLOR_B)
				z.parent.parent.setColor(RB_NODE_COLOR_R)
				z = z.parent.parent
			} else if z == z.parent.right {
				//fmt.Printf("%s\n", "case 2")
				z = z.parent
				z.leftRotate(rb)
			} else {
				//fmt.Printf("%s\n", "case 3")
				z.parent.setColor(RB_NODE_COLOR_B)
				z.parent.parent.setColor(RB_NODE_COLOR_R)
				z.parent.parent.rightRotate(rb)
			}
		} else if z.parent == z.parent.parent.right {
			y = z.parent.parent.left
			if getColor(y) == RB_NODE_COLOR_R {
				z.parent.setColor(RB_NODE_COLOR_B)
				y.setColor(RB_NODE_COLOR_B)
				z.parent.parent.setColor(RB_NODE_COLOR_R)
				z = z.parent.parent
			} else if z == z.parent.left {
				z = z.parent
				z.rightRotate(rb)
			} else {
				z.parent.setColor(RB_NODE_COLOR_B)
				z.parent.parent.setColor(RB_NODE_COLOR_R)
				z.parent.parent.leftRotate(rb)
			}
		} else {
			fmt.Printf("%s\n", "ignore")
		}
	}
	rb.root.setColor(RB_NODE_COLOR_B)
}

func (rb *RBTree) Insert(key BinKeyI, value interface{}) interface{} {
	node := rb.root
	var node_n *rbNode
	if rb.root == nil {
		node := new(rbNode)
		node.key = key
		node.value = value
		node.parent = _nil_node
		node.left = _nil_node
		node.right = _nil_node
		node.setColor(RB_NODE_COLOR_B)
		rb.root = node
		return nil
	}
J:
	for node != _nil_node {
		ret := node.key.Cmp(key)
		//fmt.Printf("%s %d\n", "tt", node.key) //
		switch {
		case ret == 0:
			old_v := node.value
			node.value = value
			return old_v
		case ret > 0:
			if node.left == _nil_node {
				node_n = new(rbNode)
				node_n.key = key
				node_n.value = value
				node_n.parent = node
				node.left = node_n
				break J
			} else {
				node = node.left
			}
		case ret < 0:
			//fmt.Printf("%s %d\n", "this", key) //
			if node.right == _nil_node {
				node_n = new(rbNode)
				node_n.key = key
				node_n.value = value
				node_n.parent = node
				node.right = node_n
				break J
			} else {
				node = node.right
			}
		}
	}
	node_n.setColor(RB_NODE_COLOR_R)
	//fmt.Printf("%s %d\n", "set attr", node_n.key)
	node_n.left = _nil_node
	node_n.right = _nil_node
	rb.insertFix(node_n)
	return nil
}

func (rb *RBTree) Lookup(key BinKeyI) interface{} {
	node := rb.root
	if node == nil {
		return nil
	}

	for node != _nil_node {
		ret := node.key.Cmp(key)
		switch {
		case ret == 0:
			return node.value
		case ret > 0:
			node = node.left
		case ret < 0:
			node = node.right
		}
	}
	return nil
}

func (rb *RBTree) Delete(key BinKeyI) interface{} {
	node := rb.root
	if rb.root == nil {
		return nil
	}
J:
	for node != _nil_node {
		ret := node.key.Cmp(key)
		switch {
		case ret == 0:
			break J
		case ret > 0:
			node = node.left
		case ret < 0:
			node = node.right
		}
	}

	if node == nil {
		return nil
	}
	value := node.value
	node.delete(rb)
	return value
}

func (rb *RBTree) Min() interface{} {
	node := rb.root
	if node == nil {
		return nil
	}

	for node.left != _nil_node {
		node = node.left
	}

	if node == _nil_node {
		return nil
	}
	return node.value
}

func (rb *RBTree) Max() interface{} {
	node := rb.root
	if node == nil {
		return nil
	}

	for node.right != _nil_node {
		node = node.right
	}

	if node == _nil_node {
		return nil
	}
	return node.value
}

func (rb *RBTree) Clear() {
	rb.root = nil
}

func (rb *RBTree) IsEmpty() bool {
	return rb.root == nil
}

func (rb *RBTree) Travel(fun func(binKey BinKeyI, value interface{}) bool) {
	if rb.root == nil {
		return
	}
	node := rb.root.minimum()
	for node != _nil_node {
		if fun(node.key, node.value) == false {
			break
		} else {
			node = node.successor()
		}
	}
}

// 遍历结点
func (rb *RBTree) TravelNode(fun func(node *rbNode) bool) {
	if rb.root == nil {
		return
	}
	node := rb.root.minimum()
	for node != _nil_node {
		if fun(node) == false {
			break
		} else {
			node = node.successor()
		}
	}
}
