package commands

import (
	"fmt"
	"runtime"
	"runtime/debug"
	"sort"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

// botStartTime records when the process started so ```botinfo``` can report
// uptime without extra plumbing.
var botStartTime = time.Now()

// requestedUserID returns the first user/mention found in args (or the author
// ID if none), mirroring how other commands resolve a target user.
func requestedUserID(m *discordgo.MessageCreate, args []string) string {
	for _, arg := range args {
		if id := ParseMentionID(arg); id != "" {
			return id
		}
	}
	return m.Author.ID
}

// ---------------------------------------------------------------- avatar ----

func avatarMessageCommandHandler(s *discordgo.Session, m *discordgo.MessageCreate, args []string) error {
	targetID := requestedUserID(m, args)
	embed, err := buildAvatarEmbed(s, m.GuildID, targetID)
	if err != nil {
		return err
	}
	_, err = s.ChannelMessageSendEmbed(m.ChannelID, embed)
	return err
}

func avatarSlashCommandHandler(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	targetID := i.Member.User.ID
	for _, opt := range i.ApplicationCommandData().Options {
		if opt.Name == "user" {
			if id, ok := opt.Value.(string); ok && id != "" {
				targetID = id
			}
		}
	}
	embed, err := buildAvatarEmbed(s, i.GuildID, targetID)
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

func buildAvatarEmbed(s *discordgo.Session, guildID, targetID string) (*discordgo.MessageEmbed, error) {
	user, err := s.User(targetID)
	if err != nil {
		return nil, err
	}

	// Prefer the member/guild avatar when available, else the global avatar.
	url := user.AvatarURL("1024")
	if guildID != "" {
		member, merr := s.GuildMember(guildID, targetID)
		if merr == nil && member.Avatar != "" {
			url = member.AvatarURL("1024")
		}
	}

	return &discordgo.MessageEmbed{
		Color: 0x5865F2,
		Title: fmt.Sprintf("%s's avatar", user.Username),
		Image: &discordgo.MessageEmbedImage{URL: url},
	}, nil
}

// ---------------------------------------------------------------- banner ----

func bannerMessageCommandHandler(s *discordgo.Session, m *discordgo.MessageCreate, args []string) error {
	targetID := requestedUserID(m, args)
	embed, err := buildBannerEmbed(s, targetID)
	if err != nil {
		return err
	}
	_, err = s.ChannelMessageSendEmbed(m.ChannelID, embed)
	return err
}

func bannerSlashCommandHandler(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	targetID := i.Member.User.ID
	for _, opt := range i.ApplicationCommandData().Options {
		if opt.Name == "user" {
			if id, ok := opt.Value.(string); ok && id != "" {
				targetID = id
			}
		}
	}
	embed, err := buildBannerEmbed(s, targetID)
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

func buildBannerEmbed(s *discordgo.Session, targetID string) (*discordgo.MessageEmbed, error) {
	user, err := s.User(targetID)
	if err != nil {
		return nil, err
	}
	if user.Banner == "" {
		return &discordgo.MessageEmbed{
			Color:       0x5865F2,
			Title:       fmt.Sprintf("%s has no banner", user.Username),
			Description: "This user does not have a profile banner set.",
		}, nil
	}
	return &discordgo.MessageEmbed{
		Color: 0x5865F2,
		Title: fmt.Sprintf("%s's banner", user.Username),
		Image: &discordgo.MessageEmbedImage{URL: user.BannerURL("1024")},
	}, nil
}

// --------------------------------------------------------------- botinfo ----

func botinfoMessageCommandHandler(s *discordgo.Session, m *discordgo.MessageCreate, _ []string) error {
	embed := buildBotInfoEmbed(s)
	_, err := s.ChannelMessageSendEmbed(m.ChannelID, embed)
	return err
}

func botinfoSlashCommandHandler(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{buildBotInfoEmbed(s)},
		},
	})
}

