package commands

import (
	"context"
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/kibetnathan/minjibot/internal/ports/dto"
)

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
		return sendModError(s, m.ChannelID, "Setup", "Usage: `-setup logchannel <#channel>` or `-setup status`")
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
		})
		if err != nil {
			return sendModError(s, m.ChannelID, "Setup", "Failed to save logging channel.")
		}
		_, err = s.ChannelMessageSendEmbed(m.ChannelID, modSuccessEmbed("Setup", fmt.Sprintf("Logging channel set to <#%s>. Deleted messages and moderation actions will be posted there.", channel.ID)))
		return err
	case "status", "show":
		val := "Not configured — deleted messages and actions are stored in the database only."
		if current.LoggingChannelID != "" {
			val = fmt.Sprintf("<#%s>", current.LoggingChannelID)
		}
		_, err := s.ChannelMessageSendEmbed(m.ChannelID, modSuccessEmbed("Setup", fmt.Sprintf("Logging channel: %s", val)))
		return err
	default:
		return sendModError(s, m.ChannelID, "Setup", "Usage: `-setup logchannel <#channel>` or `-setup status`")
	}
}

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
		return respondError("Usage: `/setup logchannel  <channel>` or `/setup status`")
	}
	sub := data.Options[0]
	name := sub.Name
	args := sub.Options

	current, err := h.SettingsRepo.Get(context.Background(), i.GuildID)
	if err != nil {
		return respondError("Could not read server settings.")
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
		})
		if err != nil {
			return respondError("Failed to save logging channel.")
		}
		return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{Embeds: []*discordgo.MessageEmbed{modSuccessEmbed("Setup", fmt.Sprintf("Logging channel set to <#%s>. Deleted messages and moderation actions will be posted there.", channel.ID))}},
		})
	case "status":
		val := "Not configured — deleted messages and actions are stored in the database only."
		if current.LoggingChannelID != "" {
			val = fmt.Sprintf("<#%s>", current.LoggingChannelID)
		}
		return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{Embeds: []*discordgo.MessageEmbed{modSuccessEmbed("Setup", fmt.Sprintf("Logging channel: %s", val))}},
		})
	default:
		return respondError("Usage: `/setup logchannel  <channel>` or `/setup status`")
	}
}
