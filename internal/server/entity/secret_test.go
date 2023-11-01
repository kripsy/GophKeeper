package entity_test

import (
	"reflect"
	"testing"

	"github.com/kripsy/GophKeeper/internal/server/entity"
)

func TestNewSecret(t *testing.T) {
	tests := []struct {
		name   string
		id     int
		data   []byte
		userID int
		want   *entity.Secret
	}{
		{
			name:   "Test Case 1",
			id:     1,
			data:   []byte("secretData1"),
			userID: 101,
			want:   &entity.Secret{ID: 1, Data: []byte("secretData1"), UserID: 101},
		},
		{
			name:   "Test Case 2",
			id:     2,
			data:   []byte("secretData2"),
			userID: 102,
			want:   &entity.Secret{ID: 2, Data: []byte("secretData2"), UserID: 102},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := entity.NewSecret(tt.id, tt.data, tt.userID); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewSecret() = %v, want %v", got, tt.want)
			}
		})
	}
}
