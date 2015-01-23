package agent

import (
	"fmt"
)

type PacketHanderI interface {
	HandleInLoop(session Session)
	HandlePac(session Session)
	HandleAll(session Session)
}

// 一个最简单地包处理:Echo
type PacketEcho struct {
}

func (pkg *PacketEcho) HandleIn(session Session)  {
	fmt.Printf("%s\n", "in")
}


func (pkg *PacketEcho) HandleOut(session Session)  {
	fmt.Printf("%s\n", "out")
}

func (pkg *PacketEcho) Handle(session Session)  {
	fmt.Printf("%s\n", "in")
}
