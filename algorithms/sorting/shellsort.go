package main

func shellSort(s []int) {
	for gap := len(s) >> 1; gap > 0; gap >>= 1 {
		for i := gap; i < len(s); i++ {
			for j := i; j >= gap && s[j-gap] > s[j]; j -= gap {
				s[j-gap], s[j] = s[j], s[j-gap]
			}
		}
	}
}
