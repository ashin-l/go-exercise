package main

import (
	"fmt"

	"github.com/gomodule/redigo/redis"
)

var pool *redis.Pool

func init() {
	pool = &redis.Pool{
		MaxIdle:     16,
		MaxActive:   0,
		IdleTimeout: 300,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", "localhost:6379")
		},
	}
}

func main() {
	c := pool.Get()
	defer c.Close()

	_, err := c.Do("set", "a", 100)
	if err != nil {
		fmt.Println("set error: ", err)
		return
	}

	a, err := redis.Int(c.Do("get", "a"))
	if err != nil {
		fmt.Println("get error: ", err)
		return
	}

	fmt.Println(a)
	pool.Close()
}
