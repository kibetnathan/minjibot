package bot

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/bwmarrin/discordgo"
	"github.com/jackc/pgx/v5"
	"github.com/kibetnathan/minjibot/internal/config"
	"github.com/kibetnathan/minjibot/internal/logger"
	"github.com/kibetnathan/minjibot/internal/bot/handlers"
	"github.com/kibetnathan/minjibot/internal/ports/repository"
	"github.com/kibetnathan/minjibot/infrastructure/postgres"
)

type App struct {
	Session     *discordgo.Session
	Cfg         *config.Config
	Conn        *pgx.Conn
	Logger      *slog.Logger
	GuildRepo   repository.GuildRepository
	SettingsRepo repository.GuildSettingsRepository
	PermRepo    repository.UserPermissionRepository
	AuditRepo   repository.AuditLogRepository
}

func NewApp() (*App, error) {
	log := logger.New()

	//  Initialise Config
	cfg, err := config.NewConfig()
	if err != nil {
		log.Error("Error parsing env file", "error", err.Error())
		return nil, err
	}

	// Initialise Discord Session
	// Note: Bot tokens require the "Bot " prefix in discordgo
	session, err := discordgo.New("Bot " + cfg.DiscordToken)
	if err != nil {
		log.Error("Error creating Discord session", "error", err.Error())
		return nil, err
	}

	// DB Conn
	conn, err := pgx.Connect(context.Background(), cfg.DBURL)
	if err != nil {
		log.Error("Error connecting to database", "error", err.Error())
		return nil, err
	}

	store := repository.NewSQLStore(postgres.New(conn))

	app := &App{
		Session:      session,
		Cfg:          cfg,
		Conn:         conn,
		Logger:       log,
		GuildRepo:    repository.NewGuildRepository(store),
		SettingsRepo: repository.NewGuildSettingsRepository(store),
		PermRepo:     repository.NewUserPermissionRepository(store),
		AuditRepo:    repository.NewAuditLogRepository(store),
	}

	// Register event handlers
	app.registerHandlers()

	return app, nil
}

func (a *App) registerHandlers() {
	// Identify required gateway intents (adjust based on your bot's needs)
	a.Session.Identify.Intents = discordgo.IntentsGuilds | discordgo.IntentsGuildMessages | discordgo.IntentMessageContent

	// Ready handler
	a.Session.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		a.Logger.Info(fmt.Sprintf("Logged in as: %s#%s", r.User.Username, r.User.Discriminator))
	})

	// Message handler
	handlers.RegisterMessageHandler(a.Session, handlers.HandlerDeps{
		Logger:       a.Logger,
		GuildRepo:    a.GuildRepo,
		SettingsRepo: a.SettingsRepo,
		AuditRepo:    a.AuditRepo,
	})
}

func (a *App) Start() error {
	defer func() {
		if err := a.Conn.Close(context.Background()); err != nil {
			a.Logger.Error("Failed to close database connection", "error", err)
		}
	}()

	// Open WebSocket connection to Discord
	if err := a.Session.Open(); err != nil {
		a.Logger.Error("Error opening Discord websocket connection", "error", err.Error())
		return err
	}

	defer func() {
		if err := a.Session.Close(); err != nil {
			a.Logger.Error("Failed to close Discord session", "error", err)
		}
	}()

	a.Logger.Info("Bot is now running. Press CTRL-C to exit.")

	// Block until a signal is received to keep the bot alive
	sc := make(chan struct{})
	<-sc

	return nil
}
