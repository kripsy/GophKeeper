package entity

import (
	"errors"
	"fmt"
	"sync"

	"github.com/google/uuid"
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
	fmt.Println("start lock r")
	val, ok := ss.syncClients[userID]
	fmt.Println("try unlock r RemoveClientSync")
	ss.rwMutex.RUnlock()

	if ok {
		if val == syncID {
			fmt.Println("sync already exist for user")
			return false, errors.New("This sync for this user already exists")
		} else {
			fmt.Println("this sync already exists")
			return false, errors.New("Sync for this user already exists")
		}
	}

	fmt.Println("try lock for w AddSync")
	ss.rwMutex.Lock()
	fmt.Println("locked for w AddSync")

	ss.syncClients[userID] = syncID
	fmt.Println("try unlock for w AddSync")
	ss.rwMutex.Unlock()
	fmt.Println("unlock for w AddSync")
	fmt.Println("add val for sync AddSync")

	return true, nil
}

func (ss *SyncStatus) RemoveClientSync(userID int, syncID uuid.UUID) error {
	ss.rwMutex.RLock()
	fmt.Println("start lock r RemoveClientSync")

	val, ok := ss.syncClients[userID]
	fmt.Println("try unlock r RemoveClientSync")
	ss.rwMutex.RUnlock()
	if !ok || val != syncID {
		return errors.New("Sync not found")
	}

	fmt.Println("try lock for w RemoveClientSync")
	ss.rwMutex.Lock()
	fmt.Println("locked for w RemoveClientSync")

	delete(ss.syncClients, userID)
	fmt.Println("try unlock for w RemoveClientSync")
	ss.rwMutex.Unlock()
	fmt.Println("unlock for w RemoveClientSync")

	return nil
}

func (ss *SyncStatus) IsSyncExists(userID int, syncID uuid.UUID) (bool, error) {
	ss.rwMutex.RLock()
	fmt.Println("start lock r IsSyncExists")
	val, ok := ss.syncClients[userID]
	fmt.Println("try unlock r IsSyncExists")
	ss.rwMutex.RUnlock()
	if ok && val == syncID {
		fmt.Println("sync already exist for user")

		return true, nil
	}
	fmt.Println("this sync is not exist")

	return false, nil
}
