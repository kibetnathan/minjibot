package auth

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"
)

const (
	// SessionCookieName is the name of the HTTP-only session cookie.
	SessionCookieName = "minji_session"
	// sessionCookieMaxAge is how long a session stays valid (30 days).
	sessionCookieMaxAge = 30 * 24 * time.Hour
)

// Session is the serialized value stored in the signed session cookie.
type Session struct {
	UserID string    `json:"user_id"`
	Expiry time.Time `json:"expiry"`
}

// SessionManager signs and verifies session cookies using an HMAC-SHA256
// message authentication code keyed by a server secret.
type SessionManager struct {
	secret []byte
}

// NewSessionManager returns a SessionManager that signs cookies with the given
// secret. An empty secret falls back to a random per-process key so that
// development servers that omit SESSION_SECRET still function.
func NewSessionManager(secret string) *SessionManager {
	key := []byte(secret)
	if len(key) == 0 {
		key = make([]byte, 32)
		_, _ = rand.Read(key)
	}
	return &SessionManager{secret: key}
}

// Create issues a new session for the given user id and returns it as a cookie.
func (sm *SessionManager) Create(userID string) (*http.Cookie, error) {
	expiry := time.Now().Add(sessionCookieMaxAge)
	raw, err := json.Marshal(Session{UserID: userID, Expiry: expiry})
	if err != nil {
		return nil, err
	}

	// The JSON payload contains characters ('"', '{', ',', ':') that are not
	// valid in an HTTP cookie value, so encode it as base64url first. Without
	// this, net/http strips those bytes from the Set-Cookie header and the
	// cookie never round-trips — logins silently fail.
	value := sm.sign(base64.RawURLEncoding.EncodeToString(raw))

	return &http.Cookie{
		Name:     SessionCookieName,
		Value:    value,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Secure:   true,
		MaxAge:   int(sessionCookieMaxAge.Seconds()),
	}, nil
}

// Verify parses and validates a session cookie, returning the session or an
// error if the signature is invalid or the session has expired.
func (sm *SessionManager) Verify(cookieValue string) (*Session, error) {
	sep := strings.Index(cookieValue, ".")
	if sep < 0 {
		return nil, errors.New("malformed session cookie")
	}
	payload, sig := cookieValue[:sep], cookieValue[sep+1:]

	if !sm.validMAC(payload, sig) {
		return nil, errors.New("invalid session signature")
	}

	raw, err := base64.RawURLEncoding.DecodeString(payload)
	if err != nil {
		return nil, errors.New("invalid session encoding")
	}

	var s Session
	if err := json.Unmarshal(raw, &s); err != nil {
		return nil, errors.New("invalid session payload")
	}
	if time.Now().After(s.Expiry) {
		return nil, errors.New("session expired")
	}
	return &s, nil
}

// Clear returns a cookie that expires the session immediately.
func (sm *SessionManager) Clear() *http.Cookie {
	return &http.Cookie{
		Name:     SessionCookieName,
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Secure:   true,
		MaxAge:   -1,
	}
}

func (sm *SessionManager) sign(payload string) string {
	mac := sm.mac(payload)
	return fmt.Sprintf("%s.%s", payload, mac)
}

func (sm *SessionManager) mac(payload string) string {
	h := hmac.New(sha256.New, sm.secret)
	h.Write([]byte(payload))
	return hex.EncodeToString(h.Sum(nil))
}

func (sm *SessionManager) validMAC(payload, signature string) bool {
	expected := sm.mac(payload)
	macA, errA := hex.DecodeString(signature)
	macB, errB := hex.DecodeString(expected)
	if errA != nil || errB != nil {
		return false
	}
	return hmac.Equal(macA, macB)
}
