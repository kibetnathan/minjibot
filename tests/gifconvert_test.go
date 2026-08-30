package tests

import (
	"bytes"
	"image"
	"image/color"
	"image/color/palette"
	"image/draw"
	"image/gif"
	"image/png"
	"strings"
	"testing"

	"github.com/kibetnathan/minjibot/internal/commands"
)

func TestIsVideoExt(t *testing.T) {
	for _, ext := range []string{"mp4", "webm", "mov", "mkv", "avi", "m4v"} {
		if !commands.IsVideoExt(ext) {
			t.Errorf("commands.IsVideoExt(%q) = false", ext)
		}
	}
	for _, ext := range []string{"png", "jpg", "gif", "", "MP4"} {
		if commands.IsVideoExt(ext) {
			t.Errorf("commands.IsVideoExt(%q) = true", ext)
		}
	}
}

func TestIsImageExt(t *testing.T) {
	for _, ext := range []string{"png", "jpg", "jpeg", "gif", "webp", "bmp"} {
		if !commands.IsImageExt(ext) {
			t.Errorf("commands.IsImageExt(%q) = false", ext)
		}
	}
	for _, ext := range []string{"mp4", "txt", ""} {
		if commands.IsImageExt(ext) {
			t.Errorf("commands.IsImageExt(%q) = true", ext)
		}
	}
}

func TestSlugify(t *testing.T) {
	cases := []struct {
		in   string
		want string
	}{
		{"Hello World!", "hello-world"},
		{"a.b c_d-e", "ab-c-d-e"},
		{"UPPER lower 123", "upper-lower-123"},
		{"", ""},
		{"---", ""},
		{" too many spaces ", "too-many-spaces"},
	}
	for _, tc := range cases {
		if got := commands.Slugify(tc.in); got != tc.want {
			t.Errorf("commands.Slugify(%q) = %q, want %q", tc.in, got, tc.want)
		}
	}
	if len(commands.Slugify(strings.Repeat("a", 100))) > 40 {
		t.Error("commands.Slugify should cap at 40 chars")
	}
}

func TestParseConvertArgs(t *testing.T) {
	start, dur, url := commands.ParseConvertArgs([]string{"start:1.5", "dur:3", "https://x/v.mp4"})
	if start != 1.5 || dur != 3 || url != "https://x/v.mp4" {
		t.Errorf("got (%v,%v,%q)", start, dur, url)
	}

	start, dur, url = commands.ParseConvertArgs([]string{"trim:2", "start:0.5"})
	if start != 0.5 || dur != 2 {
		t.Errorf("got start=%v dur=%v", start, dur)
	}

	start, dur, url = commands.ParseConvertArgs([]string{"http://x/a", "junk", "DUR:1"})
	if start != 0 || dur != 1 || url != "http://x/a" {
		t.Errorf("got (%v,%v,%q)", start, dur, url)
	}

	start, dur, url = commands.ParseConvertArgs(nil)
	if start != 0 || dur != 0 || url != "" {
		t.Errorf("got (%v,%v,%q)", start, dur, url)
	}
}

func TestImageBytesToGIFStaticPNG(t *testing.T) {
	pngBytes := makePNG(t)
	out, err := commands.ImageBytesToGIF(pngBytes, "png")
	if err != nil {
		t.Fatalf("commands.ImageBytesToGIF: %v", err)
	}
	decoded, err := gif.DecodeAll(bytes.NewReader(out))
	if err != nil {
		t.Fatalf("output is not a valid GIF: %v", err)
	}
	if len(decoded.Image) != 2 {
		t.Errorf("expected 2 frames from a static image, got %d", len(decoded.Image))
	}
}

func TestImageBytesToGIFAnimatedPassthrough(t *testing.T) {
	var buf bytes.Buffer
	anim := &gif.GIF{
		Image: []*image.Paletted{palettedFrame(255, 0, 0), palettedFrame(0, 0, 255)},
		Delay: []int{5, 5},
	}
	if err := gif.EncodeAll(&buf, anim); err != nil {
		t.Fatalf("encoding fixture: %v", err)
	}
	out, err := commands.ImageBytesToGIF(buf.Bytes(), "gif")
	if err != nil {
		t.Fatalf("commands.ImageBytesToGIF: %v", err)
	}
	decoded, err := gif.DecodeAll(bytes.NewReader(out))
	if err != nil {
		t.Fatalf("output GIF unreadable: %v", err)
	}
	if len(decoded.Image) < 2 {
		t.Errorf("expected animated GIF to keep multiple frames, got %d", len(decoded.Image))
	}
}

func TestImageBytesToGIFBadInput(t *testing.T) {
	if _, err := commands.ImageBytesToGIF([]byte("not an image"), "png"); err == nil {
		t.Error("expected an error for undecodable bytes")
	}
}

func TestOrEmpty(t *testing.T) {
	if commands.OrEmpty("") != "unknown" {
		t.Error("commands.OrEmpty(\"\") should return unknown")
	}
	if commands.OrEmpty("png") != "png" {
		t.Error("commands.OrEmpty should return the string")
	}
}

func makePNG(t *testing.T) []byte {
	t.Helper()
	img := image.NewRGBA(image.Rect(0, 0, 2, 2))
	img.Set(0, 0, color.RGBA{R: 255, A: 255})
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		t.Fatalf("png.Encode: %v", err)
	}
	return buf.Bytes()
}

func palettedFrame(r, g, b uint8) *image.Paletted {
	src := image.NewRGBA(image.Rect(0, 0, 2, 2))
	for x := 0; x < 2; x++ {
		for y := 0; y < 2; y++ {
			src.Set(x, y, color.RGBA{R: r, G: g, B: b, A: 255})
		}
	}
	dst := image.NewPaletted(src.Bounds(), palette.Plan9[:256])
	draw.Draw(dst, dst.Bounds(), src, image.Point{}, draw.Src)
	return dst
}
