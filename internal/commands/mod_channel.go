package commands

import (
	"context"
	"fmt"
	"strconv"

	"github.com/bwmarrin/discordgo"
)

const modChannelPerm = discordgo.PermissionManageChannels

// hideMessage: -hide  (deny @everyone View Channel in the current channel)
func hideMessageCommandHandler(h *CommandHandler, s *discordgo.Session, m *discordgo.MessageCreate, args []string) error {
	if ok, err := requireModForMessage(s, m, "Hide"); !ok {
		return err
	}
	if err := setEveryoneView(s, m.GuildID, m.ChannelID, false); err != nil {
		return sendModError(s, m.ChannelID, "Hide", fmt.Sprintf("Failed to hide channel: %s", err))
	}
	auditAction(h, context.Background(), m.GuildID, "HIDE", m.Author.ID, m.ChannelID, nil)
	_, err := s.ChannelMessageSendEmbed(m.ChannelID, modSuccessEmbed("Hide", "This channel is now hidden from @everyone."))
	return err
}

func hideSlashCommandHandler(h *CommandHandler, s *discordgo.Session, i *discordgo.InteractionCreate) error {
	ok, msg := requireModerator(s, i.GuildID, i.Member.User.ID, i.Member)
	if !ok {
		return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{Embeds: []*discordgo.MessageEmbed{modErrorEmbed("Hide", msg)}},
		})
	}
	chID := i.ChannelID
	if err := setEveryoneView(s, i.GuildID, chID, false); err != nil {
		return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{Embeds: []*discordgo.MessageEmbed{modErrorEmbed("Hide", fmt.Sprintf("Failed to hide channel: %s", err))}},
		})
	}
	auditAction(h, context.Background(), i.GuildID, "HIDE", i.Member.User.ID, chID, nil)
	return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{Embeds: []*discordgo.MessageEmbed{modSuccessEmbed("Hide", "This channel is now hidden from @everyone.")}},
	})
}

// revealMessage: -reveal  (restore @everyone View Channel in the current channel)
func revealMessageCommandHandler(h *CommandHandler, s *discordgo.Session, m *discordgo.MessageCreate, args []string) error {
	if ok, err := requireModForMessage(s, m, "Reveal"); !ok {
		return err
	}
	if err := setEveryoneView(s, m.GuildID, m.ChannelID, true); err != nil {
		return sendModError(s, m.ChannelID, "Reveal", fmt.Sprintf("Failed to reveal channel: %s", err))
	}
	auditAction(h, context.Background(), m.GuildID, "REVEAL", m.Author.ID, m.ChannelID, nil)
	_, err := s.ChannelMessageSendEmbed(m.ChannelID, modSuccessEmbed("Reveal", "This channel is now visible to @everyone."))
	return err
}

func revealSlashCommandHandler(h *CommandHandler, s *discordgo.Session, i *discordgo.InteractionCreate) error {
	ok, msg := requireModerator(s, i.GuildID, i.Member.User.ID, i.Member)
	if !ok {
		return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{Embeds: []*discordgo.MessageEmbed{modErrorEmbed("Reveal", msg)}},
		})
	}
	chID := i.ChannelID
	if err := setEveryoneView(s, i.GuildID, chID, true); err != nil {
		return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{Embeds: []*discordgo.MessageEmbed{modErrorEmbed("Reveal", fmt.Sprintf("Failed to reveal channel: %s", err))}},
		})
	}
	auditAction(h, context.Background(), i.GuildID, "REVEAL", i.Member.User.ID, chID, nil)
	return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{Embeds: []*discordgo.MessageEmbed{modSuccessEmbed("Reveal", "This channel is now visible to @everyone.")}},
	})
}

func setEveryoneView(s *discordgo.Session, guildID, channelID string, allow bool) error {
	perm := &discordgo.PermissionOverwrite{
		ID:   guildID,
		Type: discordgo.PermissionOverwriteTypeRole,
	}
	if allow {
		perm.Allow = discordgo.PermissionViewChannel
		perm.Deny = 0
	} else {
		perm.Deny = discordgo.PermissionViewChannel
		perm.Allow = 0
	}
	err := s.ChannelPermissionSet(channelID, guildID, discordgo.PermissionOverwriteTypeRole, perm.Allow, perm.Deny)
	return err
}

// lockdownMessage: -lockdown  (deny @everyone Send Messages in the current channel)
func lockdownMessageCommandHandler(h *CommandHandler, s *discordgo.Session, m *discordgo.MessageCreate, args []string) error {
	if ok, err := requireModForMessage(s, m, "Lockdown"); !ok {
		return err
	}
	if err := setEveryoneSend(s, m.GuildID, m.ChannelID, false); err != nil {
		return sendModError(s, m.ChannelID, "Lockdown", fmt.Sprintf("Failed to lock the channel: %s", err))
	}
	auditAction(h, context.Background(), m.GuildID, "LOCKDOWN", m.Author.ID, m.ChannelID, nil)
	_, err := s.ChannelMessageSendEmbed(m.ChannelID, modSuccessEmbed("Lockdown", "This channel is now locked."))
	return err
}

