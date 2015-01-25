package packet

import (
	"fmt"
	"testing"
)

type pktUserInfo struct {
	user_id   int32
	user_name string
	//base_a []int32
}

// 从裸包中转为逻辑包，基础类型要作强制转换.
func (lpkt *pktUserInfo) UnPack(pkt *Packet) (err error) {
	err = ((*BzSint32)(&lpkt.user_id)).UnPack(pkt)
	if err != nil {
		return err
	}
	err = ((*BzString)(&lpkt.user_name)).UnPack(pkt)
	if err != nil {
		return err
	}

	// size := BzUint16(0)
	// err = size.UnPack(pkt)
	// if err != nil {
	//	return err
	// }

	return
}

//将PacketI转为为Packet.
func (lpkt *pktUserInfo) Pack(pkt *Packet) (err error) {
	err = ((*BzSint32)(&lpkt.user_id)).Pack(pkt)
	if err != nil {
		return err
	}
	err = ((*BzString)(&lpkt.user_name)).Pack(pkt)
	if err != nil {
		return err
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

//测试基础类型是否打解包正常.
func TestBzString(t *testing.T) {
	pkt := &Packet{Data: make([]byte, 0, 64)}
	pkt.Reset()
	bz_str1 := BzString("this is a test")
	err := bz_str1.Pack(pkt)

	if err != nil {
		t.Errorf("Pack string error")
	}

	//fmt.Printf("bz_str:%v\n", pkt)

	bz_str2 := BzString("")
	pkt.Reset()
	err = bz_str2.UnPack(pkt)
	if err != nil {
		t.Errorf("Pack string error %v", err)
	}

	if bz_str1 != bz_str2 {
		t.Errorf("Pack string error %v", err)
	}
}

//测试基础类型是否打解包正常.
func TestBzSint32(t *testing.T) {
	pkt := &Packet{Data: make([]byte, 0, 64)}
	pkt.Reset()
	bz_int32_1 := BzSint32(13232323)
	err := bz_int32_1.Pack(pkt)

	if err != nil {
		t.Errorf("Pack string error")
	}

	fmt.Printf("bz_int32:%v\n", pkt)

	bz_int32_2 := BzSint32(0)
	pkt.Reset()
	err = bz_int32_2.UnPack(pkt)
	if err != nil {
		t.Errorf("Pack string error %v", err)
	}

	if bz_int32_1 != bz_int32_2 {
		t.Errorf("Pack string error %v", err)
	}
}
