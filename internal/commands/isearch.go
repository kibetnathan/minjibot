package commands

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

const iSearchMaxResults = 4

type isearchResult struct {
	Image string `json:"image"`
	Title string `json:"title"`
	URL   string `json:"url"`
}

type isearchResponse struct {
	Results []isearchResult `json:"results"`
}

type isearchClient struct {
	http *http.Client
}

func newISearchClient() *isearchClient {
	return &isearchClient{http: &http.Client{Timeout: 15 * time.Second}}
}

// search uses DuckDuckGo's public image search JSON endpoint (no API key).
func (c *isearchClient) search(query string) ([]isearchResult, error) {
	u := url.URL{Scheme: "https", Host: "duckduckgo.com", Path: "/i.js"}
	q := u.Query()
	q.Set("q", query)
	q.Set("o", "json")
	q.Set("p", "1")
	q.Set("l", "us-en")
	q.Set("f", ",,,,")
	u.RawQuery = q.Encode()

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/126.0.0.0 Safari/537.36")
	req.Header.Set("Referer", "https://duckduckgo.com/")
	req.Header.Set("Accept", "application/json")

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("image search returned status %d", resp.StatusCode)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var out isearchResponse
	if err := json.Unmarshal(body, &out); err != nil {
		return nil, err
	}
	return out.Results, nil
}

func isearchMessageCommandHandler(s *discordgo.Session, channelID string, args []string) error {
	query := strings.TrimSpace(strings.Join(args, " "))
	if query == "" {
		_, err := s.ChannelMessageSend(channelID, "Usage: `!isearch <query>`")
		return err
	}

	results, err := newISearchClient().search(query)
	if err != nil {
		return err
	}
	if len(results) == 0 {
		_, err := s.ChannelMessageSend(channelID, fmt.Sprintf("No image results for %q.", query))
		return err
	}

	count := iSearchMaxResults
	if len(results) < count {
		count = len(results)
	}

	files := make([]*discordgo.File, 0, count)
	details := make([]string, 0, count)
	for _, r := range results[:count] {
		md, err := fetchURL(r.Image)
		if err != nil {
			continue
		}
		ext := "png"
		switch md.Ext {
		case "jpg", "jpeg", "gif", "webp":
			ext = md.Ext
		case "bin":
			ext = "png"
		}
		filename := fmt.Sprintf("isearch_%d.%s", len(files), ext)
		files = append(files, &discordgo.File{Name: filename, ContentType: "image/" + ext, Reader: bytes.NewReader(md.Data)})
		details = append(details, r.URL)
	}

	if len(files) == 0 {
		_, err := s.ChannelMessageSend(channelID, fmt.Sprintf("Couldn't download any image results for %q.", query))
		return err
	}

	_, err = s.ChannelMessageSendComplex(channelID, &discordgo.MessageSend{
		Content: fmt.Sprintf("Image results for **%s**", query),
		Files:   files,
	})
	return err
}

func isearchSlashCommandHandler(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	opts := OptionMap(i.ApplicationCommandData().Options)
	query := strings.TrimSpace(OptString(opts, "query"))
	if query == "" {
		return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{Content: "Usage: `/isearch query:<query>`"},
		})
	}

	results, err := newISearchClient().search(query)
	if err != nil {
		return err
	}
	if len(results) == 0 {
		return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{Content: fmt.Sprintf("No image results for %q.", query)},
		})
	}

	count := iSearchMaxResults
	if len(results) < count {
		count = len(results)
	}

	files := make([]*discordgo.File, 0, count)
	for _, r := range results[:count] {
		md, err := fetchURL(r.Image)
		if err != nil {
			continue
		}
		ext := "png"
		switch md.Ext {
		case "jpg", "jpeg", "gif", "webp":
			ext = md.Ext
		case "bin":
			ext = "png"
		}
		files = append(files, &discordgo.File{Name: fmt.Sprintf("isearch_%d.%s", len(files), ext), ContentType: "image/" + ext, Reader: bytes.NewReader(md.Data)})
	}

	if len(files) == 0 {
		return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{Content: fmt.Sprintf("Couldn't download any image results for %q.", query)},
		})
	}

	return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf("Image results for **%s**", query),
			Files:   files,
		},
	})
}
