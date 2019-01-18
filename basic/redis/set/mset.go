package main

import (
	"fmt"

	"github.com/gomodule/redigo/redis"
)

func main() {
	c, err := redis.Dial("tcp", "localhost:6379")
	if err != nil {
		fmt.Println("connect error: ", err)
		return
	}
	defer c.Close()
	_, err = c.Do("mset", "a", 11, "b", 22)
	if err != nil {
		fmt.Println("mset error: ", err)
		return
	}
	r, err := redis.Ints(c.Do("mget", "a", "b"))
	if err != nil {
		fmt.Println("mget error: ", err)
		return
	}
	for _, v := range r {
		fmt.Println(v)
	}
}
