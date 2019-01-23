package main

import (
	"fmt"
	"os"
)

func enterMenu() {
	for {
		fmt.Println()
		fmt.Println()
		fmt.Println("------------------------------")
		fmt.Println("---                        ---")
		fmt.Println("--- 1. list online user    ---")
		fmt.Println("--- 2. push message        ---")
		fmt.Println("--- 3. exit                ---")
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
		case 3:
			fmt.Println("Bye!")
			os.Exit(0)
		}
	}
}

func ListOnlineUser() {
	fmt.Println("\nonlineusers:")
	for i, user := range onlineUser {
		fmt.Printf("%d, %s\n", i, user.NickName)
	}
}
