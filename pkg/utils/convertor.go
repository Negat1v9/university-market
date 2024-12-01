package utils

import (
	"errors"
	"strconv"
)

// Info: converts a string to a number, if the string is not a number, returns an error.
// If the string is empty, returns 0
func ConvertStringToInt64(s string) (int64, error) {
	if s == "" {
		return 0, nil
	}
	num, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0, errors.New("value is not a number")
	}
	return num, nil
}

func ConvertInt64ToString(n int64) string {
	return strconv.FormatInt(n, 10)
}

// Return 0 if string is empty or if s is not uint8 type
func ConvertStringToUint8(s string) uint8 {
	if s == "" {
		return 0
	}

	num, err := strconv.ParseUint(s, 10, 8)
	if err != nil {
		return 0
	}

	return uint8(num)
}
