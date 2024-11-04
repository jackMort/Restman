package utils

import "strings"

func Truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n-5] + "[...]"
}

func SplitLines(s string) []string {
	return strings.Split(s, "\n")
}
