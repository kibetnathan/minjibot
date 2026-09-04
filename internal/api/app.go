// Package api provides the Echo HTTP server for MinjiBot's REST API,
// including Discord OAuth authentication and dashboard data endpoints.
package api

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kibetnathan/minjibot/infrastructure/postgres"
	"github.com/kibetnathan/minjibot/internal/config"
	"github.com/kibetnathan/minjibot/internal/logger"
	"github.com/kibetnathan/minjibot/internal/ports/repository"
	authsvc "github.com/kibetnathan/minjibot/internal/services/auth"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
)

type App struct {
	Echo   *echo.Echo
	Cfg    *config.Config
	Pool   *pgxpool.Pool
	server *http.Server
}

func NewApp() (*App, error) {
	e := echo.New()
	e.Logger = logger.New()
	e.Use(middleware.Recover())

	// Health check endpoint so the bot's web service (Render free tier) can
	// report as live and receive its ping.
	var healthHandler echo.HandlerFunc = func(c *echo.Context) error {
		return c.String(http.StatusOK, "ok")
	}
	e.GET("/healthz", healthHandler)

	// Initialise Config
	cfg, err := config.NewConfig()
	if err != nil {
		e.Logger.Error("Error parsing env file", "Error:", err.Error())
		return nil, err
	}

	// CORS for cross-origin browser calls to the API (e.g. the dashboard
	// reading /api/auth/me from a separate frontend origin). Credentials are
	// allowed so the session cookie is sent, but only for explicit origins —
	// never the wildcard (echo v5 rejects "*" with AllowCredentials anyway).
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     corsOrigins(cfg),
		AllowMethods:     []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete, http.MethodOptions},
		AllowHeaders:     []string{"Authorization", "Content-Type", "X-Requested-With"},
		AllowCredentials: true,
	}))

	// DB Pool (concurrency-safe connection pool)
	pool, err := pgxpool.New(context.Background(), cfg.DBURL)
	if err != nil {
		e.Logger.Error(err.Error())
		return nil, err
	}

	port := cfg.Port
	if port == "" {
		port = "8080"
	}

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: e,
	}

	app := &App{
		Echo:   e,
		Cfg:    cfg,
		Pool:   pool,
		server: srv,
	}

	app.registerRoutes()

	return app, nil
}

// registerRoutes wires the HTTP API routes, including the Discord OAuth flow.
func (a *App) registerRoutes() {
	store := repository.NewSQLStore(postgres.New(a.Pool))
	userRepo := repository.NewUserRepository(store)

	authHandlers := &authHandlers{
		oauth: authsvc.NewDiscordOAuth(
			a.Cfg.DiscordClientID,
			a.Cfg.DiscordClientSecret,
			a.Cfg.AppURL,
		),
		sess:  authsvc.NewSessionManager(a.Cfg.SessionSecret),
		users: userRepo,
		// Redirect the user back to the frontend after login/logout. Defaults
		// to the API origin when FRONTEND_URL is unset (same-origin dev setup).
		frontend: a.Cfg.FrontendURL,
	}

	group := a.Echo.Group("/api")
	a.registerAuthRoutes(group, authHandlers)

	logHandlers := &logHandlers{
		sess:    authsvc.NewSessionManager(a.Cfg.SessionSecret),
		guilds:  repository.NewGuildRepository(store),
		audits:  repository.NewAuditLogRepository(store),
		deletes: repository.NewDeletedMessageRepository(store),
	}
	a.registerLogRoutes(group, logHandlers)
}

func (a *App) Start() error {
	defer func() {
		a.Pool.Close()
	}()
	return a.server.ListenAndServe()
}

func (a *App) Shutdown(ctx context.Context) error {
	err := a.server.Shutdown(ctx)
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
}

// corsOrigins returns the origins allowed by the CORS middleware: the frontend
// origin (FRONTEND_URL), the API origin (APP_URL), plus localhost dev origins so
// a local Vite dashboard can call the API directly.
func corsOrigins(cfg *config.Config) []string {
	seen := map[string]struct{}{}
	var origins []string
	for _, o := range []string{
		strings.TrimRight(cfg.FrontendURL, "/"),
		strings.TrimRight(cfg.AppURL, "/"),
		"http://localhost:5173",
		"http://localhost:8080",
	} {
		if o == "" {
			continue
		}
		if _, ok := seen[o]; ok {
			continue
		}
		seen[o] = struct{}{}
		origins = append(origins, o)
	}
	return origins
}
