package commands

import (
	"fmt"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

func testMessageCommandHandler(s *discordgo.Session, m *discordgo.MessageCreate, args []string) error {
	start := time.Now()
	msg, err := s.ChannelMessageSend(m.ChannelID, "**Testing... 🏓**")
	if err != nil {
		return err
	}
	content := fmt.Sprintf("**Pong! 🏓** The bot is working correctly.\n• **Round-trip:** `%dms`", time.Since(start).Milliseconds())
	if len(args) > 0 {
		content += fmt.Sprintf("\n• **You said:** `%s`", strings.Join(args, " "))
	}
	_, err = s.ChannelMessageEdit(m.ChannelID, msg.ID, content)
	return err
}

func testSlashCommandHandler(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	opts := OptionMap(i.ApplicationCommandData().Options)
	start := time.Now()

	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
	})
	if err != nil {
		return err
	}

	content := fmt.Sprintf("**Pong! 🏓** The bot is working correctly.\n• **Round-trip:** `%dms`", time.Since(start).Milliseconds())
	if t := strings.TrimSpace(OptString(opts, "text")); t != "" {
		content += fmt.Sprintf("\n• **You said:** `%s`", t)
	}

	_, err = s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{Content: &content})
	return err
}
