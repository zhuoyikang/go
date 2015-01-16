package alg

import (
	"fmt"
	"math/rand"
	"testing"
)

func check_1(t *testing.T, key int, want interface{}, btree BinTreeI) {
	v := btree.Lookup(BinInt(key))
	switch {
	case (want == nil):
		if v != nil {
			t.Errorf("want1 %d get %d key %d", want, v, key)
		}
	case want.(int) != v.(int):
		t.Errorf("want1 %d get %d", want, v)
	}
}

func check_2(t *testing.T, get interface{}, want int, btree BinTreeI) {
	if get.(int) != want {
		t.Errorf("want1 %d get %d", want, get)
	}
}

var (
	AddList =  make([]int, 0)
	DelList = make([]int, 0) //[]int{13, 16, 5}
)

func init() {
	init := []int{15, 5, 16, 3, 12, 10, 13, 6, 7, 16, 20, 18, 23}
	del := []int{13, 16, 5}

	// 增加列表
	for _,v :=  range init {
		AddList = append(AddList, v)
	}

	for i := 0; i < 200; i++ {
		value := rand.Intn(10000)
		AddList = append(AddList, value)
	}
	AddList = append(AddList, 10001)
	AddList = append(AddList, -1)

	// 删除列表
	for _,v :=  range del {
		DelList = append(DelList, v)
	}

	for i := 0; i < 100; i++ {
		value := rand.Intn(10000)
		DelList = append(DelList, value)
	}
}

func BinTreeTest(t *testing.T, btree BinTreeI) {
	for _, v := range AddList {
		btree.Insert(BinInt(v), v)
	}

	for _, v := range AddList {
		check_1(t, v, v, btree)
	}

	for _, v := range DelList {
		btree.Delete(BinInt(v))
	}

	for _, v := range DelList {
		check_1(t, v, nil, btree)
	}

	check_2(t, btree.Max(), 10001, btree)
	check_2(t, btree.Min(), -1, btree)
	TravelBinTree(btree)
}

func TravelBinTree(btree BinTreeI) {
	btree.Travel(func(binKey BinKeyI, value interface{}) bool {
		fmt.Printf("%d\n", value)
		return true
	})
}

func RandomBinTreeTest(t *testing.T, btree BinTreeI) {
	for i := 0; i < 10000; i++ {
		value := rand.Intn(1000)
		//fmt.Printf("insert %d\n", value)
		btree.Insert(BinInt(value), value)
		check_2(t, value, value, btree)
	}
	//Travel(btree)
}


