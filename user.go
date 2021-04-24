package main

import "net"

type User struct {
	Name string
	Addr net.Addr
	C    chan string
	conn net.Conn
}

func (this *User) listenMessage() {
	for {
		msg := <-this.C
		this.conn.Write([]byte(msg + "\n"))
	}
}

func NewUser(conn net.Conn) *User {
	u := &User{
		Name: conn.RemoteAddr().String(),
		Addr: conn.RemoteAddr(),
		C:    make(chan string),
		conn: conn,
	}

	go u.listenMessage()

	return u
}
