package commands

import (
	"strings"
	"testing"

	"github.com/bwmarrin/discordgo"
)

func TestStickerUploadName(t *testing.T) {
	cases := []struct {
		ext      string
		wantName string
		wantErr  bool
	}{
		{"png", "sticker.png", false},
		{"PNG", "sticker.png", false},
		{"apng", "sticker.png", false},
		{"gif", "sticker.gif", false}, // regression: GIF must not upload as .png
		{"GIF", "sticker.gif", false},
		{"json", "", true}, // Lottie is rejected, not silently mislabeled
		{"webp", "", true}, // unsupported raster type
		{"", "", true},
	}
	for _, c := range cases {
		got, err := stickerUploadName(c.ext)
		if c.wantErr {
			if err == nil {
				t.Errorf("stickerUploadName(%q): expected error, got name %q", c.ext, got)
			}
			continue
		}
		if err != nil {
			t.Errorf("stickerUploadName(%q): unexpected error: %v", c.ext, err)
			continue
		}
		if got != c.wantName {
			t.Errorf("stickerUploadName(%q) = %q, want %q", c.ext, got, c.wantName)
		}
	}
}

func TestStickerImageURL(t *testing.T) {
	const id = "123456789"
	cases := []struct {
		name   string
		format discordgo.StickerFormat
		want   string
	}{
		{"png", discordgo.StickerFormatTypePNG, "https://cdn.discordapp.com/stickers/123456789.png"},
		{"apng", discordgo.StickerFormatTypeAPNG, "https://cdn.discordapp.com/stickers/123456789.png"},
		{"gif", discordgo.StickerFormatTypeGIF, "https://cdn.discordapp.com/stickers/123456789.gif"},
		{"lottie", discordgo.StickerFormatTypeLottie, "https://cdn.discordapp.com/stickers/123456789.json"},
	}
	for _, c := range cases {
		if got := StickerImageURL(id, c.format); got != c.want {
			t.Errorf("StickerImageURL(%s) = %q, want %q", c.name, got, c.want)
		}
	}
}

// A GIF sticker's URL must round-trip to a .gif upload filename, proving the
// end-to-end path (format → CDN URL → fetched extension → upload name) no
// longer collapses animated stickers into a rejected .png upload.
func TestStickerGIFRoundTripsToGifUpload(t *testing.T) {
	url := StickerImageURL("999", discordgo.StickerFormatTypeGIF)
	ext := url[strings.LastIndex(url, ".")+1:]
	name, err := stickerUploadName(ext)
	if err != nil {
		t.Fatalf("unexpected error for gif sticker: %v", err)
	}
	if name != "sticker.gif" {
		t.Fatalf("gif sticker uploaded as %q, want sticker.gif", name)
	}
}

func TestValidateStickerName(t *testing.T) {
	cases := []struct {
		name    string
		wantErr bool
	}{
		{"", true},           // empty
		{"a", true},          // too short (1)
		{"ab", false},        // min
		{"  ab  ", false},    // trimmed to 2
		{"a nice sticker", false},
		{strings.Repeat("x", 30), false}, // max
		{strings.Repeat("x", 31), true},  // too long
	}
	for _, c := range cases {
		err := validateStickerName(c.name)
		if c.wantErr && err == nil {
			t.Errorf("validateStickerName(%q): expected error", c.name)
		}
		if !c.wantErr && err != nil {
			t.Errorf("validateStickerName(%q): unexpected error: %v", c.name, err)
		}
	}
}
