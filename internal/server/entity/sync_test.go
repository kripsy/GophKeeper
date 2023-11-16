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

func TestRemoveClientSync(t *testing.T) {
	// Инициализируем SyncStatus и добавляем тестового пользователя
	ss := entity.NewSyncStatus()
	userID := 1
	syncID := uuid.New()
	//nolint:errcheck
	ss.AddSync(userID, syncID) // Предполагаем, что AddSync работает корректно

	tests := []struct {
		name       string
		userID     int
		syncID     uuid.UUID
		setup      func() // Дополнительная настройка перед тестом
		wantErr    error
		wantExists bool
	}{
		{
			name:    "Test Case 1 - Remove Existing Sync",
			userID:  userID,
			syncID:  syncID,
			setup:   func() {}, // Нет необходимости в дополнительной настройке
			wantErr: nil,
		},
		{
			name:    "Test Case 2 - Remove Non-Existing User",
			userID:  2, // Несуществующий userID
			syncID:  syncID,
			setup:   func() {},
			wantErr: models.NewSyncError(models.ErrSyncNotFoundEnum),
		},
		{
			name:   "Test Case 3 - Remove with Incorrect SyncID",
			userID: userID,
			syncID: uuid.New(), // Новый UUID, который не совпадает с добавленным
			setup: func() {
				// Добавляем синхронизацию снова, так как она была удалена в первом тесте
				//nolint:errcheck
				ss.AddSync(userID, syncID)
			},
			wantErr: models.NewSyncError(models.ErrSyncNotFoundEnum),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Выполнение дополнительной настройки, если она есть
			tt.setup()

			// Попытка удаления синхронизации
			err := ss.RemoveClientSync(tt.userID, tt.syncID)

			// Проверка ожидаемого результата
			if tt.wantErr != nil {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			// Проверка, что синхронизация больше не существует
			exists, err := ss.IsSyncExists(tt.userID, tt.syncID)
			require.NoError(t, err)
			require.Equal(t, tt.wantExists, exists)
		})
	}
}

func TestIsSyncExists(t *testing.T) {
	ss := entity.NewSyncStatus()
	existingUserID := 1
	existingSyncID := uuid.New()
	nonExistingUserID := 2
	wrongSyncID := uuid.New()

	// Добавляем синхронизацию для существующего пользователя
	//nolint:errcheck
	ss.AddSync(existingUserID, existingSyncID) // Предполагаем, что AddSync работает корректно

	tests := []struct {
		name       string
		userID     int
		syncID     uuid.UUID
		wantExists bool
		wantErr    error
	}{
		{
			name:       "Test Case 1 - Sync Exists",
			userID:     existingUserID,
			syncID:     existingSyncID,
			wantExists: true,
			wantErr:    nil,
		},
		{
			name:       "Test Case 2 - User Does Not Exist",
			userID:     nonExistingUserID,
			syncID:     existingSyncID,
			wantExists: false,
			wantErr:    nil,
		},
		{
			name:       "Test Case 3 - SyncID is Wrong",
			userID:     existingUserID,
			syncID:     wrongSyncID,
			wantExists: false,
			wantErr:    nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			exists, err := ss.IsSyncExists(tt.userID, tt.syncID)

			require.NoError(t, err, "IsSyncExists should not return an error")
			require.Equal(t, tt.wantExists, exists, "IsSyncExists returned unexpected result")
		})
	}
}

func TestAddSync(t *testing.T) {
	ss := entity.NewSyncStatus()
	newUserID := 1
	existingUserID := 2
	newSyncID := uuid.New()
	existingSyncID := uuid.New()

	// Предварительно добавляем синхронизацию для существующего пользователя
	//nolint:errcheck
	ss.AddSync(existingUserID, existingSyncID) // Предполагаем, что AddSync работает корректно

	tests := []struct {
		name      string
		userID    int
		syncID    uuid.UUID
		wantAdded bool
		wantErr   error
	}{
		{
			name:      "Test Case 1 - Add New Sync",
			userID:    newUserID,
			syncID:    newSyncID,
			wantAdded: true,
			wantErr:   nil,
		},
		{
			name:      "Test Case 2 - Add Sync for User with Existing Sync",
			userID:    existingUserID,
			syncID:    uuid.New(), // новый syncID
			wantAdded: false,
			wantErr:   models.NewSyncError(models.ErrUserSyncExistsEnum),
		},
		{
			name:      "Test Case 3 - Add Existing Sync for User",
			userID:    existingUserID,
			syncID:    existingSyncID,
			wantAdded: false,
			wantErr:   models.NewSyncError(models.ErrSyncExistsEnum),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			added, err := ss.AddSync(tt.userID, tt.syncID)

			if tt.wantErr != nil {
				require.Error(t, err)
				require.EqualError(t, tt.wantErr, err.Error())
			} else {
				require.NoError(t, err)
			}
			require.Equal(t, tt.wantAdded, added)
		})
	}
}
