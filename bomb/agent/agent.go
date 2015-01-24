/*
思考：
所有的服务器都有以下功能
1. 监听端口
2. 收到信号后安全退出
3. 连接处理。

所以会有Agent.go来将此部分功能通用.

曾经方案1:
1) 打开并监听网络连接.
2) 收到网络连接后，为连接建立3个单独的goroutine，相互间通过channel通信
   * 一个用于收取数据报
   * 一个用于发送数据报
   * 一个用于业务逻辑处理:该模块通过chan和另外两个routine通信
3 收到退出信号时，需要等待所有的进程退出后再停止进程


当前方案:
给每个连接只分配一个routine，节约内存

good:
1.该routine需要block在conn.recv:
2.该routine接受到conn.close()退出

bad:
1.如果第3方需要发送消息给客户端,conn.Send():需要加锁，因为同一时间玩家routine也可能对其发消息。

*/

package agent

import (
	"fmt"
	"net"
	//"strconv"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"github.com/user/bomb/packet"
)


//
type Session struct {
	key   int
	Conn  net.Conn
	mutex *sync.Mutex
	agent *Agent
	pkt_handler packet.HandlerI
}

// 发送数据必须要加锁
func (session *Session) Send(b []byte) (n int, err error) {
	session.mutex.Lock()
	n, err = session.Conn.Write(b)
	session.mutex.Unlock()
	return
}

// 发送数据必须要加锁
func (session *Session) ResetKey(newkey int) (oldkey int) {
	session.mutex.Lock()
	oldkey = session.key
	session.agent.setSessionKey(oldkey, newkey, session)
	session.mutex.Unlock()
	return
}

type AgentI interface {
	//处理由网络返回的包数据
	HandlePkt(session *Session, pkt packet.Packet)
	Start(session *Session)
	Stop(session *Session)
}

type Agent struct {
	net_cls           string
	net_ipfmt         string
	agent_i           AgentI
	pkt_handler       packet.HandlerI
	listener          net.Listener
	wg                *sync.WaitGroup
	die_chan          chan bool //用来控制所有handle-routine退出.
	signal_ch         chan os.Signal
	session_id        int //用来唯一标示一个session，自增数字.
	//session_key_mutex *sync.Mutex
	session_map_mutex *sync.Mutex
	session_map       map[int]*Session
}


// 工厂.
func MakeAgent(cls string, ipfmt string, agent_i AgentI, pkt packet.HandlerI) Agent {
	agent := Agent{net_cls: cls, net_ipfmt: ipfmt,
		agent_i: agent_i, pkt_handler:pkt}
	agent.wg = &sync.WaitGroup{}
	agent.die_chan = make(chan bool)
	agent.pkt_handler = pkt
	//agent.session_key_mutex = &sync.Mutex{}
	agent.session_map_mutex = &sync.Mutex{}
	agent.session_map = make(map[int]*Session)
	return agent
}

// agent.map_lock(func(){ })
func (agent *Agent)map_lock(f func())  {
	agent.session_map_mutex.Lock()
	f()
	agent.session_map_mutex.Unlock()
}

// 互斥的获取主键，因为Agent可能在一个进程里面有多个Listener.目前没有。
// 自动生成的Key是<0的数字形式字符串.
func (agent *Agent) newSessionKey() (ret int) {
	//agent.session_key_mutex.Lock()
	agent.session_id -= 1
	ret = agent.session_id
	//agent.session_key_mutex.Unlock()
	return
}

// 重新设置key.用于业务逻辑自身处理key.
func (agent *Agent) setSessionKey(oldkey, newkey int, session *Session) {
	agent.map_lock(func(){
		delete(agent.session_map, oldkey)
		agent.session_map[newkey] = session
	})
}

// 删除key和session，用于session退出.
func (agent *Agent) delSessionKey(oldkey int) {
	agent.map_lock(func(){
		delete(agent.session_map, oldkey)
	})
}

// 通过Key查找Session.
func (agent *Agent) GetSessionByKey(key int) (session *Session){
	agent.map_lock(func(){
		session = agent.session_map[key]
	})
	return
}

// 增加key.
func (agent *Agent) addSessionKey(key int, session *Session) {
	agent.session_map_mutex.Lock()
	defer func() {
		agent.session_map_mutex.Unlock()
	}()
	if agent.session_map[key] != nil {
		panic("agent key already exist")
	}
	agent.session_map[key] = session
}

// 对外接口:开始服务器
func (agent *Agent) Run() {
	agent.wg.Add(1)
	go agent.run()
	agent.wg.Wait()
}

// 对外接口:停止服务器
func (agent *Agent) Stop() {
	agent.wg.Done()
	agent.listener.Close()
	close(agent.die_chan)
	for _, v := range agent.session_map {
		agent.agent_i.Stop(v)
		v.Conn.Close()
	}
}

// 让Agent注册信号处理.
func (agent *Agent) Signal() {
	agent.signal_ch = make(chan os.Signal, 1)
	signal.Notify(agent.signal_ch, syscall.SIGHUP, syscall.SIGTERM,
		syscall.SIGINT)
	go func() {
		for {
			msg := <-agent.signal_ch
			fmt.Printf("MSG singal %v\n", msg)
			switch msg {
			case syscall.SIGHUP: //
			case syscall.SIGINT:
				agent.Stop()
			case syscall.SIGTERM: // 关闭agent
				agent.Stop()
			}
		}
	}()
}

// 监听主循环
func (agent *Agent) run() {
	listener, err := net.Listen(agent.net_cls, agent.net_ipfmt)
	if err != nil {
		fmt.Printf("error listening: %s\n", err.Error())
		os.Exit(1)
	}
	defer listener.Close()
	agent.listener = listener
	agent.wg.Add(1)

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("Error Accept %s\n", err.Error())
			agent.wg.Done()
			return
		}
		session_key := agent.newSessionKey()
		handler := agent.pkt_handler.New()
		session := &Session{Conn: conn, mutex: &sync.Mutex{},
			key: session_key, agent: agent, pkt_handler:handler}
		//session.Buffer = make([]byte, PREALLOC_BUFSIZE)
		go agent.handle(session)
	}
}

// 循环处理
// 1.AgentI.Handle每次return一个单独的可用的包，这个包数据内容
func (agent *Agent) handle(session *Session) {
	agent.wg.Add(1)
	agent.addSessionKey(session.key, session)
	defer func() {
		agent.wg.Done()
		agent.delSessionKey(session.key)
		session.Conn.Close()
		fmt.Printf("session safe quit %d\n", session.key)
	}()

	for {
		pkt, err := session.pkt_handler.Read(session.Conn)
		if err != nil {
			fmt.Printf("Error Read %s\n", err.Error())
			return
		}
		agent.agent_i.HandlePkt(session, pkt)
	}
}
