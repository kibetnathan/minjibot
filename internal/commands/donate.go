package commands

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

// DonateURL is where users can support the bot.
const DonateURL = "https://buymeacoffee.com/kruegenn"

// WebsiteURL is the bot's public dashboard/website.
const WebsiteURL = "https://minji-bot.netlify.app/"

func donateMessageCommandHandler(s *discordgo.Session, m *discordgo.MessageCreate, _ []string) error {
	_, err := s.ChannelMessageSendComplex(m.ChannelID, &discordgo.MessageSend{
		Embeds:     []*discordgo.MessageEmbed{donateEmbed()},
		Components: donateComponents(),
	})
	return err
}

func donateSlashCommandHandler(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds:     []*discordgo.MessageEmbed{donateEmbed()},
			Components: donateComponents(),
		},
	})
}

func donateEmbed() *discordgo.MessageEmbed {
	return &discordgo.MessageEmbed{
		Color:       0xFFDD00,
		Title:       "Support MinjiBot",
		Description: fmt.Sprintf("If MinjiBot has made your server better, consider buying me a coffee to keep it running. Every bit helps!\n\n[Buy me a coffee](%s)", DonateURL),
	}
}

func donateComponents() []discordgo.MessageComponent {
	return []discordgo.MessageComponent{
		discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				discordgo.Button{
					Label: "Buy me a coffee",
					Style: discordgo.LinkButton,
					URL:   DonateURL,
				},
			},
		},
	}
}
