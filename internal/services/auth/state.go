package auth

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"net/http"
	"time"
)

const (
	// StateCookieName holds the OAuth anti-CSRF state token between the login
	// redirect and the callback.
	StateCookieName = "minji_oauth_state"
	// stateCookieMaxAge bounds how long a login attempt may stay in flight.
	stateCookieMaxAge = 10 * time.Minute
)

// NewState returns a cryptographically random, URL-safe state token for the
// OAuth flow. The token is opaque; it only needs to be unguessable and to match
// on callback.
func NewState() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

// StateCookie returns a short-lived cookie carrying the state token. SameSite is
// Lax (not Strict) so the cookie is still sent on the top-level GET navigation
// that Discord uses to redirect back to the callback.
func StateCookie(state string) *http.Cookie {
	return &http.Cookie{
		Name:     StateCookieName,
		Value:    state,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Secure:   true,
		MaxAge:   int(stateCookieMaxAge.Seconds()),
	}
}

// ClearStateCookie returns a cookie that expires the state cookie immediately.
// The state is single-use, so the callback clears it once consumed.
func ClearStateCookie() *http.Cookie {
	return &http.Cookie{
		Name:     StateCookieName,
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Secure:   true,
		MaxAge:   -1,
	}
}

// StateMatches reports whether the state returned by Discord matches the one we
// stored, using a constant-time comparison. Empty values never match.
func StateMatches(got, want string) bool {
	if got == "" || want == "" {
		return false
	}
	return subtle.ConstantTimeCompare([]byte(got), []byte(want)) == 1
}
