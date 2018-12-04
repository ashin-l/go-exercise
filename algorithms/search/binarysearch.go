package main

func binarySearch(s []int, f int) int {
	start := 0
	end := len(s)
	if start == end {
		return -1
	}
	for start <= end {
		mid := start + (end-start)/2
		if s[mid] == f {
			return mid
		}
		if s[mid] > f {
			end = mid - 1
		} else {
			start = mid + 1
		}
	}
	return -1
}
