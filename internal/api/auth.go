package api

import (
	"errors"
	"net/http"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/kibetnathan/minjibot/internal/domain/user"
	"github.com/kibetnathan/minjibot/internal/ports/dto"
	"github.com/kibetnathan/minjibot/internal/ports/repository"
	authsvc "github.com/kibetnathan/minjibot/internal/services/auth"
	"github.com/labstack/echo/v5"
)

const userContextKey = "auth.user"

// authHandlers bundles the dependencies shared by the auth routes.
type authHandlers struct {
	oauth    *authsvc.DiscordOAuth
	sess     *authsvc.SessionManager
	authz    *authorizer
	users    repository.UserRepository
	frontend string
}

// registerAuthRoutes wires the Discord OAuth flow onto the given Echo group.
func (a *App) registerAuthRoutes(group *echo.Group, h *authHandlers) {
	group.GET("/auth/discord", h.login)
	group.GET("/auth/callback/discord", h.callback)
	group.GET("/auth/logout", h.logout)
	group.GET("/auth/me", h.me)
}

// login redirects the user to Discord's authorization page. It issues a random
// state token, stores it in a short-lived cookie, and passes it to Discord so
// the callback can verify the response was not forged (OAuth anti-CSRF).
func (h *authHandlers) login(c *echo.Context) error {
	state, err := authsvc.NewState()
	if err != nil {
		return c.Redirect(http.StatusFound, strings.TrimRight(h.frontend, "/")+"/login?error=state_error")
	}
	c.SetCookie(authsvc.StateCookie(state))
	return c.Redirect(http.StatusFound, h.oauth.AuthURL(state))
}

// callback handles the redirect from Discord after the user approves access.
func (h *authHandlers) callback(c *echo.Context) error {
	frontend := strings.TrimRight(h.frontend, "/")

	code := c.QueryParam("code")
	if code == "" {
		return c.Redirect(http.StatusFound, frontend+"/login?error=missing_code")
	}

	// Verify the anti-CSRF state before doing anything with the code: the state
	// returned by Discord must match the one we stored in the login cookie.
	stateCookie, err := c.Cookie(authsvc.StateCookieName)
	if err != nil || !authsvc.StateMatches(c.QueryParam("state"), stateCookie.Value) {
		c.SetCookie(authsvc.ClearStateCookie())
		return c.Redirect(http.StatusFound, frontend+"/login?error=invalid_state")
	}
	// State is single-use; clear it now that it has been consumed.
	c.SetCookie(authsvc.ClearStateCookie())

	discordUser, err := h.oauth.Exchange(c.Request().Context(), code)
	if err != nil {
		return c.Redirect(http.StatusFound, frontend+"/login?error=oauth_failed")
	}

	u, err := h.upsertUser(c, discordUser)
	if err != nil {
		return c.Redirect(http.StatusFound, frontend+"/login?error=account_error")
	}

	cookie, err := h.sess.Create(u.UserID)
	if err != nil {
		return c.Redirect(http.StatusFound, frontend+"/login?error=session_error")
	}
	c.SetCookie(cookie)

	return c.Redirect(http.StatusFound, frontend+"/dashboard")
}

// logout clears the session cookie and returns the user to the home page.
func (h *authHandlers) logout(c *echo.Context) error {
	c.SetCookie(h.sess.Clear())
	return c.Redirect(http.StatusFound, strings.TrimRight(h.frontend, "/")+"/")
}

// me returns the currently authenticated user, or 401 if there is none.
func (h *authHandlers) me(c *echo.Context) error {
	sess, ok := h.currentSession(c)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
	}

	u, err := h.users.GetByDiscordID(c.Request().Context(), sess.UserID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "server_error"})
	}
	if !u.IsActive {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "account_disabled"})
	}

	return c.JSON(http.StatusOK, map[string]any{
		"id":       u.UserID,
		"email":    u.Email,
		"is_admin": h.authz.isAdmin(u.UserID),
	})
}

// currentSession reads and verifies the session cookie, returning the session
// and whether it was valid.
func (h *authHandlers) currentSession(c *echo.Context) (*authsvc.Session, bool) {
	return resolveSession(c, h.sess)
}

// upsertUser creates the user from their Discord profile if they do not exist
// yet, otherwise refreshes their stored email.
func (h *authHandlers) upsertUser(c *echo.Context, du *authsvc.DiscordUser) (user.User, error) {
	ctx := c.Request().Context()

	u, err := h.users.GetByDiscordID(ctx, du.ID)
	if err == nil {
		if u.Email != du.Email && du.Email != "" {
			email := du.Email
			_, uErr := h.users.Update(ctx, u.ID, dto.UpdateUserParams{
				Email:    &email,
				IsActive: u.IsActive,
			})
			if uErr != nil {
				return user.User{}, uErr
			}
			u.Email = du.Email
		}
		return u, nil
	}
	if err != nil && !isNoRows(err) {
		return user.User{}, err
	}

	email := du.Email
	created, err := h.users.Create(ctx, dto.CreateUserParams{
		UserID: du.ID,
		Email:  &email,
	})
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return h.users.GetByDiscordID(ctx, du.ID)
		}
		return user.User{}, err
	}
	return created, nil
}

func isNoRows(err error) bool {
	return errors.Is(err, pgx.ErrNoRows)
}
