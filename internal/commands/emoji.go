package commands

import (
	"encoding/base64"
	"fmt"
	"regexp"
	"strings"

	"github.com/bwmarrin/discordgo"
)

const EmojiNameMaxLen = 32

var customEmojiRe = regexp.MustCompile(`<(a?):([a-zA-Z0-9_]{2,32}):(\d+)>`)

// EmojiTarget resolves an emoji from a command argument, or — if none was
// given — from the first custom emoji in the replied-to message's content.
func EmojiTarget(m *discordgo.MessageCreate, args []string) (name, id string, animated bool, ok bool) {
	if len(args) > 0 {
		if n, i, a, parsed := ParseEmoji(args[0]); parsed {
			return n, i, a, true
		}
		return "", "", false, false
	}
	if m.ReferencedMessage != nil {
		if sub := customEmojiRe.FindStringSubmatch(m.ReferencedMessage.Content); sub != nil {
			return sub[2], sub[3], sub[1] == "a", true
		}
	}
	return "", "", false, false
}

// ParseEmoji extracts (name, id, animated) from an emoji string
// like <:wave:1234>, <a:wave:1234>, :wave:1234, or a plain id.
func ParseEmoji(arg string) (name, id string, animated bool, ok bool) {
	arg = strings.TrimSpace(arg)

	animated = strings.HasPrefix(arg, "<a:")
	arg = strings.TrimPrefix(arg, "<a:")
	arg = strings.TrimPrefix(arg, "<:")
	arg = strings.TrimSuffix(arg, ">")

	parts := strings.Split(arg, ":")
	switch {
	case len(parts) == 2:
		name, id = parts[0], parts[1]
	case len(parts) == 3 && parts[0] == "" && parts[1] != "":
		// :wave:1234
		name, id = parts[1], parts[2]
	case len(parts) == 1 && parts[0] != "":
		id = parts[0]
	}

	id = strings.TrimSpace(id)
	name = strings.TrimSpace(name)
	ok = id != "" && DigitsOnly(id)
	return
}

func DigitsOnly(s string) bool {
	if s == "" {
		return false
	}
	for _, r := range s {
		if r < '0' || r > '9' {
			return false
		}
	}
	return true
}

// EmojiImageURL returns the CDN URL for an emoji at the given size.
func EmojiImageURL(id string, animated bool, size int) string {
	ext := "png"
	if animated {
		ext = "gif"
	}
	return fmt.Sprintf("https://cdn.discordapp.com/emojis/%s.%s?size=%d", id, ext, size)
}

// Base64DataURI encodes downloaded bytes as a Discord-ready data URI.
func Base64DataURI(md *MediaData) string {
	ext := md.Ext
	switch ext {
	case "jpg":
		ext = "jpeg"
	case "webm", "mp4", "bin":
		ext = "png"
	case "":
		ext = "png"
	}
	return "data:image/" + ext + ";base64," + base64.StdEncoding.EncodeToString(md.Data)
}

func SanitizeEmojiName(name string) string {
	var b strings.Builder
	for _, r := range strings.ToLower(name) {
		switch {
		case r >= 'a' && r <= 'z', r >= '0' && r <= '9', r == '_':
			b.WriteRune(r)
		}
		if b.Len() >= EmojiNameMaxLen {
			break
		}
	}
	return b.String()
}

func emojiCreate(s *discordgo.Session, guildID, name, imageDataURI string) (*discordgo.Emoji, error) {
	return s.GuildEmojiCreate(guildID, &discordgo.EmojiParams{Name: name, Image: imageDataURI})
}

func FilepathExt(filename string) string {
	idx := strings.LastIndex(filename, ".")
	if idx < 0 {
		return ""
	}
	return filename[idx:]
}

func AttachmentEmojiName(att *discordgo.MessageAttachment) string {
	return SanitizeEmojiName(strings.TrimSuffix(att.Filename, FilepathExt(att.Filename)))
}

// OptionMap converts slash options into a map keyed by name.
func OptionMap(opts []*discordgo.ApplicationCommandInteractionDataOption) map[string]*discordgo.ApplicationCommandInteractionDataOption {
	out := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(opts))
	for _, o := range opts {
		out[o.Name] = o
	}
	return out
}

func OptString(opts map[string]*discordgo.ApplicationCommandInteractionDataOption, name string) string {
	if o, ok := opts[name]; ok && o != nil {
		// String options carry the text directly; channel options carry their
		// ID as a plain string. Assert instead of discrodgo's StringValue(),
		// which panics on any non-string option type.
		if s, ok := o.Value.(string); ok {
			return s
		}
	}
	return ""
}

