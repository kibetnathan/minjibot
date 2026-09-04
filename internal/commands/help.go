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
			{"test", "Ping-pong check that the bot is responsive"},
			{"bug", "Open a form to report a bot bug to developers"},
			{"donate", "Buy the developer a coffee to support MinjiBot"},
			{"setup logchannel|status", "Set or view server logging (owners and admins)"},
			{"help", "Show this menu"},
			{"tldr <command>", "Get a brief how-to for a command"},
			{"userinfo [@user]", "Get info about a user"},
		}},
		{Name: "Information", Items: [][2]string{
			{"avatar [@user]", "Show a user's full-resolution profile picture"},
			{"banner [@user]", "Show a user's profile banner"},
			{"botinfo", "Show bot info (version, uptime, latency)"},
			{"channelinfo [channel]", "Get info about a channel"},
			{"guild stats|icon|banner|splash", "Server stats, icon, banner, or splash"},
			{"roles", "List all server roles with member counts"},
			{"emojis", "List all custom emojis in the server"},
			{"stickers", "List all custom stickers in the server"},
			{"bans", "List all active bans in the server"},
			{"boomer [@user]", "Detect potential time-traveler users"},
			{"perms [@user] [channel]", "Show a user's effective permissions"},
			{"tz <place>", "Show the current local time in a place (e.g. Tokyo)"},
			{"urbandictionary <term>", "Search Urban Dictionary for a term"},
			{"weather <place>", "Current weather, forecast, and humidity for a place"},
		}},
		{Name: "Search", Items: [][2]string{
			{"ddg <query>", "Fetch quick search results from DuckDuckGo"},
			{"search <query> [count]", "Search chat history for a message"},
			{"isearch <url>", "Reverse image search: find similar images & sources"},
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
			{"compress <url>", "Compress an image until it's barely legible"},
			{"poll <question> <a|b|c>", "Create a reaction poll (2-10 options)"},
			{"quickpoll <question>", "Create a Yes/No poll"},
			{"birthday add|list|celebrate|channel|role", "Manage server birthdays"},
			{"diary add|view|delete", "Private per-user diary"},
			{"ttys", "Bot talks to itself until someone speaks or up to an hour"},
		}},
		{Name: "Moderation — Punishments", Items: [][2]string{
			{"ban <user> [reason]", "Ban a user from the server"},
			{"hardban <user> [reason]", "Ban a user and delete their recent messages"},
			{"softban <user> [reason]", "Ban then immediately unban a user, deleting recent messages"},
			{"kick <user> [reason]", "Kick a user from the server"},
			{"timeout <user> <duration> [reason]", "Timeout a user (e.g. 30m, 2h, 1d)"},
			{"warn <user> <reason>", "Warn a user"},
			{"history <user>", "Show a user's moderation history"},
			{"audit [limit] [actor]", "Show recent moderation actions"},
			{"jail <user>", "Jail a user, removing their roles"},
			{"unjail <user>", "Unjail a user, restoring their roles"},
		}},
		{Name: "Moderation — Roles & nick", Items: [][2]string{
			{"role add|create|edit|hoist|member", "Manage roles"},
			{"fn <user> <nickname>", "Force set a user's nickname"},
			{"nick lock|unlock <user>", "Lock or unlock a user's nickname"},
			{"staffstrip <user>", "Remove all staff roles from a user"},
		}},
		{Name: "Moderation — Channels", Items: [][2]string{
			{"purge <count> [user]", "Delete recent messages, optionally only from a user"},
			{"nuke", "Delete all messages by cloning the current channel"},
			{"hide", "Hide the current channel from @everyone"},
			{"reveal", "Make the current channel visible to @everyone again"},
			{"lockdown", "Lock the current channel for @everyone"},
			{"nsfw", "Mark the current channel as NSFW"},
			{"sfw", "Unmark the current channel as NSFW"},
			{"slowmode <seconds>", "Set a slowmode on the current channel"},
			{"topic <text>", "Set the topic of the current channel"},
		}},
		{Name: "Moderation — Permissions", Items: [][2]string{
			{"denyperm <user|role> <perm> [channel]", "Deny a permission to a user or role"},
			{"imute <user>", "Prevent a user from sending images in this channel"},
			{"gifmute <user>", "Prevent a user from images, GIFs, and embeds"},
		}},
		{Name: "RP — Emotions", Items: [][2]string{
			{"angry", "Post a GIF expressing being angry"},
			{"depressed", "Post a GIF expressing being depressed"},
			{"excited", "Post a GIF expressing being excited"},
			{"happy", "Post a GIF expressing being happy"},
			{"horny", "Post a GIF expressing being horny"},
			{"inlove", "Post a GIF expressing being in love"},
			{"sad", "Post a GIF expressing being sad"},
			{"shy", "Post a GIF expressing being shy"},
		}},
		{Name: "RP — Actions (A–H)", Items: [][2]string{
			{"baka [user]", "Call someone a baka (anime GIF)"},
			{"bite [user]", "Bite a user (anime GIF)"},
			{"cry [user]", "Cry at a user (anime GIF)"},
			{"dap [user]", "Dap up a user (anime GIF)"},
			{"eat [user]", "Munch on a user (anime GIF)"},
			{"facepalm [user]", "Facepalm at a user (anime GIF)"},
			{"feed [user]", "Feed a user (anime GIF)"},
			{"handhold [user]", "Hold hands with a user (anime GIF)"},
		}},
		{Name: "RP — Actions (K–T)", Items: [][2]string{
			{"kiss [user]", "Kiss a user (anime GIF)"},
			{"laugh [user]", "Laugh at a user (anime GIF)"},
			{"nod [user]", "Nod at a user (anime GIF)"},
			{"nutkick [user]", "Nutkick a user (anime GIF)"},
			{"pat [user]", "Pat a user (anime GIF)"},
			{"peck [user]", "Peck a user (anime GIF)"},
			{"poke [user]", "Poke a user (anime GIF)"},
			{"punch [user]", "Punch a user (anime GIF)"},
			{"run [user]", "Run away from a user (anime GIF)"},
			{"shoot [user]", "Shoot a user (anime GIF)"},
			{"shrug [user]", "Shrug at a user (anime GIF)"},
			{"slap [user]", "Slap a user (anime GIF)"},
			{"spank [user]", "Spank a user (anime GIF)"},
			{"stab [user]", "Stab a user (anime GIF)"},
			{"think [user]", "Think about a user (anime GIF)"},
			{"tickle [user]", "Tickle a user (anime GIF)"},
		}},
		{Name: "Social", Items: [][2]string{
			{"bio github|roblox|reddit|kick <name>", "Look up a user's public profile"},
			{"ship <@user> <@user>", "Calculate romance compatibility"},
			{"colors avatar [@user]", "Extract dominant colours from an avatar"},
			{"lurk", "Toggle yourself in/out of lurking mode"},
			{"lurkers", "Show who is currently lurking"},
			{"spark", "Spark the blunt before you can smoke"},
			{"smoke", "Take a hit off the blunt (spark it first)"},
			{"hits", "Show everyone's blunt hit count"},
			{"vape hit|flavor|hits|steal", "Hit, configure, or check the server vape"},
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
