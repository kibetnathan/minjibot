package commands

import (
	"context"
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/kibetnathan/minjibot/internal/config"
	"github.com/kibetnathan/minjibot/internal/ports/repository"
)

type CommandHandler struct {
	Cfg          *config.Config
	GuildRepo    repository.GuildRepository
	SettingsRepo repository.GuildSettingsRepository
	PermRepo     repository.UserPermissionRepository
	AuditRepo    repository.AuditLogRepository
	BirthdayRepo repository.BirthdayRepository
	BirthdaySett repository.GuildBirthdaySettingsRepository
	DiaryRepo    repository.DiaryRepository
}

func NewCommandHandler(cfg *config.Config, guildRepo repository.GuildRepository, settingsRepo repository.GuildSettingsRepository, permRepo repository.UserPermissionRepository, auditRepo repository.AuditLogRepository, birthdayRepo repository.BirthdayRepository, birthdaySett repository.GuildBirthdaySettingsRepository, diaryRepo repository.DiaryRepository) *CommandHandler {
	return &CommandHandler{
		Cfg:          cfg,
		GuildRepo:    guildRepo,
		SettingsRepo: settingsRepo,
		PermRepo:     permRepo,
		AuditRepo:    auditRepo,
		BirthdayRepo: birthdayRepo,
		BirthdaySett: birthdaySett,
		DiaryRepo:    diaryRepo,
	}
}

func (h *CommandHandler) Handle(ctx context.Context, s *discordgo.Session, m *discordgo.MessageCreate, cmd string, args []string) error {
	switch cmd {
	case "ping":
		return h.ping(s, m.ChannelID)
	case "help":
		return h.help(s, m, args)
	case "tldr":
		return h.tldr(s, m, args)
	case "echo":
		return h.echo(s, m, args)
	case "userinfo":
		return h.userInfo(s, m, args)
	case "avatar":
		return h.avatar(s, m, args)
	case "banner":
		return h.banner(s, m, args)
	case "botinfo":
		return h.botinfo(s, m, args)
	case "channelinfo":
		return h.channelinfo(s, m, args)
	case "roles":
		return h.roles(s, m, args)
	case "guild":
		return h.guild(s, m, args)
	case "emojis":
		return h.emojis(s, m, args)
	case "stickers":
		return h.stickers(s, m, args)
	case "bans":
		return h.bans(s, m, args)
	case "boomer":
		return h.boomer(s, m, args)
	case "perms":
		return h.perms(s, m, args)
	case "tz":
		return h.tz(s, m, args)
	case "urbandictionary":
		return h.urbandictionary(s, m, args)
	case "ddg":
		return h.ddg(s, m.ChannelID, args)
	case "search":
		return h.search(s, m, args)
	case "pinglist":
		return h.pinglist(s, m, args)
	case "gifsearch":
		return h.gifsearch(s, m.ChannelID, args)
	case "emoji":
		return h.emoji(s, m, args)
	case "sticker":
		return h.sticker(s, m, args)
	case "pin":
		return h.pin(s, m, args)
	case "unpin":
		return h.unpin(s, m, args)
	case "quote":
		return h.quote(s, m, args)
	case "translate":
		return h.translate(s, m, args)
	case "reminder":
		return h.reminder(s, m, args)
	case "isearch":
		return h.isearch(s, m, args)
	case "caption":
		return h.caption(s, m, args)
	case "img2gif":
		return h.img2gif(s, m, args)
	case "vid2gif":
		return h.vid2gif(s, m, args)
	case "autogif":
		return h.autogif(s, m, args)
	case "factcheck":
		return h.factcheck(s, m, args)
	case "howgay":
		return h.howgay(s, m, args)
	case "howautism":
		return h.howautism(s, m, args)
	case "howlesbian":
		return h.howlesbian(s, m, args)
	case "howsimp":
		return h.howsimp(s, m, args)
	case "pp":
		return h.pp(s, m, args)
	case "puh":
		return h.puh(s, m, args)
	case "iq":
		return h.iq(s, m, args)
	case "bitches":
		return h.bitches(s, m, args)
	case "choose":
		return h.choose(s, m, args)
	case "ship":
		return h.ship(s, m, args)
	case "colors":
		return h.colorsAvatar(s, m, args)
	case "lurk":
		return h.lurk(s, m, args)
	case "lurkers":
		return h.lurkers(s, m, args)
	case "spark":
		return h.spark(s, m, args)
	case "smoke":
		return h.smoke(s, m, args)
	case "hits":
		return h.hits(s, m, args)
	case "compress":
		return h.compress(s, m, args)
	case "vape":
		return h.vape(s, m, args)
	case "poll":
		return h.poll(s, m, args)
	case "quickpoll":
		return h.quickpoll(s, m, args)
	case "birthday":
		return h.birthday(s, m, args)
	case "diary":
		return h.diary(s, m, args)
	case "ttys":
		return h.ttys(s, m, args)
	case "bio":
		return h.bio(s, m, args)
	case "ban":
		return h.ban(s, m, args)
	case "hardban":
		return h.hardban(s, m, args)
	case "softban":
		return h.softban(s, m, args)
	case "kick":
		return h.kick(s, m, args)
	case "purge":
		return h.purge(s, m, args)
	case "nuke":
		return h.nuke(s, m, args)
	case "timeout":
		return h.timeout(s, m, args)
	case "warn":
		return h.warn(s, m, args)
	case "history":
		return h.history(s, m, args)
	default:
		return fmt.Errorf("unknown command: %s", cmd)
	}
}

