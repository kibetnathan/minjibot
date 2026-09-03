package api

import (
	"context"
	"errors"
	"net/http"

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

	// DB Pool (concurrency-safe connection pool)
	pool, err := pgxpool.New(context.Background(), cfg.DBURL)
	if err != nil {
		e.Logger.Error(err.Error())
		return nil, err
	}

	srv := &http.Server{
		Addr:    ":8080",
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
	}

	group := a.Echo.Group("/api")
	a.registerAuthRoutes(group, authHandlers)
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
