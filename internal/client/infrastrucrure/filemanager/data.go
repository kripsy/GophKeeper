//nolint:gochecknoglobals
package filemanager

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"

	"github.com/kripsy/GophKeeper/internal/utils"
)

const (
	NoteType = iota
	BasicAuthType
	CardDataType
	FileType
)

// Constants representing names for each data type.
const (
	NameNoteType      = "Note"
	NameBasicAuthType = "Login&Password"
	NameCardDataType  = "BankCard"
	NameFileType      = "File"
)

// Data is an interface implemented by various data types to support encryption, hashing, and displaying.
type Data interface {
	// EncryptedData methods for perform encryption using a provided key.
	EncryptedData(key []byte) ([]byte, error)
	// GetHash methods for calculate a SHA-256 hash.
	GetHash() (string, error)
	// String methods for to provide a human-readable representation.
	String() string
}

// DataTypeTable is a slice containing the names of each data type for display purposes.
var DataTypeTable = []string{NameNoteType, NameBasicAuthType, NameCardDataType, NameFileType}

// GetTypeName returns the name of a data type based on its identifier.
func GetTypeName(dataType int) string {
	if dataType < len(DataTypeTable) {
		return DataTypeTable[dataType]
	}

	return "unknown"
}

// CardData represents data related to a bank card.
type CardData struct {
	Number string `json:"number"`
	Date   string `json:"date"`
	CVV    string `json:"cvv"`
}

// BasicAuth represents data related to login credentials.
type BasicAuth struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

// Note represents textual notes.
type Note struct {
	Text string `json:"text"`
}

func (c CardData) String() string {
	return fmt.Sprintf("Number: %q, Date: %q, CVV: %q", c.Number, c.Date, c.CVV)
}

func (n Note) String() string {
	return fmt.Sprintf("Note : %q", n.Text)
}

func (a BasicAuth) String() string {
	return fmt.Sprintf("Login: %q, Password: %q", a.Login, a.Password)
}

func (c CardData) EncryptedData(key []byte) ([]byte, error) {
	data, err := json.Marshal(c)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	data, err = utils.Encrypt(data, key)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	return data, nil
}

func (n Note) EncryptedData(key []byte) ([]byte, error) {
	data, err := json.Marshal(n)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	data, err = utils.Encrypt(data, key)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	return data, nil
}

func (a BasicAuth) EncryptedData(key []byte) ([]byte, error) {
	data, err := json.Marshal(a)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	data, err = utils.Encrypt(data, key)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	return data, nil
}

func (c CardData) GetHash() (string, error) {
	data, err := json.Marshal(c)
	if err != nil {
		return "", fmt.Errorf("%w", err)
	}

	return fmt.Sprintf("%x", sha256.Sum256(data)), nil
}

func (n Note) GetHash() (string, error) {
	data, err := json.Marshal(n)
	if err != nil {
		return "", fmt.Errorf("%w", err)
	}

	return fmt.Sprintf("%x", sha256.Sum256(data)), nil
}

func (a BasicAuth) GetHash() (string, error) {
	data, err := json.Marshal(a)
	if err != nil {
		return "", fmt.Errorf("%w", err)
	}

	return fmt.Sprintf("%x", sha256.Sum256(data)), nil
}
