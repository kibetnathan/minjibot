package commands

import (
	"bytes"
	"fmt"
	"image"
	"image/gif"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func IsVideoExt(ext string) bool {
	switch ext {
	case "mp4", "webm", "mov", "mkv", "avi", "m4v":
		return true
	}
	return false
}

func IsImageExt(ext string) bool {
	switch ext {
	case "png", "jpg", "jpeg", "gif", "webp", "bmp":
		return true
	}
	return false
}

// ImageBytesToGIF converts a static image into an animated GIF (two gentle
// frames so Discord treats it as a GIF). Uses only the standard library.
func ImageBytesToGIF(data []byte, ext string) ([]byte, error) {
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

	out, err := imageToGIF(src)
	if err != nil {
		return nil, fmt.Errorf("encoding gif: %w", err)
	}
	return out, nil
}

const (
	// maxVideoBytes caps how large an uploaded video may be before conversion.
	// Videos above this are rejected outright to keep conversions fast.
	maxVideoBytes = 25 * 1024 * 1024

	// defaultClipSec caps how much of a video is converted when the caller
	// doesn't request a duration. Trimming to a short clip keeps the GIF small
	// and fast, which matters because GIF encoding scales with frame count.
	defaultClipSec = 10.0

	// Video encoding knobs tuned for aggressive compression (small GIFs, fast
	// encodes). Lower fps/scale and stronger palette dithering all trade off a
	// little fidelity for a much smaller output.
	videoFPS  = "fps=12"
	videoScale = "scale=480:-1:flags=lanczos"
)

// ffmpegToGIF converts an uploaded video to an animated GIF via ffmpeg. It
// uses a two-pass palettegen/paletteuse filter and scales down aggressively so
// source videos render quickly into small GIFs.
func ffmpegToGIF(data []byte, ext string, startSec, durSec float64) ([]byte, error) {
	if len(data) > maxVideoBytes {
		return nil, fmt.Errorf("video is too large (%d bytes > %d max), try a shorter or smaller clip", len(data), maxVideoBytes)
	}
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
	palette := filepath.Join(dir, "palette.png")

	// Cap the clip length so long videos don't churn out huge, slow GIFs.
	if durSec <= 0 {
		durSec = defaultClipSec
	}

	// Shared clip selection: frame-accurate start/time on the input.
	clip := []string{}
	if startSec > 0 {
		clip = append(clip, "-ss", fmt.Sprintf("%.2f", startSec))
	}
	clip = append(clip, "-t", fmt.Sprintf("%.2f", durSec))

	// Pass 1: generate an optimized 256-colour palette for the clip.
	args := []string{"-y", "-i", in}
	args = append(args, clip...)
	args = append(args, "-vf", videoFPS+","+videoScale+",palettegen=stats_mode=diff", palette)
	if err := runFFmpeg(args); err != nil {
		return nil, err
	}

	// Pass 2: render the GIF using that palette with heavy Bayer dithering.
	args = []string{"-y", "-i", in, "-i", palette}
	args = append(args, clip...)
	args = append(args, "-lavfi", videoFPS+","+videoScale+"[x];[x][1:v]paletteuse=dither=bayer:bayer_scale=7:diff_mode=rectangle", "-loop", "0", out)
	if err := runFFmpeg(args); err != nil {
		return nil, err
	}

	result, err := os.ReadFile(out)
	if err != nil {
		return nil, fmt.Errorf("reading converted gif: %w", err)
	}
	return result, nil
}

func runFFmpeg(args []string) error {
	cmd := exec.Command("ffmpeg", args...)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("ffmpeg failed: %v: %s", err, strings.TrimSpace(stderr.String()))
	}
	return nil
}

func Slugify(s string) string {
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

func ParseConvertArgs(args []string) (startSec, durSec float64, url string) {
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
	_, _, url := ParseConvertArgs(args)
	md, err := resolveMedia(s, m, url)
	if err != nil {
		return err
	}
	if IsVideoExt(md.Ext) {
		return fmt.Errorf("that's a video — use `-vid2gif` instead")
	}
	if !IsImageExt(md.Ext) {
		return fmt.Errorf("unsupported image type: %s", OrEmpty(md.Ext))
	}

	out, err := ImageBytesToGIF(md.Data, md.Ext)
	if err != nil {
		return err
	}

	_, err = s.ChannelFileSend(m.ChannelID, "img2gif.gif", bytes.NewReader(out))
	return err
}

func img2gifSlashCommandHandler(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	opts := OptionMap(i.ApplicationCommandData().Options)
	md, err := fetchURL(OptString(opts, "url"))
	if err != nil {
		return err
	}
	if !IsImageExt(md.Ext) {
		return fmt.Errorf("unsupported image type: %s", OrEmpty(md.Ext))
	}

	out, err := ImageBytesToGIF(md.Data, md.Ext)
	if err != nil {
		return err
	}

	return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{Files: []*discordgo.File{{Name: "img2gif.gif", Reader: bytes.NewReader(out)}}},
	})
}

