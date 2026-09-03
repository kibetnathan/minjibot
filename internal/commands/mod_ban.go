package commands

import (
	"context"
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
)

// banMessage: -ban <user> [reason]
func banMessageCommandHandler(h *CommandHandler, s *discordgo.Session, m *discordgo.MessageCreate, args []string) error {
	if ok, err := requireModForMessage(s, m, "Ban"); !ok {
		return err
	}
	if len(args) == 0 {
		return sendModError(s, m.ChannelID, "Ban", "Usage: `-ban <user> [reason]`")
	}
	targetID, name, err := resolveTargetUser(s, m.GuildID, args[0])
	if err != nil {
		return sendModError(s, m.ChannelID, "Ban", fmt.Sprintf("Could not find that user: %s", err))
	}
	reason := strings.Join(args[1:], " ")

	if err := s.GuildBanCreateWithReason(m.GuildID, targetID, reason, 0); err != nil {
		return sendModError(s, m.ChannelID, "Ban", fmt.Sprintf("Failed to ban %s: %s", name, err))
	}
	auditAction(h, context.Background(), m.GuildID, "BAN", m.Author.ID, targetID, map[string]any{"reason": reason})
	_, err = s.ChannelMessageSendEmbed(m.ChannelID, modSuccessEmbed("Ban", fmt.Sprintf("Banned **%s**.", name)))
	return err
}

// banSlash: /ban user:<user> [reason]
func banSlashCommandHandler(h *CommandHandler, s *discordgo.Session, i *discordgo.InteractionCreate) error {
	ok, msg := requireModerator(i.Member)
	if !ok {
		return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{Embeds: []*discordgo.MessageEmbed{modErrorEmbed("Ban", msg)}},
		})
	}
	opts := OptionMap(i.ApplicationCommandData().Options)
	raw := OptString(opts, "user")
	reason := OptString(opts, "reason")
	targetID, name, err := resolveTargetUser(s, i.GuildID, raw)
	if err != nil {
		return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{Embeds: []*discordgo.MessageEmbed{modErrorEmbed("Ban", fmt.Sprintf("Could not find that user: %s", err))}},
		})
	}
	if err := s.GuildBanCreateWithReason(i.GuildID, targetID, reason, 0); err != nil {
		return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{Embeds: []*discordgo.MessageEmbed{modErrorEmbed("Ban", fmt.Sprintf("Failed to ban %s: %s", name, err))}},
		})
	}
	auditAction(h, context.Background(), i.GuildID, "BAN", i.Member.User.ID, targetID, map[string]any{"reason": reason})
	return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{Embeds: []*discordgo.MessageEmbed{modSuccessEmbed("Ban", fmt.Sprintf("Banned **%s**.", name))}},
	})
}

// hardbanMessage: -hardban <user> [reason]  (ban + delete 7 days of messages)
func hardbanMessageCommandHandler(h *CommandHandler, s *discordgo.Session, m *discordgo.MessageCreate, args []string) error {
	if ok, err := requireModForMessage(s, m, "Hard Ban"); !ok {
		return err
	}
	if len(args) == 0 {
		return sendModError(s, m.ChannelID, "Hard Ban", "Usage: `-hardban <user> [reason]`")
	}
	targetID, name, err := resolveTargetUser(s, m.GuildID, args[0])
	if err != nil {
		return sendModError(s, m.ChannelID, "Hard Ban", fmt.Sprintf("Could not find that user: %s", err))
	}
	reason := strings.Join(args[1:], " ")
	if err := s.GuildBanCreateWithReason(m.GuildID, targetID, reason, 7); err != nil {
		return sendModError(s, m.ChannelID, "Hard Ban", fmt.Sprintf("Failed to hard-ban %s: %s", name, err))
	}
	auditAction(h, context.Background(), m.GuildID, "HARDBAN", m.Author.ID, targetID, map[string]any{"reason": reason})
	_, err = s.ChannelMessageSendEmbed(m.ChannelID, modSuccessEmbed("Hard Ban", fmt.Sprintf("Hard-banned **%s** and wiped their recent messages.", name)))
	return err
}

func hardbanSlashCommandHandler(h *CommandHandler, s *discordgo.Session, i *discordgo.InteractionCreate) error {
	ok, msg := requireModerator(i.Member)
	if !ok {
		return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{Embeds: []*discordgo.MessageEmbed{modErrorEmbed("Hard Ban", msg)}},
		})
	}
	opts := OptionMap(i.ApplicationCommandData().Options)
	raw := OptString(opts, "user")
	reason := OptString(opts, "reason")
	targetID, name, err := resolveTargetUser(s, i.GuildID, raw)
	if err != nil {
		return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{Embeds: []*discordgo.MessageEmbed{modErrorEmbed("Hard Ban", fmt.Sprintf("Could not find that user: %s", err))}},
		})
	}
	if err := s.GuildBanCreateWithReason(i.GuildID, targetID, reason, 7); err != nil {
		return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{Embeds: []*discordgo.MessageEmbed{modErrorEmbed("Hard Ban", fmt.Sprintf("Failed to hard-ban %s: %s", name, err))}},
		})
	}
	auditAction(h, context.Background(), i.GuildID, "HARDBAN", i.Member.User.ID, targetID, map[string]any{"reason": reason})
	return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{Embeds: []*discordgo.MessageEmbed{modSuccessEmbed("Hard Ban", fmt.Sprintf("Hard-banned **%s** and wiped their recent messages.", name))}},
	})
}

