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

func CheckValidQueryParam(params ...int) (bool, error) {
	for _, param := range params {
		if param < 0 || param > 100 {
			return false, nil
		}
	}
	return true, nil
}

func IsClosed(ch <-chan []byte) bool {
	select {
	case _, ok := <-ch:
		return !ok
	default:
	}
	return false
}
