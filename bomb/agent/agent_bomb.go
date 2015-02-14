/* 炸弹人 */
package agent

import (
	//"errors"
	"fmt"
	//"sync"
	"time"
)

/*------------------------------------------------------------------------------
 实现一个agent模块
 1.任意两个连接建立立刻为其分配一个:地图.
 2.只有1个人连接，则等待.
 3.如果任意玩家在游戏过程中断开，玩家sessin结束，则地图退出；还连接的玩家进入等待队列.
------------------------------------------------------------------------------*/

//房间管理
type RoomMgr struct {
	RoomId   int
	dieChan  chan bool
	msgChan  chan *chan BombMsg
	waitChan *chan BombMsg
}

// 新玩家玩家进入游戏，尝试匹配.
func (mgr *RoomMgr) Join(joinChan *chan BombMsg) {
	mgr.msgChan <- joinChan
	return
}

// 新玩家玩家进入游戏，尝试匹配.
func (mgr *RoomMgr) Worker() {
	var status bool;
	var joinChan *chan BombMsg;

	for {
		if mgr.waitChan == nil {
			joinChan, status = <- mgr.msgChan
		} else {
			select {
			case  <- *mgr.waitChan:
				mgr.waitChan = nil
				fmt.Printf("%s\n", "wait stop")
				continue;
			case joinChan, status = <- mgr.msgChan:
			}
		}
		fmt.Printf("join %v %v\n", joinChan, status)

		if status == false {
			break
		}
		if mgr.waitChan == nil {
			mgr.waitChan = joinChan
		} else if mgr.waitChan == joinChan {
			mgr.waitChan = nil
		} else {
			MakeRoom(*joinChan, *mgr.waitChan)
			mgr.waitChan = nil
		}
	}
}

func MakeRoomMgr() RoomMgr {
	return RoomMgr{dieChan: make(chan bool),
		msgChan: make(chan *chan BombMsg, 1),
	}
}

type AgentBomb struct {
	AgentBz
	RoomMgr
}

type BombMsg struct {
	msg_type int
	data     interface{}
}

type ClientMsg struct {
	pkt     *BzPacket
	handler func(*Session, *BzPacket)
}

type HandlerBomb struct {
	*HandlerBz
	channel chan BombMsg
	dieChan chan bool
	status  int // 玩家状态 0-等待 1-战斗
	room    *Room
}

const (
	MSG_HEART  = 0
	MSG_CLIENT = 1
	MSG_JOIN   = 2 //匹配成功
)

func (handle *HandlerBomb) New() HandlerI {
	h2 := handle.HandlerBz.New()
	ret := &HandlerBomb{HandlerBz: h2.(*HandlerBz),
		channel: make(chan BombMsg, 1),
		dieChan: make(chan bool),
	}
	return ret
}

func (handle *HandlerBomb) HandleClientMsg(session *Session, msg ClientMsg) {
	msg.handler(session, msg.pkt)
}

// 服务routine.
func (handle *HandlerBomb) Service(session *Session) {
	for {
		msg, err := <- handle.channel
		fmt.Printf("get msg %v %v\n", msg, err)
		if err == false {
			close(handle.dieChan)
			return
		}
		switch msg.msg_type {
		case MSG_HEART:
			continue;
		case MSG_JOIN:
			data, _ := BzWriteRoomNtf([]byte{}, &RoomNtf{})
			data = MakePacketData(BZ_ROOMNTF, data)
			session.Send(data)
			fmt.Printf("%s\n", "match ok")
		case MSG_CLIENT:
			// 处理客户端消息
			clientMsg := msg.data.(ClientMsg)
			handle.HandleClientMsg(session, clientMsg)
		}
	}
}

//匹配信息.
type MatchInfo struct {
	p1 *Session
	p2 *Session
}

// 开始一个Session
func (gs *AgentBomb) Start(session *Session) {
	handle := session.PktHandler.(*HandlerBomb)
	gs.RoomMgr.Join(&handle.channel)
	go handle.Service(session)
}

func (gs *AgentBomb) HandlePkt(session *Session, pkti interface{}) {
	fmt.Printf("%s\n", "bz session handle")
	pkt := pkti.(*BzPacket)
	handler := gs.handlerMap[pkt.Type]
	if handler == nil {
		return
	}
	msg := ClientMsg{handler: handler, pkt: pkt}
	handle := session.PktHandler.(*HandlerBomb)
	handle.channel <- BombMsg{msg_type: MSG_CLIENT, data: msg}
	return
}

