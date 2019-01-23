package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/ashin-l/go-exercise/chat/common"
	"github.com/ashin-l/go-exercise/chat/proto"
)

func processMessage() {
	for {
		msg, err := proto.ReadMessage(conn)
		if err != nil {
			fmt.Println("read error: ", err)
			os.Exit(0)
		}
		switch msg.Cmd {
		case proto.UserLoginRes:
			loginResp(msg.Data)
		case proto.UserRegisterRes:
			registerResp(msg.Data)
		case proto.NotifyUserStatus:
			updateOnlineUsers(msg.Data)
		}
	}
}

func loginResp(data string) {
	var logindata proto.LoginResData
	err := json.Unmarshal([]byte(data), &logindata)
	if err != nil {
		fmt.Println("unmarshal loginResp error: ", err)
		os.Exit(0)
	}
	if logindata.Error != "" {
		fmt.Println("login failed! ", logindata.Error)
		os.Exit(0)
	} else {
		fmt.Println("login success!")
		for i := range logindata.Users {
			user := logindata.Users[i]
			onlineUser[user.Id] = &user
		}
		go enterMenu()
	}
}

func registerResp(data string) {
	var regdata proto.RegisterResData
	err := json.Unmarshal([]byte(data), &regdata)
	if err != nil {
		fmt.Println("unmarshal registerResp error: ", err)
		os.Exit(0)
	}
	if regdata.Error != "" {
		fmt.Println("register failed! ", regdata.Error)
		os.Exit(0)
	} else {
		self.Id = regdata.Id
		fmt.Println("+++++++++++++++++++++++++++++++++++")
		fmt.Println()
		fmt.Println("register sucess! remember your id: ", self.Id)
		fmt.Println()
		fmt.Println("+++++++++++++++++++++++++++++++++++")
		for i := range regdata.Users {
			user := regdata.Users[i]
			onlineUser[user.Id] = &user
		}
		go enterMenu()
	}
}

func updateOnlineUsers(data string) {
	var nfdata proto.NotifyUserStatusData
	err := json.Unmarshal([]byte(data), &nfdata)
	if err != nil {
		fmt.Println("unmarshal notifydata error: ", err)
		return
	}
	if nfdata.Status == common.UserStatusOffline {
		delete(onlineUser, nfdata.User.Id)
	} else {
		onlineUser[nfdata.User.Id] = &nfdata.User
	}
}
