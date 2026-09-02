package commands

import (
	"context"
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
)

// denypermMessage: -denyperm <user|role> <permission> [channel]
func denypermMessageCommandHandler(h *CommandHandler, s *discordgo.Session, m *discordgo.MessageCreate, args []string) error {
	if ok, err := requireModForMessage(s, m, "Deny Permission"); !ok {
		return err
	}
	if len(args) < 2 {
		return sendModError(s, m.ChannelID, "Deny Permission", "Usage: `-denyperm <user|@role> <permission> [channel]`")
	}
	targetID, targetType, err := resolveChannelTarget(s, m.GuildID, args[0])
	if err != nil {
		return sendModError(s, m.ChannelID, "Deny Permission", err.Error())
	}
	perm, ok := parsePermission(args[1])
	if !ok {
		return sendModError(s, m.ChannelID, "Deny Permission", "Unknown permission. Try `view`, `send`, `attach`, `embed`, `reaction`, or `voice`.")
	}
	chID := m.ChannelID
	if len(args) >= 3 {
		ch, cerr := s.Channel(args[2])
		if cerr == nil {
			chID = ch.ID
		}
	}
	if err := setPermissionOverride(s, chID, targetID, targetType, perm, true); err != nil {
		return sendModError(s, m.ChannelID, "Deny Permission", fmt.Sprintf("Failed to deny permission: %s", err))
	}
	auditAction(h, context.Background(), m.GuildID, "DENY_PERM", m.Author.ID, targetID, map[string]any{"permission": args[1], "channel": chID})
	_, err = s.ChannelMessageSendEmbed(m.ChannelID, modSuccessEmbed("Deny Permission", fmt.Sprintf("Denied `%s` for <@%s> in <#%s>.", args[1], targetID, chID)))
	return err
}

func denypermSlashCommandHandler(h *CommandHandler, s *discordgo.Session, i *discordgo.InteractionCreate) error {
	ok, msg := requireModerator(s, i.GuildID, i.Member.User.ID, i.Member)
	if !ok {
		return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{Embeds: []*discordgo.MessageEmbed{modErrorEmbed("Deny Permission", msg)}},
		})
	}
	opts := OptionMap(i.ApplicationCommandData().Options)
	targetID := OptUser(opts, "user")
	targetType := discordgo.PermissionOverwriteTypeMember
	if targetID == "" {
		targetID = OptString(opts, "role")
		targetType = discordgo.PermissionOverwriteTypeRole
	}
	permStr := OptString(opts, "permission")
	perm, ok := parsePermission(permStr)
	if !ok {
		return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{Embeds: []*discordgo.MessageEmbed{modErrorEmbed("Deny Permission", "Unknown permission. Try `view`, `send`, `attach`, `embed`, `reaction`, or `voice`.")}},
		})
	}
	chID := OptString(opts, "channel")
	if chID == "" {
		chID = i.ChannelID
	}
	if err := setPermissionOverride(s, chID, targetID, targetType, perm, true); err != nil {
		return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{Embeds: []*discordgo.MessageEmbed{modErrorEmbed("Deny Permission", fmt.Sprintf("Failed to deny permission: %s", err))}},
		})
	}
	auditAction(h, context.Background(), i.GuildID, "DENY_PERM", i.Member.User.ID, targetID, map[string]any{"permission": permStr, "channel": chID})
	return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{Embeds: []*discordgo.MessageEmbed{modSuccessEmbed("Deny Permission", fmt.Sprintf("Denied `%s` for <@%s> in <#%s>.", permStr, targetID, chID))}},
	})
}

// imuteMessage: -imute <user>  (deny attachment uploads)
func imuteMessageCommandHandler(h *CommandHandler, s *discordgo.Session, m *discordgo.MessageCreate, args []string) error {
	return mediaMute(h, s, m.GuildID, m.ChannelID, m.Author.ID, args, "Image Mute")
}
func imuteSlashCommandHandler(h *CommandHandler, s *discordgo.Session, i *discordgo.InteractionCreate) error {
	return mediaMuteSlash(h, s, i, "Image Mute")
}

// gifmuteMessage: -gifmute <user>  (deny embed links + attachments)
func gifmuteMessageCommandHandler(h *CommandHandler, s *discordgo.Session, m *discordgo.MessageCreate, args []string) error {
	return mediaMute(h, s, m.GuildID, m.ChannelID, m.Author.ID, args, "GIF Mute")
}
func gifmuteSlashCommandHandler(h *CommandHandler, s *discordgo.Session, i *discordgo.InteractionCreate) error {
	return mediaMuteSlash(h, s, i, "GIF Mute")
}