// 停止一个Session
func (gs *AgentBomb) Stop(session *Session) {
	fmt.Printf("%s\n", "bz session stop")
	handle := session.PktHandler.(*HandlerBomb)
	close(handle.channel)
	<-handle.dieChan
}

//Step3. main
func AgentBombMain() {
	agentBz := &AgentBomb{}
	agentBz.RoomMgr = MakeRoomMgr()
	go agentBz.RoomMgr.Worker()
	agentBz.handlerMap = MakeBzGsHandler()
	agt := MakeAgent("tcp", "0.0.0.0:8080", agentBz, &HandlerBomb{})
	go func() {
		time.Sleep(time.Second * 600000)
		fmt.Printf("%s\n", "test stop")
		agt.Stop()
	}()
	agt.Signal()
	agt.Run()
	fmt.Printf("%s\n", "end")
}

// 玩家登陆包
func BzUserLoginReq(sess *Session, pkt *BzPacket) {
	_, userInfo, _ := BzReadUserLoginReq(pkt.Data)
	fmt.Printf("%v\n", userInfo)
}

// 地图请求
func BzMapReq(sess *Session, pkt *BzPacket) {
	_, mapReq, _ := BzReadMapReq(pkt.Data)
	fmt.Printf("%v\n", mapReq)
}

// 行为设置
func BzBombSetActReq(sess *Session, pkt *BzPacket) {
	_, bombSet, _ := BzReadBombSetAct(pkt.Data)
	fmt.Printf("%v\n", bombSet)
}

/*------------------------------------------------------------------------------
 游戏逻辑
------------------------------------------------------------------------------*/

// 坐标精度
const (
	MAP_COORDINATE_DETAIL = 8  //坐标精度用来跟踪玩家移动
	MAP_WIDTH             = 64 //地图固定宽度
	MAP_HIGH              = 32 //地图固定高度
)

/*
 每个位置字节表示地图上的格子的类型，有以下几种类型:
*/

const (
	MAP_T_SPACE = 0 // 空地
	MAP_T_ROCK  = 1 // 岩石，可以被炸毁
	MAP_T_STEEL = 2 // 钢筋，不可以被炸毁
	MAP_T_BOMB  = 3 //炸弹
)

/*
 每个空格的类型可以相互转换.
 1. 草地上放了一个炸弹 SPACE -> BOMB
 2. 炸弹爆炸(十字范围内+2的岩石转换为空地) BOMB -> SPACE  ROCK -> SPACE

 每个格子的变化都需要发包通知客户端 =>
 {X,Y,T}  -> MapChx 客户端更新UI.
*/

// 建立1个地图
func MakeBombMap(width int, high int) *BombMap {
	return &BombMap{mmap: make([]byte, width*high)}
}

// 设置某个格子
func (mp *BombMap) SetCell(x int, y int, t byte) {
	index := (x-1)*MAP_WIDTH + y
	mp.mmap[index] = t
	return
}

// 获取某个格子属性
func (mp *BombMap) GetCell(x int, y int) byte {
	index := (x-1)*MAP_WIDTH + y
	return mp.mmap[index]
}

// 发送给Room的消息.
type RoomMsg struct {
	bz   *BzPacket
	data interface{}
}

// 房间
type Room struct {
	msg   chan BombMsg // 玩家通过chan和Room通信.
	p1    chan BombMsg
	p2    chan BombMsg
	mmap  *BombMap  //地图
	bombs *BombList //炸弹列表
}

// 建立1个地图
func MakeBombList() *BombList {
	return &BombList{list: make([]*Bomb, 0)}
}

// 创建房间
func MakeRoom(p1 chan BombMsg, p2 chan BombMsg) *Room {
	room := Room{p1: p2, p2: p2}
	room.msg = make(chan BombMsg)
	room.mmap = MakeBombMap(MAP_WIDTH, MAP_HIGH)
	room.bombs = MakeBombList()
	p1 <- BombMsg{msg_type :MSG_JOIN}
	p2 <- BombMsg{msg_type :MSG_JOIN}
	go room.Worker()
	return &room
}

// room 工作.
func (room *Room) Worker() {
	var msg BombMsg
	var status bool
	for {
		fmt.Printf("%s\n", "room service")
		select {
		case msg, status = <-room.p1:
			fmt.Printf("msg %v status %v\n", msg, status)
		case msg, status = <-room.p2:
			fmt.Printf("msg %v status %v\n", msg, status)
		}
		if status == false {
			fmt.Printf("%s\n", "room service out")
			break
		}
	}
}

/*------------------------------------------------------------------------------
 地图处理
 将地图从txt中加载到内存.
------------------------------------------------------------------------------*/
