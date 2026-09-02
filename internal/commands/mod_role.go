package commands

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
)

const (
	modRolePerm = discordgo.PermissionManageRoles
)

// roleMessage: -role <add|create|edit|hoist|member> ...
func roleMessageCommandHandler(h *CommandHandler, s *discordgo.Session, m *discordgo.MessageCreate, args []string) error {
	if ok, err := requireModForMessage(s, m, "Role"); !ok {
		return err
	}
	if len(args) == 0 {
		return sendModError(s, m.ChannelID, "Role", "Usage: `-role add <user> <role>`, `-role create <name>`, `-role edit <role> ...`, `-role hoist <role> <on|off>`, `-role member <role>`")
	}
	switch strings.ToLower(args[0]) {
	case "add":
		return roleAddMessage(h, s, m, args[1:])
	case "create":
		return roleCreateMessage(h, s, m, args[1:])
	case "edit":
		return roleEditMessage(h, s, m, args[1:])
	case "hoist":
		return roleHoistMessage(h, s, m, args[1:])
	case "member", "members":
		return roleMemberMessage(s, m, args[1:])
	default:
		return sendModError(s, m.ChannelID, "Role", "Unknown subcommand. Use `add`, `create`, `edit`, `hoist`, or `member`.")
	}
}

func roleSlashCommandHandler(h *CommandHandler, s *discordgo.Session, i *discordgo.InteractionCreate) error {
	ok, msg := requireModerator(s, i.GuildID, i.Member.User.ID, i.Member)
	if !ok {
		return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{Embeds: []*discordgo.MessageEmbed{modErrorEmbed("Role", msg)}},
		})
	}
	subs := i.ApplicationCommandData().Options
	if len(subs) == 0 {
		return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{Embeds: []*discordgo.MessageEmbed{modErrorEmbed("Role", "No subcommand provided.")}},
		})
	}
	sub := subs[0]
	opts := OptionMap(sub.Options)
	switch sub.Name {
	case "add":
		return roleAddSlash(h, s, i, sub, opts)
	case "create":
		return roleCreateSlash(h, s, i, opts)
	case "edit":
		return roleEditSlash(h, s, i, opts)
	case "hoist":
		return roleHoistSlash(h, s, i, opts)
	case "member":
		return roleMemberSlash(s, i, opts)
	default:
		return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{Embeds: []*discordgo.MessageEmbed{modErrorEmbed("Role", "Unknown subcommand.")}},
		})
	}
}

// role add <user> <role>
func roleAddMessage(h *CommandHandler, s *discordgo.Session, m *discordgo.MessageCreate, args []string) error {
	if len(args) < 2 {
		return sendModError(s, m.ChannelID, "Role add", "Usage: `-role add <user> <role>`")
	}
	memberID, _, err := resolveTargetUser(s, m.GuildID, args[0])
	if err != nil {
		return sendModError(s, m.ChannelID, "Role add", fmt.Sprintf("Could not find that user: %s", err))
	}
	role, err := resolveTargetRole(s, m.GuildID, strings.Join(args[1:], " "))
	if err != nil {
		return sendModError(s, m.ChannelID, "Role add", err.Error())
	}
	if err := s.GuildMemberRoleAdd(m.GuildID, memberID, role.ID); err != nil {
		return sendModError(s, m.ChannelID, "Role add", fmt.Sprintf("Failed to add role: %s", err))
	}
	auditAction(h, context.Background(), m.GuildID, "ROLE_ADD", m.Author.ID, memberID, map[string]any{"role": role.ID, "role_name": role.Name})
	_, err = s.ChannelMessageSendEmbed(m.ChannelID, modSuccessEmbed("Role add", fmt.Sprintf("Added role **%s**.", role.Name)))
	return err
}

func roleAddSlash(h *CommandHandler, s *discordgo.Session, i *discordgo.InteractionCreate, sub *discordgo.ApplicationCommandInteractionDataOption, opts map[string]*discordgo.ApplicationCommandInteractionDataOption) error {
	memberID := OptUser(opts, "user")
	roleID := OptString(opts, "role")
	role, err := resolveTargetRole(s, i.GuildID, roleID)
	if err != nil {
		return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{Embeds: []*discordgo.MessageEmbed{modErrorEmbed("Role add", err.Error())}},
		})
	}
	if err := s.GuildMemberRoleAdd(i.GuildID, memberID, role.ID); err != nil {
		return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{Embeds: []*discordgo.MessageEmbed{modErrorEmbed("Role add", fmt.Sprintf("Failed to add role: %s", err))}},
		})
	}
	auditAction(h, context.Background(), i.GuildID, "ROLE_ADD", i.Member.User.ID, memberID, map[string]any{"role": role.ID, "role_name": role.Name})
	return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{Embeds: []*discordgo.MessageEmbed{modSuccessEmbed("Role add", fmt.Sprintf("Added role **%s**.", role.Name))}},
	})
}

