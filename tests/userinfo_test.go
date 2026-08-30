package tests

import (
	"testing"

	"github.com/kibetnathan/minjibot/internal/commands"
)

func TestParseMentionID(t *testing.T) {
	cases := []struct {
		in   string
		want string
	}{
		{"<@123>", "123"},
		{"<@!456>", "456"},
		{"<@&789>", "789"},
		{"  <@123>  ", "123"},
		{"123", "123"},
		{"<@0>", ""},
		{"abc", ""},
		{"-5", ""},
		{"<@!1x>", ""},
		{"", ""},
	}
	for _, tc := range cases {
		if got := commands.ParseMentionID(tc.in); got != tc.want {
			t.Errorf("commands.ParseMentionID(%q) = %q, want %q", tc.in, got, tc.want)
		}
	}
}
