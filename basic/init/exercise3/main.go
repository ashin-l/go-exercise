package main

import (
	"fmt"

	_ "github.com/ashin-l/go-exercise/basic/init/exercise3/myinit"
)

//var precomputed1 = [20]float64{}
//var precomputed2 = pcm()

//func init() {
//	fmt.Println("in init")
//	var current float64 = 1
//	precomputed1[0] = current
//	for i := 1; i < len(precomputed1); i++ {
//		precomputed1[i] = precomputed1[i-1] * 1.2
//	}
//
//	for _, v := range precomputed1 {
//		fmt.Printf("%f ", v)
//	}
//
//	fmt.Println()
//	for _, v := range precomputed2 {
//		fmt.Printf("%f ", v)
//	}
//}
//
//func pcm() [20]float64 {
//	fmt.Println("in pcm")
//	var result [20]float64
//	var current float64 = 1
//	result[0] = current
//	for i := 1; i < len(result); i++ {
//		result[i] = result[i-1] * 1.2
//	}
//	return result
//}

func main() {
	fmt.Println("\nin main")
}
