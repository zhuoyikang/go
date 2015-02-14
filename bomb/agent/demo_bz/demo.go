package agent

const (
	BZ_USERLOGINREQ = 1
	BZ_USERLOGINACK = 2
	BZ_MAPREQ = 3
	BZ_MAPACK = 4
	BZ_BOMBEXPLODEEVENT = 5
	BZ_BOMBSETACTREQ = 6
	BZ_ROOMNTF = 7
	BZ_POSITIONCHGREQ = 8
)

type UserLoginReq struct {
	UserId int32
	UserName string
	BaseArr []int32
}

type UserLoginAck struct {
	udid string
	name string
	level int32
}

type MapReq struct {
	skip int32
}

type BombMap struct {
	mmap []byte
}

type MapAck struct {
	mmap *BombMap
}

type Bomb struct {
	x int32
	y int32
	r int32
	time int32
}

type BombList struct {
	list []*Bomb
}

type BombExplodeEvent struct {
	bomb *Bomb
}

type BombSetAct struct {
	bomb *Bomb
}

type MapCell struct {
	x int32
	y int32
	t byte
}

type MapCellList struct {
	cells []*MapCell
}

type Point struct {
	x int32
	y int32
}

type RoomNtf struct {
	room_id int32
	self_id int32
	p1_id int32
	p2_id int32
	p1_point *Point
	p2_point *Point
}

type PositionChg struct {
	p *Point
	player_id int32
}

func MakeBzGsHandler() BzHandlerMAP {
	ProtocalHandler := BzHandlerMAP{
		BZ_USERLOGINREQ: BzUserLoginReq,
		BZ_MAPREQ: BzMapReq,
		BZ_BOMBSETACTREQ: BzBombSetActReq,
		BZ_POSITIONCHGREQ: BzPositionChgReq,
	}
	return ProtocalHandler
}

