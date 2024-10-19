package utils

import (
	"fmt"
	"strings"
)

func SliceToMap[T comparable](slice []T) map[T]struct{} {
	m := make(map[T]struct{})
	for _, v := range slice {
		m[v] = struct{}{}
	}
	return m
}
func ValidateTwoDecimalPlaces(amount float64) bool {
	strAmount := fmt.Sprintf("%f", amount)
	parts := strings.Split(strAmount, ".")

	if len(parts) == 2 && len(parts[1]) > 2 {
		decimalPart := parts[1]
		decimalPart = strings.TrimRight(decimalPart, "0")
		if len(decimalPart) > 2 {
			return false
		}
	}
	return true
}



