package main

import "fmt"

type Slice []int

func NewSlice() Slice {
	return make(Slice, 0)
}
func (s *Slice) Add(elem int) *Slice {
	*s = append(*s, elem)
	fmt.Print(elem)
	return s
}
func main() {
	s := NewSlice()
	fmt.Printf("%v\n", s)
	defer fmt.Println()
	defer s.Add(1).Add(2).Add(4)
	s.Add(3)
}
