package api

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	authsvc "github.com/kibetnathan/minjibot/internal/services/auth"
	"github.com/labstack/echo/v5"
)

// findCookie returns the first Set-Cookie value for name from the recorder.
func findCookie(rec *httptest.ResponseRecorder, name string) *http.Cookie {
	for _, c := range rec.Result().Cookies() {
		if c.Name == name {
			return c
		}
	}
	return nil
}

func TestLogin_SetsStateCookieAndRedirects(t *testing.T) {
	h := &authHandlers{
		oauth:    authsvc.NewDiscordOAuth("client-id", "client-secret", "http://localhost:8080"),
		frontend: "http://localhost:5173",
	}

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/auth/discord", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	if err := h.login(c); err != nil {
		t.Fatalf("login: %v", err)
	}

	if rec.Code != http.StatusFound {
		t.Fatalf("expected 302, got %d", rec.Code)
	}

	stateCookie := findCookie(rec, authsvc.StateCookieName)
	if stateCookie == nil || stateCookie.Value == "" {
		t.Fatal("login did not set a non-empty state cookie")
	}

	loc := rec.Header().Get("Location")
	u, err := url.Parse(loc)
	if err != nil {
		t.Fatalf("bad redirect URL %q: %v", loc, err)
	}
	if !strings.Contains(u.Host, "discord.com") {
		t.Errorf("expected redirect to discord.com, got %q", loc)
	}
	// The state in the auth URL must equal the state stored in the cookie.
	if got := u.Query().Get("state"); got != stateCookie.Value {
		t.Errorf("auth URL state %q != cookie state %q", got, stateCookie.Value)
	}
}

func TestCallback_StateMismatch_Rejected(t *testing.T) {
	// oauth is left nil: if the handler tried to Exchange, it would nil-panic,
	// proving the state check short-circuits before any token exchange.
	h := &authHandlers{
		frontend: "http://localhost:5173",
	}

	e := echo.New()
	// Valid-looking code, but the query state does not match the cookie state.
	req := httptest.NewRequest(http.MethodGet, "/api/auth/callback/discord?code=abc&state=attacker", nil)
	req.AddCookie(authsvc.StateCookie("legit-state"))
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	if err := h.callback(c); err != nil {
		t.Fatalf("callback: %v", err)
	}
	if rec.Code != http.StatusFound {
		t.Fatalf("expected 302 redirect, got %d", rec.Code)
	}
	if loc := rec.Header().Get("Location"); !strings.Contains(loc, "error=invalid_state") {
		t.Errorf("expected invalid_state redirect, got %q", loc)
	}
	// The state cookie should be cleared.
	if cleared := findCookie(rec, authsvc.StateCookieName); cleared == nil || cleared.MaxAge >= 0 {
		t.Error("expected state cookie to be cleared on mismatch")
	}
}

func TestCallback_MissingStateCookie_Rejected(t *testing.T) {
	h := &authHandlers{frontend: "http://localhost:5173"}

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/auth/callback/discord?code=abc&state=whatever", nil)
	// No state cookie set at all.
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	if err := h.callback(c); err != nil {
		t.Fatalf("callback: %v", err)
	}
	if loc := rec.Header().Get("Location"); !strings.Contains(loc, "error=invalid_state") {
		t.Errorf("expected invalid_state redirect without a state cookie, got %q", loc)
	}
}
