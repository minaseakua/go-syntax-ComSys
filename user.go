package main

import (
	"fmt"
	"io"
	"net"
)

type User struct {
	Name string
	Addr net.Addr
	C    chan string
	conn net.Conn
	serv *Server
}

func (this *User) listenMessage() {
	for {
		msg := <-this.C
		this.conn.Write([]byte(msg + "\n"))
	}
}

func (this *User) online() {
	this.serv.mapLock.Lock()
	this.serv.OnlineMap[this.Name] = this
	this.serv.mapLock.Unlock()
	//BroadCast
	msg := this.Name + "已上线"
	this.serv.BroadCast(this, msg)
}

func (this *User) offline() {
	//BroadCast
	msg := this.Name + "已下线"
	this.serv.BroadCast(this, msg)
	this.serv.mapLock.Lock()
	delete(this.serv.OnlineMap, this.Name)
	this.serv.mapLock.Unlock()
	return
}

func (this *User) sendMessage(msg string) {
	this.serv.BroadCast(this, msg)
}

func (this *User) getMessageFromNet() {
	buf := make([]byte, 1024)
	for {
		n, err := this.conn.Read(buf)
		if n == 0 {
			this.offline()
			return
		}
		if err != nil && err != io.EOF {
			fmt.Println("User", this.Name, "Conn Read Error:", err)
			return
		}
		msg := string(buf[:n-1])
		if msg == "who" {
			this.serv.mapLock.Lock()
			for _, user := range this.serv.OnlineMap {
				onlineMsg := "[" + user.Name + "]" + "在线\n"
				this.C <- onlineMsg
			}
			this.serv.mapLock.Unlock()
		} else {
			this.sendMessage(msg)
		}
	}
}

func NewUser(conn net.Conn, serv *Server) *User {
	u := &User{
		Name: conn.RemoteAddr().String(),
		Addr: conn.RemoteAddr(),
		C:    make(chan string),
		conn: conn,
		serv: serv,
	}
	go u.getMessageFromNet()
	go u.listenMessage()
	u.online()

	return u
}
