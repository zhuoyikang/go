/*
所有的packet类型都需要实现对Packet的转换关系。

PacketBz  ->  packet(2进制包)
PacketBz  <-  packet(2进制包)


一种网络数据序列化的方式.

每一个数据包都由三部分组成： 包长度(2字节)+包类型(2字节)+包数据(剩余字节)。
每个正常包的前两个字节都是包长，也就是你必须先收两个字节，确定之后的包的长度，等待完整包收完后再处理。

1.包类型为2个字节，最多有65535个不同的包，完全够用了。

一个完整的包收好后，将通过包类型来决定如何处理，包类型的描述我定义在了proto/api.txt文件中，每个包有以下5个属性需要配置:

1.packet_type:包类型，必须是一个唯一的正数
2.name:包名字，毕竟数字在代码里不是很友好
3.payload:包内容，这说明包类型确定后，包内容也确定了
4.desc:一句注释而已.
5.module:请求处理模块，只对以_req为后缀的包名字的包有用，可以直接映射到处理模块。

2.包内容

数据包的内容我称其为payload，他们全部在proto/protocal.txt中被定义，包内容定义比较复杂，以下是基础类型，你可以通过基础类型组合成自定义类型：

integer：数字类型
float:浮点类型
string:字符串类型，实际上处理为erlang的binary类型
boolean：布尔类型，只占1个字节
short：端整型

有了基本类型，你可以自定义一个用户类型：

pt_user=
name string
sex boolean
===
非常简单，你可以定义一个账号类型，它嵌套了用户类型:

pt_account=
user pt_user
money integer
===
很多情况下我们需要数组，你可以这样定义一个含有数组的类型：


用户累类型可以组合成新的用户类型
account_i=
pt pt_account
cls integer
===
*/

package packet

import (
	"errors"
	"math"
)

type PacketBz struct {
	Packet
}

//下面是Bz包得各种处理.
func (p *PacketBz) ReadBool() (ret bool, err error) {
	b, _err := p.ReadByte()

	if b != byte(1) {
		return false, _err
	}
	return true, _err
}

func (p *PacketBz) ReadByte() (ret byte, err error) {
	if p.Pos >= len(p.Data) {
		err = errors.New("read byte failed")
		return
	}

	ret = p.Data[p.Pos]
	p.Pos++
	return
}

func (p *PacketBz) ReadBytes() (ret []byte, err error) {
	if p.Pos+2 > len(p.Data) {
		err = errors.New("read bytes header failed")
		return
	}
	size, _ := p.ReadU16()
	if p.Pos+int(size) > len(p.Data) {
		err = errors.New("read bytes Data failed")
		return
	}

	ret = p.Data[p.Pos : p.Pos+int(size)]
	p.Pos += int(size)
	return
}

func (p *PacketBz) ReadString() (ret string, err error) {
	if p.Pos+2 > len(p.Data) {
		err = errors.New("read string header failed")
		return
	}

	size, _ := p.ReadU16()
	if p.Pos+int(size) > len(p.Data) {
		err = errors.New("read string Data failed")
		return
	}

	bytes := p.Data[p.Pos : p.Pos+int(size)]
	p.Pos += int(size)
	ret = string(bytes)
	return
}

func (p *PacketBz) ReadU16() (ret uint16, err error) {
	if p.Pos+2 > len(p.Data) {
		err = errors.New("read uint16 failed")
		return
	}

	buf := p.Data[p.Pos : p.Pos+2]
	ret = uint16(buf[0])<<8 | uint16(buf[1])
	p.Pos += 2
	return
}

func (p *PacketBz) ReadS16() (ret int16, err error) {
	_ret, _err := p.ReadU16()
	ret = int16(_ret)
	err = _err
	return
}

func (p *PacketBz) ReadU24() (ret uint32, err error) {
	if p.Pos+3 > len(p.Data) {
		err = errors.New("read uint24 failed")
		return
	}

	buf := p.Data[p.Pos : p.Pos+3]
	ret = uint32(buf[0])<<16 | uint32(buf[1])<<8 | uint32(buf[2])
	p.Pos += 3
	return
}

