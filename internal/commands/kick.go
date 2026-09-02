package commands

import (
	"encoding/json"
	"fmt"

	"github.com/bwmarrin/discordgo"
)

type kickChannelResponse struct {
	ID                  int64           `json:"id"`
	Slug                string          `json:"slug"`
	IsBanned            bool            `json:"is_banned"`
	FollowersCount      string          `json:"followers_count"`
	Verified            bool            `json:"verified"`
	SubscriptionEnabled bool            `json:"subscription_enabled"`
	BannerImage         kickImage       `json:"banner_image"`
	User                kickUser        `json:"user"`
	Livestream          json.RawMessage `json:"livestream"`
}

type kickUser struct {
	ID         int64  `json:"id"`
	Username   string `json:"username"`
	Bio        string `json:"bio"`
	ProfilePic string `json:"profile_pic"`
}

type kickImage struct {
	URL string `json:"url"`
}

type kickLivestream struct {
	ViewerCount int `json:"viewer_count"`
}

func fetchKickEmbed(slug string) (*discordgo.MessageEmbed, error) {
	client := socialHTTPClient()
	var ch kickChannelResponse
	status, err := socialFetchJSON(client, "GET", "https://kick.com/api/v2/channels/"+slug, &ch)
	if err != nil {
		return nil, err
	}
	if status == 404 {
		return socialNotFoundEmbed("Kick", slug), nil
	}
	if status != 200 {
		return socialErrorEmbed("Kick", fmt.Sprintf("Kick API returned status %d.", status)), nil
	}

	desc := ch.User.Bio
	if desc == "" {
		desc = "*(no bio)*"
	}

	fields := []*discordgo.MessageEmbedField{
		{Name: "Followers", Value: ch.FollowersCount, Inline: true},
	}
	if ch.Verified {
		fields = append(fields, &discordgo.MessageEmbedField{Name: "Verified", Value: "Yes", Inline: true})
	}

	// Livestream is `null` or `false` when offline, an object when live
	if len(ch.Livestream) > 0 && ch.Livestream[0] == '{' {
		var ls kickLivestream
		if json.Unmarshal(ch.Livestream, &ls) == nil {
			fields = append(fields, &discordgo.MessageEmbedField{
				Name: "Live", Value: fmt.Sprintf("**%d** viewers", ls.ViewerCount), Inline: true,
			})
		}
	}

	profileURL := "https://kick.com/" + ch.Slug

	return &discordgo.MessageEmbed{
		Color:       0x53F34B,
		Title:       ch.User.Username,
		URL:         profileURL,
		Description: desc,
		Thumbnail:   &discordgo.MessageEmbedThumbnail{URL: ch.User.ProfilePic},
		Fields:      fields,
		Footer:      &discordgo.MessageEmbedFooter{Text: "Kick"},
	}, nil
}
