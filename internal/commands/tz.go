package commands

import (
	"fmt"
	"strings"
	"time"
	_ "time/tzdata"

	"github.com/bwmarrin/discordgo"
)

// tzAliases maps common time zone shorthand to IANA names. Ambiguous
// abbreviations (e.g. CST) resolve to their US interpretation.
var tzAliases = map[string]string{
	"utc":  "UTC",
	"gmt":  "UTC",
	"est":  "America/New_York",
	"edt":  "America/New_York",
	"cst":  "America/Chicago",
	"cdt":  "America/Chicago",
	"mst":  "America/Denver",
	"mdt":  "America/Denver",
	"pst":  "America/Los_Angeles",
	"pdt":  "America/Los_Angeles",
	"et":   "America/New_York",
	"ct":   "America/Chicago",
	"mt":   "America/Denver",
	"pt":   "America/Los_Angeles",
	"bst":  "Europe/London",
	"cet":  "Europe/Berlin",
	"cest": "Europe/Berlin",
	"hkt":  "Asia/Hong_Kong",
	"sgt":  "Asia/Singapore",
	"ist":  "Asia/Kolkata",
	"jst":  "Asia/Tokyo",
	"kst":  "Asia/Seoul",
	"aest": "Australia/Sydney",
	"acst": "Australia/Adelaide",
	"awst": "Australia/Perth",
	"nzst": "Pacific/Auckland",
}

// commonTz is used for fuzzy suggestions on unknown zone names.
var commonTz = []string{
	"UTC",
	"America/New_York", "America/Chicago", "America/Denver", "America/Los_Angeles",
	"America/Anchorage", "America/Phoenix", "America/Toronto", "America/Mexico_City",
	"America/Bogota", "America/Sao_Paulo", "America/Buenos_Aires",
	"Europe/London", "Europe/Paris", "Europe/Berlin", "Europe/Madrid", "Europe/Rome",
	"Europe/Amsterdam", "Europe/Stockholm", "Europe/Warsaw", "Europe/Moscow",
	"Europe/Istanbul", "Africa/Cairo", "Africa/Lagos", "Africa/Nairobi",
	"Asia/Dubai", "Asia/Karachi", "Asia/Kolkata", "Asia/Dhaka", "Asia/Bangkok",
	"Asia/Jakarta", "Asia/Singapore", "Asia/Shanghai", "Asia/Hong_Kong",
	"Asia/Tokyo", "Asia/Seoul", "Australia/Sydney", "Australia/Melbourne",
	"Australia/Brisbane", "Australia/Adelaide", "Australia/Perth",
	"Pacific/Auckland", "Pacific/Honolulu",
}

func tzMessageCommandHandler(s *discordgo.Session, m *discordgo.MessageCreate, args []string) error {
	if len(args) == 0 {
		_, err := s.ChannelMessageSend(m.ChannelID, "Usage: `-tz <timezone>` — e.g. `-tz America/New_York` or `-tz pst`")
		return err
	}
	embed := buildTzEmbed(strings.Join(args, "/"))
	_, err := s.ChannelMessageSendEmbed(m.ChannelID, embed)
	return err
}

func tzSlashCommandHandler(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	opts := OptionMap(i.ApplicationCommandData().Options)
	name := strings.TrimSpace(OptString(opts, "timezone"))
	if name == "" {
		return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{Content: "Usage: `/tz timezone:<zone>` — e.g. `America/New_York` or `pst`"},
		})
	}
	return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{buildTzEmbed(name)},
		},
	})
}

func buildTzEmbed(name string) *discordgo.MessageEmbed {
	key := strings.ToLower(strings.TrimSpace(name))
	if mapped, ok := tzAliases[key]; ok {
		name = mapped
	}

	loc, err := time.LoadLocation(name)
	if err != nil {
		msg := fmt.Sprintf("Unknown time zone `%s`.", name)
		if sug := suggestTz(name); sug != "" {
			msg += fmt.Sprintf("\nDid you mean `%s`?", sug)
		}
		return &discordgo.MessageEmbed{
			Color:       0xED4245,
			Title:       "Time Zone",
			Description: msg,
		}
	}

	now := time.Now().In(loc)
	zoneName, offset := now.Zone()
	sign := "+"
	if offset < 0 {
		sign = "-"
		offset = -offset
	}
	hh := offset / 3600
	mm := (offset % 3600) / 60

	return &discordgo.MessageEmbed{
		Color:       0x5865F2,
		Title:       "Current Time",
		Description: fmt.Sprintf("It is currently **%s** in **%s** (`%s`).", now.Format("15:04"), zoneName, loc.String()),
		Fields: []*discordgo.MessageEmbedField{
			{Name: "Local Time", Value: now.Format("3:04 PM"), Inline: true},
			{Name: "24-hour", Value: now.Format("15:04"), Inline: true},
			{Name: "UTC Offset", Value: fmt.Sprintf("UTC%s%02d:%02d", sign, hh, mm), Inline: true},
			{Name: "Date", Value: now.Format("Mon, 02 Jan 2006"), Inline: true},
		},
	}
}

// suggestTz tries to find a common time zone whose name contains the input,
// tolerating spaces in place of underscores.
func suggestTz(name string) string {
	key := strings.ToLower(strings.ReplaceAll(strings.TrimSpace(name), " ", "_"))
	if key == "" {
		return ""
	}
	for _, z := range commonTz {
		if strings.Contains(strings.ToLower(z), key) {
			return z
		}
	}
	return ""
}
