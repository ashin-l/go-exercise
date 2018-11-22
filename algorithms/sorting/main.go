package main

import (
	"fmt"

	"github.com/ashin-l/go-exercise/algorithms"
)

func main() {
	arr := util.RandArray(13)
	fmt.Println(arr)
	//bubbleSort(arr)
	//selectSort(arr)
	//arr = mergeSort(arr)
	//cocktailSort(arr)
	//gnomeSort(arr)
	//arr = quickSort(arr)
	//heapSort(arr)
	//shellSort(arr)
	arr = countSort(arr)
	fmt.Println(arr)
}
