package main

func cocktailSort(s []int) {
	l, r := 0, len(s)-1
	for l < r {
		for i := l; i != r; i++ {
			if s[i] > s[i+1] {
				s[i], s[i+1] = s[i+1], s[i]
			}
		}
		r--
		for i := r; i != l; i-- {
			if s[i] < s[i-1] {
				s[1], s[i-1] = s[i-1], s[i]
			}
		}
		l++
	}
}
