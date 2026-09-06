package commands

import (
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/kibetnathan/minjibot/internal/safe"
)

const lurkAutoExpire = 1 * time.Hour

// lurkState holds who is currently lurking per guild and tracks auto-unlurk
// timers. State lives only in memory and resets on restart.
type lurkState struct {
	mu      sync.Mutex
	lurkers map[string]map[string]bool // guildID -> userID set
	timers  map[string]func()          // "guildID:userID" -> cancel func
}

var lurkStateSingleton = &lurkState{
	lurkers: make(map[string]map[string]bool),
	timers:  make(map[string]func()),
}

func (ls *lurkState) guildLurkers(guildID string) map[string]bool {
	set, ok := ls.lurkers[guildID]
	if !ok {
		set = make(map[string]bool)
		ls.lurkers[guildID] = set
	}
	return set
}

// IsLurking reports whether a user is currently lurking in the given guild.
func IsLurking(guildID, userID string) bool {
	lurkStateSingleton.mu.Lock()
	defer lurkStateSingleton.mu.Unlock()
	return lurkStateSingleton.guildLurkers(guildID)[userID]
}

// scheduleAutoUnlurk starts a timer that removes a user from the lurk list
// after lurkAutoExpire. If the user stops lurking before the timer fires, the
// timer is cancelled. Calling scheduleAutoUnlurk for an already-lurking user
// replaces the previous timer.
func scheduleAutoUnlurk(guildID, userID string) {
	key := guildID + ":" + userID
	lurkStateSingleton.mu.Lock()

	// Cancel any previous timer for this user.
	if cancel, ok := lurkStateSingleton.timers[key]; ok {
		cancel()
	}

	done := make(chan struct{})
	lurkStateSingleton.timers[key] = func() { close(done) }
	lurkStateSingleton.mu.Unlock()

	safe.Go(nil, "lurkAutoExpire", func() {
		select {
		case <-time.After(lurkAutoExpire):
			lurkStateSingleton.mu.Lock()
			set := lurkStateSingleton.guildLurkers(guildID)
			delete(set, userID)
			delete(lurkStateSingleton.timers, key)
			lurkStateSingleton.mu.Unlock()
		case <-done:
			// Timer was cancelled (user manually stopped lurking).
		}
	})
}

func lurkMessageCommandHandler(s *discordgo.Session, m *discordgo.MessageCreate, args []string) error {
	guildID := m.GuildID
	userID := m.Author.ID

	lurkStateSingleton.mu.Lock()
	set := lurkStateSingleton.guildLurkers(guildID)
	wasLurking := set[userID]
	var schedule bool
	if wasLurking {
		delete(set, userID)
		// Cancel any pending auto-unlurk timer.
		key := guildID + ":" + userID
		if cancel, ok := lurkStateSingleton.timers[key]; ok {
			cancel()
			delete(lurkStateSingleton.timers, key)
		}
	} else {
		set[userID] = true
		schedule = true
	}
	count := len(set)
	lurkStateSingleton.mu.Unlock()

	// Schedule the auto-unlurk timer outside the lock: scheduleAutoUnlurk
	// acquires the mutex itself, and sync.Mutex is not reentrant.
	if schedule {
		scheduleAutoUnlurk(guildID, userID)
	}

	if wasLurking {
		_, err := s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("<@%s> stops lurking. %d member(s) still lurking.", userID, count))
		return err
	}
	_, err := s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("<@%s> is now lurking in the shadows (auto-expires in 1 hour). %d member(s) lurking.", userID, count))
	return err
}

func lurkersMessageCommandHandler(s *discordgo.Session, m *discordgo.MessageCreate, args []string) error {
	guildID := m.GuildID

	lurkStateSingleton.mu.Lock()
	set := lurkStateSingleton.guildLurkers(guildID)
	ids := make([]string, 0, len(set))
	for id := range set {
		ids = append(ids, id)
	}
	lurkStateSingleton.mu.Unlock()

	if len(ids) == 0 {
		_, err := s.ChannelMessageSend(m.ChannelID, "Nobody is lurking right now.")
		return err
	}

	sort.Strings(ids)
	var b strings.Builder
	for _, id := range ids {
		b.WriteString(fmt.Sprintf("<@%s>\n", id))
	}
	_, err := s.ChannelMessageSendEmbed(m.ChannelID, &discordgo.MessageEmbed{
		Title:       "Lurking members",
		Description: b.String(),
		Color:       0x2f3136,
	})
	return err
}

func lurkSlashCommandHandler(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	guildID := i.GuildID
	userID := i.Member.User.ID

	lurkStateSingleton.mu.Lock()
	set := lurkStateSingleton.guildLurkers(guildID)
	wasLurking := set[userID]
	var schedule bool
	if wasLurking {
		delete(set, userID)
		key := guildID + ":" + userID
		if cancel, ok := lurkStateSingleton.timers[key]; ok {
			cancel()
			delete(lurkStateSingleton.timers, key)
		}
	} else {
		set[userID] = true
		schedule = true
	}
	count := len(set)
	lurkStateSingleton.mu.Unlock()

	// Schedule outside the lock (sync.Mutex is not reentrant).
	if schedule {
		scheduleAutoUnlurk(guildID, userID)
	}

	var content string
	if wasLurking {
		content = fmt.Sprintf("<@%s> stops lurking. %d member(s) still lurking.", userID, count)
	} else {
		content = fmt.Sprintf("<@%s> is now lurking in the shadows (auto-expires in 1 hour). %d member(s) lurking.", userID, count)
	}
	return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{Content: content},
	})
}

func lurkersSlashCommandHandler(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	guildID := i.GuildID

	lurkStateSingleton.mu.Lock()
	set := lurkStateSingleton.guildLurkers(guildID)
	ids := make([]string, 0, len(set))
	for id := range set {
		ids = append(ids, id)
	}
	lurkStateSingleton.mu.Unlock()

	if len(ids) == 0 {
		return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{Content: "Nobody is lurking right now."},
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
				Title:       "Lurking members",
				Description: b.String(),
				Color:       0x2f3136,
			}},
		},
	})
}
