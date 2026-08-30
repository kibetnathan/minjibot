package tests

import (
	"strings"
	"testing"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/kibetnathan/minjibot/internal/commands"
)

func TestParseMessageLink(t *testing.T) {
	cases := []struct {
		name        string
		in          string
		wantGuild   string
		wantChannel string
		wantMessage string
		wantOK      bool
	}{
		{"full link", "https://discord.com/channels/111/222/333", "111", "222", "333", true},
		{"app domain", "https://discordapp.com/channels/111/222/333", "111", "222", "333", true},
		{"dm at-me", "https://discord.com/channels/@me/222/333", "@me", "222", "333", true},
		{"embedded in text", "see https://discord.com/channels/111/222/333 ok", "111", "222", "333", true},
		{"not a link", "just some text", "", "", "", false},
		{"wrong domain", "https://example.com/channels/1/2/3", "", "", "", false},
		{"missing group", "https://discord.com/channels/111/222", "", "", "", false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			g, c, m, ok := commands.ParseMessageLink(tc.in)
			if g != tc.wantGuild || c != tc.wantChannel || m != tc.wantMessage || ok != tc.wantOK {
				t.Errorf("commands.ParseMessageLink(%q) = (%q,%q,%q,%v), want (%q,%q,%q,%v)",
					tc.in, g, c, m, ok, tc.wantGuild, tc.wantChannel, tc.wantMessage, tc.wantOK)
			}
		})
	}
}

func TestMessageToLink(t *testing.T) {
	withGuild := &discordgo.Message{GuildID: "g", ChannelID: "c", ID: "m"}
	if got := commands.MessageToLink(withGuild); got != "https://discord.com/channels/g/c/m" {
		t.Errorf("got %q", got)
	}
	dm := &discordgo.Message{ChannelID: "c", ID: "m"}
	if got := commands.MessageToLink(dm); got != "https://discord.com/channels/@me/c/m" {
		t.Errorf("got %q", got)
	}
}

func TestQuoteEmbed(t *testing.T) {
	msg := &discordgo.Message{
		ID:          "m",
		ChannelID:   "c",
		Author:      &discordgo.User{ID: "u", Username: "name", Avatar: "hash"},
		Content:     "the content",
		Timestamp:   time.Date(2024, 3, 4, 5, 6, 7, 0, time.UTC),
		Attachments: []*discordgo.MessageAttachment{{ID: "a1", URL: "https://cdn/att.png"}},
	}
	embed := commands.QuoteEmbed(msg)
	if embed.Author == nil || embed.Author.Name != "name" {
		t.Errorf("author = %+v", embed.Author)
	}
	if embed.Author.IconURL != "https://cdn.discordapp.com/avatars/u/hash.png?size=64" {
		t.Errorf("icon = %q", embed.Author.IconURL)
	}
	if embed.Description != "the content" {
		t.Errorf("description = %q", embed.Description)
	}
	if embed.Timestamp != "2024-03-04T05:06:07+00:00" {
		t.Errorf("timestamp = %q", embed.Timestamp)
	}
	if embed.Image == nil || embed.Image.URL != "https://cdn/att.png" {
		t.Errorf("image = %+v", embed.Image)
	}
	if embed.Footer == nil || !strings.Contains(embed.Footer.Text, "https://discord.com/channels/@me/c/m") {
		t.Errorf("footer = %+v", embed.Footer)
	}
}

func TestQuoteEmbedEmbedOverridesAttachment(t *testing.T) {
	msg := &discordgo.Message{
		Author:      &discordgo.User{Username: "name"},
		Content:     "x",
		Attachments: []*discordgo.MessageAttachment{{ID: "a1", URL: "https://cdn/att.png"}},
		Embeds:      []*discordgo.MessageEmbed{{Image: &discordgo.MessageEmbedImage{URL: "https://cdn/embed.png"}}},
	}
	embed := commands.QuoteEmbed(msg)
	if embed.Image == nil || embed.Image.URL != "https://cdn/embed.png" {
		t.Errorf("embed image should win, got %+v", embed.Image)
	}
	if !strings.HasPrefix(embed.Timestamp, "0001-01-01T00:00:00") {
		t.Errorf("zero timestamp should still format, got %q", embed.Timestamp)
	}
}

func TestQuoteEmbedNilAuthor(t *testing.T) {
	msg := &discordgo.Message{Content: "lonely"}
	embed := commands.QuoteEmbed(msg)
	if embed.Author == nil || embed.Author.Name != "Unknown" || embed.Author.IconURL != "" {
		t.Errorf("author = %+v", embed.Author)
	}
}