func (h *CommandHandler) HandleSlash(ctx context.Context, s *discordgo.Session, i *discordgo.InteractionCreate) error {
	switch i.ApplicationCommandData().Name {
	case "ping":
		return h.pingSlash(s, i)
	case "help":
		return h.helpSlash(s, i)
	case "tldr":
		return h.tldrSlash(s, i)
	case "echo":
		return h.echoSlash(s, i)
	case "userinfo":
		return h.userInfoSlash(s, i)
	case "avatar":
		return h.avatarSlash(s, i)
	case "banner":
		return h.bannerSlash(s, i)
	case "botinfo":
		return h.botinfoSlash(s, i)
	case "channelinfo":
		return h.channelinfoSlash(s, i)
	case "roles":
		return h.rolesSlash(s, i)
	case "guild":
		return h.guildSlash(s, i)
	case "emojis":
		return h.emojisSlash(s, i)
	case "stickers":
		return h.stickersSlash(s, i)
	case "bans":
		return h.bansSlash(s, i)
	case "boomer":
		return h.boomerSlash(s, i)
	case "perms":
		return h.permsSlash(s, i)
	case "tz":
		return h.tzSlash(s, i)
	case "urbandictionary":
		return h.urbandictionarySlash(s, i)
	case "ddg":
		return h.ddgSlash(s, i)
	case "search":
		return h.searchSlash(s, i)
	case "pinglist":
		return h.pinglistSlash(s, i)
	case "gifsearch":
		return h.gifsearchSlash(s, i)
	case "emoji":
		return h.emojiSlash(s, i)
	case "sticker":
		return h.stickerSlash(s, i)
	case "pin":
		return h.pinSlash(s, i)
	case "unpin":
		return h.unpinSlash(s, i)
	case "quote":
		return h.quoteSlash(s, i)
	case "translate":
		return h.translateSlash(s, i)
	case "reminder":
		return h.reminderSlash(s, i)
	case "isearch":
		return h.isearchSlash(s, i)
	case "caption":
		return h.captionSlash(s, i)
	case "img2gif":
		return h.img2gifSlash(s, i)
	case "vid2gif":
		return h.vid2gifSlash(s, i)
	case "autogif":
		return h.autogifSlash(s, i)
	case "factcheck":
		return h.factcheckSlash(s, i)
	case "howgay":
		return h.howgaySlash(s, i)
	case "howautism":
		return h.howautismSlash(s, i)
	case "howlesbian":
		return h.howlesbianSlash(s, i)
	case "howsimp":
		return h.howsimpSlash(s, i)
	case "pp":
		return h.ppSlash(s, i)
	case "puh":
		return h.puhSlash(s, i)
	case "iq":
		return h.iqSlash(s, i)
	case "bitches":
		return h.bitchesSlash(s, i)
	case "choose":
		return h.chooseSlash(s, i)
	case "ship":
		return h.shipSlash(s, i)
	case "colors":
		return h.colorsAvatarSlash(s, i)
	case "lurk":
		return h.lurkSlash(s, i)
	case "lurkers":
		return h.lurkersSlash(s, i)
	case "spark":
		return h.sparkSlash(s, i)
	case "smoke":
		return h.smokeSlash(s, i)
	case "hits":
		return h.hitsSlash(s, i)
	case "compress":
		return h.compressSlash(s, i)
	case "vape":
		return h.vapeSlash(s, i)
	case "poll":
		return h.pollSlash(s, i)
	case "quickpoll":
		return h.quickpollSlash(s, i)
	case "birthday":
		return h.birthdaySlash(s, i)
	case "diary":
		return h.diarySlash(s, i)
	case "ttys":
		return h.ttysSlash(s, i)
	case "bio":
		return h.bioSlash(s, i)
	case "ban":
		return h.banSlash(s, i)
	case "hardban":
		return h.hardbanSlash(s, i)
	case "softban":
		return h.softbanSlash(s, i)
	case "kick":
		return h.kickSlash(s, i)
	case "purge":
		return h.purgeSlash(s, i)
	case "nuke":
		return h.nukeSlash(s, i)
	case "timeout":
		return h.timeoutSlash(s, i)
	case "warn":
		return h.warnSlash(s, i)
	case "history":
		return h.historySlash(s, i)
	default:
		return fmt.Errorf("unknown command: %s", i.ApplicationCommandData().Name)
	}
}

