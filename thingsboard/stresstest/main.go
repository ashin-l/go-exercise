package main

import (
	"fmt"

	"github.com/ashin-l/go-exercise/thingsboard/stresstest/common"
)

func main() {
	err := common.InitConfig("ini", "app.conf")
	if err != nil {
		fmt.Println(err)
		return
	}
}
