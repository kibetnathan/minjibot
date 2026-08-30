package tests

import (
	"strings"
	"testing"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/kibetnathan/minjibot/internal/commands"
)

func TestParseSearchArgs(t *testing.T) {
	cases := []struct {
		name    string
		args    []string
		wantQ   string
		wantN   int
		wantErr bool
	}{
		{"simple", []string{"hello", "world"}, "hello world", commands.SearchDefaultLimit, false},
		{"trailing number", []string{"hello", "50"}, "hello", 50, false},
		{"messages flag", []string{"messages:75", "hello"}, "hello", 75, false},
		{"clamp low", []string{"hello", "0"}, "hello", commands.SearchDefaultLimit, false},
		{"clamp high", []string{"messages:9999", "hi"}, "hi", commands.SearchDefaultLimit, false},
		{"bad messages flag", []string{"messages:abc", "hi"}, "hi", commands.SearchDefaultLimit, false},
		{"bare number only", []string{"50"}, "", 50, false},
		{"number mid query stays", []string{"version", "3", "released"}, "version released", 3, false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			q, n, err := commands.ParseSearchArgs(&discordgo.Session{}, "", tc.args)
			if (err != nil) != tc.wantErr {
				t.Fatalf("err = %v, wantErr %v", err, tc.wantErr)
			}
			if q != tc.wantQ {
				t.Errorf("query = %q, want %q", q, tc.wantQ)
			}
			if n != tc.wantN {
				t.Errorf("limit = %d, want %d", n, tc.wantN)
			}
		})
	}
}

func TestBuildSearchEmbedNoMatches(t *testing.T) {
	embed := commands.BuildSearchEmbed("xyzzy", 100, nil)
	if !strings.Contains(embed.Description, "No matches found.") {
		t.Errorf("expected no-matches hint in description, got %q", embed.Description)
	}
	if embed.Footer != nil {
		t.Errorf("unexpected footer for no matches: %+v", embed.Footer)
	}
}

func TestBuildSearchEmbedMatches(t *testing.T) {
	long := strings.Repeat("x", commands.SearchMaxFieldLength+100)
	match := func(id, author string, bot bool, content string) *discordgo.Message {
		return &discordgo.Message{
			ID:        id,
			GuildID:   "guild1",
			ChannelID: "chan1",
			Author:    &discordgo.User{Username: author, Bot: bot},
			Content:   content,
			Timestamp: time.Date(2024, 3, 4, 5, 6, 7, 0, time.UTC),
		}
	}
	msgs := []*discordgo.Message{
		match("1", "Alice", false, "found one"),
		match("2", "Bob", true, long),
		match("3", "Carol", false, "found three"),
	}

	embed := commands.BuildSearchEmbed("found", 200, msgs)
	if embed.Footer == nil || !strings.Contains(embed.Footer.Text, "3 match") {
		t.Errorf("unexpected footer: %+v", embed.Footer)
	}
	if len(embed.Fields) != 3 {
		t.Fatalf("len(fields) = %d, want 3", len(embed.Fields))
	}
	if !strings.HasPrefix(embed.Fields[0].Name, "Alice • 2024-03-04T05:06:07Z") {
		t.Errorf("field name = %q", embed.Fields[0].Name)
	}
	if !strings.Contains(embed.Fields[1].Name, "Bob [bot]") {
		t.Errorf("bot author should be flagged, got %q", embed.Fields[1].Name)
	}
	contentPart := strings.SplitN(embed.Fields[1].Value, "\n", 2)[0]
	if len(contentPart) > commands.SearchMaxFieldLength || !strings.HasSuffix(contentPart, "...") {
		t.Errorf("content not truncated to %d: len=%d value=%q", commands.SearchMaxFieldLength, len(contentPart), contentPart[:min(len(contentPart), 80)])
	}
	if !strings.Contains(embed.Fields[0].Value, "https://discord.com/channels/guild1/chan1/1") {
		t.Errorf("missing guild link: %q", embed.Fields[0].Value)
	}
}

func TestBuildSearchEmbedTruncatesResults(t *testing.T) {
	var msgs []*discordgo.Message
	for i := 0; i < commands.SearchMaxResults+2; i++ {
		msgs = append(msgs, &discordgo.Message{
			ID:        "id",
			GuildID:   "",
			ChannelID: "chan1",
			Author:    &discordgo.User{Username: "u"},
			Content:   "hit",
		})
	}
	embed := commands.BuildSearchEmbed("hit", 200, msgs)
	if len(embed.Fields) != commands.SearchMaxResults {
		t.Errorf("len(fields) = %d, want %d", len(embed.Fields), commands.SearchMaxResults)
	}
	if !strings.Contains(embed.Fields[0].Value, "https://discord.com/channels/@me/chan1/id") {
		t.Errorf("missing @me link fallback: %q", embed.Fields[0].Value)
	}
}
