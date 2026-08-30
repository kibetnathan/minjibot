package commands

import (
	"bytes"
	"fmt"
	"image"
	"image/png"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

	"github.com/bwmarrin/discordgo"
)

// Pixelate blocks an image into a coarse grid of solid-colour cells — the
// "compressed until barely legible" effect. It downsamples every `cells`-wide
// region to a single (nearest-neighbour) colour, then scales each cell back up
// to the original dimensions so the result is a low-res mosaic at full size.
// Returns the re-encoded PNG bytes.
func Pixelate(data []byte, cells int) ([]byte, error) {
	if cells < 1 {
		cells = 1
	}

	src, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("decoding image: %w", err)
	}
	bounds := src.Bounds()
	w, h := bounds.Dx(), bounds.Dy()
	if w <= 0 || h <= 0 {
		return nil, fmt.Errorf("empty image")
	}

	cellW := maxInt(w/cells, 1)
	cellH := maxInt(h/cells, 1)

	// Nearest-neighbour downsample to the coarse grid.
	gridW, gridH := (w+cellW-1)/cellW, (h+cellH-1)/cellH
	grid := image.NewRGBA(image.Rect(0, 0, gridW, gridH))
	for gy := 0; gy < gridH; gy++ {
		for gx := 0; gx < gridW; gx++ {
			sx := minInt((gx+1)*cellW-1, w-1)
			sy := minInt((gy+1)*cellH-1, h-1)
			grid.Set(gx, gy, src.At(sx, sy))
		}
	}

	// Upscale the grid back to the original dimensions (hard blocky edges).
	out := image.NewRGBA(image.Rect(0, 0, w, h))
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			out.Set(x, y, grid.At(minInt(x/cellW, gridW-1), minInt(y/cellH, gridH-1)))
		}
	}

	var buf bytes.Buffer
	if err := png.Encode(&buf, out); err != nil {
		return nil, fmt.Errorf("encoding png: %w", err)
	}
	return buf.Bytes(), nil
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// compressCells is the default grid width used to make an image low-quality
// but still readable; the exact block size scales with the source dimensions.
const compressCells = 80

func compressMessageCommandHandler(s *discordgo.Session, m *discordgo.MessageCreate, args []string) error {
	if m.ReferencedMessage == nil && len(m.Attachments) == 0 && len(args) == 0 {
		_, err := s.ChannelMessageSend(m.ChannelID, "Usage: `-compress` on a message with an image, or `-compress <url>`")
		return err
	}

	url := ""
	for _, arg := range args {
		if len(arg) > 4 && (arg[:4] == "http" || arg[:5] == "https") {
			url = arg
			break
		}
	}

	md, err := resolveMedia(s, m, url)
	if err != nil {
		return err
	}
	if !IsImageExt(md.Ext) {
		return fmt.Errorf("unsupported image type: %s", OrEmpty(md.Ext))
	}

	out, err := Pixelate(md.Data, compressCells)
	if err != nil {
		return err
	}

	_, err = s.ChannelFileSend(m.ChannelID, "compressed.png", bytes.NewReader(out))
	return err
}

func compressSlashCommandHandler(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	opts := OptionMap(i.ApplicationCommandData().Options)
	md, err := fetchURL(OptString(opts, "url"))
	if err != nil {
		return err
	}
	if !IsImageExt(md.Ext) {
		return fmt.Errorf("unsupported image type: %s", OrEmpty(md.Ext))
	}

	out, err := Pixelate(md.Data, compressCells)
	if err != nil {
		return err
	}

	return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Files: []*discordgo.File{{Name: "compressed.png", Reader: bytes.NewReader(out)}},
		},
	})
}

