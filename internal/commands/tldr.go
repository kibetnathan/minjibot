package commands

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
)

// TldrEntry is a brief, plain explanation of how to use a command.
type TldrEntry struct {
	Name    string
	Usage   string
	Explain string
}

// tldrEntries maps command names to their brief usage descriptions.
var tldrEntries = map[string]TldrEntry{
	"help":            {"help", "`help` or `help <category>`", "Show the command menu, optionally filtered to a category (e.g. `fun`)."},
	"tldr":            {"tldr", "`tldr <command>`", "Show a brief how-to for a command (e.g. `-tldr echo`)."},
	"ping":            {"ping", "`ping`", "Check the bot's round-trip and WebSocket latency."},
	"echo":            {"echo", "`echo <text>`", "Repeat the given text back into the channel."},
	"userinfo":        {"userinfo", "`userinfo [@user]`", "Show account/server details for a user."},
	"ddg":             {"ddg", "`ddg <query>`", "Fetch quick search results from DuckDuckGo."},
	"search":          {"search", "`search <query> [count]`", "Search recent chat history for a message matching the query."},
	"gifsearch":       {"gifsearch", "`gifsearch <query> [creator]`", "Post a relevant GIF from Giphy for the query."},
	"isearch":         {"isearch", "`isearch <url>` or reply to/an attach an image", "Reverse image search: find similar images and the pages they come from."},
	"pinglist":        {"pinglist", "`pinglist [@user|role]`", "List all pings for a user or role."},
	"emoji":           {"emoji", "`emoji add|list|remove|steal ...`", "Manage server emojis: add, bulk-add, list, enlarge, remove, or steal."},
	"sticker":         {"sticker", "`sticker add|remove ...`", "Upload or remove custom server stickers."},
	"pin":             {"pin", "`pin <message|link>`", "Pin a message by ID or link."},
	"unpin":           {"unpin", "`unpin <message|link>`", "Unpin a previously pinned message."},
	"quote":           {"quote", "`quote <message|link>`", "Render a message as a styled embed quote."},
	"translate":       {"translate", "`translate [to:<lang>] <text>`", "Translate text into a target language (defaults to English)."},
	"reminder":        {"reminder", "`reminder <text> <delay>`", "Ping yourself after a delay (e.g. `30m`, `2h`)."},
	"caption":         {"caption", "`caption [top] [bottom] [url]`", "Add meme text to an image, from a URL or attachment."},
	"img2gif":         {"img2gif", "`img2gif <url>`", "Convert a static image into an animated GIF."},
	"vid2gif":         {"vid2gif", "`vid2gif <url> [start:<s>] [dur:<s>]`", "Convert a video clip into a GIF (≤25MB, clips to 10s by default)."},
	"autogif":         {"autogif", "`autogif <url>`", "Convert any media (image or video) into a GIF."},
	"factcheck":       {"factcheck", "`factcheck <claim>`", "Fact-check a claim against searchable ratings."},
	"howgay":          {"howgay", "`howgay [@user]`", "Returns a random 0-100% gay reading for a member."},
	"howautism":       {"howautism", "`howautism [@user]`", "Returns a random 0-100% autistic reading for a member."},
	"howlesbian":      {"howlesbian", "`howlesbian [@user]`", "Returns a random 0-100% lesbian reading for a member."},
	"howsimp":         {"howsimp", "`howsimp [@user]`", "Returns a random 0-100% simp reading for a member."},
	"pp":              {"pp", "`pp [@user]`", "Returns a (fake) randomized pp length for a member."},
	"puh":             {"puh", "`puh`", "Returns a random puh tightness percentage."},
	"iq":              {"iq", "`iq [@user]`", "Returns a randomized IQ score for a member."},
	"bitches":         {"bitches", "`bitches [@user]`", "Returns how many bitches a member has."},
	"choose":          {"choose", "`choose <a, b, c>`", "Randomly picks one option from a comma-separated list."},
	"ship":            {"ship", "`ship <@a> <@b>`", "Calculates a romance compatibility score between two members."},
	"colors":          {"colors", "`colors avatar [@user]`", "Extract the dominant colours from a member's avatar."},
	"lurk":            {"lurk", "`lurk`", "Toggle yourself in/out of lurking mode (in-memory, per server)."},
	"lurkers":         {"lurkers", "`lurkers`", "List who is currently lurking in this server."},
	"spark":           {"spark", "`spark`", "Spark (light) the blunt before you can smoke it. Resets to whoever sparks."},
	"smoke":           {"smoke", "`smoke`", "Take a hit off the blunt and add one to your hit count (must spark first)."},
	"hits":            {"hits", "`hits`", "Show everyone's blunt hit count in this server."},
	"compress":        {"compress", "`compress <url>`", "Pixelate/compress an image until it's barely legible."},
	"avatar":          {"avatar", "`avatar [@user]`", "Show a user's full-resolution profile picture."},
	"banner":          {"banner", "`banner [@user]`", "Show a user's profile banner image."},
	"botinfo":         {"botinfo", "`botinfo`", "Show bot info: version, Go runtime, uptime, latency, and server count."},
	"channelinfo":     {"channelinfo", "`channelinfo [channel]`", "Show a channel's creation date, ID, type, topic, slowmode, and category."},
	"roles":           {"roles", "`roles`", "List all server roles with their member counts."},
	"guild":           {"guild", "`guild stats|icon|banner|splash`", "Show server stats, the server icon, the server banner, or the invite splash."},
	"emojis":          {"emojis", "`emojis`", "List every custom emoji in the server."},
	"stickers":        {"stickers", "`stickers`", "List every custom sticker in the server."},
	"bans":            {"bans", "`bans`", "List all active bans in the server (paginated)."},
	"boomer":          {"boomer", "`boomer [@user]`", "Score a user's account age as a lighthearted 0-100 boomer rating."},
	"perms":           {"perms", "`perms [@user] [channel]`", "Show a user's effective server and channel permissions."},
	"tz":              {"tz", "`tz <place>`", "Show the current local time for a place (e.g. `-tz Tokyo` or `-tz New York`)."},
	"urbandictionary": {"urbandictionary", "`urbandictionary <term>`", "Search Urban Dictionary and show the top definition for a term."},
	"ttys":            {"ttys", "`ttys`", "The bot talks to itself in the channel until someone speaks or an hour passes."},
	"vape":            {"vape", "`vape hit|flavor|hits|steal`", "Hit the server vape, set/clear its flavour, show hit counts, or steal it."},
	"poll":            {"poll", "`poll <question> <opt1|opt2|...>`", "Create a reaction-based poll with 2-10 options."},
	"quickpoll":       {"quickpoll", "`quickpoll <question>`", "Create a quick Yes/No reaction poll."},
	"birthday":        {"birthday", "`birthday add|list|celebrate|channel|role`", "Manage server birthdays: save, list, celebrate, set the channel, or set the celebration role."},
	"diary":           {"diary", "`diary add|view|delete`", "Private per-user diary: save, view (DMed), or delete entries."},
	"bio":             {"bio", "`bio github|roblox|reddit|kick <name>`", "Look up a user's public profile on GitHub, Roblox, Reddit, or Kick."},
	"ban":             {"ban", "`ban <user> [reason]`", "Ban a user from the server. Requires a moderation permission."},
	"hardban":         {"hardban", "`hardban <user> [reason]`", "Ban a user and delete their recent messages. Requires a moderation permission."},
	"softban":         {"softban", "`softban <user> [reason]`", "Ban then immediately unban a user, deleting their recent messages. Requires a moderation permission."},
	"kick":            {"kick", "`kick <user> [reason]`", "Kick a user from the server. Requires a moderation permission."},
	"purge":           {"purge", "`purge <count> [user]`", "Delete a number of recent messages, optionally from a specific user. Requires a moderation permission."},
	"nuke":            {"nuke", "`nuke`", "Delete all messages by cloning the current channel. Requires a moderation permission."},
	"timeout":         {"timeout", "`timeout <user> <duration> [reason]`", "Timeout a user for a duration. Requires a moderation permission."},
	"warn":            {"warn", "`warn <user> <reason>`", "Issue a warning to a user. Requires a moderation permission."},
	"history":         {"history", "`history <user>`", "Show a user's moderation history. Requires a moderation permission."},
	"audit":           {"audit", "`audit [limit] [actor]`", "Show recent moderation actions. Requires a moderation permission."},
	"role":            {"role", "`role add|create|edit|hoist|member ...`", "Manage roles. Requires a moderation permission."},
	"fn":              {"fn", "`fn <user> <nickname>`", "Force set a user's nickname. Requires a moderation permission."},
	"nick":            {"nick", "`nick lock|unlock <user>`", "Lock or unlock a user's nickname. Requires a moderation permission."},
}

