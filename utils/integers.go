package utils

import "strconv"

func MaxInt(x, y int) int {
	if x > y {
		return x
	}
	return y
}

func MinInt(x, y int) int {
	if x < y {
		return x
	}
	return y
}

func Join(a string, b int) string {
	return a + strconv.Itoa(b)
}
