package util

import (
	"math/rand"
	"time"
)

func RandArray(n int) []int {
	rand.Seed(time.Now().Unix())
	arr := make([]int, n)
	for i := 0; i != n; i++ {
		arr[i] = rand.Intn(n * 2)
	}
	return arr
}
