package main

import (
	"fmt"
	"io"
	"net"
	"sync"
	"time"
)

//服务端结构体
type Server struct {
	Ip   string
	Port int
	//在线用户列表
	OnlineMap map[string]*User
	mapLock   sync.RWMutex
	//消息广播的Channel
	Message chan string
}

//初始化服务端
func NewServer(ip string, port int) *Server {
	server := &Server{
		Ip:        ip,
		Port:      port,
		OnlineMap: make(map[string]*User),
		Message:   make(chan string),
	}
	return server
}

//业务逻辑
func (s *Server) Handler(conn net.Conn) {

	user := NewUser(conn, s)

	user.Online()

	//监听用户是否活跃
	isLive := make(chan bool)

	go func() {
		buf := make([]byte, 4096)
		for {
			n, error := conn.Read(buf)
			if n == 0 {
				user.Offline()
				return
			}
			if error != nil && error != io.EOF {
				fmt.Println("read error")
			}
			msg := string(buf[:n-1])

			user.DoMessage(msg)

			//表示用户活跃
			isLive <- true
		}
	}()

	for {
		select {
		case <-isLive:
			//当前用户活跃，不做任何时，激活select，重置定时器

		case <-time.After(time.Second * 300):
			//超时，将user强制关闭
			user.SendMsg("连接超时,请重启客户端")
			close(user.C)
			conn.Close()
			return
		}
	}
}

//群发监听端
func (s *Server) ListenMessager() {
	for {
		msg := <-s.Message

		s.mapLock.Lock()
		for _, user := range s.OnlineMap {
			user.C <- msg
		}
		s.mapLock.Unlock()
	}
}

//广播函数
func (s *Server) BroadCast(user *User, msg string) {
	sendMsg := "%BC|/$/" + "广播" + ":" + "[" + user.Addr + "]" + user.Name + ":" + msg
	s.Message <- sendMsg
}

//启动服务端
func (s *Server) Start() {
	listener, error := net.Listen("tcp", fmt.Sprintf("%s:%d", s.Ip, s.Port))
	if error != nil {
		fmt.Println("listener error...")
		return
	}
	defer listener.Close()

	go s.ListenMessager()
	for {
		conn, error := listener.Accept()
		if error != nil {
			fmt.Println("accept error...")
			continue
		}
		go s.Handler(conn)
	}

}

func main() {

	// server := NewServer("101.43.19.232", 82)

	server := NewServer("127.0.0.1", 8888)
	server.Start()
}