func lockdownSlashCommandHandler(h *CommandHandler, s *discordgo.Session, i *discordgo.InteractionCreate) error {
	ok, msg := requireModerator(s, i.GuildID, i.Member.User.ID, i.Member)
	if !ok {
		return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{Embeds: []*discordgo.MessageEmbed{modErrorEmbed("Lockdown", msg)}},
		})
	}
	chID := i.ChannelID
	if err := setEveryoneSend(s, i.GuildID, chID, false); err != nil {
		return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{Embeds: []*discordgo.MessageEmbed{modErrorEmbed("Lockdown", fmt.Sprintf("Failed to lock the channel: %s", err))}},
		})
	}
	auditAction(h, context.Background(), i.GuildID, "LOCKDOWN", i.Member.User.ID, chID, nil)
	return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{Embeds: []*discordgo.MessageEmbed{modSuccessEmbed("Lockdown", "This channel is now locked.")}},
	})
}

func setEveryoneSend(s *discordgo.Session, guildID, channelID string, allow bool) error {
	perm := &discordgo.PermissionOverwrite{ID: guildID, Type: discordgo.PermissionOverwriteTypeRole}
	if allow {
		perm.Allow = discordgo.PermissionSendMessages
		perm.Deny = 0
	} else {
		perm.Deny = discordgo.PermissionSendMessages
		perm.Allow = 0
	}
	err := s.ChannelPermissionSet(channelID, guildID, discordgo.PermissionOverwriteTypeRole, perm.Allow, perm.Deny)
	return err
}

// nsfwMessage: -nsfw  (mark the current channel as NSFW)
func nsfwMessageCommandHandler(h *CommandHandler, s *discordgo.Session, m *discordgo.MessageCreate, args []string) error {
	if ok, err := requireModForMessage(s, m, "NSFW"); !ok {
		return err
	}
	if _, err := s.ChannelEditComplex(m.ChannelID, &discordgo.ChannelEdit{NSFW: boolPtr(true)}); err != nil {
		return sendModError(s, m.ChannelID, "NSFW", fmt.Sprintf("Failed to update channel: %s", err))
	}
	auditAction(h, context.Background(), m.GuildID, "NSFW", m.Author.ID, m.ChannelID, nil)
	_, err := s.ChannelMessageSendEmbed(m.ChannelID, modSuccessEmbed("NSFW", "This channel is now marked as NSFW."))
	return err
}

func nsfwSlashCommandHandler(h *CommandHandler, s *discordgo.Session, i *discordgo.InteractionCreate) error {
	ok, msg := requireModerator(s, i.GuildID, i.Member.User.ID, i.Member)
	if !ok {
		return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{Embeds: []*discordgo.MessageEmbed{modErrorEmbed("NSFW", msg)}},
		})
	}
	chID := i.ChannelID
	if _, err := s.ChannelEditComplex(chID, &discordgo.ChannelEdit{NSFW: boolPtr(true)}); err != nil {
		return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{Embeds: []*discordgo.MessageEmbed{modErrorEmbed("NSFW", fmt.Sprintf("Failed to update channel: %s", err))}},
		})
	}
	auditAction(h, context.Background(), i.GuildID, "NSFW", i.Member.User.ID, chID, nil)
	return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{Embeds: []*discordgo.MessageEmbed{modSuccessEmbed("NSFW", "This channel is now marked as NSFW.")}},
	})
}

// sfwMessage: -sfw  (remove the NSFW flag from the current channel)
func sfwMessageCommandHandler(h *CommandHandler, s *discordgo.Session, m *discordgo.MessageCreate, args []string) error {
	if ok, err := requireModForMessage(s, m, "SFW"); !ok {
		return err
	}
	if _, err := s.ChannelEditComplex(m.ChannelID, &discordgo.ChannelEdit{NSFW: boolPtr(false)}); err != nil {
		return sendModError(s, m.ChannelID, "SFW", fmt.Sprintf("Failed to update channel: %s", err))
	}
	auditAction(h, context.Background(), m.GuildID, "SFW", m.Author.ID, m.ChannelID, nil)
	_, err := s.ChannelMessageSendEmbed(m.ChannelID, modSuccessEmbed("SFW", "This channel is no longer marked as NSFW."))
	return err
}

func sfwSlashCommandHandler(h *CommandHandler, s *discordgo.Session, i *discordgo.InteractionCreate) error {
	ok, msg := requireModerator(s, i.GuildID, i.Member.User.ID, i.Member)
	if !ok {
		return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{Embeds: []*discordgo.MessageEmbed{modErrorEmbed("SFW", msg)}},
		})
	}
	chID := i.ChannelID
	if _, err := s.ChannelEditComplex(chID, &discordgo.ChannelEdit{NSFW: boolPtr(false)}); err != nil {
		return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{Embeds: []*discordgo.MessageEmbed{modErrorEmbed("SFW", fmt.Sprintf("Failed to update channel: %s", err))}},
		})
	}
	auditAction(h, context.Background(), i.GuildID, "SFW", i.Member.User.ID, chID, nil)
	return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{Embeds: []*discordgo.MessageEmbed{modSuccessEmbed("SFW", "This channel is no longer marked as NSFW.")}},
	})
}

