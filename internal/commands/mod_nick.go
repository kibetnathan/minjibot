package commands

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/kibetnathan/minjibot/internal/ports/dto"
)

const nickLockRole = "NICK_LOCK"

// fnMessage: -fn <user> <new nickname>  (force a nickname)
func fnMessageCommandHandler(h *CommandHandler, s *discordgo.Session, m *discordgo.MessageCreate, args []string) error {
	if ok, err := requireModForMessage(s, m, "Force Nickname"); !ok {
		return err
	}
	if len(args) < 2 {
		return sendModError(s, m.ChannelID, "Force Nickname", "Usage: `-fn <user> <new nickname>`")
	}
	targetID, _, err := resolveTargetUser(s, m.GuildID, args[0])
	if err != nil {
		return sendModError(s, m.ChannelID, "Force Nickname", fmt.Sprintf("Could not find that user: %s", err))
	}
	nick := joinArgs(args[1:])
	if err := s.GuildMemberNickname(m.GuildID, targetID, nick); err != nil {
		return sendModError(s, m.ChannelID, "Force Nickname", fmt.Sprintf("Failed to set nickname: %s", err))
	}
	auditAction(h, context.Background(), m.GuildID, "FORCE_NICK", m.Author.ID, targetID, map[string]any{"nickname": nick})
	_, err = s.ChannelMessageSendEmbed(m.ChannelID, modSuccessEmbed("Force Nickname", fmt.Sprintf("Set nickname to **%s**.", nick)))
	return err
}

func fnSlashCommandHandler(h *CommandHandler, s *discordgo.Session, i *discordgo.InteractionCreate) error {
	ok, msg := requireModerator(s, i.GuildID, i.Member.User.ID, i.Member)
	if !ok {
		return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{Embeds: []*discordgo.MessageEmbed{modErrorEmbed("Force Nickname", msg)}},
		})
	}
	opts := OptionMap(i.ApplicationCommandData().Options)
	raw := OptString(opts, "user")
	nick := OptString(opts, "nickname")
	targetID, _, err := resolveTargetUser(s, i.GuildID, raw)
	if err != nil {
		return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{Embeds: []*discordgo.MessageEmbed{modErrorEmbed("Force Nickname", fmt.Sprintf("Could not find that user: %s", err))}},
		})
	}
	if err := s.GuildMemberNickname(i.GuildID, targetID, nick); err != nil {
		return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{Embeds: []*discordgo.MessageEmbed{modErrorEmbed("Force Nickname", fmt.Sprintf("Failed to set nickname: %s", err))}},
		})
	}
	auditAction(h, context.Background(), i.GuildID, "FORCE_NICK", i.Member.User.ID, targetID, map[string]any{"nickname": nick})
	return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{Embeds: []*discordgo.MessageEmbed{modSuccessEmbed("Force Nickname", fmt.Sprintf("Set nickname to **%s**.", nick))}},
	})
}

// nickMessage: -nick lock <user> | -nick unlock <user>
func nickMessageCommandHandler(h *CommandHandler, s *discordgo.Session, m *discordgo.MessageCreate, args []string) error {
	if ok, err := requireModForMessage(s, m, "Nick lock"); !ok {
		return err
	}
	if len(args) < 2 {
		return sendModError(s, m.ChannelID, "Nick lock", "Usage: `-nick lock <user>` or `-nick unlock <user>`")
	}
	action := args[0]
	targetID, name, err := resolveTargetUser(s, m.GuildID, args[1])
	if err != nil {
		return sendModError(s, m.ChannelID, "Nick lock", fmt.Sprintf("Could not find that user: %s", err))
	}
	switch action {
	case "lock":
		if err := setNickLock(h, m.GuildID, targetID, true); err != nil {
			return sendModError(s, m.ChannelID, "Nick lock", fmt.Sprintf("Failed to store lock state: %s", err))
		}
		auditAction(h, context.Background(), m.GuildID, "NICK_LOCK", m.Author.ID, targetID, nil)
		_, err = s.ChannelMessageSendEmbed(m.ChannelID, modSuccessEmbed("Nick lock", fmt.Sprintf("Locked **%s**'s nickname.", name)))
		return err
	case "unlock":
		if err := setNickLock(h, m.GuildID, targetID, false); err != nil {
			return sendModError(s, m.ChannelID, "Nick lock", fmt.Sprintf("Failed to store lock state: %s", err))
		}
		auditAction(h, context.Background(), m.GuildID, "NICK_UNLOCK", m.Author.ID, targetID, nil)
		_, err = s.ChannelMessageSendEmbed(m.ChannelID, modSuccessEmbed("Nick lock", fmt.Sprintf("Unlocked **%s**'s nickname.", name)))
		return err
	default:
		return sendModError(s, m.ChannelID, "Nick lock", "Usage: `-nick lock <user>` or `-nick unlock <user>`")
	}
}

func nickSlashCommandHandler(h *CommandHandler, s *discordgo.Session, i *discordgo.InteractionCreate) error {
	ok, msg := requireModerator(s, i.GuildID, i.Member.User.ID, i.Member)
	if !ok {
		return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{Embeds: []*discordgo.MessageEmbed{modErrorEmbed("Nick lock", msg)}},
		})
	}
	subs := i.ApplicationCommandData().Options
	if len(subs) == 0 {
		return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{Embeds: []*discordgo.MessageEmbed{modErrorEmbed("Nick lock", "No subcommand provided.")}},
		})
	}
	sub := subs[0]
	opts := OptionMap(sub.Options)
	raw := OptString(opts, "user")
	targetID, name, err := resolveTargetUser(s, i.GuildID, raw)
	if err != nil {
		return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{Embeds: []*discordgo.MessageEmbed{modErrorEmbed("Nick lock", fmt.Sprintf("Could not find that user: %s", err))}},
		})
	}
	lock := sub.Name == "lock"
	if err := setNickLock(h, i.GuildID, targetID, lock); err != nil {
		return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{Embeds: []*discordgo.MessageEmbed{modErrorEmbed("Nick lock", fmt.Sprintf("Failed to store lock state: %s", err))}},
		})
	}
	action := "NICK_UNLOCK"
	desc := fmt.Sprintf("Unlocked **%s**'s nickname.", name)
	if lock {
		action = "NICK_LOCK"
		desc = fmt.Sprintf("Locked **%s**'s nickname.", name)
	}
	auditAction(h, context.Background(), i.GuildID, action, i.Member.User.ID, targetID, nil)
	return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{Embeds: []*discordgo.MessageEmbed{modSuccessEmbed("Nick lock", desc)}},
	})
}

// setNickLock persists a nick-lock flag for a user in user_permissions.
func setNickLock(h *CommandHandler, guildID, userID string, locked bool) error {
	if h == nil || h.PermRepo == nil {
		return fmt.Errorf("permission store unavailable")
	}
	payload, _ := json.Marshal(map[string]bool{"nick_locked": locked})
	_, err := h.PermRepo.Upsert(context.Background(), dto.UpsertUserPermissionParams{
		UserID:          userID,
		GuildID:         guildID,
		Role:            nickLockRole,
		PermissionsJSON: payload,
	})
	return err
}
