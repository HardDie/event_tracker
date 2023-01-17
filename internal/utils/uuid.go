package utils

import (
	"fmt"
	"strings"

	"github.com/google/uuid"
)

func UUIDGenerate() (string, error) {
	uid, err := uuid.NewRandom()
	if err != nil {
		return "", fmt.Errorf("error generating uuid: %w", err)
	}
	return strings.ToLower(uid.String()), nil
}
