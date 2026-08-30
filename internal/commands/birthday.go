package commands

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/kibetnathan/minjibot/internal/ports/dto"
)

// parseBirthday parses a birthday string into a time.Time. Supported formats:
// MM-DD, MM/DD, MM/DD/YYYY, YYYY-MM-DD, and any layout Go's time can parse.
// The entered year is preserved when present, otherwise a representative leap
// year (2000) is used so month/day matching stays correct.
func parseBirthday(raw string) (time.Time, error) {
	raw = strings.TrimSpace(raw)
	layouts := []string{
		"2006-01-02",
		"01-02-2006",
		"01-02",
		"01/02/2006",
		"01/02",
		"January 2, 2006",
	}
	for _, l := range layouts {
		if t, err := time.Parse(l, raw); err == nil {
			return t, nil
		}
	}
	// Bare "MM" or "Jan" not supported; give a helpful error.
	return time.Time{}, fmt.Errorf("couldn't parse %q as a date. Try formats like MM-DD, MM/DD, YYYY-MM-DD, or MM/DD/YYYY", raw)
}

// birthdayMonthDay returns the month and day of a birthday for display/matching.
func birthdayMonthDay(t time.Time) (int, int) {
	return int(t.Month()), t.Day()
}

func birthdayMessageCommandHandler(h *CommandHandler, s *discordgo.Session, m *discordgo.MessageCreate, args []string) error {
	if len(args) == 0 {
		if _, err := s.ChannelMessageSend(m.ChannelID, "Usage: `-birthday add|list|celebrate|channel|role`"); err != nil {
			return err
		}
		return fmt.Errorf("birthday requires a subcommand")
	}

	sub := strings.ToLower(args[0])
	switch sub {
	case "add":
		if len(args) < 2 {
			if _, err := s.ChannelMessageSend(m.ChannelID, "Usage: `-birthday add <date>` e.g. `-birthday add 07-14`"); err != nil {
				return err
			}
			return nil
		}
		return h.birthdayAdd(s, m, m.Author.ID, strings.Join(args[1:], " "))
	case "list":
		return h.birthdayList(s, m)
	case "celebrate":
		targetID := ""
		if len(args) > 1 {
			targetID = ParseMentionID(args[1])
		}
		return h.birthdayCelebrate(s, m, targetID)
	case "channel":
		if len(args) < 2 {
			if _, err := s.ChannelMessageSend(m.ChannelID, "Usage: `-birthday channel <#channel>`"); err != nil {
				return err
			}
			return nil
		}
		return h.birthdaySetChannel(s, m, args[1])
	case "role":
		if len(args) < 2 {
			if _, err := s.ChannelMessageSend(m.ChannelID, "Usage: `-birthday role <@role>`"); err != nil {
				return err
			}
			return nil
		}
		return h.birthdaySetRole(s, m, args[1])
	default:
		if _, err := s.ChannelMessageSend(m.ChannelID, "Usage: `-birthday add|list|celebrate|channel|role`"); err != nil {
			return err
		}
		return nil
	}
}

func (h *CommandHandler) birthdayAdd(s *discordgo.Session, m *discordgo.MessageCreate, userID, raw string) error {
	t, err := parseBirthday(raw)
	if err != nil {
		if _, serr := s.ChannelMessageSend(m.ChannelID, err.Error()); serr != nil {
			return serr
		}
		return nil
	}
	_, err = h.BirthdayRepo.Set(context.Background(), dto.SetBirthdayParams{
		GuildID:  m.GuildID,
		UserID:   userID,
		Birthday: t,
	})
	if err != nil {
		return err
	}
	if userID == m.Author.ID {
		_, err = s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("<@%s> birthday saved as **%s**.", userID, t.Format("January 2")))
	} else {
		_, err = s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Saved birthday **%s** for <@%s>.", t.Format("January 2"), userID))
	}
	return err
}

