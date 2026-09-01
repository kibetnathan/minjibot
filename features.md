# MinjiBot — Feature Reference

Every feature ships as both a **Slash Command** (`/`) and a **Classic Prefix Command** (`,`).

## 1. General

| Command    | Description                                                         |
| ---------- | ------------------------------------------------------------------- |
| `bug`      | Opens a modal or prompt to report a bot bug directly to developers. |
| `echo`     | Repeats the user's input text back into the channel.                |
| `ping`     | Displays real-time API latency and WebSocket response times.        |
| `test`     | Ping-pong functionality check to confirm bot responsiveness.        |
| `userinfo` | Shows detailed account data (creation date, join date, roles).      |

## 2. Moderation

| Command       | Description                                                           |
| ------------- | --------------------------------------------------------------------- |
| `audit`       | Fetches recent server audit log entries with parameter filters.       |
| `ban`         | Bans a user from the server with an optional reason.                  |
| `denyperm`    | Explicitly denies a specific permission for a user/role in a channel. |
| `fn`          | Force Nickname: locks a forced nickname onto a targeted user.         |
| `gifmute`     | Denies a targeted user permission to embed links or attach media.     |
| `hardban`     | Bans a user and wipes all their message history instantly.            |
| `hide`        | Hides a channel by denying `View Channel` perms for `@everyone`.      |
| `history`     | Pulls a user's entire server punishment and infraction log.           |
| `imute`       | Image Mute: revokes image/file upload permissions for a user.         |
| `jail`        | Strips user roles and isolates them in a designated jail channel.     |
| `kick`        | Removes a user from the server without issuing a permanent ban.       |
| `lockdown`    | Revokes `Send Messages` for `@everyone` across channels.              |
| `nick lock`   | Prevents a targeted user from changing their server nickname.         |
| `nsfw`        | Marks the current text channel as NSFW.                               |
| `nuke`        | Clones and deletes the current channel to clear all past messages.    |
| `purge`       | Deletes a specified batch of recent messages (`1-100`).               |
| `reveal`      | Unhides a previously hidden channel by resetting permissions.         |
| `role add`    | Assigns a designated role to a target member.                         |
| `role create` | Creates a new server role with custom parameters.                     |
| `role edit`   | Modifies an existing role's name, color, or permissions.              |
| `role hoist`  | Toggles whether a role displays separately in the member list.        |
| `role member` | Lists every member currently assigned to a specific role.             |
| `sfw`         | Removes the NSFW flag from the current text channel.                  |
| `slowmode`    | Sets a message rate limit delay for the current channel.              |
| `softban`     | Bans and immediately unbans a user to wipe their recent messages.     |
| `staffstrip`  | Instantly strips all administrative/mod roles from a user.            |
| `timeout`     | Temporarily mutes/isolates a member using Discord's timeout feature.  |
| `topic`       | Updates the description/topic text of the current channel.            |
| `warn`        | Issues a formal recorded warning to a target member.                  |

## 3. Utility

| Command          | Description                                                   |
| ---------------- | ------------------------------------------------------------- |
| `autogif`        | Automatically fetches and converts relevant media to GIFs.    |
| `caption`        | Generates a meme image by attaching custom text to an image.  |
| `compress`       | Pixelates an image to a low-quality but still readable mosaic. |
| `ddg`            | Fetches quick search results directly from DuckDuckGo.        |
| `emoji add`      | Uploads an emoji to the server via URL or attachment.         |
| `emoji add many` | Bulk-uploads multiple custom emojis at once.                  |
| `emoji enlarge`  | Renders a high-resolution, full-size version of an emoji.     |
| `emoji list`     | Displays an inventory of all custom emojis in the server.     |
| `emoji remove`   | Deletes a target custom emoji from the server.                |
| `emoji steal`    | Copies an emoji from another server using its message format. |
| `gifsearch`      | Searches Tenor/Giphy and posts a relevant GIF.                |
| `img2gif`        | Converts an image file or static link into an animated GIF.   |
| `isearch`        | Reverse image search: finds similar images and the pages they're from. |
| `pin`            | Pins a target message to the channel's pinned list.           |
| `quote`          | Converts a specific message link into a stylized embed.       |
| `reminder`       | Sets a delayed system ping for a user-specified task/time.    |
| `search`         | Searches for a specific message in chat history               |
| `sticker add`    | Uploads a custom sticker to the server.                       |
| `sticker remove` | Removes a custom sticker from the server inventory.           |
| `translate`      | Translates provided text into a target language.              |
| `unpin`          | Removes a message from the channel's pinned list.             |
| `vid2gif`        | Trims and converts an uploaded video file into a GIF.         |
| `pinglist`       | Returns all the pings for a certain user/role                 |
| `factcheck`      | Factchecks someones message                                   |

## 4. Fun

