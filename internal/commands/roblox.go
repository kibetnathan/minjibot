package commands

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/bwmarrin/discordgo"
)

type robloxUsernameRequest struct {
	Usernames          []string `json:"usernames"`
	ExcludeBannedUsers bool     `json:"excludeBannedUsers"`
}

type robloxUsernameResponse struct {
	Data []struct {
		ID          int64  `json:"id"`
		Name        string `json:"name"`
		DisplayName string `json:"displayName"`
	} `json:"data"`
}

type robloxUserInfo struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	DisplayName string `json:"displayName"`
	Description string `json:"description"`
	Created     string `json:"created"`
	IsBanned    bool   `json:"isBanned"`
}

type robloxAvatarResponse struct {
	Data []struct {
		ImageURL string `json:"imageUrl"`
	} `json:"data"`
}

type robloxFollowersResponse struct {
	Count int64 `json:"count"`
}

func fetchRobloxEmbed(username string) (*discordgo.MessageEmbed, error) {
	client := socialHTTPClient()

	// Step 1: resolve username to ID
	reqBody, _ := json.Marshal(robloxUsernameRequest{
		Usernames:          []string{username},
		ExcludeBannedUsers: true,
	})
	req, err := http.NewRequest("POST", "https://users.roblox.com/v1/usernames/users", bytes.NewReader(reqBody))
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", socialUserAgent)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode == 400 {
		return socialNotFoundEmbed("Roblox", username), nil
	}
	if resp.StatusCode != 200 {
		return socialErrorEmbed("Roblox", fmt.Sprintf("Roblox API returned status %d.", resp.StatusCode)), nil
	}

	var nameResp robloxUsernameResponse
	if err := json.Unmarshal(body, &nameResp); err != nil {
		return nil, err
	}
	if len(nameResp.Data) == 0 {
		return socialNotFoundEmbed("Roblox", username), nil
	}
	userID := nameResp.Data[0].ID

	// Step 2: get full user info
	var info robloxUserInfo
	status, err := socialFetchJSON(client, "GET", fmt.Sprintf("https://users.roblox.com/v1/users/%d", userID), &info)
	if err != nil {
		return nil, err
	}
	if status != 200 {
		return socialErrorEmbed("Roblox", fmt.Sprintf("Could not fetch user info (status %d).", status)), nil
	}

	// Step 3: avatar headshot
	var avatar robloxAvatarResponse
	avatarURL := ""
	socialFetchJSON(client, "GET", fmt.Sprintf("https://thumbnails.roblox.com/v1/users/avatar-headshot?userIds=%d&size=48x48&format=Png&isCircular=false", userID), &avatar)
	if len(avatar.Data) > 0 {
		avatarURL = avatar.Data[0].ImageURL
	}

	// Step 4: follower count
	var followers robloxFollowersResponse
	socialFetchJSON(client, "GET", fmt.Sprintf("https://friends.roblox.com/v1/users/%d/followers/count", userID), &followers)

	desc := info.Description
	if desc == "" {
		desc = "*(no description)*"
	}

	fields := []*discordgo.MessageEmbedField{
		{Name: "User ID", Value: fmt.Sprintf("%d", userID), Inline: true},
		{Name: "Followers", Value: fmt.Sprintf("%d", followers.Count), Inline: true},
	}
	if t, err := time.Parse(time.RFC3339, info.Created); err == nil {
		fields = append(fields, &discordgo.MessageEmbedField{
			Name: "Joined", Value: fmt.Sprintf("%d-%02d-%02d", t.Year(), t.Month(), t.Day()), Inline: true,
		})
	}
	if info.IsBanned {
		fields = append(fields, &discordgo.MessageEmbedField{Name: "Status", Value: "Banned", Inline: true})
	}

	profileURL := fmt.Sprintf("https://www.roblox.com/users/%d/profile", userID)
	title := info.DisplayName + " (@" + info.Name + ")"

	return &discordgo.MessageEmbed{
		Color:       0x5865F2,
		Title:       title,
		URL:         profileURL,
		Description: desc,
		Thumbnail:   &discordgo.MessageEmbedThumbnail{URL: avatarURL},
		Fields:      fields,
		Footer:      &discordgo.MessageEmbedFooter{Text: "Roblox"},
	}, nil
}
