package agent

const (
	BZ_USERLOGINREQ = 0
)

type PktUserLoginReq struct {
	UserId int32
	UserName string
	BaseArr []int32
}

type PktUserLoginAck struct {
	udid string
	name string
	level int32
}

func MakeBzGsHandler() BzHandlerMAP {
	ProtocalHandler := BzHandlerMAP{
		BZ_USERLOGINREQ: BzUserLoginReq,
	}
	return ProtocalHandler
}

func BzReadPktUserLoginReq(datai []byte) (data []byte, ret *PktUserLoginReq, err error) {
	data = datai
	ret = &PktUserLoginReq{}
	data, ret.UserId, err = BzReadint32(data)
 	if err != nil {
 		return
 	}
	data, ret.UserName, err = BzReadstring(data)
 	if err != nil {
 		return
 	}
	var BaseArr_v int32
	data, BaseArr_size, err := BzReaduint16(data)
	for i := 0; i < int(BaseArr_size); i++ {
		data, BaseArr_v, err = BzReadint32(data)
	 	if err != nil {
	 		return
	 	}
		ret.BaseArr = append(ret.BaseArr, BaseArr_v)
	}
 	if err != nil {
 		return
 	}
	return
}
func BzWritePktUserLoginReq(datai []byte, ret *PktUserLoginReq) (data []byte, err error) {
	data = datai
	data, err = BzWriteint32(data, ret.UserId)
	data, err = BzWritestring(data, ret.UserName)
	data, err = BzWriteuint16(data, uint16(len(ret.BaseArr)))
	for _, BaseArr_v := range ret.BaseArr {
		data, err = BzWriteint32(data, BaseArr_v)
	}
	return
}
func BzReadPktUserLoginAck(datai []byte) (data []byte, ret *PktUserLoginAck, err error) {
	data = datai
	ret = &PktUserLoginAck{}
	data, ret.udid, err = BzReadstring(data)
 	if err != nil {
 		return
 	}
	data, ret.name, err = BzReadstring(data)
 	if err != nil {
 		return
 	}
	data, ret.level, err = BzReadint32(data)
 	if err != nil {
 		return
 	}
	return
}
func BzWritePktUserLoginAck(datai []byte, ret *PktUserLoginAck) (data []byte, err error) {
	data = datai
	data, err = BzWritestring(data, ret.udid)
	data, err = BzWritestring(data, ret.name)
	data, err = BzWriteint32(data, ret.level)
	return
}
