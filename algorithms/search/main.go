package main

import (
	"fmt"

	"github.com/ashin-l/go-exercise/algorithms"
)

func main() {
	arr := util.RandArray(13)
	arr = quickSort(arr)
	fmt.Println(arr)
	index := binarySearch(arr, 7)
	if index == -1 {
		fmt.Println("not find 7!")
	} else {
		fmt.Printf("find 7, index: %d\n", index)
	}
}

func quickSort(s []int) []int {
	if len(s) <= 1 {
		return s
	}
	l := make([]int, 0, len(s))
	r := make([]int, 0, len(s))
	m := make([]int, 0, len(s))
	median := s[0]
	for _, v := range s {
		if v < median {
			l = append(l, v)
		} else if v == median {
			m = append(m, v)
		} else {
			r = append(r, v)
		}
	}
	l = quickSort(l)
	r = quickSort(r)
	l = append(l, m...)
	l = append(l, r...)
	return l
}
