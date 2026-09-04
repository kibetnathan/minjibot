package commands

import "github.com/bwmarrin/discordgo"

var SlashCommands = []*discordgo.ApplicationCommand{
	{
		Name:        "ping",
		Description: "Check bot latency",
	},
	{
		Name:        "help",
		Description: "Show available commands",
	},
	{
		Name:        "tldr",
		Description: "Get a brief how-to for a command",
		Options: []*discordgo.ApplicationCommandOption{
			{Type: discordgo.ApplicationCommandOptionString, Name: "command", Description: "The command to explain", Required: true},
		},
	},
	{
		Name:        "echo",
		Description: "Repeat back a message",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "text",
				Description: "Text to echo",
				Required:    true,
			},
		},
	},
	{
		Name:        "userinfo",
		Description: "Get info about a user",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionUser,
				Name:        "user",
				Description: "The user to look up (defaults to yourself)",
				Required:    false,
			},
		},
	},
	{
		Name:        "avatar",
		Description: "Show a user's full-resolution profile picture",
		Options:     []*discordgo.ApplicationCommandOption{userOption(false)},
	},
	{
		Name:        "banner",
		Description: "Show a user's profile banner",
		Options:     []*discordgo.ApplicationCommandOption{userOption(false)},
	},
	{
		Name:        "botinfo",
		Description: "Show bot info (version, uptime, latency)",
	},
	{
		Name:        "channelinfo",
		Description: "Get info about a channel",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionChannel,
				Name:        "channel",
				Description: "The channel to inspect (defaults to this one)",
				Required:    false,
			},
		},
	},
	{
		Name:        "roles",
		Description: "List all server roles with member counts",
	},
	{
		Name:        "guild",
		Description: "Server info (stats, icon, banner, splash)",
		Options: []*discordgo.ApplicationCommandOption{
			{Type: discordgo.ApplicationCommandOptionSubCommand, Name: "stats", Description: "Server stats (members, boosts, owner)"},
			{Type: discordgo.ApplicationCommandOptionSubCommand, Name: "icon", Description: "Server icon image"},
			{Type: discordgo.ApplicationCommandOptionSubCommand, Name: "banner", Description: "Server banner image"},
			{Type: discordgo.ApplicationCommandOptionSubCommand, Name: "splash", Description: "Server invite splash image"},
		},
	},
	{
		Name:        "emojis",
		Description: "List all custom emojis in the server",
	},
	{
		Name:        "stickers",
		Description: "List all custom stickers in the server",
	},
	{
		Name:        "bans",
		Description: "List all active bans in the server",
	},
	{
		Name:        "boomer",
		Description: "Detect potential time-traveler users (spammer detection)",
		Options:     []*discordgo.ApplicationCommandOption{userOption(false)},
	},
	{
		Name:        "perms",
		Description: "Show a user's effective permissions in a channel",
		Options: []*discordgo.ApplicationCommandOption{
			userOption(false),
			{
				Type:        discordgo.ApplicationCommandOptionChannel,
				Name:        "channel",
				Description: "The channel to check (defaults to this one)",
				Required:    false,
			},
		},
	},
	{
		Name:        "tz",
		Description: "Show the current local time for a place",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "place",
				Description: "A city or town, e.g. Tokyo or New York",
				Required:    true,
			},
		},
	},
	{
		Name:        "urbandictionary",
		Description: "Search Urban Dictionary for a term",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "term",
				Description: "The term to look up",
				Required:    true,
			},
		},
	},
	{
		Name:        "weather",
		Description: "Current weather, forecast, and humidity for a location",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "location",
				Description: "City, town, or place to check the weather",
				Required:    true,
			},
		},
	},
	{
		Name:        "test",
		Description: "Ping-pong functionality check to confirm the bot is responsive",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "text",
				Description: "Optional text to echo back",
				Required:    false,
			},
		},
	},
	{
		Name:        "bug",
		Description: "Open a form to report a bot bug to the developers",
	},
	{
		Name:        "donate",
		Description: "Buy the developer a coffee to support MinjiBot",
	},
	{
		Name:        "setup",
		Description: "Set up server logging (owners and admins)",
		Options: []*discordgo.ApplicationCommandOption{
			{Type: discordgo.ApplicationCommandOptionSubCommand, Name: "logchannel", Description: "Set the channel where deleted messages and moderation actions are logged", Options: []*discordgo.ApplicationCommandOption{
				{Type: discordgo.ApplicationCommandOptionChannel, Name: "channel", Description: "The channel to log to", Required: true},
			}},
			{Type: discordgo.ApplicationCommandOptionSubCommand, Name: "status", Description: "Show the current logging configuration"},
		},
	},
	{
		Name:        "ddg",
		Description: "Fetch quick search results from DuckDuckGo",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "query",
				Description: "The search query",
				Required:    true,
			},
		},
	},
	{
		Name:        "search",
		Description: "Search chat history for a specific message",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "query",
				Description: "Text to search for in chat history",
				Required:    true,
			},
			{
				Type:        discordgo.ApplicationCommandOptionInteger,
				Name:        "messages",
				Description: "How many recent messages to search (default 200, max 1000)",
				Required:    false,
			},
		},
	},
	{
		Name:        "pinglist",
		Description: "Returns all the pings for a certain user/role",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionUser,
				Name:        "user",
				Description: "User to show pings for",
				Required:    false,
			},
			{
				Type:        discordgo.ApplicationCommandOptionRole,
				Name:        "role",
				Description: "Role to show pings for",
				Required:    false,
			},
			{
				Type:        discordgo.ApplicationCommandOptionInteger,
				Name:        "messages",
				Description: "How many recent messages to search (default 1000, max 1000)",
				Required:    false,
			},
		},
	},
	{
		Name:        "gifsearch",
		Description: "Searches Giphy and posts a relevant GIF",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "query",
				Description: "The GIF search query",
				Required:    true,
			},
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "creator",
				Description: "Only show GIFs by a specific creator",
				Required:    false,
			},
		},
	},
	{
		Name:        "emoji",
		Description: "Manage server emojis",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionSubCommand,
				Name:        "add",
				Description: "Upload an emoji to the server",
				Options: []*discordgo.ApplicationCommandOption{
					{Type: discordgo.ApplicationCommandOptionString, Name: "name", Description: "Emoji name", Required: true},
					{Type: discordgo.ApplicationCommandOptionString, Name: "url", Description: "Emoji image URL", Required: false},
				},
			},
			{
				Type:        discordgo.ApplicationCommandOptionSubCommand,
				Name:        "enlarge",
				Description: "Enlarge an emoji",
				Options: []*discordgo.ApplicationCommandOption{
					{Type: discordgo.ApplicationCommandOptionString, Name: "emoji", Description: "The emoji to enlarge", Required: true},
				},
			},
			{
				Type:        discordgo.ApplicationCommandOptionSubCommand,
				Name:        "list",
				Description: "List all server emojis",
			},
			{
				Type:        discordgo.ApplicationCommandOptionSubCommand,
				Name:        "remove",
				Description: "Remove a server emoji",
				Options: []*discordgo.ApplicationCommandOption{
					{Type: discordgo.ApplicationCommandOptionString, Name: "emoji", Description: "Emoji to remove", Required: true},
				},
			},
			{
				Type:        discordgo.ApplicationCommandOptionSubCommand,
				Name:        "steal",
				Description: "Copy an emoji into this server",
				Options: []*discordgo.ApplicationCommandOption{
					{Type: discordgo.ApplicationCommandOptionString, Name: "emoji", Description: "Emoji to steal", Required: true},
				},
			},
		},
	},
	{
		Name:        "sticker",
		Description: "Manage server stickers",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionSubCommand,
				Name:        "add",
				Description: "Upload a sticker to the server",
				Options: []*discordgo.ApplicationCommandOption{
					{Type: discordgo.ApplicationCommandOptionString, Name: "name", Description: "Sticker name", Required: true},
					{Type: discordgo.ApplicationCommandOptionString, Name: "url", Description: "Sticker image URL", Required: true},
				},
			},
			{
				Type:        discordgo.ApplicationCommandOptionSubCommand,
				Name:        "steal",
				Description: "Copy a sticker from a message link or ID",
				Options: []*discordgo.ApplicationCommandOption{
					{Type: discordgo.ApplicationCommandOptionString, Name: "sticker", Description: "Sticker ID, message link, or CDN URL", Required: true},
				},
			},
			{
				Type:        discordgo.ApplicationCommandOptionSubCommand,
				Name:        "remove",
				Description: "Remove a server sticker",
				Options: []*discordgo.ApplicationCommandOption{
					{Type: discordgo.ApplicationCommandOptionString, Name: "sticker_id", Description: "Sticker ID or message link", Required: true},
				},
			},
		},
	},
	{
		Name:        "pin",
		Description: "Pin a message",
		Options: []*discordgo.ApplicationCommandOption{
			{Type: discordgo.ApplicationCommandOptionString, Name: "message", Description: "Message ID or link", Required: true},
			{Type: discordgo.ApplicationCommandOptionChannel, Name: "channel", Description: "Channel (defaults to this one)", Required: false},
		},
	},
	{
		Name:        "unpin",
		Description: "Unpin a message",
		Options: []*discordgo.ApplicationCommandOption{
			{Type: discordgo.ApplicationCommandOptionString, Name: "message", Description: "Message ID or link", Required: true},
			{Type: discordgo.ApplicationCommandOptionChannel, Name: "channel", Description: "Channel (defaults to this one)", Required: false},
		},
	},
	{
		Name:        "quote",
		Description: "Quote a message as a styled embed",
		Options: []*discordgo.ApplicationCommandOption{
			{Type: discordgo.ApplicationCommandOptionString, Name: "message", Description: "Message ID or link", Required: true},
			{Type: discordgo.ApplicationCommandOptionChannel, Name: "channel", Description: "Channel (defaults to this one)", Required: false},
		},
	},
	{
		Name:        "translate",
		Description: "Translate text into a target language",
		Options: []*discordgo.ApplicationCommandOption{
			{Type: discordgo.ApplicationCommandOptionString, Name: "text", Description: "Text to translate", Required: true},
			{Type: discordgo.ApplicationCommandOptionString, Name: "target", Description: "Target language code (default: en)", Required: false},
		},
	},
	{
		Name:        "reminder",
		Description: "Set a delayed reminder ping",
		Options: []*discordgo.ApplicationCommandOption{
			{Type: discordgo.ApplicationCommandOptionString, Name: "text", Description: "What to remind you about", Required: true},
			{Type: discordgo.ApplicationCommandOptionString, Name: "delay", Description: "When, e.g. 30m / 2h / 1h30m", Required: true},
		},
	},
	{
		Name:        "isearch",
		Description: "Reverse image search: find where an image appears online",
		Options: []*discordgo.ApplicationCommandOption{
			{Type: discordgo.ApplicationCommandOptionString, Name: "image_url", Description: "Image URL to search (jpg/png/webp)", Required: true},
		},
	},
	{
		Name:        "caption",
		Description: "Add meme text to an image",
		Options: []*discordgo.ApplicationCommandOption{
			{Type: discordgo.ApplicationCommandOptionString, Name: "top", Description: "Top text", Required: false},
			{Type: discordgo.ApplicationCommandOptionString, Name: "bottom", Description: "Bottom text", Required: false},
			{Type: discordgo.ApplicationCommandOptionString, Name: "image_url", Description: "Background image URL", Required: false},
		},
	},
	{
		Name:        "img2gif",
		Description: "Convert an image into a GIF",
		Options: []*discordgo.ApplicationCommandOption{
			{Type: discordgo.ApplicationCommandOptionString, Name: "url", Description: "Image URL", Required: true},
		},
	},
	{
		Name:        "vid2gif",
		Description: "Convert a video into a GIF (≤25MB, clips to 10s)",
		Options: []*discordgo.ApplicationCommandOption{
			{Type: discordgo.ApplicationCommandOptionString, Name: "url", Description: "Video URL", Required: true},
		},
	},
	{
		Name:        "autogif",
		Description: "Convert any media into a GIF",
		Options: []*discordgo.ApplicationCommandOption{
			{Type: discordgo.ApplicationCommandOptionString, Name: "url", Description: "Image or video URL", Required: true},
		},
	},
	{
		Name:        "factcheck",
		Description: "Fact-check a claim against searchable ratings",
		Options: []*discordgo.ApplicationCommandOption{
			{Type: discordgo.ApplicationCommandOptionString, Name: "claim", Description: "The claim to fact-check", Required: true},
		},
	},
	{
		Name:        "howgay",
		Description: "Measure how gay a member is",
		Options:     []*discordgo.ApplicationCommandOption{userOption(false)},
	},
	{
		Name:        "howautism",
		Description: "Measure how autistic a member is",
		Options:     []*discordgo.ApplicationCommandOption{userOption(false)},
	},
	{
		Name:        "howlesbian",
		Description: "Measure how lesbian a member is",
		Options:     []*discordgo.ApplicationCommandOption{userOption(false)},
	},
	{
		Name:        "howsimp",
		Description: "Measure how much of a simp a member is",
		Options:     []*discordgo.ApplicationCommandOption{userOption(false)},
	},
	{
		Name:        "pp",
		Description: "Measure a member's pp length",
		Options:     []*discordgo.ApplicationCommandOption{userOption(false)},
	},
	{
		Name:        "puh",
		Description: "Check the puh tightness",
		Options:     []*discordgo.ApplicationCommandOption{userOption(false)},
	},
	{
		Name:        "iq",
		Description: "Measure a member's IQ",
		Options:     []*discordgo.ApplicationCommandOption{userOption(false)},
	},
	{
		Name:        "bitches",
		Description: "See how many bitches a member has",
		Options:     []*discordgo.ApplicationCommandOption{userOption(false)},
	},
	{
		Name:        "choose",
		Description: "Pick an option from a comma-separated list",
		Options: []*discordgo.ApplicationCommandOption{
			{Type: discordgo.ApplicationCommandOptionString, Name: "choices", Description: "Comma-separated options, e.g. a, b, c", Required: true},
		},
	},
	{
		Name:        "ship",
		Description: "Calculate romance compatibility between two members",
		Options: []*discordgo.ApplicationCommandOption{
			{Type: discordgo.ApplicationCommandOptionUser, Name: "user1", Description: "First user", Required: true},
			{Type: discordgo.ApplicationCommandOptionUser, Name: "user2", Description: "Second user", Required: true},
		},
	},
	{
		Name:        "colors",
		Description: "Extract dominant colours from an avatar",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionSubCommand,
				Name:        "avatar",
				Description: "Dominant colours from a member's avatar",
				Options:     []*discordgo.ApplicationCommandOption{userOption(false)},
			},
		},
	},
	{
		Name:        "lurk",
		Description: "Toggle yourself in/out of lurking mode",
	},
	{
		Name:        "lurkers",
		Description: "Show who is currently lurking",
	},
	{
		Name:        "spark",
		Description: "Spark the blunt before you can smoke",
	},
	{
		Name:        "smoke",
		Description: "Take a hit off the blunt (spark it first)",
	},
	{
		Name:        "hits",
		Description: "Show everyone's blunt hit count",
	},
	{
		Name:        "compress",
		Description: "Compress an image until it's barely legible",
		Options: []*discordgo.ApplicationCommandOption{
			{Type: discordgo.ApplicationCommandOptionString, Name: "url", Description: "Image URL", Required: false},
		},
	},
	{
		Name:        "vape",
		Description: "Hit, configure, or check the server vape",
		Options: []*discordgo.ApplicationCommandOption{
			{Type: discordgo.ApplicationCommandOptionSubCommand, Name: "hit", Description: "Take a hit off the vape"},
			{Type: discordgo.ApplicationCommandOptionSubCommand, Name: "flavor", Description: "Set or clear the vape flavour", Options: []*discordgo.ApplicationCommandOption{
				{Type: discordgo.ApplicationCommandOptionString, Name: "flavour", Description: "Flavour text", Required: false},
			}},
			{Type: discordgo.ApplicationCommandOptionSubCommand, Name: "hits", Description: "Show everyone's vape hit count"},
			{Type: discordgo.ApplicationCommandOptionSubCommand, Name: "steal", Description: "Steal the vape"},
		},
	},
	{
		Name:        "poll",
		Description: "Create a reaction-based poll",
		Options: []*discordgo.ApplicationCommandOption{
			{Type: discordgo.ApplicationCommandOptionString, Name: "question", Description: "The poll question", Required: true},
			{Type: discordgo.ApplicationCommandOptionString, Name: "options", Description: "Options separated by | (2-10)", Required: true},
		},
	},
	{
		Name:        "quickpoll",
		Description: "Create a quick Yes/No poll",
		Options: []*discordgo.ApplicationCommandOption{
			{Type: discordgo.ApplicationCommandOptionString, Name: "question", Description: "The poll question", Required: true},
		},
	},
	{
		Name:        "birthday",
		Description: "Manage server birthdays",
		Options: []*discordgo.ApplicationCommandOption{
			{Type: discordgo.ApplicationCommandOptionSubCommand, Name: "add", Description: "Save a birthday", Options: []*discordgo.ApplicationCommandOption{
				{Type: discordgo.ApplicationCommandOptionString, Name: "date", Description: "Date e.g. 14-07 or 1998-14-07", Required: true},
				userOption(false),
			}},
			{Type: discordgo.ApplicationCommandOptionSubCommand, Name: "list", Description: "List upcoming birthdays"},
			{Type: discordgo.ApplicationCommandOptionSubCommand, Name: "celebrate", Description: "Celebrate a birthday", Options: []*discordgo.ApplicationCommandOption{
				userOption(false),
			}},
			{Type: discordgo.ApplicationCommandOptionSubCommand, Name: "channel", Description: "Set the birthday celebration channel", Options: []*discordgo.ApplicationCommandOption{
				{Type: discordgo.ApplicationCommandOptionChannel, Name: "channel", Description: "Channel", Required: true},
			}},
			{Type: discordgo.ApplicationCommandOptionSubCommand, Name: "role", Description: "Set the temporary birthday role", Options: []*discordgo.ApplicationCommandOption{
				{Type: discordgo.ApplicationCommandOptionRole, Name: "role", Description: "Role", Required: true},
			}},
		},
	},
	{
		Name:        "diary",
		Description: "Private per-user diary",
		Options: []*discordgo.ApplicationCommandOption{
			{Type: discordgo.ApplicationCommandOptionSubCommand, Name: "add", Description: "Save a diary entry", Options: []*discordgo.ApplicationCommandOption{
				{Type: discordgo.ApplicationCommandOptionString, Name: "text", Description: "Entry text", Required: true},
			}},
			{Type: discordgo.ApplicationCommandOptionSubCommand, Name: "view", Description: "View your diary (DMed privately)"},
			{Type: discordgo.ApplicationCommandOptionSubCommand, Name: "delete", Description: "Delete a diary entry", Options: []*discordgo.ApplicationCommandOption{
				{Type: discordgo.ApplicationCommandOptionInteger, Name: "id", Description: "Entry ID", Required: true},
			}},
		},
	},
	{
		Name:        "ttys",
		Description: "Bot talks to itself until someone speaks or an hour passes",
	},
	{
		Name:        "bio",
		Description: "Look up a user's public profile on a supported platform",
		Options: []*discordgo.ApplicationCommandOption{
			{Type: discordgo.ApplicationCommandOptionSubCommand, Name: "github", Description: "Look up a GitHub profile", Options: []*discordgo.ApplicationCommandOption{
				{Type: discordgo.ApplicationCommandOptionString, Name: "username", Description: "GitHub username", Required: true},
			}},
			{Type: discordgo.ApplicationCommandOptionSubCommand, Name: "roblox", Description: "Look up a Roblox profile", Options: []*discordgo.ApplicationCommandOption{
				{Type: discordgo.ApplicationCommandOptionString, Name: "username", Description: "Roblox username", Required: true},
			}},
			{Type: discordgo.ApplicationCommandOptionSubCommand, Name: "reddit", Description: "Look up a Reddit profile", Options: []*discordgo.ApplicationCommandOption{
				{Type: discordgo.ApplicationCommandOptionString, Name: "username", Description: "Reddit username", Required: true},
			}},
			{Type: discordgo.ApplicationCommandOptionSubCommand, Name: "kick", Description: "Look up a Kick channel", Options: []*discordgo.ApplicationCommandOption{
				{Type: discordgo.ApplicationCommandOptionString, Name: "username", Description: "Kick channel slug", Required: true},
			}},
		},
	},
	{
		Name:        "ban",
		Description: "Ban a user from the server",
		Options: []*discordgo.ApplicationCommandOption{
			{Type: discordgo.ApplicationCommandOptionUser, Name: "user", Description: "User to ban", Required: true},
			{Type: discordgo.ApplicationCommandOptionString, Name: "reason", Description: "Ban reason", Required: false},
		},
	},
	{
		Name:        "hardban",
		Description: "Ban a user and delete their recent messages",
		Options: []*discordgo.ApplicationCommandOption{
			{Type: discordgo.ApplicationCommandOptionUser, Name: "user", Description: "User to hard ban", Required: true},
			{Type: discordgo.ApplicationCommandOptionString, Name: "reason", Description: "Ban reason", Required: false},
		},
	},
	{
		Name:        "softban",
		Description: "Ban then immediately unban a user, deleting their recent messages",
		Options: []*discordgo.ApplicationCommandOption{
			{Type: discordgo.ApplicationCommandOptionUser, Name: "user", Description: "User to soft ban", Required: true},
			{Type: discordgo.ApplicationCommandOptionString, Name: "reason", Description: "Ban reason", Required: false},
		},
	},
	{
		Name:        "kick",
		Description: "Kick a user from the server",
		Options: []*discordgo.ApplicationCommandOption{
			{Type: discordgo.ApplicationCommandOptionUser, Name: "user", Description: "User to kick", Required: true},
			{Type: discordgo.ApplicationCommandOptionString, Name: "reason", Description: "Kick reason", Required: false},
		},
	},
	{
		Name:        "purge",
		Description: "Delete a number of recent messages, optionally from a specific user",
		Options: []*discordgo.ApplicationCommandOption{
			{Type: discordgo.ApplicationCommandOptionInteger, Name: "count", Description: "Number of messages to delete (1-100)", Required: true},
			{Type: discordgo.ApplicationCommandOptionUser, Name: "user", Description: "Only delete messages from this user", Required: false},
		},
	},
	{
		Name:        "nuke",
		Description: "Delete all messages by cloning the current channel",
	},
	{
		Name:        "timeout",
		Description: "Timeout a user for a duration",
		Options: []*discordgo.ApplicationCommandOption{
			{Type: discordgo.ApplicationCommandOptionString, Name: "user", Description: "User to timeout", Required: true},
			{Type: discordgo.ApplicationCommandOptionString, Name: "duration", Description: "Duration, e.g. 30m, 2h, 1d", Required: true},
			{Type: discordgo.ApplicationCommandOptionString, Name: "reason", Description: "Timeout reason", Required: false},
		},
	},
	{
		Name:        "warn",
		Description: "Warn a user",
		Options: []*discordgo.ApplicationCommandOption{
			{Type: discordgo.ApplicationCommandOptionString, Name: "user", Description: "User to warn", Required: true},
			{Type: discordgo.ApplicationCommandOptionString, Name: "reason", Description: "Warning reason", Required: true},
		},
	},
	{
		Name:        "history",
		Description: "Show a user's moderation history",
		Options: []*discordgo.ApplicationCommandOption{
			{Type: discordgo.ApplicationCommandOptionString, Name: "user", Description: "User to look up", Required: true},
		},
	},
	{
		Name:        "audit",
		Description: "Show recent moderation actions",
		Options: []*discordgo.ApplicationCommandOption{
			{Type: discordgo.ApplicationCommandOptionInteger, Name: "limit", Description: "Number of entries (1-50)", Required: false},
			{Type: discordgo.ApplicationCommandOptionUser, Name: "actor", Description: "Only show actions by this user", Required: false},
		},
	},
	{
		Name:        "role",
		Description: "Manage roles",
		Options: []*discordgo.ApplicationCommandOption{
			{Type: discordgo.ApplicationCommandOptionSubCommand, Name: "add", Description: "Add a role to a user", Options: []*discordgo.ApplicationCommandOption{
				{Type: discordgo.ApplicationCommandOptionUser, Name: "user", Description: "User to give the role", Required: true},
				{Type: discordgo.ApplicationCommandOptionString, Name: "role", Description: "Role name or ID", Required: true},
			}},
			{Type: discordgo.ApplicationCommandOptionSubCommand, Name: "create", Description: "Create a role", Options: []*discordgo.ApplicationCommandOption{
				{Type: discordgo.ApplicationCommandOptionString, Name: "name", Description: "Role name", Required: true},
			}},
			{Type: discordgo.ApplicationCommandOptionSubCommand, Name: "edit", Description: "Edit a role", Options: []*discordgo.ApplicationCommandOption{
				{Type: discordgo.ApplicationCommandOptionString, Name: "role", Description: "Role name or ID", Required: true},
				{Type: discordgo.ApplicationCommandOptionString, Name: "name", Description: "New name", Required: false},
				{Type: discordgo.ApplicationCommandOptionString, Name: "color", Description: "New hex color", Required: false},
			}},
			{Type: discordgo.ApplicationCommandOptionSubCommand, Name: "hoist", Description: "Toggle role hoisting", Options: []*discordgo.ApplicationCommandOption{
				{Type: discordgo.ApplicationCommandOptionString, Name: "role", Description: "Role name or ID", Required: true},
				{Type: discordgo.ApplicationCommandOptionBoolean, Name: "hoist", Description: "Hoist on/off", Required: true},
			}},
			{Type: discordgo.ApplicationCommandOptionSubCommand, Name: "member", Description: "List members with a role", Options: []*discordgo.ApplicationCommandOption{
				{Type: discordgo.ApplicationCommandOptionString, Name: "role", Description: "Role name or ID", Required: true},
			}},
		},
	},
	{
		Name:        "fn",
		Description: "Force set a user's nickname",
		Options: []*discordgo.ApplicationCommandOption{
			{Type: discordgo.ApplicationCommandOptionString, Name: "user", Description: "User to rename", Required: true},
			{Type: discordgo.ApplicationCommandOptionString, Name: "nickname", Description: "New nickname", Required: true},
		},
	},
	{
		Name:        "nick",
		Description: "Lock or unlock a user's nickname",
		Options: []*discordgo.ApplicationCommandOption{
			{Type: discordgo.ApplicationCommandOptionSubCommand, Name: "lock", Description: "Lock a user's nickname", Options: []*discordgo.ApplicationCommandOption{
				{Type: discordgo.ApplicationCommandOptionString, Name: "user", Description: "User to lock", Required: true},
			}},
			{Type: discordgo.ApplicationCommandOptionSubCommand, Name: "unlock", Description: "Unlock a user's nickname", Options: []*discordgo.ApplicationCommandOption{
				{Type: discordgo.ApplicationCommandOptionString, Name: "user", Description: "User to unlock", Required: true},
			}},
		},
	},
	{
		Name:        "jail",
		Description: "Jail a user, removing their roles and restricting them",
		Options: []*discordgo.ApplicationCommandOption{
			{Type: discordgo.ApplicationCommandOptionString, Name: "user", Description: "User to jail", Required: true},
		},
	},
	{
		Name:        "unjail",
		Description: "Unjail a user, restoring their previous roles",
		Options: []*discordgo.ApplicationCommandOption{
			{Type: discordgo.ApplicationCommandOptionString, Name: "user", Description: "User to unjail", Required: true},
		},
	},
	{
		Name:        "staffstrip",
		Description: "Remove all staff roles from a user",
		Options: []*discordgo.ApplicationCommandOption{
			{Type: discordgo.ApplicationCommandOptionString, Name: "user", Description: "User to strip roles from", Required: true},
		},
	},
	{
		Name:        "hide",
		Description: "Hide the current channel from @everyone",
	},
	{
		Name:        "reveal",
		Description: "Make the current channel visible to @everyone again",
	},
	{
		Name:        "lockdown",
		Description: "Lock the current channel for @everyone",
	},
	{
		Name:        "nsfw",
		Description: "Mark the current channel as NSFW",
	},
	{
		Name:        "sfw",
		Description: "Unmark the current channel as NSFW",
	},
	{
		Name:        "slowmode",
		Description: "Set a slowmode on the current channel",
		Options: []*discordgo.ApplicationCommandOption{
			{Type: discordgo.ApplicationCommandOptionInteger, Name: "seconds", Description: "Slowmode in seconds (0 to disable)", Required: true},
		},
	},
	{
		Name:        "topic",
		Description: "Set the topic of the current channel",
		Options: []*discordgo.ApplicationCommandOption{
			{Type: discordgo.ApplicationCommandOptionString, Name: "topic", Description: "New channel topic", Required: true},
		},
	},
	{
		Name:        "denyperm",
		Description: "Deny a permission to a user or role in a channel",
		Options: []*discordgo.ApplicationCommandOption{
			{Type: discordgo.ApplicationCommandOptionUser, Name: "user", Description: "User to deny permission to", Required: false},
			{Type: discordgo.ApplicationCommandOptionString, Name: "role", Description: "Role to deny permission to", Required: false},
			{Type: discordgo.ApplicationCommandOptionString, Name: "permission", Description: "Permission to deny (view, send, attach, embed, reaction, voice)", Required: true},
			{Type: discordgo.ApplicationCommandOptionString, Name: "channel", Description: "Channel to apply in (default current)", Required: false},
		},
	},
	{
		Name:        "imute",
		Description: "Prevent a user from sending images in this channel",
		Options: []*discordgo.ApplicationCommandOption{
			{Type: discordgo.ApplicationCommandOptionString, Name: "user", Description: "User to image mute", Required: true},
		},
	},
	{
		Name:        "gifmute",
		Description: "Prevent a user from sending images, GIFs, and embeds in this channel",
		Options: []*discordgo.ApplicationCommandOption{
			{Type: discordgo.ApplicationCommandOptionString, Name: "user", Description: "User to GIF mute", Required: true},
		},
	},
}

func init() {
	SlashCommands = append(SlashCommands, emotionSlashCommands...)
	SlashCommands = append(SlashCommands, actionSlashCommands...)
}

func userOption(required bool) *discordgo.ApplicationCommandOption {
	return &discordgo.ApplicationCommandOption{
		Type:        discordgo.ApplicationCommandOptionUser,
		Name:        "user",
		Description: "The user (defaults to yourself)",
		Required:    required,
	}
}
