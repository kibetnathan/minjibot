// Package commands implements all of MinjiBot's bot commands. Each command
// has both a prefix handler (e.g. foo(s, m, args)) and a slash handler
// (e.g. fooSlash(s, i)). The CommandHandler struct holds repository
// dependencies and dispatches to the appropriate handler.
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
	case "weather":
		return h.weather(s, m, args)
	case "test":
		return h.test(s, m, args)
	case "donate":
		return h.donate(s, m, args)
	case "setup":
		return h.setup(s, m, args)
	case "bug":
		return h.bug(s, m, args)
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
	case "audit":
		return h.audit(s, m, args)
	case "role":
		return h.role(s, m, args)
	case "fn":
		return h.fn(s, m, args)
	case "nick":
		return h.nick(s, m, args)
	case "jail":
		return h.jail(s, m, args)
	case "unjail":
		return h.unjail(s, m, args)
	case "staffstrip":
		return h.staffstrip(s, m, args)
	case "hide":
		return h.hide(s, m, args)
	case "reveal":
		return h.reveal(s, m, args)
	case "lockdown":
		return h.lockdown(s, m, args)
	case "nsfw":
		return h.nsfw(s, m, args)
	case "sfw":
		return h.sfw(s, m, args)
	case "slowmode":
		return h.slowmode(s, m, args)
	case "topic":
		return h.topic(s, m, args)
	case "denyperm":
		return h.denyperm(s, m, args)
	case "imute":
		return h.imute(s, m, args)
	case "gifmute":
		return h.gifmute(s, m, args)
	case "angry":
		return rpEmotionMessage(s, m, "angry")
	case "depressed":
		return rpEmotionMessage(s, m, "depressed")
	case "excited":
		return rpEmotionMessage(s, m, "excited")
	case "happy":
		return rpEmotionMessage(s, m, "happy")
	case "horny":
		return rpEmotionMessage(s, m, "horny")
	case "inlove":
		return rpEmotionMessage(s, m, "inlove")
	case "sad":
		return rpEmotionMessage(s, m, "sad")
	case "shy":
		return rpEmotionMessage(s, m, "shy")
	case "baka":
		return rpActionMessage(s, m, args, "baka")
	case "bite":
		return rpActionMessage(s, m, args, "bite")
	case "cry":
		return rpActionMessage(s, m, args, "cry")
	case "dap":
		return rpActionMessage(s, m, args, "dap")
	case "eat":
		return rpActionMessage(s, m, args, "eat")
	case "facepalm":
		return rpActionMessage(s, m, args, "facepalm")
	case "feed":
		return rpActionMessage(s, m, args, "feed")
	case "handhold":
		return rpActionMessage(s, m, args, "handhold")
	case "kiss":
		return rpActionMessage(s, m, args, "kiss")
	case "laugh":
		return rpActionMessage(s, m, args, "laugh")
	case "nod":
		return rpActionMessage(s, m, args, "nod")
	case "nutkick":
		return rpActionMessage(s, m, args, "nutkick")
	case "pat":
		return rpActionMessage(s, m, args, "pat")
	case "peck":
		return rpActionMessage(s, m, args, "peck")
	case "poke":
		return rpActionMessage(s, m, args, "poke")
	case "punch":
		return rpActionMessage(s, m, args, "punch")
	case "run":
		return rpActionMessage(s, m, args, "run")
	case "shoot":
		return rpActionMessage(s, m, args, "shoot")
	case "shrug":
		return rpActionMessage(s, m, args, "shrug")
	case "slap":
		return rpActionMessage(s, m, args, "slap")
	case "spank":
		return rpActionMessage(s, m, args, "spank")
	case "stab":
		return rpActionMessage(s, m, args, "stab")
	case "think":
		return rpActionMessage(s, m, args, "think")
	case "tickle":
		return rpActionMessage(s, m, args, "tickle")
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
	case "weather":
		return h.weatherSlash(s, i)
	case "test":
		return h.testSlash(s, i)
	case "donate":
		return h.donateSlash(s, i)
	case "setup":
		return h.setupSlash(s, i)
	case "bug":
		return h.bugSlash(s, i)
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
	case "audit":
		return h.auditSlash(s, i)
	case "role":
		return h.roleSlash(s, i)
	case "fn":
		return h.fnSlash(s, i)
	case "nick":
		return h.nickSlash(s, i)
	case "jail":
		return h.jailSlash(s, i)
	case "unjail":
		return h.unjailSlash(s, i)
	case "staffstrip":
		return h.staffstripSlash(s, i)
	case "hide":
		return h.hideSlash(s, i)
	case "reveal":
		return h.revealSlash(s, i)
	case "lockdown":
		return h.lockdownSlash(s, i)
	case "nsfw":
		return h.nsfwSlash(s, i)
	case "sfw":
		return h.sfwSlash(s, i)
	case "slowmode":
		return h.slowmodeSlash(s, i)
	case "topic":
		return h.topicSlash(s, i)
	case "denyperm":
		return h.denypermSlash(s, i)
	case "imute":
		return h.imuteSlash(s, i)
	case "gifmute":
		return h.gifmuteSlash(s, i)
	case "angry":
		return rpEmotionSlash(s, i)
	case "depressed":
		return rpEmotionSlash(s, i)
	case "excited":
		return rpEmotionSlash(s, i)
	case "happy":
		return rpEmotionSlash(s, i)
	case "horny":
		return rpEmotionSlash(s, i)
	case "inlove":
		return rpEmotionSlash(s, i)
	case "sad":
		return rpEmotionSlash(s, i)
	case "shy":
		return rpEmotionSlash(s, i)
	case "baka":
		return rpActionSlash(s, i)
	case "bite":
		return rpActionSlash(s, i)
	case "cry":
		return rpActionSlash(s, i)
	case "dap":
		return rpActionSlash(s, i)
	case "eat":
		return rpActionSlash(s, i)
	case "facepalm":
		return rpActionSlash(s, i)
	case "feed":
		return rpActionSlash(s, i)
	case "handhold":
		return rpActionSlash(s, i)
	case "kiss":
		return rpActionSlash(s, i)
	case "laugh":
		return rpActionSlash(s, i)
	case "nod":
		return rpActionSlash(s, i)
	case "nutkick":
		return rpActionSlash(s, i)
	case "pat":
		return rpActionSlash(s, i)
	case "peck":
		return rpActionSlash(s, i)
	case "poke":
		return rpActionSlash(s, i)
	case "punch":
		return rpActionSlash(s, i)
	case "run":
		return rpActionSlash(s, i)
	case "shoot":
		return rpActionSlash(s, i)
	case "shrug":
		return rpActionSlash(s, i)
	case "slap":
		return rpActionSlash(s, i)
	case "spank":
		return rpActionSlash(s, i)
	case "stab":
		return rpActionSlash(s, i)
	case "think":
		return rpActionSlash(s, i)
	case "tickle":
		return rpActionSlash(s, i)
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
	// -help <category> opens on that category; plain -help opens on General.
	// Navigation is always via buttons, matching the slash command.
	page := 0
	if len(args) > 0 {
		if idx := FindHelpSection(args[0]); idx >= 0 {
			page = idx
		}
	}
	msg, err := s.ChannelMessageSendComplex(m.ChannelID, &discordgo.MessageSend{
		Embeds:     []*discordgo.MessageEmbed{BuildHelpPageEmbed(page)},
		Components: helpButtonsRow(),
	})
	if err != nil {
		return err
	}
	helpPageMu.Lock()
	helpPageByMsg[msg.ID] = page
	helpPageMu.Unlock()
	return nil
}