func (h *CommandHandler) ping(s *discordgo.Session, channelID string) error {
	return pingMessageCommandHandler(s, channelID)
}

func (h *CommandHandler) pingSlash(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	return pingSlashCommandHandler(s, i)
}

func (h *CommandHandler) help(s *discordgo.Session, m *discordgo.MessageCreate, args []string) error {
	// -help <category> shows a single page; plain -help paginates.
	if len(args) > 0 {
		if idx := FindHelpSection(args[0]); idx >= 0 {
			_, err := s.ChannelMessageSendEmbed(m.ChannelID, BuildHelpPageEmbed(idx))
			return err
		}
	}
	return paginateReactions(s, m.ChannelID, m.Author.ID, NumHelpPages(), BuildHelpPageEmbed)
}

func (h *CommandHandler) helpSlash(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	components := []discordgo.MessageComponent{
		discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				discordgo.Button{Label: "◀", Style: discordgo.SecondaryButton, CustomID: helpPrevCustomID},
				discordgo.Button{Label: "▶", Style: discordgo.SecondaryButton, CustomID: helpNextCustomID},
			},
		},
	}
	return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds:     []*discordgo.MessageEmbed{BuildHelpPageEmbed(0)},
			Components: components,
		},
	})
}

func (h *CommandHandler) tldr(s *discordgo.Session, m *discordgo.MessageCreate, args []string) error {
	return tldrMessageCommandHandler(s, m, args)
}

func (h *CommandHandler) tldrSlash(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	return tldrSlashCommandHandler(s, i)
}

func (h *CommandHandler) echo(s *discordgo.Session, m *discordgo.MessageCreate, args []string) error {
	return echoMessageCommandHandler(s, m, args)
}

func (h *CommandHandler) echoSlash(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	return echoSlashCommandHandler(s, i)
}

func (h *CommandHandler) userInfo(s *discordgo.Session, m *discordgo.MessageCreate, args []string) error {
	return userInfoMessageCommandHandler(s, m, args)
}

func (h *CommandHandler) userInfoSlash(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	return userInfoSlashCommandHandler(s, i)
}

func (h *CommandHandler) avatar(s *discordgo.Session, m *discordgo.MessageCreate, args []string) error {
	return avatarMessageCommandHandler(s, m, args)
}
func (h *CommandHandler) avatarSlash(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	return avatarSlashCommandHandler(s, i)
}

func (h *CommandHandler) banner(s *discordgo.Session, m *discordgo.MessageCreate, args []string) error {
	return bannerMessageCommandHandler(s, m, args)
}
func (h *CommandHandler) bannerSlash(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	return bannerSlashCommandHandler(s, i)
}

func (h *CommandHandler) botinfo(s *discordgo.Session, m *discordgo.MessageCreate, args []string) error {
	return botinfoMessageCommandHandler(s, m, args)
}
func (h *CommandHandler) botinfoSlash(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	return botinfoSlashCommandHandler(s, i)
}

