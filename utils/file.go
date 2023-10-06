package utils

import (
	"encoding/json"

	"github.com/aronlt/toolkit/terror"
	"github.com/aronlt/toolkit/tio"
)

func WriteToFile[T any](data []T, filepath string) error {
	content, err := json.MarshalIndent(data, "", " ")
	if err != nil {
		return terror.Wrap(err, "call MarshalIndent fail")
	}
	_, err = tio.WriteFile(filepath, content, false)
	if err != nil {
		return terror.Wrap(err, "call WriteFile fail")
	}
	return nil
}
