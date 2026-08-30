package commands

import (
	"fmt"
	"sort"
	"strings"
	"sync"

	"github.com/bwmarrin/discordgo"
)

// vapeState holds the shared vape per guild: its flavour, who currently holds
// it, and per-user hit counts. State lives only in memory and resets on
// restart.
type vapeState struct {
	mu    sync.Mutex
	vapes map[string]*guildVape // guildID -> vape
}

type guildVape struct {
	holder string
	hits   map[string]int
	flavor string
}

var vapeStateSingleton = &vapeState{
	vapes: make(map[string]*guildVape),
}

func (vs *vapeState) guildVape(guildID string) *guildVape {
	v, ok := vs.vapes[guildID]
	if !ok {
		v = &guildVape{hits: make(map[string]int)}
		vs.vapes[guildID] = v
	}
	return v
}

// vapeMessageCommandHandler handles the `-vape` prefix command and its
// subcommands: plain `-vape`, `-vape flavour <text>`, `-vape hits`, and
// `-vape steal [@user]`.
func vapeMessageCommandHandler(s *discordgo.Session, m *discordgo.MessageCreate, args []string) error {
	userID := m.Author.ID

	sub := ""
	if len(args) > 0 {
		sub = strings.ToLower(args[0])
	}

	switch sub {
	case "flavor", "flavour":
		return vapeSetFlavorMessage(s, m, args[1:])
	case "hits", "count":
		return vapeHitsMessage(s, m)
	case "steal":
		return vapeStealMessage(s, m, args[1:])
	default:
		return vapeHitMessage(s, m, userID)
	}
}

func vapeHitMessage(s *discordgo.Session, m *discordgo.MessageCreate, userID string) error {
	guildID := m.GuildID

	vapeStateSingleton.mu.Lock()
	v := vapeStateSingleton.guildVape(guildID)
	flavor := v.flavor
	var msg string
	if v.holder == "" || v.holder == userID {
		// Nobody holds it or you already do — just hit.
		if v.holder == "" {
			v.holder = userID
		}
		v.hits[userID]++
		msg = "takes a hit off the vape"
	} else {
		// Someone else holds the vape: it passes to you.
		v.holder = userID
		v.hits[userID]++
		msg = "grabs the vape and takes a hit"
	}
	total := v.hits[userID]
	vapeStateSingleton.mu.Unlock()

	flavorTxt := ""
	if flavor != "" {
		flavorTxt = fmt.Sprintf(" (flavour: %s)", flavor)
	}
	_, err := s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("<@%s> %s%s — **%d** total hit(s).", userID, msg, flavorTxt, total))
	return err
}

func vapeSetFlavorMessage(s *discordgo.Session, m *discordgo.MessageCreate, args []string) error {
	guildID := m.GuildID
	flavor := strings.TrimSpace(strings.Join(args, " "))

	vapeStateSingleton.mu.Lock()
	v := vapeStateSingleton.guildVape(guildID)
	v.flavor = flavor
	vapeStateSingleton.mu.Unlock()

	if flavor == "" {
		_, err := s.ChannelMessageSend(m.ChannelID, "The vape flavour has been cleared.")
		return err
	}
	_, err := s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("The vape flavour is now **%s**.", flavor))
	return err
}

func vapeHitsMessage(s *discordgo.Session, m *discordgo.MessageCreate) error {
	guildID := m.GuildID

	vapeStateSingleton.mu.Lock()
	v := vapeStateSingleton.guildVape(guildID)
	type kv struct {
		id string
		n  int
	}
	entries := make([]kv, 0, len(v.hits))
	for id, n := range v.hits {
		entries = append(entries, kv{id, n})
	}
	flavor := v.flavor
	holder := v.holder
	vapeStateSingleton.mu.Unlock()

	sort.Slice(entries, func(i, j int) bool { return entries[i].id < entries[j].id })

	var b strings.Builder
	if flavor != "" {
		fmt.Fprintf(&b, "Flavour: %s\n\n", flavor)
	}
	if holder != "" {
		fmt.Fprintf(&b, "Currently held by <@%s>.\n\n", holder)
	}
	if len(entries) == 0 {
		b.WriteString("No one has taken a hit yet. Hit the vape!")
	} else {
		for _, e := range entries {
			fmt.Fprintf(&b, "<@%s> — **%d** hit(s)\n", e.id, e.n)
		}
	}

	_, err := s.ChannelMessageSendEmbed(m.ChannelID, &discordgo.MessageEmbed{
		Title:       "Vape hit counter",
		Description: b.String(),
		Color:       0x5865F2,
	})
	return err
}

