package main

import (
	"fmt"
	"os"
	_ "taosSql"

	"github.com/ashin-l/go-exercise/thermometer/common"

	"github.com/ashin-l/go-exercise/thermometer/device"

	"github.com/ashin-l/go-exercise/thermometer/persist"
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
		}
	}

	err = common.InitLogger()
	if err != nil {
		fmt.Println(err)
		return
	}

	persist.InitDB()
	err = device.Run()
	if err != nil {
		fmt.Println(err)
		return
	}

	persist.CloseDB()
}
