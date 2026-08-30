package commands

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/kibetnathan/minjibot/internal/ports/dto"
)

func diaryMessageCommandHandler(h *CommandHandler, s *discordgo.Session, m *discordgo.MessageCreate, args []string) error {
	if len(args) == 0 {
		if _, err := s.ChannelMessageSend(m.ChannelID, "Usage: `-diary add <text>` | `-diary view` | `-diary delete <id>`"); err != nil {
			return err
		}
		return fmt.Errorf("diary requires a subcommand")
	}

	sub := strings.ToLower(args[0])
	switch sub {
	case "add":
		if len(args) < 2 {
			if _, err := s.ChannelMessageSend(m.ChannelID, "Usage: `-diary add <text>`"); err != nil {
				return err
			}
			return nil
		}
		content := strings.TrimSpace(strings.Join(args[1:], " "))
		_, err := h.DiaryRepo.Create(context.Background(), dto.CreateDiaryEntryParams{
			UserID:  m.Author.ID,
			Content: content,
		})
		if err != nil {
			return err
		}
		if _, err := s.ChannelMessageSend(m.ChannelID, "Diary entry saved. Use `-diary view` to see your entries (they're private)."); err != nil {
			return err
		}
		return nil
	case "view":
		return h.diaryViewMessage(s, m)
	case "delete", "del", "remove":
		if len(args) < 2 {
			if _, err := s.ChannelMessageSend(m.ChannelID, "Usage: `-diary delete <id>`"); err != nil {
				return err
			}
			return nil
		}
		id, err := strconv.ParseInt(args[1], 10, 64)
		if err != nil {
			if _, serr := s.ChannelMessageSend(m.ChannelID, "That doesn't look like a valid entry ID."); serr != nil {
				return serr
			}
			return nil
		}
		if err := h.DiaryRepo.Delete(context.Background(), id, m.Author.ID); err != nil {
			return err
		}
		if _, err := s.ChannelMessageSend(m.ChannelID, "Diary entry deleted."); err != nil {
			return err
		}
		return nil
	default:
		if _, err := s.ChannelMessageSend(m.ChannelID, "Usage: `-diary add <text>` | `-diary view` | `-diary delete <id>`"); err != nil {
			return err
		}
		return nil
	}
}

// sendPrivateEntries DMs the user their diary entries. Falls back to posting a
// visible warning if DMs are blocked.
func (h *CommandHandler) sendPrivateEntries(s *discordgo.Session, userID, channelID string, entries []diaryEntryView) error {
	var b strings.Builder
	if len(entries) == 0 {
		b.WriteString("Your diary is empty. Use `-diary add <text>` to save your first entry.")
	} else {
		for _, e := range entries {
			b.WriteString(fmt.Sprintf("**#%d** (%s)\n%s\n\n", e.ID, e.Created.Format("Jan 2, 2006 15:04"), e.Content))
		}
	}

	dm, err := s.UserChannelCreate(userID)
	if err == nil {
		if _, err := s.ChannelMessageSend(dm.ID, b.String()); err == nil {
			_, _ = s.ChannelMessageSend(channelID, fmt.Sprintf("<@%s>, your diary entries have been sent to you in a DM (they're private).", userID))
			return nil
		}
	}
	// DMs unavailable — don't leak the private content publicly.
	_, serr := s.ChannelMessageSend(channelID, fmt.Sprintf("<@%s>, I couldn't DM you your diary (DMs may be closed). I won't post it here for privacy.", userID))
	return serr
}

type diaryEntryView struct {
	ID      int64
	Content string
	Created time.Time
}

func (h *CommandHandler) diaryViewMessage(s *discordgo.Session, m *discordgo.MessageCreate) error {
	rowEntries, err := h.DiaryRepo.ListByUser(context.Background(), m.Author.ID)
	if err != nil {
		return err
	}
	entries := make([]diaryEntryView, len(rowEntries))
	for i, e := range rowEntries {
		entries[i] = diaryEntryView{ID: e.ID, Content: e.Content, Created: e.CreatedAt}
	}
	return h.sendPrivateEntries(s, m.Author.ID, m.ChannelID, entries)
}

func diarySlashCommandHandler(h *CommandHandler, s *discordgo.Session, i *discordgo.InteractionCreate) error {
	data := i.ApplicationCommandData()
	sub := data.Options
	if len(sub) == 0 {
		return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{Content: "Usage: `/diary add|view|delete`"},
		})
	}
	opt := sub[0]
	args := opt.Options

	respond := func(content string) error {
		return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{Content: content},
		})
	}

	switch opt.Name {
	case "add":
		content := ""
		for _, a := range args {
			if a.Name == "text" {
				content = a.StringValue()
			}
		}
		content = strings.TrimSpace(content)
		if content == "" {
			return respond("Usage: `/diary add text:<text>`")
		}
		if _, err := h.DiaryRepo.Create(context.Background(), dto.CreateDiaryEntryParams{
			UserID: i.Member.User.ID, Content: content,
		}); err != nil {
			return err
		}
		return respond("Diary entry saved. Use `/diary view` to see your entries (they're private).")
	case "view":
		rowEntries, err := h.DiaryRepo.ListByUser(context.Background(), i.Member.User.ID)
		if err != nil {
			return err
		}
		if len(rowEntries) == 0 {
			return respond("Your diary is empty. Use `/diary add text:<text>` to save your first entry.")
		}
		var b strings.Builder
		for _, e := range rowEntries {
			b.WriteString(fmt.Sprintf("**#%d** (%s)\n%s\n\n", e.ID, e.CreatedAt.Format("Jan 2, 2006 15:04"), e.Content))
		}
		if dm, err := s.UserChannelCreate(i.Member.User.ID); err == nil {
			if _, err := s.ChannelMessageSend(dm.ID, b.String()); err == nil {
				return respond("Your diary entries have been sent to you in a DM (they're private).")
			}
		}
		return respond("I couldn't DM you your diary (DMs may be closed). I won't post it here for privacy.")
	case "delete", "del", "remove":
		id := int64(0)
		for _, a := range args {
			if a.Name == "id" {
				id = a.IntValue()
			}
		}
		if id == 0 {
			return respond("Usage: `/diary delete id:<entry-id>`")
		}
		if err := h.DiaryRepo.Delete(context.Background(), id, i.Member.User.ID); err != nil {
			return respond("Couldn't delete that entry (it may not exist or belong to you).")
		}
		return respond("Diary entry deleted.")
	default:
		return respond("Usage: `/diary add|view|delete`")
	}
}
