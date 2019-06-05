package main

import (
	"os"
	"fmt"
	"github.com/ashin-l/go-exercise/thingsboard/stresstest/persist"
	"github.com/ashin-l/go-exercise/thingsboard/stresstest/device"

	"github.com/ashin-l/go-exercise/thingsboard/stresstest/common"
)

func main() {
	err := common.InitConfig("ini", "app.conf")
	if err != nil {
		fmt.Println(err)
		return
	}

	if len(os.Args) == 2 {
		if os.Args[1] == "-syncdb" {
			persist.Syncdb()
			return
		} else if os.Args[1] == "-del" {
			err = persist.InitDB()
			if err != nil {
				fmt.Println(err)
				return
			}
			device.Delall()
			os.Exit(0)
		}
	}

	err = common.InitLogger()
	if err != nil {
		fmt.Println(err)
		return
	}

	err = persist.InitDB()
	if err != nil {
		fmt.Println(err)
		return
	}
	common.Logger.Info("InitDB down ...")

	/*
	err = device.Create(1)
	if err != nil {
		fmt.Println(err)
	}
	*/
	err = device.Run()
	if err != nil {
		fmt.Println(err)
	}
}
