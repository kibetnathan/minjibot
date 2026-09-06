package auth

import (
	"strings"
	"testing"
)

func TestNewState_UniqueAndURLSafe(t *testing.T) {
	seen := map[string]struct{}{}
	for i := 0; i < 100; i++ {
		s, err := NewState()
		if err != nil {
			t.Fatalf("NewState: %v", err)
		}
		if s == "" {
			t.Fatal("NewState returned empty string")
		}
		// base64url alphabet only: A-Za-z0-9-_
		for _, r := range s {
			ok := (r >= 'A' && r <= 'Z') || (r >= 'a' && r <= 'z') ||
				(r >= '0' && r <= '9') || r == '-' || r == '_'
			if !ok {
				t.Fatalf("state contains non-URL-safe char %q: %s", r, s)
			}
		}
		if _, dup := seen[s]; dup {
			t.Fatalf("NewState produced a duplicate: %s", s)
		}
		seen[s] = struct{}{}
	}
}

func TestStateMatches(t *testing.T) {
	cases := []struct {
		got, want string
		match     bool
	}{
		{"abc", "abc", true},
		{"abc", "abd", false},
		{"", "abc", false},
		{"abc", "", false},
		{"", "", false},
		{"abc", "ab", false},
	}
	for _, tc := range cases {
		if got := StateMatches(tc.got, tc.want); got != tc.match {
			t.Errorf("StateMatches(%q, %q) = %v, want %v", tc.got, tc.want, got, tc.match)
		}
	}
}

func TestStateCookie_IsHTTPSafeAndShortLived(t *testing.T) {
	s, _ := NewState()
	cookie := StateCookie(s)
	if cookie.Name != StateCookieName {
		t.Errorf("cookie name = %q, want %q", cookie.Name, StateCookieName)
	}
	if cookie.Value != s {
		t.Errorf("cookie value = %q, want %q", cookie.Value, s)
	}
	if !cookie.HttpOnly || !cookie.Secure {
		t.Error("state cookie must be HttpOnly and Secure")
	}
	if cookie.MaxAge <= 0 {
		t.Errorf("state cookie should be short-lived, got MaxAge=%d", cookie.MaxAge)
	}
	if strings.ContainsAny(cookie.Value, "\"{},: ") {
		t.Errorf("state cookie value not HTTP-safe: %s", cookie.Value)
	}
}