// slowmodeMessage: -slowmode <seconds>  (set the channel rate limit, 0 to clear)
func slowmodeMessageCommandHandler(h *CommandHandler, s *discordgo.Session, m *discordgo.MessageCreate, args []string) error {
	if ok, err := requireModForMessage(s, m, "Slowmode"); !ok {
		return err
	}
	if len(args) == 0 {
		return sendModError(s, m.ChannelID, "Slowmode", "Usage: `-slowmode <seconds>` (0 to disable, max 21600)")
	}
	secs, err := strconv.Atoi(args[0])
	if err != nil || secs < 0 {
		return sendModError(s, m.ChannelID, "Slowmode", "Seconds must be a whole number >= 0.")
	}
	if secs > 21600 {
		secs = 21600
	}
	if _, err := s.ChannelEditComplex(m.ChannelID, &discordgo.ChannelEdit{RateLimitPerUser: &secs}); err != nil {
		return sendModError(s, m.ChannelID, "Slowmode", fmt.Sprintf("Failed to set slowmode: %s", err))
	}
	auditAction(h, context.Background(), m.GuildID, "SLOWMODE", m.Author.ID, m.ChannelID, map[string]any{"seconds": secs})
	_, err = s.ChannelMessageSendEmbed(m.ChannelID, modSuccessEmbed("Slowmode", fmt.Sprintf("Slowmode set to **%d** second%s.", secs, plural(secs))))
	return err
}

func slowmodeSlashCommandHandler(h *CommandHandler, s *discordgo.Session, i *discordgo.InteractionCreate) error {
	ok, msg := requireModerator(s, i.GuildID, i.Member.User.ID, i.Member)
	if !ok {
		return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{Embeds: []*discordgo.MessageEmbed{modErrorEmbed("Slowmode", msg)}},
		})
	}
	opts := OptionMap(i.ApplicationCommandData().Options)
	secs := OptInt(opts, "seconds", 0)
	if secs < 0 {
		secs = 0
	}
	if secs > 21600 {
		secs = 21600
	}
	chID := i.ChannelID
	if _, err := s.ChannelEditComplex(chID, &discordgo.ChannelEdit{RateLimitPerUser: &secs}); err != nil {
		return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{Embeds: []*discordgo.MessageEmbed{modErrorEmbed("Slowmode", fmt.Sprintf("Failed to set slowmode: %s", err))}},
		})
	}
	auditAction(h, context.Background(), i.GuildID, "SLOWMODE", i.Member.User.ID, chID, map[string]any{"seconds": secs})
	return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{Embeds: []*discordgo.MessageEmbed{modSuccessEmbed("Slowmode", fmt.Sprintf("Slowmode set to **%d** second%s.", secs, plural(secs)))}},
	})
}

// topicMessage: -topic <text>  (update the current channel topic)
func topicMessageCommandHandler(h *CommandHandler, s *discordgo.Session, m *discordgo.MessageCreate, args []string) error {
	if ok, err := requireModForMessage(s, m, "Topic"); !ok {
		return err
	}
	topic := joinArgs(args)
	if _, err := s.ChannelEditComplex(m.ChannelID, &discordgo.ChannelEdit{Topic: topic}); err != nil {
		return sendModError(s, m.ChannelID, "Topic", fmt.Sprintf("Failed to update topic: %s", err))
	}
	auditAction(h, context.Background(), m.GuildID, "TOPIC", m.Author.ID, m.ChannelID, map[string]any{"topic": topic})
	_, err := s.ChannelMessageSendEmbed(m.ChannelID, modSuccessEmbed("Topic", fmt.Sprintf("Topic updated to:\n%s", topic)))
	return err
}

func topicSlashCommandHandler(h *CommandHandler, s *discordgo.Session, i *discordgo.InteractionCreate) error {
	ok, msg := requireModerator(s, i.GuildID, i.Member.User.ID, i.Member)
	if !ok {
		return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{Embeds: []*discordgo.MessageEmbed{modErrorEmbed("Topic", msg)}},
		})
	}
	opts := OptionMap(i.ApplicationCommandData().Options)
	topic := OptString(opts, "topic")
	chID := i.ChannelID
	if _, err := s.ChannelEditComplex(chID, &discordgo.ChannelEdit{Topic: topic}); err != nil {
		return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{Embeds: []*discordgo.MessageEmbed{modErrorEmbed("Topic", fmt.Sprintf("Failed to update topic: %s", err))}},
		})
	}
	auditAction(h, context.Background(), i.GuildID, "TOPIC", i.Member.User.ID, chID, map[string]any{"topic": topic})
	return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{Embeds: []*discordgo.MessageEmbed{modSuccessEmbed("Topic", fmt.Sprintf("Topic updated to:\n%s", topic))}},
	})
}

func boolPtr(b bool) *bool { return &b }
