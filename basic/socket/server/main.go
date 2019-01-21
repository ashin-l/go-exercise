package main

import (
	"fmt"
	"net"
)

func main() {
	fmt.Println("server start...")
	l, err := net.Listen("tcp", "0.0.0.0:9000")
	if err != nil {
		fmt.Println("server error: ", err.Error())
	}
	defer l.Close()
	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("connect error: ", err.Error())
			continue
		}
		go process(conn)
	}
}

func process(conn net.Conn) {
	defer conn.Close()
	for {
		buf := make([]byte, 512)
		_, err := conn.Read(buf)
		if err != nil {
			fmt.Println("read error: ", err.Error())
			return
		}
		fmt.Println("read: ", string(buf))
	}

}