func (h *CommandHandler) channelinfo(s *discordgo.Session, m *discordgo.MessageCreate, args []string) error {
	return channelinfoMessageCommandHandler(s, m, args)
}
func (h *CommandHandler) channelinfoSlash(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	return channelinfoSlashCommandHandler(s, i)
}

func (h *CommandHandler) roles(s *discordgo.Session, m *discordgo.MessageCreate, args []string) error {
	return rolesMessageCommandHandler(s, m, args)
}
func (h *CommandHandler) rolesSlash(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	return rolesSlashCommandHandler(s, i)
}

func (h *CommandHandler) guild(s *discordgo.Session, m *discordgo.MessageCreate, args []string) error {
	return guildMessageCommandHandler(s, m, args)
}
func (h *CommandHandler) guildSlash(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	return guildSlashCommandHandler(s, i)
}

func (h *CommandHandler) emojis(s *discordgo.Session, m *discordgo.MessageCreate, args []string) error {
	return emojisMessageCommandHandler(s, m, args)
}
func (h *CommandHandler) emojisSlash(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	return emojisSlashCommandHandler(s, i)
}

func (h *CommandHandler) stickers(s *discordgo.Session, m *discordgo.MessageCreate, args []string) error {
	return stickersMessageCommandHandler(s, m, args)
}
func (h *CommandHandler) stickersSlash(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	return stickersSlashCommandHandler(s, i)
}

func (h *CommandHandler) bans(s *discordgo.Session, m *discordgo.MessageCreate, args []string) error {
	return bansMessageCommandHandler(s, m, args)
}
func (h *CommandHandler) bansSlash(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	return bansSlashCommandHandler(s, i)
}

func (h *CommandHandler) boomer(s *discordgo.Session, m *discordgo.MessageCreate, args []string) error {
	return boomerMessageCommandHandler(s, m, args)
}
func (h *CommandHandler) boomerSlash(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	return boomerSlashCommandHandler(s, i)
}

func (h *CommandHandler) perms(s *discordgo.Session, m *discordgo.MessageCreate, args []string) error {
	return permsMessageCommandHandler(s, m, args)
}
func (h *CommandHandler) permsSlash(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	return permsSlashCommandHandler(s, i)
}

func (h *CommandHandler) tz(s *discordgo.Session, m *discordgo.MessageCreate, args []string) error {
	return tzMessageCommandHandler(s, m, args)
}
func (h *CommandHandler) tzSlash(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	return tzSlashCommandHandler(s, i)
}

func (h *CommandHandler) urbandictionary(s *discordgo.Session, m *discordgo.MessageCreate, args []string) error {
	return urbandictionaryMessageCommandHandler(s, m, args)
}
func (h *CommandHandler) urbandictionarySlash(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	return urbandictionarySlashCommandHandler(s, i)
}

func (h *CommandHandler) ttys(s *discordgo.Session, m *discordgo.MessageCreate, args []string) error {
	return ttysMessageCommandHandler(s, m, args)
}
func (h *CommandHandler) ttysSlash(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	return ttysSlashCommandHandler(s, i)
}

func (h *CommandHandler) ddg(s *discordgo.Session, channelID string, args []string) error {
	return ddgMessageCommandHandler(s, channelID, args)
}

func (h *CommandHandler) ddgSlash(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	return ddgSlashCommandHandler(s, i)
}

func (h *CommandHandler) search(s *discordgo.Session, m *discordgo.MessageCreate, args []string) error {
	return searchMessageCommandHandler(s, m, args)
}

func (h *CommandHandler) searchSlash(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	return searchSlashCommandHandler(s, i)
}

func (h *CommandHandler) pinglist(s *discordgo.Session, m *discordgo.MessageCreate, args []string) error {
	return pinglistMessageCommandHandler(s, m, args)
}

func (h *CommandHandler) pinglistSlash(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	return pinglistSlashCommandHandler(s, i)
}

func (h *CommandHandler) gifsearch(s *discordgo.Session, channelID string, args []string) error {
	return gifsearchMessageCommandHandler(s, channelID, args)
}

func (h *CommandHandler) gifsearchSlash(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	return gifsearchSlashCommandHandler(s, i)
}

func (h *CommandHandler) emoji(s *discordgo.Session, m *discordgo.MessageCreate, args []string) error {
	return emojiMessageCommandHandler(s, m, args)
}