func buildBotInfoEmbed(s *discordgo.Session) *discordgo.MessageEmbed {
	name := "MinjiBot"
	icon := ""
	if s.State != nil && s.State.User != nil {
		name = s.State.User.Username
		icon = s.State.User.AvatarURL("256")
	}

	version := "unknown"
	if bi, ok := debug.ReadBuildInfo(); ok {
		version = bi.Main.Version
		if version == "" || version == "(devel)" {
			version = "dev"
		}
	}

	uptime := time.Since(botStartTime).Round(time.Second)
	wsLatency := s.HeartbeatLatency().Milliseconds()

	embed := &discordgo.MessageEmbed{
		Color:     0x5865F2,
		Title:     name,
		Thumbnail: &discordgo.MessageEmbedThumbnail{URL: icon},
		Fields: []*discordgo.MessageEmbedField{
			{Name: "Version", Value: version, Inline: true},
			{Name: "Go", Value: runtime.Version(), Inline: true},
			{Name: "Uptime", Value: uptime.String(), Inline: true},
			{Name: "WebSocket Latency", Value: fmt.Sprintf("%dms", wsLatency), Inline: true},
			{Name: "Servers", Value: fmt.Sprintf("%d", len(s.State.Guilds)), Inline: true},
		},
	}
	return embed
}

// ----------------------------------------------------------- channelinfo ----

func channelinfoMessageCommandHandler(s *discordgo.Session, m *discordgo.MessageCreate, args []string) error {
	channelID := m.ChannelID
	for _, arg := range args {
		if id := ParseChannelMention(arg); id != "" {
			channelID = id
			break
		}
	}
	embed, err := buildChannelInfoEmbed(s, channelID)
	if err != nil {
		return err
	}
	_, err = s.ChannelMessageSendEmbed(m.ChannelID, embed)
	return err
}

