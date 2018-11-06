package main

import "fmt"

type T struct {
	x int
}

func (t T) Value() { //value receiver
	t.x++
}
func (t *T) Pointer() { //pointer receiver
	t.x++ //Go没有->运算符，编译器会自动把t转成(*t)
}

func main() {
	var t *T = &T{1}

	t.Value()
	fmt.Println(t.x)
	t.Pointer()
	fmt.Println(t.x)
}
