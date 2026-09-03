package commands

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/kibetnathan/minjibot/internal/config"
)

type redditAboutResponse struct {
	Data struct {
		Name         string  `json:"name"`
		IconImg      string  `json:"icon_img"`
		LinkKarma    int     `json:"link_karma"`
		CommentKarma int     `json:"comment_karma"`
		CreatedUTC   float64 `json:"created_utc"`
		Subreddit    struct {
			Title             string `json:"title"`
			PublicDescription string `json:"public_description"`
		} `json:"subreddit"`
	} `json:"data"`
}

type redditTokenCache struct {
	mu        sync.Mutex
	token     string
	expiresAt time.Time
}

var redditCache redditTokenCache

func redditAccessToken(cfg *config.Config) (string, error) {
	redditCache.mu.Lock()
	defer redditCache.mu.Unlock()

	if redditCache.token != "" && time.Now().Before(redditCache.expiresAt) {
		return redditCache.token, nil
	}

	body := url.Values{}
	body.Set("grant_type", "client_credentials")

	req, err := http.NewRequest("POST", "https://www.reddit.com/api/v1/access_token", strings.NewReader(body.Encode()))
	if err != nil {
		return "", err
	}
	req.SetBasicAuth(cfg.RedditClientID, cfg.RedditClientSecret)
	req.Header.Set("User-Agent", socialUserAgent)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := socialHTTPClient().Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		return "", fmt.Errorf("reddit token request failed (status %d)", resp.StatusCode)
	}

	var tok struct {
		AccessToken string `json:"access_token"`
		ExpiresIn   int    `json:"expires_in"`
	}
	if err := json.Unmarshal(respBody, &tok); err != nil {
		return "", err
	}
	if tok.AccessToken == "" {
		return "", fmt.Errorf("reddit returned empty access token")
	}

	redditCache.token = tok.AccessToken
	redditCache.expiresAt = time.Now().Add(time.Duration(tok.ExpiresIn)*time.Second - 30*time.Second)
	return tok.AccessToken, nil
}

func fetchRedditEmbed(username string, cfg *config.Config) (*discordgo.MessageEmbed, error) {
	if cfg == nil || cfg.RedditClientID == "" || cfg.RedditClientSecret == "" {
		return socialErrorEmbed("Reddit", "Reddit lookup requires `REDDIT_CLIENT_ID` and `REDDIT_CLIENT_SECRET` in `.env`. Create a script app at https://www.reddit.com/prefs/apps."), nil
	}

	token, err := redditAccessToken(cfg)
	if err != nil {
		return socialErrorEmbed("Reddit", fmt.Sprintf("Failed to authenticate with Reddit: %s", err)), nil
	}

	client := socialHTTPClient()
	req, err := socialNewRequest("GET", "https://oauth.reddit.com/user/"+username+"/about")
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 404 {
		return socialNotFoundEmbed("Reddit", username), nil
	}
	if resp.StatusCode == 403 {
		return socialErrorEmbed("Reddit", "Reddit API returned 403. Check your credentials in `.env`."), nil
	}
	if resp.StatusCode == 429 {
		return socialErrorEmbed("Reddit", "Reddit rate limit hit. Try again shortly."), nil
	}
	if resp.StatusCode != 200 {
		return socialErrorEmbed("Reddit", fmt.Sprintf("Reddit API returned status %d.", resp.StatusCode)), nil
	}

	var about redditAboutResponse
	if err := json.NewDecoder(resp.Body).Decode(&about); err != nil {
		return nil, err
	}

	d := about.Data
	name := "u/" + d.Name
	desc := d.Subreddit.PublicDescription
	if desc == "" {
		desc = "*(no description)*"
	}

	iconURL := d.IconImg
	if iconURL == "" {
		iconURL = "https://www.redditstatic.com/avatars/avatar_default_1.png"
	}

	created := ""
	if d.CreatedUTC > 0 {
		t := time.Unix(int64(d.CreatedUTC), 0)
		created = fmt.Sprintf("%d-%02d-%02d", t.Year(), t.Month(), t.Day())
	}

	fields := []*discordgo.MessageEmbedField{
		{Name: "Post Karma", Value: fmt.Sprintf("%d", d.LinkKarma), Inline: true},
		{Name: "Comment Karma", Value: fmt.Sprintf("%d", d.CommentKarma), Inline: true},
	}
	if created != "" {
		fields = append(fields, &discordgo.MessageEmbedField{Name: "Cake Day", Value: created, Inline: true})
	}

	return &discordgo.MessageEmbed{
		Color:       0xFF4500,
		Title:       name,
		URL:         "https://www.reddit.com/user/" + username,
		Description: desc,
		Thumbnail:   &discordgo.MessageEmbedThumbnail{URL: iconURL},
		Fields:      fields,
		Footer:      &discordgo.MessageEmbedFooter{Text: "Reddit"},
	}, nil
}
