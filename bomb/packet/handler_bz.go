// 对bz来说，头两个字节是无符号包长度。
// 后两个字节是包类型
// 最后是包内容.

package packet

import (
	"errors"
	//"math"
	"io"
	"fmt"
	//"net"
)

const (
	PREALLOC_BUFSIZE = 1024
)

// 包处理器实现1，前2个字节视为包长，分配一个缓冲区收数据.
type HandlerBz struct {
	buffer []byte
}

//从Reader中获取一个完整的包
func (handle *HandlerBz) Read(reader io.Reader) (Packet, error) {
	pkt := Packet{}
	header := []byte{0, 0}

	//前两个字节包长度.
	n, err := io.ReadFull(reader, header)
	if err != nil || n != 2 {
		return pkt, err
	}

	ret := uint16(handle.buffer[0]) << 8 | uint16(handle.buffer[1])
	if ret < PREALLOC_BUFSIZE {
		buffer := make([]byte, ret)
		n, err := io.ReadFull(reader, buffer)
		if err != nil || uint16(n) != ret {
			return pkt, errors.New("pkt length error")
		} else {
			pkt.Data = buffer
			pkt.Len = n
			return pkt, nil
		}
	} else {
		return pkt, errors.New("pkt too long")
	}
}

func (handle *HandlerBz) Write(writer io.Writer, pkt Packet) error {
	want := len(pkt.Data)
	n := 0
	for {
		ret, err := writer.Write(pkt.Data)
		switch {
		case (err != nil) :
			return err
		case n == want:
			return nil
		case n > want:
			fmt.Printf("%s\n", "write bug")
			return nil
		}
		n += ret
	}
}
