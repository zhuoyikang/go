package agent

import (
	"errors"
	"fmt"
	"io"
	//"time"
	//"math"
)

//Step1. 实现一个Handle，处理网络包.

// Bz类型的包
type BzPacket struct {
	Data []byte
	Type uint16
}

// 每一次Pack和UnPack前需要Reset()
func (pkt *BzPacket) Reset() {
	pkt.Data = []byte{}
}

// 包处理器实现1，前2个字节视为包长，分配一个缓冲区收数据.
type HandlerBz struct {
	buffer []byte
}

//从Reader中获取一个完整的包
func (handle *HandlerBz) Read(reader io.Reader) (interface{}, error) {
	pkt := &BzPacket{}

	//前两个字节包长度.
	n, err := io.ReadFull(reader, handle.buffer)
	if err != nil || n != 2 {
		return pkt, err
	}

	pkt_length := uint16(handle.buffer[0])<<8 | uint16(handle.buffer[1])
	fmt.Printf("pkt_length: %v\n", pkt_length)
	//后两个字节包类型
	n, err = io.ReadFull(reader, handle.buffer)
	//fmt.Printf("n2: %v\n", n)
	if err != nil || n != 2 {
		return pkt, err
	}
	pkt_type := uint16(handle.buffer[0])<<8 | uint16(handle.buffer[1])
	fmt.Printf("pkt_type: %v\n", pkt_type)

	buffer := make([]byte, pkt_length-4)
	n, err = io.ReadFull(reader, buffer)
	if err != nil || uint16(n) != (pkt_length-4) {
		return pkt, errors.New("pkt length error")
	} else {
		pkt.Data = buffer
		pkt.Type = pkt_type
		return pkt, nil
	}
}

func (handle *HandlerBz) Write(writer io.Writer, pkt_i interface{}) error {
	pkt := pkt_i.(*BzPacket)
	want := len(pkt.Data)
	n := 0
	for {
		ret, err := writer.Write(pkt.Data)
		switch {
		case (err != nil):
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

func (handle *HandlerBz) New() HandlerI {
	return &HandlerBz{buffer: make([]byte, 2)}
}

//------------------------------------------------------------------------------
// 基础类型的Pack和UnPack
//------------------------------------------------------------------------------

func BzReadbyte(datai []byte) (data []byte, ret byte, err error) {
	data = datai
	if 1 > len(data) {
		err = errors.New("read byte failed")
		return
	}
	ret = data[0]
	data = data[1:]
	return
}

func BzWritebyte(datai []byte, v byte) (data []byte, err error) {
	data = datai
	data = append(data, byte(v))
	return
}

func BzReaduint16(datai []byte) (data []byte, ret uint16, err error) {
	data = datai
	if 2 > len(data) {
		err = errors.New("read uint16 failed")
		return
	}

	buf := data[0:2]
	ret = uint16(buf[0])<<8 | uint16(buf[1])
	data = data[2:]
	return
}

func BzWriteuint16(datai []byte, v uint16) (data []byte, err error) {
	data = datai
	data = append(data, byte(v>>8), byte(v))
	return
}

func BzReadint16(datai []byte) (data []byte, ret int16, err error) {
	if 2 > len(data) {
		err = errors.New("read uint16 failed")
		return
	}

	buf := data[0:2]
	ret = int16(buf[0])<<8 | int16(buf[1])
	data = data[2:]
	return
}

func BzWriteint16(datai []byte, v int16) (data []byte, err error) {
	data = datai
	data = append(data, byte(v>>8), byte(v))
	return
}

func BzReaduint32(datai []byte) (data []byte, ret uint32, err error) {
	data = datai
	if 4 > len(data) {
		err = errors.New("read uint32 failed")
		return
	}

	buf := data[0:4]
	ret = uint32(buf[0])<<24 | uint32(buf[1])<<16 | uint32(buf[2])<<8 |
		uint32(buf[3])

	data = data[4:]
	return
}

func BzWriteuint32(datai []byte, v uint32) (data []byte, err error) {
	data = datai
	data = append(data, byte(v>>24), byte(v>>16),
		byte(v>>8), byte(v))
	return
}

func BzReadint32(datai []byte) (data []byte, ret int32, err error) {
	data, ret1, err := BzReaduint32(datai)
	ret = int32(ret1)
	return
}

func BzWriteint32(datai []byte, v int32) (data []byte, err error) {
	return BzWriteuint32(datai, uint32(v))
}

func BzReadstring(datai []byte) (data []byte, ret string, err error) {
	data, size, err := BzReaduint16(datai)
	if err != nil {
		return
	}
	if int(size) > len(data) {
		err = errors.New("read string failed")
	}

	bytes := data[0:int(size)]
	ret = string(bytes)
	data = data[int(size):]
	return
}

func BzWritestring(datai []byte, str string) (data []byte, err error) {
	bytes := []byte(str)
	data, err = BzWriteuint16(datai, uint16(len(bytes)))
	data = append(data, bytes...)
	return
}

// 创建一个完整的数据包
func MakePacketData(api uint16, datai []byte) (data []byte) {
	length := 4 + len(datai)
	data = append(data, byte(length>>8), byte(length))
	data = append(data, byte(api>>8), byte(api))
	data = append(data, datai...)

	return
}

type BzHandlerMAP map[uint16]func(*Session, *BzPacket)

// 所有的Bz包都需要继承这个Agent.
type AgentBz struct {
	Agent
	handlerMap BzHandlerMAP
}
