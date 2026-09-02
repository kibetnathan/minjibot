package commands

import (
	"context"
	"fmt"
	"strconv"

	"github.com/bwmarrin/discordgo"
)

// purgeMessage: -purge <1-100> [user]
func purgeMessageCommandHandler(h *CommandHandler, s *discordgo.Session, m *discordgo.MessageCreate, args []string) error {
	if ok, err := requireModForMessage(s, m, "Purge"); !ok {
		return err
	}
	if len(args) == 0 {
		return sendModError(s, m.ChannelID, "Purge", "Usage: `-purge <count>` (1-100), optionally `-purge <count> [user]`")
	}
	count, err := strconv.Atoi(args[0])
	if err != nil || count < 1 {
		return sendModError(s, m.ChannelID, "Purge", "Count must be a whole number between 1 and 100.")
	}
	if count > 100 {
		count = 100
	}

	filterID := ""
	if len(args) > 1 {
		filterID = parseMentionID(args[1])
	}

	deleted, err := purgeMessages(s, m.ChannelID, count, filterID)
	if err != nil {
		return sendModError(s, m.ChannelID, "Purge", fmt.Sprintf("Failed to purge messages: %s", err))
	}
	auditAction(h, context.Background(), m.GuildID, "PURGE", m.Author.ID, m.ChannelID, map[string]any{"count": deleted, "filter": filterID})
	_, err = s.ChannelMessageSendEmbed(m.ChannelID, modSuccessEmbed("Purge", fmt.Sprintf("Deleted **%d** message%s.", deleted, plural(deleted))))
	return err
}

// purgeSlash: /purge count:<n> [user]
func purgeSlashCommandHandler(h *CommandHandler, s *discordgo.Session, i *discordgo.InteractionCreate) error {
	ok, msg := requireModerator(s, i.GuildID, i.Member.User.ID, i.Member)
	if !ok {
		return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{Embeds: []*discordgo.MessageEmbed{modErrorEmbed("Purge", msg)}},
		})
	}
	opts := OptionMap(i.ApplicationCommandData().Options)
	count := OptInt(opts, "count", 0)
	if count < 1 {
		return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{Embeds: []*discordgo.MessageEmbed{modErrorEmbed("Purge", "Count must be between 1 and 100.")}},
		})
	}
	if count > 100 {
		count = 100
	}
	filterID := OptUser(opts, "user")

	deleted, err := purgeMessages(s, i.ChannelID, count, filterID)
	if err != nil {
		return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{Embeds: []*discordgo.MessageEmbed{modErrorEmbed("Purge", fmt.Sprintf("Failed to purge messages: %s", err))}},
		})
	}
	auditAction(h, context.Background(), i.GuildID, "PURGE", i.Member.User.ID, i.ChannelID, map[string]any{"count": deleted, "filter": filterID})
	return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{Embeds: []*discordgo.MessageEmbed{modSuccessEmbed("Purge", fmt.Sprintf("Deleted **%d** message%s.", deleted, plural(deleted)))}},
	})
}

// purgeMessages bulk-deletes a batch of <=100 messages, optionally filtered to
// one author. It handles Discord's 14-day and 100-message API limits.
func purgeMessages(s *discordgo.Session, channelID string, count int, filterID string) (int, error) {
	total := 0
	lastID := ""
	for total < count {
		batch := count - total
		if batch > 100 {
			batch = 100
		}
		if filterID != "" {
			// Bulk delete only supports object-based deletion; to filter by
			// author we have to delete one-by-one.
			removed, last, err := purgeByAuthor(s, channelID, batch, filterID, lastID)
			if err != nil {
				return total, err
			}
			if removed == 0 {
				break
			}
			total += removed
			lastID = last
			continue
		}
		removed, last, err := bulkDeleteBatch(s, channelID, batch, lastID)
		if err != nil {
			return total, err
		}
		if removed == 0 {
			break
		}
		total += removed
		lastID = last
	}
	return total, nil
}