func (h *CommandHandler) emojiSlash(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	return emojiSlashCommandHandler(s, i)
}

func (h *CommandHandler) sticker(s *discordgo.Session, m *discordgo.MessageCreate, args []string) error {
	return stickerMessageCommandHandler(s, m, args)
}

func (h *CommandHandler) stickerSlash(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	return stickerSlashCommandHandler(s, i)
}

func (h *CommandHandler) pin(s *discordgo.Session, m *discordgo.MessageCreate, args []string) error {
	return pinMessageCommandHandler(s, m, args)
}

func (h *CommandHandler) pinSlash(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	return pinSlashCommandHandler(s, i)
}

func (h *CommandHandler) unpin(s *discordgo.Session, m *discordgo.MessageCreate, args []string) error {
	return unpinMessageCommandHandler(s, m, args)
}

func (h *CommandHandler) unpinSlash(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	return unpinSlashCommandHandler(s, i)
}

func (h *CommandHandler) quote(s *discordgo.Session, m *discordgo.MessageCreate, args []string) error {
	return quoteMessageCommandHandler(s, m, args)
}

func (h *CommandHandler) quoteSlash(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	return quoteSlashCommandHandler(s, i)
}

func (h *CommandHandler) translate(s *discordgo.Session, m *discordgo.MessageCreate, args []string) error {
	return translateMessageCommandHandler(s, m, args)
}

func (h *CommandHandler) translateSlash(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	return translateSlashCommandHandler(s, i)
}

func (h *CommandHandler) reminder(s *discordgo.Session, m *discordgo.MessageCreate, args []string) error {
	return reminderMessageCommandHandler(s, m, args)
}

func (h *CommandHandler) reminderSlash(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	return reminderSlashCommandHandler(s, i)
}

func (h *CommandHandler) isearch(s *discordgo.Session, m *discordgo.MessageCreate, args []string) error {
	return isearchMessageCommandHandler(s, m, args)
}

func (h *CommandHandler) isearchSlash(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	return isearchSlashCommandHandler(s, i)
}

func (h *CommandHandler) caption(s *discordgo.Session, m *discordgo.MessageCreate, args []string) error {
	return captionMessageCommandHandler(s, m, args)
}

func (h *CommandHandler) captionSlash(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	return captionSlashCommandHandler(s, i)
}

func (h *CommandHandler) img2gif(s *discordgo.Session, m *discordgo.MessageCreate, args []string) error {
	return img2gifMessageCommandHandler(s, m, args)
}

func (h *CommandHandler) img2gifSlash(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	return img2gifSlashCommandHandler(s, i)
}

func (h *CommandHandler) vid2gif(s *discordgo.Session, m *discordgo.MessageCreate, args []string) error {
	return vid2gifMessageCommandHandler(s, m, args)
}

func (h *CommandHandler) vid2gifSlash(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	return vid2gifSlashCommandHandler(s, i)
}

func (h *CommandHandler) autogif(s *discordgo.Session, m *discordgo.MessageCreate, args []string) error {
	return autogifMessageCommandHandler(s, m, args)
}

func (h *CommandHandler) autogifSlash(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	return autogifSlashCommandHandler(s, i)
}

func (h *CommandHandler) factcheck(s *discordgo.Session, m *discordgo.MessageCreate, args []string) error {
	return factcheckMessageCommandHandler(s, m, args, h.Cfg)
}

func (h *CommandHandler) factcheckSlash(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	return factcheckSlashCommandHandler(s, i, h.Cfg)
}

func (h *CommandHandler) howgay(s *discordgo.Session, m *discordgo.MessageCreate, args []string) error {
	return howgayMessageCommandHandler(s, m, args)
}
func (h *CommandHandler) howgaySlash(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	return howgaySlashCommandHandler(s, i)
}

func (h *CommandHandler) howautism(s *discordgo.Session, m *discordgo.MessageCreate, args []string) error {
	return howautismMessageCommandHandler(s, m, args)
}
func (h *CommandHandler) howautismSlash(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	return howautismSlashCommandHandler(s, i)
}

func (h *CommandHandler) howlesbian(s *discordgo.Session, m *discordgo.MessageCreate, args []string) error {
	return howlesbianMessageCommandHandler(s, m, args)
}
func (h *CommandHandler) howlesbianSlash(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	return howlesbianSlashCommandHandler(s, i)
}

