package main

func selectSort(s []int) {
	n := len(s)
	for i := 0; i != n; i++ {
		min := i
		for j := i + 1; j != n; j++ {
			if s[j] < s[min] {
				min = j
			}
		}
		if min != i {
			s[i], s[min] = s[min], s[i]
		}
	}
}
