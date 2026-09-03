import type { ComponentType } from "react"
import {
  Home,
  Info,
  Search,
  Server,
  Wrench,
  Image,
  Smile,
  Shield,
  Heart,
  Globe,
} from "lucide-react"
import { SiGithub, SiReddit, SiRoblox, SiKick } from "react-icons/si"

export type Icon = ComponentType<{ className?: string }>

export type Command = {
  name: string
  description: string
  icon?: Icon
}

export type CommandCategory = {
  id: string
  name: string
  description: string
  icon: Icon
  commands: Command[]
}

export const commandCategories: CommandCategory[] = [
  {
    id: "general",
    name: "General",
    description: "Everyday commands to check and interact with the bot.",
    icon: Home,
    commands: [
      { name: "ping", description: "Check bot latency" },
      {
        name: "test",
        description: "Ping-pong check that the bot is responsive",
      },
      { name: "bug", description: "Open a form to report a bot bug" },
      { name: "help", description: "Show the interactive command menu" },
      { name: "tldr", description: "Get a brief how-to for a command" },
      { name: "userinfo", description: "Get info about a user" },
    ],
  },
  {
    id: "information",
    name: "Information",
    description: "Profile, server, channel, and weather details on demand.",
    icon: Info,
    commands: [
      {
        name: "avatar",
        description: "Show a user's full-resolution profile picture",
      },
      { name: "banner", description: "Show a user's profile banner" },
      {
        name: "botinfo",
        description: "Show bot info (version, uptime, latency)",
      },
      { name: "channelinfo", description: "Get info about a channel" },
      {
        name: "guild stats / icon / banner / splash",
        description: "Server stats, icon, banner, or splash",
      },
      {
        name: "roles",
        description: "List all server roles with member counts",
      },
      { name: "emojis", description: "List all custom emojis in the server" },
      {
        name: "stickers",
        description: "List all custom stickers in the server",
      },
      { name: "bans", description: "List all active bans in the server" },
      { name: "boomer", description: "Detect potential time-traveler users" },
      { name: "perms", description: "Show a user's effective permissions" },
      {
        name: "tz",
        description: "Show the current local time in a place (e.g. Tokyo)",
      },
      {
        name: "urbandictionary",
        description: "Search Urban Dictionary for a term",
      },
      {
        name: "weather",
        description: "Current weather, forecast, and humidity for a place",
      },
    ],
  },
  {
    id: "search",
    name: "Search",
    description: "Search the web, chat history, images, and GIFs.",
    icon: Search,
    commands: [
      {
        name: "ddg",
        description: "Fetch quick search results from DuckDuckGo",
      },
      { name: "search", description: "Search chat history for a message" },
      {
        name: "isearch",
        description: "Reverse image search: find similar images & sources",
      },
      { name: "gifsearch", description: "Post a relevant GIF from Giphy" },
    ],
  },
  {
    id: "messages",
    name: "Server & Messages",
    description: "Manage messages, emojis, stickers, and pins.",
    icon: Server,
    commands: [
      { name: "pinglist", description: "Show pings for a user or role" },
      {
        name: "emoji",
        description: "Manage server emojis (add / list / remove / steal)",
      },
      {
        name: "sticker",
        description: "Manage server stickers (add / steal / remove)",
      },
      { name: "pin", description: "Pin a message" },
      { name: "unpin", description: "Unpin a message" },
      { name: "quote", description: "Quote a message as a styled embed" },
    ],
  },
  {
    id: "utilities",
    name: "Utilities",
    description: "Handy helpers for text, translation, and reminders.",
    icon: Wrench,
    commands: [
      { name: "echo", description: "Repeat back a message" },
      {
        name: "translate",
        description: "Translate text into a target language",
      },
      { name: "reminder", description: "Set a delayed reminder ping" },
      { name: "caption", description: "Add meme text to an image" },
    ],
  },
  {
    id: "media",
    name: "Media & AI",
    description: "Convert media and fact-check claims.",
    icon: Image,
    commands: [
      { name: "img2gif", description: "Convert an image into a GIF" },
      {
        name: "vid2gif",
        description: "Convert a video into a GIF (≤25MB, clips to 10s)",
      },
      { name: "autogif", description: "Convert any media into a GIF" },
      {
        name: "factcheck",
        description: "Fact-check a claim against searchable ratings",
      },
    ],
  },
  {
    id: "fun",
    name: "Fun",
    description: "Silly readings, polls, birthdays, and more.",
    icon: Smile,
    commands: [
      { name: "howgay", description: "Measure how gay a member is" },
      { name: "howautism", description: "Measure how autistic a member is" },
      { name: "howlesbian", description: "Measure how lesbian a member is" },
      {
        name: "howsimp",
        description: "Measure how much of a simp a member is",
      },
      { name: "pp", description: "Measure a member's pp length" },
      { name: "puh", description: "Check the puh tightness" },
      { name: "iq", description: "Measure a member's IQ" },
      { name: "bitches", description: "See how many bitches a member has" },
      {
        name: "choose",
        description: "Pick an option from a comma-separated list",
      },
      {
        name: "compress",
        description: "Compress an image until it's barely legible",
      },
      { name: "poll", description: "Create a reaction poll (2-10 options)" },
      { name: "quickpoll", description: "Create a Yes/No poll" },
      {
        name: "birthday",
        description: "Manage server birthdays (add / list / celebrate)",
      },
      {
        name: "diary",
        description: "Private per-user diary (add / view / delete)",
      },
      { name: "ttys", description: "Bot talks to itself until someone speaks" },
    ],
  },
  {
    id: "moderation",
    name: "Moderation",
    description:
      "Keep your server clean and safe (punishments, roles, channels).",
    icon: Shield,
    commands: [
      { name: "ban", description: "Ban a user from the server" },
      {
        name: "hardban",
        description: "Ban and delete a user's recent messages",
      },
      {
        name: "softban",
        description: "Ban then immediately unban, deleting recent messages",
      },
      { name: "kick", description: "Kick a user from the server" },
      { name: "timeout", description: "Timeout a user (e.g. 30m, 2h, 1d)" },
      { name: "warn", description: "Warn a user" },
      { name: "history", description: "Show a user's moderation history" },
      { name: "audit", description: "Show recent moderation actions" },
      { name: "jail", description: "Jail a user, removing their roles" },
      { name: "unjail", description: "Unjail a user, restoring their roles" },
      {
        name: "role",
        description: "Manage roles (add / create / edit / hoist / member)",
      },
      { name: "fn", description: "Force set a user's nickname" },
      { name: "nick", description: "Lock or unlock a user's nickname" },
      { name: "staffstrip", description: "Remove all staff roles from a user" },
      {
        name: "purge",
        description: "Delete recent messages, optionally only from a user",
      },
      {
        name: "nuke",
        description: "Delete all messages by cloning the current channel",
      },
      { name: "hide", description: "Hide the current channel from @everyone" },
      { name: "reveal", description: "Make the current channel visible again" },
      {
        name: "lockdown",
        description: "Lock the current channel for @everyone",
      },
      { name: "nsfw", description: "Mark the current channel as NSFW" },
      { name: "sfw", description: "Unmark the current channel as NSFW" },
      {
        name: "slowmode",
        description: "Set a slowmode on the current channel",
      },
      { name: "topic", description: "Set the topic of the current channel" },
      { name: "denyperm", description: "Deny a permission to a user or role" },
      {
        name: "imute",
        description: "Prevent a user from sending images in this channel",
      },
      {
        name: "gifmute",
        description: "Prevent a user from images, GIFs, and embeds",
      },
    ],
  },
  {
    id: "roleplay",
    name: "Roleplay",
    description: "React with emotions and perform anime-style actions as GIFs.",
    icon: Heart,
    commands: [
      { name: "angry", description: "Post a GIF expressing being angry" },
      {
        name: "depressed",
        description: "Post a GIF expressing being depressed",
      },
      { name: "excited", description: "Post a GIF expressing being excited" },
      { name: "happy", description: "Post a GIF expressing being happy" },
      { name: "horny", description: "Post a GIF expressing being horny" },
      { name: "inlove", description: "Post a GIF expressing being in love" },
      { name: "sad", description: "Post a GIF expressing being sad" },
      { name: "shy", description: "Post a GIF expressing being shy" },
      { name: "baka", description: "Call someone a baka (anime GIF)" },
      { name: "bite", description: "Bite a user (anime GIF)" },
      { name: "cry", description: "Cry at a user (anime GIF)" },
      { name: "dap", description: "Dap up a user (anime GIF)" },
      { name: "eat", description: "Munch on a user (anime GIF)" },
      { name: "facepalm", description: "Facepalm at a user (anime GIF)" },
      { name: "feed", description: "Feed a user (anime GIF)" },
      { name: "handhold", description: "Hold hands with a user (anime GIF)" },
      { name: "kiss", description: "Kiss a user (anime GIF)" },
      { name: "laugh", description: "Laugh at a user (anime GIF)" },
      { name: "nod", description: "Nod at a user (anime GIF)" },
      { name: "nutkick", description: "Nutkick a user (anime GIF)" },
      { name: "pat", description: "Pat a user (anime GIF)" },
      { name: "peck", description: "Peck a user (anime GIF)" },
      { name: "poke", description: "Poke a user (anime GIF)" },
      { name: "punch", description: "Punch a user (anime GIF)" },
      { name: "run", description: "Run away from a user (anime GIF)" },
      { name: "shoot", description: "Shoot a user (anime GIF)" },
      { name: "shrug", description: "Shrug at a user (anime GIF)" },
      { name: "slap", description: "Slap a user (anime GIF)" },
      { name: "spank", description: "Spank a user (anime GIF)" },
      { name: "stab", description: "Stab a user (anime GIF)" },
      { name: "think", description: "Think about a user (anime GIF)" },
      { name: "tickle", description: "Tickle a user (anime GIF)" },
    ],
  },
  {
    id: "social",
    name: "Social",
    description: "Public profile lookups and server social mini-games.",
    icon: Globe,
    commands: [
      {
        name: "bio github",
        description: "Look up a GitHub profile's public details",
        icon: SiGithub,
      },
      {
        name: "bio reddit",
        description: "Look up a Reddit user's public profile",
        icon: SiReddit,
      },
      {
        name: "bio roblox",
        description: "Look up a Roblox profile's details",
        icon: SiRoblox,
      },
      {
        name: "bio kick",
        description: "Look up a Kick streamer's live status",
        icon: SiKick,
      },
      { name: "ship", description: "Calculate romance compatibility" },
      {
        name: "colors avatar",
        description: "Extract dominant colours from an avatar",
      },
      { name: "lurk", description: "Toggle yourself in/out of lurking mode" },
      { name: "lurkers", description: "Show who is currently lurking" },
      { name: "spark", description: "Spark the blunt before you can smoke" },
      {
        name: "smoke",
        description: "Take a hit off the blunt (spark it first)",
      },
      { name: "hits", description: "Show everyone's blunt hit count" },
      { name: "vape", description: "Hit, configure, or check the server vape" },
    ],
  },
]
