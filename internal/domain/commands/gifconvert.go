package commands

import (
	"bytes"
	"fmt"
	"image"
	"image/color/palette"
	"image/draw"
	"image/gif"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func isVideoExt(ext string) bool {
	switch ext {
	case "mp4", "webm", "mov", "mkv", "avi", "m4v":
		return true
	}
	return false
}

func isImageExt(ext string) bool {
	switch ext {
	case "png", "jpg", "jpeg", "gif", "webp", "bmp":
		return true
	}
	return false
}

// imageBytesToGIF converts a static image into an animated GIF (two gentle
// frames so Discord treats it as a GIF). Uses only the standard library.
func imageBytesToGIF(data []byte, ext string) ([]byte, error) {
	src, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("decoding image: %w", err)
	}

	// If the source is already an animated GIF, pass it through untouched.
	if ext == "gif" {
		if g, err := gif.DecodeAll(bytes.NewReader(data)); err == nil && len(g.Image) > 1 {
			var buf bytes.Buffer
			if err := gif.EncodeAll(&buf, g); err == nil {
				return buf.Bytes(), nil
			}
		}
	}

	bounds := src.Bounds()
	paletted := image.NewPaletted(bounds, palette.Plan9[:256])
	draw.Draw(paletted, bounds, src, image.Point{}, draw.Src)

	anim := &gif.GIF{
		Image: []*image.Paletted{paletted, paletted},
		Delay: []int{5, 5},
	}
	var buf bytes.Buffer
	if err := gif.EncodeAll(&buf, anim); err != nil {
		return nil, fmt.Errorf("encoding gif: %w", err)
	}
	return buf.Bytes(), nil
}

// ffmpegToGIF converts an uploaded video to an animated GIF via ffmpeg.
func ffmpegToGIF(data []byte, ext string, startSec, durSec float64) ([]byte, error) {
	if _, err := exec.LookPath("ffmpeg"); err != nil {
		return nil, fmt.Errorf("ffmpeg is not installed on this machine — can't convert videos to GIF")
	}

	dir, err := os.MkdirTemp("", "minjibot-gif-*")
	if err != nil {
		return nil, err
	}
	defer os.RemoveAll(dir)

	in := filepath.Join(dir, "input."+ext)
	if err := os.WriteFile(in, data, 0o600); err != nil {
		return nil, err
	}
	out := filepath.Join(dir, "output.gif")

	args := []string{"-y", "-i", in}
	if startSec > 0 {
		args = append(args, "-ss", fmt.Sprintf("%.2f", startSec))
	}
	if durSec > 0 {
		args = append(args, "-t", fmt.Sprintf("%.2f", durSec))
	}
	args = append(args, "-vf", "fps=15,scale=480:-1:flags=lanczos", "-loop", "0", out)

	cmd := exec.Command("ffmpeg", args...)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("ffmpeg failed: %v: %s", err, strings.TrimSpace(stderr.String()))
	}

	result, err := os.ReadFile(out)
	if err != nil {
		return nil, fmt.Errorf("reading converted gif: %w", err)
	}
	return result, nil
}

func slugify(s string) string {
	var b strings.Builder
	for _, r := range strings.ToLower(s) {
		switch {
		case r >= 'a' && r <= 'z', r >= '0' && r <= '9':
			b.WriteRune(r)
		case r == ' ', r == '-', r == '_':
			b.WriteRune('-')
		}
		if b.Len() >= 40 {
			break
		}
	}
	return strings.Trim(b.String(), "-")
}

func parseConvertArgs(args []string) (startSec, durSec float64, url string) {
	for _, arg := range args {
		switch {
		case strings.HasPrefix(strings.ToLower(arg), "start:"):
			fmt.Sscanf(strings.TrimSpace(arg[len("start:"):]), "%f", &startSec)
		case strings.HasPrefix(strings.ToLower(arg), "dur:"), strings.HasPrefix(strings.ToLower(arg), "trim:"):
			i := strings.Index(arg, ":")
			fmt.Sscanf(strings.TrimSpace(arg[i+1:]), "%f", &durSec)
		case strings.HasPrefix(arg, "http://"), strings.HasPrefix(arg, "https://"):
			url = arg
		}
	}
	return
}

