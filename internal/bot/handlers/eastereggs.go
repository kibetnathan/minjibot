package handlers

import (
	"strings"

	"github.com/bwmarrin/discordgo"
)

// checkEasterEggs applies one-off easter egg reactions to certain messages.
// Returning true means the message was handled and the normal handler should
// stop. Add new easter eggs here in one place.
func checkEasterEggs(s *discordgo.Session, m *discordgo.MessageCreate) bool {
	return scatEasterEgg(s, m)
}

// scatEasterEgg deletes messages containing "scat" AND a mention of the owner
// (Kruegen / Nathan / @Kruegenn), then replies with 👀.
func scatEasterEgg(s *discordgo.Session, m *discordgo.MessageCreate) bool {
	content := strings.ToLower(m.Content)
	if strings.Contains(content, "scat") &&
		(strings.Contains(content, "kruegen") || strings.Contains(content, "nathan") || strings.Contains(content, "@kruegenn")) {
		_ = s.ChannelMessageDelete(m.ChannelID, m.ID)
		_, _ = s.ChannelMessageSend(m.ChannelID, "👀")
		return true
	}
	return false
}