func BzReadUserLoginReq(datai []byte) (data []byte, ret *UserLoginReq, err error) {
	data = datai
	ret = &UserLoginReq{}
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
func BzWriteUserLoginReq(datai []byte, ret *UserLoginReq) (data []byte, err error) {
	data = datai
	data, err = BzWriteint32(data, ret.UserId)
	data, err = BzWritestring(data, ret.UserName)
	data, err = BzWriteuint16(data, uint16(len(ret.BaseArr)))
	for _, BaseArr_v := range ret.BaseArr {
		data, err = BzWriteint32(data, BaseArr_v)
	}
	return
}
func BzReadUserLoginAck(datai []byte) (data []byte, ret *UserLoginAck, err error) {
	data = datai
	ret = &UserLoginAck{}
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
func BzWriteUserLoginAck(datai []byte, ret *UserLoginAck) (data []byte, err error) {
	data = datai
	data, err = BzWritestring(data, ret.udid)
	data, err = BzWritestring(data, ret.name)
	data, err = BzWriteint32(data, ret.level)
	return
}
func BzReadMapReq(datai []byte) (data []byte, ret *MapReq, err error) {
	data = datai
	ret = &MapReq{}
	data, ret.skip, err = BzReadint32(data)
 	if err != nil {
 		return
 	}
	return
}
func BzWriteMapReq(datai []byte, ret *MapReq) (data []byte, err error) {
	data = datai
	data, err = BzWriteint32(data, ret.skip)
	return
}
func BzReadBombMap(datai []byte) (data []byte, ret *BombMap, err error) {
	data = datai
	ret = &BombMap{}
	var mmap_v byte
	data, mmap_size, err := BzReaduint16(data)
	for i := 0; i < int(mmap_size); i++ {
		data, mmap_v, err = BzReadbyte(data)
	 	if err != nil {
	 		return
	 	}
		ret.mmap = append(ret.mmap, mmap_v)
	}
 	if err != nil {
 		return
 	}
	return
}
func BzWriteBombMap(datai []byte, ret *BombMap) (data []byte, err error) {
	data = datai
	data, err = BzWriteuint16(data, uint16(len(ret.mmap)))
	for _, mmap_v := range ret.mmap {
		data, err = BzWritebyte(data, mmap_v)
	}
	return
}
func BzReadMapAck(datai []byte) (data []byte, ret *MapAck, err error) {
	data = datai
	ret = &MapAck{}
	data, ret.mmap, err = BzReadBombMap(data)
 	if err != nil {
 		return
 	}
	return
}
func BzWriteMapAck(datai []byte, ret *MapAck) (data []byte, err error) {
	data = datai
	data, err = BzWriteBombMap(data, ret.mmap)
	return
}
func BzReadBomb(datai []byte) (data []byte, ret *Bomb, err error) {
	data = datai
	ret = &Bomb{}
	data, ret.x, err = BzReadint32(data)
 	if err != nil {
 		return
 	}
	data, ret.y, err = BzReadint32(data)
 	if err != nil {
 		return
 	}
	data, ret.r, err = BzReadint32(data)
 	if err != nil {
 		return
 	}
	data, ret.time, err = BzReadint32(data)
 	if err != nil {
 		return
 	}
	return
}
func BzWriteBomb(datai []byte, ret *Bomb) (data []byte, err error) {
	data = datai
	data, err = BzWriteint32(data, ret.x)
	data, err = BzWriteint32(data, ret.y)
	data, err = BzWriteint32(data, ret.r)
	data, err = BzWriteint32(data, ret.time)
	return
}
func BzReadBombList(datai []byte) (data []byte, ret *BombList, err error) {
	data = datai
	ret = &BombList{}
	var list_v *Bomb
	data, list_size, err := BzReaduint16(data)
	for i := 0; i < int(list_size); i++ {
		data, list_v, err = BzReadBomb(data)
	 	if err != nil {
	 		return
	 	}
		ret.list = append(ret.list, list_v)
	}
 	if err != nil {
 		return
 	}
	return
}
func BzWriteBombList(datai []byte, ret *BombList) (data []byte, err error) {
	data = datai
	data, err = BzWriteuint16(data, uint16(len(ret.list)))
	for _, list_v := range ret.list {
		data, err = BzWriteBomb(data, list_v)
	}
	return
}
func BzReadBombExplodeEvent(datai []byte) (data []byte, ret *BombExplodeEvent, err error) {
	data = datai
	ret = &BombExplodeEvent{}
	data, ret.bomb, err = BzReadBomb(data)
 	if err != nil {
 		return
 	}
	return
}
func BzWriteBombExplodeEvent(datai []byte, ret *BombExplodeEvent) (data []byte, err error) {
	data = datai
	data, err = BzWriteBomb(data, ret.bomb)
	return
}
func BzReadBombSetAct(datai []byte) (data []byte, ret *BombSetAct, err error) {
	data = datai
	ret = &BombSetAct{}
	data, ret.bomb, err = BzReadBomb(data)
 	if err != nil {
 		return
 	}
	return
}
func BzWriteBombSetAct(datai []byte, ret *BombSetAct) (data []byte, err error) {
	data = datai
	data, err = BzWriteBomb(data, ret.bomb)
	return
}
func BzReadMapCell(datai []byte) (data []byte, ret *MapCell, err error) {
	data = datai
	ret = &MapCell{}
	data, ret.x, err = BzReadint32(data)
 	if err != nil {
 		return
 	}
	data, ret.y, err = BzReadint32(data)
 	if err != nil {
 		return
 	}
	data, ret.t, err = BzReadbyte(data)
 	if err != nil {
 		return
 	}
	return
}
func BzWriteMapCell(datai []byte, ret *MapCell) (data []byte, err error) {
	data = datai
	data, err = BzWriteint32(data, ret.x)
	data, err = BzWriteint32(data, ret.y)
	data, err = BzWritebyte(data, ret.t)
	return
}
func BzReadMapCellList(datai []byte) (data []byte, ret *MapCellList, err error) {
	data = datai
	ret = &MapCellList{}
	var cells_v *MapCell
	data, cells_size, err := BzReaduint16(data)
	for i := 0; i < int(cells_size); i++ {
		data, cells_v, err = BzReadMapCell(data)
	 	if err != nil {
	 		return
	 	}
		ret.cells = append(ret.cells, cells_v)
	}
 	if err != nil {
 		return
 	}
	return
}
func BzWriteMapCellList(datai []byte, ret *MapCellList) (data []byte, err error) {
	data = datai
	data, err = BzWriteuint16(data, uint16(len(ret.cells)))
	for _, cells_v := range ret.cells {
		data, err = BzWriteMapCell(data, cells_v)
	}
	return
}
func BzReadPoint(datai []byte) (data []byte, ret *Point, err error) {
	data = datai
	ret = &Point{}
	data, ret.x, err = BzReadint32(data)
 	if err != nil {
 		return
 	}
	data, ret.y, err = BzReadint32(data)
 	if err != nil {
 		return
 	}
	return
}
func BzWritePoint(datai []byte, ret *Point) (data []byte, err error) {
	data = datai
	data, err = BzWriteint32(data, ret.x)
	data, err = BzWriteint32(data, ret.y)
	return
}
func BzReadRoomNtf(datai []byte) (data []byte, ret *RoomNtf, err error) {
	data = datai
	ret = &RoomNtf{}
	data, ret.room_id, err = BzReadint32(data)
 	if err != nil {
 		return
 	}
	data, ret.self_id, err = BzReadint32(data)
 	if err != nil {
 		return
 	}
	data, ret.p1_id, err = BzReadint32(data)
 	if err != nil {
 		return
 	}
	data, ret.p2_id, err = BzReadint32(data)
 	if err != nil {
 		return
 	}
	data, ret.p1_point, err = BzReadPoint(data)
 	if err != nil {
 		return
 	}
	data, ret.p2_point, err = BzReadPoint(data)
 	if err != nil {
 		return
 	}
	return
}
func BzWriteRoomNtf(datai []byte, ret *RoomNtf) (data []byte, err error) {
	data = datai
	data, err = BzWriteint32(data, ret.room_id)
	data, err = BzWriteint32(data, ret.self_id)
	data, err = BzWriteint32(data, ret.p1_id)
	data, err = BzWriteint32(data, ret.p2_id)
	data, err = BzWritePoint(data, ret.p1_point)
	data, err = BzWritePoint(data, ret.p2_point)
	return
}
func BzReadPositionChg(datai []byte) (data []byte, ret *PositionChg, err error) {
	data = datai
	ret = &PositionChg{}
	data, ret.p, err = BzReadPoint(data)
 	if err != nil {
 		return
 	}
	data, ret.player_id, err = BzReadint32(data)
 	if err != nil {
 		return
 	}
	return
}
func BzWritePositionChg(datai []byte, ret *PositionChg) (data []byte, err error) {
	data = datai
	data, err = BzWritePoint(data, ret.p)
	data, err = BzWriteint32(data, ret.player_id)
	return
}
