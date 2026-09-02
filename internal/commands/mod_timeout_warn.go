package commands

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

// timeoutMessage: -timeout <user> <duration> [reason]   e.g. -timeout @x 30m
func timeoutMessageCommandHandler(h *CommandHandler, s *discordgo.Session, m *discordgo.MessageCreate, args []string) error {
	if ok, err := requireModForMessage(s, m, "Timeout"); !ok {
		return err
	}
	if len(args) < 2 {
		return sendModError(s, m.ChannelID, "Timeout", "Usage: `-timeout <user> <duration>` (e.g. `30m`, `2h`, `1d`); optionally append a reason.")
	}
	targetID, name, err := resolveTargetUser(s, m.GuildID, args[0])
	if err != nil {
		return sendModError(s, m.ChannelID, "Timeout", fmt.Sprintf("Could not find that user: %s", err))
	}
	dur, err := parseDuration(args[1])
	if err != nil || dur <= 0 {
		return sendModError(s, m.ChannelID, "Timeout", "Invalid duration. Use e.g. `30m`, `2h`, `1d`.")
	}
	if dur > 28*24*time.Hour {
		dur = 28 * 24 * time.Hour
	}
	reason := strings.Join(args[2:], " ")
	until := time.Now().Add(dur)
	if err := s.GuildMemberTimeout(m.GuildID, targetID, &until); err != nil {
		return sendModError(s, m.ChannelID, "Timeout", fmt.Sprintf("Failed to timeout %s: %s", name, err))
	}
	auditAction(h, context.Background(), m.GuildID, "TIMEOUT", m.Author.ID, targetID, map[string]any{"duration": dur.String(), "reason": reason})
	_, err = s.ChannelMessageSendEmbed(m.ChannelID, modSuccessEmbed("Timeout", fmt.Sprintf("Timed out **%s** for %s.", name, dur.String())))
	return err
}

func timeoutSlashCommandHandler(h *CommandHandler, s *discordgo.Session, i *discordgo.InteractionCreate) error {
	ok, msg := requireModerator(s, i.GuildID, i.Member.User.ID, i.Member)
	if !ok {
		return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{Embeds: []*discordgo.MessageEmbed{modErrorEmbed("Timeout", msg)}},
		})
	}
	opts := OptionMap(i.ApplicationCommandData().Options)
	raw := OptString(opts, "user")
	reason := OptString(opts, "reason")
	durStr := OptString(opts, "duration")
	targetID, name, err := resolveTargetUser(s, i.GuildID, raw)
	if err != nil {
		return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{Embeds: []*discordgo.MessageEmbed{modErrorEmbed("Timeout", fmt.Sprintf("Could not find that user: %s", err))}},
		})
	}
	dur, err := parseDuration(durStr)
	if err != nil || dur <= 0 {
		return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{Embeds: []*discordgo.MessageEmbed{modErrorEmbed("Timeout", "Invalid duration. Use e.g. `30m`, `2h`, `1d`.")}},
		})
	}
	if dur > 28*24*time.Hour {
		dur = 28 * 24 * time.Hour
	}
	until := time.Now().Add(dur)
	if err := s.GuildMemberTimeout(i.GuildID, targetID, &until); err != nil {
		return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{Embeds: []*discordgo.MessageEmbed{modErrorEmbed("Timeout", fmt.Sprintf("Failed to timeout %s: %s", name, err))}},
		})
	}
	auditAction(h, context.Background(), i.GuildID, "TIMEOUT", i.Member.User.ID, targetID, map[string]any{"duration": dur.String(), "reason": reason})
	return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{Embeds: []*discordgo.MessageEmbed{modSuccessEmbed("Timeout", fmt.Sprintf("Timed out **%s** for %s.", name, dur.String()))}},
	})
}

func parseDuration(s string) (time.Duration, error) {
	return time.ParseDuration(s)
}

// warnMessage: -warn <user> [reason]
func warnMessageCommandHandler(h *CommandHandler, s *discordgo.Session, m *discordgo.MessageCreate, args []string) error {
	if ok, err := requireModForMessage(s, m, "Warn"); !ok {
		return err
	}
	if len(args) == 0 {
		return sendModError(s, m.ChannelID, "Warn", "Usage: `-warn <user> [reason]`")
	}
	targetID, name, err := resolveTargetUser(s, m.GuildID, args[0])
	if err != nil {
		return sendModError(s, m.ChannelID, "Warn", fmt.Sprintf("Could not find that user: %s", err))
	}
	reason := strings.Join(args[1:], " ")
	auditAction(h, context.Background(), m.GuildID, "WARN", m.Author.ID, targetID, map[string]any{"reason": reason})
	_, err = s.ChannelMessageSendEmbed(m.ChannelID, modSuccessEmbed("Warn", fmt.Sprintf("Warned **%s**.", name)))
	return err
}

