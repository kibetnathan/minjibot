package commands

import (
	"fmt"
	"sort"
	"strings"
	"sync"

	"github.com/bwmarrin/discordgo"
)

// bluntState holds the shared blunt per guild: who sparked it and hit counts.
// State lives only in memory and resets on restart.
type bluntState struct {
	mu     sync.Mutex
	blunts map[string]*guildBlunt // guildID -> blunt
}

type guildBlunt struct {
	sparkedBy string
	hits      map[string]int
}

var bluntStateSingleton = &bluntState{
	blunts: make(map[string]*guildBlunt),
}

func (bs *bluntState) guildBlunt(guildID string) *guildBlunt {
	b, ok := bs.blunts[guildID]
	if !ok {
		b = &guildBlunt{hits: make(map[string]int)}
		bs.blunts[guildID] = b
	}
	return b
}

func sparkMessageCommandHandler(s *discordgo.Session, m *discordgo.MessageCreate, args []string) error {
	guildID := m.GuildID
	userID := m.Author.ID

	bluntStateSingleton.mu.Lock()
	b := bluntStateSingleton.guildBlunt(guildID)
	prev := b.sparkedBy
	b.sparkedBy = userID
	bluntStateSingleton.mu.Unlock()

	if prev == userID {
		_, err := s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("<@%s> keeps the flame lit", userID))
		return err
	}
	_, err := s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("<@%s> sparked the blunt", userID))
	return err
}

func smokeMessageCommandHandler(s *discordgo.Session, m *discordgo.MessageCreate, args []string) error {
	guildID := m.GuildID
	userID := m.Author.ID

	bluntStateSingleton.mu.Lock()
	b := bluntStateSingleton.guildBlunt(guildID)
	if b.sparkedBy == "" {
		bluntStateSingleton.mu.Unlock()
		_, err := s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("<@%s>, the blunt isn't sparked - run `-spark` first.", userID))
		return err
	}

	content := ""
	if b.sparkedBy != userID {
		b.sparkedBy = userID
		b.hits[userID]++
		content = fmt.Sprintf("<@%s> grabs the blunt (the spark resets to them) and takes a hit", userID)
	} else {
		b.hits[userID]++
		content = fmt.Sprintf("<@%s> takes a hit off the blunt", userID)
	}
	b.sparkedBy = ""
	total := b.hits[userID]
	bluntStateSingleton.mu.Unlock()

	_, err := s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s - **%d** total hit(s).", content, total))
	return err
}

func hitsMessageCommandHandler(s *discordgo.Session, m *discordgo.MessageCreate, args []string) error {
	guildID := m.GuildID

	bluntStateSingleton.mu.Lock()
	b := bluntStateSingleton.guildBlunt(guildID)
	entries := make([]string, 0, len(b.hits))
	for id, n := range b.hits {
		entries = append(entries, fmt.Sprintf("<@%s> - **%d** hit(s)", id, n))
	}
	sparked := b.sparkedBy
	bluntStateSingleton.mu.Unlock()

	var body string
	if len(entries) == 0 {
		body = "No one has taken a hit yet. Spark one up and smoke!"
	} else {
		sort.Strings(entries)
		body = strings.Join(entries, "\n")
	}
	if sparked != "" {
		body = fmt.Sprintf("Blunt is sparked by <@%s>.\n\n%s", sparked, body)
	}
	_, err := s.ChannelMessageSendEmbed(m.ChannelID, &discordgo.MessageEmbed{
		Title:       "Blunt hit counter",
		Description: body,
		Color:       0x2ed573,
	})
	return err
}

func sparkSlashCommandHandler(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	guildID := i.GuildID
	userID := i.Member.User.ID

	bluntStateSingleton.mu.Lock()
	b := bluntStateSingleton.guildBlunt(guildID)
	prev := b.sparkedBy
	b.sparkedBy = userID
	bluntStateSingleton.mu.Unlock()

	content := fmt.Sprintf("<@%s> sparked the blunt", userID)
	if prev == userID {
		content = fmt.Sprintf("<@%s> keeps the flame lit", userID)
	}
	return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{Content: content},
	})
}

func smokeSlashCommandHandler(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	guildID := i.GuildID
	userID := i.Member.User.ID

	bluntStateSingleton.mu.Lock()
	b := bluntStateSingleton.guildBlunt(guildID)
	if b.sparkedBy == "" {
		bluntStateSingleton.mu.Unlock()
		return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: fmt.Sprintf("<@%s>, the blunt isn't sparked - run `/spark` first.", userID),
			},
		})
	}

	content := ""
	if b.sparkedBy != userID {
		b.sparkedBy = userID
		b.hits[userID]++
		content = fmt.Sprintf("<@%s> grabs the blunt (the spark resets to them) and takes a hit", userID)
	} else {
		b.hits[userID]++
		content = fmt.Sprintf("<@%s> takes a hit off the blunt", userID)
	}
	b.sparkedBy = ""
	total := b.hits[userID]
	bluntStateSingleton.mu.Unlock()

	return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf("%s - **%d** total hit(s).", content, total),
		},
	})
}

func hitsSlashCommandHandler(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	guildID := i.GuildID

	bluntStateSingleton.mu.Lock()
	b := bluntStateSingleton.guildBlunt(guildID)
	entries := make([]string, 0, len(b.hits))
	for id, n := range b.hits {
		entries = append(entries, fmt.Sprintf("<@%s> - **%d** hit(s)", id, n))
	}
	sparked := b.sparkedBy
	bluntStateSingleton.mu.Unlock()

	var body string
	if len(entries) == 0 {
		body = "No one has taken a hit yet. Spark one up and smoke!"
	} else {
		sort.Strings(entries)
		body = strings.Join(entries, "\n")
	}
	if sparked != "" {
		body = fmt.Sprintf("Blunt is sparked by <@%s>.\n\n%s", sparked, body)
	}
	return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{{
				Title:       "Blunt hit counter",
				Description: body,
				Color:       0x2ed573,
			}},
		},
	})
}
