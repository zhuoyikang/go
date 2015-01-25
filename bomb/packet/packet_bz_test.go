package packet

import (
	"fmt"
	"testing"
)

type pktUserInfo struct {
	user_id   int32
	user_name string
	base_a    []int32
}

// 从裸包中转为逻辑包，基础类型要作强制转换.
func (lpkt *pktUserInfo) UnPack(pkt *Packet) (err error) {
	lpkt.user_id, err = BzReadS32(pkt)
	if err != nil {
		return err
	}

	lpkt.user_name, err = BzReadString(pkt)
	if err != nil {
		return err
	}

	var int_v int32
	size, err := BzReadU16(pkt)
	for i := 0; i < int(size); i++ {
		int_v, err = BzReadS32(pkt)
		if err != nil {
			return
		}
		lpkt.base_a = append(lpkt.base_a, int_v)
	}

	return
}

//将PacketI转为为Packet.
func (lpkt *pktUserInfo) Pack(pkt *Packet) (err error) {
	BzWriteS32(pkt, lpkt.user_id)
	BzWriteString(pkt, lpkt.user_name)

	BzWriteU16(pkt, uint16(len(lpkt.base_a)))
	for _, value := range lpkt.base_a {
		BzWriteS32(pkt, int32(value))
	}

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
	userInfo1 := pktUserInfo{2, "zyk", []int32{10, 20, 30, 40}}
	err := userInfo1.Pack(pkt)

	if err != nil {
		t.Errorf("Pack userInfo1 error")
	}

	fmt.Printf("pkt: %v\n", pkt)

	userInfo2 := pktUserInfo{}
	pkt.Reset()
	err = userInfo2.UnPack(pkt)

	if err != nil {
		t.Errorf("Pack userInfo1 error %v", err)
	}
	fmt.Printf("user_1: %v\n", userInfo1)
	fmt.Printf("user_2: %v\n", userInfo2)

	// if userInfo1 != userInfo2 {
	//	t.Errorf("Pack And UnPack UserInfo error")
	// }

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

	// if acc1 != acc2 {
	//	t.Errorf("Pack And UnPack Acc error")
	// }
}
