package packet

import (
	//"errors"
	//"math"
	"io"
	//"fmt"
	//"net"
)

const (
	PREALLOC_BUFSIZE = 1024
)

// 包处理器实现1，前2个字节视为包长，分配一个缓冲区收数据.
type HandlerEcho struct {
	buffer []byte
}

//从Reader中获取一个完整的包
func (handle *HandlerEcho) Read(reader io.Reader) (Packet, error) {
	pkt := Packet{}
	_, err := reader.Read(handle.buffer)
	pkt.Data=handle.buffer
	pkt.Pos=0
	return pkt, err
}

func (handle *HandlerEcho) Write(writer io.Writer, pkt Packet) {
}

func (handle *HandlerEcho) New() (HandlerI) {
	h := HandlerEcho{}
	h.buffer = make([]byte, PREALLOC_BUFSIZE)
	return &h
}
