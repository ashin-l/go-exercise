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
	_, err = c.Do("hset", "book", "a", 100)
	if err != nil {
		fmt.Println("hset error: ", err)
		return
	}
	a, err := redis.Int(c.Do("hget", "book", "a"))
	fmt.Println("a: ", a)
}