func warnSlashCommandHandler(h *CommandHandler, s *discordgo.Session, i *discordgo.InteractionCreate) error {
	ok, msg := requireModerator(s, i.GuildID, i.Member.User.ID, i.Member)
	if !ok {
		return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{Embeds: []*discordgo.MessageEmbed{modErrorEmbed("Warn", msg)}},
		})
	}
	opts := OptionMap(i.ApplicationCommandData().Options)
	raw := OptString(opts, "user")
	reason := OptString(opts, "reason")
	targetID, name, err := resolveTargetUser(s, i.GuildID, raw)
	if err != nil {
		return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{Embeds: []*discordgo.MessageEmbed{modErrorEmbed("Warn", fmt.Sprintf("Could not find that user: %s", err))}},
		})
	}
	auditAction(h, context.Background(), i.GuildID, "WARN", i.Member.User.ID, targetID, map[string]any{"reason": reason})
	return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{Embeds: []*discordgo.MessageEmbed{modSuccessEmbed("Warn", fmt.Sprintf("Warned **%s**.", name))}},
	})
}

// historyMessage: -history <user>  (show the user's punishment/infraction log)
func historyMessageCommandHandler(h *CommandHandler, s *discordgo.Session, m *discordgo.MessageCreate, args []string) error {
	if ok, err := requireModForMessage(s, m, "History"); !ok {
		return err
	}
	if len(args) == 0 {
		return sendModError(s, m.ChannelID, "History", "Usage: `-history <user>`")
	}
	targetID, _, err := resolveTargetUser(s, m.GuildID, args[0])
	if err != nil {
		return sendModError(s, m.ChannelID, "History", fmt.Sprintf("Could not find that user: %s", err))
	}
	embed, err := buildHistoryEmbed(h, s, m.GuildID, targetID)
	if err != nil {
		return sendModError(s, m.ChannelID, "History", err.Error())
	}
	_, err = s.ChannelMessageSendEmbed(m.ChannelID, embed)
	return err
}

func historySlashCommandHandler(h *CommandHandler, s *discordgo.Session, i *discordgo.InteractionCreate) error {
	ok, msg := requireModerator(s, i.GuildID, i.Member.User.ID, i.Member)
	if !ok {
		return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{Embeds: []*discordgo.MessageEmbed{modErrorEmbed("History", msg)}},
		})
	}
	opts := OptionMap(i.ApplicationCommandData().Options)
	raw := OptString(opts, "user")
	targetID, _, err := resolveTargetUser(s, i.GuildID, raw)
	if err != nil {
		return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{Embeds: []*discordgo.MessageEmbed{modErrorEmbed("History", fmt.Sprintf("Could not find that user: %s", err))}},
		})
	}
	embed, err := buildHistoryEmbed(h, s, i.GuildID, targetID)
	if err != nil {
		return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{Embeds: []*discordgo.MessageEmbed{modErrorEmbed("History", err.Error())}},
		})
	}
	return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{Embeds: []*discordgo.MessageEmbed{embed}},
	})
}

// buildHistoryEmbed assembles a user's punishment history from audit log
// entries where the target is the user and the action is a punishment.
func buildHistoryEmbed(h *CommandHandler, s *discordgo.Session, guildID, targetID string) (*discordgo.MessageEmbed, error) {
	user, err := s.User(targetID)
	if err != nil {
		return nil, fmt.Errorf("could not fetch user: %s", err)
	}

	var lines []string
	// Collect punishment entries targeting this user by walking the guild log.
	if h != nil && h.AuditRepo != nil {
		logs, err := h.AuditRepo.ListForGuild(context.Background(), guildID, 100, 0)
		if err == nil {
			for _, l := range logs {
				if l.TargetID != targetID {
					continue
				}
				if !isPunishmentAction(l.Action) {
					continue
				}
				lines = append(lines, fmt.Sprintf("`%s` — %s", l.CreatedAt.Format("2006-01-02 15:04"), l.Action))
			}
		}
	}

	if len(lines) == 0 {
		lines = append(lines, "No recorded infractions.")
	}

	return &discordgo.MessageEmbed{
		Color:       modColor,
		Title:       fmt.Sprintf("History — %s", user.Username),
		Description: strings.Join(lines, "\n"),
	}, nil
}

func isPunishmentAction(action string) bool {
	switch action {
	case "BAN", "HARDBAN", "SOFTBAN", "KICK", "TIMEOUT", "WARN", "JAIL", "STAFFSTRIP":
		return true
	}
	return false
}
