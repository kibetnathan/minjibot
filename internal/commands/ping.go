package commands

import (
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
)

func pingSlashCommandHandler(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	startTime := time.Now()

	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{Type: discordgo.InteractionResponseDeferredChannelMessageWithSource})
	if err != nil {
		return err
	}

	roundTrip := time.Since(startTime).Milliseconds()

	wsLatency := s.HeartbeatLatency().Milliseconds()
	content := fmt.Sprintf("**Pong!**\n"+
		"• **Round-trip Latency:** `%dms`\n"+
		"• **WebSocket Latency:** `%dms`", roundTrip, wsLatency)

	_, err = s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{Content: &content})

	return nil
}

func pingMessageCommandHandler(s *discordgo.Session, channelID string) error {
	startTime := time.Now()

	msg, err := s.ChannelMessageSend(channelID, "**Pong!**\n• **Pinging...**")
	if err != nil {
		return err
	}

	roundTrip := time.Since(startTime).Milliseconds()

	wsLatency := s.HeartbeatLatency().Milliseconds()
	content := fmt.Sprintf("**Pong!**\n"+
		"• **Round-trip Latency:** `%dms`\n"+
		"• **WebSocket Latency:** `%dms`", roundTrip, wsLatency)

	_, err = s.ChannelMessageEdit(channelID, msg.ID, content)

	return err
}
