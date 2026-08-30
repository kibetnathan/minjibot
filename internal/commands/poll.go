package commands

import (
	"fmt"
	"strings"
	"sync"

	"github.com/bwmarrin/discordgo"
)

// Poll option emoji (1-10) used for reaction voting.
var pollEmojis = []string{"1️⃣", "2️⃣", "3️⃣", "4️⃣", "5️⃣", "6️⃣", "7️⃣", "8️⃣", "9️⃣", "🔟"}

const maxPollOptions = 10

// PollContext describes an active poll: the question, option labels, and how
// many votes each option currently has.
type PollContext struct {
	Question string
	Options  []string
	Votes    []int
	Closed   bool
	IsQuick  bool
}

// pollManager tracks active polls keyed by message ID and counts votes by
// reacting to message reactions.
type pollManager struct {
	mu       sync.Mutex
	polls    map[string]*PollContext
	channels map[string]string // messageID -> channelID
	once     sync.Once
}

// pollManagerSingleton is the shared poll state for reaction-based polls.
var pollManagerSingleton = &pollManager{
	polls:    make(map[string]*PollContext),
	channels: make(map[string]string),
}

// ParsePollOptions parses poll args into a question and a list of options.
// Options may be separated by a `|` or given as bare arguments. At least one
// option is required.
func ParsePollOptions(args []string) (question string, options []string, err error) {
	if len(args) == 0 {
		return "", nil, fmt.Errorf("Usage: `-poll <question>` with options, e.g. `-poll \"Favourite colour\" red blue green`")
	}

	var rest []string
	question = args[0]
	rest = args[1:]

	// If any option contains a `|`, split on it.
	if len(rest) > 0 {
		var split []string
		for _, r := range rest {
			for _, part := range strings.Split(r, "|") {
				if t := strings.TrimSpace(part); t != "" {
					split = append(split, t)
				}
			}
		}
		if len(split) > 0 {
			rest = split
		}
	}

	options = make([]string, 0, len(rest))
	for _, o := range rest {
		if t := strings.TrimSpace(o); t != "" {
			options = append(options, t)
		}
	}

	if len(options) == 0 {
		return "", nil, fmt.Errorf("A poll needs at least one option (separate options with `|` or spaces).")
	}
	if len(options) > maxPollOptions {
		return "", nil, fmt.Errorf("A poll supports at most %d options.", maxPollOptions)
	}
	return question, options, nil
}

// BuildPollEmbed renders the poll with its current vote counts.
func BuildPollEmbed(q string, options []string, votes []int) *discordgo.MessageEmbed {
	var b strings.Builder
	for i, opt := range options {
		votesTxt := ""
		if i < len(votes) && votes[i] > 0 {
			votesTxt = fmt.Sprintf(" — **%d**", votes[i])
		}
		fmt.Fprintf(&b, "%s %s%s\n", pollEmojis[i], opt, votesTxt)
	}
	return &discordgo.MessageEmbed{
		Color:       0x5865F2,
		Title:       "📊 " + q,
		Description: b.String(),
		Footer:      &discordgo.MessageEmbedFooter{Text: "Vote with the reactions below"},
	}
}

func pollMessageCommandHandler(s *discordgo.Session, m *discordgo.MessageCreate, args []string) error {
	question, options, err := ParsePollOptions(args)
	if err != nil {
		if _, serr := s.ChannelMessageSend(m.ChannelID, err.Error()); serr != nil {
			return serr
		}
		return nil
	}

	ctx := &PollContext{Question: question, Options: options, Votes: make([]int, len(options))}
	msg, err := s.ChannelMessageSendEmbed(m.ChannelID, BuildPollEmbed(question, options, ctx.Votes))
	if err != nil {
		return err
	}

	pollManagerSingleton.register(s, msg.ID, m.ChannelID, ctx)
	for idx := 0; idx < len(options); idx++ {
		_ = s.MessageReactionAdd(m.ChannelID, msg.ID, pollEmojis[idx])
	}
	return nil
}

func quickpollMessageCommandHandler(s *discordgo.Session, m *discordgo.MessageCreate, args []string) error {
	if len(args) == 0 {
		if _, err := s.ChannelMessageSend(m.ChannelID, "Usage: `-quickpoll <question>`"); err != nil {
			return err
		}
		return fmt.Errorf("quickpoll requires a question")
	}

	question := strings.TrimSpace(strings.Join(args, " "))
	ctx := &PollContext{Question: question, Options: []string{"Yes", "No"}, Votes: []int{0, 0}, IsQuick: true}

	msg, err := s.ChannelMessageSendEmbed(m.ChannelID, BuildPollEmbed(question, ctx.Options, ctx.Votes))
	if err != nil {
		return err
	}

	pollManagerSingleton.register(s, msg.ID, m.ChannelID, ctx)
	_ = s.MessageReactionAdd(m.ChannelID, msg.ID, "👍")
	_ = s.MessageReactionAdd(m.ChannelID, msg.ID, "👎")
	return nil
}

