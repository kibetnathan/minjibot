package commands

import (
	"context"
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/kibetnathan/minjibot/internal/ports/dto"
)

// parseToggle interprets a human on/off value into a boolean. The second return
// reports whether the input was recognised.
func parseToggle(v string) (value bool, ok bool) {
	switch strings.ToLower(strings.TrimSpace(v)) {
	case "on", "enable", "enabled", "true", "yes", "1":
		return true, true
	case "off", "disable", "disabled", "false", "no", "0":
		return false, true
	}
	return false, false
}

// isOwnerOrAdminPrefix is the prefix-command gate for /setup. It computes the
// user's server-wide permissions via the Discord API and requires the
// Administrator permission.
func isOwnerOrAdminPrefix(s *discordgo.Session, guildID, userID string, channelID string) (bool, error) {
	perms, err := effectiveModPerms(s, guildID, userID)
	if err != nil {
		return false, sendModError(s, channelID, "Setup", fmt.Sprintf("Could not check permissions: %s", err))
	}
	if perms&discordgo.PermissionAdministrator != 0 {
		return true, nil
	}
	return false, sendModError(s, channelID, "Setup", "You need the Administrator permission to use this command.")
}

func isOwnerOrAdminSlash(m *discordgo.Member) (bool, string) {
	if m == nil || m.Permissions&discordgo.PermissionAdministrator == 0 {
		return false, "You need the Administrator permission to use this command."
	}
	return true, ""
}

// setupMessageCommandHandler lets a guild owner/administrator point logging at
// a channel:
//
//	-setup logchannel <#channel|channelid>   set where deleted messages & mod actions are posted
//	-setup status                            show the current logging configuration
func setupMessageCommandHandler(h *CommandHandler, s *discordgo.Session, m *discordgo.MessageCreate, args []string) error {
	ok, err := isOwnerOrAdminPrefix(s, m.GuildID, m.Author.ID, m.ChannelID)
	if !ok {
		return err
	}
	if len(args) == 0 {
		return sendModError(s, m.ChannelID, "Setup", "Usage: `-setup logchannel <#channel>`, `-setup messagelogging <on|off>`, or `-setup status`")
	}

	current, err := h.SettingsRepo.Get(context.Background(), m.GuildID)
	if err != nil {
		return sendModError(s, m.ChannelID, "Setup", "Could not read server settings.")
	}

	switch args[0] {
	case "logchannel":
		if len(args) < 2 {
			return sendModError(s, m.ChannelID, "Setup", "Usage: `-setup logchannel <#channel>`")
		}
		raw := strings.Join(args[1:], " ")
		channelID := parseMentionID(raw)
		channel, err := s.Channel(channelID)
		if err != nil {
			return sendModError(s, m.ChannelID, "Setup", fmt.Sprintf("Could not find that channel: %s", err))
		}
		_, err = h.SettingsRepo.Upsert(context.Background(), dto.UpsertGuildSettingsParams{
			GuildID:               m.GuildID,
			Prefix:                current.Prefix,
			Language:              current.Language,
			AutoModerationEnabled: current.AutoModerationEnabled,
			LoggingChannelID:      channel.ID,
			MessageLoggingEnabled: current.MessageLoggingEnabled,
		})
		if err != nil {
			return sendModError(s, m.ChannelID, "Setup", "Failed to save logging channel.")
		}
		_, err = s.ChannelMessageSendEmbed(m.ChannelID, modSuccessEmbed("Setup", fmt.Sprintf("Logging channel set to <#%s>. Deleted messages and moderation actions will be posted there.", channel.ID)))
		return err
	case "messagelogging", "messagelog":
		if len(args) < 2 {
			return sendModError(s, m.ChannelID, "Setup", "Usage: `-setup messagelogging <on|off>`")
		}
		enabled, ok := parseToggle(args[1])
		if !ok {
			return sendModError(s, m.ChannelID, "Setup", "Please specify `on` or `off`.")
		}
		_, err = h.SettingsRepo.Upsert(context.Background(), dto.UpsertGuildSettingsParams{
			GuildID:               m.GuildID,
			Prefix:                current.Prefix,
			Language:              current.Language,
			AutoModerationEnabled: current.AutoModerationEnabled,
			LoggingChannelID:      current.LoggingChannelID,
			MessageLoggingEnabled: enabled,
		})
		if err != nil {
			return sendModError(s, m.ChannelID, "Setup", "Failed to save message logging setting.")
		}
		_, err = s.ChannelMessageSendEmbed(m.ChannelID, modSuccessEmbed("Setup", messageLoggingConfirm(enabled)))
		return err
	case "status", "show":
		_, err := s.ChannelMessageSendEmbed(m.ChannelID, modSuccessEmbed("Setup", setupStatusText(current.LoggingChannelID, current.MessageLoggingEnabled)))
		return err
	default:
		return sendModError(s, m.ChannelID, "Setup", "Usage: `-setup logchannel <#channel>`, `-setup messagelogging <on|off>`, or `-setup status`")
	}
}

