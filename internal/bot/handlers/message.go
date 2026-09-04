package handlers

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/kibetnathan/minjibot/internal/commands"
	"github.com/kibetnathan/minjibot/internal/ports/dto"
	"github.com/kibetnathan/minjibot/internal/ports/repository"
	"log/slog"
)

const DefaultPrefix = "-"

type MessageHandlerDeps struct {
	Logger       *slog.Logger
	GuildRepo    repository.GuildRepository
	SettingsRepo repository.GuildSettingsRepository
	PermRepo     repository.UserPermissionRepository
	AuditRepo    repository.AuditLogRepository
}

func RegisterMessageHandler(s *discordgo.Session, deps MessageHandlerDeps, cmdHandler *commands.CommandHandler) {
	s.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
		onMessageCreate(s, m, deps, cmdHandler)
	})
}

func onMessageCreate(s *discordgo.Session, m *discordgo.MessageCreate, deps MessageHandlerDeps, cmdHandler *commands.CommandHandler) {
	if m.Author.Bot {
		return
	}

	ctx := context.Background()

	// Resolve the guild's command prefix (fall back to the default). Looked up
	// early so both the auto-delete guard and command dispatch see the same
	// prefix.
	settings, sErr := deps.SettingsRepo.Get(ctx, m.GuildID)
	prefix := DefaultPrefix
	if sErr == nil && settings.Prefix != "" {
		prefix = settings.Prefix
	}

	// Check for one-off easter eggs (scat, etc.).
	if checkEasterEggs(s, m) {
		return
	}

	// Auto-delete messages from lurkers after 0.5 seconds — but not lurk
	// commands themselves (otherwise they can never stop lurking).
	if !isLurkCommand(m.Content, prefix) {
		if commands.IsLurking(m.GuildID, m.Author.ID) {
			chID, msgID := m.ChannelID, m.ID
			go func() {
				time.Sleep(500 * time.Millisecond)
				_ = s.ChannelMessageDelete(chID, msgID)
			}()
		}
	}

	// Ensure guild exists in DB
	guild, err := deps.GuildRepo.GetByID(ctx, m.GuildID)
	if err != nil {
		guild, err = deps.GuildRepo.Create(ctx, dto.CreateGuildParams{
			ID:          m.GuildID,
			Name:        "",
			PremiumTier: 0,
		})
		if err != nil {
			deps.Logger.Error("Failed to create guild", "error", err, "guild_id", m.GuildID)
			return
		}
	}

	// Log the message content only when the guild has explicitly opted in.
	// Message-content logging is off by default: it stores every message and,
	// left unbounded, both bloats the audit table and is a privacy liability.
	// Guilds that enable it are pruned by the retention job (see bot.App).
	if sErr == nil && settings.MessageLoggingEnabled {
		_, err = deps.AuditRepo.Create(ctx, dto.CreateAuditLogParams{
			GuildID:  m.GuildID,
			Action:   "MESSAGE_CREATE",
			ActorID:  m.Author.ID,
			TargetID: m.ChannelID,
			Metadata: []byte(fmt.Sprintf(`{"message_id":%q,"content":%q,"channel_id":%q}`, m.ID, m.Content, m.ChannelID)),
		})
		if err != nil {
			deps.Logger.Error("Failed to create audit log", "error", err)
		}
	}

	_ = guild
	_ = settings

	// Check for command
	if !strings.HasPrefix(m.Content, prefix) {
		return
	}

	// Commands can be chained together with "&&", e.g. "-spark && -smoke".
	// Each segment is trimmed and dispatched in order; a failing segment does
	// not stop the remaining ones.
	chain := strings.Split(m.Content, "&&")
	for _, segment := range chain {
		dispatchCommand(ctx, s, m, prefix, segment, cmdHandler)
	}
}

// dispatchCommand runs a single command segment (already extracted from any
// "&&" chain) if it is a valid command invocation.
func dispatchCommand(
	ctx context.Context,
	s *discordgo.Session,
	m *discordgo.MessageCreate,
	prefix string,
	segment string,
	cmdHandler *commands.CommandHandler,
) {
	segment = strings.TrimSpace(segment)
	if !strings.HasPrefix(segment, prefix) {
		return
	}

	args := strings.Fields(strings.TrimPrefix(segment, prefix))
	if len(args) == 0 {
		return
	}

	if err := cmdHandler.Handle(ctx, s, m, args[0], args[1:]); err != nil {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Error: %v", err))
	}
}

// isLurkCommand reports whether a raw message content is a lurk or lurkers
// command using the given prefix. This lets lurkers toggle their state without
// the auto-delete kicking in.
func isLurkCommand(content, prefix string) bool {
	low := strings.ToLower(strings.TrimSpace(content))
	if prefix != "" && strings.HasPrefix(low, prefix) {
		low = strings.TrimSpace(strings.TrimPrefix(low, prefix))
	}
	return low == "lurk" || strings.HasPrefix(low, "lurk ") ||
		low == "lurkers" || strings.HasPrefix(low, "lurkers ")
}
