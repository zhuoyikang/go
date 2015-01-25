/*
  一个简单地实现:Echo
*/

package agent

import (
	"fmt"
	"io"
	//"strconv"
	"time"
)

//Step 1.实现1个agent:应用程序App
type AgentEcho struct {
}

// 启动一个session
func (gs *AgentEcho) Start(session *Session) {
	fmt.Printf("%s\n", "echo session start")
}

// 将接受到得包，又发送给客户端.
func (gs *AgentEcho) HandlePkt(session *Session, pkti interface{}) {
	fmt.Printf("%s\n", "echo session handle")
	session.Send(pkti)
	return
}

// 停止一个session.
func (gs *AgentEcho) Stop(session *Session) {
	fmt.Printf("%s\n", "echo session stop")
}

//Step 2.实现1个handler:网络包处理
type HandlerEcho struct {
}

//从Reader中获取一个完整的包
func (handle *HandlerEcho) Read(reader io.Reader) (interface{}, error) {
	buffer := make([]byte, 1024)
	_, err := reader.Read(buffer)
	return buffer, err
}

func (handle *HandlerEcho) Write(writer io.Writer, pkti interface{}) error {
	buffer := pkti.([]byte)
	_, err := writer.Write(buffer)
	return err
}

func (handle *HandlerEcho) New() HandlerI {
	return &HandlerEcho{}
}

//Step3:测试一个服务器.
func AgentEchoMain() {
	agt := MakeAgent("tcp", "0.0.0.0:8080", &AgentEcho{}, &HandlerEcho{})
	go func() {
		// 10s后安全停止所有agent进程.
		time.Sleep(time.Second * 10)
		fmt.Printf("%s\n", "test stop")
		agt.Stop()
	}()
	agt.Signal() //让Agt处理信号。Ctrl-C
	agt.Run()
	fmt.Printf("%s\n", "end")
}