// messageLoggingConfirm is the confirmation shown after toggling message logging.
func messageLoggingConfirm(enabled bool) string {
	if enabled {
		return fmt.Sprintf("Message-content logging is now **on**. Messages are stored for %d days, then automatically deleted.", int(messageLogRetentionDays))
	}
	return "Message-content logging is now **off**. New messages will not be stored."
}

// setupStatusText renders the current logging configuration.
func setupStatusText(loggingChannelID string, messageLoggingEnabled bool) string {
	channel := "Not configured — deleted messages and actions are stored in the database only."
	if loggingChannelID != "" {
		channel = fmt.Sprintf("<#%s>", loggingChannelID)
	}
	msgLog := "off"
	if messageLoggingEnabled {
		msgLog = "on"
	}
	return fmt.Sprintf("Logging channel: %s\nMessage-content logging: %s", channel, msgLog)
}

// messageLogRetentionDays mirrors the bot's retention window for display in
// setup confirmations. Kept in sync with bot.messageLogRetention.
const messageLogRetentionDays = 30

// setupSlashCommandHandler is the slash-command version of /setup. It exposes a
// subcommand for setting the log channel and one for viewing the status.
func setupSlashCommandHandler(h *CommandHandler, s *discordgo.Session, i *discordgo.InteractionCreate) error {
	ok, msg := isOwnerOrAdminSlash(i.Member)
	if !ok {
		return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{Embeds: []*discordgo.MessageEmbed{modErrorEmbed("Setup", msg)}},
		})
	}

	respondError := func(embedMsg string) error {
		return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{Embeds: []*discordgo.MessageEmbed{modErrorEmbed("Setup", embedMsg)}},
		})
	}

	data := i.ApplicationCommandData()
	if len(data.Options) == 0 {
		return respondError("Usage: `/setup logchannel <channel>`, `/setup messagelogging <enabled>`, or `/setup status`")
	}
	sub := data.Options[0]
	name := sub.Name
	args := sub.Options

	current, err := h.SettingsRepo.Get(context.Background(), i.GuildID)
	if err != nil {
		return respondError("Could not read server settings.")
	}

	respondOK := func(msg string) error {
		return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{Embeds: []*discordgo.MessageEmbed{modSuccessEmbed("Setup", msg)}},
		})
	}

	switch name {
	case "logchannel":
		chID := ""
		for _, a := range args {
			if a.Name == "channel" {
				chID = a.StringValue()
			}
		}
		if chID == "" {
			return respondError("Please pick a channel.")
		}
		channel, err := s.Channel(chID)
		if err != nil {
			return respondError(fmt.Sprintf("Could not find that channel: %s", err))
		}
		_, err = h.SettingsRepo.Upsert(context.Background(), dto.UpsertGuildSettingsParams{
			GuildID:               i.GuildID,
			Prefix:                current.Prefix,
			Language:              current.Language,
			AutoModerationEnabled: current.AutoModerationEnabled,
			LoggingChannelID:      channel.ID,
			MessageLoggingEnabled: current.MessageLoggingEnabled,
		})
		if err != nil {
			return respondError("Failed to save logging channel.")
		}
		return respondOK(fmt.Sprintf("Logging channel set to <#%s>. Deleted messages and moderation actions will be posted there.", channel.ID))
	case "messagelogging":
		enabled := false
		for _, a := range args {
			if a.Name == "enabled" {
				enabled = a.BoolValue()
			}
		}
		_, err = h.SettingsRepo.Upsert(context.Background(), dto.UpsertGuildSettingsParams{
			GuildID:               i.GuildID,
			Prefix:                current.Prefix,
			Language:              current.Language,
			AutoModerationEnabled: current.AutoModerationEnabled,
			LoggingChannelID:      current.LoggingChannelID,
			MessageLoggingEnabled: enabled,
		})
		if err != nil {
			return respondError("Failed to save message logging setting.")
		}
		return respondOK(messageLoggingConfirm(enabled))
	case "status":
		return respondOK(setupStatusText(current.LoggingChannelID, current.MessageLoggingEnabled))
	default:
		return respondError("Usage: `/setup logchannel <channel>`, `/setup messagelogging <enabled>`, or `/setup status`")
	}
}
