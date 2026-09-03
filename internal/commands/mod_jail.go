package commands

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/kibetnathan/minjibot/internal/ports/dto"
)

const jailRoleName = "Jail"
const jailRoleKey = "JAIL_SENTINEL"
const jailStoredRolesKey = "JAIL_ORIGINAL_ROLES"

// jailMessage: -jail <user>  (strip roles and isolate the user in a jail role)
func jailMessageCommandHandler(h *CommandHandler, s *discordgo.Session, m *discordgo.MessageCreate, args []string) error {
	if ok, err := requireModForMessage(s, m, "Jail"); !ok {
		return err
	}
	if len(args) == 0 {
		return sendModError(s, m.ChannelID, "Jail", "Usage: `-jail <user>`")
	}
	targetID, name, err := resolveTargetUser(s, m.GuildID, args[0])
	if err != nil {
		return sendModError(s, m.ChannelID, "Jail", fmt.Sprintf("Could not find that user: %s", err))
	}
	if err := jailUser(h, s, m.GuildID, targetID, m.Author.ID); err != nil {
		return sendModError(s, m.ChannelID, "Jail", err.Error())
	}
	_, err = s.ChannelMessageSendEmbed(m.ChannelID, modSuccessEmbed("Jail", fmt.Sprintf("Jailed **%s** (roles stripped, moved to jail role).", name)))
	return err
}

func jailSlashCommandHandler(h *CommandHandler, s *discordgo.Session, i *discordgo.InteractionCreate) error {
	ok, msg := requireModerator(i.Member)
	if !ok {
		return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{Embeds: []*discordgo.MessageEmbed{modErrorEmbed("Jail", msg)}},
		})
	}
	opts := OptionMap(i.ApplicationCommandData().Options)
	raw := OptString(opts, "user")
	targetID, name, err := resolveTargetUser(s, i.GuildID, raw)
	if err != nil {
		return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{Embeds: []*discordgo.MessageEmbed{modErrorEmbed("Jail", fmt.Sprintf("Could not find that user: %s", err))}},
		})
	}
	if err := jailUser(h, s, i.GuildID, targetID, i.Member.User.ID); err != nil {
		return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{Embeds: []*discordgo.MessageEmbed{modErrorEmbed("Jail", err.Error())}},
		})
	}
	return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{Embeds: []*discordgo.MessageEmbed{modSuccessEmbed("Jail", fmt.Sprintf("Jailed **%s** (roles stripped, moved to jail role).", name))}},
	})
}

// jailUser strips a member's roles, stores them for restore, and assigns the
// jail role that denies chatting everywhere.
func jailUser(h *CommandHandler, s *discordgo.Session, guildID, targetID, actorID string) error {
	member, err := s.GuildMember(guildID, targetID)
	if err != nil {
		return fmt.Errorf("could not fetch member: %s", err)
	}

	// Find or create the jail role.
	jasRole, err := findOrCreateJailRole(s, guildID)
	if err != nil {
		return err
	}

	// Remember the member's roles so a later release can restore them.
	payload, _ := json.Marshal(map[string]any{"roles": member.Roles})
	if h != nil && h.PermRepo != nil {
		_, _ = h.PermRepo.Upsert(context.Background(), dto.UpsertUserPermissionParams{
			UserID:          targetID,
			GuildID:         guildID,
			Role:            jailRoleKey,
			PermissionsJSON: payload,
		})
	}

	// Remove every role (except the jail role and the everyone role).
	for _, rid := range member.Roles {
		if rid == jasRole.ID || rid == guildID {
			continue
		}
		_ = s.GuildMemberRoleRemove(guildID, targetID, rid)
	}
	// Assign the jail role.
	if err := s.GuildMemberRoleAdd(guildID, targetID, jasRole.ID); err != nil {
		return fmt.Errorf("failed to assign jail role: %s", err)
	}

	auditAction(h, context.Background(), guildID, "JAIL", actorID, targetID, nil)
	return nil
}

// unjailMessage: -unjail <user>  (release from jail and attempt role restore)
func unjailMessageCommandHandler(h *CommandHandler, s *discordgo.Session, m *discordgo.MessageCreate, args []string) error {
	if ok, err := requireModForMessage(s, m, "Unjail"); !ok {
		return err
	}
	if len(args) == 0 {
		return sendModError(s, m.ChannelID, "Unjail", "Usage: `-unjail <user>`")
	}
	targetID, name, err := resolveTargetUser(s, m.GuildID, args[0])
	if err != nil {
		return sendModError(s, m.ChannelID, "Unjail", fmt.Sprintf("Could not find that user: %s", err))
	}
	restored, err := unjailUser(h, s, m.GuildID, targetID, m.Author.ID)
	if err != nil {
		return sendModError(s, m.ChannelID, "Unjail", err.Error())
	}
	_, err = s.ChannelMessageSendEmbed(m.ChannelID, modSuccessEmbed("Unjail", fmt.Sprintf("Released **%s**. Original roles restored: %v.", name, restored)))
	return err
}

func unjailSlashCommandHandler(h *CommandHandler, s *discordgo.Session, i *discordgo.InteractionCreate) error {
	ok, msg := requireModerator(i.Member)
	if !ok {
		return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{Embeds: []*discordgo.MessageEmbed{modErrorEmbed("Unjail", msg)}},
		})
	}
	opts := OptionMap(i.ApplicationCommandData().Options)
	raw := OptString(opts, "user")
	targetID, name, err := resolveTargetUser(s, i.GuildID, raw)
	if err != nil {
		return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{Embeds: []*discordgo.MessageEmbed{modErrorEmbed("Unjail", fmt.Sprintf("Could not find that user: %s", err))}},
		})
	}
	restored, err := unjailUser(h, s, i.GuildID, targetID, i.Member.User.ID)
	if err != nil {
		return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{Embeds: []*discordgo.MessageEmbed{modErrorEmbed("Unjail", err.Error())}},
		})
	}
	return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{Embeds: []*discordgo.MessageEmbed{modSuccessEmbed("Unjail", fmt.Sprintf("Released **%s**. Original roles restored: %v.", name, restored))}},
	})
}

