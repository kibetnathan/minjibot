package auth

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

// encodeSession mirrors Create's payload encoding (base64url of the JSON) so a
// test can construct a signable payload directly.
func encodeSession(t *testing.T, s Session) string {
	t.Helper()
	raw, err := json.Marshal(s)
	if err != nil {
		t.Fatalf("marshal session: %v", err)
	}
	return base64.RawURLEncoding.EncodeToString(raw)
}

// TestSession_RoundTripThroughHTTP is the regression test for the cookie
// encoding bug: the cookie must survive being written to a Set-Cookie header
// and read back from a request, then verify to the original user id.
func TestSession_RoundTripThroughHTTP(t *testing.T) {
	sm := NewSessionManager("a-test-secret")

	cookie, err := sm.Create("user-123")
	if err != nil {
		t.Fatalf("Create: %v", err)
	}

	// Emit via Set-Cookie the way the server does...
	rec := httptest.NewRecorder()
	http.SetCookie(rec, cookie)
	setCookie := rec.Header().Get("Set-Cookie")
	if setCookie == "" {
		t.Fatal("no Set-Cookie header emitted")
	}

	// ...and read it back the way a browser would send it on the next request.
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Cookie", setCookie[:strings.Index(setCookie, ";")])
	got, err := req.Cookie(SessionCookieName)
	if err != nil {
		t.Fatalf("cookie did not round-trip: %v", err)
	}

	s, err := sm.Verify(got.Value)
	if err != nil {
		t.Fatalf("Verify after round-trip: %v", err)
	}
	if s.UserID != "user-123" {
		t.Errorf("user id = %q, want %q", s.UserID, "user-123")
	}
}

func TestSession_CookieValueIsHTTPSafe(t *testing.T) {
	sm := NewSessionManager("secret")
	cookie, err := sm.Create("user-123")
	if err != nil {
		t.Fatalf("Create: %v", err)
	}
	// Cookie values must not contain characters that net/http would strip.
	for _, bad := range []string{"\"", "{", "}", ",", ":", " "} {
		if strings.Contains(cookie.Value, bad) {
			t.Errorf("cookie value contains invalid byte %q: %s", bad, cookie.Value)
		}
	}
}

func TestSession_CreateVerify(t *testing.T) {
	sm := NewSessionManager("secret")
	cookie, err := sm.Create("abc")
	if err != nil {
		t.Fatalf("Create: %v", err)
	}
	s, err := sm.Verify(cookie.Value)
	if err != nil {
		t.Fatalf("Verify: %v", err)
	}
	if s.UserID != "abc" {
		t.Errorf("user id = %q, want abc", s.UserID)
	}
}

func TestSession_TamperedPayloadRejected(t *testing.T) {
	sm := NewSessionManager("secret")
	cookie, _ := sm.Create("abc")

	// Flip a character in the payload portion (before the '.') — the MAC no
	// longer matches.
	sep := strings.Index(cookie.Value, ".")
	if sep < 1 {
		t.Fatalf("unexpected cookie shape: %q", cookie.Value)
	}
	tampered := "X" + cookie.Value[1:]
	if _, err := sm.Verify(tampered); err == nil {
		t.Error("expected tampered payload to be rejected")
	}
}

func TestSession_WrongSecretRejected(t *testing.T) {
	a := NewSessionManager("secret-a")
	b := NewSessionManager("secret-b")
	cookie, _ := a.Create("abc")
	if _, err := b.Verify(cookie.Value); err == nil {
		t.Error("expected a session signed with a different secret to be rejected")
	}
}

func TestSession_ExpiredRejected(t *testing.T) {
	sm := NewSessionManager("secret")
	// Hand-craft an already-expired session using the same signing scheme.
	expired := Session{UserID: "abc", Expiry: time.Now().Add(-time.Hour)}
	value := sm.sign(encodeSession(t, expired))
	if _, err := sm.Verify(value); err == nil {
		t.Error("expected expired session to be rejected")
	}
}

func TestSession_MalformedRejected(t *testing.T) {
	sm := NewSessionManager("secret")
	for _, v := range []string{"", "nodot", ".", "abc.", ".abc", "notbase64!.deadbeef"} {
		if _, err := sm.Verify(v); err == nil {
			t.Errorf("expected malformed cookie %q to be rejected", v)
		}
	}
}
