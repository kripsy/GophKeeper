package entity

type Secret struct {
	ID     int
	Data   []byte
	UserID int
}

func NewSecret(id int, data []byte, userID int) *Secret {
	return &Secret{
		ID:     id,
		Data:   data,
		UserID: userID,
	}
}
