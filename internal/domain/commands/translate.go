package commands

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

type translateClient struct {
	http *http.Client
}

func newTranslateClient() *translateClient {
	return &translateClient{http: &http.Client{Timeout: 10 * time.Second}}
}

// translate uses the public Google Translate gtx endpoint (no API key needed).
func (c *translateClient) translate(text, target, source string) (string, error) {
	u := url.URL{Scheme: "https", Host: "translate.googleapis.com", Path: "/translate_a/single"}
	q := u.Query()
	q.Set("client", "gtx")
	q.Set("sl", source) // "auto" when the caller wants auto-detection
	q.Set("tl", target)
	q.Set("dt", "t")
	q.Set("q", text)
	u.RawQuery = q.Encode()

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/126.0.0.0 Safari/537.36")

	resp, err := c.http.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("translate returned status %d", resp.StatusCode)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	// Response is a nested array; the first element is a list of
	// [translated-segment, original-segment, ...] tuples.
	var raw []json.RawMessage
	if err := json.Unmarshal(body, &raw); err != nil {
		return "", err
	}
	var segments []json.RawMessage
	if err := json.Unmarshal(raw[0], &segments); err != nil {
		return "", err
	}

	var b strings.Builder
	for _, segRaw := range segments {
		var seg []json.RawMessage
		if err := json.Unmarshal(segRaw, &seg); err != nil {
			continue
		}
		var part string
		if len(seg) > 0 && json.Unmarshal(seg[0], &part) == nil {
			b.WriteString(part)
		}
	}
	return strings.TrimSpace(b.String()), nil
}

func parseTranslateArgs(args []string) (target, text string) {
	target = "en"
	var rest []string
	for _, arg := range args {
		switch {
		case strings.HasPrefix(strings.ToLower(arg), "to:"):
			target = strings.ToLower(strings.TrimSpace(arg[len("to:"):]))
		case strings.HasPrefix(strings.ToLower(arg), "lang:"):
			target = strings.ToLower(strings.TrimSpace(arg[len("lang:"):]))
		default:
			rest = append(rest, arg)
		}
	}
	return target, strings.TrimSpace(strings.Join(rest, " "))
}

func translateMessageCommandHandler(s *discordgo.Session, channelID string, args []string) error {
	target, text := parseTranslateArgs(args)
	if text == "" {
		_, err := s.ChannelMessageSend(channelID, "Usage: `!translate [to:<lang>] <text>`")
		return err
	}

	result, err := newTranslateClient().translate(text, target, "auto")
	if err != nil {
		return err
	}

	embed := &discordgo.MessageEmbed{
		Color: 0x5865F2,
		Title: fmt.Sprintf("Translated to %s", target),
		Fields: []*discordgo.MessageEmbedField{
			{Name: "Original", Value: truncateForEmbed(text, 1024)},
			{Name: "Translation", Value: truncateForEmbed(result, 1024)},
		},
	}
	_, err = s.ChannelMessageSendEmbed(channelID, embed)
	return err
}

func translateSlashCommandHandler(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	opts := optionMap(i.ApplicationCommandData().Options)
	text := strings.TrimSpace(optString(opts, "text"))
	if text == "" {
		return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{Content: "Usage: `/translate text:<text> [target:<lang>]`"},
		})
	}
	target := optString(opts, "target")
	if target == "" {
		target = "en"
	}

	result, err := newTranslateClient().translate(text, target, "auto")
	if err != nil {
		return err
	}

	embed := &discordgo.MessageEmbed{
		Color: 0x5865F2,
		Title: fmt.Sprintf("Translated to %s", target),
		Fields: []*discordgo.MessageEmbedField{
			{Name: "Original", Value: truncateForEmbed(text, 1024)},
			{Name: "Translation", Value: truncateForEmbed(result, 1024)},
		},
	}
	return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{Embeds: []*discordgo.MessageEmbed{embed}},
	})
}

func truncateForEmbed(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max-3] + "..."
}