func channelinfoSlashCommandHandler(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	channelID := i.ChannelID
	for _, opt := range i.ApplicationCommandData().Options {
		if opt.Name == "channel" {
			if id, ok := opt.Value.(string); ok && id != "" {
				channelID = id
			}
		}
	}
	embed, err := buildChannelInfoEmbed(s, channelID)
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

func buildChannelInfoEmbed(s *discordgo.Session, channelID string) (*discordgo.MessageEmbed, error) {
	ch, err := s.Channel(channelID)
	if err != nil {
		return nil, err
	}

	created, _ := discordgo.SnowflakeTimestamp(channelID)
	createdStr := "Unknown"
	if !created.IsZero() {
		createdStr = created.Format("January 2, 2006")
	}

	var b strings.Builder
	fmt.Fprintf(&b, "**Name:** %s\n", ch.Name)
	fmt.Fprintf(&b, "**ID:** %s\n", ch.ID)
	fmt.Fprintf(&b, "**Type:** %s\n", channelTypeName(ch.Type))
	fmt.Fprintf(&b, "**Created:** %s\n", createdStr)
	if ch.Topic != "" {
		fmt.Fprintf(&b, "**Topic:** %s\n", ch.Topic)
	}
	if ch.RateLimitPerUser > 0 {
		fmt.Fprintf(&b, "**Slowmode:** %ds\n", ch.RateLimitPerUser)
	} else {
		fmt.Fprintf(&b, "**Slowmode:** Off\n")
	}
	if ch.ParentID != "" {
		fmt.Fprintf(&b, "**Category:** <#%s>\n", ch.ParentID)
	}

	return &discordgo.MessageEmbed{
		Color:       0x5865F2,
		Title:       "Channel Info",
		Description: b.String(),
	}, nil
}

func channelTypeName(t discordgo.ChannelType) string {
	switch t {
	case discordgo.ChannelTypeGuildText:
		return "Text Channel"
	case discordgo.ChannelTypeGuildVoice:
		return "Voice Channel"
	case discordgo.ChannelTypeGuildCategory:
		return "Category"
	case discordgo.ChannelTypeGuildNews:
		return "Announcement Channel"
	case discordgo.ChannelTypeGuildForum:
		return "Forum"
	case discordgo.ChannelTypeGuildStageVoice:
		return "Stage Channel"
	default:
		return "Channel"
	}
}

// ---------------------------------------------------------------- roles ----

func rolesMessageCommandHandler(s *discordgo.Session, m *discordgo.MessageCreate, _ []string) error {
	pages, err := buildRolesPages(s, m.GuildID)
	if err != nil {
		return err
	}
	return sendPagedEmbeds(s, m.ChannelID, m.Author.ID, pages)
}

func rolesSlashCommandHandler(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	pages, err := buildRolesPages(s, i.GuildID)
	if err != nil {
		return err
	}
	return respondPaged(s, i, pages)
}

// maxEmbedTextBudget is a conservative character cap for each embed description
// so Discord's 6000-char embed limit is never exceeded. Titles/footers/fields
// also count toward the limit, so we stay well under.
const maxEmbedTextBudget = 5000

// pageEmbeds chunks raw lines into one embed page per chunk, each whose
// description stays within budget. Returns a single placeholder embed if lines
// is empty (callers usually override or check len first).
func pageEmbeds(title string, lines []string, budget int, color int) []*discordgo.MessageEmbed {
	var pages []*discordgo.MessageEmbed
	var b strings.Builder
	for _, line := range lines {
		if b.Len()+len(line) > budget {
			pages = append(pages, &discordgo.MessageEmbed{
				Color:       color,
				Title:       title,
				Description: b.String(),
				Footer:      &discordgo.MessageEmbedFooter{Text: fmt.Sprintf("Page %d", len(pages)+1)},
			})
			b.Reset()
		}
		b.WriteString(line)
	}
	if b.Len() > 0 || len(pages) == 0 {
		pages = append(pages, &discordgo.MessageEmbed{
			Color:       color,
			Title:       title,
			Description: b.String(),
			Footer:      &discordgo.MessageEmbedFooter{Text: fmt.Sprintf("Page %d", len(pages)+1)},
		})
	}
	return pages
}

// sendPagedEmbeds sends a multi-page embed list: paginated with reactions for a
// prefix/message context, or just the first page for a slash interaction
// (which can't use the reaction paginator without extra plumbing).
func sendPagedEmbeds(s *discordgo.Session, channelID, authorID string, pages []*discordgo.MessageEmbed) error {
	if len(pages) == 0 {
		return nil
	}
	if len(pages) == 1 {
		_, err := s.ChannelMessageSendEmbed(channelID, pages[0])
		return err
	}
	return paginateReactions(s, channelID, authorID, len(pages), func(page int) *discordgo.MessageEmbed {
		return pages[page]
	})
}

func respondPaged(s *discordgo.Session, i *discordgo.InteractionCreate, pages []*discordgo.MessageEmbed) error {
	if len(pages) == 0 {
		return nil
	}
	return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{pages[0]},
		},
	})
}

// buildRolesPages returns one embed page per chunk of server roles, sorted by
// position (descending), each with member counts.
func buildRolesPages(s *discordgo.Session, guildID string) ([]*discordgo.MessageEmbed, error) {
	roles, err := s.GuildRoles(guildID)
	if err != nil {
		return nil, err
	}

	// Count members per role from the state cache (may be partial for large
	// guilds unless the member chunk is complete).
	counts := map[string]int{}
	guild, _ := s.State.Guild(guildID)
	if guild != nil {
		for _, member := range guild.Members {
			for _, rid := range member.Roles {
				counts[rid]++
			}
		}
	}

	sort.Slice(roles, func(a, b int) bool {
		if roles[a].Position != roles[b].Position {
			return roles[a].Position > roles[b].Position
		}
		return roles[a].ID < roles[b].ID
	})

	lines := make([]string, 0, len(roles))
	for _, r := range roles {
		lines = append(lines, fmt.Sprintf("<@&%s> — `%d` members\n", r.ID, counts[r.ID]))
	}
	if len(lines) == 0 {
		return []*discordgo.MessageEmbed{{
			Color:       0x5865F2,
			Title:       "Server Roles",
			Description: "This server has no roles.",
		}}, nil
	}
	return pageEmbeds(fmt.Sprintf("Server Roles (%d)", len(roles)), lines, maxEmbedTextBudget, 0x5865F2), nil
}