func (h *CommandHandler) helpSlash(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds:     []*discordgo.MessageEmbed{BuildHelpPageEmbed(0)},
			Components: helpButtonsRow(),
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

func (h *CommandHandler) weather(s *discordgo.Session, m *discordgo.MessageCreate, args []string) error {
	return weatherMessageCommandHandler(s, m, args)
}
func (h *CommandHandler) weatherSlash(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	return weatherSlashCommandHandler(s, i)
}

func (h *CommandHandler) test(s *discordgo.Session, m *discordgo.MessageCreate, args []string) error {
	return testMessageCommandHandler(s, m, args)
}
func (h *CommandHandler) testSlash(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	return testSlashCommandHandler(s, i)
}

func (h *CommandHandler) donate(s *discordgo.Session, m *discordgo.MessageCreate, args []string) error {
	return donateMessageCommandHandler(s, m, args)
}
func (h *CommandHandler) donateSlash(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	return donateSlashCommandHandler(s, i)
}

func (h *CommandHandler) bug(s *discordgo.Session, m *discordgo.MessageCreate, args []string) error {
	return bugMessageCommandHandler(s, m, args)
}
func (h *CommandHandler) bugSlash(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	return bugSlashCommandHandler(s, i)
}

func (h *CommandHandler) setup(s *discordgo.Session, m *discordgo.MessageCreate, args []string) error {
	return setupMessageCommandHandler(h, s, m, args)
}
func (h *CommandHandler) setupSlash(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	return setupSlashCommandHandler(h, s, i)
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
func (h *CommandHandler) audit(s *discordgo.Session, m *discordgo.MessageCreate, args []string) error {
	return auditMessageCommandHandler(h, s, m, args)
}
func (h *CommandHandler) auditSlash(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	return auditSlashCommandHandler(h, s, i)
}
func (h *CommandHandler) role(s *discordgo.Session, m *discordgo.MessageCreate, args []string) error {
	return roleMessageCommandHandler(h, s, m, args)
}
func (h *CommandHandler) roleSlash(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	return roleSlashCommandHandler(h, s, i)
}
func (h *CommandHandler) fn(s *discordgo.Session, m *discordgo.MessageCreate, args []string) error {
	return fnMessageCommandHandler(h, s, m, args)
}
func (h *CommandHandler) fnSlash(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	return fnSlashCommandHandler(h, s, i)
}
func (h *CommandHandler) nick(s *discordgo.Session, m *discordgo.MessageCreate, args []string) error {
	return nickMessageCommandHandler(h, s, m, args)
}
func (h *CommandHandler) nickSlash(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	return nickSlashCommandHandler(h, s, i)
}
func (h *CommandHandler) jail(s *discordgo.Session, m *discordgo.MessageCreate, args []string) error {
	return jailMessageCommandHandler(h, s, m, args)
}
func (h *CommandHandler) jailSlash(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	return jailSlashCommandHandler(h, s, i)
}
func (h *CommandHandler) unjail(s *discordgo.Session, m *discordgo.MessageCreate, args []string) error {
	return unjailMessageCommandHandler(h, s, m, args)
}
func (h *CommandHandler) unjailSlash(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	return unjailSlashCommandHandler(h, s, i)
}
func (h *CommandHandler) staffstrip(s *discordgo.Session, m *discordgo.MessageCreate, args []string) error {
	return staffstripMessageCommandHandler(h, s, m, args)
}
func (h *CommandHandler) staffstripSlash(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	return staffstripSlashCommandHandler(h, s, i)
}
func (h *CommandHandler) hide(s *discordgo.Session, m *discordgo.MessageCreate, args []string) error {
	return hideMessageCommandHandler(h, s, m, args)
}
func (h *CommandHandler) hideSlash(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	return hideSlashCommandHandler(h, s, i)
}
func (h *CommandHandler) reveal(s *discordgo.Session, m *discordgo.MessageCreate, args []string) error {
	return revealMessageCommandHandler(h, s, m, args)
}
func (h *CommandHandler) revealSlash(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	return revealSlashCommandHandler(h, s, i)
}
func (h *CommandHandler) lockdown(s *discordgo.Session, m *discordgo.MessageCreate, args []string) error {
	return lockdownMessageCommandHandler(h, s, m, args)
}
func (h *CommandHandler) lockdownSlash(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	return lockdownSlashCommandHandler(h, s, i)
}
func (h *CommandHandler) nsfw(s *discordgo.Session, m *discordgo.MessageCreate, args []string) error {
	return nsfwMessageCommandHandler(h, s, m, args)
}
func (h *CommandHandler) nsfwSlash(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	return nsfwSlashCommandHandler(h, s, i)
}
func (h *CommandHandler) sfw(s *discordgo.Session, m *discordgo.MessageCreate, args []string) error {
	return sfwMessageCommandHandler(h, s, m, args)
}
func (h *CommandHandler) sfwSlash(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	return sfwSlashCommandHandler(h, s, i)
}
func (h *CommandHandler) slowmode(s *discordgo.Session, m *discordgo.MessageCreate, args []string) error {
	return slowmodeMessageCommandHandler(h, s, m, args)
}
func (h *CommandHandler) slowmodeSlash(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	return slowmodeSlashCommandHandler(h, s, i)
}
func (h *CommandHandler) topic(s *discordgo.Session, m *discordgo.MessageCreate, args []string) error {
	return topicMessageCommandHandler(h, s, m, args)
}
func (h *CommandHandler) topicSlash(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	return topicSlashCommandHandler(h, s, i)
}
func (h *CommandHandler) denyperm(s *discordgo.Session, m *discordgo.MessageCreate, args []string) error {
	return denypermMessageCommandHandler(h, s, m, args)
}
func (h *CommandHandler) denypermSlash(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	return denypermSlashCommandHandler(h, s, i)
}
func (h *CommandHandler) imute(s *discordgo.Session, m *discordgo.MessageCreate, args []string) error {
	return imuteMessageCommandHandler(h, s, m, args)
}
func (h *CommandHandler) imuteSlash(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	return imuteSlashCommandHandler(h, s, i)
}
func (h *CommandHandler) gifmute(s *discordgo.Session, m *discordgo.MessageCreate, args []string) error {
	return gifmuteMessageCommandHandler(h, s, m, args)
}
func (h *CommandHandler) gifmuteSlash(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	return gifmuteSlashCommandHandler(h, s, i)
}
