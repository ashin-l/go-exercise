package main

import (
	"fmt"
	"time"
)

func main() {
	fmt.Println("input:")
	var data string
	fmt.Scanln(&data)
	fmt.Println(data)
	fmt.Println(time.Now())
}
