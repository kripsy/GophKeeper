package infrastructure

import "github.com/kripsy/GophKeeper/internal/server/entity"

type BlockRepository struct {
	syncStatus *entity.SyncStatus
}
