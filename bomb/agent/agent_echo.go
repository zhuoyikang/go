/*
  一个简单地agent实现
  对所有客户端回复数字:2

  网络层自定义 =>

  使用者可以按自己的意愿在这里重写网络层HandleIn和HandleOut处理。
  1.可以选择自定义的序列化方式
  2.也可以选择其他，可以在这里特化使用不同的协议.

*/

package agent

import (
	"fmt"
	//"io"
	//"strconv"
	"time"
	. "github.com/user/bomb/packet"
)

type AgentEcho struct {
}


// 启动一个session
func (gs* AgentEcho) Start(session *Session) {
	fmt.Printf("%s\n", "begin stop")
}

// 当然这里接受到1各数字2，返回给客户端.
func (gs* AgentEcho) HandlePkt(session *Session, pkt Packet) {
	session.Send(pkt)
	return
}

// 停止一个session.
func (gs* AgentEcho) Stop(session *Session) {
	fmt.Printf("%s\n", "begin stop")
}


func AgentEchoMain() {
	agt := MakeAgent("tcp", "0.0.0.0:8080", &AgentEcho{}, &HandlerEcho{})
	go func() {
		// 10s后安全停止所有agent进程.
		time.Sleep(time.Second * 10)
		fmt.Printf("%s\n", "test stop")
		agt.Stop()
	}()
	agt.Signal() //让Agt处理信号。
	agt.Run()
	fmt.Printf("%s\n", "end")
}