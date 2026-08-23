package api

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/kibetnathan/minjibot/internal/config"
	"github.com/kibetnathan/minjibot/internal/logger"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
)

type App struct {
	Echo *echo.Echo
	Cfg  *config.Config
	Conn *pgx.Conn
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
	return &App{
		Echo: e,
		Cfg:  cfg,
		Conn: conn,
	}, nil
}

func (a *App) Start() error {
	defer func() {
		if err := a.Conn.Close(context.Background()); err != nil {
			a.Echo.Logger.Error("Failed to close database connection", "error", err)
		}
	}()
	return a.Echo.Start(":8080")
}
