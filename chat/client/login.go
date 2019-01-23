package main

import (
	"fmt"
	"os"

	"github.com/ashin-l/go-exercise/chat/proto"
)

func login() {
	fmt.Println("\n请输入id： ")
	fmt.Scanln(&self.Id)
	fmt.Println("请输入密码： ")
	fmt.Scanln(&self.Password)
	var loginData proto.LoginReqData
	loginData.Id = self.Id
	loginData.Password = self.Password
	err := proto.WriteMessage(proto.UserLoginReq, loginData, conn)
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
}
