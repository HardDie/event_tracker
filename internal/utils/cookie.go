package utils

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/http"
	"time"

	"github.com/HardDie/event_tracker/internal/logger"
)

func SetSessionCookie(session string, w http.ResponseWriter) {
	cookie := http.Cookie{
		Name:     "session",
		Path:     "/",
		Value:    session,
		HttpOnly: true,
	}
	http.SetCookie(w, &cookie)
}
func DeleteSessionCookie(w http.ResponseWriter) {
	cookie := http.Cookie{
		Name:     "session",
		Path:     "/",
		Value:    "",
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
	}
	http.SetCookie(w, &cookie)
}

func GetCookie(r *http.Request) *http.Cookie {
	cookie, err := r.Cookie("session")
	if err != nil {
		logger.Error.Println("Can't read cookie from request:", err.Error())
		return nil
	}
	return cookie
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
