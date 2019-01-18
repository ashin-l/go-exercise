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
	c.Do("del", "book_list")
	_, err = c.Do("lpush", "book_list", "a", "b", 100)
	if err != nil {
		fmt.Println("lpush error: ", err)
		return
	}

	r, err := redis.Strings(c.Do("lrange", "book_list", 0, 100))
	if err != nil {
		fmt.Println("lpop error: ", err)
	}
	for _, v := range r {
		fmt.Println(v)
	}

	str, err := redis.String(c.Do("rpop", "book_list"))
	if err != nil {
		fmt.Println("lpop error: ", err)
	}
	fmt.Println(str)
}
