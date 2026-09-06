// Diary endpoints for the dashboard: only the authenticated user's own private
// journal entries are ever returned or mutated. All handlers require a valid
// session cookie (Discord OAuth).
package api

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/kibetnathan/minjibot/internal/domain/diary"
	"github.com/kibetnathan/minjibot/internal/ports/dto"
	"github.com/kibetnathan/minjibot/internal/ports/repository"
	authsvc "github.com/kibetnathan/minjibot/internal/services/auth"
	"github.com/labstack/echo/v5"
)

// diaryHandlers bundles dependencies for the diary endpoints.
type diaryHandlers struct {
	sess  *authsvc.SessionManager
	diary repository.DiaryRepository
}

func (a *App) registerDiaryRoutes(group *echo.Group, h *diaryHandlers) {
	group.GET("/diary", h.listEntries)
	group.POST("/diary", h.createEntry)
	group.DELETE("/diary/:entryId", h.deleteEntry)
}

// requireUser rejects the request with 401 when no valid session is present,
// returning the session's Discord user ID otherwise.
func (h *diaryHandlers) requireUser(c *echo.Context) (string, bool) {
	sess, ok := resolveSession(c, h.sess)
	if !ok {
		c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return "", false
	}
	return sess.UserID, true
}

type diaryEntryView struct {
	ID        int64  `json:"id"`
	Content   string `json:"content"`
	CreatedAt string `json:"created_at"`
}

func toDiaryEntryView(e diary.Entry) diaryEntryView {
	return diaryEntryView{
		ID:        e.ID,
		Content:   e.Content,
		CreatedAt: e.CreatedAt.UTC().Format("2006-01-02T15:04:05Z07:00"),
	}
}

// listEntries returns the current user's diary entries, newest first.
func (h *diaryHandlers) listEntries(c *echo.Context) error {
	userID, ok := h.requireUser(c)
	if !ok {
		return nil
	}
	ctx := c.Request().Context()

	entries, err := h.diary.ListByUser(ctx, userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "server_error"})
	}
	out := make([]diaryEntryView, 0, len(entries))
	for _, e := range entries {
		out = append(out, toDiaryEntryView(e))
	}
	return c.JSON(http.StatusOK, map[string]any{"items": out})
}

// createEntry adds a diary entry for the current user.
func (h *diaryHandlers) createEntry(c *echo.Context) error {
	userID, ok := h.requireUser(c)
	if !ok {
		return nil
	}
	ctx := c.Request().Context()

	var req diaryCreateRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid_body"})
	}
	content := strings.TrimSpace(req.Content)
	if content == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "content_required"})
	}

	entry, err := h.diary.Create(ctx, dto.CreateDiaryEntryParams{UserID: userID, Content: content})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "server_error"})
	}
	return c.JSON(http.StatusCreated, toDiaryEntryView(entry))
}

// diaryCreateRequest is the JSON body accepted by createEntry.
type diaryCreateRequest struct {
	Content string `json:"content"`
}

// deleteEntry removes the current user's diary entry. Entries are scoped to the
// owner so deleting someone else's entry is impossible (no-op).
func (h *diaryHandlers) deleteEntry(c *echo.Context) error {
	userID, ok := h.requireUser(c)
	if !ok {
		return nil
	}
	ctx := c.Request().Context()

	id, err := strconv.ParseInt(c.Param("entryId"), 10, 64)
	if err != nil || id <= 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid_entry_id"})
	}
	if err := h.diary.Delete(ctx, id, userID); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "server_error"})
	}
	return c.NoContent(http.StatusNoContent)
}
