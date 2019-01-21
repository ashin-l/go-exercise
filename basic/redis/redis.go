package main

import (
	"fmt"

	"github.com/gomodule/redigo/redis"
)

func main() {
	c, err := redis.Dial("tcp", "140.143.223.156:6379", redis.DialPassword("Tallqc612$"))
	if err != nil {
		fmt.Println("connect redis failed, ", err)
		return
	}
	defer c.Close()
}
