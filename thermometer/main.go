package main

import (
	"os"
	_ "taosSql"

	"github.com/ashin-l/go-exercise/thermometer/persist"
)

func main() {
	if len(os.Args) == 2 {
		if os.Args[1] == "-syncdb" {
			persist.Syncdb()
			return
		}
	}

}
