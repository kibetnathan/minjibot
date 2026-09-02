package commands

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/kibetnathan/minjibot/internal/config"
)

type bioProvider struct {
	name  string
	build func(string, *config.Config) (*discordgo.MessageEmbed, error)
}

var bioProviders = []bioProvider{
	{"github", func(u string, _ *config.Config) (*discordgo.MessageEmbed, error) { return fetchGitHubEmbed(u) }},
	{"roblox", func(u string, _ *config.Config) (*discordgo.MessageEmbed, error) { return fetchRobloxEmbed(u) }},
	{"reddit", func(u string, cfg *config.Config) (*discordgo.MessageEmbed, error) { return fetchRedditEmbed(u, cfg) }},
	{"kick", func(u string, _ *config.Config) (*discordgo.MessageEmbed, error) { return fetchKickEmbed(u) }},
}

func bioProviderByName(name string) (bioProvider, bool) {
	for _, p := range bioProviders {
		if p.name == name {
			return p, true
		}
	}
	return bioProvider{}, false
}

func bioUsageEmbed() *discordgo.MessageEmbed {
	var names []string
	for _, p := range bioProviders {
		names = append(names, p.name)
	}
	return &discordgo.MessageEmbed{
		Color:       0xED4245,
		Title:       "bio",
		Description: fmt.Sprintf("Usage: `-bio %s <username>`. To add a provider, implement its fetcher and register it in `bioProviders`.", strings.Join(names, "|")),
	}
}

// bioMessage usage: -bio {github|roblox|reddit|kick} <username>
func bioMessageCommandHandler(s *discordgo.Session, m *discordgo.MessageCreate, args []string) error {
	if len(args) < 2 {
		_, err := s.ChannelMessageSendEmbed(m.ChannelID, bioUsageEmbed())
		return err
	}
	provider, ok := bioProviderByName(strings.ToLower(args[0]))
	if !ok {
		_, err := s.ChannelMessageSendEmbed(m.ChannelID, bioUsageEmbed())
		return err
	}
	username := strings.Join(args[1:], "")
	embed, err := provider.build(username, nil)
	if err != nil {
		return err
	}
	_, err = s.ChannelMessageSendEmbed(m.ChannelID, embed)
	return err
}

// bioSlash reads the subcommand (<provider>) and its "username" option.
func bioSlashCommandHandler(s *discordgo.Session, i *discordgo.InteractionCreate, cfg *config.Config) error {
	subs := i.ApplicationCommandData().Options
	if len(subs) == 0 {
		return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{Embeds: []*discordgo.MessageEmbed{bioUsageEmbed()}},
		})
	}
	sub := subs[0]
	opts := OptionMap(sub.Options)
	username := strings.TrimSpace(OptString(opts, "username"))
	if username == "" {
		return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{Content: fmt.Sprintf("Usage: `/bio %s username:<name>`", sub.Name)},
		})
	}

	provider, ok := bioProviderByName(sub.Name)
	if !ok {
		return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{Embeds: []*discordgo.MessageEmbed{bioUsageEmbed()}},
		})
	}

	embed, err := provider.build(username, cfg)
	if err != nil {
		return err
	}
	return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{Embeds: []*discordgo.MessageEmbed{embed}},
	})
}