func img2gifMessageCommandHandler(s *discordgo.Session, m *discordgo.MessageCreate, args []string) error {
	_, _, url := parseConvertArgs(args)
	md, err := resolveMedia(s, m, url)
	if err != nil {
		return err
	}
	if isVideoExt(md.Ext) {
		return fmt.Errorf("that's a video — use `!vid2gif` instead")
	}
	if !isImageExt(md.Ext) {
		return fmt.Errorf("unsupported image type: %s", orEmpty(md.Ext))
	}

	out, err := imageBytesToGIF(md.Data, md.Ext)
	if err != nil {
		return err
	}

	_, err = s.ChannelFileSend(m.ChannelID, "img2gif.gif", bytes.NewReader(out))
	return err
}

func img2gifSlashCommandHandler(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	opts := optionMap(i.ApplicationCommandData().Options)
	md, err := fetchURL(optString(opts, "url"))
	if err != nil {
		return err
	}
	if !isImageExt(md.Ext) {
		return fmt.Errorf("unsupported image type: %s", orEmpty(md.Ext))
	}

	out, err := imageBytesToGIF(md.Data, md.Ext)
	if err != nil {
		return err
	}

	return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{Files: []*discordgo.File{{Name: "img2gif.gif", Reader: bytes.NewReader(out)}}},
	})
}

func vid2gifMessageCommandHandler(s *discordgo.Session, m *discordgo.MessageCreate, args []string) error {
	startSec, durSec, url := parseConvertArgs(args)
	md, err := resolveMedia(s, m, url)
	if err != nil {
		return err
	}
	if !isVideoExt(md.Ext) {
		return fmt.Errorf("that file isn't a video (got %q) — use `!img2gif` for images", orEmpty(md.Ext))
	}

	out, err := ffmpegToGIF(md.Data, md.Ext, startSec, durSec)
	if err != nil {
		return err
	}

	_, err = s.ChannelFileSend(m.ChannelID, "vid2gif.gif", bytes.NewReader(out))
	return err
}

func vid2gifSlashCommandHandler(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	opts := optionMap(i.ApplicationCommandData().Options)
	md, err := fetchURL(optString(opts, "url"))
	if err != nil {
		return err
	}
	if !isVideoExt(md.Ext) {
		return fmt.Errorf("that file isn't a video (got %q)", orEmpty(md.Ext))
	}

	out, err := ffmpegToGIF(md.Data, md.Ext, 0, 0)
	if err != nil {
		return err
	}

	return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{Files: []*discordgo.File{{Name: "vid2gif.gif", Reader: bytes.NewReader(out)}}},
	})
}

// autogif fetches media from a reply/attachment/URL and posts it as a GIF,
// converting images or videos accordingly. Already-animated GIFs pass through.
func autogifMessageCommandHandler(s *discordgo.Session, m *discordgo.MessageCreate, args []string) error {
	_, _, url := parseConvertArgs(args)
	if url == "" && m.ReferencedMessage == nil && len(m.Attachments) == 0 {
		_, err := s.ChannelMessageSend(m.ChannelID, "Usage: `!autogif` on a message with media, or `!autogif <url>`")
		return err
	}

	md, err := resolveMedia(s, m, url)
	if err != nil {
		return err
	}

	var out []byte
	switch {
	case isVideoExt(md.Ext):
		out, err = ffmpegToGIF(md.Data, md.Ext, 0, 0)
	case isImageExt(md.Ext):
		out, err = imageBytesToGIF(md.Data, md.Ext)
	default:
		return fmt.Errorf("unsupported media type: %s", orEmpty(md.Ext))
	}
	if err != nil {
		return err
	}

	_, err = s.ChannelFileSend(m.ChannelID, "autogif.gif", bytes.NewReader(out))
	return err
}

func autogifSlashCommandHandler(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	opts := optionMap(i.ApplicationCommandData().Options)
	md, err := fetchURL(optString(opts, "url"))
	if err != nil {
		return err
	}

	var out []byte
	switch {
	case isVideoExt(md.Ext):
		out, err = ffmpegToGIF(md.Data, md.Ext, 0, 0)
	case isImageExt(md.Ext):
		out, err = imageBytesToGIF(md.Data, md.Ext)
	default:
		return fmt.Errorf("unsupported media type: %s", orEmpty(md.Ext))
	}
	if err != nil {
		return err
	}

	return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{Files: []*discordgo.File{{Name: "autogif.gif", Reader: bytes.NewReader(out)}}},
	})
}

func orEmpty(s string) string {
	if s == "" {
		return "unknown"
	}
	return s
}
