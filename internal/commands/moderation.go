package commands

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/kibetnathan/minjibot/internal/ports/dto"
)

const modColor = 0x5865F2
const modErrorColor = 0xED4245

// modPerformance is the set of permissions that grant moderation power. Any of
// these allows a user to run destructive moderation commands.
const defaultModPerm = discordgo.PermissionManageMessages |
	discordgo.PermissionKickMembers |
	discordgo.PermissionBanMembers |
	discordgo.PermissionManageRoles |
	discordgo.PermissionManageChannels |
	discordgo.PermissionAdministrator

// effectiveModPerms returns the union of server-wide permissions held by a
// user (everyone role + their roles). Channel-level perms are intentionally
// excluded so moderation cannot be granted by a single channel override.
func effectiveModPerms(s *discordgo.Session, guildID, userID string) (int64, error) {
	member, err := s.GuildMember(guildID, userID)
	if err != nil {
		return 0, err
	}
	roles, err := s.GuildRoles(guildID)
	if err != nil {
		return 0, err
	}

	byID := make(map[string]*discordgo.Role, len(roles))
	for _, r := range roles {
		byID[r.ID] = r
	}

	var perms int64
	if everyone, ok := byID[guildID]; ok {
		perms |= everyone.Permissions
	}
	for _, rid := range member.Roles {
		if r, ok := byID[rid]; ok {
			perms |= r.Permissions
		}
	}
	return perms, nil
}

// hasAnyPerm reports whether perms contains at least one of the given
// permissions (or Administrator).
func hasAnyPerm(perms int64, wanted ...int64) bool {
	if perms&discordgo.PermissionAdministrator != 0 {
		return true
	}
	for _, w := range wanted {
		if perms&w != 0 {
			return true
		}
	}
	return false
}

// requireModerator checks that the user holds moderation permissions. For
// slash commands it prefers the per-member permissions Discord attaches to the
// interaction; otherwise it recomputes from guild roles. Returns (true, "") if
// allowed, otherwise (false, errorMessage).
func requireModerator(s *discordgo.Session, guildID, userID string, m *discordgo.Member) (bool, string) {
	perms := modPermsFromMember(m)
	ok := hasAnyPerm(perms, defaultModPerm)
	return ok, "You need a moderation permission (Manage Messages, Kick/Ban Members, Manage Roles/Channels, or Administrator) to use this command."
}

func modPermsFromMember(m *discordgo.Member) int64 {
	if m == nil {
		return 0
	}
	return m.Permissions
}

// resolveTargetUser returns the userID and display name for a target. It
// accepts a plain ID string or a raw <@id> / <@!id> mention. Used by moderation
// commands that take a "user" argument.
func resolveTargetUser(s *discordgo.Session, guildID, raw string) (string, string, error) {
	id := parseMentionID(raw)
	member, err := s.GuildMember(guildID, id)
	if err != nil {
		return "", "", err
	}
	name := member.User.Username
	if member.Nick != "" {
		name = member.Nick
	}
	return id, name, nil
}

// parseMentionID strips Discord user-mention wrappers from raw.
func parseMentionID(raw string) string {
	out := raw
	if len(out) >= 3 && out[0] == '<' && out[len(out)-1] == '>' {
		inner := out[1 : len(out)-1]
		if len(inner) > 0 && inner[0] == '@' {
			if len(inner) > 1 && inner[1] == '!' {
				inner = inner[2:]
			} else {
				inner = inner[1:]
			}
			if len(inner) > 0 && inner[0] == '&' {
				inner = inner[1:]
			}
			out = inner
		}
	}
	return out
}

// resolveTargetRole resolves a role ID or mention to a role in the guild.
func resolveTargetRole(s *discordgo.Session, guildID, raw string) (*discordgo.Role, error) {
	id := parseMentionID(raw)
	roles, err := s.GuildRoles(guildID)
	if err != nil {
		return nil, err
	}
	for _, r := range roles {
		if r.ID == id {
			return r, nil
		}
	}
	return nil, fmt.Errorf("role `%s` not found", raw)
}

// auditAction writes a moderation action into the audit log table when an
// audit repository is available. Metadata is optional and JSON-encoded.
func auditAction(h *CommandHandler, ctx context.Context, guildID, action, actorID, targetID string, metadata map[string]any) {
	if h == nil || h.AuditRepo == nil {
		return
	}
	var meta []byte
	if metadata != nil {
		if b, err := json.Marshal(metadata); err == nil {
			meta = b
		}
	}
	if _, err := h.AuditRepo.Create(ctx, dto.CreateAuditLogParams{
		GuildID:  guildID,
		Action:   action,
		ActorID:  actorID,
		TargetID: targetID,
		Metadata: meta,
	}); err != nil {
		return
	}
}

func modErrorEmbed(title, msg string) *discordgo.MessageEmbed {
	return &discordgo.MessageEmbed{
		Color:       modErrorColor,
		Title:       title,
		Description: msg,
	}
}

func modSuccessEmbed(title, msg string) *discordgo.MessageEmbed {
	return &discordgo.MessageEmbed{
		Color:       modColor,
		Title:       title,
		Description: msg,
	}
}

// sendModError replies to a message-command channel with an error embed.
func sendModError(s *discordgo.Session, channelID, title, msg string) error {
	_, err := s.ChannelMessageSendEmbed(channelID, modErrorEmbed(title, msg))
	return err
}

// plural returns "s" for counts != 1, "" otherwise.
func plural(n int) string {
	if n == 1 {
		return ""
	}
	return "s"
}

// joinArgs joins words into a single space-separated string.
func joinArgs(args []string) string {
	return strings.Join(args, " ")
}

// OptInt reads an integer-valued interaction option, returning def when absent.
func OptInt(opts map[string]*discordgo.ApplicationCommandInteractionDataOption, name string, def int) int {
	if o, ok := opts[name]; ok && o != nil {
		if v, ok := o.Value.(float64); ok {
			return int(v)
		}
	}
	return def
}

// OptBool reads a boolean-valued interaction option.
func OptBool(opts map[string]*discordgo.ApplicationCommandInteractionDataOption, name string) bool {
	if o, ok := opts[name]; ok && o != nil {
		if v, ok := o.Value.(bool); ok {
			return v
		}
	}
	return false
}

// OptUser reads a user-valued interaction option, returning its ID.
func OptUser(opts map[string]*discordgo.ApplicationCommandInteractionDataOption, name string) string {
	if o, ok := opts[name]; ok && o != nil {
		if v, ok := o.Value.(string); ok {
			return parseMentionID(v)
		}
	}
	return ""
}

// requireModForMessage is a convenience for prefix commands: it computes the
// invoking user's server perms and returns a friendly error embed if they lack
// moderation permission. Callers should return immediately when ok is false.
func requireModForMessage(s *discordgo.Session, m *discordgo.MessageCreate, title string) (bool, error) {
	return requireModForUser(s, m.GuildID, m.Author.ID, m.ChannelID, title)
}

// requireModForUser is requireModForMessage without needing the full message;
// used by helpers that operate on an arbitrary guild/channel.
func requireModForUser(s *discordgo.Session, guildID, userID, channelID, title string) (bool, error) {
	perms, err := effectiveModPerms(s, guildID, userID)
	if err != nil {
		return false, sendModError(s, channelID, title, fmt.Sprintf("Could not check permissions: %s", err))
	}
	if !hasAnyPerm(perms, defaultModPerm) {
		return false, sendModError(s, channelID, title, "You need a moderation permission (Manage Messages, Kick/Ban Members, Manage Roles/Channels, or Administrator) to use this command.")
	}
	return true, nil
}
