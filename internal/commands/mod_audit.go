package commands

import (
	"context"
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
)

// auditMessage: -audit [limit] [actor:<user>]
func auditMessageCommandHandler(h *CommandHandler, s *discordgo.Session, m *discordgo.MessageCreate, args []string) error {
	if ok, err := requireModForMessage(s, m, "Audit"); !ok {
		return err
	}
	limit := 10
	var actor string
	for _, a := range args {
		if strings.HasPrefix(a, "actor:") {
			actor = parseMentionID(strings.TrimPrefix(a, "actor:"))
		}
	}
	embed, err := buildAuditEmbed(h, s, m.GuildID, actor, limit)
	if err != nil {
		return sendModError(s, m.ChannelID, "Audit", err.Error())
	}
	_, err = s.ChannelMessageSendEmbed(m.ChannelID, embed)
	return err
}

func auditSlashCommandHandler(h *CommandHandler, s *discordgo.Session, i *discordgo.InteractionCreate) error {
	ok, msg := requireModerator(i.Member)
	if !ok {
		return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{Embeds: []*discordgo.MessageEmbed{modErrorEmbed("Audit", msg)}},
		})
	}
	opts := OptionMap(i.ApplicationCommandData().Options)
	limit := OptInt(opts, "limit", 10)
	if limit < 1 {
		limit = 1
	}
	if limit > 50 {
		limit = 50
	}
	actor := OptUser(opts, "actor")

	embed, err := buildAuditEmbed(h, s, i.GuildID, actor, limit)
	if err != nil {
		return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{Embeds: []*discordgo.MessageEmbed{modErrorEmbed("Audit", err.Error())}},
		})
	}
	return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{Embeds: []*discordgo.MessageEmbed{embed}},
	})
}

func buildAuditEmbed(h *CommandHandler, s *discordgo.Session, guildID, actor string, limit int) (*discordgo.MessageEmbed, error) {
	if h == nil || h.AuditRepo == nil {
		return nil, fmt.Errorf("audit logging is not available")
	}
	logs, err := h.AuditRepo.ListForGuild(context.Background(), guildID, int32(limit), 0)
	if err != nil {
		return nil, fmt.Errorf("could not fetch audit log: %s", err)
	}

	var lines []string
	for _, l := range logs {
		if actor != "" && l.ActorID != actor {
			continue
		}
		actName := l.ActorID
		if u, err := s.User(l.ActorID); err == nil {
			actName = u.Username
		}
		lines = append(lines, fmt.Sprintf("%s — `%s` by **%s**", l.CreatedAt.Format("2006-01-02 15:04"), l.Action, actName))
	}
	if len(lines) == 0 {
		lines = append(lines, "No audit log entries.")
	}

	title := "Audit Log"
	if actor != "" {
		title += " (filtered by actor)"
	}
	return &discordgo.MessageEmbed{
		Color:       modColor,
		Title:       title,
		Description: strings.Join(lines, "\n"),
	}, nil
}
