// 基础类型定义

package packet

import (
	//"errors"
	//"math"
	"io"
	//"net"
)

// 裸包:从网络层获取，里面是2进制
type Packet struct {
	Data []byte
	Pos  int
}

//
func (pkt *Packet) Reset() {
	pkt.Pos = 0
}

//
func (pkt *Packet) Clear() {
	pkt.Data = make([]byte, 0, 64)
	pkt.Pos = 0
}

// 包处理器接口
// handle需要为每个session分配一个实例
type HandlerI interface {
	//从Reader中获取一个完整的包
	Read(io.Reader) (Packet, error)
	//向Writer写入一个完整的包
	Write(io.Writer, Packet) error
	//分配一个新的handler.
	New() HandlerI
}

// 逻辑包:所有逻辑包都要实现对裸包的转换.
// 指针作为参数是方便循环调用.
type PacketI interface {
	//从Packet中转为换PacketI
	UnPack(*Packet) error
	//将PacektI转为为Packet.
	Pack(*Packet) error
}
