package main

import (
	"fmt"
	"os"

	"github.com/ashin-l/go-exercise/chat/proto"
)

func register() {
	fmt.Println("欢迎注册！")
	fmt.Printf("请输入昵称: ")
	fmt.Scanln(&self.NickName)
	fmt.Printf("请输入密码: ")
	fmt.Scanln(&self.Password)
	fmt.Printf("请输入性别: ")
	fmt.Scanln(&self.Sex)
	data := proto.RegisterReqData{self}
	err := proto.WriteMessage(proto.UserRegisterReq, data, conn)
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
}