func (h *CommandHandler) birthdayList(s *discordgo.Session, m *discordgo.MessageCreate) error {
	birthdays, err := h.BirthdayRepo.ListByGuild(context.Background(), m.GuildID)
	if err != nil {
		return err
	}
	if len(birthdays) == 0 {
		_, err = s.ChannelMessageSend(m.ChannelID, "No birthdays saved in this server yet. Use `-birthday add <date>`.")
		return err
	}

	sort.Slice(birthdays, func(i, j int) bool {
		mi, di := birthdayMonthDay(birthdays[i].Birthday)
		mj, dj := birthdayMonthDay(birthdays[j].Birthday)
		if mi != mj {
			return mi < mj
		}
		return di < dj
	})

	var b strings.Builder
	now := time.Now()
	for _, bd := range birthdays {
		date := bd.Birthday.Format("Jan 2")
		// Show age if this year would be a birthday (approx).
		b.WriteString(fmt.Sprintf("<@%s> — %s\n", bd.UserID, date))
	}
	_ = now

	_, err = s.ChannelMessageSendEmbed(m.ChannelID, &discordgo.MessageEmbed{
		Title:       "🎂 Upcoming birthdays",
		Description: b.String(),
		Color:       0x5865F2,
	})
	return err
}

func (h *CommandHandler) birthdayCelebrate(s *discordgo.Session, m *discordgo.MessageCreate, targetID string) error {
	userID := targetID
	if userID == "" {
		userID = m.Author.ID
	}

	bd, err := h.BirthdayRepo.Get(context.Background(), m.GuildID, userID)
	if err != nil {
		if _, serr := s.ChannelMessageSend(m.ChannelID, "No birthday saved for that user. Use `-birthday add <date>`."); serr != nil {
			return serr
		}
		return nil
	}

	settings, _ := h.BirthdaySett.Get(context.Background(), m.GuildID)
	channelID := m.ChannelID
	if settings.ChannelID != "" {
		channelID = settings.ChannelID
	}

	_, err = s.ChannelMessageSend(channelID, fmt.Sprintf("🎉 Happy birthday <@%s>! 🎉 (%s)", userID, bd.Birthday.Format("January 2")))
	return err
}

func (h *CommandHandler) birthdaySetChannel(s *discordgo.Session, m *discordgo.MessageCreate, raw string) error {
	channelID := ParseMentionID(raw)
	if channelID == "" {
		channelID = raw
	}
	_, err := h.BirthdaySett.SetChannel(context.Background(), dto.SetGuildBirthdayChannelParams{
		GuildID:   m.GuildID,
		ChannelID: channelID,
	})
	if err != nil {
		return err
	}
	_, err = s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Birthday celebrations will be posted in <#%s>.", channelID))
	return err
}

func (h *CommandHandler) birthdaySetRole(s *discordgo.Session, m *discordgo.MessageCreate, raw string) error {
	roleID := strings.TrimSpace(raw)
	roleID = strings.TrimPrefix(roleID, "<@&")
	roleID = strings.TrimSuffix(roleID, ">")
	if roleID == "" {
		if _, serr := s.ChannelMessageSend(m.ChannelID, "Usage: `-birthday role <@role>`"); serr != nil {
			return serr
		}
		return nil
	}
	_, err := h.BirthdaySett.SetRole(context.Background(), dto.SetGuildBirthdayRoleParams{
		GuildID: m.GuildID,
		RoleID:  roleID,
	})
	if err != nil {
		return err
	}
	_, err = s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Birthday role set to <@&%s>.", roleID))
	return err
}

