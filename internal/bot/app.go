package bot

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kibetnathan/minjibot/infrastructure/postgres"
	"github.com/kibetnathan/minjibot/internal/bot/handlers"
	"github.com/kibetnathan/minjibot/internal/commands"
	"github.com/kibetnathan/minjibot/internal/config"
	"github.com/kibetnathan/minjibot/internal/logger"
	"github.com/kibetnathan/minjibot/internal/ports/repository"
)

type App struct {
	Session      *discordgo.Session
	Cfg          *config.Config
	Pool         *pgxpool.Pool
	Logger       *slog.Logger
	GuildRepo    repository.GuildRepository
	SettingsRepo repository.GuildSettingsRepository
	PermRepo     repository.UserPermissionRepository
	AuditRepo    repository.AuditLogRepository
	UserRepo     repository.UserRepository
	BirthdayRepo repository.BirthdayRepository
	BirthdaySett repository.GuildBirthdaySettingsRepository
	DiaryRepo    repository.DiaryRepository
	cmdHandler   *commands.CommandHandler
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
	session, err := discordgo.New("Bot " + cfg.DiscordToken)
	if err != nil {
		log.Error("Error creating Discord session", "error", err.Error())
		return nil, err
	}

	// DB Pool (concurrency-safe connection pool for parallel event handlers)
	pool, err := pgxpool.New(context.Background(), cfg.DBURL)
	if err != nil {
		log.Error("Error connecting to database", "error", err.Error())
		return nil, err
	}
	if err := pool.Ping(context.Background()); err != nil {
		log.Error("Error pinging database", "error", err.Error())
		return nil, err
	}

	store := repository.NewSQLStore(postgres.New(pool))

	app := &App{
		Session:      session,
		Cfg:          cfg,
		Pool:         pool,
		Logger:       log,
		GuildRepo:    repository.NewGuildRepository(store),
		SettingsRepo: repository.NewGuildSettingsRepository(store),
		PermRepo:     repository.NewUserPermissionRepository(store),
		AuditRepo:    repository.NewAuditLogRepository(store),
		UserRepo:     repository.NewUserRepository(store),
		BirthdayRepo: repository.NewBirthdayRepository(store),
		BirthdaySett: repository.NewGuildBirthdaySettingsRepository(store),
		DiaryRepo:    repository.NewDiaryRepository(store),
	}

	// Initialize command handler
	app.cmdHandler = commands.NewCommandHandler(app.Cfg, app.GuildRepo, app.SettingsRepo, app.PermRepo, app.BirthdayRepo, app.BirthdaySett, app.DiaryRepo)

	// Register event handlers
	app.RegisterHandlers()

	return app, nil
}

func (a *App) RegisterHandlers() {
	// Identify required gateway intents (adjust based on your bot's needs)
	a.Session.Identify.Intents = discordgo.IntentsGuilds | discordgo.IntentsGuildMessages | discordgo.IntentMessageContent | discordgo.IntentsGuildMessageReactions

	// Keep a rolling cache of recent messages so content of deleted
	// messages is still available to log to the database.
	a.Session.State.MaxMessageCount = 2000

	// Ready handler - register slash commands
	a.Session.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		a.Logger.Info(fmt.Sprintf("Logged in as: %s#%s", r.User.Username, r.User.Discriminator))

		// Register global slash commands
		appID := a.Cfg.DiscordClientID
		if appID == "" {
			appID = s.State.User.ID
			a.Logger.Warn("DISCORD_CLIENT_ID not set, falling back to bot user ID for command registration")
		}
		a.Logger.Info("Registering slash commands", "app_id", appID)
		for _, cmd := range commands.SlashCommands {
			_, err := s.ApplicationCommandCreate(appID, "", cmd)
			if err != nil {
				a.Logger.Error("Failed to register slash command", "command", cmd.Name, "error", err, "app_id", appID)
			} else {
				a.Logger.Info("Registered slash command", "command", cmd.Name)
			}
		}
	})

	// Message handler
	handlers.RegisterMessageHandler(a.Session, handlers.MessageHandlerDeps{
		Logger:       a.Logger,
		GuildRepo:    a.GuildRepo,
		SettingsRepo: a.SettingsRepo,
		PermRepo:     a.PermRepo,
		AuditRepo:    a.AuditRepo,
	}, a.cmdHandler)

	// Message delete handler
	handlers.RegisterMessageDeleteHandler(a.Session, handlers.MessageDeleteHandlerDeps{
		Logger:    a.Logger,
		GuildRepo: a.GuildRepo,
		AuditRepo: a.AuditRepo,
	})

	// Interaction handler
	handlers.RegisterInteractionHandler(a.Session, handlers.InteractionHandlerDeps{
		Logger:       a.Logger,
		GuildRepo:    a.GuildRepo,
		SettingsRepo: a.SettingsRepo,
		PermRepo:     a.PermRepo,
		AuditRepo:    a.AuditRepo,
	}, a.cmdHandler)
}
func (a *App) Start() error {
	// Open WebSocket connection to Discord
	if err := a.Session.Open(); err != nil {
		a.Logger.Error("Error opening Discord websocket connection", "error", err.Error())
		return err
	}

	a.Logger.Info("Bot is now running. Press CTRL-C to exit.")

	// Block until SIGINT or SIGTERM is received
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	a.Logger.Info("Shutting down bot...")

	if err := a.Session.Close(); err != nil {
		a.Logger.Error("Failed to close Discord session", "error", err)
	}

	a.Pool.Close()
	a.Logger.Info("Closed database pool")

	return nil
}
