package middleware

import (
	"context"
	"errors"
	"net/http"

	"github.com/HardDie/event_tracker/internal/service"
	"github.com/HardDie/event_tracker/internal/utils"
)

type AuthMiddleware struct {
	authService service.IAuth
}

func NewAuthMiddleware(authService service.IAuth) *AuthMiddleware {
	return &AuthMiddleware{
		authService: authService,
	}
}
func (m *AuthMiddleware) RequestMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie := utils.GetCookie(r)

		// If we got no cookie
		if cookie == nil || cookie.Value == "" {
			http.Error(w, "Invalid session", http.StatusBadRequest)
			return
		}

		// Validate if cookie is active
		ctx := r.Context()
		session, err := m.authService.ValidateCookie(ctx, cookie.Value)
		if err != nil || session == nil {
			if errors.Is(err, service.ErrorSessionHasExpired) {
				http.Error(w, "Session has expired", http.StatusUnauthorized)
			} else {
				http.Error(w, "Invalid session", http.StatusUnauthorized)
			}
			return
		}

		ctx = context.WithValue(ctx, "userID", session.UserID)
		ctx = context.WithValue(ctx, "session", session)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