// ---------------------------------------------------------------- guild -----

func guildMessageCommandHandler(s *discordgo.Session, m *discordgo.MessageCreate, args []string) error {
	sub := ""
	if len(args) > 0 {
		sub = strings.ToLower(args[0])
	}
	switch sub {
	case "stats":
		embed, err := buildGuildStatsEmbed(s, m.GuildID)
		if err != nil {
			return err
		}
		_, err = s.ChannelMessageSendEmbed(m.ChannelID, embed)
		return err
	case "icon":
		embed, err := buildGuildIconEmbed(s, m.GuildID)
		if err != nil {
			return err
		}
		_, err = s.ChannelMessageSendEmbed(m.ChannelID, embed)
		return err
	case "banner":
		embed, err := buildGuildBannerEmbed(s, m.GuildID)
		if err != nil {
			return err
		}
		_, err = s.ChannelMessageSendEmbed(m.ChannelID, embed)
		return err
	default:
		_, err := s.ChannelMessageSend(m.ChannelID, "Usage: `-guild stats|icon|banner`")
		return err
	}
}

func guildSlashCommandHandler(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	sub := i.ApplicationCommandData().Options
	var embed *discordgo.MessageEmbed
	var err error
	if len(sub) > 0 {
		switch sub[0].Name {
		case "stats":
			embed, err = buildGuildStatsEmbed(s, i.GuildID)
		case "icon":
			embed, err = buildGuildIconEmbed(s, i.GuildID)
		case "banner":
			embed, err = buildGuildBannerEmbed(s, i.GuildID)
		}
	}
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

func buildGuildStatsEmbed(s *discordgo.Session, guildID string) (*discordgo.MessageEmbed, error) {
	g, err := s.GuildWithCounts(guildID)
	if err != nil {
		return nil, err
	}

	created, _ := discordgo.SnowflakeTimestamp(guildID)
	createdStr := "Unknown"
	if !created.IsZero() {
		createdStr = created.Format("January 2, 2006")
	}

	boostTier := fmt.Sprintf("Level %d", g.PremiumTier)

	var b strings.Builder
	fmt.Fprintf(&b, "**Members:** %d\n", g.ApproximateMemberCount)
	fmt.Fprintf(&b, "**Online:** %d\n", g.ApproximatePresenceCount)
	fmt.Fprintf(&b, "**Owner:** <@%s>\n", g.OwnerID)
	fmt.Fprintf(&b, "**Boost Tier:** %s\n", boostTier)
	fmt.Fprintf(&b, "**Boosts:** %d\n", g.PremiumSubscriptionCount)
	fmt.Fprintf(&b, "**Created:** %s\n", createdStr)
	if g.Description != "" {
		fmt.Fprintf(&b, "**Description:** %s\n", g.Description)
	}

	embed := &discordgo.MessageEmbed{
		Color:       0x5865F2,
		Title:       g.Name,
		Thumbnail:   &discordgo.MessageEmbedThumbnail{URL: g.IconURL("256")},
		Description: b.String(),
	}
	return embed, nil
}

func buildGuildIconEmbed(s *discordgo.Session, guildID string) (*discordgo.MessageEmbed, error) {
	g, err := s.Guild(guildID)
	if err != nil {
		return nil, err
	}
	if g.Icon == "" {
		return &discordgo.MessageEmbed{
			Color:       0x5865F2,
			Title:       g.Name,
			Description: "This server has no icon set.",
		}, nil
	}
	return &discordgo.MessageEmbed{
		Color: 0x5865F2,
		Title: fmt.Sprintf("%s's icon", g.Name),
		Image: &discordgo.MessageEmbedImage{URL: g.IconURL("1024")},
	}, nil
}

func buildGuildBannerEmbed(s *discordgo.Session, guildID string) (*discordgo.MessageEmbed, error) {
	g, err := s.Guild(guildID)
	if err != nil {
		return nil, err
	}
	if g.Banner == "" {
		return &discordgo.MessageEmbed{
			Color:       0x5865F2,
			Title:       g.Name,
			Description: "This server has no banner set.",
		}, nil
	}
	return &discordgo.MessageEmbed{
		Color: 0x5865F2,
		Title: fmt.Sprintf("%s's banner", g.Name),
		Image: &discordgo.MessageEmbedImage{URL: g.BannerURL("1024")},
	}, nil
}

// --------------------------------------------------------------- emojis -----

func emojisMessageCommandHandler(s *discordgo.Session, m *discordgo.MessageCreate, _ []string) error {
	pages, err := buildEmojisPages(s, m.GuildID)
	if err != nil {
		return err
	}
	return sendPagedEmbeds(s, m.ChannelID, m.Author.ID, pages)
}

func emojisSlashCommandHandler(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	pages, err := buildEmojisPages(s, i.GuildID)
	if err != nil {
		return err
	}
	return respondPaged(s, i, pages)
}

func buildEmojisPages(s *discordgo.Session, guildID string) ([]*discordgo.MessageEmbed, error) {
	emojis, err := s.GuildEmojis(guildID)
	if err != nil {
		return nil, err
	}

	sort.Slice(emojis, func(a, b int) bool {
		return emojis[a].Name < emojis[b].Name
	})

	lines := make([]string, 0, len(emojis))
	for _, e := range emojis {
		lines = append(lines, fmt.Sprintf("<%s:%s:%s>\n", emojiNamePrefix(e.Animated), e.Name, e.ID))
	}
	if len(emojis) == 0 {
		return []*discordgo.MessageEmbed{{
			Color:       0x5865F2,
			Title:       "Server Emojis",
			Description: "No custom emojis in this server.",
		}}, nil
	}
	return pageEmbeds(fmt.Sprintf("Server Emojis (%d)", len(emojis)), lines, maxEmbedTextBudget, 0x5865F2), nil
}

func emojiNamePrefix(animated bool) string {
	if animated {
		return "a"
	}
	return ""
}

// -------------------------------------------------------------- stickers ----

func stickersMessageCommandHandler(s *discordgo.Session, m *discordgo.MessageCreate, _ []string) error {
	pages, err := buildStickersPages(s, m.GuildID)
	if err != nil {
		return err
	}
	return sendPagedEmbeds(s, m.ChannelID, m.Author.ID, pages)
}

func stickersSlashCommandHandler(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	pages, err := buildStickersPages(s, i.GuildID)
	if err != nil {
		return err
	}
	return respondPaged(s, i, pages)
}

func buildStickersPages(s *discordgo.Session, guildID string) ([]*discordgo.MessageEmbed, error) {
	g, err := s.Guild(guildID)
	if err != nil {
		return nil, err
	}

	var lines []string
	for _, st := range g.Stickers {
		lines = append(lines, fmt.Sprintf("**%s** — %s\n", st.Name, stickerFormatName(st.FormatType)))
	}
	if len(g.Stickers) == 0 {
		return []*discordgo.MessageEmbed{{
			Color:       0x5865F2,
			Title:       "Server Stickers",
			Description: "No custom stickers in this server.",
		}}, nil
	}
	return pageEmbeds(fmt.Sprintf("Server Stickers (%d)", len(g.Stickers)), lines, maxEmbedTextBudget, 0x5865F2), nil
}

func stickerFormatName(f discordgo.StickerFormat) string {
	switch f {
	case discordgo.StickerFormatTypePNG:
		return "PNG"
	case discordgo.StickerFormatTypeAPNG:
		return "Animated PNG"
	case discordgo.StickerFormatTypeLottie:
		return "Lottie"
	case discordgo.StickerFormatTypeGIF:
		return "GIF"
	default:
		return "Sticker"
	}
}