func unjailUser(h *CommandHandler, s *discordgo.Session, guildID, targetID, actorID string) (int, error) {
	jasRole, err := findOrCreateJailRole(s, guildID)
	if err != nil {
		return 0, err
	}
	// Remove the jail role.
	_ = s.GuildMemberRoleRemove(guildID, targetID, jasRole.ID)

	restored := 0
	if h != nil && h.PermRepo != nil {
		perm, gerr := h.PermRepo.Get(context.Background(), targetID, guildID, jailRoleKey)
		if gerr == nil {
			var data struct {
				Roles []string `json:"roles"`
			}
			if json.Unmarshal(perm.PermissionsJSON, &data) == nil {
				for _, rid := range data.Roles {
					if rid == jasRole.ID || rid == guildID {
						continue
					}
					if err := s.GuildMemberRoleAdd(guildID, targetID, rid); err == nil {
						restored++
					}
				}
			}
			_ = h.PermRepo.Delete(context.Background(), targetID, guildID, jailRoleKey)
		}
	}
	auditAction(h, context.Background(), guildID, "UNJAIL", actorID, targetID, nil)
	return restored, nil
}

// findOrCreateJailRole returns the guild's "Jail" role, creating it if absent,
// configured to deny chatting in every channel (given the channel for setup).
func findOrCreateJailRole(s *discordgo.Session, guildID string) (*discordgo.Role, error) {
	roles, err := s.GuildRoles(guildID)
	if err != nil {
		return nil, fmt.Errorf("could not fetch roles: %s", err)
	}
	for _, r := range roles {
		if r.Name == jailRoleName {
			return r, nil
		}
	}
	role, err := s.GuildRoleCreate(guildID, &discordgo.RoleParams{
		Name:        jailRoleName,
		Permissions: &[]int64{0}[0],
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create jail role: %s", err)
	}
	return role, nil
}

// staffstripMessage: -staffstrip <user>
func staffstripMessageCommandHandler(h *CommandHandler, s *discordgo.Session, m *discordgo.MessageCreate, args []string) error {
	if ok, err := requireModForMessage(s, m, "Staff Strip"); !ok {
		return err
	}
	if len(args) == 0 {
		return sendModError(s, m.ChannelID, "Staff Strip", "Usage: `-staffstrip <user>`")
	}
	targetID, name, err := resolveTargetUser(s, m.GuildID, args[0])
	if err != nil {
		return sendModError(s, m.ChannelID, "Staff Strip", fmt.Sprintf("Could not find that user: %s", err))
	}
	removed, err := stripStaffRoles(s, m.GuildID, targetID)
	if err != nil {
		return sendModError(s, m.ChannelID, "Staff Strip", err.Error())
	}
	auditAction(h, context.Background(), m.GuildID, "STAFFSTRIP", m.Author.ID, targetID, map[string]any{"removed": removed})
	_, err = s.ChannelMessageSendEmbed(m.ChannelID, modSuccessEmbed("Staff Strip", fmt.Sprintf("Stripped **%d** staff role%s from **%s**.", removed, plural(removed), name)))
	return err
}

func staffstripSlashCommandHandler(h *CommandHandler, s *discordgo.Session, i *discordgo.InteractionCreate) error {
	ok, msg := requireModerator(i.Member)
	if !ok {
		return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{Embeds: []*discordgo.MessageEmbed{modErrorEmbed("Staff Strip", msg)}},
		})
	}
	opts := OptionMap(i.ApplicationCommandData().Options)
	raw := OptString(opts, "user")
	targetID, name, err := resolveTargetUser(s, i.GuildID, raw)
	if err != nil {
		return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{Embeds: []*discordgo.MessageEmbed{modErrorEmbed("Staff Strip", fmt.Sprintf("Could not find that user: %s", err))}},
		})
	}
	removed, err := stripStaffRoles(s, i.GuildID, targetID)
	if err != nil {
		return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{Embeds: []*discordgo.MessageEmbed{modErrorEmbed("Staff Strip", err.Error())}},
		})
	}
	auditAction(h, context.Background(), i.GuildID, "STAFFSTRIP", i.Member.User.ID, targetID, map[string]any{"removed": removed})
	return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{Embeds: []*discordgo.MessageEmbed{modSuccessEmbed("Staff Strip", fmt.Sprintf("Stripped **%d** staff role%s from **%s**.", removed, plural(removed), name))}},
	})
}

// stripStaffRoles removes roles that grant moderation/management permissions.
func stripStaffRoles(s *discordgo.Session, guildID, targetID string) (int, error) {
	member, err := s.GuildMember(guildID, targetID)
	if err != nil {
		return 0, fmt.Errorf("could not fetch member: %s", err)
	}
	roles, err := s.GuildRoles(guildID)
	if err != nil {
		return 0, fmt.Errorf("could not fetch roles: %s", err)
	}
	removed := 0
	for _, rid := range member.Roles {
		if rid == guildID {
			continue
		}
		for _, role := range roles {
			if role.ID != rid {
				continue
			}
			if role.Permissions&defaultModPerm != 0 {
				if err := s.GuildMemberRoleRemove(guildID, targetID, rid); err == nil {
					removed++
				}
			}
			break
		}
	}
	return removed, nil
}
