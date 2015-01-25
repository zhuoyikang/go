package agent

import (
	"errors"
	"fmt"
	"io"
	"time"
	//"math"
)

//Step1. 实现一个Handle，处理网络包.

// Bz类型的包
type BzPacket struct {
	Data []byte
	Pos  int
	Type uint16
}

// 每一次Pack和UnPack前需要Reset()
func (pkt *BzPacket) Reset() {
	pkt.Pos = 0
}

func (pkt *BzPacket) Clear() {
	pkt.Data = make([]byte, 0, 64)
	pkt.Pos = 0
}

// 包处理器实现1，前2个字节视为包长，分配一个缓冲区收数据.
type HandlerBz struct {
	buffer []byte
}

//从Reader中获取一个完整的包
func (handle *HandlerBz) Read(reader io.Reader) (interface{}, error) {
	pkt := &BzPacket{}

	//fmt.Printf("%v\n", handle.buffer)

	//前两个字节包长度.
	n, err := io.ReadFull(reader, handle.buffer)
	if err != nil || n != 2 {
		return pkt, err
	}
	//fmt.Printf("n1: %v\n", n)
	pkt_length := uint16(handle.buffer[0])<<8 | uint16(handle.buffer[1])

	//后两个字节包类型
	n, err = io.ReadFull(reader, handle.buffer)
	//fmt.Printf("n2: %v\n", n)
	if err != nil || n != 2 {
		return pkt, err
	}
	pkt_type := uint16(handle.buffer[0])<<8 | uint16(handle.buffer[1])
	//fmt.Printf("p1: %v\n", pkt_type)

	buffer := make([]byte, pkt_length-4)
	n, err = io.ReadFull(reader, buffer)
	if err != nil || uint16(n) != (pkt_length-4) {
		return pkt, errors.New("pkt length error")
	} else {
		pkt.Data = buffer
		pkt.Pos = 0
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

func BzReadu16(p *BzPacket) (ret uint16, err error) {
	if p.Pos+2 > len(p.Data) {
		err = errors.New("read uint16 failed")
		return
	}

	buf := p.Data[p.Pos : p.Pos+2]
	ret = uint16(buf[0])<<8 | uint16(buf[1])
	p.Pos += 2
	return
}

func BzWriteu16(p *BzPacket, v uint16) (err error) {
	p.Data = append(p.Data, byte(v>>8), byte(v))
	p.Pos += 2
	return
}

func BzReads16(p *BzPacket) (ret int16, err error) {
	if p.Pos+2 > len(p.Data) {
		err = errors.New("read uint16 failed")
		return
	}

	buf := p.Data[p.Pos : p.Pos+2]
	ret = int16(buf[0])<<8 | int16(buf[1])
	p.Pos += 2
	return
}

func BzWrites16(p *BzPacket, v int16) (err error) {
	p.Data = append(p.Data, byte(v>>8), byte(v))
	p.Pos += 2
	return
}

func BzReadu32(p *BzPacket) (ret uint32, err error) {
	if p.Pos+4 > len(p.Data) {
		err = errors.New("read uint32 failed")
		return
	}

	buf := p.Data[p.Pos : p.Pos+4]
	ret = uint32(buf[0])<<24 | uint32(buf[1])<<16 | uint32(buf[2])<<8 |
		uint32(buf[3])
	p.Pos += 4
	return
}

func BzWriteu32(p *BzPacket, v uint32) (err error) {
	p.Data = append(p.Data, byte(v>>24), byte(v>>16),
		byte(v>>8), byte(v))
	p.Pos += 4
	return
}

func BzReads32(p *BzPacket) (ret int32, err error) {
	_ret, _err := BzReadu32(p)
	ret = int32(_ret)
	err = _err
	return
}

func BzWrites32(p *BzPacket, v int32) (err error) {
	BzWriteu32(p, uint32(v))
	return
}

func BzReadstring(p *BzPacket) (ret string, err error) {
	if p.Pos+2 > len(p.Data) {
		err = errors.New("read string header failed")
		return
	}

	size, _ := BzReadu16(p)
	if p.Pos+int(size) > len(p.Data) {
		err = errors.New("read string Data failed")
		return
	}

	bytes := p.Data[p.Pos : p.Pos+int(size)]
	p.Pos += int(size)
	ret = string(bytes)
	return
}

func BzWritestring(p *BzPacket, v string) (err error) {
	bytes := []byte(v)
	BzWriteu16(p, uint16(len(bytes)))
	p.Data = append(p.Data, bytes...)
	p.Pos += len(bytes)
	return
}

// 创建一个完整的数据包
func MakePacketData(api uint16, p *BzPacket) (data []byte) {
	length := 2 + len(p.Data)
	data = append(data, byte(length>>8), byte(length))
	data = append(data, byte(api>>8), byte(api))
	data = append(data, p.Data...)

	return
}

type BzHandlerMAP map[uint16]func(*Session, *BzPacket)

//Step2. 实现一个Agent.
type AgentBz struct {
	handlerMap BzHandlerMAP
}

func (gs *AgentBz) Start(session *Session) {
	fmt.Printf("%s\n", "bz session start")
}

func (gs *AgentBz) HandlePkt(session *Session, pkti interface{}) {
	fmt.Printf("%s\n", "bz session handle")
	pkt := pkti.(*BzPacket)
	handler := gs.handlerMap[pkt.Type]
	if handler == nil {
		return
	}
	handler(session, pkt)
	return
}

func (gs *AgentBz) Stop(session *Session) {
	fmt.Printf("%s\n", "bz session stop")
}

//Step3. main
func AgentBzMain() {
	agentBz := &AgentBz{}
	agentBz.handlerMap = MakeBzDemoHandler()
	agt := MakeAgent("tcp", "0.0.0.0:8080", agentBz, &HandlerBz{})
	go func() {
		time.Sleep(time.Second * 60)
		fmt.Printf("%s\n", "test stop")
		agt.Stop()
	}()
	agt.Signal()
	agt.Run()
	fmt.Printf("%s\n", "end")
}
