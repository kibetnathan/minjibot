package api

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kibetnathan/minjibot/internal/domain/guild"
	"github.com/kibetnathan/minjibot/internal/ports/repository"
	authsvc "github.com/kibetnathan/minjibot/internal/services/auth"
	"github.com/labstack/echo/v5"
)

// fakeGuildRepo satisfies GuildRepository by embedding the interface (all other
// methods are nil and would panic if called, but the tested path only calls List).
type fakeGuildRepo struct {
	repository.GuildRepository
	guilds []guild.Guild
}

func (f *fakeGuildRepo) List(context.Context) ([]guild.Guild, error) { return f.guilds, nil }

type fakeDeletedRepo struct {
	repository.DeletedMessageRepository
}

func (f *fakeDeletedRepo) CountForAllGuilds(context.Context) (map[string]int64, error) {
	return map[string]int64{}, nil
}

type fakeAuditRepo struct {
	repository.AuditLogRepository
}

func (f *fakeAuditRepo) CountForAllGuilds(context.Context) (map[string]int64, error) {
	return map[string]int64{}, nil
}

// newTestLogHandlers builds logHandlers wired to a session manager with a known
// secret and an allowlist containing only "admin-id".
func newTestLogHandlers() (*logHandlers, *authsvc.SessionManager) {
	sess := authsvc.NewSessionManager("test-secret")
	return &logHandlers{
		sess:    sess,
		authz:   newAuthorizer([]string{"admin-id"}),
		guilds:  &fakeGuildRepo{guilds: []guild.Guild{{ID: "g1", Name: "Guild One"}}},
		audits:  &fakeAuditRepo{},
		deletes: &fakeDeletedRepo{},
	}, sess
}

// doListGuilds runs listGuilds with an optional session cookie for userID (empty
// = no cookie) and returns the recorded response.
func doListGuilds(t *testing.T, h *logHandlers, sess *authsvc.SessionManager, userID string) *httptest.ResponseRecorder {
	t.Helper()
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/guilds", nil)
	if userID != "" {
		cookie, err := sess.Create(userID)
		if err != nil {
			t.Fatalf("create session: %v", err)
		}
		req.AddCookie(cookie)
	}
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	if err := h.listGuilds(c); err != nil {
		t.Fatalf("listGuilds returned error: %v", err)
	}
	return rec
}

func TestListGuilds_NoSession_401(t *testing.T) {
	h, sess := newTestLogHandlers()
	rec := doListGuilds(t, h, sess, "")
	if rec.Code != http.StatusUnauthorized {
		t.Errorf("expected 401 without a session, got %d", rec.Code)
	}
}

func TestListGuilds_NonAdmin_403(t *testing.T) {
	h, sess := newTestLogHandlers()
	rec := doListGuilds(t, h, sess, "not-an-admin")
	if rec.Code != http.StatusForbidden {
		t.Errorf("expected 403 for a non-admin session, got %d", rec.Code)
	}
}

func TestListGuilds_Admin_200(t *testing.T) {
	h, sess := newTestLogHandlers()
	rec := doListGuilds(t, h, sess, "admin-id")
	if rec.Code != http.StatusOK {
		t.Errorf("expected 200 for an admin session, got %d (body: %s)", rec.Code, rec.Body.String())
	}
}

// A session signed with a different secret must not be accepted (tamper / wrong key).
func TestListGuilds_WrongSecret_401(t *testing.T) {
	h, _ := newTestLogHandlers()
	other := authsvc.NewSessionManager("a-different-secret")
	rec := doListGuilds(t, h, other, "admin-id")
	if rec.Code != http.StatusUnauthorized {
		t.Errorf("expected 401 for a session signed with the wrong secret, got %d", rec.Code)
	}
}
