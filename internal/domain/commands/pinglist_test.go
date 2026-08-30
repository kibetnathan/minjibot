package commands

import (
	"strings"
	"testing"
	"time"

	"github.com/bwmarrin/discordgo"
)

func TestSnowflakeFromValue(t *testing.T) {
	cases := []struct {
		name string
		in   any
		want string
	}{
		{"string", "123456", "123456"},
		{"int id via float", float64(123456789), "123456789"},
		{"rounds float", float64(123.9), "124"},
		{"nil", nil, ""},
		{"int", 5, ""},
		{"zero", float64(0), "0"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := snowflakeFromValue(tc.in); got != tc.want {
				t.Errorf("snowflakeFromValue(%v) = %q, want %q", tc.in, got, tc.want)
			}
		})
	}
}

func TestBuildPinglistEmbedNoMatches(t *testing.T) {
	embed := buildPinglistEmbed(&pingTarget{Kind: "user", ID: "123"}, nil)
	if !strings.Contains(embed.Description, "No pings found.") {
		t.Errorf("description = %q", embed.Description)
	}
	// Name falls back to ID.
	if !strings.Contains(embed.Title, "123") {
		t.Errorf("title = %q", embed.Title)
	}
	if embed.Footer != nil {
		t.Errorf("unexpected footer: %+v", embed.Footer)
	}
}

func TestBuildPinglistEmbedMatches(t *testing.T) {
	target := &pingTarget{Kind: "role", ID: "999", Name: "admins"}
	msgs := []*discordgo.Message{
		{ID: "m1", GuildID: "", ChannelID: "c", Author: &discordgo.User{Username: "alice"}, Content: "hey <@&999> here",
			Timestamp: time.Date(2024, 1, 2, 3, 4, 5, 0, time.UTC)},
	}
	embed := buildPinglistEmbed(target, msgs)
	if !strings.Contains(embed.Title, "admins") {
		t.Errorf("title should prefer name: %q", embed.Title)
	}
	if len(embed.Fields) != 1 {
		t.Fatalf("len(fields) = %d", len(embed.Fields))
	}
	if !strings.Contains(embed.Fields[0].Value, "https://discord.com/channels/@me/c/m1") {
		t.Errorf("field value = %q", embed.Fields[0].Value)
	}
	if !strings.HasPrefix(embed.Fields[0].Name, "alice •") {
		t.Errorf("field name = %q", embed.Fields[0].Name)
	}
	if embed.Footer == nil || !strings.Contains(embed.Footer.Text, "1 ping") {
		t.Errorf("footer = %+v", embed.Footer)
	}
}

func TestBuildPinglistEmbedTruncates(t *testing.T) {
	var msgs []*discordgo.Message
	for i := 0; i < pinglistMaxResults+5; i++ {
		msgs = append(msgs, &discordgo.Message{
			ID: "m", GuildID: "g", ChannelID: "c",
			Author:  &discordgo.User{Username: "u"},
			Content: "x",
		})
	}
	embed := buildPinglistEmbed(&pingTarget{Kind: "user", ID: "1"}, msgs)
	if len(embed.Fields) != pinglistMaxResults {
		t.Errorf("len(fields) = %d, want %d", len(embed.Fields), pinglistMaxResults)
	}
}