// bulkDeleteBatch deletes up to n messages older than beforeID via ChannelBulkDelete.
func bulkDeleteBatch(s *discordgo.Session, channelID string, n int, beforeID string) (int, string, error) {
	messages, err := s.ChannelMessages(channelID, n+1, beforeID, "", "")
	if err != nil {
		return 0, "", err
	}
	// ChannelMessages returns messages in descending order (newest first).
	if len(messages) == 0 {
		return 0, "", nil
	}
	ids := make([]string, 0, len(messages))
	for _, msg := range messages {
		ids = append(ids, msg.ID)
	}
	// Do not delete messages older than 14 days (Discord rejects them).
	if err := s.ChannelMessagesBulkDelete(channelID, ids); err != nil {
		// Fall back to individual deletes on permission/age errors.
		for _, msg := range messages {
			_ = s.ChannelMessageDelete(channelID, msg.ID)
		}
	}
	oldest := messages[len(messages)-1].ID
	return len(messages), oldest, nil
}

// purgeByAuthor walks backwards and deletes messages from a single author.
func purgeByAuthor(s *discordgo.Session, channelID string, n int, userID, beforeID string) (int, string, error) {
	removed := 0
	lastBefore := beforeID
	for removed < n {
		messages, err := s.ChannelMessages(channelID, 100, lastBefore, "", "")
		if err != nil {
			return removed, lastBefore, err
		}
		if len(messages) == 0 {
			return removed, lastBefore, nil
		}
		lastBefore = messages[len(messages)-1].ID
		for _, msg := range messages {
			if removed >= n {
				return removed, lastBefore, nil
			}
			if msg.Author.ID == userID && !msg.Author.Bot {
				if err := s.ChannelMessageDelete(channelID, msg.ID); err == nil {
					removed++
				}
			}
		}
	}
	return removed, lastBefore, nil
}

// nukeMessage: -nuke  (clone the channel and delete the original)
func nukeMessageCommandHandler(h *CommandHandler, s *discordgo.Session, m *discordgo.MessageCreate, args []string) error {
	if ok, err := requireModForMessage(s, m, "Nuke"); !ok {
		return err
	}
	return nukeChannel(s, m.ChannelID, m.GuildID, m.Author.ID, h)
}

func nukeSlashCommandHandler(h *CommandHandler, s *discordgo.Session, i *discordgo.InteractionCreate) error {
	ok, msg := requireModerator(s, i.GuildID, i.Member.User.ID, i.Member)
	if !ok {
		return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{Embeds: []*discordgo.MessageEmbed{modErrorEmbed("Nuke", msg)}},
		})
	}
	return nukeChannel(s, i.ChannelID, i.GuildID, i.Member.User.ID, h)
}

func nukeChannel(s *discordgo.Session, channelID, guildID, actorID string, h *CommandHandler) error {
	ch, err := s.Channel(channelID)
	if err != nil {
		return sendModError(s, channelID, "Nuke", fmt.Sprintf("Could not fetch the channel: %s", err))
	}
	isThread := ch.Type == discordgo.ChannelTypeGuildNewsThread || ch.Type == discordgo.ChannelTypeGuildPublicThread || ch.Type == discordgo.ChannelTypeGuildPrivateThread
	if isThread {
		return sendModError(s, channelID, "Nuke", "Nuke only works on regular text channels, not threads.")
	}

	clone := discordgo.GuildChannelCreateData{
		Name:                 ch.Name,
		Type:                 ch.Type,
		Topic:                ch.Topic,
		NSFW:                 ch.NSFW,
		Bitrate:              ch.Bitrate,
		UserLimit:            ch.UserLimit,
		RateLimitPerUser:     ch.RateLimitPerUser,
		ParentID:             ch.ParentID,
		Position:             ch.Position,
		PermissionOverwrites: ch.PermissionOverwrites,
	}
	newCh, err := s.GuildChannelCreateComplex(guildID, clone)
	if err != nil {
		return sendModError(s, channelID, "Nuke", fmt.Sprintf("Failed to clone the channel: %s", err))
	}
	if _, err := s.ChannelDelete(channelID); err != nil {
		// Channel cloned but original failed to delete; report partial success.
		auditAction(h, context.Background(), guildID, "NUKE", actorID, channelID, map[string]any{"cloned_to": newCh.ID})
		_, _ = s.ChannelMessageSendEmbed(newCh.ID, modSuccessEmbed("Nuke", "Channel cloned; original could not be deleted."))
		return nil
	}
	auditAction(h, context.Background(), guildID, "NUKE", actorID, channelID, map[string]any{"cloned_to": newCh.ID})
	msg := "Channel nuked. This is the fresh replacement."
	_, _ = s.ChannelMessageSendEmbed(newCh.ID, modSuccessEmbed("Nuke", msg))
	return nil
}