func mediaMute(h *CommandHandler, s *discordgo.Session, guildID, channelID, actorID string, args []string, title string) error {
	if ok, err := requireModForUser(s, guildID, actorID, channelID, title); !ok {
		return err
	}
	if len(args) == 0 {
		return sendModError(s, channelID, title, fmt.Sprintf("Usage: `-%s <user>`", strings.ToLower(strings.ReplaceAll(title, " ", ""))))
	}
	targetID, name, err := resolveTargetUser(s, guildID, args[0])
	if err != nil {
		return sendModError(s, channelID, title, fmt.Sprintf("Could not find that user: %s", err))
	}
	perm := discordgo.PermissionAttachFiles
	if title == "GIF Mute" {
		perm = discordgo.PermissionEmbedLinks | discordgo.PermissionAttachFiles
	}
	if err := setPermissionOverride(s, channelID, targetID, discordgo.PermissionOverwriteTypeMember, int64(perm), true); err != nil {
		return sendModError(s, channelID, title, fmt.Sprintf("Failed to update permissions: %s", err))
	}
	auditAction(h, context.Background(), guildID, "MEDIA_MUTE", actorID, targetID, map[string]any{"type": title, "channel": channelID})
	_, err = s.ChannelMessageSendEmbed(channelID, modSuccessEmbed(title, fmt.Sprintf("Muted **%s** in this channel.", name)))
	return err
}

func mediaMuteSlash(h *CommandHandler, s *discordgo.Session, i *discordgo.InteractionCreate, title string) error {
	ok, msg := requireModerator(s, i.GuildID, i.Member.User.ID, i.Member)
	if !ok {
		return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{Embeds: []*discordgo.MessageEmbed{modErrorEmbed(title, msg)}},
		})
	}
	opts := OptionMap(i.ApplicationCommandData().Options)
	raw := OptString(opts, "user")
	targetID, name, err := resolveTargetUser(s, i.GuildID, raw)
	if err != nil {
		return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{Embeds: []*discordgo.MessageEmbed{modErrorEmbed(title, fmt.Sprintf("Could not find that user: %s", err))}},
		})
	}
	perm := discordgo.PermissionAttachFiles
	if title == "GIF Mute" {
		perm = discordgo.PermissionEmbedLinks | discordgo.PermissionAttachFiles
	}
	if err := setPermissionOverride(s, i.ChannelID, targetID, discordgo.PermissionOverwriteTypeMember, int64(perm), true); err != nil {
		return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{Embeds: []*discordgo.MessageEmbed{modErrorEmbed(title, fmt.Sprintf("Failed to update permissions: %s", err))}},
		})
	}
	auditAction(h, context.Background(), i.GuildID, "MEDIA_MUTE", i.Member.User.ID, targetID, map[string]any{"type": title, "channel": i.ChannelID})
	return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{Embeds: []*discordgo.MessageEmbed{modSuccessEmbed(title, fmt.Sprintf("Muted **%s** in this channel.", name))}},
	})
}

func resolveChannelTarget(s *discordgo.Session, guildID, raw string) (string, discordgo.PermissionOverwriteType, error) {
	id := parseMentionID(raw)
	// Try role first if it looks like a role mention, else member.
	roles, _ := s.GuildRoles(guildID)
	for _, r := range roles {
		if r.ID == id {
			return id, discordgo.PermissionOverwriteTypeRole, nil
		}
	}
	if _, err := s.GuildMember(guildID, id); err == nil {
		return id, discordgo.PermissionOverwriteTypeMember, nil
	}
	return "", 0, fmt.Errorf("could not find that user or role")
}

func parsePermission(s string) (int64, bool) {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "view":
		return discordgo.PermissionViewChannel, true
	case "send":
		return discordgo.PermissionSendMessages, true
	case "attach", "file":
		return discordgo.PermissionAttachFiles, true
	case "embed", "embedlinks":
		return discordgo.PermissionEmbedLinks, true
	case "reaction", "addreactions":
		return discordgo.PermissionAddReactions, true
	case "voice", "speak":
		return discordgo.PermissionVoiceSpeak, true
	}
	return 0, false
}

func setPermissionOverride(s *discordgo.Session, channelID, targetID string, targetType discordgo.PermissionOverwriteType, deny int64, denyMode bool) error {
	current, err := s.Channel(channelID)
	if err != nil {
		return err
	}
	var allow, denyVal int64
	for _, o := range current.PermissionOverwrites {
		if o.ID == targetID {
			allow = o.Allow
			denyVal = o.Deny
			break
		}
	}
	if denyMode {
		denyVal |= deny
		allow &^= deny
	} else {
		allow |= deny
		denyVal &^= deny
	}
	err = s.ChannelPermissionSet(channelID, targetID, targetType, allow, denyVal)
	return err
}
