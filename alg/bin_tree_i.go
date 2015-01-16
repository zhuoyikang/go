package alg

import (
//	"fmt"
)

//------------------------------------------------------------------------------
// 2叉树的key，仿sort实现
//------------------------------------------------------------------------------

type BinKeyI interface {
	Cmp(a BinKeyI) int
}

type BinInt int64

func (x BinInt) Cmp(y BinKeyI) int {
	y1 := y.(BinInt)
	return int(x - y1)
}

type BinString string

func (x BinString) Cmp(y BinKeyI) int {
	y1 := y.(BinString)
	switch {
	case x < y1:
		return -1
	case x > y1:
		return 1
	}
	return 0
}

type BinFloat float64

func (x BinFloat) Cmp(y BinKeyI) int {
	y1 := y.(BinFloat)
	switch {
	case x < y1:
		return -1
	case x > y1:
		return 1
	}
	return 0
}

//------------------------------------------------------------------------------
// 基本的2叉树接口
//------------------------------------------------------------------------------

type BinTreeI interface {
	Insert(key BinKeyI, value interface{}) interface{}
	Lookup(key BinKeyI) interface{}
	Delete(key BinKeyI) interface{}
	Min() interface{}
	Max() interface{}
	Clear()
	IsEmpty() bool
	Travel(fun func(binKey BinKeyI, value interface{}) bool)
}
