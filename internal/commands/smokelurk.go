package commands

import (
	"fmt"
	"sort"
	"strings"
	"sync"

	"github.com/bwmarrin/discordgo"
)

// MemberState holds the ephemeral (in-memory) state for the local fun commands:
// who is lurking, and the shared "blunt" used by the smoke/spark/hits commands.
// State is keyed by guild so separate servers don't share a blunt or a lurk
// list. Everything lives only in memory — it resets when the bot restarts,
// which is fine for these jokes (there's no DB persistence involved).
type MemberState struct {
	mu sync.Mutex

	// lurkerSet maps guildID -> set of userIDs currently lurking.
	lurkerSet map[string]map[string]bool

	// blunt tracks the shared blunt per guild: the user who last sparked it
	// (empty string = not sparked) and each member's lifetime hit count.
	blunt  map[string]*guildBlunt
}

type guildBlunt struct {
	sparkedBy string
	hits      map[string]int
}

var memberState = &MemberState{
	lurkerSet: make(map[string]map[string]bool),
	blunt:     make(map[string]*guildBlunt),
}

func (s *MemberState) guildLurkers(guildID string) map[string]bool {
	set, ok := s.lurkerSet[guildID]
	if !ok {
		set = make(map[string]bool)
		s.lurkerSet[guildID] = set
	}
	return set
}

func (s *MemberState) guildBlunt(guildID string) *guildBlunt {
	b, ok := s.blunt[guildID]
	if !ok {
		b = &guildBlunt{hits: make(map[string]int)}
		s.blunt[guildID] = b
	}
	return b
}

// ---- lurk ----

func lurkMessageCommandHandler(s *discordgo.Session, m *discordgo.MessageCreate, args []string) error {
	guildID := m.GuildID
	userID := m.Author.ID
	_ = s.ChannelMessageDelete(m.ChannelID, m.ID)

	memberState.mu.Lock()
	set := memberState.guildLurkers(guildID)
	now := set[userID]
	if now {
		delete(set, userID)
	} else {
		set[userID] = true
	}
	count := len(set)
	memberState.mu.Unlock()

	if now {
		_, err := s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("<@%s> stops lurking. %d member(s) still lurking.", userID, count))
		return err
	}
	_, err := s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("<@%s> is now lurking in the shadows 👻. %d member(s) lurking.", userID, count))
	return err
}

func lurkersMessageCommandHandler(s *discordgo.Session, m *discordgo.MessageCreate, args []string) error {
	guildID := m.GuildID
	_ = s.ChannelMessageDelete(m.ChannelID, m.ID)

	memberState.mu.Lock()
	set := memberState.guildLurkers(guildID)
	ids := make([]string, 0, len(set))
	for id := range set {
		ids = append(ids, id)
	}
	memberState.mu.Unlock()

	if len(ids) == 0 {
		_, err := s.ChannelMessageSend(m.ChannelID, "Nobody is lurking right now. 🤷")
		return err
	}

	sort.Strings(ids)
	var b strings.Builder
	for _, id := range ids {
		b.WriteString(fmt.Sprintf("<@%s>\n", id))
	}
	_, err := s.ChannelMessageSendEmbed(m.ChannelID, &discordgo.MessageEmbed{
		Title:       "👀 Lurking members",
		Description: b.String(),
		Color:       0x2f3136,
	})
	return err
}

func lurkSlashCommandHandler(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	guildID := i.GuildID
	userID := i.Member.User.ID

	memberState.mu.Lock()
	set := memberState.guildLurkers(guildID)
	now := set[userID]
	if now {
		delete(set, userID)
	} else {
		set[userID] = true
	}
	count := len(set)
	memberState.mu.Unlock()

	var content string
	if now {
		content = fmt.Sprintf("<@%s> stops lurking. %d member(s) still lurking.", userID, count)
	} else {
		content = fmt.Sprintf("<@%s> is now lurking in the shadows 👻. %d member(s) lurking.", userID, count)
	}
	return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{Content: content},
	})
}

func lurkersSlashCommandHandler(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	guildID := i.GuildID

	memberState.mu.Lock()
	set := memberState.guildLurkers(guildID)
	ids := make([]string, 0, len(set))
	for id := range set {
		ids = append(ids, id)
	}
	memberState.mu.Unlock()

	if len(ids) == 0 {
		return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{Content: "Nobody is lurking right now. 🤷"},
		})
	}

	sort.Strings(ids)
	var b strings.Builder
	for _, id := range ids {
		b.WriteString(fmt.Sprintf("<@%s>\n", id))
	}
	return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{{
				Title:       "👀 Lurking members",
				Description: b.String(),
				Color:       0x2f3136,
			}},
		},
	})
}

// ---- smoke ----