// Tldr returns the brief usage entry for a command name (case-insensitive).
func Tldr(name string) (TldrEntry, bool) {
	e, ok := tldrEntries[strings.ToLower(name)]
	return e, ok
}

func tldrMessageCommandHandler(s *discordgo.Session, m *discordgo.MessageCreate, args []string) error {
	if len(args) == 0 {
		_, err := s.ChannelMessageSend(m.ChannelID, "Usage: `-tldr <command>` — e.g. `-tldr echo`")
		return err
	}
	return sendTldr(s, m.ChannelID, args[0])
}

func tldrSlashCommandHandler(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	opts := OptionMap(i.ApplicationCommandData().Options)
	name := strings.TrimSpace(OptString(opts, "command"))
	if name == "" {
		return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{Content: "Usage: `/tldr command:<name>`"},
		})
	}
	return respondTldr(s, i, name)
}

func sendTldr(s *discordgo.Session, channelID, name string) error {
	entry, ok := Tldr(name)
	if !ok {
		_, err := s.ChannelMessageSend(channelID, fmt.Sprintf("No tldr found for `%s`. Try `-help` to list commands.", name))
		return err
	}
	embed := &discordgo.MessageEmbed{
		Color: 0x5865F2,
		Title: fmt.Sprintf("tldr: %s", entry.Name),
		Fields: []*discordgo.MessageEmbedField{
			{Name: "Usage", Value: entry.Usage},
			{Name: "What it does", Value: entry.Explain},
		},
	}
	_, err := s.ChannelMessageSendEmbed(channelID, embed)
	return err
}

func respondTldr(s *discordgo.Session, i *discordgo.InteractionCreate, name string) error {
	entry, ok := Tldr(name)
	if !ok {
		return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{Content: fmt.Sprintf("No tldr found for `%s`. Try `/help` to list commands.", name)},
		})
	}
	embed := &discordgo.MessageEmbed{
		Color: 0x5865F2,
		Title: fmt.Sprintf("tldr: %s", entry.Name),
		Fields: []*discordgo.MessageEmbedField{
			{Name: "Usage", Value: entry.Usage},
			{Name: "What it does", Value: entry.Explain},
		},
	}
	return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{Embeds: []*discordgo.MessageEmbed{embed}},
	})
}
