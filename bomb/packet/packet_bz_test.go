package packet

import (
	"testing"
	"fmt"
)

type pktUserInfo struct {
	user_id   int32
	user_name string
}

// 从裸包中转为逻辑包
func (lpkt *pktUserInfo) UnPack(pkt Packet) (err error) {
	pkt_bz := PacketBz{pkt}
	lpkt.user_id, err = pkt_bz.ReadS32()
	if err != nil {
		return
	}
	lpkt.user_name, err = pkt_bz.ReadString()
	if err != nil {
		return
	}
	return
}

//将PacketI转为为Packet.
func (lpkt *pktUserInfo) Pack() (Packet, error) {
	pkt_bz := PacketBz{Packet{Data: make([]byte, 0,64)}}
	pkt_bz.WriteS32(lpkt.user_id)
	pkt_bz.WriteString(lpkt.user_name)
	return pkt_bz.Packet, nil
}

func TestPktUserInfo(t *testing.T) {
	userInfo1 := pktUserInfo{2, "zyk"}
	pkt, err := userInfo1.Pack()

	if err != nil {
		t.Errorf("Pack userInfo1 error")
	}

	fmt.Printf("%v\n", pkt)

	userInfo2 := pktUserInfo{}
	err = userInfo2.UnPack(pkt)

	if err != nil {
		t.Errorf("Pack userInfo1 error %v", err)
	}

	if userInfo1 != userInfo2 {
		t.Errorf("Pack And UnPack error")
	}
}