func sparkMessageCommandHandler(s *discordgo.Session, m *discordgo.MessageCreate, args []string) error {
	guildID := m.GuildID
	userID := m.Author.ID
	_ = s.ChannelMessageDelete(m.ChannelID, m.ID)

	memberState.mu.Lock()
	b := memberState.guildBlunt(guildID)
	prev := b.sparkedBy
	b.sparkedBy = userID
	memberState.mu.Unlock()

	if prev == userID {
		_, err := s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("<@%s> keeps the flame lit 🌿🔥", userID))
		return err
	}
	_, err := s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("<@%s> sparked the blunt 🌿🔥", userID))
	return err
}

func smokeMessageCommandHandler(s *discordgo.Session, m *discordgo.MessageCreate, args []string) error {
	guildID := m.GuildID
	userID := m.Author.ID
	_ = s.ChannelMessageDelete(m.ChannelID, m.ID)

	memberState.mu.Lock()
	b := memberState.guildBlunt(guildID)
	if b.sparkedBy == "" {
		memberState.mu.Unlock()
		_, err := s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("<@%s>, the blunt isn't sparked — run `-spark` first. 🌱", userID))
		return err
	}

	content := ""
	if b.sparkedBy != userID {
		b.sparkedBy = userID
		b.hits[userID]++
		content = fmt.Sprintf("🔥 <@%s> grabs the blunt (the spark resets to them) and takes a hit 🚬", userID)
	} else {
		b.hits[userID]++
		content = fmt.Sprintf("<@%s> takes a hit off the blunt 🚬", userID)
	}
	b.sparkedBy = ""
	total := b.hits[userID]
	memberState.mu.Unlock()

	_, err := s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s — **%d** total hit(s).", content, total))
	return err
}

func hitsMessageCommandHandler(s *discordgo.Session, m *discordgo.MessageCreate, args []string) error {
	guildID := m.GuildID

	memberState.mu.Lock()
	b := memberState.guildBlunt(guildID)
	entries := make([]string, 0, len(b.hits))
	for id, n := range b.hits {
		entries = append(entries, fmt.Sprintf("<@%s> — **%d** hit(s)", id, n))
	}
	sparked := b.sparkedBy
	memberState.mu.Unlock()

	var body string
	if len(entries) == 0 {
		body = "No one has taken a hit yet. Spark one up and smoke! 🌿"
	} else {
		sort.Strings(entries)
		body = strings.Join(entries, "\n")
	}
	if sparked != "" {
		body = fmt.Sprintf("🌿 Blunt is sparked by <@%s>.\n\n%s", sparked, body)
	}
	_, err := s.ChannelMessageSendEmbed(m.ChannelID, &discordgo.MessageEmbed{
		Title:       "🚬 Blunt hit counter",
		Description: body,
		Color:       0x2ed573,
	})
	return err
}

func sparkSlashCommandHandler(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	guildID := i.GuildID
	userID := i.Member.User.ID

	memberState.mu.Lock()
	b := memberState.guildBlunt(guildID)
	prev := b.sparkedBy
	b.sparkedBy = userID
	memberState.mu.Unlock()

	content := fmt.Sprintf("<@%s> sparked the blunt 🌿🔥", userID)
	if prev == userID {
		content = fmt.Sprintf("<@%s> keeps the flame lit 🌿🔥", userID)
	}
	return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{Content: content},
	})
}

func smokeSlashCommandHandler(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	guildID := i.GuildID
	userID := i.Member.User.ID

	memberState.mu.Lock()
	b := memberState.guildBlunt(guildID)
	if b.sparkedBy == "" {
		memberState.mu.Unlock()
		return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: fmt.Sprintf("<@%s>, the blunt isn't sparked — run `/spark` first. 🌱", userID),
			},
		})
	}

	content := ""
	if b.sparkedBy != userID {
		b.sparkedBy = userID
		b.hits[userID]++
		content = fmt.Sprintf("🔥 <@%s> grabs the blunt (the spark resets to them) and takes a hit 🚬", userID)
	} else {
		b.hits[userID]++
		content = fmt.Sprintf("<@%s> takes a hit off the blunt 🚬", userID)
	}
	b.sparkedBy = ""
	total := b.hits[userID]
	memberState.mu.Unlock()

	return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf("%s — **%d** total hit(s).", content, total),
		},
	})
}

func hitsSlashCommandHandler(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	guildID := i.GuildID

	memberState.mu.Lock()
	b := memberState.guildBlunt(guildID)
	entries := make([]string, 0, len(b.hits))
	for id, n := range b.hits {
		entries = append(entries, fmt.Sprintf("<@%s> — **%d** hit(s)", id, n))
	}
	sparked := b.sparkedBy
	memberState.mu.Unlock()

	var body string
	if len(entries) == 0 {
		body = "No one has taken a hit yet. Spark one up and smoke! 🌿"
	} else {
		sort.Strings(entries)
		body = strings.Join(entries, "\n")
	}
	if sparked != "" {
		body = fmt.Sprintf("🌿 Blunt is sparked by <@%s>.\n\n%s", sparked, body)
	}
	return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{{
				Title:       "🚬 Blunt hit counter",
				Description: body,
				Color:       0x2ed573,
			}},
		},
	})
}
