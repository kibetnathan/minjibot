package commands

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

// rpEmotions maps each emotion command to the text shown in the embed title.
var rpEmotions = map[string]string{
	"angry":     "angry",
	"depressed": "depressed",
	"excited":   "excited",
	"happy":     "happy",
	"horny":     "horny",
	"inlove":    "in love",
	"sad":       "sad",
	"shy":       "shy",
}

// rpActions maps each action command to the past-tense text used in the embed.
var rpActions = map[string]string{
	"baka":     "calls",
	"bite":     "bites",
	"cry":      "cries at",
	"dap":      "daps up",
	"eat":      "munches on",
	"facepalm": "facepalms at",
	"feed":     "feeds",
	"handhold": "holds hands with",
	"kiss":     "kisses",
	"laugh":    "laughs at",
	"nod":      "nods at",
	"nutkick":  "nutkicks",
	"pat":      "pats",
	"peck":     "pecks",
	"poke":     "pokes",
	"punch":    "punches",
	"run":      "runs away from",
	"shoot":    "shoots",
	"shrug":    "shrugs at",
	"slap":     "slaps",
	"spank":    "spanks",
	"stab":     "stabs",
	"think":    "thinks about",
	"tickle":   "tickles",
}

// rpEmotionMessage handles a prefix emotion command (e.g. -angry).
func rpEmotionMessage(s *discordgo.Session, m *discordgo.MessageCreate, cmd string) error {
	label, ok := rpEmotions[cmd]
	if !ok {
		return fmt.Errorf("unknown emotion: %s", cmd)
	}
	embed, err := buildRPEmbed(s, "anime "+label, fmt.Sprintf("**%s** is feeling %s", displayName(m), label), 0xF58CBA)
	if err != nil {
		return err
	}
	_, err = s.ChannelMessageSendEmbed(m.ChannelID, embed)
	return err
}

// rpActionMessage handles a prefix action command (e.g. -slap @user).
func rpActionMessage(s *discordgo.Session, m *discordgo.MessageCreate, args []string, cmd string) error {
	verb, ok := rpActions[cmd]
	if !ok {
		return fmt.Errorf("unknown action: %s", cmd)
	}
	target := m.Author.Username
	if m.Author.Username == "" {
		target = "themselves"
	}
	if len(args) > 0 {
		if id, name, err := resolveTargetUser(s, m.GuildID, args[0]); err == nil {
			_ = id
			target = name
		}
	}
	embed, err := buildRPEmbed(s, "anime "+cmd, fmt.Sprintf("%s %s **%s**", displayName(m), verb, target), 0x5865F2)
	if err != nil {
		return err
	}
	_, err = s.ChannelMessageSendEmbed(m.ChannelID, embed)
	return err
}

// rpEmotionSlash handles a slash emotion command.
func rpEmotionSlash(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	cmd := i.ApplicationCommandData().Name
	label, ok := rpEmotions[cmd]
	if !ok {
		label = cmd
	}
	embed, err := buildRPEmbed(s, "anime "+label, fmt.Sprintf("**%s** is feeling %s", interactionDisplayName(i), label), 0xF58CBA)
	if err != nil {
		return err
	}
	return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{Embeds: []*discordgo.MessageEmbed{embed}},
	})
}

// rpActionSlash handles a slash action command.
func rpActionSlash(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	cmd := i.ApplicationCommandData().Name
	verb, ok := rpActions[cmd]
	if !ok {
		verb = cmd
	}
	target := "themselves"
	if len(i.ApplicationCommandData().Options) > 0 {
		raw := i.ApplicationCommandData().Options[0].StringValue()
		if id, name, err := resolveTargetUser(s, i.GuildID, raw); err == nil {
			_ = id
			target = name
		}
	}
	embed, err := buildRPEmbed(s, "anime "+cmd, fmt.Sprintf("%s %s **%s**", interactionDisplayName(i), verb, target), 0x5865F2)
	if err != nil {
		return err
	}
	return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{Embeds: []*discordgo.MessageEmbed{embed}},
	})
}

// buildRPEmbed fetches the first matching Giphy result and wraps it in an embed.
func buildRPEmbed(s *discordgo.Session, query, description string, color int) (*discordgo.MessageEmbed, error) {
	results, err := newGiphyClient().search(query, 8)
	if err != nil || len(results) == 0 {
		return &discordgo.MessageEmbed{
			Color:       color,
			Description: description,
		}, nil
	}
	return &discordgo.MessageEmbed{
		Color:       color,
		Description: description,
		Image:       &discordgo.MessageEmbedImage{URL: results[0].GIFURL()},
		Footer:      &discordgo.MessageEmbedFooter{Text: "GIF via Giphy"},
	}, nil
}

func displayName(m *discordgo.MessageCreate) string {
	if m.Author != nil {
		if m.Member != nil && m.Member.Nick != "" {
			return m.Member.Nick
		}
		return m.Author.Username
	}
	return "Someone"
}

func interactionDisplayName(i *discordgo.InteractionCreate) string {
	if i.Member != nil {
		if i.Member.Nick != "" {
			return i.Member.Nick
		}
		if i.Member.User != nil {
			return i.Member.User.Username
		}
	}
	if i.User != nil {
		return i.User.Username
	}
	return "Someone"
}

func rpUserOption() *discordgo.ApplicationCommandOption {
	return &discordgo.ApplicationCommandOption{
		Type:        discordgo.ApplicationCommandOptionUser,
		Name:        "user",
		Description: "The user to direct the action at",
		Required:    false,
	}
}

var emotionSlashCommands = func() []*discordgo.ApplicationCommand {
	out := make([]*discordgo.ApplicationCommand, 0, len(rpEmotions))
	for name := range rpEmotions {
		out = append(out, &discordgo.ApplicationCommand{
			Name:        name,
			Description: fmt.Sprintf("Post a GIF expressing being %s", name),
		})
	}
	return out
}()

var actionSlashCommands = func() []*discordgo.ApplicationCommand {
	out := make([]*discordgo.ApplicationCommand, 0, len(rpActions))
	for name := range rpActions {
		out = append(out, &discordgo.ApplicationCommand{
			Name:        name,
			Description: fmt.Sprintf("Perform %s on a user", name),
			Options:     []*discordgo.ApplicationCommandOption{rpUserOption()},
		})
	}
	return out
}()
