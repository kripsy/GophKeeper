package entity_test

import (
	"testing"

	"github.com/kripsy/GophKeeper/internal/server/entity"
	"github.com/stretchr/testify/require"
)

func TestInitNewUser(t *testing.T) {
	tests := []struct {
		name     string
		username string
		password string
		want     entity.User
		wantErr  bool
	}{
		{
			name:     "Test Case 1 - Initialize User",
			username: "testuser",
			password: "testpassword",
			want: entity.User{
				Username: "testuser",
				Password: "testpassword",
			},
			wantErr: false,
		},
		{
			name:     "Test Case 2 - Initialize User without Password",
			username: "testuser2",
			password: "",
			want:     entity.User{},
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := entity.InitNewUser(tt.username, tt.password)

			if !tt.wantErr {
				require.Equal(t, tt.want, got)
				require.NoError(t, err)
			} else {
				require.Error(t, err)
			}
		})
	}
}
