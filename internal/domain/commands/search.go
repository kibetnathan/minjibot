package commands

import (
	"fmt"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

const (
	searchBatchSize      = 100
	searchDefaultLimit   = 200
	searchMaxLimit       = 1000
	searchMaxResults     = 5
	searchMaxFieldLength = 1024
)

func searchMessageCommandHandler(s *discordgo.Session, m *discordgo.MessageCreate, args []string) error {
	query, limit, err := parseSearchArgs(s, m.ChannelID, args)
	if err != nil {
		return err
	}

	matches, err := searchChatHistory(s, m.GuildID, m.ChannelID, query, limit)
	if err != nil {
		return err
	}

	embed := buildSearchEmbed(query, limit, matches)
	_, err = s.ChannelMessageSendEmbed(m.ChannelID, embed)
	return err
}

func searchSlashCommandHandler(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	query := ""
	limit := searchDefaultLimit
	for _, opt := range i.ApplicationCommandData().Options {
		switch opt.Name {
		case "query":
			query = opt.StringValue()
		case "messages":
			limit = int(opt.IntValue())
		}
	}

	query = strings.TrimSpace(query)
	if query == "" {
		return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Usage: `/search query:<text>`",
			},
		})
	}
	if limit < 1 || limit > searchMaxLimit {
		limit = searchDefaultLimit
	}

	matches, err := searchChatHistory(s, i.GuildID, i.ChannelID, query, limit)
	if err != nil {
		return err
	}

	embed := buildSearchEmbed(query, limit, matches)
	return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{embed},
		},
	})
}

// parseSearchArgs handles `!search <query>` and optional trailing
// `messages:<n>` or `<n>` to control how far back to look.
func parseSearchArgs(s *discordgo.Session, channelID string, args []string) (string, int, error) {
	if len(args) == 0 {
		if _, err := s.ChannelMessageSend(channelID, "Usage: `!search <query>`"); err != nil {
			return "", 0, err
		}
		return "", 0, fmt.Errorf("search requires a query")
	}

	query := ""
	limit := searchDefaultLimit
	rest := make([]string, 0, len(args))
	for _, arg := range args {
		lower := strings.ToLower(arg)
		if strings.HasPrefix(lower, "messages:") {
			if n := strings.TrimPrefix(lower, "messages:"); n != "" {
				fmt.Sscanf(n, "%d", &limit)
			}
			continue
		}
		if n := 0; len(arg) > 0 {
			if _, err := fmt.Sscanf(arg, "%d", &n); err == nil {
				limit = n
				continue
			}
		}
		rest = append(rest, arg)
	}

	query = strings.TrimSpace(strings.Join(rest, " "))
	if limit < 1 || limit > searchMaxLimit {
		limit = searchDefaultLimit
	}
	return query, limit, nil
}

func searchChatHistory(s *discordgo.Session, guildID, channelID, query string, limit int) ([]*discordgo.Message, error) {
	need := searchBatchSize
	if limit < searchBatchSize {
		need = limit
	}

	var all []*discordgo.Message
	beforeID := ""
	for len(all) < limit {
		batch := need
		remaining := limit - len(all)
		if remaining < batch {
			batch = remaining
		}

		msgs, err := s.ChannelMessages(channelID, batch, beforeID, "", "")
		if err != nil {
			return nil, err
		}
		if len(msgs) == 0 {
			break
		}

		all = append(all, msgs...)
		beforeID = msgs[len(msgs)-1].ID
	}

	lowerQuery := strings.ToLower(query)
	matches := make([]*discordgo.Message, 0, len(all))
	for _, msg := range all {
		if msg.Content != "" && strings.Contains(strings.ToLower(msg.Content), lowerQuery) {
			matches = append(matches, msg)
		}
	}
	return matches, nil
}

func buildSearchEmbed(query string, limit int, matches []*discordgo.Message) *discordgo.MessageEmbed {
	embed := &discordgo.MessageEmbed{
		Color:      0x5865F2,
		Title:      fmt.Sprintf("Search results for %q", query),
		Description: fmt.Sprintf("Searched the last %d messages in this channel.", limit),
	}

	if len(matches) == 0 {
		embed.Description += "\n\nNo matches found."
		return embed
	}

	embed.Footer = &discordgo.MessageEmbedFooter{
		Text: fmt.Sprintf("%d match(es) found", len(matches)),
	}

	count := searchMaxResults
	if len(matches) < count {
		count = len(matches)
	}
	for _, msg := range matches[:count] {
		author := msg.Author.Username
		if msg.Author.Bot {
			author += " [bot]"
		}
		content := msg.Content
		if len(content) > searchMaxFieldLength {
			content = content[:searchMaxFieldLength-3] + "..."
		}

		value := content
		link := fmt.Sprintf("https://discord.com/channels/%s/%s/%s", msg.GuildID, msg.ChannelID, msg.ID)
		if msg.GuildID == "" {
			link = fmt.Sprintf("https://discord.com/channels/@me/%s/%s", msg.ChannelID, msg.ID)
		}
		if value != "" {
			value += "\n"
		}
		value += link

		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
			Name:  fmt.Sprintf("%s • %s", author, msg.Timestamp.Format(time.RFC3339)),
			Value: value,
		})
	}

	return embed
}