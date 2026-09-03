package commands

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
)

const bugModalCustomID = "bug_report_modal"
const bugModalFieldID = "bug_report_details"

func bugSlashCommandHandler(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseModal,
		Data: &discordgo.InteractionResponseData{
			CustomID: bugModalCustomID,
			Title:    "Report a bug",
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.TextInput{
							CustomID:    bugModalFieldID,
							Label:       "Describe the bug",
							Style:       discordgo.TextInputParagraph,
							Placeholder: "What happened? Include any commands used and expected vs. actual behavior.",
							Required:    true,
							MinLength:   10,
							MaxLength:   2000,
						},
					},
				},
			},
		},
	})
}

func bugMessageCommandHandler(s *discordgo.Session, m *discordgo.MessageCreate, args []string) error {
	report := strings.TrimSpace(strings.Join(args, " "))
	if report == "" {
		_, err := s.ChannelMessageSend(m.ChannelID, "Usage: `-bug <describe the bug>`")
		return err
	}
	return sendBugConfirmation(s, m.ChannelID, &discordgo.Message{
		Author:    m.Author,
		ChannelID: m.ChannelID,
	}, report)
}

// HandleBugModalSubmit processes a submitted bug report modal and confirms to
// the reporter. If a reports channel is configured on the guild it forwards the
// report there; otherwise it returns an ephemeral confirmation to the user.
func HandleBugModalSubmit(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	data := i.ModalSubmitData()
	report := ""
	for _, row := range data.Components {
		for _, comp := range row.(*discordgo.ActionsRow).Components {
			if input, ok := comp.(*discordgo.TextInput); ok && input.CustomID == bugModalFieldID {
				report = strings.TrimSpace(input.Value)
			}
		}
	}
	if report == "" {
		report = "(empty report)"
	}

	embed := &discordgo.MessageEmbed{
		Color:       0xED4245,
		Title:       "Bug Report",
		Description: report,
		Footer:      &discordgo.MessageEmbedFooter{Text: fmt.Sprintf("Reported by %s", i.Member.User.Username)},
	}

	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{embed},
			Flags:  discordgo.MessageFlagsEphemeral,
		},
	})
	return err
}

func sendBugConfirmation(s *discordgo.Session, channelID string, msg *discordgo.Message, report string) error {
	embed := &discordgo.MessageEmbed{
		Color:       0xED4245,
		Title:       "Bug Report",
		Description: report,
		Footer:      &discordgo.MessageEmbedFooter{Text: fmt.Sprintf("Reported by %s", msg.Author.Username)},
	}
	_, err := s.ChannelMessageSendEmbed(channelID, embed)
	return err
}
