package commands

import (
	"sync"
	"time"

	"github.com/bwmarrin/discordgo"
)

const (
	paginationPrevEmoji = "◀️"
	paginationNextEmoji = "▶️"
	paginationTTL       = 2 * time.Minute
)

// paginationSession tracks an active reaction-paginated message.
type paginationSession struct {
	channelID string
	authorID  string
	page      int
	total     int
	pageEmbed func(int) *discordgo.MessageEmbed
}

// reactionPaginator manages active reaction paginations, keyed by message ID.
type reactionPaginator struct {
	mu       sync.Mutex
	sessions map[string]*paginationSession
	once     sync.Once
}

// paginator is the shared singleton used by paginated commands.
var paginator = &reactionPaginator{sessions: make(map[string]*paginationSession)}

// paginateReactions sends pageEmbed(0) and lets author flip pages with
// reactions within a timeout. It is fire-and-forget: it returns once the
// message is sent and reactions are added.
func paginateReactions(s *discordgo.Session, channelID, authorID string, total int, pageEmbed func(int) *discordgo.MessageEmbed) error {
	if total <= 0 {
		return nil
	}
	msg, err := s.ChannelMessageSendEmbed(channelID, pageEmbed(0))
	if err != nil {
		return err
	}

	sess := &paginationSession{
		channelID: channelID,
		authorID:  authorID,
		page:      0,
		total:     total,
		pageEmbed: pageEmbed,
	}

	paginator.mu.Lock()
	paginator.sessions[msg.ID] = sess
	paginator.mu.Unlock()

	paginator.once.Do(func() { s.AddHandler(paginator.onReaction) })

	if total > 1 {
		_ = s.MessageReactionAdd(channelID, msg.ID, paginationPrevEmoji)
		_ = s.MessageReactionAdd(channelID, msg.ID, paginationNextEmoji)
	}
	return nil
}

func (p *reactionPaginator) onReaction(s *discordgo.Session, r *discordgo.MessageReactionAdd) {
	if r.UserID == s.State.User.ID {
		return
	}

	p.mu.Lock()
	sess, ok := p.sessions[r.MessageID]
	if !ok {
		p.mu.Unlock()
		return
	}
	if sess.authorID != "" && r.UserID != sess.authorID {
		p.mu.Unlock()
		return
	}
	page := sess.page
	switch r.Emoji.Name {
	case paginationPrevEmoji:
		page--
	case paginationNextEmoji:
		page++
	default:
		p.mu.Unlock()
		return
	}
	if page < 0 || page >= sess.total {
		p.mu.Unlock()
		return
	}
	sess.page = page
	p.mu.Unlock()

	_, _ = s.ChannelMessageEditEmbed(r.ChannelID, r.MessageID, sess.pageEmbed(page))
}

const helpPrevCustomID = "help:page:prev"
const helpNextCustomID = "help:page:next"

// helpButtonsRow returns the prev/next navigation buttons for help. It is shared
// by both the slash and prefix help commands so they navigate identically.
func helpButtonsRow() []discordgo.MessageComponent {
	return []discordgo.MessageComponent{
		discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				discordgo.Button{Label: "◀", Style: discordgo.SecondaryButton, CustomID: helpPrevCustomID},
				discordgo.Button{Label: "▶", Style: discordgo.SecondaryButton, CustomID: helpNextCustomID},
			},
		},
	}
}

// helpPageState tracks the current page of a slash-help message so button
// clicks can advance it. Keyed by message ID.
var (
	helpPageMu    sync.Mutex
	helpPageByMsg = map[string]int{}
)

// HelpButtonHandler processes clicks on the help pagination buttons. It returns
// nil (no error) for non-help components so callers can ignore unrelated ones.
func HelpButtonHandler(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	cd := i.MessageComponentData().CustomID
	if cd != helpPrevCustomID && cd != helpNextCustomID {
		return nil
	}
	msgID := i.Message.ID

	helpPageMu.Lock()
	page := helpPageByMsg[msgID]
	helpPageMu.Unlock()

	switch cd {
	case helpPrevCustomID:
		page--
	case helpNextCustomID:
		page++
	}
	total := NumHelpPages()
	if page < 0 {
		page = 0
	}
	if page >= total {
		page = total - 1
	}

	helpPageMu.Lock()
	helpPageByMsg[msgID] = page
	helpPageMu.Unlock()

	components := i.Message.Components
	return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseUpdateMessage,
		Data: &discordgo.InteractionResponseData{
			Embeds:     []*discordgo.MessageEmbed{BuildHelpPageEmbed(page)},
			Components: components,
		},
	})
}
