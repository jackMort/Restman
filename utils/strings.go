package utils

import (
	"regexp"
	"strings"

	"github.com/muesli/ansi"
)

func Truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n-5] + "[...]"
}

func SplitLines(s string) []string {
	return strings.Split(s, "\n")
}

func GetStartColRow(content string, bgRaw string) (int, int) {

	bg := strings.Split(bgRaw, "\n")
	bgWidth := ansi.PrintableRuneWidth(bg[0])
	bgHeight := len(bg)

	cnt := strings.Split(content, "\n")
	width := ansi.PrintableRuneWidth(cnt[0])
	height := len(cnt)

	if height > bgHeight {
		height = bgHeight
	}
	if width > bgWidth {
		width = bgWidth
	}

	startRow := (bgHeight - height) / 2
	startCol := (bgWidth - width) / 2

	return startCol, startRow
}

// Regex to match ANSI escape codes
var ansiReg = regexp.MustCompile(`\x1b\[[0-9;]*m`)

func RemoveANSI(str string) string {
	return ansiReg.ReplaceAllString(str, "")
}
