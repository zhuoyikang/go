// 基础类型定义

package packet

import (
	//"errors"
	//"math"
	"io"
	//"net"
)

// 从网络层获取的一个完整的包
type Packet struct {
	Data []byte
	Pos  int
}

// 包处理器接口
type HandlerI interface {
	//从Reader中获取一个完整的包
	Read(io.Reader) (Packet, error)
	//向Writer写入一个完整的包
	Write(io.Writer, Packet)
	//返回一个新的Handle.每个session一个
	New() HandlerI
}
