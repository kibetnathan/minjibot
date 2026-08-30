package commands

import (
	"testing"

	"github.com/bwmarrin/discordgo"
)

func TestExtensionFromContentType(t *testing.T) {
	cases := []struct {
		in   string
		want string
	}{
		{"image/png", "png"},
		{"image/png; charset=utf-8", "png"},
		{"image/jpeg", "jpg"},
		{"image/jpg", "jpg"},
		{"IMAGE/JPEG", "jpg"},
		{"image/gif", "gif"},
		{"image/webp", "webp"},
		{"video/mp4", "mp4"},
		{"video/webm", "webm"},
		{"image/avif", "avif"},
		{"text/html", "bin"},
		{"", "bin"},
		{"application/octet-stream", "bin"},
	}
	for _, tc := range cases {
		if got := extensionFromContentType(tc.in); got != tc.want {
			t.Errorf("extensionFromContentType(%q) = %q, want %q", tc.in, got, tc.want)
		}
	}
}

func TestFirstAttachment(t *testing.T) {
	if got := firstAttachment(nil); got != nil {
		t.Errorf("expected nil for nil message")
	}
	empty := &discordgo.Message{}
	if got := firstAttachment(empty); got != nil {
		t.Errorf("expected nil for no attachments")
	}
	emptyURL := &discordgo.Message{Attachments: []*discordgo.MessageAttachment{{ID: "1"}}}
	if got := firstAttachment(emptyURL); got != nil {
		t.Errorf("expected nil when attachment has no URL")
	}
	with := []*discordgo.MessageAttachment{{ID: "1", URL: "u1"}, {ID: "2", URL: "u2"}}
	msg := &discordgo.Message{Attachments: with}
	if got := firstAttachment(msg); got == nil || got.URL != "u1" {
		t.Errorf("expected first attachment, got %+v", got)
	}
}
