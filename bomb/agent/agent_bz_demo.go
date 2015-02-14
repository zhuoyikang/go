/*  一个啥也不干的agent demo */
package agent

import (
	"fmt"
	"time"
)

type AgentBzDemo struct {
	AgentBz
}


func (gs *AgentBzDemo) Start(session *Session) {
	fmt.Printf("%s\n", "bz session start")
}

func (gs *AgentBzDemo) HandlePkt(session *Session, pkti interface{}) {
	fmt.Printf("%s\n", "bz session handle")
	pkt := pkti.(*BzPacket)
	handler := gs.handlerMap[pkt.Type]
	if handler == nil {
		return
	}
	handler(session, pkt)
	return
}

func (gs *AgentBzDemo) Stop(session *Session) {
	fmt.Printf("%s\n", "bz session stop")
}

//Step3. main
func AgentBzDemoMain() {
	agentBz := &AgentBzDemo{}
	agentBz.handlerMap = MakeBzGsHandler()
	agt := MakeAgent("tcp", "0.0.0.0:8080", agentBz, &HandlerBz{})
	go func() {
		time.Sleep(time.Second * 60)
		fmt.Printf("%s\n", "test stop")
		agt.Stop()
	}()
	agt.Signal()
	agt.Run()
	fmt.Printf("%s\n", "end")
}


// // 玩家登陆包
// func BzUserLoginReq(sess *Session, pkt *BzPacket) {
// 	_, userInfo, _ := BzReadPktUserLoginReq(pkt.Data)
// 	fmt.Printf("%v\n", userInfo)
// }
