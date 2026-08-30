package commands

import (
	"strings"

	"github.com/bwmarrin/discordgo"
)

func echoMessageCommandHandler(s *discordgo.Session, channelID string, args []string) error {
	text := strings.Join(args, " ")
	if strings.TrimSpace(text) == "" {
		_, err := s.ChannelMessageSend(channelID, "Usage: `-echo <text>`")
		return err
	}

	_, err := s.ChannelMessageSend(channelID, text)
	return err
}

func echoSlashCommandHandler(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	data := i.ApplicationCommandData()

	var text string
	for _, opt := range data.Options {
		if opt.Name == "text" {
			text = opt.StringValue()
		}
	}

	if strings.TrimSpace(text) == "" {
		text = "Usage: `/echo text:<message>`"
	}

	return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: text,
		},
	})
}
