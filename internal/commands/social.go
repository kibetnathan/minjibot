package commands

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/bwmarrin/discordgo"
)

const socialUserAgent = "MinjiBot/1.0 (https://github.com/kibetnathan/minjibot)"

func socialHTTPClient() *http.Client {
	return &http.Client{Timeout: 15 * time.Second}
}

func socialNewRequest(method, url string) (*http.Request, error) {
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", socialUserAgent)
	return req, nil
}

func socialFetchJSON(client *http.Client, method, url string, out any) (int, error) {
	req, err := socialNewRequest(method, url)
	if err != nil {
		return 0, err
	}
	resp, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return resp.StatusCode, nil
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return resp.StatusCode, err
	}
	return resp.StatusCode, json.Unmarshal(body, out)
}

func socialNotFoundEmbed(platform, username string) *discordgo.MessageEmbed {
	return &discordgo.MessageEmbed{
		Color:       0xED4245,
		Title:       platform,
		Description: fmt.Sprintf("User `%s` not found on %s.", username, platform),
	}
}

func socialErrorEmbed(platform, msg string) *discordgo.MessageEmbed {
	return &discordgo.MessageEmbed{
		Color:       0xED4245,
		Title:       platform,
		Description: msg,
	}
}
