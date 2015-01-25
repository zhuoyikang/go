/* 这部分功能将使用awk自动实现生产 */
package agent

import (
	"fmt"
)

type PktUserInfo struct {
	UserId   int32
	UserName string
	BaseArr  []int32
}

const (
	BZ_USER_REGISTER_REQ = 0
)

func MakeBzDemoHandler() BzHandlerMAP {
	ProtocalHandler := BzHandlerMAP{
		BZ_USER_REGISTER_REQ: UserRegisterReq,
	}
	return ProtocalHandler
}

func BzReadPktUserInfo(datai []byte) (data []byte, ret *PktUserInfo, err error) {
	data = datai
	ret = &PktUserInfo{}
	data, ret.UserId, err = BzReads32(data)
	if err != nil {
		return
	}

	data, ret.UserName, err = BzReadstring(data)
	if err != nil {
		return
	}

	var int_v int32
	data, size, err := BzReadu16(data)
	fmt.Printf("xxxx %d\n", data)
	for i := 0; i < int(size); i++ {
		data, int_v, err = BzReads32(data)
		if err != nil {
			return
		}
		ret.BaseArr = append(ret.BaseArr, int_v)
	}
	return
}

func BzWritePktUserInfo(datai []byte, ret *PktUserInfo) (data []byte, err error) {
	data = datai
	data, err = BzWrites32(data, ret.UserId)
	data, err = BzWritestring(data, ret.UserName)

	data, err = BzWriteu16(data, uint16(len(ret.BaseArr)))
	for _, v := range ret.BaseArr {
		data, err = BzWrites32(data, v)
	}
	return
}

// 玩家登陆包
func UserRegisterReq(sess *Session, pkt *BzPacket) {
	_, userInfo, _ := BzReadPktUserInfo(pkt.Data)
	fmt.Printf("%v\n", userInfo)
}
