package bin_tree

import (
	"testing"
	"fmt"
)

func TestRbTest(t *testing.T) {
	//fmt.Printf("%s\n", "test")
}

// 检查颜色设置是否正确。
func check_color(t *testing.T, get uint8, want uint8) {
	if get != want {
		t.Errorf("coolor want1 %d get %d", want, get)
	}
}


func TravelRbTree(rbtree RBTree) {
	rbtree.TravelNode(func(node *rbNode) bool {
		fmt.Printf("%d %d %d\n", node.key, node.value, node.getColor())
		return true
	})
}


func TestRbNodeTest(t *testing.T) {
	//颜色测试
	node := rbNode{}
	check_color(t, node.getColor(), RB_NODE_COLOR_R)
	node.setColor(RB_NODE_COLOR_B)
	check_color(t, node.getColor(), RB_NODE_COLOR_B)
	node.setColor(RB_NODE_COLOR_R)
	check_color(t, node.getColor(), RB_NODE_COLOR_R)

	//基础测试
	rbTree := RBTree{}
	rbTree.Insert(BinInt(10),10)
	rbTree.Insert(BinInt(9),9)
	rbTree.Insert(BinInt(8),8)
	rbTree.Insert(BinInt(7),7)
	rbTree.Insert(BinInt(6),6)
	//TravelRbTree(rbTree)

	rbTree2 := RBTree{}
	RandomBinTreeTest(t, &rbTree2)
	// BinTreeTest(t, &RBTree{})
}
