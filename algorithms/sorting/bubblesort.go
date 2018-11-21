package main

func bubbleSort(s []int) {
	for j := len(s) - 1; j > 0; j-- {
		for i := 0; i != j; i++ {
			if s[i] > s[i+1] {
				s[i], s[i+1] = s[i+1], s[i]
			}
		}
	}
}
