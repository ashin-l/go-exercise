package main

import (
	"fmt"
	"time"

	"github.com/gomodule/redigo/redis"
)

func main() {
	c, err := redis.Dial("tcp", "localhost:6379")
	if err != nil {
		fmt.Println("connect error: ", err)
		return
	}
	defer c.Close()
	_, err = c.Do("set", "abc", 3)
	if err != nil {
		fmt.Println("set error: ", err)
	}
	_, err = c.Do("expire", "abc", 1)
	if err != nil {
		fmt.Println("expire error: ", err)
	}
	a, _ := c.Do("get", "abc")
	fmt.Println("a: ", a)
	time.Sleep(2 * time.Second)
	a, _ = c.Do("get", "abc")
	fmt.Println("a: ", a)
}
