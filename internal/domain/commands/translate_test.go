package commands

import (
	"strings"
	"testing"
)

func TestParseTranslateArgs(t *testing.T) {
	cases := []struct {
		name       string
		args       []string
		wantTarget string
		wantText   string
	}{
		{"default english", []string{"hola"}, "en", "hola"},
		{"to flag", []string{"to:fr", "bonjour", "monde"}, "fr", "bonjour monde"},
		{"lang flag", []string{"lang:ES", "hola"}, "es", "hola"},
		{"uppercase prefix", []string{"TO:de", "hi"}, "de", "hi"},
		{"target only then text", []string{"to:ja", "日本"}, "ja", "日本"},
		{"nothing", nil, "en", ""},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			target, text := parseTranslateArgs(tc.args)
			if target != tc.wantTarget || text != tc.wantText {
				t.Errorf("got (%q,%q), want (%q,%q)", target, text, tc.wantTarget, tc.wantText)
			}
		})
	}
}

func TestTruncateForEmbed(t *testing.T) {
	if got := truncateForEmbed("short", 1024); got != "short" {
		t.Errorf("got %q", got)
	}
	long := strings.Repeat("a", 100)
	if got := truncateForEmbed(long, 10); got != "aaaaaaa..." || len(got) != 10 {
		t.Errorf("got %q (len %d)", got, len(got))
	}
	if got := truncateForEmbed(strings.Repeat("a", 10), 10); got == "" || strings.HasSuffix(got, "...") {
		t.Errorf("exact-length input should be untouched, got %q", got)
	}
	if got := truncateForEmbed(long, 3); got != "..." {
		t.Errorf("max smaller than suffix got %q", got)
	}
}
