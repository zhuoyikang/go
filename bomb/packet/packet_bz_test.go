package packet

import (
	"fmt"
	"testing"
)

type pktUserInfo struct {
	user_id   int32
	user_name string
}

// 从裸包中转为逻辑包
func (lpkt *pktUserInfo) UnPack(pkt *Packet) (err error) {
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
func (lpkt *pktUserInfo) Pack(pkt *Packet) error {
	pkt_bz := PacketBz{pkt}
	pkt_bz.WriteS32(lpkt.user_id)
	pkt_bz.WriteString(lpkt.user_name)
	return nil
}

type pktAccount struct {
	user_info1 pktUserInfo
	user_info2 pktUserInfo
}

// 固定的逻辑包对裸包的处理模式：递归。

// 从裸包中转为逻辑包
func (lpkt *pktAccount) UnPack(pkt *Packet) (err error) {
	err = lpkt.user_info1.UnPack(pkt)
	if err != nil {
		return err
	}
	err = lpkt.user_info2.UnPack(pkt)
	if err != nil {
		return err
	}
	return
}

//将PacketI转为为Packet.
func (lpkt *pktAccount) Pack(pkt *Packet) (err error) {
	err = lpkt.user_info1.Pack(pkt)
	if err != nil {
		return err
	}
	err = lpkt.user_info2.Pack(pkt)
	if err != nil {
		return err
	}
	return
}


func TestPktUserInfo(t *testing.T) {
	pkt := &Packet{Data: make([]byte, 0, 64)}
	pkt.Reset()
	userInfo1 := pktUserInfo{2, "zyk"}
	err := userInfo1.Pack(pkt)

	if err != nil {
		t.Errorf("Pack userInfo1 error")
	}

	fmt.Printf("%v\n", pkt)

	userInfo2 := pktUserInfo{}
	pkt.Reset()
	err = userInfo2.UnPack(pkt)

	if err != nil {
		t.Errorf("Pack userInfo1 error %v", err)
	}

	if userInfo1 != userInfo2 {
		t.Errorf("Pack And UnPack UserInfo error")
	}

	pkt.Clear()
	acc1 := pktAccount{userInfo1, userInfo2}
	err = acc1.Pack(pkt)
	if err != nil {
		t.Errorf("Pack userInfo1 error")
	}
	acc2 := pktAccount{}
	pkt.Reset()
	err = acc2.UnPack(pkt)
	if err != nil {
		t.Errorf("Pack userInfo1 error %v", err)
	}

	if acc1 != acc2 {
		t.Errorf("Pack And UnPack Acc error")
	}

}
