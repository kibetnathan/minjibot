package commands

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
	_ "time/tzdata"

	"github.com/bwmarrin/discordgo"
)

// openMeteoGeocodeURL resolves a place name to a time zone. Free, no API key.
const openMeteoGeocodeURL = "https://geocoding-api.open-meteo.com/v1/search"

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

// tzPlace is a geocoded time zone lookup result.
type tzPlace struct {
	Name     string
	Timezone string
}

func tzMessageCommandHandler(s *discordgo.Session, m *discordgo.MessageCreate, args []string) error {
	if len(args) == 0 {
		_, err := s.ChannelMessageSend(m.ChannelID, "Usage: `-tz <place>` — e.g. `-tz Tokyo` or `-tz New York`")
		return err
	}
	embed := buildTzEmbed(strings.Join(args, " "))
	_, err := s.ChannelMessageSendEmbed(m.ChannelID, embed)
	return err
}

func tzSlashCommandHandler(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	opts := OptionMap(i.ApplicationCommandData().Options)
	place := strings.TrimSpace(OptString(opts, "place"))
	if place == "" {
		return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{Content: "Usage: `/tz place:<place>` — e.g. `Tokyo` or `New York`"},
		})
	}
	embed := buildTzEmbed(place)
	return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{embed},
		},
	})
}

// buildTzEmbed shows the current local time for a place. It prefers a direct
// zone name / alias (fast, offline), then falls back to geocoding the place.
func buildTzEmbed(place string) *discordgo.MessageEmbed {
	query := strings.TrimSpace(place)

	key := strings.ToLower(query)
	if mapped, ok := tzAliases[key]; ok {
		return tzEmbedFor(query, mapped, query)
	}
	if strings.Contains(query, "/") {
		if loc := loadTzLocation(query); loc != nil {
			return tzEmbedFor(query, loc.String(), query)
		}
	}

	client := &http.Client{Timeout: 15 * time.Second}
	p, err := geocodePlace(client, query)
	if err != nil {
		return &discordgo.MessageEmbed{
			Color:       0xED4245,
			Title:       "Time Zone",
			Description: "Couldn't reach the time zone lookup service. Try again in a moment.",
		}
	}
	if p != nil {
		if loc := loadTzLocation(p.Timezone); loc != nil {
			return tzEmbedFor(p.Name, loc.String(), query)
		}
	}

	if loc := loadTzLocation(query); loc != nil {
		return tzEmbedFor(query, loc.String(), query)
	}

	msg := fmt.Sprintf("Couldn't find a place called `%s`.", query)
	if sug := suggestTz(query); sug != "" {
		msg += fmt.Sprintf("\nDid you mean `%s`?", sug)
	}
	return &discordgo.MessageEmbed{
		Color:       0xED4245,
		Title:       "Time Zone",
		Description: msg,
	}
}

// geocodePlace resolves a place name via Open-Meteo's free geocoding API,
// returning the best match with its IANA time zone. Returns (nil, nil) when
// the service has no match for the query.
func geocodePlace(client *http.Client, query string) (*tzPlace, error) {
	endpoint := fmt.Sprintf("%s?name=%s&count=5&language=en", openMeteoGeocodeURL, url.QueryEscape(query))
	resp, err := client.Get(endpoint)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("geocoding service returned status %d", resp.StatusCode)
	}

	var out struct {
		Results []struct {
			Name     string `json:"name"`
			Timezone string `json:"timezone"`
		} `json:"results"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, err
	}
	for _, r := range out.Results {
		if r.Timezone != "" && r.Name != "" {
			return &tzPlace{Name: r.Name, Timezone: r.Timezone}, nil
		}
	}
	return nil, nil
}

func loadTzLocation(name string) *time.Location {
	loc, err := time.LoadLocation(name)
	if err != nil {
		return nil
	}
	return loc
}

func tzEmbedFor(placeName, zoneName, query string) *discordgo.MessageEmbed {
	loc, err := time.LoadLocation(zoneName)
	if err != nil {
		// Shouldn't happen, but degrade gracefully rather than panic.
		return &discordgo.MessageEmbed{
			Color:       0xED4245,
			Title:       "Time Zone",
			Description: fmt.Sprintf("Unknown time zone `%s`.", zoneName),
		}
	}

	now := time.Now().In(loc)
	zoneAbbrev, offset := now.Zone()
	sign := "+"
	if offset < 0 {
		sign = "-"
		offset = -offset
	}
	hh := offset / 3600
	mm := (offset % 3600) / 60

	return &discordgo.MessageEmbed{
		Color:       0x5865F2,
		Title:       fmt.Sprintf("Current Time in %s", placeName),
		Description: fmt.Sprintf("It is currently **%s** (`%s`).", now.Format("15:04"), zoneName),
		Fields: []*discordgo.MessageEmbedField{
			{Name: "Local Time", Value: now.Format("3:04 PM"), Inline: true},
			{Name: "24-hour", Value: now.Format("15:04"), Inline: true},
			{Name: "Zone", Value: zoneAbbrev, Inline: true},
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
