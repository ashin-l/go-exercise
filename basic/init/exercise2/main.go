package main

import "fmt"

func init() {
	fmt.Println("init")
}

// init() 函数不可以调用
func main() {
	init() // 编译时报错：undefined: init
}
