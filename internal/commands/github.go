package commands

import (
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
)

type githubUser struct {
	Login       string `json:"login"`
	Name        string `json:"name"`
	AvatarURL   string `json:"avatar_url"`
	Bio         string `json:"bio"`
	Followers   int    `json:"followers"`
	Following   int    `json:"following"`
	PublicRepos int    `json:"public_repos"`
	Location    string `json:"location"`
	Company     string `json:"company"`
	Blog        string `json:"blog"`
	HTMLURL     string `json:"html_url"`
	CreatedAt   string `json:"created_at"`
}

func fetchGitHubEmbed(username string) (*discordgo.MessageEmbed, error) {
	client := socialHTTPClient()
	var gh githubUser
	status, err := socialFetchJSON(client, "GET", "https://api.github.com/users/"+username, &gh)
	if err != nil {
		return nil, err
	}
	if status == 404 {
		return socialNotFoundEmbed("GitHub", username), nil
	}
	if status == 403 {
		return socialErrorEmbed("GitHub", "GitHub unauthenticated rate limit (60 req/hr) reached. Try again shortly."), nil
	}
	if status != 200 {
		return socialErrorEmbed("GitHub", fmt.Sprintf("GitHub API returned status %d.", status)), nil
	}

	name := gh.Name
	if name == "" {
		name = gh.Login
	}
	title := name + " (@" + gh.Login + ")"
	desc := gh.Bio
	if desc == "" {
		desc = "*(no bio)*"
	}

	fields := []*discordgo.MessageEmbedField{
		{Name: "Followers", Value: fmt.Sprintf("%d", gh.Followers), Inline: true},
		{Name: "Following", Value: fmt.Sprintf("%d", gh.Following), Inline: true},
		{Name: "Repos", Value: fmt.Sprintf("%d", gh.PublicRepos), Inline: true},
	}
	if gh.Location != "" {
		fields = append(fields, &discordgo.MessageEmbedField{Name: "Location", Value: gh.Location, Inline: true})
	}
	if gh.Company != "" {
		fields = append(fields, &discordgo.MessageEmbedField{Name: "Company", Value: gh.Company, Inline: true})
	}
	if gh.Blog != "" {
		fields = append(fields, &discordgo.MessageEmbedField{Name: "Blog", Value: gh.Blog, Inline: true})
	}

	var created string
	if t, err := time.Parse(time.RFC3339, gh.CreatedAt); err == nil {
		created = fmt.Sprintf("%d-%02d-%02d", t.Year(), t.Month(), t.Day())
	}
	if created != "" {
		fields = append(fields, &discordgo.MessageEmbedField{Name: "Joined", Value: created, Inline: true})
	}

	return &discordgo.MessageEmbed{
		Color:       0x6e40c9,
		Title:       title,
		URL:         gh.HTMLURL,
		Description: desc,
		Thumbnail:   &discordgo.MessageEmbedThumbnail{URL: gh.AvatarURL},
		Fields:      fields,
		Footer:      &discordgo.MessageEmbedFooter{Text: "GitHub"},
	}, nil
}
