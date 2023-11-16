package utils_test

import (
	"context"
	"testing"

	"github.com/kripsy/GophKeeper/internal/utils"
	"github.com/stretchr/testify/require"
)

func TestFromUser2BucketNameWithAsserts(t *testing.T) {
	require := require.New(t)
	prefix := "ilovesber"
	tests := []struct {
		name     string
		username string
		userID   int
		want     string
		wantErr  bool
	}{
		{
			name:     "Valid username and userID",
			username: "TestUser",
			userID:   123,
			want:     prefix + "testuser123",
			wantErr:  false,
		},
		{
			name:     "Empty username",
			username: "",
			userID:   123,
			want:     "",
			wantErr:  true,
		},
		{
			name:     "Zero userID",
			username: "TestUser",
			userID:   0,
			want:     prefix + "testuser0",
			wantErr:  false,
		},
		{
			name:     "Negative userID",
			username: "TestUser",
			userID:   -1,
			want:     "",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := utils.FromUser2BucketName(context.Background(), tt.username, tt.userID)
			if tt.wantErr {
				require.Error(err, "FromUser2BucketName() should return an error")
			} else {
				require.NoError(err, "FromUser2BucketName() should not return an error")
			}
			require.Equal(tt.want, got, "FromUser2BucketName() returned unexpected bucket name")
		})
	}
}