func vapeStealMessage(s *discordgo.Session, m *discordgo.MessageCreate, args []string) error {
	guildID := m.GuildID
	userID := m.Author.ID

	vapeStateSingleton.mu.Lock()
	v := vapeStateSingleton.guildVape(guildID)
	prev := v.holder
	v.holder = userID
	v.hits[userID]++
	total := v.hits[userID]
	vapeStateSingleton.mu.Unlock()

	msg := ""
	if prev == "" || prev == userID {
		msg = fmt.Sprintf("<@%s> grabs the vape and takes a hit — **%d** total hit(s).", userID, total)
	} else {
		msg = fmt.Sprintf("<@%s> steals the vape from <@%s> and takes a hit — **%d** total hit(s).", userID, prev, total)
	}
	_, err := s.ChannelMessageSend(m.ChannelID, msg)
	return err
}

// vapeSlashCommandHandler handles the `/vape` slash command and its
// subcommands.
func vapeSlashCommandHandler(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	guildID := i.GuildID
	userID := i.Member.User.ID

	sub := ""
	flavor := ""
	for _, opt := range i.ApplicationCommandData().Options {
		if opt.Name == "flavour" {
			sub = "flavor"
			flavor = opt.StringValue()
		}
	}
	if sub != "flavor" && len(i.ApplicationCommandData().Options) > 0 {
		sub = strings.ToLower(i.ApplicationCommandData().Options[0].Name)
	}

	var content string
	switch sub {
	case "flavor":
		vapeStateSingleton.mu.Lock()
		v := vapeStateSingleton.guildVape(guildID)
		v.flavor = flavor
		vapeStateSingleton.mu.Unlock()
		if flavor == "" {
			content = "The vape flavour has been cleared."
		} else {
			content = fmt.Sprintf("The vape flavour is now **%s**.", flavor)
		}
	case "hits":
		vapeStateSingleton.mu.Lock()
		v := vapeStateSingleton.guildVape(guildID)
		var b strings.Builder
		if v.flavor != "" {
			fmt.Fprintf(&b, "Flavour: %s\n\n", v.flavor)
		}
		if v.holder != "" {
			fmt.Fprintf(&b, "Currently held by <@%s>.\n\n", v.holder)
		}
		if len(v.hits) == 0 {
			b.WriteString("No one has taken a hit yet. Hit the vape!")
		} else {
			hits := make([]struct {
				id string
				n  int
			}, 0, len(v.hits))
			for id, n := range v.hits {
				hits = append(hits, struct {
					id string
					n  int
				}{id, n})
			}
			sort.Slice(hits, func(a, b int) bool { return hits[a].id < hits[b].id })
			for _, h := range hits {
				fmt.Fprintf(&b, "<@%s> — **%d** hit(s)\n", h.id, h.n)
			}
		}
		vapeStateSingleton.mu.Unlock()
		content = b.String()
	case "steal":
		vapeStateSingleton.mu.Lock()
		v := vapeStateSingleton.guildVape(guildID)
		prev := v.holder
		v.holder = userID
		v.hits[userID]++
		total := v.hits[userID]
		vapeStateSingleton.mu.Unlock()
		if prev == "" || prev == userID {
			content = fmt.Sprintf("<@%s> grabs the vape and takes a hit — **%d** total hit(s).", userID, total)
		} else {
			content = fmt.Sprintf("<@%s> steals the vape from <@%s> and takes a hit — **%d** total hit(s).", userID, prev, total)
		}
	default:
		vapeStateSingleton.mu.Lock()
		v := vapeStateSingleton.guildVape(guildID)
		var msg string
		if v.holder == "" || v.holder == userID {
			if v.holder == "" {
				v.holder = userID
			}
			v.hits[userID]++
			msg = "takes a hit off the vape"
		} else {
			v.holder = userID
			v.hits[userID]++
			msg = "grabs the vape and takes a hit"
		}
		total := v.hits[userID]
		flavour := v.flavor
		vapeStateSingleton.mu.Unlock()

		flavourTxt := ""
		if flavour != "" {
			flavourTxt = fmt.Sprintf(" (flavour: %s)", flavour)
		}
		content = fmt.Sprintf("<@%s> %s%s — **%d** total hit(s).", userID, msg, flavourTxt, total)
	}

	return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{Content: content},
	})
}
