package commands

import (
	"testing"
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
		if got := parseMentionID(tc.in); got != tc.want {
			t.Errorf("parseMentionID(%q) = %q, want %q", tc.in, got, tc.want)
		}
	}
}
