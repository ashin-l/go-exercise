package main

import (
	"fmt"
	"time"

	"github.com/hpcloud/tail"
)

func main() {
	filename := "mylog.log"
	t, err := tail.TailFile(filename, tail.Config{
		ReOpen:    true,
		MustExist: false,
		Poll:      true,
		Follow:    true,
		Location:  &tail.SeekInfo{Offset: 0, Whence: 2},
	})
	if err != nil {
		panic(err)
	}

	var msg *tail.Line
	var ok bool
	for {
		msg, ok = <-t.Lines
		if !ok {
			fmt.Println("tail file close reopen, filename:", t.Filename)
			time.Sleep(100 * time.Millisecond)
			continue
		}
		fmt.Println("msg:", msg)
	}
}
