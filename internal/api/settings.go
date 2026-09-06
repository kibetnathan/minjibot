// Per-guild settings endpoints for the dashboard: reading and updating a
// guild's bot configuration (prefix, language, auto-moderation, log channel).
// All handlers require a valid session cookie.
package api

import (
	"errors"
	"net/http"

	"github.com/jackc/pgx/v5"
	"github.com/kibetnathan/minjibot/internal/domain/guildsettings"
	"github.com/kibetnathan/minjibot/internal/ports/dto"
	"github.com/kibetnathan/minjibot/internal/ports/repository"
	authsvc "github.com/kibetnathan/minjibot/internal/services/auth"
	"github.com/labstack/echo/v5"
)

// settingsHandlers bundles dependencies for the server settings endpoints.
type settingsHandlers struct {
	sess     *authsvc.SessionManager
	settings repository.GuildSettingsRepository
}

func (a *App) registerSettingsRoutes(group *echo.Group, h *settingsHandlers) {
	group.GET("/guilds/:guildId/settings", h.getSettings)
	group.PUT("/guilds/:guildId/settings", h.updateSettings)
}

// requireUser rejects the request with 401 when no valid session is present,
// returning the session's Discord user ID otherwise.
func (h *settingsHandlers) requireUser(c *echo.Context) (string, bool) {
	sess, ok := resolveSession(c, h.sess)
	if !ok {
		c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return "", false
	}
	return sess.UserID, true
}

type settingsView struct {
	GuildID               string `json:"guild_id"`
	Prefix                string `json:"prefix"`
	Language              string `json:"language"`
	AutoModerationEnabled bool   `json:"auto_moderation_enabled"`
	LoggingChannelID      string `json:"logging_channel_id"`
}

// settingsUpdateRequest is the partial JSON body accepted by updateSettings.
// Pointer fields distinguish "not provided" from an explicit empty value.
type settingsUpdateRequest struct {
	Prefix                *string `json:"prefix"`
	Language              *string `json:"language"`
	AutoModerationEnabled *bool   `json:"auto_moderation_enabled"`
	LoggingChannelID      *string `json:"logging_channel_id"`
}

func toSettingsView(s guildsettings.GuildSettings) settingsView {
	return settingsView{
		GuildID:               s.GuildID,
		Prefix:                s.Prefix,
		Language:              s.Language,
		AutoModerationEnabled: s.AutoModerationEnabled,
		LoggingChannelID:      s.LoggingChannelID,
	}
}

// getSettings returns the stored settings for a guild, or a default empty
// settings object when none have been configured yet.
func (h *settingsHandlers) getSettings(c *echo.Context) error {
	if _, ok := h.requireUser(c); !ok {
		return nil
	}
	ctx := c.Request().Context()

	settings, err := h.settings.Get(ctx, c.Param("guildId"))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return c.JSON(http.StatusOK, settingsView{
				GuildID: c.Param("guildId"),
				Prefix:  "-",
			})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "server_error"})
	}
	return c.JSON(http.StatusOK, toSettingsView(settings))
}

// updateSettings merges a partial update onto a guild's existing settings (or
// defaults) and persists it, so a request only needs to send the fields it
// changes. Empty values for prefix/log channel are rejected to avoid clearing
// them accidentally via a partial request.
func (h *settingsHandlers) updateSettings(c *echo.Context) error {
	if _, ok := h.requireUser(c); !ok {
		return nil
	}
	ctx := c.Request().Context()
	guildID := c.Param("guildId")

	var req settingsUpdateRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid_body"})
	}

	// Merge onto current stored values so partial updates behave sanely.
	current := guildsettings.GuildSettings{
		GuildID:  guildID,
		Prefix:   "-",
		Language: "en",
	}
	if s, err := h.settings.Get(ctx, guildID); err == nil {
		current = s
	}

	if req.Prefix != nil {
		if *req.Prefix == "" {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "prefix_empty"})
		}
		current.Prefix = *req.Prefix
	}
	if req.Language != nil {
		current.Language = *req.Language
	}
	if req.AutoModerationEnabled != nil {
		current.AutoModerationEnabled = *req.AutoModerationEnabled
	}
	if req.LoggingChannelID != nil {
		current.LoggingChannelID = *req.LoggingChannelID
	}

	updated, err := h.settings.Upsert(ctx, dto.UpsertGuildSettingsParams{
		GuildID:               current.GuildID,
		Prefix:                current.Prefix,
		Language:              current.Language,
		AutoModerationEnabled: current.AutoModerationEnabled,
		LoggingChannelID:      current.LoggingChannelID,
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "server_error"})
	}
	return c.JSON(http.StatusOK, toSettingsView(updated))
}