func vid2gifMessageCommandHandler(s *discordgo.Session, m *discordgo.MessageCreate, args []string) error {
	startSec, durSec, url := ParseConvertArgs(args)
	if url == "" {
		if att := mediaAttachment(m); att != nil && att.Size > maxVideoBytes {
			return fmt.Errorf("video is too large (%d bytes > %d max), try a shorter or smaller clip", att.Size, maxVideoBytes)
		}
	}
	md, err := resolveMedia(s, m, url)
	if err != nil {
		return err
	}
	if !IsVideoExt(md.Ext) {
		return fmt.Errorf("that file isn't a video (got %q) — use `-img2gif` for images", OrEmpty(md.Ext))
	}

	out, err := ffmpegToGIF(md.Data, md.Ext, startSec, durSec)
	if err != nil {
		return err
	}

	_, err = s.ChannelFileSend(m.ChannelID, "vid2gif.gif", bytes.NewReader(out))
	return err
}

// mediaAttachment mirrors resolveMedia's attachment selection (referenced
// message first, then the issuing message) and returns the chosen attachment.
func mediaAttachment(m *discordgo.MessageCreate) *discordgo.MessageAttachment {
	if att := FirstAttachment(m.ReferencedMessage); att != nil {
		return att
	}
	return FirstAttachment(m.Message)
}

func vid2gifSlashCommandHandler(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	opts := OptionMap(i.ApplicationCommandData().Options)
	md, err := fetchURL(OptString(opts, "url"))
	if err != nil {
		return err
	}
	if !IsVideoExt(md.Ext) {
		return fmt.Errorf("that file isn't a video (got %q)", OrEmpty(md.Ext))
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
	_, _, url := ParseConvertArgs(args)
	if url == "" && m.ReferencedMessage == nil && len(m.Attachments) == 0 {
		_, err := s.ChannelMessageSend(m.ChannelID, "Usage: `-autogif` on a message with media, or `-autogif <url>`")
		return err
	}
	if url == "" {
		if att := mediaAttachment(m); att != nil && att.Size > maxVideoBytes {
			return fmt.Errorf("video is too large (%d bytes > %d max), try a shorter or smaller clip", att.Size, maxVideoBytes)
		}
	}

	md, err := resolveMedia(s, m, url)
	if err != nil {
		return err
	}

	var out []byte
	switch {
	case IsVideoExt(md.Ext):
		out, err = ffmpegToGIF(md.Data, md.Ext, 0, 0)
	case IsImageExt(md.Ext):
		out, err = ImageBytesToGIF(md.Data, md.Ext)
	default:
		return fmt.Errorf("unsupported media type: %s", OrEmpty(md.Ext))
	}
	if err != nil {
		return err
	}

	_, err = s.ChannelFileSend(m.ChannelID, "autogif.gif", bytes.NewReader(out))
	return err
}

func autogifSlashCommandHandler(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	opts := OptionMap(i.ApplicationCommandData().Options)
	md, err := fetchURL(OptString(opts, "url"))
	if err != nil {
		return err
	}

	var out []byte
	switch {
	case IsVideoExt(md.Ext):
		out, err = ffmpegToGIF(md.Data, md.Ext, 0, 0)
	case IsImageExt(md.Ext):
		out, err = ImageBytesToGIF(md.Data, md.Ext)
	default:
		return fmt.Errorf("unsupported media type: %s", OrEmpty(md.Ext))
	}
	if err != nil {
		return err
	}

	return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{Files: []*discordgo.File{{Name: "autogif.gif", Reader: bytes.NewReader(out)}}},
	})
}

func OrEmpty(s string) string {
	if s == "" {
		return "unknown"
	}
	return s
}
