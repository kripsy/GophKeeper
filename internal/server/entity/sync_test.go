package entity_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/kripsy/GophKeeper/internal/models"
	"github.com/kripsy/GophKeeper/internal/server/entity"
	"github.com/stretchr/testify/require"
)

func TestSyncStatus(t *testing.T) {
	tests := []struct {
		name       string
		userID     int
		syncID     uuid.UUID
		wantErr    error
		wantExists bool
	}{
		{
			name:       "Test Case 1 - Add New Sync",
			userID:     1,
			syncID:     uuid.New(),
			wantErr:    nil,
			wantExists: true,
		},
		{
			name:       "Test Case 2 - Add Existing Sync",
			userID:     1,
			syncID:     uuid.New(),
			wantErr:    models.NewSyncError(models.ErrUserSyncExistsEnum),
			wantExists: false,
		},
	}

	ss := entity.NewSyncStatus()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			exists, err := ss.AddSync(tt.userID, tt.syncID)

			if tt.wantErr != nil {
				require.EqualError(t, tt.wantErr, err.Error())
			}
			require.Equal(t, tt.wantExists, exists)

			if tt.wantExists {
				isExists, _ := ss.IsSyncExists(tt.userID, tt.syncID)
				require.True(t, isExists)
			}
		})
	}
}
