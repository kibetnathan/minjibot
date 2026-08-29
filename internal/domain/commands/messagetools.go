package commands

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/bwmarrin/discordgo"
)

var messageLinkRe = regexp.MustCompile(`discord(?:app)?\.com/channels/(\d+)/(\d+)/(\d+)`)

// parseMessageLink extracts (guildID, channelID, messageID) from a Discord
// message link. The first component may be @me for DMs.
func parseMessageLink(link string) (guildID, channelID, messageID string, ok bool) {
	m := messageLinkRe.FindStringSubmatch(link)
	if len(m) != 4 {
		return "", "", "", false
	}
	return m[1], m[2], m[3], true
}

// resolveMessageTarget finds a message from a message link, a bare message ID
// in the current channel, or the referenced (replied-to) message.
func resolveMessageTarget(s *discordgo.Session, m *discordgo.MessageCreate, args []string) (*discordgo.Message, error) {
	if len(args) > 0 {
		if strings.Contains(args[0], "discord.com/channels/") {
			_, channelID, messageID, ok := parseMessageLink(args[0])
			if !ok {
				return nil, fmt.Errorf("couldn't parse message link: %q", args[0])
			}
			msg, err := s.ChannelMessage(channelID, messageID)
			if err != nil {
				return nil, err
			}
			return msg, nil
		}

		if digitsOnly(args[0]) {
			msg, err := s.ChannelMessage(m.ChannelID, args[0])
			if err != nil {
				return nil, err
			}
			return msg, nil
		}
	}

	if m.ReferencedMessage != nil {
		return m.ReferencedMessage, nil
	}

	return nil, fmt.Errorf("no message provided — reply to a message, pass a message link, or pass a message ID")
}

func pinMessageCommandHandler(s *discordgo.Session, m *discordgo.MessageCreate, args []string) error {
	msg, err := resolveMessageTarget(s, m, args)
	if err != nil {
		return err
	}

	if err := s.ChannelMessagePin(msg.ChannelID, msg.ID); err != nil {
		return err
	}

	_, err = s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Pinned [message](%s).", messageToLink(msg)))
	return err
}

func unpinMessageCommandHandler(s *discordgo.Session, m *discordgo.MessageCreate, args []string) error {
	msg, err := resolveMessageTarget(s, m, args)
	if err != nil {
		return err
	}

	if err := s.ChannelMessageUnpin(msg.ChannelID, msg.ID); err != nil {
		return err
	}

	_, err = s.ChannelMessageSend(m.ChannelID, "Unpinned that message.")
	return err
}

func pinSlashCommandHandler(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	return pinUnpinSlashCommandHandler(s, i, true)
}

func unpinSlashCommandHandler(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	return pinUnpinSlashCommandHandler(s, i, false)
}

func pinUnpinSlashCommandHandler(s *discordgo.Session, i *discordgo.InteractionCreate, pin bool) error {
	opts := optionMap(i.ApplicationCommandData().Options)
	channelID := optString(opts, "channel")
	if channelID == "" {
		channelID = i.ChannelID
	}
	messageID := optString(opts, "message")
	if messageID == "" {
		return fmt.Errorf("provide a message ID or a message link")
	}

	// Allow pasting a full message link into the message option.
	if strings.Contains(messageID, "discord.com/channels/") {
		_, ch, msg, ok := parseMessageLink(messageID)
		if !ok {
			return fmt.Errorf("couldn't parse message link")
		}
		channelID, messageID = ch, msg
	}

	var err error
	verb := "pinned"
	if pin {
		err = s.ChannelMessagePin(channelID, messageID)
	} else {
		err = s.ChannelMessageUnpin(channelID, messageID)
		verb = "unpinned"
	}
	if err != nil {
		return err
	}

	return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{Content: fmt.Sprintf("Message %s.", verb)},
	})
}

func quoteMessageCommandHandler(s *discordgo.Session, m *discordgo.MessageCreate, args []string) error {
	msg, err := resolveMessageTarget(s, m, args)
	if err != nil {
		return err
	}

	_, err = s.ChannelMessageSendEmbed(m.ChannelID, quoteEmbed(msg))
	return err
}

func quoteSlashCommandHandler(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	opts := optionMap(i.ApplicationCommandData().Options)
	channelID := optString(opts, "channel")
	if channelID == "" {
		channelID = i.ChannelID
	}
	messageID := optString(opts, "message")
	if messageID == "" {
		return fmt.Errorf("provide a message ID or a message link")
	}

	if strings.Contains(messageID, "discord.com/channels/") {
		_, ch, msg, ok := parseMessageLink(messageID)
		if !ok {
			return fmt.Errorf("couldn't parse message link")
		}
		channelID, messageID = ch, msg
	}

	msg, err := s.ChannelMessage(channelID, messageID)
	if err != nil {
		return err
	}

	return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{Embeds: []*discordgo.MessageEmbed{quoteEmbed(msg)}},
	})
}

func messageToLink(msg *discordgo.Message) string {
	guild := msg.GuildID
	if guild == "" {
		guild = "@me"
	}
	return fmt.Sprintf("https://discord.com/channels/%s/%s/%s", guild, msg.ChannelID, msg.ID)
}

func quoteEmbed(msg *discordgo.Message) *discordgo.MessageEmbed {
	author := "Unknown"
	if msg.Author != nil {
		author = msg.Author.Username
	}

	embed := &discordgo.MessageEmbed{
		Color:       0x5865F2,
		Author:      &discordgo.MessageEmbedAuthor{Name: author},
		Description: msg.Content,
		Timestamp:   msg.Timestamp.Format("2006-01-02T15:04:05-07:00"),
	}

	if msg.Author != nil && msg.Author.Avatar != "" {
		embed.Author.IconURL = "https://cdn.discordapp.com/avatars/" + msg.Author.ID + "/" + msg.Author.Avatar + ".png?size=64"
	}

	if att := firstAttachment(msg); att != nil {
		embed.Image = &discordgo.MessageEmbedImage{URL: att.URL}
	}

	if len(msg.Embeds) > 0 && msg.Embeds[0].Image != nil {
		embed.Image = &discordgo.MessageEmbedImage{URL: msg.Embeds[0].Image.URL}
	}

	embed.Footer = &discordgo.MessageEmbedFooter{Text: "Jump to message • " + messageToLink(msg)}
	return embed
}
