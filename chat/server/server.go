package main

import (
	"fmt"
	"net"

	"github.com/ashin-l/go-exercise/chat/server/models"
)

func runServer(addr string) (err error) {
	l, err := net.Listen("tcp", addr)
	if err != nil {
		fmt.Println("listen failed!")
		return
	}

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("accept error!")
			continue
		}
		go process(conn)
	}
}

func process(conn net.Conn) {
	defer conn.Close()
	c := models.NewClient(conn)
	err := c.Process()
	if err != nil {
		fmt.Println("parse process failed, ", err)
	}
}
