package main

import (
	"fmt"
	"net"

	"github.com/ashin-l/go-exercise/chat/common"
)

var (
	conn       net.Conn
	self       common.User
	onlineUser map[int]*common.User
)

func init() {
	onlineUser = make(map[int]*common.User)
}

func main() {
	var err error
	conn, err = net.Dial("tcp", "localhost:10000")
	if err != nil {
		fmt.Println("Error dialing", err.Error())
		return
	}

	defer conn.Close()
	var menuid int
	fmt.Println("##############################")
	fmt.Println("1. 登录")
	fmt.Println("2. 注册")
	fmt.Println("##############################")
	fmt.Print("\n请输入序号： ")
	fmt.Scanln(&menuid)
	switch menuid {
	case 1:
		login()
	case 2:
		register()
	}

	go processMessage()
	var nc chan struct{}
	<-nc

	//err = login(conn)
	//if err != nil {
	//	fmt.Println("login failed, err:", err)
	//	return
	//}

}