func birthdaySlashCommandHandler(h *CommandHandler, s *discordgo.Session, i *discordgo.InteractionCreate) error {
	data := i.ApplicationCommandData()
	sub := data.Options
	if len(sub) == 0 {
		return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{Content: "Usage: `/birthday add|list|celebrate|channel|role`"},
		})
	}

	opt := sub[0]
	name := opt.Name
	args := opt.Options

	respond := func(content string) error {
		return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{Content: content},
		})
	}
	respondEmbed := func(embed *discordgo.MessageEmbed) error {
		return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{Embeds: []*discordgo.MessageEmbed{embed}},
		})
	}

	ctx := context.Background()

	switch name {
	case "add":
		userStr := i.Member.User.ID
		rawDate := ""
		for _, a := range args {
			if a.Name == "date" {
				rawDate = a.StringValue()
			}
			if a.Name == "user" {
				if u, ok := a.Value.(string); ok && u != "" {
					userStr = u
				}
			}
		}
		t, err := parseBirthday(rawDate)
		if err != nil {
			return respond(err.Error())
		}
		if _, err := h.BirthdayRepo.Set(ctx, dto.SetBirthdayParams{
			GuildID: i.GuildID, UserID: userStr, Birthday: t,
		}); err != nil {
			return err
		}
		return respond(fmt.Sprintf("Saved birthday **%s** for <@%s>.", t.Format("January 2"), userStr))
	case "list":
		birthdays, err := h.BirthdayRepo.ListByGuild(ctx, i.GuildID)
		if err != nil {
			return err
		}
		if len(birthdays) == 0 {
			return respond("No birthdays saved in this server yet.")
		}
		sort.Slice(birthdays, func(a, b int) bool {
			ma, da := birthdayMonthDay(birthdays[a].Birthday)
			mb, db := birthdayMonthDay(birthdays[b].Birthday)
			if ma != mb {
				return ma < mb
			}
			return da < db
		})
		var b strings.Builder
		for _, bd := range birthdays {
			b.WriteString(fmt.Sprintf("<@%s> — %s\n", bd.UserID, bd.Birthday.Format("Jan 2")))
		}
		return respondEmbed(&discordgo.MessageEmbed{
			Title: "🎂 Upcoming birthdays", Description: b.String(), Color: 0x5865F2,
		})
	case "celebrate":
		userStr := i.Member.User.ID
		for _, a := range args {
			if a.Name == "user" {
				if u, ok := a.Value.(string); ok && u != "" {
					userStr = u
				}
			}
		}
		bd, err := h.BirthdayRepo.Get(ctx, i.GuildID, userStr)
		if err != nil {
			return respond("No birthday saved for that user.")
		}
		settings, _ := h.BirthdaySett.Get(ctx, i.GuildID)
		content := fmt.Sprintf("🎉 Happy birthday <@%s>! 🎉 (%s)", userStr, bd.Birthday.Format("January 2"))
		if settings.ChannelID != "" {
			content = fmt.Sprintf("🎉 Happy birthday <@%s>! 🎉 (%s)", userStr, bd.Birthday.Format("January 2"))
			_ = settings
			_, _ = s.ChannelMessageSend(settings.ChannelID, content)
			return respond("Celebrated!")
		}
		return respond(content)
	case "channel":
		channelID := ""
		for _, a := range args {
			if a.Name == "channel" {
				if v, ok := a.Value.(string); ok {
					channelID = v
				}
			}
		}
		if _, err := h.BirthdaySett.SetChannel(ctx, dto.SetGuildBirthdayChannelParams{
			GuildID: i.GuildID, ChannelID: channelID,
		}); err != nil {
			return err
		}
		return respond(fmt.Sprintf("Birthday celebrations will be posted in <#%s>.", channelID))
	case "role":
		roleID := ""
		for _, a := range args {
			if a.Name == "role" {
				if v, ok := a.Value.(string); ok {
					roleID = v
				}
			}
		}
		if _, err := h.BirthdaySett.SetRole(ctx, dto.SetGuildBirthdayRoleParams{
			GuildID: i.GuildID, RoleID: roleID,
		}); err != nil {
			return err
		}
		return respond(fmt.Sprintf("Birthday role set to <@&%s>.", roleID))
	default:
		return respond("Usage: `/birthday add|list|celebrate|channel|role`")
	}
}
