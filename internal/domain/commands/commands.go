package commands

import (
	"context"
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/kibetnathan/minjibot/internal/ports/repository"
)

type CommandHandler struct {
	GuildRepo    repository.GuildRepository
	SettingsRepo repository.GuildSettingsRepository
	PermRepo     repository.UserPermissionRepository
}

func NewCommandHandler(guildRepo repository.GuildRepository, settingsRepo repository.GuildSettingsRepository, permRepo repository.UserPermissionRepository) *CommandHandler {
	return &CommandHandler{
		GuildRepo:    guildRepo,
		SettingsRepo: settingsRepo,
		PermRepo:     permRepo,
	}
}

func (h *CommandHandler) Handle(ctx context.Context, s *discordgo.Session, m *discordgo.MessageCreate, cmd string, args []string) error {
	switch cmd {
	case "ping":
		return h.ping(s, m.ChannelID)
	case "help":
		return h.help(s, m.ChannelID)
	case "echo":
		return h.echo(s, m.ChannelID, args)
	case "userinfo":
		return h.userInfo(s, m, args)
	case "ddg":
		return h.ddg(s, m.ChannelID, args)
	case "search":
		return h.search(s, m, args)
	default:
		return fmt.Errorf("unknown command: %s", cmd)
	}
}

func (h *CommandHandler) HandleSlash(ctx context.Context, s *discordgo.Session, i *discordgo.InteractionCreate) error {
	switch i.ApplicationCommandData().Name {
	case "ping":
		return h.pingSlash(s, i)
	case "help":
		return h.helpSlash(s, i)
	case "echo":
		return h.echoSlash(s, i)
	case "userinfo":
		return h.userInfoSlash(s, i)
	case "ddg":
		return h.ddgSlash(s, i)
	case "search":
		return h.searchSlash(s, i)
	default:
		return fmt.Errorf("unknown command: %s", i.ApplicationCommandData().Name)
	}
}

func (h *CommandHandler) ping(s *discordgo.Session, channelID string) error {
	return pingMessageCommandHandler(s, channelID)
}

func (h *CommandHandler) pingSlash(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	return pingSlashCommandHandler(s, i)
}

func (h *CommandHandler) help(s *discordgo.Session, channelID string) error {
	_, err := s.ChannelMessageSend(channelID, "Available commands: `ping`, `help`, `echo`, `userinfo`, `ddg`, `search`")
	return err
}

func (h *CommandHandler) helpSlash(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Available commands: `/ping`, `/help`, `/echo`, `/userinfo`, `/ddg`, `/search`",
		},
	})
}

func (h *CommandHandler) echo(s *discordgo.Session, channelID string, args []string) error {
	return echoMessageCommandHandler(s, channelID, args)
}

func (h *CommandHandler) echoSlash(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	return echoSlashCommandHandler(s, i)
}

func (h *CommandHandler) userInfo(s *discordgo.Session, m *discordgo.MessageCreate, args []string) error {
	return userInfoMessageCommandHandler(s, m, args)
}

func (h *CommandHandler) userInfoSlash(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	return userInfoSlashCommandHandler(s, i)
}

func (h *CommandHandler) ddg(s *discordgo.Session, channelID string, args []string) error {
	return ddgMessageCommandHandler(s, channelID, args)
}

func (h *CommandHandler) ddgSlash(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	return ddgSlashCommandHandler(s, i)
}

func (h *CommandHandler) search(s *discordgo.Session, m *discordgo.MessageCreate, args []string) error {
	return searchMessageCommandHandler(s, m, args)
}

func (h *CommandHandler) searchSlash(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	return searchSlashCommandHandler(s, i)
}

var SlashCommands = []*discordgo.ApplicationCommand{
	{
		Name:        "ping",
		Description: "Check bot latency",
	},
	{
		Name:        "help",
		Description: "Show available commands",
	},
	{
		Name:        "echo",
		Description: "Repeat back a message",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "text",
				Description: "Text to echo",
				Required:    true,
			},
		},
	},
	{
		Name:        "userinfo",
		Description: "Get info about a user",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionUser,
				Name:        "user",
				Description: "The user to look up (defaults to yourself)",
				Required:    false,
			},
		},
	},
	{
		Name:        "ddg",
		Description: "Fetch quick search results from DuckDuckGo",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "query",
				Description: "The search query",
				Required:    true,
			},
		},
	},
	{
		Name:        "search",
		Description: "Search chat history for a specific message",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "query",
				Description: "Text to search for in chat history",
				Required:    true,
			},
			{
				Type:        discordgo.ApplicationCommandOptionInteger,
				Name:        "messages",
				Description: "How many recent messages to search (default 200, max 1000)",
				Required:    false,
			},
		},
	},
}

