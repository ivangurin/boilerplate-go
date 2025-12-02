package utils

import "fmt"

func ToText[T any](val *T) string {
	if val != nil {
		return fmt.Sprint(*val)
	}
	return ""
}