func emojiMessageCommandHandler(s *discordgo.Session, m *discordgo.MessageCreate, args []string) error {
	if len(args) == 0 {
		_, err := s.ChannelMessageSend(m.ChannelID, "Usage: `-emoji add <name> [url]`, `-emoji add many name1=url1 name2=url2`, `-emoji enlarge <emoji>`, `-emoji list`, `-emoji remove <emoji>`, `-emoji steal <emoji>`")
		return err
	}

	switch args[0] {
	case "add":
		return emojiAdd(s, m, args[1:])
	case "addmany", "many":
		return emojiAddMany(s, m, args[1:])
	case "enlarge":
		return emojiEnlarge(s, m, args[1:])
	case "list":
		return emojiList(s, m)
	case "remove", "delete":
		return emojiRemove(s, m, args[1:])
	case "steal":
		return emojiSteal(s, m, args[1:])
	default:
		_, err := s.ChannelMessageSend(m.ChannelID, "Unknown emoji subcommand. Use `add`, `many`, `enlarge`, `list`, `remove`, or `steal`.")
		return err
	}
}

func emojiAdd(s *discordgo.Session, m *discordgo.MessageCreate, args []string) error {
	guildID := m.GuildID
	if guildID == "" {
		return fmt.Errorf("emoji commands can only be used in a server")
	}
	if len(args) == 0 {
		_, err := s.ChannelMessageSend(m.ChannelID, "Usage: `-emoji add <name> [url]` (attach the emoji image or pass a URL)")
		return err
	}

	name := SanitizeEmojiName(args[0])
	if name == "" {
		return fmt.Errorf("invalid emoji name: %q", args[0])
	}

	url := ""
	if len(args) > 1 {
		url = args[1]
	}

	md, err := resolveMedia(s, m, url)
	if err != nil {
		return err
	}

	emoji, err := emojiCreate(s, guildID, name, Base64DataURI(md))
	if err != nil {
		return err
	}

	_, err = s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Added emoji %s", emoji.MessageFormat()))
	return err
}

func emojiAddMany(s *discordgo.Session, m *discordgo.MessageCreate, args []string) error {
	guildID := m.GuildID
	if guildID == "" {
		return fmt.Errorf("emoji commands can only be used in a server")
	}

	emojis := make([]*discordgo.Emoji, 0, len(args))

	// name=url pairs
	for _, arg := range args {
		parts := strings.SplitN(arg, "=", 2)
		if len(parts) != 2 {
			continue
		}
		name := SanitizeEmojiName(parts[0])
		if name == "" {
			continue
		}
		md, err := fetchURL(parts[1])
		if err != nil {
			continue
		}
		emoji, err := emojiCreate(s, guildID, name, Base64DataURI(md))
		if err != nil {
			continue
		}
		emojis = append(emojis, emoji)
	}

	// Any attachments on the message
	for _, att := range m.Attachments {
		md, err := fetchURL(att.URL)
		if err != nil {
			continue
		}
		name := AttachmentEmojiName(att)
		if name == "" {
			continue
		}
		emoji, err := emojiCreate(s, guildID, name, Base64DataURI(md))
		if err != nil {
			continue
		}
		emojis = append(emojis, emoji)
	}

	if len(emojis) == 0 {
		return fmt.Errorf("no emojis were added — use `-emoji add many name1=url1 name2=url2` or attach images")
	}

	formatted := make([]string, 0, len(emojis))
	for _, e := range emojis {
		formatted = append(formatted, e.Name)
	}

	_, err := s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Added %d emojis: %s", len(emojis), strings.Join(formatted, " ")))
	return err
}

func emojiEnlarge(s *discordgo.Session, m *discordgo.MessageCreate, args []string) error {
	_, id, animated, ok := EmojiTarget(m, args)
	if !ok {
		_, err := s.ChannelMessageSend(m.ChannelID, "Usage: `-emoji enlarge <emoji>` — or reply to a message containing an emoji")
		return err
	}

	_, err := s.ChannelMessageSend(m.ChannelID, EmojiImageURL(id, animated, 256))
	return err
}

func emojiList(s *discordgo.Session, m *discordgo.MessageCreate) error {
	guildID := m.GuildID
	if guildID == "" {
		return fmt.Errorf("emoji commands can only be used in a server")
	}

	pages, err := buildEmojiListPages(s, guildID)
	if err != nil {
		return err
	}
	return sendPagedEmbeds(s, m.ChannelID, m.Author.ID, pages)
}

// buildEmojiListPages returns paged embeds describing the server's emojis,
// chunked so each page stays within Discord's embed size limit.
func buildEmojiListPages(s *discordgo.Session, guildID string) ([]*discordgo.MessageEmbed, error) {
	emojis, err := s.GuildEmojis(guildID)
	if err != nil {
		return nil, err
	}
	if len(emojis) == 0 {
		return []*discordgo.MessageEmbed{{
			Color:       0x5865F2,
			Title:       "Server Emojis",
			Description: "This server has no custom emojis.",
		}}, nil
	}

	lines := make([]string, 0, len(emojis))
	for _, e := range emojis {
		lines = append(lines, fmt.Sprintf("%s `:%s:`\n", e.MessageFormat(), e.Name))
	}
	return pageEmbeds(fmt.Sprintf("Server Emojis (%d)", len(emojis)), lines, maxEmbedTextBudget, 0x5865F2), nil
}

