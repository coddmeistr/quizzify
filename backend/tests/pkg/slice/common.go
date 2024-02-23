package slice

import "golang.org/x/exp/constraints"

func Contains[T comparable](arr []T, item T) bool {
	for _, v := range arr {
		if v == item {
			return true
		}
	}
	return false
}

func MaxInt[T constraints.Integer](arr []T) T {
	if len(arr) == 0 {
		return 0
	}
	m := arr[0]
	for _, v := range arr {
		if v > m {
			m = v
		}
	}
	return m
}
