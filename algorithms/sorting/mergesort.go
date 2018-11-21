package main

func mergeSort(s []int) []int {
	if len(s) <= 1 {
		return s
	}
	middle := len(s) / 2
	l := mergeSort(s[:middle])
	r := mergeSort(s[middle:])

	result := merge(l, r)
	return result
}

func merge(l, r []int) []int {
	result := make([]int, 0, len(l)+len(r))
	for len(l) > 0 || len(r) > 0 {
		if len(l) == 0 {
			return append(result, r...)
		}
		if len(r) == 0 {
			return append(result, l...)
		}
		if l[0] <= r[0] {
			result = append(result, l[0])
			l = l[1:]
		} else {
			result = append(result, r[0])
			r = r[1:]
		}
	}
	return result
}
