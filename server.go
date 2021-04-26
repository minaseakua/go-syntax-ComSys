package main

import (
	"fmt"
	"io"
	"net"
	"sync"
)

type Server struct {
	Ip        string
	Port      int
	OnlineMap map[string]*User
	Message   chan string
	mapLock   sync.RWMutex
}

func NewServer(ip string, port int) *Server {
	s := &Server{
		Ip:        ip,
		Port:      port,
		OnlineMap: make(map[string]*User),
		Message:   make(chan string),
	}
	return s
}

func (this *Server) ListenMessage() {
	for {
		msg := <-this.Message

		this.mapLock.Lock()
		for _, cli := range this.OnlineMap {
			cli.C <- msg
		}
		this.mapLock.Unlock()
	}
}

func (this *Server) BroadCast(u *User, msg string) {
	sendMsg := "[" + u.Name + "]:" + msg + "\n"
	//广播消息
	this.Message <- sendMsg
}

func (this *Server) Handler(conn net.Conn) {
	//创建新的User,加入OnlineMap
	u := NewUser(conn)
	this.mapLock.Lock()
	this.OnlineMap[u.Name] = u
	this.mapLock.Unlock()
	//handler
	fmt.Println(u.Name, "链接建立成功")
	//BroadCast
	msg := u.Name + "已上线"
	this.BroadCast(u, msg)
	go func() {
		buf := make([]byte, 1024)
		for {
			n, err := conn.Read(buf)
			if n == 0 {
				msg := u.Name + "已下线"
				this.BroadCast(u, msg)
				return
			}
			if err != nil && err != io.EOF {
				fmt.Println("User", u.Name, "Conn Read Error:", err)
				return
			}
			//提取消息并广播
			msg := string(buf[:n-1])
			this.BroadCast(u, msg)
		}
	}()

}

func (this *Server) Start() {
	//socket listen
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", this.Ip, this.Port))
	if err != nil {
		fmt.Println("net listen error:", err)
		return
	}
	defer listener.Close()
	//启动广播用的ListenMessage
	go this.ListenMessage()
	for {
		//socket accept
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("net accept error:", err)
			continue
		}
		//go handler
		go this.Handler(conn)

	}

}
