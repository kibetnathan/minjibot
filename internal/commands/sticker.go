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
		_, err := s.ChannelMessageSend(m.ChannelID, "Usage: `-sticker add <name> [url]`, `-sticker steal <sticker-id|message-link|sticker>`, `-sticker remove <sticker-id|message-link>`")
		return err
	}

	switch args[0] {
	case "add":
		return stickerAddMessage(s, m, args[1:])
	case "steal":
		return stickerStealMessage(s, m, args[1:])
	case "remove", "delete":
		return stickerRemoveMessage(s, m, args[1:])
	default:
		_, err := s.ChannelMessageSend(m.ChannelID, "Unknown sticker subcommand. Use `add`, `steal`, or `remove`.")
		return err
	}
}

// resolveStickerSource returns a sticker ID and CDN URL from a raw argument:
// a message link containing a sticker, a bare sticker ID, or a sticker CDN URL.
func resolveStickerSource(s *discordgo.Session, raw string) (string, string, error) {
	raw = strings.TrimSpace(raw)

	// Message link → look up the message and grab its sticker.
	if strings.Contains(raw, "discord.com/channels/") {
		_, channelID, messageID, ok := ParseMessageLink(raw)
		if !ok {
			return "", "", fmt.Errorf("couldn't parse message link")
		}
		msg, err := s.ChannelMessage(channelID, messageID)
		if err != nil {
			return "", "", err
		}
		if len(msg.StickerItems) == 0 {
			return "", "", fmt.Errorf("that message doesn't contain a sticker")
		}
		item := msg.StickerItems[0]
		return item.ID, StickerImageURL(item.ID, item.FormatType), nil
	}

	// Bare sticker ID. Look the sticker up so we use the right CDN extension for
	// its format — a GIF sticker lives at a .gif URL, not .png, so assuming PNG
	// here 404s for animated stickers.
	if DigitsOnly(raw) {
		st, err := fetchSticker(s, raw)
		if err != nil {
			return "", "", fmt.Errorf("couldn't look up sticker %s: %w", raw, err)
		}
		return raw, StickerImageURL(raw, st.FormatType), nil
	}

	// Sticker CDN URL.
	if strings.Contains(raw, "cdn.discordapp.com/stickers/") {
		return "", raw, nil
	}

	return "", "", fmt.Errorf("provide a message link, sticker ID, or sticker URL to steal")
}

func stickerStealMessage(s *discordgo.Session, m *discordgo.MessageCreate, args []string) error {
	if len(args) == 0 {
		_, err := s.ChannelMessageSend(m.ChannelID, "Usage: `-sticker steal <sticker-id|message-link|sticker>`")
		return err
	}
	if err := stealStickerToGuild(s, m.GuildID, args[0]); err != nil {
		return err
	}
	_, err := s.ChannelMessageSend(m.ChannelID, "Sticker stolen and added to this server!")
	return err
}

// stealStickerToGuild resolves a sticker source and uploads it to the guild.
func stealStickerToGuild(s *discordgo.Session, guildID, raw string) error {
	stickerID, url, err := resolveStickerSource(s, raw)
	if err != nil {
		return err
	}

	label := "sticker"
	if stickerID != "" {
		label = stickerID
	}

	md, err := fetchURL(url)
	if err != nil {
		return err
	}

	if _, err := guildStickerCreate(s, guildID, label, "Stolen via steal command", label, md); err != nil {
		return err
	}
	return nil
}

func stickerAddMessage(s *discordgo.Session, m *discordgo.MessageCreate, args []string) error {
	if len(args) == 0 {
		_, err := s.ChannelMessageSend(m.ChannelID, "Usage: `-sticker add <name> [url]` (attach the sticker image or pass a URL)")
		return err
	}
	name := strings.TrimSpace(args[0])
	if err := validateStickerName(name); err != nil {
		return err
	}

	url := ""
	if len(args) > 1 {
		url = args[1]
	}
	md, err := resolveMedia(s, m, url)
	if err != nil {
		return err
	}

	sticker, err := guildStickerCreate(s, m.GuildID, name, "Added by "+m.Author.Username, name, md)
	if err != nil {
		return err
	}

	_, err = s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Added sticker `%s`", sticker.Name))
	return err
}

