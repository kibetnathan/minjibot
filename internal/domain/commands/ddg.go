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

// ddgTopic models a DuckDuckGo Instant Answer topic. The API interleaves
// two shapes in RelatedTopics: a plain topic (Text/FirstURL) and a category
// wrapper (Name/Topics), so both sets of fields live on one struct.
type ddgTopic struct {
	Text     string `json:"Text"`
	FirstURL string `json:"FirstURL"`
	Icon     struct {
		URL string `json:"URL"`
	} `json:"Icon"`
	Name   string     `json:"Name"`
	Topics []ddgTopic `json:"Topics"`
}

type ddgResponse struct {
	Heading        string     `json:"Heading"`
	AbstractText   string     `json:"AbstractText"`
	AbstractURL    string     `json:"AbstractURL"`
	AbstractSource string     `json:"AbstractSource"`
	Image          string     `json:"Image"`
	Answer         string     `json:"Answer"`
	Definition     string     `json:"Definition"`
	DefinitionURL  string     `json:"DefinitionURL"`
	RelatedTopics  []ddgTopic `json:"RelatedTopics"`
	Results        []ddgTopic `json:"Results"`
}

type ddgClient struct {
	http *http.Client
}

func newDDGClient() *ddgClient {
	return &ddgClient{http: &http.Client{Timeout: 10 * time.Second}}
}

func (c *ddgClient) Search(query string) (*ddgResponse, error) {
	u := url.URL{
		Scheme: "https",
		Host:   "api.duckduckgo.com",
	}
	q := u.Query()
	q.Set("q", query)
	q.Set("format", "json")
	q.Set("no_html", "1")
	q.Set("skip_disambig", "1")
	u.RawQuery = q.Encode()

	resp, err := c.do(u.String())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var out ddgResponse
	if err := json.Unmarshal(body, &out); err != nil {
		return nil, err
	}

	return &out, nil
}

func (c *ddgClient) do(rawURL string) (*http.Response, error) {
	for attempt := 1; ; attempt++ {
		req, err := http.NewRequest(http.MethodGet, rawURL, nil)
		if err != nil {
			return nil, err
		}
		req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/126.0.0.0 Safari/537.36")
		req.Header.Set("Accept", "application/json")
		req.Header.Set("Connection", "keep-alive")

		resp, err := c.http.Do(req)
		if err != nil {
			return nil, err
		}

		// DuckDuckGo answers 202 when it wants to run an anomaly/redirect
		// check. Retry once per docs; if it persists, surface the status.
		if resp.StatusCode == http.StatusAccepted && attempt < 2 {
			resp.Body.Close()
			time.Sleep(1 * time.Second)
			continue
		}

		if resp.StatusCode != http.StatusOK {
			defer resp.Body.Close()
			return nil, fmt.Errorf("duckduckgo returned status %d", resp.StatusCode)
		}

		return resp, nil
	}
}

func (r *ddgResponse) FlattenedTopics() []ddgTopic {
	var out []ddgTopic
	var walk func(items []ddgTopic)
	walk = func(items []ddgTopic) {
		for _, t := range items {
			if len(t.Topics) > 0 {
				walk(t.Topics)
			} else if t.Text != "" {
				out = append(out, t)
			}
		}
	}
	walk(r.RelatedTopics)
	return out
}

func ddgMessageCommandHandler(s *discordgo.Session, channelID string, args []string) error {
	query := strings.TrimSpace(strings.Join(args, " "))
	if query == "" {
		_, err := s.ChannelMessageSend(channelID, "Usage: `!ddg <query>`")
		return err
	}

	res, err := newDDGClient().Search(query)
	if err != nil {
		return err
	}

	embed := buildDDGEmbed(query, res)
	_, err = s.ChannelMessageSendEmbed(channelID, embed)
	return err
}

func ddgSlashCommandHandler(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	var query string
	for _, opt := range i.ApplicationCommandData().Options {
		if opt.Name == "query" {
			query = opt.StringValue()
		}
	}

	query = strings.TrimSpace(query)
	if query == "" {
		return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Usage: `/ddg query:<query>`",
			},
		})
	}

	res, err := newDDGClient().Search(query)
	if err != nil {
		return err
	}

	embed := buildDDGEmbed(query, res)
	return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{embed},
		},
	})
}

func buildDDGEmbed(query string, res *ddgResponse) *discordgo.MessageEmbed {
	embed := &discordgo.MessageEmbed{
		Color: 0xDE5833, // DuckDuckGo orange
		Title: fmt.Sprintf("DuckDuckGo: %s", query),
	}

	if res.Image != "" {
		embed.Thumbnail = &discordgo.MessageEmbedThumbnail{URL: res.Image}
	}

	if res.AbstractText != "" {
		abstract := res.AbstractText
		if res.AbstractURL != "" {
			abstract = fmt.Sprintf("[%s](%s)", abstract, res.AbstractURL)
		}
		embed.Description = abstract
	} else if res.Answer != "" {
		embed.Description = res.Answer
	} else if res.Definition != "" {
		definition := res.Definition
		if res.DefinitionURL != "" {
			definition = fmt.Sprintf("[%s](%s)", definition, res.DefinitionURL)
		}
		embed.Description = definition
	}

	topics := res.FlattenedTopics()
	if len(topics) == 0 {
		topics = res.Results
	}

	limit := 5
	if len(topics) < limit {
		limit = len(topics)
	}
	embed.Fields = make([]*discordgo.MessageEmbedField, 0, limit)
	for _, t := range topics[:limit] {
		title := t.Text
		desc := ""
		if idx := strings.Index(t.Text, " - "); idx != -1 {
			title = t.Text[:idx]
			desc = strings.TrimSpace(t.Text[idx+3:])
		}
		value := desc
		if t.FirstURL != "" {
			if value != "" {
				value += "\n"
			}
			value += t.FirstURL
		}
		if value == "" {
			value = "—"
		}
		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
			Name:  title,
			Value: value,
		})
	}

	if len(topics) == 0 {
		if embed.Description == "" {
			embed.Description = "No results found."
		}
		embed.Fields = nil
	}

	embed.Footer = &discordgo.MessageEmbedFooter{
		Text: "Results from DuckDuckGo",
	}
	return embed
}
