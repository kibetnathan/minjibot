package commands

import (
	"bytes"
	"encoding/json"
	"fmt"
	"mime/multipart"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func stickerMessageCommandHandler(s *discordgo.Session, m *discordgo.MessageCreate, args []string) error {
	if len(args) == 0 {
		_, err := s.ChannelMessageSend(m.ChannelID, "Usage: `!sticker add <name> [url]`, `!sticker remove <sticker-id|message-link>`")
		return err
	}

	switch args[0] {
	case "add":
		return stickerAddMessage(s, m, args[1:])
	case "remove", "delete":
		return stickerRemoveMessage(s, m, args[1:])
	default:
		_, err := s.ChannelMessageSend(m.ChannelID, "Unknown sticker subcommand. Use `add` or `remove`.")
		return err
	}
}

func stickerAddMessage(s *discordgo.Session, m *discordgo.MessageCreate, args []string) error {
	if len(args) == 0 {
		_, err := s.ChannelMessageSend(m.ChannelID, "Usage: `!sticker add <name> [url]` (attach the sticker image or pass a URL)")
		return err
	}
	name := strings.TrimSpace(args[0])
	if name == "" {
		return fmt.Errorf("invalid sticker name")
	}

	url := ""
	if len(args) > 1 {
		url = args[1]
	}
	md, err := resolveMedia(s, m, url)
	if err != nil {
		return err
	}

	sticker, err := guildStickerCreate(s, m.GuildID, name, "Added by "+m.Author.Username, name, md.Data)
	if err != nil {
		return err
	}

	_, err = s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Added sticker `%s`", sticker.Name))
	return err
}

func stickerRemoveMessage(s *discordgo.Session, m *discordgo.MessageCreate, args []string) error {
	if len(args) == 0 {
		_, err := s.ChannelMessageSend(m.ChannelID, "Usage: `!sticker remove <sticker-id|message-link>`")
		return err
	}
	id := strings.TrimSpace(args[0])

	// Accept a message link that contains a sticker.
	if strings.Contains(id, "discord.com/channels/") {
		_, channelID, messageID, ok := parseMessageLink(id)
		if !ok {
			return fmt.Errorf("couldn't parse message link")
		}
		msg, err := s.ChannelMessage(channelID, messageID)
		if err != nil {
			return err
		}
		if len(msg.StickerItems) == 0 {
			return fmt.Errorf("that message doesn't contain a sticker")
		}
		id = msg.StickerItems[0].ID
	}

	if err := guildStickerDelete(s, m.GuildID, id); err != nil {
		return err
	}

	_, err := s.ChannelMessageSend(m.ChannelID, "Removed sticker.")
	return err
}

func stickerSlashCommandHandler(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	data := i.ApplicationCommandData()
	if len(data.Options) == 0 {
		return fmt.Errorf("missing sticker subcommand")
	}
	sub := data.Options[0]
	opts := optionMap(sub.Options)

	guildID := i.GuildID
	if guildID == "" {
		return fmt.Errorf("sticker commands can only be used in a server")
	}

	var content string

	switch sub.Name {
	case "add":
		name := strings.TrimSpace(optString(opts, "name"))
		if name == "" {
			return fmt.Errorf("invalid sticker name")
		}
		md, err := fetchURL(optString(opts, "url"))
		if err != nil {
			return err
		}
		sticker, err := guildStickerCreate(s, guildID, name, "Added via slash command", name, md.Data)
		if err != nil {
			return err
		}
		content = "Added sticker `" + sticker.Name + "`"

	case "remove":
		id := strings.TrimSpace(optString(opts, "sticker_id"))
		if id == "" {
			return fmt.Errorf("provide a sticker ID or message link")
		}
		if strings.Contains(id, "discord.com/channels/") {
			_, channelID, messageID, ok := parseMessageLink(id)
			if !ok {
				return fmt.Errorf("couldn't parse message link")
			}
			msg, err := s.ChannelMessage(channelID, messageID)
			if err != nil {
				return err
			}
			if len(msg.StickerItems) == 0 {
				return fmt.Errorf("that message doesn't contain a sticker")
			}
			id = msg.StickerItems[0].ID
		}
		if err := guildStickerDelete(s, guildID, id); err != nil {
			return err
		}
		content = "Removed sticker."

	default:
		return fmt.Errorf("unknown sticker subcommand: %s", sub.Name)
	}

	return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{Content: content},
	})
}

// guildStickerCreate uploads a sticker via multipart/form-data, since the
// discordgo version in use lacks a wrapper for the sticker endpoints.
func guildStickerCreate(s *discordgo.Session, guildID, name, desc, tags string, data []byte) (*discordgo.Sticker, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	for k, v := range map[string]string{
		"name":        name,
		"description": desc,
		"tags":        tags,
	} {
		if err := writer.WriteField(k, v); err != nil {
			return nil, err
		}
	}

	part, err := writer.CreateFormFile("file", "sticker.png")
	if err != nil {
		return nil, err
	}
	if _, err := part.Write(data); err != nil {
		return nil, err
	}
	if err := writer.Close(); err != nil {
		return nil, err
	}

	url := discordgo.EndpointAPI + "guilds/" + guildID + "/stickers"
	bodyBytes := body.Bytes()
	contentType := writer.FormDataContentType()

	respBytes, err := s.RequestRaw("POST", url, contentType, bodyBytes, guildID+":stickers", 0)
	if err != nil {
		return nil, err
	}

	var sticker discordgo.Sticker
	if err := json.Unmarshal(respBytes, &sticker); err != nil {
		return nil, err
	}
	return &sticker, nil
}

func guildStickerDelete(s *discordgo.Session, guildID, stickerID string) error {
	if !digitsOnly(stickerID) {
		return fmt.Errorf("invalid sticker ID: %q", stickerID)
	}
	url := discordgo.EndpointAPI + "guilds/" + guildID + "/stickers/" + stickerID
	_, err := s.RequestWithBucketID("DELETE", url, nil, guildID+":stickers")
	return err
}
