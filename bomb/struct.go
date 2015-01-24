package main

import (
	"fmt"
	"time"
)

// 游戏物体，所有物体均有2维坐标属性
type EBase struct {
	X         int
	Y         int
	CanBlock  bool //是否阻挡其他物体.
	CanBombed bool //是否可以被砸毁
}

const (
	CAN_BLOCK_Y  = true     //是否会阻挡其他物体
	CAN_BLOCK_N  = false    //
	CAN_BOMBED_Y = true     //是否可以被炸毁
	CAN_BOMBED_N = false	//
)

// 空地
const (
	ESPACE_T_BASE   = 0    //基地：初始化地点
	ESPACE_T_NORMAL = iota //普通空地：自由穿越
)

// 空地，玩家可以自由行走.
type ESpace struct {
	_type int
	EBase
}

// 创建一块新的空地
func MakeESpace(x int, y int, t int) ESpace {
	s := ESpace{_type: t}
	s.X = x
	s.Y = y
	s.CanBlock = CAN_BLOCK_N
	s.CanBombed = CAN_BOMBED_N
	return s
}

// 创建一个岩石
func MakeERock(x int, y int, t int) ESpace {
	s := ESpace{_type: t}
	s.X = x
	s.Y = y
	s.CanBlock = CAN_BLOCK_Y
	s.CanBombed = CAN_BOMBED_Y
	return s
}

// 炸弹类型
type EBomb struct {
	Power          int   //炸弹威力，十字型摧毁Rock
	DetonationTime int64 //爆炸时间,utc_s
	EBase
}

// 创建一个deTime秒后爆炸的炸弹.
func MakeEBomb(x int, y int, power int, deTime int64) EBomb {
	s := EBomb{Power: power}
	s.X = x
	s.Y = y
	s.DetonationTime = time.Now().Unix() + deTime
	return s
}

// 角色
// 可以移动
// 可以放炸弹
type ESprite struct {
	EBase
}

//
func (sprite *ESprite) SetPosition(x, y int) {
	sprite.X = x
	sprite.Y = y
}

// 精灵
const (
	ESPRITE_T_PLAYER = 0 //基地：初始化地点
	ESPRITE_T_ROBOT  = iota
)

//
func MakeESprite(x int, y int, t int) {
	s := ESpace{_type: t}
	s.X = x
	s.Y = y
	s.CanBlock = CAN_BLOCK_Y
	s.CanBombed = CAN_BOMBED_Y
}

func main() {
	x := MakeESpace(1, 21, ESPACE_T_BASE)
	fmt.Printf("%v %d %d\n", x, x.X, x.Y)

	b := MakeEBomb(1, 1, 2, 2)
	fmt.Printf("%v\n", b)
}