func (h *CommandHandler) howsimp(s *discordgo.Session, m *discordgo.MessageCreate, args []string) error {
	return howsimpMessageCommandHandler(s, m, args)
}
func (h *CommandHandler) howsimpSlash(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	return howsimpSlashCommandHandler(s, i)
}

func (h *CommandHandler) pp(s *discordgo.Session, m *discordgo.MessageCreate, args []string) error {
	return ppMessageCommandHandler(s, m, args)
}
func (h *CommandHandler) ppSlash(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	return ppSlashCommandHandler(s, i)
}

func (h *CommandHandler) puh(s *discordgo.Session, m *discordgo.MessageCreate, args []string) error {
	return puhMessageCommandHandler(s, m, args)
}
func (h *CommandHandler) puhSlash(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	return puhSlashCommandHandler(s, i)
}

func (h *CommandHandler) iq(s *discordgo.Session, m *discordgo.MessageCreate, args []string) error {
	return iqMessageCommandHandler(s, m, args)
}
func (h *CommandHandler) iqSlash(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	return iqSlashCommandHandler(s, i)
}

func (h *CommandHandler) bitches(s *discordgo.Session, m *discordgo.MessageCreate, args []string) error {
	return bitchesMessageCommandHandler(s, m, args)
}
func (h *CommandHandler) bitchesSlash(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	return bitchesSlashCommandHandler(s, i)
}

func (h *CommandHandler) choose(s *discordgo.Session, m *discordgo.MessageCreate, args []string) error {
	return chooseMessageCommandHandler(s, m, args)
}
func (h *CommandHandler) chooseSlash(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	return chooseSlashCommandHandler(s, i)
}

func (h *CommandHandler) ship(s *discordgo.Session, m *discordgo.MessageCreate, args []string) error {
	return shipMessageCommandHandler(s, m, args)
}
func (h *CommandHandler) shipSlash(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	return shipSlashCommandHandler(s, i)
}

func (h *CommandHandler) colorsAvatar(s *discordgo.Session, m *discordgo.MessageCreate, args []string) error {
	return colorsAvatarMessageCommandHandler(s, m, args)
}
func (h *CommandHandler) colorsAvatarSlash(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	return colorsAvatarSlashCommandHandler(s, i)
}

func (h *CommandHandler) lurk(s *discordgo.Session, m *discordgo.MessageCreate, args []string) error {
	return lurkMessageCommandHandler(s, m, args)
}
func (h *CommandHandler) lurkSlash(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	return lurkSlashCommandHandler(s, i)
}

func (h *CommandHandler) lurkers(s *discordgo.Session, m *discordgo.MessageCreate, args []string) error {
	return lurkersMessageCommandHandler(s, m, args)
}
func (h *CommandHandler) lurkersSlash(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	return lurkersSlashCommandHandler(s, i)
}

func (h *CommandHandler) spark(s *discordgo.Session, m *discordgo.MessageCreate, args []string) error {
	return sparkMessageCommandHandler(s, m, args)
}
func (h *CommandHandler) sparkSlash(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	return sparkSlashCommandHandler(s, i)
}

func (h *CommandHandler) smoke(s *discordgo.Session, m *discordgo.MessageCreate, args []string) error {
	return smokeMessageCommandHandler(s, m, args)
}
func (h *CommandHandler) smokeSlash(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	return smokeSlashCommandHandler(s, i)
}

func (h *CommandHandler) hits(s *discordgo.Session, m *discordgo.MessageCreate, args []string) error {
	return hitsMessageCommandHandler(s, m, args)
}
func (h *CommandHandler) hitsSlash(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	return hitsSlashCommandHandler(s, i)
}

func (h *CommandHandler) compress(s *discordgo.Session, m *discordgo.MessageCreate, args []string) error {
	return compressMessageCommandHandler(s, m, args)
}
func (h *CommandHandler) compressSlash(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	return compressSlashCommandHandler(s, i)
}

func (h *CommandHandler) vape(s *discordgo.Session, m *discordgo.MessageCreate, args []string) error {
	return vapeMessageCommandHandler(s, m, args)
}
func (h *CommandHandler) vapeSlash(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	return vapeSlashCommandHandler(s, i)
}

