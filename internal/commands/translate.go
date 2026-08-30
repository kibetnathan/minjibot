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
	return &translateClient{http: &http.Client{Timeout: 15 * time.Second}}
}

// google translates using the public (keyless) Google Translate web endpoints.
// The gtx client is preferred, but Google occasionally 403s datacenter IPs, so
// we fall back across alternative hosts/clients before giving up.
func (c *translateClient) translate(text, target, source string) (string, error) {
	attempts := []struct {
		host   string
		path   string
		client string
	}{
		{"translate.googleapis.com", "/translate_a/single", "gtx"},
		{"translate.googleapis.com", "/translate_a/single", "dict-chrome-ex"},
		{"clients5.google.com", "/translate_a/t", "dict-chrome-ex"},
	}

	var lastErr error
	for _, a := range attempts {
		out, err := c.call(a.host, a.path, a.client, text, target, source)
		if err == nil {
			return out, nil
		}
		lastErr = err
	}
	return "", lastErr
}

func (c *translateClient) call(host, path, client, text, target, source string) (string, error) {
	u := url.URL{Scheme: "https", Host: host, Path: path}
	q := u.Query()
	q.Set("client", client)
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
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("Referer", "https://translate.google.com/")
	req.Header.Set("Origin", "https://translate.google.com")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")

	resp, err := c.http.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		// 403/429 are retryable against a different endpoint; anything else fatal.
		if resp.StatusCode == http.StatusForbidden || resp.StatusCode == http.StatusTooManyRequests {
			return "", errRetryable
		}
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
	if len(raw) == 0 {
		return "", errRetryable
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
	if b.Len() == 0 {
		return "", errRetryable
	}
	return strings.TrimSpace(b.String()), nil
}

// errRetryable signals an unsupported path that should be retried against a
// different translate endpoint.
var errRetryable = fmt.Errorf("translate endpoint rejected request")

func ParseTranslateArgs(args []string) (target, text string) {
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

func translateMessageCommandHandler(s *discordgo.Session, m *discordgo.MessageCreate, args []string) error {
	target, text := ParseTranslateArgs(args)
	if text == "" {
		text = referencedMessageContent(s, m)
	}
	if text == "" {
		_, err := s.ChannelMessageSend(m.ChannelID, "Usage: `-translate [to:<lang>] <text>` — or reply to a message to translate it")
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
			{Name: "Original", Value: TruncateForEmbed(text, 1024)},
			{Name: "Translation", Value: TruncateForEmbed(result, 1024)},
		},
	}
	_, err = s.ChannelMessageSendEmbed(m.ChannelID, embed)
	return err
}

func translateSlashCommandHandler(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	opts := OptionMap(i.ApplicationCommandData().Options)
	text := strings.TrimSpace(OptString(opts, "text"))
	if text == "" {
		return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{Content: "Usage: `/translate text:<text> [target:<lang>]`"},
		})
	}
	target := OptString(opts, "target")
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
			{Name: "Original", Value: TruncateForEmbed(text, 1024)},
			{Name: "Translation", Value: TruncateForEmbed(result, 1024)},
		},
	}
	return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{Embeds: []*discordgo.MessageEmbed{embed}},
	})
}

func TruncateForEmbed(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max-3] + "..."
}
