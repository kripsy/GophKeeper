package utils

import (
	"strconv"
	"strings"
)

// ExtractPartNumber extracts the part number from a given object name.
//
// The function assumes that the object name follows a specific format,
// with the part number appended after '-part-'.
//
// Parameters:
// - objectName: The name of the object from which the part number is to be extracted.
//
// Returns:
// - An integer representing the extracted part number.
// - Returns 0 if the part number cannot be extracted or parsed.
func ExtractPartNumber(objectName string) int {
	minLen := 2
	parts := strings.Split(objectName, "-part-")
	if len(parts) < minLen {
		return 0
	}
	num, err := strconv.Atoi(parts[1])
	if err != nil {
		return 0
	}

	return num
}
