package utils

import (
	"slices"
	"strconv"
)

func StringToInt(value string) int {
	intValue, err := strconv.Atoi(value)
	if err != nil {
		return 0
	}
	return intValue
}

func Contains(slice []string, str string) bool {
	return slices.Contains(slice, str)
}
