package commands

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

const urbanDictionaryAPI = "https://api.urbandictionary.com/v0/define"

type urbanDictionaryResponse struct {
	List []struct {
		Definition string `json:"definition"`
		Word       string `json:"word"`
		Example    string `json:"example"`
		Permalink  string `json:"permalink"`
		ThumbsUp   int    `json:"thumbs_up"`
		ThumbsDown int    `json:"thumbs_down"`
	} `json:"list"`
}

func urbandictionaryMessageCommandHandler(s *discordgo.Session, m *discordgo.MessageCreate, args []string) error {
	if len(args) == 0 {
		_, err := s.ChannelMessageSend(m.ChannelID, "Usage: `-urbandictionary <term>`")
		return err
	}
	term := strings.Join(args, " ")
	embed, err := buildUrbanDictionaryEmbed(term)
	if err != nil {
		return err
	}
	_, err = s.ChannelMessageSendEmbed(m.ChannelID, embed)
	return err
}

func urbandictionarySlashCommandHandler(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	opts := OptionMap(i.ApplicationCommandData().Options)
	term := strings.TrimSpace(OptString(opts, "term"))
	if term == "" {
		return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{Content: "Usage: `/urbandictionary term:<term>`"},
		})
	}
	embed, err := buildUrbanDictionaryEmbed(term)
	if err != nil {
		return err
	}
	return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{embed},
		},
	})
}

func buildUrbanDictionaryEmbed(term string) (*discordgo.MessageEmbed, error) {
	endpoint := fmt.Sprintf("%s?term=%s", urbanDictionaryAPI, url.QueryEscape(term))
	client := &http.Client{Timeout: 15 * time.Second}

	resp, err := client.Get(endpoint)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return &discordgo.MessageEmbed{
			Color:       0xED4245,
			Title:       "Urban Dictionary",
			Description: fmt.Sprintf("Urban Dictionary returned status %d. Try again in a bit.", resp.StatusCode),
		}, nil
	}

	var out urbanDictionaryResponse
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, err
	}
	if len(out.List) == 0 {
		return &discordgo.MessageEmbed{
			Color:       0xED4245,
			Title:       "Urban Dictionary",
			Description: fmt.Sprintf("No definitions found for `%s`.", term),
		}, nil
	}

	def := out.List[0]
	return &discordgo.MessageEmbed{
		Color: 0x5865F2,
		Title: def.Word,
		URL:   def.Permalink,
		Fields: []*discordgo.MessageEmbedField{
			{Name: "Definition", Value: truncateEmbedText(def.Definition, "*no definition given*")},
			{Name: "Example", Value: truncateEmbedText(def.Example, "*(no example given)*")},
		},
		Footer: &discordgo.MessageEmbedFooter{Text: fmt.Sprintf("👍 %d    👎 %d", def.ThumbsUp, def.ThumbsDown)},
	}, nil
}

// truncateEmbedText trims s and caps it at the embed field limit (1024),
// returning fallback when the result would be empty.
func truncateEmbedText(s, fallback string) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return fallback
	}
	if len(s) > 1024 {
		s = s[:1023] + "…"
	}
	return s
}
