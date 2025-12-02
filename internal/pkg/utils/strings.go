package utils

import (
	"strconv"
	"strings"
)

func JoinInt(slice []int, sep string) string {
	if len(slice) == 0 {
		return ""
	}

	result := strconv.Itoa(slice[0])
	for _, v := range slice[1:] {
		result += sep + strconv.Itoa(v)
	}
	return result
}

func HideText(s string, visibleChars int) string {
	if len(s) <= visibleChars {
		return s
	}
	hidden := strings.Builder{}
	for i := 0; i < len(s)-visibleChars; i++ {
		hidden.WriteString("*")
	}
	return s[:visibleChars] + hidden.String()
}
