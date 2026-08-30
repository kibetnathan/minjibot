package tests

import (
	"strings"
	"testing"

	"github.com/bwmarrin/discordgo"
	"github.com/kibetnathan/minjibot/internal/commands"
)

func TestRandomReadingRange(t *testing.T) {
	for i := 0; i < 500; i++ {
		v := commands.RandomReading()
		if v < 0 || v > 100 {
			t.Fatalf("RandomReading() = %d out of range", v)
		}
	}
}

func TestRandomPPRange(t *testing.T) {
	seen := map[int]bool{}
	for i := 0; i < 2000; i++ {
		v := commands.RandomPP()
		if v < 0 {
			t.Fatalf("RandomPP() = %d < 0", v)
		}
		seen[v] = true
	}
	// A small range should exercise multiple distinct values.
	if len(seen) < 3 {
		t.Errorf("RandomPP() only produced %d distinct values", len(seen))
	}
}

func TestRandomPuhAndIQ(t *testing.T) {
	for i := 0; i < 500; i++ {
		if v := commands.RandomPuh(); v < 0 || v > 100 {
			t.Fatalf("RandomPuh() = %d out of range", v)
		}
		if v := commands.RandomIQ(); v < 50 || v > 160 {
			t.Fatalf("RandomIQ() = %d out of range", v)
		}
	}
}

func TestRandomBitchesRange(t *testing.T) {
	for i := 0; i < 500; i++ {
		if v := commands.RandomBitches(); v < 0 || v > 10 {
			t.Fatalf("RandomBitches() = %d out of range", v)
		}
	}
}

func TestParseChooseItems(t *testing.T) {
	items := commands.ParseChooseItems(" a , b , c ")
	if len(items) != 3 || items[0] != "a" || items[1] != "b" || items[2] != "c" {
		t.Errorf("ParseChooseItems = %v", items)
	}
	if got := commands.ParseChooseItems("  , , "); len(got) != 0 {
		t.Errorf("expected no items, got %v", got)
	}
	if got := commands.ParseChooseItems(""); len(got) != 0 {
		t.Errorf("expected no items for empty input, got %v", got)
	}
}

func TestPickChoose(t *testing.T) {
	items := []string{"a", "b", "c"}
	seen := map[string]bool{}
	for i := 0; i < 500; i++ {
		pick := commands.PickChoose(items)
		if pick != "a" && pick != "b" && pick != "c" {
			t.Fatalf("PickChoose returned unknown item %q", pick)
		}
		seen[pick] = true
	}
	if len(seen) < 3 {
		t.Errorf("PickChoose only returned %v", seen)
	}
	if commands.PickChoose(nil) != "" {
		t.Error("PickChoose(nil) should return empty string")
	}
}

func TestDominantColors(t *testing.T) {
	colors := commands.DominantColors(makePNG(t), 5)
	if len(colors) > 5 {
		t.Errorf("DominantColors returned %d colors (max 5)", len(colors))
	}
	for _, c := range colors {
		if len(c) != 7 || !strings.HasPrefix(c, "#") {
			t.Errorf("color %q not in #rrggbb form", c)
		}
	}
	if commands.DominantColors([]byte("not an image"), 5) != nil {
		t.Error("DominantColors should return nil for undecodable bytes")
	}
}

func TestFunSlashCommandsRegistered(t *testing.T) {
	seen := map[string]*discordgo.ApplicationCommand{}
	for _, c := range commands.SlashCommands {
		seen[c.Name] = c
	}
	for _, name := range []string{"howgay", "howautism", "howlesbian", "howsimp", "pp", "iq", "bitches", "ship"} {
		cmd, ok := seen[name]
		if !ok {
			t.Fatalf("slash command %q missing", name)
		}
		if len(cmd.Options) == 0 {
			t.Errorf("%q should have a user option", name)
		}
	}
	if c, ok := seen["choose"]; !ok {
		t.Error("choose slash command missing")
	} else if len(c.Options) == 0 {
		t.Error("choose should have a choices option")
	}
	if c, ok := seen["colors"]; !ok {
		t.Error("colors slash command missing")
	} else if len(c.Options) == 0 || c.Options[0].Type != discordgo.ApplicationCommandOptionSubCommand {
		t.Error("colors should expose an avatar subcommand")
	}
}
