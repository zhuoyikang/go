/* 这部分功能将使用awk自动实现生产 */
package agent

import (
	"fmt"
)

type PktUserInfo struct {
	UserId   int32
	UserName string
	BaseArr    []int32
}

const (
	BZ_USER_REGISTER_REQ = 0
)


func MakeBzDemoHandler() BzHandlerMAP {
	ProtocalHandler := BzHandlerMAP {
		BZ_USER_REGISTER_REQ: UserRegisterReq,
	}
	return ProtocalHandler
}

func BzReadPktUserInfo(pkt *BzPacket) (ret *PktUserInfo, err error) {
	ret = &PktUserInfo{}
	ret.UserId, err = BzReads32(pkt)
	if err != nil {
		return
	}

	ret.UserName, err = BzReadstring(pkt)
	if err != nil {
		return
	}

	var int_v int32
	size, err := BzReadu16(pkt)
	for i := 0; i < int(size); i++ {
		int_v, err = BzReads32(pkt)
		if err != nil {
			return
		}
		ret.BaseArr = append(ret.BaseArr, int_v)
	}
	return
}

func BzWritePktUserInfo(pkt *BzPacket, ret *PktUserInfo) (err error) {
	err = BzWrites32(pkt, ret.UserId)
	err = BzWritestring(pkt, ret.UserName)

	BzWriteu16(pkt, uint16(len(ret.BaseArr)))
	for _, v := range ret.BaseArr {
		BzWrites32(pkt, v)
	}
	return
}

// 玩家登陆包
func UserRegisterReq(sess *Session, pkt *BzPacket) {
	userInfo, _ := BzReadPktUserInfo(pkt)
	fmt.Printf("%v\n", userInfo)
}
