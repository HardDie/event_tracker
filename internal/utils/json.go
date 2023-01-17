package utils

import (
	"encoding/json"
	"fmt"
	"io"
)

func ParseJsonFromHTTPRequest(r io.ReadCloser, data interface{}) error {
	bytes, err := io.ReadAll(r)
	if err != nil {
		return fmt.Errorf("read data from HTTP request: %w", err)
	}
	if err = r.Close(); err != nil {
		return fmt.Errorf("closing request body: %w", err)
	}
	if err = json.Unmarshal(bytes, data); err != nil {
		return fmt.Errorf("parse json HTTP request: %w", err)
	}
	return nil
}
