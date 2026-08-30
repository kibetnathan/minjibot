package commands

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

const (
	giphyAPIHost     = "api.giphy.com"
	giphyDefaultKey  = "dc6zaTOxFJmzC" // Giphy's public beta key
	giphySearchLimit = 20
	GifSearchMaxDesc = 256
)

type GiphyMedia struct {
	URL string `json:"url"`
}

type GiphyImages struct {
	Original  GiphyMedia `json:"original"`
	Downsized GiphyMedia `json:"downsized"`
}

type GiphyUser struct {
	Username    string `json:"username"`
	DisplayName string `json:"display_name"`
}

type GiphyResult struct {
	ID     string      `json:"id"`
	Title  string      `json:"title"`
	URL    string      `json:"url"`
	Images GiphyImages `json:"images"`
	User   GiphyUser   `json:"user"`
}

type giphyResponse struct {
	Data []GiphyResult `json:"data"`
}

// GIFURL returns the raw GIF media URL (downsized is the largest embed-safe
// variant Giphy provides; original is used when downsized is missing).
func (r GiphyResult) GIFURL() string {
	if r.Images.Downsized.URL != "" {
		return r.Images.Downsized.URL
	}
	return r.Images.Original.URL
}

func (r GiphyResult) Description() string {
	return strings.TrimSpace(r.Title)
}

func (r GiphyResult) CreatorName() string {
	if s := strings.TrimSpace(r.User.DisplayName); s != "" {
		return s
	}
	return strings.TrimSpace(r.User.Username)
}

type giphyClient struct {
	http *http.Client
	key  string
}

func newGiphyClient() *giphyClient {
	key := strings.TrimSpace(os.Getenv("GIPHY_API_KEY"))
	if key == "" {
		key = giphyDefaultKey
	}
	return &giphyClient{
		http: &http.Client{Timeout: 10 * time.Second},
		key:  key,
	}
}

func (c *giphyClient) search(query string, limit int) ([]GiphyResult, error) {
	u := url.URL{Scheme: "https", Host: giphyAPIHost, Path: "/v1/gifs/search"}
	q := u.Query()
	q.Set("api_key", c.key)
	q.Set("q", query)
	q.Set("limit", fmt.Sprintf("%d", limit))
	q.Set("rating", "g")
	q.Set("lang", "en")
	q.Set("bundle", "messaging_non_clips")
	u.RawQuery = q.Encode()

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/126.0.0.0 Safari/537.36")
	req.Header.Set("Accept", "application/json")

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("giphy returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var out giphyResponse
	if err := json.Unmarshal(body, &out); err != nil {
		return nil, err
	}
	return out.Data, nil
}

// FilterByCreator keeps results whose uploader or title matches the creator
// query (case-insensitive substring).
func FilterByCreator(results []GiphyResult, creator string) []GiphyResult {
	creator = strings.ToLower(strings.TrimSpace(creator))
	if creator == "" {
		return results
	}

	out := make([]GiphyResult, 0, len(results))
	for _, r := range results {
		if strings.Contains(strings.ToLower(r.User.Username), creator) ||
			strings.Contains(strings.ToLower(r.User.DisplayName), creator) ||
			strings.Contains(strings.ToLower(r.Title), creator) {
			out = append(out, r)
		}
	}
	return out
}

// ParseGifSearchArgs splits out trailing creator:/by: tokens from the query.
func ParseGifSearchArgs(args []string) (query, creator string) {
	rest := make([]string, 0, len(args))
	for _, arg := range args {
		switch {
		case strings.HasPrefix(strings.ToLower(arg), "creator:"):
			creator = strings.TrimSpace(arg[len("creator:"):])
		case strings.HasPrefix(strings.ToLower(arg), "by:"):
			creator = strings.TrimSpace(arg[len("by:"):])
		default:
			rest = append(rest, arg)
		}
	}
	return strings.TrimSpace(strings.Join(rest, " ")), creator
}

func gifsearchMessageCommandHandler(s *discordgo.Session, channelID string, args []string) error {
	query, creator := ParseGifSearchArgs(args)
	if query == "" {
		_, err := s.ChannelMessageSend(channelID, "Usage: `!gifsearch <query> [creator:<name>]`")
		return err
	}

	results, err := newGiphyClient().search(query, giphySearchLimit)
	if err != nil {
		return err
	}
	results = FilterByCreator(results, creator)

	if len(results) == 0 {
		msg := fmt.Sprintf("No GIFs found for %q.", query)
		if creator != "" {
			msg = fmt.Sprintf("No GIFs found for %q by %q.", query, creator)
		}
		_, err := s.ChannelMessageSend(channelID, msg)
		return err
	}

	embed := BuildGifEmbed(query, creator, results[0])
	_, err = s.ChannelMessageSendEmbed(channelID, embed)
	return err
}

func gifsearchSlashCommandHandler(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	var query, creator string
	for _, opt := range i.ApplicationCommandData().Options {
		switch opt.Name {
		case "query":
			query = opt.StringValue()
		case "creator":
			creator = opt.StringValue()
		}
	}

	query = strings.TrimSpace(query)
	if query == "" {
		return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Usage: `/gifsearch query:<query> [creator:<name>]`",
			},
		})
	}

	results, err := newGiphyClient().search(query, giphySearchLimit)
	if err != nil {
		return err
	}
	results = FilterByCreator(results, creator)

	if len(results) == 0 {
		msg := fmt.Sprintf("No GIFs found for %q.", query)
		if creator != "" {
			msg = fmt.Sprintf("No GIFs found for %q by %q.", query, creator)
		}
		return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{Content: msg},
		})
	}

	embed := BuildGifEmbed(query, creator, results[0])
	return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{embed},
		},
	})
}

func BuildGifEmbed(query, creator string, r GiphyResult) *discordgo.MessageEmbed {
	title := r.Description()
	if title == "" {
		title = query
	}
	if len(title) > GifSearchMaxDesc {
		title = title[:GifSearchMaxDesc-3] + "..."
	}

	embed := &discordgo.MessageEmbed{
		Color: 0x46E07A, // Giphy green
		Title: title,
	}
	if u := r.GIFURL(); u != "" {
		embed.Image = &discordgo.MessageEmbedImage{URL: u}
	}

	credit := "GIF via [Giphy](https://giphy.com)"
	if r.User.Username != "" {
		credit += fmt.Sprintf(" • by [%s](https://giphy.com/@%s)", r.CreatorName(), r.User.Username)
	}
	if creator != "" {
		credit += fmt.Sprintf(" • creator filter: %q", creator)
	}
	embed.Footer = &discordgo.MessageEmbedFooter{Text: credit}
	return embed
}
