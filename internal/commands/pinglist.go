package commands

import (
	"fmt"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

const (
	pinglistBatchSize  = searchBatchSize
	pinglistMaxLimit   = 500
	PinglistMaxResults = 10
)

type PingTarget struct {
	Kind string // "user" or "role"
	ID   string
	Name string
}

func pinglistMessageCommandHandler(s *discordgo.Session, m *discordgo.MessageCreate, args []string) error {
	target, err := resolvePingTarget(s, m.GuildID, args)
	if err != nil {
		if _, serr := s.ChannelMessageSend(m.ChannelID, err.Error()); serr != nil {
			return serr
		}
		return nil
	}

	matches, err := searchMentions(s, m.GuildID, m.ChannelID, target, pinglistMaxLimit)
	if err != nil {
		return err
	}

	embed := BuildPinglistEmbed(target, matches)
	_, err = s.ChannelMessageSendEmbed(m.ChannelID, embed)
	return err
}

func pinglistSlashCommandHandler(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	target := &PingTarget{}
	for _, opt := range i.ApplicationCommandData().Options {
		switch opt.Name {
		case "user":
			target.Kind = "user"
			target.ID = SnowflakeFromValue(opt.Value)
		case "role":
			target.Kind = "role"
			target.ID = SnowflakeFromValue(opt.Value)
		}
	}

	if target.ID == "" {
		return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Usage: `/pinglist user:<user>` or `/pinglist role:<role>`",
			},
		})
	}
	target = resolvePingTargetByName(s, i.GuildID, target)

	matches, err := searchMentions(s, i.GuildID, i.ChannelID, target, pinglistMaxLimit)
	if err != nil {
		return err
	}

	embed := BuildPinglistEmbed(target, matches)
	return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{embed},
		},
	})
}

// resolvePingTarget figures out the target from a prefix command's args:
// a user mention, role mention, plain ID, or role name.
func resolvePingTarget(s *discordgo.Session, guildID string, args []string) (*PingTarget, error) {
	if len(args) == 0 {
		return nil, fmt.Errorf("Usage: `-pinglist @user` or `-pinglist @role`")
	}

	arg := args[0]

	switch {
	case strings.HasPrefix(arg, "<@&"):
		return &PingTarget{Kind: "role", ID: ParseMentionID(arg)}, nil
	case strings.HasPrefix(arg, "<@"):
		return &PingTarget{Kind: "user", ID: ParseMentionID(arg)}, nil
	}

	// Plain snowflake or role name: disambiguate against guild roles.
	roles, err := s.GuildRoles(guildID)
	if err == nil {
		for _, r := range roles {
			if r.ID == arg || strings.EqualFold(r.Name, arg) {
				return &PingTarget{Kind: "role", ID: r.ID, Name: r.Name}, nil
			}
		}
	}

	if id := ParseMentionID(arg); id != "" {
		return resolvePingTargetByName(s, guildID, &PingTarget{Kind: "user", ID: id}), nil
	}

	return nil, fmt.Errorf("Couldn't find a user or role matching %q", arg)
}

// resolvePingTargetByName fills in the display name for slash targets.
func resolvePingTargetByName(s *discordgo.Session, guildID string, t *PingTarget) *PingTarget {
	if t.Name != "" {
		return t
	}
	if t.Kind == "role" {
		if roles, err := s.GuildRoles(guildID); err == nil {
			for _, r := range roles {
				if r.ID == t.ID {
					t.Name = r.Name
					break
				}
			}
		}
	} else if t.Kind == "user" {
		if u, err := s.User(t.ID); err == nil && u != nil {
			t.Name = u.Username
		}
	}
	return t
}

func searchMentions(s *discordgo.Session, guildID, channelID string, target *PingTarget, limit int) ([]*discordgo.Message, error) {
	needles := []string{}
	switch target.Kind {
	case "user":
		needles = []string{"<@" + target.ID, "<@!" + target.ID}
	case "role":
		needles = []string{"<@&" + target.ID}
	}

	all, err := fetchChannelMessages(s, channelID, limit)
	if err != nil {
		return nil, err
	}

	matches := make([]*discordgo.Message, 0, len(all))
	for _, msg := range all {
		for _, n := range needles {
			if strings.Contains(msg.Content, n) {
				matches = append(matches, msg)
				break
			}
		}
	}
	return matches, nil
}

func fetchChannelMessages(s *discordgo.Session, channelID string, limit int) ([]*discordgo.Message, error) {
	need := pinglistBatchSize
	if limit < pinglistBatchSize {
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
	return all, nil
}

func BuildPinglistEmbed(target *PingTarget, matches []*discordgo.Message) *discordgo.MessageEmbed {
	name := target.Name
	if name == "" {
		name = target.ID
	}

	embed := &discordgo.MessageEmbed{
		Color:       0xED4245,
		Title:       fmt.Sprintf("Pings for %s", name),
		Description: fmt.Sprintf("Recent messages mentioning %s in this channel (last %d messages).", target.Kind, pinglistMaxLimit),
	}

	if len(matches) == 0 {
		embed.Description += "\n\nNo pings found."
		return embed
	}

	embed.Footer = &discordgo.MessageEmbedFooter{
		Text: fmt.Sprintf("%d ping(s) found", len(matches)),
	}

	count := PinglistMaxResults
	if len(matches) < count {
		count = len(matches)
	}
	for _, msg := range matches[:count] {
		author := msg.Author.Username
		if msg.Author.Bot {
			author += " [bot]"
		}
		content := msg.Content
		if len(content) > SearchMaxFieldLength {
			content = content[:SearchMaxFieldLength-3] + "..."
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

func SnowflakeFromValue(v any) string {
	switch val := v.(type) {
	case string:
		return val
	case float64:
		return fmt.Sprintf("%.0f", val)
	}
	return ""
}