// role create <name>
func roleCreateMessage(h *CommandHandler, s *discordgo.Session, m *discordgo.MessageCreate, args []string) error {
	if len(args) == 0 {
		return sendModError(s, m.ChannelID, "Role create", "Usage: `-role create <name>`")
	}
	role, err := s.GuildRoleCreate(m.GuildID, &discordgo.RoleParams{Name: strings.Join(args, " ")})
	if err != nil {
		return sendModError(s, m.ChannelID, "Role create", fmt.Sprintf("Failed to create role: %s", err))
	}
	auditAction(h, context.Background(), m.GuildID, "ROLE_CREATE", m.Author.ID, role.ID, map[string]any{"role_name": role.Name})
	_, err = s.ChannelMessageSendEmbed(m.ChannelID, modSuccessEmbed("Role create", fmt.Sprintf("Created role **%s**.", role.Name)))
	return err
}

func roleCreateSlash(h *CommandHandler, s *discordgo.Session, i *discordgo.InteractionCreate, opts map[string]*discordgo.ApplicationCommandInteractionDataOption) error {
	name := OptString(opts, "name")
	if name == "" {
		return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{Embeds: []*discordgo.MessageEmbed{modErrorEmbed("Role create", "Name is required.")}},
		})
	}
	role, err := s.GuildRoleCreate(i.GuildID, &discordgo.RoleParams{Name: name})
	if err != nil {
		return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{Embeds: []*discordgo.MessageEmbed{modErrorEmbed("Role create", fmt.Sprintf("Failed to create role: %s", err))}},
		})
	}
	auditAction(h, context.Background(), i.GuildID, "ROLE_CREATE", i.Member.User.ID, role.ID, map[string]any{"role_name": role.Name})
	return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{Embeds: []*discordgo.MessageEmbed{modSuccessEmbed("Role create", fmt.Sprintf("Created role **%s**.", role.Name))}},
	})
}

// role edit <role> [name:<name>] [color:<hex>]
func roleEditMessage(h *CommandHandler, s *discordgo.Session, m *discordgo.MessageCreate, args []string) error {
	if len(args) < 2 {
		return sendModError(s, m.ChannelID, "Role edit", "Usage: `-role edit <role> [name:<name>] [color:<hex>]`")
	}
	role, err := resolveTargetRole(s, m.GuildID, args[0])
	if err != nil {
		return sendModError(s, m.ChannelID, "Role edit", err.Error())
	}
	params := roleEditParams(args[1:], role)
	updated, err := s.GuildRoleEdit(m.GuildID, role.ID, params)
	if err != nil {
		return sendModError(s, m.ChannelID, "Role edit", fmt.Sprintf("Failed to edit role: %s", err))
	}
	auditAction(h, context.Background(), m.GuildID, "ROLE_EDIT", m.Author.ID, role.ID, map[string]any{"role_name": updated.Name})
	_, err = s.ChannelMessageSendEmbed(m.ChannelID, modSuccessEmbed("Role edit", fmt.Sprintf("Edited role **%s**.", updated.Name)))
	return err
}

func roleEditSlash(h *CommandHandler, s *discordgo.Session, i *discordgo.InteractionCreate, opts map[string]*discordgo.ApplicationCommandInteractionDataOption) error {
	roleID := OptString(opts, "role")
	role, err := resolveTargetRole(s, i.GuildID, roleID)
	if err != nil {
		return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{Embeds: []*discordgo.MessageEmbed{modErrorEmbed("Role edit", err.Error())}},
		})
	}
	var params discordgo.RoleParams
	if n := OptString(opts, "name"); n != "" {
		params.Name = n
	}
	if c := OptString(opts, "color"); c != "" {
		col := parseHexColor(c)
		params.Color = &col
	}
	updated, err := s.GuildRoleEdit(i.GuildID, role.ID, &params)
	if err != nil {
		return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{Embeds: []*discordgo.MessageEmbed{modErrorEmbed("Role edit", fmt.Sprintf("Failed to edit role: %s", err))}},
		})
	}
	auditAction(h, context.Background(), i.GuildID, "ROLE_EDIT", i.Member.User.ID, role.ID, map[string]any{"role_name": updated.Name})
	return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{Embeds: []*discordgo.MessageEmbed{modSuccessEmbed("Role edit", fmt.Sprintf("Edited role **%s**.", updated.Name))}},
	})
}

func roleEditParams(args []string, role *discordgo.Role) *discordgo.RoleParams {
	params := &discordgo.RoleParams{}
	set := false
	for _, a := range args {
		if strings.HasPrefix(a, "name:") {
			params.Name = strings.TrimPrefix(a, "name:")
			set = true
		}
		if strings.HasPrefix(a, "color:") {
			c := parseHexColor(strings.TrimPrefix(a, "color:"))
			params.Color = &c
			set = true
		}
	}
	if !set {
		// No changes given; keep existing values so the edit is a no-op safe call.
		params.Name = role.Name
		c := role.Color
		params.Color = &c
	}
	return params
}

func parseHexColor(s string) int {
	s = strings.TrimPrefix(strings.TrimSpace(s), "#")
	if s == "" {
		return 0
	}
	v, err := strconv.ParseUint(s, 16, 32)
	if err != nil {
		return 0
	}
	return int(v)
}