// softbanMessage: -softban <user> [reason]  (ban then immediately unban)
func softbanMessageCommandHandler(h *CommandHandler, s *discordgo.Session, m *discordgo.MessageCreate, args []string) error {
	if ok, err := requireModForMessage(s, m, "Soft Ban"); !ok {
		return err
	}
	if len(args) == 0 {
		return sendModError(s, m.ChannelID, "Soft Ban", "Usage: `-softban <user> [reason]`")
	}
	targetID, name, err := resolveTargetUser(s, m.GuildID, args[0])
	if err != nil {
		return sendModError(s, m.ChannelID, "Soft Ban", fmt.Sprintf("Could not find that user: %s", err))
	}
	reason := strings.Join(args[1:], " ")
	if err := s.GuildBanCreateWithReason(m.GuildID, targetID, reason, 1); err != nil {
		return sendModError(s, m.ChannelID, "Soft Ban", fmt.Sprintf("Failed to soft-ban %s: %s", name, err))
	}
	if err := s.GuildBanDelete(m.GuildID, targetID); err != nil {
		return sendModError(s, m.ChannelID, "Soft Ban", fmt.Sprintf("Unbanned **%s** but could not remove the ban record: %s", name, err))
	}
	auditAction(h, context.Background(), m.GuildID, "SOFTBAN", m.Author.ID, targetID, map[string]any{"reason": reason})
	_, err = s.ChannelMessageSendEmbed(m.ChannelID, modSuccessEmbed("Soft Ban", fmt.Sprintf("Soft-banned **%s** (banned then unbanned, recent messages wiped).", name)))
	return err
}

func softbanSlashCommandHandler(h *CommandHandler, s *discordgo.Session, i *discordgo.InteractionCreate) error {
	ok, msg := requireModerator(i.Member)
	if !ok {
		return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{Embeds: []*discordgo.MessageEmbed{modErrorEmbed("Soft Ban", msg)}},
		})
	}
	opts := OptionMap(i.ApplicationCommandData().Options)
	raw := OptString(opts, "user")
	reason := OptString(opts, "reason")
	targetID, name, err := resolveTargetUser(s, i.GuildID, raw)
	if err != nil {
		return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{Embeds: []*discordgo.MessageEmbed{modErrorEmbed("Soft Ban", fmt.Sprintf("Could not find that user: %s", err))}},
		})
	}
	if err := s.GuildBanCreateWithReason(i.GuildID, targetID, reason, 1); err != nil {
		return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{Embeds: []*discordgo.MessageEmbed{modErrorEmbed("Soft Ban", fmt.Sprintf("Failed to soft-ban %s: %s", name, err))}},
		})
	}
	if err := s.GuildBanDelete(i.GuildID, targetID); err != nil {
		return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{Embeds: []*discordgo.MessageEmbed{modErrorEmbed("Soft Ban", fmt.Sprintf("Unbanned **%s** but could not remove the ban record: %s", name, err))}},
		})
	}
	auditAction(h, context.Background(), i.GuildID, "SOFTBAN", i.Member.User.ID, targetID, map[string]any{"reason": reason})
	return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{Embeds: []*discordgo.MessageEmbed{modSuccessEmbed("Soft Ban", fmt.Sprintf("Soft-banned **%s** (banned then unbanned, recent messages wiped).", name))}},
	})
}

// kickMessage: -kick <user> [reason]
func kickMessageCommandHandler(h *CommandHandler, s *discordgo.Session, m *discordgo.MessageCreate, args []string) error {
	if ok, err := requireModForMessage(s, m, "Kick"); !ok {
		return err
	}
	if len(args) == 0 {
		return sendModError(s, m.ChannelID, "Kick", "Usage: `-kick <user> [reason]`")
	}
	targetID, name, err := resolveTargetUser(s, m.GuildID, args[0])
	if err != nil {
		return sendModError(s, m.ChannelID, "Kick", fmt.Sprintf("Could not find that user: %s", err))
	}
	reason := strings.Join(args[1:], " ")
	if err := s.GuildMemberDeleteWithReason(m.GuildID, targetID, reason); err != nil {
		return sendModError(s, m.ChannelID, "Kick", fmt.Sprintf("Failed to kick %s: %s", name, err))
	}
	auditAction(h, context.Background(), m.GuildID, "KICK", m.Author.ID, targetID, map[string]any{"reason": reason})
	_, err = s.ChannelMessageSendEmbed(m.ChannelID, modSuccessEmbed("Kick", fmt.Sprintf("Kicked **%s**.", name)))
	return err
}

func kickSlashCommandHandler(h *CommandHandler, s *discordgo.Session, i *discordgo.InteractionCreate) error {
	ok, msg := requireModerator(i.Member)
	if !ok {
		return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{Embeds: []*discordgo.MessageEmbed{modErrorEmbed("Kick", msg)}},
		})
	}
	opts := OptionMap(i.ApplicationCommandData().Options)
	raw := OptString(opts, "user")
	reason := OptString(opts, "reason")
	targetID, name, err := resolveTargetUser(s, i.GuildID, raw)
	if err != nil {
		return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{Embeds: []*discordgo.MessageEmbed{modErrorEmbed("Kick", fmt.Sprintf("Could not find that user: %s", err))}},
		})
	}
	if err := s.GuildMemberDeleteWithReason(i.GuildID, targetID, reason); err != nil {
		return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{Embeds: []*discordgo.MessageEmbed{modErrorEmbed("Kick", fmt.Sprintf("Failed to kick %s: %s", name, err))}},
		})
	}
	auditAction(h, context.Background(), i.GuildID, "KICK", i.Member.User.ID, targetID, map[string]any{"reason": reason})
	return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{Embeds: []*discordgo.MessageEmbed{modSuccessEmbed("Kick", fmt.Sprintf("Kicked **%s**.", name))}},
	})
}
