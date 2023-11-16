// Package entity defines core data structures and entities for the GophKeeper application.
// It includes structures for managing secrets, synchronization status, and user information.
package entity

// Secret represents a secret object in the GophKeeper application.
type Secret struct {
	ID     int
	Data   []byte
	UserID int
}

// NewSecret creates a new instance of a Secret.
func NewSecret(id int, data []byte, userID int) *Secret {
	return &Secret{
		ID:     id,
		Data:   data,
		UserID: userID,
	}
}
