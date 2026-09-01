package commands

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
)

// HelpSection is one paginated category of the help menu.
type HelpSection struct {
	Name  string
	Items [][2]string
}

// HelpSections returns every help category and its command entries. It is the
// single source of truth used by the paginated help command.
func HelpSections() []HelpSection {
	return []HelpSection{
		{Name: "General", Items: [][2]string{
			{"ping", "Check bot latency"},
			{"help", "Show this menu"},
			{"tldr <command>", "Get a brief how-to for a command"},
			{"userinfo [@user]", "Get info about a user"},
		}},
		{Name: "Information", Items: [][2]string{
			{"avatar [@user]", "Show a user's full-resolution profile picture"},
			{"banner [@user]", "Show a user's profile banner"},
			{"botinfo", "Show bot info (version, uptime, latency)"},
			{"channelinfo [channel]", "Get info about a channel"},
			{"guild stats|icon|banner", "Server stats, icon, or banner"},
			{"roles", "List all server roles with member counts"},
			{"emojis", "List all custom emojis in the server"},
			{"stickers", "List all custom stickers in the server"},
		}},
		{Name: "Search", Items: [][2]string{
			{"ddg <query>", "Fetch quick search results from DuckDuckGo"},
			{"search <query> [count]", "Search chat history for a message"},
			{"isearch <query>", "Search the web for images"},
			{"gifsearch <query> [creator]", "Post a relevant GIF from Giphy"},
		}},
		{Name: "Server & messages", Items: [][2]string{
			{"pinglist [user|role]", "Show pings for a user or role"},
			{"emoji add|list|remove|steal", "Manage server emojis"},
			{"sticker add|steal|remove", "Manage server stickers"},
			{"pin <message> [channel]", "Pin a message"},
			{"unpin <message> [channel]", "Unpin a message"},
			{"quote <message> [channel]", "Quote a message as a styled embed"},
		}},
		{Name: "Utilities", Items: [][2]string{
			{"echo <text>", "Repeat back a message"},
			{"translate <text> [lang]", "Translate text into a target language"},
			{"reminder <text> <delay>", "Set a delayed reminder ping"},
			{"caption [top] [bottom] [url]", "Add meme text to an image"},
		}},
		{Name: "Media & AI", Items: [][2]string{
			{"img2gif <url>", "Convert an image into a GIF"},
			{"vid2gif <url>", "Convert a video into a GIF (≤25MB, clips to 10s)"},
			{"autogif <url>", "Convert any media into a GIF"},
			{"factcheck <claim>", "Fact-check a claim against searchable ratings"},
		}},
		{Name: "Fun", Items: [][2]string{
			{"howgay [@user]", "Measure how gay a member is"},
			{"howautism [@user]", "Measure how autistic a member is"},
			{"howlesbian [@user]", "Measure how lesbian a member is"},
			{"howsimp [@user]", "Measure how much of a simp a member is"},
			{"pp [@user]", "Measure a member's pp length"},
			{"puh", "Check the puh tightness"},
			{"iq [@user]", "Measure a member's IQ"},
			{"bitches [@user]", "See how many bitches a member has"},
			{"choose <a, b, c>", "Pick an option from a comma-separated list"},
			{"ship <@user> <@user>", "Calculate romance compatibility"},
			{"colors avatar [@user]", "Extract dominant colours from an avatar"},
			{"lurk", "Toggle yourself in/out of lurking mode"},
			{"lurkers", "Show who is currently lurking"},
			{"spark", "Spark the blunt before you can smoke"},
			{"smoke", "Take a hit off the blunt (spark it first)"},
			{"hits", "Show everyone's blunt hit count"},
			{"compress <url>", "Compress an image until it's barely legible"},
			{"vape hit|flavor|hits|steal", "Hit, configure, or check the server vape"},
			{"poll <question> <a|b|c>", "Create a reaction poll (2-10 options)"},
			{"quickpoll <question>", "Create a Yes/No poll"},
			{"birthday add|list|celebrate|channel|role", "Manage server birthdays"},
			{"diary add|view|delete", "Private per-user diary"},
		}},
	}
}

// NumHelpPages returns the number of help categories/pages.
func NumHelpPages() int {
	return len(HelpSections())
}

// FindHelpSection returns the zero-based index of the category whose name
// matches name (case-insensitive), or -1 if no match.
func FindHelpSection(name string) int {
	for i, section := range HelpSections() {
		if strings.EqualFold(section.Name, name) {
			return i
		}
	}
	return -1
}

// BuildHelpEmbed renders the full command list in a single embed. Kept for
// convenience and tests; the interactive help command is paginated instead.
func BuildHelpEmbed() *discordgo.MessageEmbed {
	fields := make([]*discordgo.MessageEmbedField, 0, len(HelpSections()))
	for _, section := range HelpSections() {
		fields = append(fields, helpField(section.Name, section.Items))
	}
	return &discordgo.MessageEmbed{
		Color: 0x5865F2,
		Title: "MinjiBot commands",
		Description: "Commands work as prefix commands (`-help`) or slash commands (`/help`). " +
			"Angle brackets are required, square brackets optional.",
		Fields: fields,
		Footer: &discordgo.MessageEmbedFooter{Text: "MinjiBot"},
	}
}

// BuildHelpPageEmbed renders a single help category as one page. page is the
// zero-based section index. A page outside range clamps to the nearest valid
// index.
func BuildHelpPageEmbed(page int) *discordgo.MessageEmbed {
	sections := HelpSections()
	if len(sections) == 0 {
		return BuildHelpEmbed()
	}
	if page < 0 {
		page = 0
	}
	if page >= len(sections) {
		page = len(sections) - 1
	}
	section := sections[page]
	return &discordgo.MessageEmbed{
		Color: 0x5865F2,
		Title: fmt.Sprintf("MinjiBot commands — %s", section.Name),
		Description: "Commands work as prefix commands (`-help`) or slash commands (`/help`). " +
			"Angle brackets are required, square brackets optional.",
		Fields: []*discordgo.MessageEmbedField{helpField(section.Name, section.Items)},
		Footer: &discordgo.MessageEmbedFooter{
			Text: fmt.Sprintf("Page %d/%d — use the arrows to navigate", page+1, len(sections)),
		},
	}
}

func helpField(name string, items [][2]string) *discordgo.MessageEmbedField {
	var b strings.Builder
	for _, it := range items {
		fmt.Fprintf(&b, "`%s` — %s\n", it[0], it[1])
	}
	return &discordgo.MessageEmbedField{Name: name, Value: b.String()}
}
