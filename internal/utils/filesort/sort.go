package filesort

import (
	"strconv"
	"strings"
)

// Helper function to extract the part number from the object name.
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
