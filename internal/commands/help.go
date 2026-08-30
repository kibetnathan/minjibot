package commands

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
)

// BuildHelpEmbed renders the command list grouped by category. Both prefix
// (-) and slash (/) invocations share it.
func BuildHelpEmbed() *discordgo.MessageEmbed {
	return &discordgo.MessageEmbed{
		Color: 0x5865F2,
		Title: "MinjiBot commands",
		Description: "Commands work as prefix commands (`-help`) or slash commands (`/help`). " +
			"Angle brackets are required, square brackets optional.",
		Fields: []*discordgo.MessageEmbedField{
			helpField("General", [][2]string{
				{"ping", "Check bot latency"},
				{"help", "Show this menu"},
				{"userinfo [@user]", "Get info about a user"},
			}),
			helpField("Search", [][2]string{
				{"ddg <query>", "Fetch quick search results from DuckDuckGo"},
				{"search <query> [count]", "Search chat history for a message"},
				{"isearch <query>", "Search the web for images"},
				{"gifsearch <query> [creator]", "Post a relevant GIF from Giphy"},
			}),
			helpField("Server & messages", [][2]string{
				{"pinglist [user|role]", "Show pings for a user or role"},
				{"emoji add|list|remove|steal", "Manage server emojis"},
				{"sticker add|remove", "Manage server stickers"},
				{"pin <message> [channel]", "Pin a message"},
				{"unpin <message> [channel]", "Unpin a message"},
				{"quote <message> [channel]", "Quote a message as a styled embed"},
			}),
			helpField("Utilities", [][2]string{
				{"echo <text>", "Repeat back a message"},
				{"translate <text> [lang]", "Translate text into a target language"},
				{"reminder <text> <delay>", "Set a delayed reminder ping"},
				{"caption [top] [bottom] [url]", "Add meme text to an image"},
			}),
			helpField("Media & AI", [][2]string{
				{"img2gif <url>", "Convert an image into a GIF"},
				{"vid2gif <url>", "Convert a video into a GIF"},
				{"autogif <url>", "Convert any media into a GIF"},
				{"factcheck <claim>", "Fact-check a claim against searchable ratings"},
			}),
		},
		Footer: &discordgo.MessageEmbedFooter{Text: "MinjiBot"},
	}
}

func helpField(name string, items [][2]string) *discordgo.MessageEmbedField {
	var b strings.Builder
	for _, it := range items {
		fmt.Fprintf(&b, "`%s` — %s\n", it[0], it[1])
	}
	return &discordgo.MessageEmbedField{Name: name, Value: b.String()}
}
