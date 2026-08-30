package commands

import (
	"bytes"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"sort"
	"strings"

	"github.com/bwmarrin/discordgo"
)

// DominantColors extracts the n most frequent colours from image bytes,
// returned as #rrggbb hex strings. Ignores transparent/low-alpha pixels.
func DominantColors(data []byte, n int) []string {
	src, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return nil
	}
	// Coarse quantization: 8 levels per channel -> 512 buckets.
	type bucket struct {
		r, g, b int
		count   int
	}
	buckets := make(map[int]*bucket)
	bounds := src.Bounds()
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			cr, cg, cb, ca := src.At(x, y).RGBA()
			if ca>>8 < 128 {
				continue
			}
			r := int(cr>>8) / 32
			g := int(cg>>8) / 32
			b := int(cb>>8) / 32
			key := (r << 16) | (g << 8) | b
			bk, ok := buckets[key]
			if !ok {
				bk = &bucket{r: r * 32, g: g * 32, b: b * 32}
				buckets[key] = bk
			}
			bk.count++
		}
	}
	all := make([]*bucket, 0, len(buckets))
	for _, bk := range buckets {
		all = append(all, bk)
	}
	sort.Slice(all, func(i, j int) bool { return all[i].count > all[j].count })
	if len(all) > n {
		all = all[:n]
	}
	out := make([]string, 0, len(all))
	for _, bk := range all {
		out = append(out, fmt.Sprintf("#%02x%02x%02x", bk.r, bk.g, bk.b))
	}
	return out
}

func colorsAvatarMessageCommandHandler(s *discordgo.Session, m *discordgo.MessageCreate, args []string) error {
	if len(args) == 0 || args[0] != "avatar" {
		_, err := s.ChannelMessageSend(m.ChannelID, "Usage: `-colors avatar [user]`")
		return err
	}
	id := funTargetID(m, args[1:])
	return sendAvatarColors(s, m.ChannelID, id)
}

func colorsAvatarSlashCommandHandler(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	id := targetIDFromSlash(i)
	if id == "" {
		id = optMemberUserID(i)
	}

	user, err := s.User(id)
	if err != nil {
		return err
	}
	md, err := fetchURL(user.AvatarURL("512"))
	if err != nil {
		return err
	}
	embed := buildColorEmbed(user, DominantColors(md.Data, 5))
	return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{Embeds: []*discordgo.MessageEmbed{embed}},
	})
}

// targetIDFromSlash reads the "user" option, unwrapping a single subcommand
// level if present.
func targetIDFromSlash(i *discordgo.InteractionCreate) string {
	opts := i.ApplicationCommandData().Options
	for _, o := range opts {
		if o.Name == "user" {
			if id, ok := o.Value.(string); ok && id != "" {
				return id
			}
		}
		if o.Type == discordgo.ApplicationCommandOptionSubCommand {
			for _, sub := range o.Options {
				if sub.Name == "user" {
					if id, ok := sub.Value.(string); ok && id != "" {
						return id
					}
				}
			}
		}
	}
	return ""
}

func optMemberUserID(i *discordgo.InteractionCreate) string {
	if i != nil && i.Member != nil && i.Member.User != nil {
		return i.Member.User.ID
	}
	return ""
}

func sendAvatarColors(s *discordgo.Session, channelID, userID string) error {
	user, err := s.User(userID)
	if err != nil {
		return err
	}
	md, err := fetchURL(user.AvatarURL("512"))
	if err != nil {
		return err
	}
	embed := buildColorEmbed(user, DominantColors(md.Data, 5))
	_, err = s.ChannelMessageSendEmbed(channelID, embed)
	return err
}

func buildColorEmbed(user *discordgo.User, colors []string) *discordgo.MessageEmbed {
	var b strings.Builder
	for _, c := range colors {
		b.WriteString(fmt.Sprintf("▰ **%s**\n", strings.ToUpper(c)))
	}
	return &discordgo.MessageEmbed{
		Color:       0x5865F2,
		Title:       user.Username + "'s avatar colours",
		Thumbnail:   &discordgo.MessageEmbedThumbnail{URL: user.AvatarURL("256")},
		Description: b.String(),
	}
}
