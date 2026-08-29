package commands

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

func userInfoMessageCommandHandler(s *discordgo.Session, m *discordgo.MessageCreate, args []string) error {
	targetID := m.Author.ID
	for _, arg := range args {
		if id := parseMentionID(arg); id != "" {
			targetID = id
			break
		}
	}

	embed, err := buildUserInfoEmbed(s, m.GuildID, targetID)
	if err != nil {
		return err
	}

	_, err = s.ChannelMessageSendEmbed(m.ChannelID, embed)
	return err
}

func userInfoSlashCommandHandler(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	targetID := i.Member.User.ID
	for _, opt := range i.ApplicationCommandData().Options {
		if opt.Name == "user" {
			if id, ok := opt.Value.(string); ok && id != "" {
				targetID = id
			}
		}
	}

	guildID := i.GuildID
	embed, err := buildUserInfoEmbed(s, guildID, targetID)
	if err != nil {
		return err
	}

	return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{embed},
		},
	})
}

func buildUserInfoEmbed(s *discordgo.Session, guildID, targetID string) (*discordgo.MessageEmbed, error) {
	user, err := s.User(targetID)
	if err != nil {
		return nil, err
	}

	member, err := s.GuildMember(guildID, targetID)
	if err != nil {
		return nil, err
	}

	nickname := member.Nick
	if nickname == "" {
		nickname = "—"
	}

	createdAt, err := discordgo.SnowflakeTimestamp(targetID)
	createdStr := "Unknown"
	if err == nil {
		createdStr = createdAt.Format(time.RFC3339)
	}

	joinedAt := member.JoinedAt.Format(time.RFC3339)

	roles := formatRoleNames(s, guildID, member.Roles)

	embed := &discordgo.MessageEmbed{
		Color: 0x5865F2,
		Title: fmt.Sprintf("%s#%s", user.Username, user.Discriminator),
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: user.AvatarURL("256"),
		},
		Description: fmt.Sprintf("<@%s>", targetID),
		Fields: []*discordgo.MessageEmbedField{
			{Name: "User ID", Value: targetID, Inline: true},
			{Name: "Nickname", Value: nickname, Inline: true},
			{Name: "Bot", Value: strconv.FormatBool(user.Bot), Inline: true},
			{Name: "Account Created", Value: createdStr},
			{Name: "Joined Server", Value: joinedAt},
			{Name: "Roles", Value: roles, Inline: true},
		},
	}

	return embed, nil
}

func formatRoleNames(s *discordgo.Session, guildID string, roleIDs []string) string {
	if len(roleIDs) == 0 {
		return "—"
	}

	roles, err := s.GuildRoles(guildID)
	if err != nil {
		return fmt.Sprintf("%d role(s)", len(roleIDs))
	}

	names := make([]string, 0, len(roleIDs))
	for _, id := range roleIDs {
		for _, r := range roles {
			if r.ID == id {
				names = append(names, r.Name)
				break
			}
		}
	}

	return strings.Join(names, ", ")
}

func parseMentionID(arg string) string {
	arg = strings.TrimSpace(arg)
	arg = strings.TrimPrefix(arg, "<@&")
	arg = strings.TrimPrefix(arg, "<@!")
	arg = strings.TrimPrefix(arg, "<@")
	arg = strings.TrimSuffix(arg, ">")
	if id, err := strconv.ParseInt(arg, 10, 64); err == nil && id > 0 {
		return strconv.FormatInt(id, 10)
	}
	return ""
}
