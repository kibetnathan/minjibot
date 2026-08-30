package commands

import (
	"context"
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/kibetnathan/minjibot/internal/config"
	"github.com/kibetnathan/minjibot/internal/ports/repository"
)

type CommandHandler struct {
	Cfg          *config.Config
	GuildRepo    repository.GuildRepository
	SettingsRepo repository.GuildSettingsRepository
	PermRepo     repository.UserPermissionRepository
}

func NewCommandHandler(cfg *config.Config, guildRepo repository.GuildRepository, settingsRepo repository.GuildSettingsRepository, permRepo repository.UserPermissionRepository) *CommandHandler {
	return &CommandHandler{
		Cfg:          cfg,
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
	case "pinglist":
		return h.pinglist(s, m, args)
	case "gifsearch":
		return h.gifsearch(s, m.ChannelID, args)
	case "emoji":
		return h.emoji(s, m, args)
	case "sticker":
		return h.sticker(s, m, args)
	case "pin":
		return h.pin(s, m, args)
	case "unpin":
		return h.unpin(s, m, args)
	case "quote":
		return h.quote(s, m, args)
	case "translate":
		return h.translate(s, m, args)
	case "reminder":
		return h.reminder(s, m, args)
	case "isearch":
		return h.isearch(s, m.ChannelID, args)
	case "caption":
		return h.caption(s, m, args)
	case "img2gif":
		return h.img2gif(s, m, args)
	case "vid2gif":
		return h.vid2gif(s, m, args)
	case "autogif":
		return h.autogif(s, m, args)
	case "factcheck":
		return h.factcheck(s, m, args)
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
	case "pinglist":
		return h.pinglistSlash(s, i)
	case "gifsearch":
		return h.gifsearchSlash(s, i)
	case "emoji":
		return h.emojiSlash(s, i)
	case "sticker":
		return h.stickerSlash(s, i)
	case "pin":
		return h.pinSlash(s, i)
	case "unpin":
		return h.unpinSlash(s, i)
	case "quote":
		return h.quoteSlash(s, i)
	case "translate":
		return h.translateSlash(s, i)
	case "reminder":
		return h.reminderSlash(s, i)
	case "isearch":
		return h.isearchSlash(s, i)
	case "caption":
		return h.captionSlash(s, i)
	case "img2gif":
		return h.img2gifSlash(s, i)
	case "vid2gif":
		return h.vid2gifSlash(s, i)
	case "autogif":
		return h.autogifSlash(s, i)
	case "factcheck":
		return h.factcheckSlash(s, i)
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
	_, err := s.ChannelMessageSendEmbed(channelID, BuildHelpEmbed())
	return err
}

func (h *CommandHandler) helpSlash(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{BuildHelpEmbed()},
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

func (h *CommandHandler) pinglist(s *discordgo.Session, m *discordgo.MessageCreate, args []string) error {
	return pinglistMessageCommandHandler(s, m, args)
}

func (h *CommandHandler) pinglistSlash(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	return pinglistSlashCommandHandler(s, i)
}

func (h *CommandHandler) gifsearch(s *discordgo.Session, channelID string, args []string) error {
	return gifsearchMessageCommandHandler(s, channelID, args)
}

func (h *CommandHandler) gifsearchSlash(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	return gifsearchSlashCommandHandler(s, i)
}

func (h *CommandHandler) emoji(s *discordgo.Session, m *discordgo.MessageCreate, args []string) error {
	return emojiMessageCommandHandler(s, m, args)
}

func (h *CommandHandler) emojiSlash(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	return emojiSlashCommandHandler(s, i)
}

func (h *CommandHandler) sticker(s *discordgo.Session, m *discordgo.MessageCreate, args []string) error {
	return stickerMessageCommandHandler(s, m, args)
}

func (h *CommandHandler) stickerSlash(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	return stickerSlashCommandHandler(s, i)
}

func (h *CommandHandler) pin(s *discordgo.Session, m *discordgo.MessageCreate, args []string) error {
	return pinMessageCommandHandler(s, m, args)
}

func (h *CommandHandler) pinSlash(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	return pinSlashCommandHandler(s, i)
}

func (h *CommandHandler) unpin(s *discordgo.Session, m *discordgo.MessageCreate, args []string) error {
	return unpinMessageCommandHandler(s, m, args)
}

func (h *CommandHandler) unpinSlash(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	return unpinSlashCommandHandler(s, i)
}

func (h *CommandHandler) quote(s *discordgo.Session, m *discordgo.MessageCreate, args []string) error {
	return quoteMessageCommandHandler(s, m, args)
}

func (h *CommandHandler) quoteSlash(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	return quoteSlashCommandHandler(s, i)
}

func (h *CommandHandler) translate(s *discordgo.Session, m *discordgo.MessageCreate, args []string) error {
	return translateMessageCommandHandler(s, m, args)
}

func (h *CommandHandler) translateSlash(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	return translateSlashCommandHandler(s, i)
}

func (h *CommandHandler) reminder(s *discordgo.Session, m *discordgo.MessageCreate, args []string) error {
	return reminderMessageCommandHandler(s, m, args)
}

func (h *CommandHandler) reminderSlash(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	return reminderSlashCommandHandler(s, i)
}

func (h *CommandHandler) isearch(s *discordgo.Session, channelID string, args []string) error {
	return isearchMessageCommandHandler(s, channelID, args)
}

func (h *CommandHandler) isearchSlash(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	return isearchSlashCommandHandler(s, i)
}

func (h *CommandHandler) caption(s *discordgo.Session, m *discordgo.MessageCreate, args []string) error {
	return captionMessageCommandHandler(s, m, args)
}

func (h *CommandHandler) captionSlash(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	return captionSlashCommandHandler(s, i)
}

func (h *CommandHandler) img2gif(s *discordgo.Session, m *discordgo.MessageCreate, args []string) error {
	return img2gifMessageCommandHandler(s, m, args)
}

func (h *CommandHandler) img2gifSlash(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	return img2gifSlashCommandHandler(s, i)
}

func (h *CommandHandler) vid2gif(s *discordgo.Session, m *discordgo.MessageCreate, args []string) error {
	return vid2gifMessageCommandHandler(s, m, args)
}

func (h *CommandHandler) vid2gifSlash(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	return vid2gifSlashCommandHandler(s, i)
}

func (h *CommandHandler) autogif(s *discordgo.Session, m *discordgo.MessageCreate, args []string) error {
	return autogifMessageCommandHandler(s, m, args)
}

func (h *CommandHandler) autogifSlash(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	return autogifSlashCommandHandler(s, i)
}

func (h *CommandHandler) factcheck(s *discordgo.Session, m *discordgo.MessageCreate, args []string) error {
	return factcheckMessageCommandHandler(s, m, args, h.Cfg)
}

func (h *CommandHandler) factcheckSlash(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	return factcheckSlashCommandHandler(s, i, h.Cfg)
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
	{
		Name:        "pinglist",
		Description: "Returns all the pings for a certain user/role",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionUser,
				Name:        "user",
				Description: "User to show pings for",
				Required:    false,
			},
			{
				Type:        discordgo.ApplicationCommandOptionRole,
				Name:        "role",
				Description: "Role to show pings for",
				Required:    false,
			},
		},
	},
	{
		Name:        "gifsearch",
		Description: "Searches Giphy and posts a relevant GIF",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "query",
				Description: "The GIF search query",
				Required:    true,
			},
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "creator",
				Description: "Only show GIFs by a specific creator",
				Required:    false,
			},
		},
	},
	{
		Name:        "emoji",
		Description: "Manage server emojis",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionSubCommand,
				Name:        "add",
				Description: "Upload an emoji to the server",
				Options: []*discordgo.ApplicationCommandOption{
					{Type: discordgo.ApplicationCommandOptionString, Name: "name", Description: "Emoji name", Required: true},
					{Type: discordgo.ApplicationCommandOptionString, Name: "url", Description: "Emoji image URL", Required: false},
				},
			},
			{
				Type:        discordgo.ApplicationCommandOptionSubCommand,
				Name:        "enlarge",
				Description: "Enlarge an emoji",
				Options: []*discordgo.ApplicationCommandOption{
					{Type: discordgo.ApplicationCommandOptionString, Name: "emoji", Description: "The emoji to enlarge", Required: true},
				},
			},
			{
				Type:        discordgo.ApplicationCommandOptionSubCommand,
				Name:        "list",
				Description: "List all server emojis",
			},
			{
				Type:        discordgo.ApplicationCommandOptionSubCommand,
				Name:        "remove",
				Description: "Remove a server emoji",
				Options: []*discordgo.ApplicationCommandOption{
					{Type: discordgo.ApplicationCommandOptionString, Name: "emoji", Description: "Emoji to remove", Required: true},
				},
			},
			{
				Type:        discordgo.ApplicationCommandOptionSubCommand,
				Name:        "steal",
				Description: "Copy an emoji into this server",
				Options: []*discordgo.ApplicationCommandOption{
					{Type: discordgo.ApplicationCommandOptionString, Name: "emoji", Description: "Emoji to steal", Required: true},
				},
			},
		},
	},
	{
		Name:        "sticker",
		Description: "Manage server stickers",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionSubCommand,
				Name:        "add",
				Description: "Upload a sticker to the server",
				Options: []*discordgo.ApplicationCommandOption{
					{Type: discordgo.ApplicationCommandOptionString, Name: "name", Description: "Sticker name", Required: true},
					{Type: discordgo.ApplicationCommandOptionString, Name: "url", Description: "Sticker image URL", Required: true},
				},
			},
			{
				Type:        discordgo.ApplicationCommandOptionSubCommand,
				Name:        "remove",
				Description: "Remove a server sticker",
				Options: []*discordgo.ApplicationCommandOption{
					{Type: discordgo.ApplicationCommandOptionString, Name: "sticker_id", Description: "Sticker ID or message link", Required: true},
				},
			},
		},
	},
	{
		Name:        "pin",
		Description: "Pin a message",
		Options: []*discordgo.ApplicationCommandOption{
			{Type: discordgo.ApplicationCommandOptionString, Name: "message", Description: "Message ID or link", Required: true},
			{Type: discordgo.ApplicationCommandOptionChannel, Name: "channel", Description: "Channel (defaults to this one)", Required: false},
		},
	},
	{
		Name:        "unpin",
		Description: "Unpin a message",
		Options: []*discordgo.ApplicationCommandOption{
			{Type: discordgo.ApplicationCommandOptionString, Name: "message", Description: "Message ID or link", Required: true},
			{Type: discordgo.ApplicationCommandOptionChannel, Name: "channel", Description: "Channel (defaults to this one)", Required: false},
		},
	},
	{
		Name:        "quote",
		Description: "Quote a message as a styled embed",
		Options: []*discordgo.ApplicationCommandOption{
			{Type: discordgo.ApplicationCommandOptionString, Name: "message", Description: "Message ID or link", Required: true},
			{Type: discordgo.ApplicationCommandOptionChannel, Name: "channel", Description: "Channel (defaults to this one)", Required: false},
		},
	},
	{
		Name:        "translate",
		Description: "Translate text into a target language",
		Options: []*discordgo.ApplicationCommandOption{
			{Type: discordgo.ApplicationCommandOptionString, Name: "text", Description: "Text to translate", Required: true},
			{Type: discordgo.ApplicationCommandOptionString, Name: "target", Description: "Target language code (default: en)", Required: false},
		},
	},
	{
		Name:        "reminder",
		Description: "Set a delayed reminder ping",
		Options: []*discordgo.ApplicationCommandOption{
			{Type: discordgo.ApplicationCommandOptionString, Name: "text", Description: "What to remind you about", Required: true},
			{Type: discordgo.ApplicationCommandOptionString, Name: "delay", Description: "When, e.g. 30m / 2h / 1h30m", Required: true},
		},
	},
	{
		Name:        "isearch",
		Description: "Search the web for images",
		Options: []*discordgo.ApplicationCommandOption{
			{Type: discordgo.ApplicationCommandOptionString, Name: "query", Description: "Image search query", Required: true},
		},
	},
	{
		Name:        "caption",
		Description: "Add meme text to an image",
		Options: []*discordgo.ApplicationCommandOption{
			{Type: discordgo.ApplicationCommandOptionString, Name: "top", Description: "Top text", Required: false},
			{Type: discordgo.ApplicationCommandOptionString, Name: "bottom", Description: "Bottom text", Required: false},
			{Type: discordgo.ApplicationCommandOptionString, Name: "image_url", Description: "Background image URL", Required: false},
		},
	},
	{
		Name:        "img2gif",
		Description: "Convert an image into a GIF",
		Options: []*discordgo.ApplicationCommandOption{
			{Type: discordgo.ApplicationCommandOptionString, Name: "url", Description: "Image URL", Required: true},
		},
	},
	{
		Name:        "vid2gif",
		Description: "Convert a video into a GIF",
		Options: []*discordgo.ApplicationCommandOption{
			{Type: discordgo.ApplicationCommandOptionString, Name: "url", Description: "Video URL", Required: true},
		},
	},
	{
		Name:        "autogif",
		Description: "Convert any media into a GIF",
		Options: []*discordgo.ApplicationCommandOption{
			{Type: discordgo.ApplicationCommandOptionString, Name: "url", Description: "Image or video URL", Required: true},
		},
	},
	{
		Name:        "factcheck",
		Description: "Fact-check a claim against searchable ratings",
		Options: []*discordgo.ApplicationCommandOption{
			{Type: discordgo.ApplicationCommandOptionString, Name: "claim", Description: "The claim to fact-check", Required: true},
		},
	},
}
