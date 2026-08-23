package api

import (
	"github.com/jackc/pgx/v5"
	"github.com/kibetnathan/minjibot/internal/config"
	"github.com/labstack/echo/v5"
)

type App struct {
	Echo *echo.Echo
	Cfg  *config.Config
	Conn *pgx.Conn
}
