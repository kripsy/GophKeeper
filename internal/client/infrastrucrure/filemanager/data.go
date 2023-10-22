package filemanager

import (
	"encoding/json"
	"github.com/kripsy/GophKeeper/internal/utils"
)

//type DataType int

const (
	NoteType = iota
	BasicAuthType
	CardDataType
	FileType
)

const (
	NameNoteType      = "Note"
	NameBasicAuthType = "Login&Password"
	NameCardDataType  = "BankCard"
	NameFileType      = "File"
)

var DataTypeTable = []string{NameNoteType, NameBasicAuthType, NameCardDataType, NameFileType}

func GetTypeName(dataType int) string {
	if dataType < len(DataTypeTable) {
		return DataTypeTable[dataType]
	}

	return "unknown"
}

type CardData struct {
	Number string `json:"number"`
	Date   string `json:"date"`
	CVV    string `json:"cvv"`
}

type BasicAuth struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type Note struct {
	Text string `json:"text"`
}

type File struct {
	Data []byte `json:"Data"`
}

func (f File) EncryptedData(key []byte) ([]byte, error) {
	data, err := json.Marshal(f)
	if err != nil {
		return nil, err
	}

	return utils.Encrypt(data, key)
}

func (c CardData) EncryptedData(key []byte) ([]byte, error) {
	data, err := json.Marshal(c)
	if err != nil {
		return nil, err
	}

	return utils.Encrypt(data, key)
}

func (n Note) EncryptedData(key []byte) ([]byte, error) {
	data, err := json.Marshal(n)
	if err != nil {
		return nil, err
	}

	return utils.Encrypt(data, key)
}

func (a BasicAuth) EncryptedData(key []byte) ([]byte, error) {
	data, err := json.Marshal(a)
	if err != nil {
		return nil, err
	}

	return utils.Encrypt(data, key)
}

type Data interface {
	EncryptedData(key []byte) ([]byte, error)
}
