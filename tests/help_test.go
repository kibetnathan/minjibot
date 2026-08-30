package tests

import (
	"strings"
	"testing"

	"github.com/kibetnathan/minjibot/internal/commands"
)

func TestHelpEmbedListsEverySlashCommand(t *testing.T) {
	embed := commands.BuildHelpEmbed()
	var haystack strings.Builder
	for _, f := range embed.Fields {
		haystack.WriteString(f.Name)
		haystack.WriteString("\n")
		haystack.WriteString(f.Value)
		haystack.WriteString("\n")
	}
	if embed.Title == "" || embed.Description == "" {
		t.Errorf("help embed missing title/description: %+v", embed)
	}

	for _, cmd := range commands.SlashCommands {
		name := cmd.Name
		if !containsCommandHelp(haystack.String(), name) {
			t.Errorf("help embed does not mention slash command %q", name)
		}
	}
}

func containsCommandHelp(haystack, name string) bool {
	for _, token := range strings.Fields(haystack) {
		if strings.Trim(token, "`<>[]|@,.") == name {
			return true
		}
	}
	return false
}

func TestHelpEmbedFieldValuesUnderLimit(t *testing.T) {
	for _, f := range commands.BuildHelpEmbed().Fields {
		if len([]rune(f.Name)) > 256 {
			t.Errorf("field name %q too long", f.Name)
		}
		if len([]rune(f.Value)) > 1024 {
			t.Errorf("field value for %q too long", f.Name)
		}
	}
}
