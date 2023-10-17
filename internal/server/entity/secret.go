package entity

import (
	"bytes"
	"encoding/json"

	"github.com/go-playground/validator/v10"
	"github.com/kripsy/GophKeeper/internal/utils"
)

type Secret struct {
	ID          int    `json:"id,omitempty"`
	Type        string `json:"type" validate:"required,oneof=text binary card login_password"`
	Data        []byte `json:"data" validate:"required"`
	Meta        string `json:"meta,omitempty"`
	ChunkNum    int    `json:"chunk_num,omitempty"`
	TotalChunks int    `json:"total_chunks,omitempty"`
}

func InitNewSecret(data []byte) (Secret, error) {
	var buf bytes.Buffer
	buf.Write(data)
	decoder := json.NewDecoder(&buf)
	decoder.DisallowUnknownFields()

	s := Secret{}
	err := decoder.Decode(&s)
	if err != nil {
		return Secret{}, err
	}

	validate := validator.New()
	err = validate.Struct(&s)
	if err != nil {
		return Secret{}, err
	}

	err = utils.CheckDuplicateFields(data, &Secret{})
	if err != nil {
		return Secret{}, err
	}

	return s, nil
}
