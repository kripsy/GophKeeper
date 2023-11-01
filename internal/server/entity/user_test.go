package entity_test

import (
	"testing"

	"github.com/kripsy/GophKeeper/internal/server/entity"
	"github.com/stretchr/testify/assert"
)

func TestInitNewUser(t *testing.T) {
	tests := []struct {
		name     string
		username string
		password string
		want     entity.User
	}{
		{
			name:     "Test Case 1 - Initialize User",
			username: "testuser",
			password: "testpassword",
			want: entity.User{
				Username: "testuser",
				Password: "testpassword",
			},
		},
		{
			name:     "Test Case 2 - Initialize User without Password",
			username: "testuser2",
			password: "",
			want: entity.User{
				Username: "testuser2",
				Password: "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := entity.InitNewUser(tt.username, tt.password)

			assert.NoError(t, err)

			assert.Equal(t, tt.want, got)
		})
	}
}
