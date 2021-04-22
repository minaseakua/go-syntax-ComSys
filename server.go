package main

import (
	"fmt"
	"net"
)

type Server struct {
	Ip   string
	Port int
}

func NewServer(ip string, port int) (s *Server) {
	s = &Server{
		Ip:   ip,
		Port: port,
	}
	return
}

func (this *Server) Handler(conn net.Conn) {
	//handler
	fmt.Println("链接建立成功")
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
		//go handler
		go this.Handler(conn)

	}

}
