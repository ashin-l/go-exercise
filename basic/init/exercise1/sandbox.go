package main

import "fmt"

var _ int64 = s()

func init() {
	fmt.Println("init1 in sandbox.go")
}

func init() {
	fmt.Println("init2 in sandbox.go")
}

func s() int64 {
	fmt.Println("calling s() in sandbox.go")
	return 1
}

func main() {
	fmt.Println("main")
}
