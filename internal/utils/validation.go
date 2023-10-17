package utils

import (
	"fmt"

	jsonv2 "github.com/go-json-experiment/json"
)

func CheckDuplicateFields(data []byte, v interface{}) error {
	fmt.Println(string(data))
	err := jsonv2.UnmarshalOptions{}.Unmarshal(jsonv2.DecodeOptions{
		AllowDuplicateNames: false,
	}, data, v)

	if err != nil {
		return err
	}

	return nil
}
