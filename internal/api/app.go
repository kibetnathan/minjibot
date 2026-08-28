package api

import (
	"context"
	"errors"
	"net/http"

	"github.com/jackc/pgx/v5"
	"github.com/kibetnathan/minjibot/internal/config"
	"github.com/kibetnathan/minjibot/internal/logger"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
)

type App struct {
	Echo   *echo.Echo
	Cfg    *config.Config
	Conn   *pgx.Conn
	server *http.Server
}

func NewApp() (*App, error) {
	e := echo.New()
	e.Logger = logger.New()
	e.Use(middleware.Recover())

	// Initialise Config
	cfg, err := config.NewConfig()
	if err != nil {
		e.Logger.Error("Error parsing env file", "Error:", err.Error())
		return nil, err
	}

	// DB Conn
	conn, err := pgx.Connect(context.Background(), cfg.DBURL)
	if err != nil {
		e.Logger.Error(err.Error())
		return nil, err
	}

	srv := &http.Server{
		Addr:    ":8080",
		Handler: e,
	}

	return &App{
		Echo:   e,
		Cfg:    cfg,
		Conn:   conn,
		server: srv,
	}, nil
}

func (a *App) Start() error {
	defer func() {
		if err := a.Conn.Close(context.Background()); err != nil {
			a.Echo.Logger.Error("Failed to close database connection", "error", err)
		}
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