// register stores a poll and registers the shared reaction handler once.
func (pm *pollManager) register(s *discordgo.Session, messageID, channelID string, ctx *PollContext) {
	pm.mu.Lock()
	pm.polls[messageID] = ctx
	pm.channels[messageID] = channelID
	pm.mu.Unlock()

	pm.once.Do(func() {
		s.AddHandler(pm.onReactionAdd)
		s.AddHandler(pm.onReactionRemove)
	})
}

// onReactionAdd increments the matching option's vote count.
func (pm *pollManager) onReactionAdd(s *discordgo.Session, r *discordgo.MessageReactionAdd) {
	if r.UserID == s.State.User.ID {
		return
	}
	pm.mu.Lock()
	ctx, ok := pm.polls[r.MessageID]
	if !ok || ctx.Closed {
		pm.mu.Unlock()
		return
	}
	idx := pollEmojiIndex(ctx, r.Emoji.Name)
	if idx >= 0 && idx < len(ctx.Votes) {
		ctx.Votes[idx]++
	}
	chid := pm.channels[r.MessageID]
	pm.mu.Unlock()

	if idx >= 0 {
		_, _ = s.ChannelMessageEditEmbed(chid, r.MessageID, BuildPollEmbed(ctx.Question, ctx.Options, ctx.Votes))
	}
}

// onReactionRemove decrements the matching option's vote count.
func (pm *pollManager) onReactionRemove(s *discordgo.Session, r *discordgo.MessageReactionRemove) {
	if r.UserID == s.State.User.ID {
		return
	}
	pm.mu.Lock()
	ctx, ok := pm.polls[r.MessageID]
	if !ok || ctx.Closed {
		pm.mu.Unlock()
		return
	}
	idx := pollEmojiIndex(ctx, r.Emoji.Name)
	if idx >= 0 && idx < len(ctx.Votes) && ctx.Votes[idx] > 0 {
		ctx.Votes[idx]--
	}
	chid := pm.channels[r.MessageID]
	pm.mu.Unlock()

	if idx >= 0 {
		_, _ = s.ChannelMessageEditEmbed(chid, r.MessageID, BuildPollEmbed(ctx.Question, ctx.Options, ctx.Votes))
	}
}

// pollEmojiIndex maps a reaction emoji to an option index, or -1 if the emoji
// isn't a valid poll option.
func pollEmojiIndex(ctx *PollContext, emoji string) int {
	if ctx.IsQuick {
		if emoji == "👍" {
			return 0
		}
		if emoji == "👎" {
			return 1
		}
		return -1
	}
	for i := 0; i < len(ctx.Options) && i < maxPollOptions; i++ {
		if pollEmojis[i] == emoji {
			return i
		}
	}
	return -1
}

func pollSlashCommandHandler(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	question := ""
	options := make([]string, 0)
	for _, opt := range i.ApplicationCommandData().Options {
		switch opt.Name {
		case "question":
			question = opt.StringValue()
		case "options":
			for _, part := range strings.Split(opt.StringValue(), "|") {
				if t := strings.TrimSpace(part); t != "" {
					options = append(options, t)
				}
			}
		}
	}

	question = strings.TrimSpace(question)
	if question == "" {
		return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{Content: "Usage: `/poll question:<question> options:<opts>`"},
		})
	}
	if len(options) < 1 || len(options) > maxPollOptions {
		return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{Content: "Please provide 1-10 options."},
		})
	}

	ctx := &PollContext{Question: question, Options: options, Votes: make([]int, len(options))}

	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{BuildPollEmbed(question, options, ctx.Votes)},
		},
	})
	if err != nil {
		return err
	}

	msg, err := s.InteractionResponse(i.Interaction)
	if err != nil {
		return err
	}
	pollManagerSingleton.register(s, msg.ID, i.ChannelID, ctx)
	for idx := 0; idx < len(options); idx++ {
		_ = s.MessageReactionAdd(i.ChannelID, msg.ID, pollEmojis[idx])
	}
	return nil
}

func quickpollSlashCommandHandler(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	question := ""
	for _, opt := range i.ApplicationCommandData().Options {
		if opt.Name == "question" {
			question = opt.StringValue()
		}
	}
	question = strings.TrimSpace(question)
	if question == "" {
		return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{Content: "Usage: `/quickpoll question:<question>`"},
		})
	}

	ctx := &PollContext{Question: question, Options: []string{"Yes", "No"}, Votes: []int{0, 0}, IsQuick: true}
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{BuildPollEmbed(question, ctx.Options, ctx.Votes)},
		},
	})
	if err != nil {
		return err
	}

	msg, err := s.InteractionResponse(i.Interaction)
	if err != nil {
		return err
	}
	pollManagerSingleton.register(s, msg.ID, i.ChannelID, ctx)
	_ = s.MessageReactionAdd(i.ChannelID, msg.ID, "👍")
	_ = s.MessageReactionAdd(i.ChannelID, msg.ID, "👎")
	return nil
}
