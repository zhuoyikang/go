/*
  一个简单地agent实现
  对所有客户端回复2
*/

package agent

import (
	"fmt"
	"io"
	"strconv"
	"time"
)

type AgentDemo1 struct {
}

// 阻塞在网络层，收到完整的包后返回包数据.
func (gs AgentDemo1) HandleIn(session *Session) (n interface{}, err error) {
	header := make([]byte, 2)
	n, err = io.ReadFull(session.Conn, header)
	return
}

func (gs AgentDemo1) HandleOut(session *Session, i interface{}) {
	n := i.(int)
	n_str := strconv.Itoa(n)
	session.Send([]byte(n_str))
	return
}

func (gs AgentDemo1) Stop() {
	fmt.Printf("%s\n", "begin stop")
}

func AgentDemo1Main() {
	agt := MakeAgent("tcp", "0.0.0.0:8080", AgentDemo1{})
	go func() {
		time.Sleep(time.Second * 10)
		fmt.Printf("%s\n", "test stop")
		agt.Stop()
	}()
	agt.Signal() //让Agt处理信号。
	agt.Run()
	fmt.Printf("%s\n", "end")
}