| Command              | Description                                                      |
| -------------------- | ---------------------------------------------------------------- |
| `birthday add`       | Saves a user's birthday to the bot's database.                   |
| `birthday celebrate` | Triggers an immediate birthday celebration message/embed.        |
| `birthday channel`   | Sets the channel for automated birthday notices.                 |
| `birthday list`      | Displays upcoming server birthdays sorted by date.               |
| `birthday role`      | Configures the temporary role awarded on a member's birthday.    |
| `bitches`            | Gives a random number for how many bitches a specific member has |
| `choose`             | Picks an option from a comma-separated list.                     |
| `colors avatar`      | Extracts dominant hex colors from an avatar.                     |
| `diary add`          | Saves a private journal entry into the user's database.          |
| `diary delete`       | Removes a specific saved diary entry by ID.                      |
| `diary view`         | Shows past diary entries in a private embed.                     |
| `hits`               | Displays everyone's blunt hit count.                             |
| `howautism`          | Randomized fun percentage reading.                               |
| `howgay`             | Randomized fun percentage reading.                               |
| `howlesbian`         | Randomized fun percentage reading.                               |
| `howsimp`            | Randomized fun percentage reading.                               |
| `IQ`                 | Generates a fun, randomized IQ score readout.                    |
| `lurk`               | Toggle lurk mode: your messages auto-delete after 2s and you auto-unlurk after 1 hour. |
| `lurkers`            | Shows who is currently lurking.                                 |
| `poll`               | Creates an interactive reaction- or button-based poll embed.     |
| `pp`                 | Generates a randomized pp length                                 |
| `puh`                | Generates random puh tightness percentage.                       |
| `quickpoll`          | Instantly creates a simple Yes/No reaction poll.                 |
| `ship`               | Calculates a romance compatibility score between two users.      |
| `smoke`              | Take a hit off the blunt (spark it first; grabbing resets the spark). |
| `spark`              | Light the blunt before anyone can smoke.                         |
| `vape`               | Hit the vape.                                                    |
| `vape flavor`        | Sets or changes your virtual vapeflavor.                         |
| `vape hits`          | Displays total recorded virtual vape hits.                       |
| `vape steal`         | Steal the vape.                                                  |

## 5. Roleplay (RP)

### Emotions

| Command     | Description                         |
| ----------- | ----------------------------------- |
| `angry`     | Animated GIF expressing anger.      |
| `depressed` | Animated GIF expressing depression. |
| `excited`   | Animated GIF expressing excitement. |
| `happy`     | Animated GIF expressing happiness.  |
| `horny`     | Animated GIF expressing horniness.  |
| `inlove`    | Animated GIF expressing love.       |
| `sad`       | Animated GIF expressing sadness.    |
| `shy`       | Animated GIF expressing shyness.    |

### Actions

All action commands trigger a GIF of the action performed towards another user (e.g. `/slap @user`).

| Command | Command    | Command | Command    |
| ------- | ---------- | ------- | ---------- |
| `baka`  | `bite`     | `cry`   | `dap`      |
| `eat`   | `facepalm` | `feed`  | `handhold` |
| `kiss`  | `laugh`    | `nod`   | `nutkick`  |
| `pat`   | `peck`     | `poke`  | `punch`    |
| `run`   | `shoot`    | `shrug` | `slap`     |
| `spank` | `stab`     | `think` | `tickle`   |

### Character

| Command             | Description                                                                        |
| ------------------- | ---------------------------------------------------------------------------------- |
| `character act`     | Displays a specific action/narration performable by the created persona.           |
| `character autosay` | Toggles automatic webhook proxying for all messages sent by the user.              |
| `character create`  | Creates a custom roleplay persona stored in the database.                          |
| `character say`     | Sends a message styled under the active character's name and avatar (via Webhook). |

## 6. Information

| Command                 | Description                                                  |
| ----------------------- | ------------------------------------------------------------ |
| `avatar`                | Fetches a targeted user's full-resolution profile picture.   |
| `avatarhistory disable` | Disables tracking and purges stored avatar history.          |
| `avatarhistory enable`  | Enables automatic logging of past avatar changes.            |
| `bans`                  | Paginated list of all currently banned server members.       |
| `banner`                | Displays the profile banner image of a specified user.       |
| `boomer`                | Evaluates account age into a lighthearted "boomer score."    |
| `botinfo`               | System uptime, RAM usage, ping, and bot architecture.        |
| `channelinfo`           | Creation date, ID, slowmode, and topic of a channel.         |
| `emojis`                | Inventory list of every custom emoji in the server.          |
| `guild banner`          | High-resolution server banner image.                         |
| `guild icon`            | Full-resolution icon of the server.                          |
| `guild splash`          | Invite splash background image of the server.                |
| `guild stats`           | Member counts, online status ratios, and boost levels.       |
| `help`                  | Interactive, categorized list of all bot commands.           |
| `perms`                 | Explicit server and channel permissions for a targeted user. |
| `roles`                 | All roles in the server with their member counts.            |
| `stickers`              | Complete list of custom stickers in the server.              |
| `tldr`                  | Summarizes recent channel chat activity.                     |
| `tz`                    | Current local time or timezone conversion.                   |
| `urbandictionary`       | Top slang definitions from Urban Dictionary.                 |
| `weather`               | Current weather, forecasts, and humidity for a location.     |

## 7. Social Lookup

| Command             | Description                                           |
| ------------------- | ----------------------------------------------------- |
| `github repository` | Searches GitHub repositories and displays key stats.  |
| `guns`              | Looks up profiles on Guns.lol link-in-bio platform.   |
| `instagram`         | Public profile stats, bios, and follower counts.      |
| `kick`              | Live stream status and profile info from Kick.com.    |
| `linkedin`          | Basic public profile or company information.          |
| `linktree`          | Links hosted on a user's Linktree URL.                |
| `pinterest`         | Pinterest profiles, public boards, and top pins.      |
| `reddit`            | User karma, top posts, or subreddit information.      |
| `roblox`            | User details, rap value, badges, and account age.     |
| `tiktok`            | Public profile stats, likes, and follower counts.     |
| `twitch`            | Live status, game category, and channel stats.        |
| `twitter`           | Public X/Twitter profile details and recent activity. |
| `youtube`           | Channels, subscriber counts, or video stats.          |