func ChunkEmojiList(emojis []*discordgo.Emoji, maxLen int) []*discordgo.MessageEmbedField {
	var fields []*discordgo.MessageEmbedField
	var cur string
	for _, e := range emojis {
		line := fmt.Sprintf("%s `:%s:`\n", e.MessageFormat(), e.Name)
		if cur != "" && len(cur+line) > maxLen {
			fields = append(fields, &discordgo.MessageEmbedField{Name: "—", Value: cur})
			cur = ""
		}
		cur += line
	}
	if cur != "" {
		fields = append(fields, &discordgo.MessageEmbedField{Name: "—", Value: cur})
	}
	return fields
}

func emojiRemove(s *discordgo.Session, m *discordgo.MessageCreate, args []string) error {
	guildID := m.GuildID
	if guildID == "" {
		return fmt.Errorf("emoji commands can only be used in a server")
	}
	if len(args) == 0 {
		_, err := s.ChannelMessageSend(m.ChannelID, "Usage: `-emoji remove <emoji>` (or an emoji ID)")
		return err
	}

	_, id, _, ok := ParseEmoji(args[0])
	if !ok {
		return fmt.Errorf("couldn't parse an emoji from %q", args[0])
	}

	if err := s.GuildEmojiDelete(guildID, id); err != nil {
		return err
	}

	_, err := s.ChannelMessageSend(m.ChannelID, "Removed emoji.")
	return err
}

func emojiSteal(s *discordgo.Session, m *discordgo.MessageCreate, args []string) error {
	guildID := m.GuildID
	if guildID == "" {
		return fmt.Errorf("emoji commands can only be used in a server")
	}

	name, id, animated, ok := EmojiTarget(m, args)
	if !ok {
		_, err := s.ChannelMessageSend(m.ChannelID, "Usage: `-emoji steal <emoji>` — or reply to a message containing an emoji")
		return err
	}

	md, err := fetchURL(EmojiImageURL(id, animated, 128))
	if err != nil {
		return err
	}

	emoji, err := emojiCreate(s, guildID, name, Base64DataURI(md))
	if err != nil {
		return err
	}

	_, err = s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Stole %s", emoji.MessageFormat()))
	return err
}

func emojiSlashCommandHandler(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	data := i.ApplicationCommandData()
	if len(data.Options) == 0 {
		return fmt.Errorf("missing emoji subcommand")
	}
	sub := data.Options[0]
	opts := OptionMap(sub.Options)

	guildID := i.GuildID
	if guildID == "" {
		return fmt.Errorf("emoji commands can only be used in a server")
	}

	var content string
	var embeds []*discordgo.MessageEmbed

	switch sub.Name {
	case "add":
		name := SanitizeEmojiName(OptString(opts, "name"))
		if name == "" {
			return fmt.Errorf("invalid emoji name")
		}
		md, err := fetchURL(OptString(opts, "url"))
		if err != nil {
			return err
		}
		emoji, err := emojiCreate(s, guildID, name, Base64DataURI(md))
		if err != nil {
			return err
		}
		content = "Added " + emoji.MessageFormat()

	case "enlarge":
		_, id, animated, ok := ParseEmoji(OptString(opts, "emoji"))
		if !ok {
			return fmt.Errorf("couldn't parse an emoji")
		}
		content = EmojiImageURL(id, animated, 256)

	case "list":
		pages, err := buildEmojiListPages(s, guildID)
		if err != nil {
			return err
		}
		return respondPaged(s, i, pages)

	case "remove":
		_, id, _, ok := ParseEmoji(OptString(opts, "emoji"))
		if !ok {
			return fmt.Errorf("couldn't parse an emoji")
		}
		if err := s.GuildEmojiDelete(guildID, id); err != nil {
			return err
		}
		content = "Removed emoji."

	case "steal":
		name, id, animated, ok := ParseEmoji(OptString(opts, "emoji"))
		if !ok {
			return fmt.Errorf("couldn't parse an emoji")
		}
		md, err := fetchURL(EmojiImageURL(id, animated, 128))
		if err != nil {
			return err
		}
		emoji, err := emojiCreate(s, guildID, name, Base64DataURI(md))
		if err != nil {
			return err
		}
		content = "Stole " + emoji.MessageFormat()

	default:
		return fmt.Errorf("unknown emoji subcommand: %s", sub.Name)
	}

	return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{Content: content, Embeds: embeds},
	})
}
