package main

import (
	"net"
	"strings"
)

//用户结构体
type User struct {
	Name   string
	Addr   string
	C      chan string
	conn   net.Conn
	server *Server
}

//初始化用户
func NewUser(conn net.Conn, server *Server) *User {
	userName := "无名小卒"
	userAddr := conn.RemoteAddr().String()
	user := &User{
		Name:   userName,
		Addr:   userAddr,
		C:      make(chan string),
		conn:   conn,
		server: server,
	}
	go user.ListenMessage()
	return user
}

//用户上线
func (u *User) Online() {
	u.server.mapLock.Lock()
	u.server.OnlineMap[u.Name] = u
	u.server.mapLock.Unlock()

	u.server.BroadCast(u, "上线")
}

//用户下线
func (u *User) Offline() {
	u.server.mapLock.Lock()
	delete(u.server.OnlineMap, u.Name)
	u.server.mapLock.Unlock()

	u.server.BroadCast(u, "下线")
}

//给当前user的客户端发送消息
func (u *User) SendMsg(msg string) {
	u.conn.Write([]byte(msg))
}

//处理消息
func (u *User) DoMessage(msg string) {

	if msg == "onlist|/$/" {
		//查询当前在线用户
		u.server.mapLock.Lock()
		var TotalStrings string
		result := []string{""}
		for _, user := range u.server.OnlineMap {
			TotalStrings += "%OL|/$/" + user.Addr + user.Name
			result[0] = TotalStrings
		}
		TotalStrings = result[0]

		u.SendMsg(TotalStrings)
		u.server.mapLock.Unlock()

	} else if len(msg) > 7 && msg[:7] == "rename|" {
		//修改用户名 rename|xxx
		newName := strings.Split(msg, "|")[1]
		//判断名字是否已经存在
		_, ok := u.server.OnlineMap[newName]
		if ok {
			u.SendMsg("用户名已存在\n")
		} else {
			u.server.mapLock.Lock()
			delete(u.server.OnlineMap, u.Name)
			u.server.OnlineMap[newName] = u
			u.server.mapLock.Unlock()

			u.Name = newName
			u.SendMsg("用户名成功修改为:" + newName + "\n")
		}

	} else if len(msg) > 4 && msg[:3] == "to|" {
		//私聊  to|zhangsan|你好

		//获取对方用户名
		remoteName := strings.Split(msg, "|")[1]
		if remoteName == "" {
			u.SendMsg("用户名格式不对\n")
			return
		}
		//获取对方user
		remoteUser, ok := u.server.OnlineMap[remoteName]
		if !ok {
			u.SendMsg("用户不存在\n")
			return
		}
		//获取消息
		msg := strings.Split(msg, "|")[2]
		if msg == "" {
			u.SendMsg("无消息内容，重新发送\n")
		}
		//发送消息
		scMsg := "%SC|/$/" + u.Name + "@你：" + msg
		remoteUser.SendMsg(scMsg)

	} else if len(msg) > 5 && msg[:4] == "all|" {

		msg := strings.Split(msg, "|")[1]
		u.server.BroadCast(u, msg)

	}
}

func (u *User) ListenMessage() {
	for {

		msg := <-u.C
		u.conn.Write([]byte(msg + "\n"))

	}
}
