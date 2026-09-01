package commands

import (
	"context"
	"math/rand"
	"sync"
	"time"

	"github.com/bwmarrin/discordgo"
)

// ttys makes the bot talk to itself in a channel until someone else sends a
// message there, or until an hour passes. Each channel runs at most one
// self-talk loop at a time.
var (
	ttysMu           sync.Mutex
	ttysChats        = map[string]context.CancelFunc{}
	ttysListenerOnce sync.Once
)

const ttysMaxDuration = time.Hour

const ttyIntroMessage = "I'll talk to myself for a while. Say anything in this channel to shut me up — or wait an hour and I'll stop on my own."

var ttyLines = []string{
	"you guys ever just",
	"no",
	"wait actually hold on",
	"yeah exactly",
	"anyway so",
	"moving on",
	"...nobody?",
	"right, back to what i was saying",
	"brain go brrr",
	"did i already say this",
	"so like hypothetically",
	"agreed",
	"great, glad we talked",
	"ok ok real talk now",
	"i was gonna say something but i forgot",
	"...and that's the story",
	"yep, that's what i thought",
	"heard",
	"noted",
	"can't believe it's taken this long",
	"stay with me, this is good",
	"ok ok, anyway",
	"someone back me up here",
	"no wait, wrong channel",
	"crickets. ok.",
	"punctuation is overrated",
	"has anyone actually been reading these",
	"remind me why i do this",
	"the silence is loud",
	"fine, i'll just talk to myself then",
}

func ttysMessageCommandHandler(s *discordgo.Session, m *discordgo.MessageCreate, _ []string) error {
	if !ttysBegin(s, m.ChannelID) {
		_, err := s.ChannelMessageSend(m.ChannelID, "I'm already talking to myself here — say something to stop me.")
		return err
	}
	_, err := s.ChannelMessageSend(m.ChannelID, ttyIntroMessage)
	return err
}

func ttysSlashCommandHandler(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	if !ttysBegin(s, i.ChannelID) {
		return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{Content: "I'm already talking to myself here — say something to stop me."},
		})
	}
	return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{Content: ttyIntroMessage},
	})
}

// ttysBegin registers a self-talk loop for the channel, returning false if one
// is already running there.
func ttysBegin(s *discordgo.Session, channelID string) bool {
	ttysMu.Lock()
	if _, active := ttysChats[channelID]; active {
		ttysMu.Unlock()
		return false
	}
	ctx, cancel := context.WithCancel(context.Background())
	ttysChats[channelID] = cancel
	ttysMu.Unlock()

	if s != nil {
		ttysListenerOnce.Do(func() {
			s.AddHandler(ttyOnMessage)
		})
	}
	go ttyLoop(s, channelID, ctx)
	return true
}

// ttyLoop keeps sending self-conversation lines until someone speaks in the
// channel (ctx cancelled) or the hour-long timer expires.
func ttyLoop(s *discordgo.Session, channelID string, ctx context.Context) {
	defer func() {
		ttysMu.Lock()
		delete(ttysChats, channelID)
		ttysMu.Unlock()
	}()

	duration := time.NewTimer(ttysMaxDuration)
	defer duration.Stop()

	for {
		select {
		case <-ctx.Done():
			_, _ = s.ChannelMessageSend(channelID, "phew — someone said something, I'll shut up.")
			return
		case <-duration.C:
			_, _ = s.ChannelMessageSend(channelID, "ok, that's a full hour of me talking to myself. taking a break.")
			return
		case <-time.After(15 + time.Duration(rand.Intn(30))*time.Second):
			_, err := s.ChannelMessageSend(channelID, ttyNextLine())
			if err != nil {
				return
			}
		}
	}
}

func ttyNextLine() string {
	if rand.Intn(3) == 0 {
		return "... " + ttyLines[rand.Intn(len(ttyLines))]
	}
	return ttyLines[rand.Intn(len(ttyLines))]
}

// ttyOnMessage stops any self-talk loop running in a channel the moment a
// non-bot message is sent there.
func ttyOnMessage(_ *discordgo.Session, mc *discordgo.MessageCreate) {
	if mc.Author != nil && mc.Author.Bot {
		return
	}
	ttysMu.Lock()
	cancel, active := ttysChats[mc.ChannelID]
	ttysMu.Unlock()
	if active {
		cancel()
	}
}
