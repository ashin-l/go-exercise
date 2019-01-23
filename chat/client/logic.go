package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/ashin-l/go-exercise/chat/proto"
)

func enterMenu() {
	for {
		fmt.Println()
		fmt.Println()
		fmt.Println("------------------------------")
		fmt.Println("---                        ---")
		fmt.Println("--- 1. list online user    ---")
		fmt.Println("--- 2. push message        ---")
		fmt.Println("--- 3. list message        ---")
		fmt.Println("--- 4. exit                ---")
		fmt.Println("---                        ---")
		fmt.Println("------------------------------")
		var menuid int
		fmt.Print("请输入序号： ")
		fmt.Scanln(&menuid)
		switch menuid {
		case 1:
			ListOnlineUser()
		case 2:
			fmt.Println("push message")
			PushMessage()
		case 3:
			ListMessage()
		case 4:
			fmt.Println("Bye!")
			os.Exit(0)
		default:
			fmt.Println("error menuid!")
		}
	}
}

func ListOnlineUser() {
	fmt.Println("\nonlineusers:")
	for i, user := range onlineUser {
		fmt.Printf("%d, %s\n", i, user.NickName)
	}
}

func PushMessage() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("\nSay something:\n")
	content, _, _ := reader.ReadLine()
	proto.WriteMessage(proto.SendMessage, string(content), conn)
}

func ListMessage() {
	fmt.Println("\n按回车返回主菜单！")
	fmt.Println("\nlist messages:")
	exitchan := make(chan byte)
	go func() {
		reader := bufio.NewReader(os.Stdin)
		reader.ReadLine()
		exitchan <- 'e'
	}()
	for {
		select {
		case msg := <-msgchan:
			fmt.Println(msg.UserInfo.NickName, ":", msg.Content)
		case <-exitchan:
			return
		}
	}
}
