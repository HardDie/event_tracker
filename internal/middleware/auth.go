package middleware

import (
	"context"
	"errors"
	"net/http"

	"github.com/HardDie/event_tracker/internal/errs"
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
		bearer := utils.GetBearer(r)

		// If we got no session
		if bearer == "" {
			http.Error(w, "Invalid session token", http.StatusBadRequest)
			return
		}

		// Validate if session is active
		ctx := r.Context()
		session, err := m.authService.ValidateCookie(ctx, bearer)
		if err != nil {
			if errors.Is(err, errs.SessionInvalid) {
				http.Error(w, err.Error(), http.StatusUnauthorized)
			} else {
				http.Error(w, "Internal error", http.StatusInternalServerError)
			}
			return
		}

		ctx = context.WithValue(ctx, "userID", session.UserID)
		ctx = context.WithValue(ctx, "session", session)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
