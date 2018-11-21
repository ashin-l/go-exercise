package main

import (
	"fmt"

	"github.com/ashin-l/go-exercise/algorithms"
)

func main() {
	arr := util.RandArray(7)
	fmt.Println(arr)
	//bubbleSort(arr)
	//selectSort(arr)
	//arr = mergeSort(arr)
	//cocktailSort(arr)
	//gnomeSort(arr)
	arr = quickSort(arr)
	fmt.Println(arr)
}
