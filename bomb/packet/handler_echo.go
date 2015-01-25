package packet

import (
	//"errors"
	//"math"
	"io"
	//"fmt"
	//"net"
)

const (
	ECHO_PREALLOC_BUFSIZE = 1024
)

// 包处理器实现1，前2个字节视为包长，分配一个缓冲区收数据.
type HandlerEcho struct {
}

//从Reader中获取一个完整的包
func (handle *HandlerEcho) Read(reader io.Reader) (Packet, error) {
	pkt := Packet{}
	buffer := make([]byte, ECHO_PREALLOC_BUFSIZE)
	_, err := reader.Read(buffer)
	pkt.Data=buffer
	pkt.Pos=0
	return pkt, err
}

func (handle *HandlerEcho) Write(writer io.Writer, pkt Packet) error {
	_, err := writer.Write(pkt.Data)
	return err
}


func (handle *HandlerEcho) New() HandlerI {
	return &HandlerEcho{}
}
