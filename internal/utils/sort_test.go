package utils_test

import (
	"testing"

	"github.com/kripsy/GophKeeper/internal/utils"
	"github.com/stretchr/testify/assert"
)

func TestExtractPartNumber(t *testing.T) {
	assert := assert.New(t)

	tests := []struct {
		name         string
		objectName   string
		expectedPart int
	}{
		{
			name:         "Valid part number",
			objectName:   "file-part-1",
			expectedPart: 1,
		},
		{
			name:         "No part number",
			objectName:   "file-part-",
			expectedPart: 0,
		},
		{
			name:         "Invalid part number",
			objectName:   "file-part-abc",
			expectedPart: 0,
		},
		{
			name:         "Missing part number",
			objectName:   "file",
			expectedPart: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			partNumber := utils.ExtractPartNumber(tt.objectName)
			assert.Equal(tt.expectedPart, partNumber, "ExtractPartNumber() did not return the expected part number")
		})
	}
}
