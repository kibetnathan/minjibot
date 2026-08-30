package commands

import (
	"fmt"
	"math/rand/v2"
	"strings"

	"github.com/bwmarrin/discordgo"
)

// RandomReading returns a 0-100 "reading" for the fun percentage commands.
func RandomReading() int {
	return rand.IntN(101)
}

// RandomPP returns a randomized "pp length" in cm (range enforced for humour).
func RandomPP() int {
	return rand.IntN(21) // 0..20cm
}

// RandomPuh returns a randomized "puh tightness" percentage.
func RandomPuh() int {
	return rand.IntN(101)
}

// RandomIQ returns a randomized IQ score in a fun range.
func RandomIQ() int {
	return rand.IntN(111) + 50 // 50..160
}

// RandomBitches returns a randomized number of "bitches".
func RandomBitches() int {
	return rand.IntN(11) // 0..10
}

// RandomShipScore returns a 0-100 romance compatibility score.
func RandomShipScore() int {
	return rand.IntN(101)
}

// ParseChooseItems splits a raw "a, b, c" input into trimmed choices.
func ParseChooseItems(raw string) []string {
	parts := strings.Split(raw, ",")
	items := make([]string, 0, len(parts))
	for _, p := range parts {
		if t := strings.TrimSpace(p); t != "" {
			items = append(items, t)
		}
	}
	return items
}

// PickChoose returns a random choice from items, or "" if empty.
func PickChoose(items []string) string {
	if len(items) == 0 {
		return ""
	}
	return items[rand.IntN(len(items))]
}

// funTargetID extracts a user ID from a mention in args, defaulting to the
// message author when none is given.
func funTargetID(m *discordgo.MessageCreate, args []string) string {
	for _, arg := range args {
		if id := ParseMentionID(arg); id != "" {
			return id
		}
	}
	return m.Author.ID
}

// resolveMentionArg returns the first mention found in args, or "".
func funMentionArg(args []string) string {
	for _, arg := range args {
		if id := ParseMentionID(arg); id != "" {
			return id
		}
	}
	return ""
}

// readingMessage posts a "X is N% <label>" reading for a target user.
func readingMessage(s *discordgo.Session, m *discordgo.MessageCreate, args []string, label string) error {
	id := funTargetID(m, args)
	_, err := s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("<@%s> is **%d%%** %s", id, RandomReading(), label))
	return err
}

// readingSlash responds to a slash invocation with a reading.
func readingSlash(s *discordgo.Session, i *discordgo.InteractionCreate, label string) error {
	opts := OptionMap(i.ApplicationCommandData().Options)
	id := OptUserID(opts, "user", i)
	content := fmt.Sprintf("<@%s> is **%d%%** %s", id, RandomReading(), label)
	return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{Content: content},
	})
}

func howgayMessageCommandHandler(s *discordgo.Session, m *discordgo.MessageCreate, args []string) error {
	return readingMessage(s, m, args, "gay")
}
func howgaySlashCommandHandler(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	return readingSlash(s, i, "gay")
}

func howautismMessageCommandHandler(s *discordgo.Session, m *discordgo.MessageCreate, args []string) error {
	return readingMessage(s, m, args, "autistic")
}
func howautismSlashCommandHandler(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	return readingSlash(s, i, "autistic")
}

func howlesbianMessageCommandHandler(s *discordgo.Session, m *discordgo.MessageCreate, args []string) error {
	return readingMessage(s, m, args, "lesbian")
}
func howlesbianSlashCommandHandler(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	return readingSlash(s, i, "lesbian")
}

func howsimpMessageCommandHandler(s *discordgo.Session, m *discordgo.MessageCreate, args []string) error {
	return readingMessage(s, m, args, "a simp")
}
func howsimpSlashCommandHandler(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	return readingSlash(s, i, "a simp")
}

func ppMessageCommandHandler(s *discordgo.Session, m *discordgo.MessageCreate, args []string) error {
	id := funTargetID(m, args)
	cm := RandomPP()
	visual := strings.Repeat("=", cm) + "D"
	_, err := s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("<@%s>'s pp: `%s` **%dcm**", id, visual, cm))
	return err
}
func ppSlashCommandHandler(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	opts := OptionMap(i.ApplicationCommandData().Options)
	id := OptUserID(opts, "user", i)
	cm := RandomPP()
	visual := strings.Repeat("=", cm) + "D"
	content := fmt.Sprintf("<@%s>'s pp: `%s` **%dcm**", id, visual, cm)
	return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{Content: content},
	})
}

