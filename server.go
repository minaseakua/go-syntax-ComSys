package main

import (
	"fmt"
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

func (this *Server) BroadCast(u *User, msg string) {
	sendMsg := "[" + u.Name + "]:" + msg + "\n"
	this.mapLock.Lock()
	for _, cli := range this.OnlineMap {
		cli.C <- sendMsg
	}
	this.mapLock.Unlock()
}

func (this *Server) Start() {
	//socket listen
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", this.Ip, this.Port))
	if err != nil {
		fmt.Println("net listen error:", err)
		return
	}
	defer listener.Close()
	for {
		//socket accept
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("net accept error:", err)
			continue
		}
		NewUser(conn, this)
	}
}