func stickerRemoveMessage(s *discordgo.Session, m *discordgo.MessageCreate, args []string) error {
	if len(args) == 0 {
		_, err := s.ChannelMessageSend(m.ChannelID, "Usage: `-sticker remove <sticker-id|message-link>`")
		return err
	}
	id := strings.TrimSpace(args[0])

	// Accept a message link that contains a sticker.
	if strings.Contains(id, "discord.com/channels/") {
		_, channelID, messageID, ok := ParseMessageLink(id)
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
	opts := OptionMap(sub.Options)

	guildID := i.GuildID
	if guildID == "" {
		return fmt.Errorf("sticker commands can only be used in a server")
	}

	var content string

	switch sub.Name {
	case "add":
		name := strings.TrimSpace(OptString(opts, "name"))
		if err := validateStickerName(name); err != nil {
			return err
		}
		md, err := fetchURL(OptString(opts, "url"))
		if err != nil {
			return err
		}
		sticker, err := guildStickerCreate(s, guildID, name, "Added via slash command", name, md)
		if err != nil {
			return err
		}
		content = "Added sticker `" + sticker.Name + "`"

	case "remove":
		id := strings.TrimSpace(OptString(opts, "sticker_id"))
		if id == "" {
			return fmt.Errorf("provide a sticker ID or message link")
		}
		if strings.Contains(id, "discord.com/channels/") {
			_, channelID, messageID, ok := ParseMessageLink(id)
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

	case "steal":
		raw := strings.TrimSpace(OptString(opts, "sticker"))
		if raw == "" {
			return fmt.Errorf("provide a sticker ID, message link, or sticker URL to steal")
		}
		if err := stealStickerToGuild(s, guildID, raw); err != nil {
			return err
		}
		content = "Sticker stolen and added to this server!"

	default:
		return fmt.Errorf("unknown sticker subcommand: %s", sub.Name)
	}

	return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{Content: content},
	})
}

// validateStickerName enforces Discord's sticker name length (2–30 characters)
// up front so callers get a clear message instead of a raw 400 from the API.
func validateStickerName(name string) error {
	n := len([]rune(strings.TrimSpace(name)))
	if n < 2 || n > 30 {
		return fmt.Errorf("sticker name must be 2-30 characters (got %d)", n)
	}
	return nil
}

// stickerUploadName returns the multipart filename Discord expects for a
// sticker of the given image type. Discord infers the sticker's format from the
// upload's file extension, so a GIF sent as "sticker.png" is rejected — the
// filename has to match the bytes. Lottie stickers can't be created from a
// raster/animation upload, so they're refused with a clear message.
func stickerUploadName(ext string) (string, error) {
	switch strings.ToLower(ext) {
	case "png", "apng":
		return "sticker.png", nil
	case "gif":
		return "sticker.gif", nil
	case "json":
		return "", fmt.Errorf("Lottie (animated JSON) stickers can't be uploaded — only PNG, APNG, or GIF")
	default:
		return "", fmt.Errorf("unsupported sticker image type %q — Discord stickers must be PNG, APNG, or GIF", OrEmpty(ext))
	}
}

// fetchSticker looks up a sticker object by ID. discordgo v0.29.0 ships no
// Session wrapper for this endpoint, so we call it directly and decode the
// result to learn the sticker's format.
func fetchSticker(s *discordgo.Session, stickerID string) (*discordgo.Sticker, error) {
	if !DigitsOnly(stickerID) {
		return nil, fmt.Errorf("invalid sticker ID: %q", stickerID)
	}
	body, err := s.RequestWithBucketID("GET", discordgo.EndpointSticker(stickerID), nil, "stickers")
	if err != nil {
		return nil, err
	}
	var st discordgo.Sticker
	if err := json.Unmarshal(body, &st); err != nil {
		return nil, err
	}
	return &st, nil
}

// guildStickerCreate uploads a sticker via multipart/form-data, since the
// discordgo version in use lacks a wrapper for the sticker endpoints. The
// upload filename is derived from the media's real type so Discord accepts the
// format (see stickerUploadName).
func guildStickerCreate(s *discordgo.Session, guildID, name, desc, tags string, md *MediaData) (*discordgo.Sticker, error) {
	filename, err := stickerUploadName(md.Ext)
	if err != nil {
		return nil, err
	}

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

	part, err := writer.CreateFormFile("file", filename)
	if err != nil {
		return nil, err
	}
	if _, err := part.Write(md.Data); err != nil {
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
	if !DigitsOnly(stickerID) {
		return fmt.Errorf("invalid sticker ID: %q", stickerID)
	}
	url := discordgo.EndpointAPI + "guilds/" + guildID + "/stickers/" + stickerID
	_, err := s.RequestWithBucketID("DELETE", url, nil, guildID+":stickers")
	return err
}

// StickerImageURL returns the CDN image URL for a sticker, matching the
// extension to the sticker's format. PNG and APNG both serve from .png; GIF
// from .gif; Lottie from .json (which guildStickerCreate then rejects, rather
// than silently fetching a .png that doesn't exist for that sticker).
func StickerImageURL(id string, format discordgo.StickerFormat) string {
	switch format {
	case discordgo.StickerFormatTypeGIF:
		return fmt.Sprintf("https://cdn.discordapp.com/stickers/%s.gif", id)
	case discordgo.StickerFormatTypeLottie:
		return fmt.Sprintf("https://cdn.discordapp.com/stickers/%s.json", id)
	default:
		return fmt.Sprintf("https://cdn.discordapp.com/stickers/%s.png", id)
	}
}
