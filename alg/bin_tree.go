package alg

import (
//	"fmt"
)

//------------------------------------------------------------------------------
// 普通2叉树实现
//------------------------------------------------------------------------------

type binNode struct {
	left   *binNode
	right  *binNode
	parent *binNode
	key    BinKeyI
	value  interface{}
}

//简单删除
func (z *binNode) delete(binTree *BinTree) *binNode {
	var y *binNode
	var x *binNode
	if z.left == nil || z.right == nil {
		y = z
	} else {
		y = z.successor()
	}

	if y.left != nil {
		x = y.left
	} else {
		x = y.right
	}

	if x != nil {
		x.parent = y.parent
	}

	if y.parent == nil {
		binTree.root = x
	} else if y == y.parent.left {
		y.parent.left = x
	} else {
		y.parent.right = x
	}

	if y != z {
		z.key = y.key
		z.value = y.value
	}
	return y
}

//查找一棵树中最小的节点
func (node *binNode) minimum() *binNode {
	for node.left != nil {
		node = node.left
	}
	return node
}

//查找一棵树中最小的节点
func (node *binNode) maximum() *binNode {
	for node.right != nil {
		node = node.right
	}
	return node
}

//查找后继
func (node *binNode) successor() *binNode {
	if node.right != nil {
		return node.right.minimum()
	}

	left := node
	parent := node.parent

	for parent != nil && node == parent.right {
		left = node
		node = parent
		parent = node.parent
	}

	switch {
	case parent == nil:
		if left == node.right {
			return nil
		}
	case node == parent.left:
		return parent
	}

	return nil
}

type BinTree struct {
	root *binNode
}

func (bin *BinTree) Insert(key BinKeyI, value interface{}) interface{} {
	node := bin.root
	if node == nil {
		node := new(binNode)
		node.key = key
		node.value = value
		bin.root = node
		return nil
	}
	for node != nil {
		ret := node.key.Cmp(key)
		switch {
		case ret == 0:
			old_v := node.value
			node.value = value
			return old_v
		case ret > 0:
			if node.left == nil {
				node_n := new(binNode)
				node_n.key = key
				node_n.value = value
				node_n.parent = node
				node.left = node_n
				return nil
			} else {
				node = node.left
			}
		case ret < 0:
			if node.right == nil {
				node_n := new(binNode)
				node_n.key = key
				node_n.value = value
				node_n.parent = node
				node.right = node_n
				return nil
			} else {
				node = node.right
			}
		}
	}
	panic("show not be here")
}

func (bin *BinTree) Lookup(key BinKeyI) interface{} {
	node := bin.root
	for node != nil {
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

func (bin *BinTree) Delete(key BinKeyI) interface{} {
	node := bin.root
J:
	for node != nil {
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
	node.delete(bin)
	return value
}

func (bin *BinTree) Min() interface{} {
	node := bin.root
	if node == nil {
		return nil
	}

	for node.left != nil {
		node = node.left
	}

	if node == nil {
		return nil
	}
	return node.value
}

func (bin *BinTree) Max() interface{} {
	node := bin.root
	if node == nil {
		return nil
	}

	for node.right != nil {
		node = node.right
	}

	if node == nil {
		return nil
	}
	return node.value
}

func (bin *BinTree) Clear() {
	bin.root = nil
}

func (bin *BinTree) IsEmpty() bool {
	return bin.root == nil
}

func (bin *BinTree) Travel(fun func(binKey BinKeyI, value interface{}) bool) {
	node := bin.root
	if node == nil {
		return
	}
	node = bin.root.minimum()
	for node != nil {
		if fun(node.key, node.value) == false {
			break
		} else {
			node = node.successor()
		}
	}
}
