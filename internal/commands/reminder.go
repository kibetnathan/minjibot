package commands

import (
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/bwmarrin/discordgo"
)

var (
	reminderMu    sync.Mutex
	reminderCount int
)

// scheduleReminder sends a direct message to the user after the given delay.
func scheduleReminder(s *discordgo.Session, userID, channelID, text string, delay time.Duration) {
	reminderMu.Lock()
	reminderCount++
	reminderMu.Unlock()

	go func() {
		timer := time.NewTimer(delay)
		defer timer.Stop()
		<-timer.C

		dm, err := s.UserChannelCreate(userID)
		if err != nil {
			return
		}

		embed := &discordgo.MessageEmbed{
			Color: 0xFEE75C,
			Title: "Reminder",
			Fields: []*discordgo.MessageEmbedField{
				{Name: "Task", Value: TruncateForEmbed(text, 1024)},
				{Name: "Asked in", Value: "<#" + channelID + ">"},
			},
			Timestamp: time.Now().Format(time.RFC3339),
		}
		s.ChannelMessageSendEmbed(dm.ID, embed)
	}()
}

// ParseDuration extracts durations like "30s", "5m", "2h", "1d", or a
// combination like "1h30m". Returns 0 if nothing parseable is found.
func ParseDuration(s string) time.Duration {
	var total time.Duration

	for len(s) > 0 {
		i := 0
		for i < len(s) && s[i] >= '0' && s[i] <= '9' {
			i++
		}
		if i == 0 {
			return 0
		}
		n, err := strconv.Atoi(s[:i])
		if err != nil || n < 0 {
			return 0
		}
		s = s[i:]

		if len(s) == 0 {
			return 0
		}
		unit, rest := s[0], s[1:]
		var d time.Duration
		switch unit {
		case 's':
			d = time.Duration(n) * time.Second
		case 'm':
			d = time.Duration(n) * time.Minute
		case 'h':
			d = time.Duration(n) * time.Hour
		case 'd':
			d = time.Duration(n) * 24 * time.Hour
		default:
			return 0
		}
		total += d
		s = rest
	}

	return total
}

// ParseReminderArgs returns the reminder text and its delay.
// Supported forms: "in 30m", "in 2h30m", "in 30m text...",
// "30m text...", or "in 14:00 text..." for an absolute time today.
func ParseReminderArgs(args []string) (text string, delay time.Duration) {
	if len(args) == 0 {
		return "", 0
	}

	pieces := make([]string, 0, len(args))
	for _, arg := range args {
		pieces = append(pieces, strings.TrimPrefix(arg, ","))
	}

	joined := strings.Join(pieces, " ")
	rest := joined

	if strings.HasPrefix(strings.ToLower(joined), "in ") {
		rest = strings.TrimSpace(joined[3:])
	}

	fields := strings.Fields(rest)

	// Absolute HH:MM time.
	if len(fields) >= 2 && IsClock(fields[0]) {
		at, err := time.ParseInLocation("15:04", fields[0], time.Local)
		if err == nil {
			now := time.Now()
			when := time.Date(now.Year(), now.Month(), now.Day(), at.Hour(), at.Minute(), 0, 0, time.Local)
			if !when.After(now) {
				when = when.Add(24 * time.Hour)
			}
			return strings.TrimSpace(rest[len(fields[0]):]), when.Sub(now)
		}
	}

	// Relative duration: a single field made of digits+units.
	if d := ParseDuration(fields[0]); d > 0 {
		return strings.TrimSpace(rest[len(fields[0]):]), d
	}
	// Duration plus "from now" semantics: "1h from now".
	if len(fields) >= 3 && fields[1] == "from" {
		if d := ParseDuration(fields[0]); d > 0 {
			return strings.TrimSpace(rest[len(fields[0])+len(" from now"):]), d
		}
	}

	return "", 0
}

func IsClock(s string) bool {
	if len(s) != 5 || s[2] != ':' {
		return false
	}
	return s[0] >= '0' && s[0] <= '2' && s[1] >= '0' && s[1] <= '9' && s[3] >= '0' && s[3] <= '9' && s[4] >= '0' && s[4] <= '9'
}

func reminderMessageCommandHandler(s *discordgo.Session, m *discordgo.MessageCreate, args []string) error {
	text, delay := ParseReminderArgs(args)
	if text == "" || delay <= 0 {
		_, err := s.ChannelMessageSend(m.ChannelID, "Usage: `!reminder in 30m buy milk`, `!reminder 2h take a break`, or `!reminder in 14:00 lunch`")
		return err
	}

	scheduleReminder(s, m.Author.ID, m.ChannelID, text, delay)

	_, err := s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Reminder set — I'll ping you in %s.", HumanizeDuration(delay)))
	return err
}

func reminderSlashCommandHandler(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	opts := OptionMap(i.ApplicationCommandData().Options)
	text := OptString(opts, "text")
	if text == "" {
		text = strings.TrimSpace(OptString(opts, "when"))
	}
	delayStr := OptString(opts, "delay")
	if delayStr == "" {
		delayStr = OptString(opts, "time")
	}

	delay := ParseDuration(delayStr)
	if text == "" || delay <= 0 {
		return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{Content: "Usage: `/reminder text:<what> delay:<2h30m>`"},
		})
	}

	userID := i.Member.User.ID
	if userID == "" {
		userID = i.User.ID
	}

	scheduleReminder(s, userID, i.ChannelID, text, delay)

	return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{Content: fmt.Sprintf("Reminder set — I'll ping you in %s.", HumanizeDuration(delay))},
	})
}

func HumanizeDuration(d time.Duration) string {
	switch {
	case d < time.Minute:
		return fmt.Sprintf("%ds", int(d.Seconds()))
	case d < time.Hour:
		return fmt.Sprintf("%dm", int(d.Minutes()))
	case d < 24*time.Hour:
		h := int(d.Hours())
		m := int(d.Minutes()) % 60
		if m == 0 {
			return fmt.Sprintf("%dh", h)
		}
		return fmt.Sprintf("%dh%dm", h, m)
	default:
		dd := int(d.Hours() / 24)
		h := int(d.Hours()) % 24
		if h == 0 {
			return fmt.Sprintf("%dd", dd)
		}
		return fmt.Sprintf("%dd%dh", dd, h)
	}
}
