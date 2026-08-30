package tests

import (
	"testing"

	"github.com/bwmarrin/discordgo"
	"github.com/kibetnathan/minjibot/internal/commands"
)

func TestSlashCommandsRegistered(t *testing.T) {
	if len(commands.SlashCommands) == 0 {
		t.Fatal("no slash commands registered")
	}

	seen := map[string]bool{}
	for _, c := range commands.SlashCommands {
		if c.Name == "" || c.Description == "" {
			t.Errorf("command %+v is missing name or description", c)
		}
		if seen[c.Name] {
			t.Errorf("duplicate slash command name %q", c.Name)
		}
		seen[c.Name] = true
	}

	want := []string{
		"ping", "help", "echo", "userinfo", "ddg", "search", "pinglist",
		"gifsearch", "emoji", "sticker", "pin", "unpin", "quote", "translate",
		"reminder", "isearch", "caption", "img2gif", "vid2gif", "autogif", "factcheck",
	}
	for _, name := range want {
		if !seen[name] {
			t.Errorf("expected slash command %q to be registered", name)
		}
	}
}

func TestSubcommandSlashShapes(t *testing.T) {
	for _, name := range []string{"emoji", "sticker"} {
		var cmd *discordgo.ApplicationCommand
		for _, c := range commands.SlashCommands {
			if c.Name == name {
				cmd = c
				break
			}
		}
		if cmd == nil {
			t.Fatalf("%q command missing", name)
		}
		groups := 0
		for _, o := range cmd.Options {
			if o.Type == discordgo.ApplicationCommandOptionSubCommand {
				groups++
			}
		}
		if groups == 0 {
			t.Errorf("%q should expose subcommands", name)
		}
	}
}

func TestCommandHandlerDispatchUnknown(t *testing.T) {
	h := commands.NewCommandHandler(nil, nil, nil, nil)
	if err := h.Handle(nil, nil, &discordgo.MessageCreate{}, "does-not-exist", nil); err == nil {
		t.Error("expected error for unknown prefix command")
	}
	if err := h.HandleSlash(nil, nil, slashInteraction("does-not-exist")); err == nil {
		t.Error("expected error for unknown slash command")
	}
}

func slashInteraction(name string) *discordgo.InteractionCreate {
	return &discordgo.InteractionCreate{
		Interaction: &discordgo.Interaction{
			Type: discordgo.InteractionApplicationCommand,
			Data: discordgo.ApplicationCommandInteractionData{Name: name},
		},
	}
}
