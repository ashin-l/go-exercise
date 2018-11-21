package main

func gnomeSort(s []int) {
	for i := 0; i != len(s); {
		if i == 0 || s[i] >= s[i-1] {
			i++
		} else {
			s[i], s[i-1] = s[i-1], s[i]
			i--
		}
	}
}
