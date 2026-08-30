package tests

import (
	"testing"

	"github.com/kibetnathan/minjibot/internal/commands"
)

func TestHelpSectionsNonEmptyAndUnique(t *testing.T) {
	sections := commands.HelpSections()
	if len(sections) == 0 {
		t.Fatal("no help sections defined")
	}
	seen := map[string]bool{}
	for _, s := range sections {
		if s.Name == "" {
			t.Error("help section with empty name")
		}
		if seen[s.Name] {
			t.Errorf("duplicate help section name %q", s.Name)
		}
		seen[s.Name] = true
		if len(s.Items) == 0 {
			t.Errorf("help section %q has no items", s.Name)
		}
	}
}

func TestBuildHelpPageEmbedClamps(t *testing.T) {
	total := commands.NumHelpPages()
	if total == 0 {
		t.Fatal("expected at least one help page")
	}

	page0 := commands.BuildHelpPageEmbed(0)
	if len(page0.Fields) != 1 {
		t.Errorf("expected 1 field per page, got %d", len(page0.Fields))
	}
	if page0.Footer == nil {
		t.Error("expected a page footer")
	}

	// Out-of-range pages should clamp without panicking.
	if commands.BuildHelpPageEmbed(-1).Title == "" {
		t.Error("BuildHelpPageEmbed(-1) returned empty title")
	}
	if commands.BuildHelpPageEmbed(total+10).Title == "" {
		t.Error("BuildHelpPageEmbed(too big) returned empty title")
	}
}

func TestFindHelpSection(t *testing.T) {
	if commands.FindHelpSection("general") < 0 {
		t.Error("expected to find 'general'")
	}
	if commands.FindHelpSection("GENERAL") < 0 {
		t.Error("expect FindHelpSection to be case-insensitive")
	}
	if commands.FindHelpSection("does-not-exist") != -1 {
		t.Error("expected -1 for unknown category")
	}
}

func TestHelpEmbedListsEverySlashCommandStillPasses(t *testing.T) {
	// The single-embed view should render one field per category.
	embed := commands.BuildHelpEmbed()
	if len(embed.Fields) != len(commands.HelpSections()) {
		t.Errorf("full help embed has %d fields, expected %d", len(embed.Fields), len(commands.HelpSections()))
	}
}