func (h *CommandHandler) poll(s *discordgo.Session, m *discordgo.MessageCreate, args []string) error {
	return pollMessageCommandHandler(s, m, args)
}
func (h *CommandHandler) pollSlash(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	return pollSlashCommandHandler(s, i)
}

func (h *CommandHandler) quickpoll(s *discordgo.Session, m *discordgo.MessageCreate, args []string) error {
	return quickpollMessageCommandHandler(s, m, args)
}
func (h *CommandHandler) quickpollSlash(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	return quickpollSlashCommandHandler(s, i)
}

func (h *CommandHandler) birthday(s *discordgo.Session, m *discordgo.MessageCreate, args []string) error {
	return birthdayMessageCommandHandler(h, s, m, args)
}
func (h *CommandHandler) birthdaySlash(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	return birthdaySlashCommandHandler(h, s, i)
}

func (h *CommandHandler) diary(s *discordgo.Session, m *discordgo.MessageCreate, args []string) error {
	return diaryMessageCommandHandler(h, s, m, args)
}
func (h *CommandHandler) diarySlash(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	return diarySlashCommandHandler(h, s, i)
}

func (h *CommandHandler) bio(s *discordgo.Session, m *discordgo.MessageCreate, args []string) error {
	return bioMessageCommandHandler(s, m, args)
}
func (h *CommandHandler) bioSlash(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	return bioSlashCommandHandler(s, i, h.Cfg)
}

func (h *CommandHandler) ban(s *discordgo.Session, m *discordgo.MessageCreate, args []string) error {
	return banMessageCommandHandler(h, s, m, args)
}
func (h *CommandHandler) banSlash(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	return banSlashCommandHandler(h, s, i)
}
func (h *CommandHandler) hardban(s *discordgo.Session, m *discordgo.MessageCreate, args []string) error {
	return hardbanMessageCommandHandler(h, s, m, args)
}
func (h *CommandHandler) hardbanSlash(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	return hardbanSlashCommandHandler(h, s, i)
}
func (h *CommandHandler) softban(s *discordgo.Session, m *discordgo.MessageCreate, args []string) error {
	return softbanMessageCommandHandler(h, s, m, args)
}
func (h *CommandHandler) softbanSlash(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	return softbanSlashCommandHandler(h, s, i)
}
func (h *CommandHandler) kick(s *discordgo.Session, m *discordgo.MessageCreate, args []string) error {
	return kickMessageCommandHandler(h, s, m, args)
}
func (h *CommandHandler) kickSlash(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	return kickSlashCommandHandler(h, s, i)
}
func (h *CommandHandler) purge(s *discordgo.Session, m *discordgo.MessageCreate, args []string) error {
	return purgeMessageCommandHandler(h, s, m, args)
}
func (h *CommandHandler) purgeSlash(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	return purgeSlashCommandHandler(h, s, i)
}
func (h *CommandHandler) nuke(s *discordgo.Session, m *discordgo.MessageCreate, args []string) error {
	return nukeMessageCommandHandler(h, s, m, args)
}
func (h *CommandHandler) nukeSlash(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	return nukeSlashCommandHandler(h, s, i)
}
func (h *CommandHandler) timeout(s *discordgo.Session, m *discordgo.MessageCreate, args []string) error {
	return timeoutMessageCommandHandler(h, s, m, args)
}
func (h *CommandHandler) timeoutSlash(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	return timeoutSlashCommandHandler(h, s, i)
}
func (h *CommandHandler) warn(s *discordgo.Session, m *discordgo.MessageCreate, args []string) error {
	return warnMessageCommandHandler(h, s, m, args)
}
func (h *CommandHandler) warnSlash(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	return warnSlashCommandHandler(h, s, i)
}
func (h *CommandHandler) history(s *discordgo.Session, m *discordgo.MessageCreate, args []string) error {
	return historyMessageCommandHandler(h, s, m, args)
}
func (h *CommandHandler) historySlash(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	return historySlashCommandHandler(h, s, i)
}

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
}

func userOption(required bool) *discordgo.ApplicationCommandOption {
	return &discordgo.ApplicationCommandOption{
		Type:        discordgo.ApplicationCommandOptionUser,
		Name:        "user",
		Description: "The user (defaults to yourself)",
		Required:    required,
	}
}
