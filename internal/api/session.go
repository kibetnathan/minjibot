package api

import (
	authsvc "github.com/kibetnathan/minjibot/internal/services/auth"
	"github.com/labstack/echo/v5"
)

// resolveSession reads and verifies the session cookie on a request, returning
// the signed session when valid. It is the shared auth gate for API handlers.
func resolveSession(c *echo.Context, sess *authsvc.SessionManager) (*authsvc.Session, bool) {
	cookie, err := c.Cookie(authsvc.SessionCookieName)
	if err != nil {
		return nil, false
	}
	s, err := sess.Verify(cookie.Value)
	if err != nil {
		return nil, false
	}
	return s, true
}
