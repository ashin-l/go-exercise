package main

import (
	"fmt"

	"github.com/gomodule/redigo/redis"
)

func main() {
	c, err := redis.Dial("tcp", "localhost:6379")
	if err != nil {
		fmt.Println("connect redis error: ", err)
		return
	}
	defer c.Close()
	_, err = c.Do("set", "a", 33)
	if err != nil {
		fmt.Println("set error: ", err)
		return
	}
	a, err := redis.Int(c.Do("get", "a"))
	if err != nil {
		fmt.Println("get error: ", err)
		return
	}
	fmt.Println("a: ", a)
}
