package utils

import "strings"

func PrepareStringToLike(s string) string {
	if s == "" {
		return s
	}
	arr := strings.Split(s, " ")
	if len(arr) == 1 {
		return "%" + s + "%"
	}
	return "%" + strings.Join(arr, "%") + "%"
}