func (p *PacketBz) ReadS24() (ret int32, err error) {
	_ret, _err := p.ReadU24()
	ret = int32(_ret)
	err = _err
	return
}

func (p *PacketBz) ReadU32() (ret uint32, err error) {
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

func (p *PacketBz) ReadS32() (ret int32, err error) {
	_ret, _err := p.ReadU32()
	ret = int32(_ret)
	err = _err
	return
}

func (p *PacketBz) ReadU64() (ret uint64, err error) {
	if p.Pos+8 > len(p.Data) {
		err = errors.New("read uint64 failed")
		return
	}

	ret = 0
	buf := p.Data[p.Pos : p.Pos+8]
	for i, v := range buf {
		ret |= uint64(v) << uint((7-i)*8)
	}
	p.Pos += 8
	return
}

func (p *PacketBz) ReadS64() (ret int64, err error) {
	_ret, _err := p.ReadU64()
	ret = int64(_ret)
	err = _err
	return
}

func (p *PacketBz) ReadFloat32() (ret float32, err error) {
	bits, _err := p.ReadU32()
	if _err != nil {
		return float32(0), _err
	}

	ret = math.Float32frombits(bits)
	if math.IsNaN(float64(ret)) || math.IsInf(float64(ret), 0) {
		return 0, nil
	}

	return ret, nil
}

func (p *PacketBz) ReadFloat64() (ret float64, err error) {
	bits, _err := p.ReadU64()
	if _err != nil {
		return float64(0), _err
	}

	ret = math.Float64frombits(bits)
	if math.IsNaN(ret) || math.IsInf(ret, 0) {
		return 0, nil
	}

	return ret, nil
}

func (p *PacketBz) WriteBool(v bool) {
	if v {
		p.Data = append(p.Data, byte(1))
	} else {
		p.Data = append(p.Data, byte(0))
	}
	p.Len += 1
}

func (p *PacketBz) WriteByte(v byte) (err error) {
	p.Data = append(p.Data, v)
	p.Len += 1
	return
}

func (p *PacketBz) WriteBytes(v []byte) (err error) {
	p.WriteU16(uint16(len(v)))
	p.Data = append(p.Data, v...)
	p.Len += len(v)
	return
}

func (p *PacketBz) WriteString(v string) (err error) {
	bytes := []byte(v)
	p.WriteU16(uint16(len(bytes)))
	p.Data = append(p.Data, bytes...)
	p.Len += len(bytes)
	return
}

func (p *PacketBz) WriteU16(v uint16) (err error) {
	p.Data = append(p.Data, byte(v>>8), byte(v))
	p.Len += 2
	return
}

func (p *PacketBz) WriteS16(v int16) (err error) {
	p.WriteU16(uint16(v))
	return
}

func (p *PacketBz) WriteU24(v uint32) (err error) {
	p.Data = append(p.Data, byte(v>>16), byte(v>>8), byte(v))
	p.Len += 3
	return
}

func (p *PacketBz) WriteU32(v uint32) (err error) {
	p.Data = append(p.Data, byte(v>>24), byte(v>>16),
		byte(v>>8), byte(v))
	p.Len += 4
	return
}

func (p *PacketBz) WriteS32(v int32) (err error) {
	p.WriteU32(uint32(v))
	return
}

func (p *PacketBz) WriteU64(v uint64) (err error) {
	p.Data = append(p.Data, byte(v>>56), byte(v>>48),
		byte(v>>40), byte(v>>32), byte(v>>24),
		byte(v>>16), byte(v>>8), byte(v))

	p.Len += 8
	return
}

func (p *PacketBz) WriteS64(v int64) (err error) {
	p.WriteU64(uint64(v))
	return
}

func (p *PacketBz) WriteFloat32(f float32) (err error) {
	v := math.Float32bits(f)
	p.WriteU32(v)
	return
}

func (p *PacketBz) WriteFloat64(f float64) (err error) {
	v := math.Float64bits(f)
	p.WriteU64(v)
	return
}
