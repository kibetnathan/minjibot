package tests

import (
	"strings"
	"testing"

	"github.com/kibetnathan/minjibot/internal/commands"
)

func TestTldrKnownAndUnknown(t *testing.T) {
	if _, ok := commands.Tldr("echo"); !ok {
		t.Error("expected tldr to exist for echo")
	}
	if _, ok := commands.Tldr("ECHO"); !ok {
		t.Error("expected tldr to be case-insensitive")
	}
	if _, ok := commands.Tldr("nonexistent-command"); ok {
		t.Error("expected no tldr for unknown command")
	}
}

func TestTldrContent(t *testing.T) {
	entries := []string{
		"ping", "help", "tldr", "echo", "userinfo", "ddg", "search", "gifsearch",
		"isearch", "pinglist", "emoji", "sticker", "pin", "unpin", "quote",
		"translate", "reminder", "caption", "img2gif", "vid2gif", "autogif",
		"factcheck", "howgay", "howautism", "howlesbian", "howsimp", "pp", "puh",
		"iq", "bitches", "choose", "ship", "colors",
	}
	for _, name := range entries {
		e, ok := commands.Tldr(name)
		if !ok {
			t.Errorf("command %q missing a tldr entry", name)
			continue
		}
		if strings.TrimSpace(e.Usage) == "" {
			t.Errorf("tldr for %q has empty usage", name)
		}
		if strings.TrimSpace(e.Explain) == "" {
			t.Errorf("tldr for %q has empty explanation", name)
		}
	}

	// Every registered slash command (except subcommand-only parents) should
	// have a tldr.
	for _, c := range commands.SlashCommands {
		if _, ok := commands.Tldr(c.Name); !ok {
			t.Errorf("slash command %q has no tldr entry", c.Name)
		}
	}
}
