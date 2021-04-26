package main

import (
	"fmt"
	"io"
	"net"
	"strings"
	"time"
)

type User struct {
	Name    string
	Addr    net.Addr
	C       chan string
	conn    net.Conn
	serv    *Server
	isAlive chan int
}

func (this *User) listenMessage() {
	for {
		select {
		case <-this.isAlive:
		case msg := <-this.C:
			this.conn.Write([]byte(msg + "\n"))
		case <-time.After(10 * time.Second):
			this.logout()
		}
	}
}

func (this *User) online() {
	this.serv.mapLock.Lock()
	this.serv.OnlineMap[this.Name] = this
	this.serv.mapLock.Unlock()
	//BroadCast
	msg := "已上线"
	this.sendMessage(msg, nil)
}

func (this *User) offline() {
	//BroadCast
	msg := "已下线"
	this.sendMessage(msg, nil)
	this.serv.mapLock.Lock()
	delete(this.serv.OnlineMap, this.Name)
	this.serv.mapLock.Unlock()
	return
}

func (this *User) sendMessage(msg string, user *User) {
	this.serv.SendMessage(msg, this, user)
}

func (this *User) logout() {
	this.conn.Write([]byte("您已登出\n"))
	this.conn.Close()
}

func (this *User) getMessageFromNet() {
	buf := make([]byte, 1024)
	for {
		n, err := this.conn.Read(buf)
		this.isAlive <- 1
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
				this.sendMessage(onlineMsg, this)
			}
			this.serv.mapLock.Unlock()
		} else if msg == "exit" {
			this.logout()
		} else if len(msg) > 7 && msg[:7] == "rename:" {
			newName := strings.Split(msg, ":")[1]
			_, ok := this.serv.OnlineMap[newName]
			if ok {
				this.sendMessage("用户名重复", this)
			} else {
				this.serv.mapLock.Lock()
				this.serv.OnlineMap[newName] = this
				delete(this.serv.OnlineMap, this.Name)
				this.serv.mapLock.Unlock()
				this.Name = newName
				this.sendMessage("用户名更改成功", this)
			}
		} else {
			this.sendMessage(msg, nil)
		}
	}
}

func NewUser(conn net.Conn, serv *Server) *User {
	u := &User{
		Name:    conn.RemoteAddr().String(),
		Addr:    conn.RemoteAddr(),
		C:       make(chan string),
		conn:    conn,
		serv:    serv,
		isAlive: make(chan int, 1),
	}
	go u.getMessageFromNet()
	go u.listenMessage()
	u.online()

	return u
}
