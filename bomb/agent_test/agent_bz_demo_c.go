package main

import (
	"fmt"
	"github.com/user/bomb/agent"
	"net"
	"time"
)


func test1(data []byte) {
	data = append(data, byte(1), byte(2), byte(3))
	fmt.Printf("t1 %v\n", data)
}

func test2() {
	data := []byte{}
	test1(data)
	fmt.Printf("t2 %v\n", data)
}

func main() {
	test2()
	conn, err := net.Dial("tcp", "0.0.0.0:8080")
	if err != nil {
		fmt.Println("连接服务端失败:", err.Error())
		return
	}
	fmt.Println("已连接服务器")
	defer conn.Close()

	userInfo := &agent.PktUserInfo{
		UserId:   23,
		UserName: "good",
		BaseArr:  []int32{1, 2, 4, 5},
	}

	fmt.Printf("%v\n", userInfo)
	data, err := agent.BzWritePktUserInfo([]byte{}, userInfo)
	data = agent.MakePacketData(0, data)

	conn.Write(data)
	time.Sleep(1 * time.Second)
	fmt.Printf("%s %v %d\n", "client quit", data, len(data))
}
