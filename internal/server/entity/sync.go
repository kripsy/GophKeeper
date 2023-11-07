package entity

import (
	"fmt"
	"sync"

	"github.com/google/uuid"
	"github.com/kripsy/GophKeeper/internal/models"
)

type SyncStatus struct {
	syncClients map[int]uuid.UUID
	rwMutex     *sync.RWMutex
}

func NewSyncStatus() *SyncStatus {
	instance := &SyncStatus{
		syncClients: make(map[int]uuid.UUID),
		rwMutex:     &sync.RWMutex{},
	}

	return instance
}

func (ss *SyncStatus) AddSync(userID int, syncID uuid.UUID) (bool, error) {
	ss.rwMutex.RLock()
	// fmt.Println("start lock r")
	val, ok := ss.syncClients[userID]
	// fmt.Println("try unlock r RemoveClientSync")
	ss.rwMutex.RUnlock()

	if ok {
		if val == syncID {
			// fmt.Println("sync already exist for user")
			return false, fmt.Errorf("%w", models.NewSyncError(models.ErrSyncExistsEnum))
		}

		return false, fmt.Errorf("%w", models.NewSyncError(models.ErrUserSyncExistsEnum))
	}

	// fmt.Println("try lock for w AddSync")
	ss.rwMutex.Lock()
	// fmt.Println("locked for w AddSync")

	ss.syncClients[userID] = syncID
	// fmt.Println("try unlock for w AddSync")
	ss.rwMutex.Unlock()
	// fmt.Println("unlock for w AddSync")
	// fmt.Println("add val for sync AddSync")

	return true, nil
}

func (ss *SyncStatus) RemoveClientSync(userID int, syncID uuid.UUID) error {
	ss.rwMutex.RLock()
	// fmt.Println("start lock r RemoveClientSync")

	val, ok := ss.syncClients[userID]
	// fmt.Println("try unlock r RemoveClientSync")
	ss.rwMutex.RUnlock()
	if !ok || val != syncID {
		return fmt.Errorf("%w", models.NewSyncError(models.ErrSyncNotFoundEnum))
	}

	// fmt.Println("try lock for w RemoveClientSync")
	ss.rwMutex.Lock()
	// fmt.Println("locked for w RemoveClientSync")

	delete(ss.syncClients, userID)
	// fmt.Println("try unlock for w RemoveClientSync")
	ss.rwMutex.Unlock()
	// fmt.Println("unlock for w RemoveClientSync")

	return nil
}

func (ss *SyncStatus) IsSyncExists(userID int, syncID uuid.UUID) (bool, error) {
	ss.rwMutex.RLock()
	// fmt.Println("start lock r IsSyncExists")
	val, ok := ss.syncClients[userID]
	// fmt.Println("try unlock r IsSyncExists")
	ss.rwMutex.RUnlock()
	if ok && val == syncID {
		// fmt.Println("sync already exist for user")

		return true, nil
	}
	// fmt.Println("this sync is not exist")

	return false, nil
}
