package httputil

import (
	"encoding/json"
	"io"
)

func UnmarshalJSONBody(r io.Reader, dest any) error {
	bytes, err := io.ReadAll(r)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(bytes, dest); err != nil {
		return err
	}

	return nil
}