func puhMessageCommandHandler(s *discordgo.Session, m *discordgo.MessageCreate, args []string) error {
	_, err := s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Puh tightness: **%d%%**", RandomPuh()))
	return err
}
func puhSlashCommandHandler(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{Content: fmt.Sprintf("Puh tightness: **%d%%**", RandomPuh())},
	})
}

func iqMessageCommandHandler(s *discordgo.Session, m *discordgo.MessageCreate, args []string) error {
	id := funTargetID(m, args)
	_, err := s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("<@%s> has an IQ of **%d**", id, RandomIQ()))
	return err
}
func iqSlashCommandHandler(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	opts := OptionMap(i.ApplicationCommandData().Options)
	id := OptUserID(opts, "user", i)
	content := fmt.Sprintf("<@%s> has an IQ of **%d**", id, RandomIQ())
	return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{Content: content},
	})
}

func bitchesMessageCommandHandler(s *discordgo.Session, m *discordgo.MessageCreate, args []string) error {
	id := funTargetID(m, args)
	_, err := s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("<@%s> has **%d** bitches", id, RandomBitches()))
	return err
}
func bitchesSlashCommandHandler(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	opts := OptionMap(i.ApplicationCommandData().Options)
	id := OptUserID(opts, "user", i)
	content := fmt.Sprintf("<@%s> has **%d** bitches", id, RandomBitches())
	return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{Content: content},
	})
}

func chooseMessageCommandHandler(s *discordgo.Session, m *discordgo.MessageCreate, args []string) error {
	items := ParseChooseItems(strings.Join(args, " "))
	if len(items) == 0 {
		_, err := s.ChannelMessageSend(m.ChannelID, "Usage: `-choose option1, option2, option3`")
		return err
	}
	_, err := s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("I choose: **%s**", PickChoose(items)))
	return err
}
func chooseSlashCommandHandler(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	opts := OptionMap(i.ApplicationCommandData().Options)
	items := ParseChooseItems(OptString(opts, "choices"))
	if len(items) == 0 {
		return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{Content: "Usage: `/choose choices:option1, option2, option3`"},
		})
	}
	content := fmt.Sprintf("I choose: **%s**", PickChoose(items))
	return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{Content: content},
	})
}

func shipMessageCommandHandler(s *discordgo.Session, m *discordgo.MessageCreate, args []string) error {
	a, b := shipUsers(m, args)
	_, err := s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("❤️ <@%s> × <@%s> — compatibility **%d%%**", a, b, RandomShipScore()))
	return err
}
func shipSlashCommandHandler(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	opts := OptionMap(i.ApplicationCommandData().Options)
	a := OptUserID(opts, "user1", i)
	b := OptUserID(opts, "user2", i)
	content := fmt.Sprintf("❤️ <@%s> × <@%s> — compatibility **%d%%**", a, b, RandomShipScore())
	return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{Content: content},
	})
}

// shipUsers returns the two target user IDs for the ship command, defaulting
// to the caller when fewer than two mentions are provided.
func shipUsers(m *discordgo.MessageCreate, args []string) (string, string) {
	first := funMentionArg(args)
	if first == "" {
		a := m.Author.ID
		return a, a
	}
	b := ""
	for _, arg := range args {
		if id := ParseMentionID(arg); id != "" && id != first {
			b = id
			break
		}
	}
	if b == "" {
		b = m.Author.ID
	}
	return first, b
}

// OptUserID returns a user option's snowflake ID, defaulting to the invoking
// member's user ID when absent.
func OptUserID(opts map[string]*discordgo.ApplicationCommandInteractionDataOption, name string, i *discordgo.InteractionCreate) string {
	if o, ok := opts[name]; ok && o != nil {
		if id, ok := o.Value.(string); ok && id != "" {
			return id
		}
	}
	if i != nil && i.Member != nil && i.Member.User != nil {
		return i.Member.User.ID
	}
	return ""
}