// role hoist <role> <on|off>
func roleHoistMessage(h *CommandHandler, s *discordgo.Session, m *discordgo.MessageCreate, args []string) error {
	if len(args) < 2 {
		return sendModError(s, m.ChannelID, "Role hoist", "Usage: `-role hoist <role> <on|off>`")
	}
	role, err := resolveTargetRole(s, m.GuildID, args[0])
	if err != nil {
		return sendModError(s, m.ChannelID, "Role hoist", err.Error())
	}
	hoist := strings.EqualFold(args[1], "on")
	if _, err := s.GuildRoleEdit(m.GuildID, role.ID, &discordgo.RoleParams{Hoist: &hoist}); err != nil {
		return sendModError(s, m.ChannelID, "Role hoist", fmt.Sprintf("Failed to update role: %s", err))
	}
	auditAction(h, context.Background(), m.GuildID, "ROLE_HOIST", m.Author.ID, role.ID, map[string]any{"hoist": hoist})
	state := "enabled"
	if !hoist {
		state = "disabled"
	}
	_, err = s.ChannelMessageSendEmbed(m.ChannelID, modSuccessEmbed("Role hoist", fmt.Sprintf("Hoist %s for **%s**.", state, role.Name)))
	return err
}

func roleHoistSlash(h *CommandHandler, s *discordgo.Session, i *discordgo.InteractionCreate, opts map[string]*discordgo.ApplicationCommandInteractionDataOption) error {
	roleID := OptString(opts, "role")
	role, err := resolveTargetRole(s, i.GuildID, roleID)
	if err != nil {
		return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{Embeds: []*discordgo.MessageEmbed{modErrorEmbed("Role hoist", err.Error())}},
		})
	}
	hoist := OptBool(opts, "hoist")
	if _, err := s.GuildRoleEdit(i.GuildID, role.ID, &discordgo.RoleParams{Hoist: &hoist}); err != nil {
		return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{Embeds: []*discordgo.MessageEmbed{modErrorEmbed("Role hoist", fmt.Sprintf("Failed to update role: %s", err))}},
		})
	}
	auditAction(h, context.Background(), i.GuildID, "ROLE_HOIST", i.Member.User.ID, role.ID, map[string]any{"hoist": hoist})
	state := "disabled"
	if hoist {
		state = "enabled"
	}
	return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{Embeds: []*discordgo.MessageEmbed{modSuccessEmbed("Role hoist", fmt.Sprintf("Hoist %s for **%s**.", state, role.Name))}},
	})
}

// role member <role>
func roleMemberMessage(s *discordgo.Session, m *discordgo.MessageCreate, args []string) error {
	if len(args) == 0 {
		return sendModError(s, m.ChannelID, "Role member", "Usage: `-role member <role>`")
	}
	role, err := resolveTargetRole(s, m.GuildID, strings.Join(args, " "))
	if err != nil {
		return sendModError(s, m.ChannelID, "Role member", err.Error())
	}
	embed, err := buildRoleMembersEmbed(s, m.GuildID, role)
	if err != nil {
		return sendModError(s, m.ChannelID, "Role member", err.Error())
	}
	_, err = s.ChannelMessageSendEmbed(m.ChannelID, embed)
	return err
}

func roleMemberSlash(s *discordgo.Session, i *discordgo.InteractionCreate, opts map[string]*discordgo.ApplicationCommandInteractionDataOption) error {
	roleID := OptString(opts, "role")
	role, err := resolveTargetRole(s, i.GuildID, roleID)
	if err != nil {
		return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{Embeds: []*discordgo.MessageEmbed{modErrorEmbed("Role member", err.Error())}},
		})
	}
	embed, err := buildRoleMembersEmbed(s, i.GuildID, role)
	if err != nil {
		return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{Embeds: []*discordgo.MessageEmbed{modErrorEmbed("Role member", err.Error())}},
		})
	}
	return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{Embeds: []*discordgo.MessageEmbed{embed}},
	})
}

func buildRoleMembersEmbed(s *discordgo.Session, guildID string, role *discordgo.Role) (*discordgo.MessageEmbed, error) {
	members, err := s.GuildMembers(guildID, "", 1000)
	if err != nil {
		return nil, fmt.Errorf("could not fetch members: %s", err)
	}
	var names []string
	for _, m := range members {
		for _, rid := range m.Roles {
			if rid == role.ID {
				n := m.User.Username
				if m.Nick != "" {
					n = m.Nick
				}
				names = append(names, n)
				break
			}
		}
	}
	if len(names) == 0 {
		names = append(names, "No members.")
	}
	if len(names) > 30 {
		names = names[:30]
		names = append(names, fmt.Sprintf("... and %d more", len(members)-30))
	}
	return &discordgo.MessageEmbed{
		Color:       modColor,
		Title:       fmt.Sprintf("Role members — %s", role.Name),
		Description: strings.Join(names, "\n"),
	}, nil
}
