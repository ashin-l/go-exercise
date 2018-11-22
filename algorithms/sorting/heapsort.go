package main

import (
	"fmt"
)

func heapSort(s []int) {
	for i := len(s)/2 - 1; i >= 0; i-- {
		maxHeap(s, i, len(s))
	}
	for i := len(s) - 1; i > 0; i-- {
		s[0], s[i] = s[i], s[0]
		maxHeap(s, 0, i)
	}
}

func maxHeap(s []int, start, end int) {
	parent := start
	child := parent*2 + 1
	for child < end {
		if (child+1) < end && s[child+1] > s[child] {
			child++
		}
		if s[parent] >= s[child] {
			return
		}
		s[parent], s[child] = s[child], s[parent]
		parent = child
		child = parent*2 + 1
		fmt.Println(s)
	}
}
