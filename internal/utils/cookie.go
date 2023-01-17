package utils

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"
)

func GetBearer(r *http.Request) string {
	header := r.Header.Get("Authorization")
	return strings.ReplaceAll(header, "Bearer ", "")
}

func GenerateSessionKey() (string, error) {
	sessionLen := 32
	b := make([]byte, sessionLen)
	nRead, err := rand.Read(b)
	if err != nil {
		return "", fmt.Errorf("read random: %w", err)
	}
	if nRead != sessionLen {
		return "", fmt.Errorf("bad length")
	}
	return base64.URLEncoding.EncodeToString(b), nil
}
