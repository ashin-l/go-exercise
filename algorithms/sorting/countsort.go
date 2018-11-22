package main

func countSort(s []int) []int {
	result := make([]int, len(s))
	max := getMax(s)
	count := make([]int, max)
	for _, v := range s {
		count[v]++
	}
	i := 0
	for k, v := range count {
		for v > 0 {
			result[i] = k
			i++
			v--
		}
	}
	return result
}

func getMax(s []int) int {
	max := s[0]
	for _, v := range s {
		if v > max {
			max = v
		}
	}
	return max + 1
}
