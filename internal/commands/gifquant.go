package commands

import (
	"bytes"
	"image"
	"image/color"
	"image/draw"
	"image/gif"
	"sort"
)

// imageToGIF converts a still image into a quality-animated GIF using a
// median-cut palette built from the source colours, plus Floyd-Steinberg
// dithering to minimise banding. Two gentle frames make Discord treat the
// result as a GIF.
func imageToGIF(img image.Image) ([]byte, error) {
	paletted := quantizeImage(img)

	anim := &gif.GIF{
		Image: []*image.Paletted{paletted, paletted},
		Delay: []int{5, 5},
	}
	var buf bytes.Buffer
	if err := gif.EncodeAll(&buf, anim); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// quantizeImage builds a 256-colour paletted version of img via median cut,
// dithering with the standard library's Floyd-Steinberg drawer.
func quantizeImage(img image.Image) *image.Paletted {
	bounds := img.Bounds()
	pal := medianCutPalette(img, 255)
	paletted := image.NewPaletted(bounds, pal)
	draw.FloydSteinberg.Draw(paletted, bounds, img, image.Point{})
	return paletted
}

// medianCutPalette returns a palette of up to maxColours colours sampled from
// img using a simple median cut.
func medianCutPalette(img image.Image, maxColours int) color.Palette {
	seen := make(map[uint32]color.RGBA)
	for y := img.Bounds().Min.Y; y < img.Bounds().Max.Y; y++ {
		for x := img.Bounds().Min.X; x < img.Bounds().Max.X; x++ {
			r, g, b, a := img.At(x, y).RGBA()
			if a == 0 {
				continue
			}
			c := color.RGBA{R: uint8(r >> 8), G: uint8(g >> 8), B: uint8(b >> 8), A: 255}
			key := uint32(c.R)<<16 | uint32(c.G)<<8 | uint32(c.B)
			seen[key] = c
		}
	}

	cols := make([]color.RGBA, 0, len(seen))
	for _, c := range seen {
		cols = append(cols, c)
	}
	if len(cols) == 0 {
		return color.Palette{color.RGBA{}, color.Black, color.White}
	}
	if len(cols) <= maxColours {
		return toPalette(cols)
	}
	return medianCut(cols, maxColours)
}

func toPalette(cols []color.RGBA) color.Palette {
	p := make(color.Palette, 0, len(cols)+1)
	p = append(p, color.RGBA{})
	for _, c := range cols {
		p = append(p, c)
	}
	return p
}

type mcBox struct {
	colors []color.RGBA
}

// medianCut repeatedly splits the box with the largest colour range until
// maxColours buckets remain, then averages each bucket into one palette entry.
func medianCut(cols []color.RGBA, maxColours int) color.Palette {
	boxes := []*mcBox{{colors: cols}}
	for len(boxes) < maxColours {
		idx := widestBox(boxes)
		if idx < 0 {
			break
		}
		a, b := splitBox(boxes[idx])
		if a == nil || b == nil {
			break
		}
		boxes[idx] = a
		boxes = append(boxes, b)
	}

	pal := make(color.Palette, 0, len(boxes)+1)
	pal = append(pal, color.RGBA{})
	for _, bx := range boxes {
		pal = append(pal, averageColor(bx.colors))
	}
	return pal
}

func widestBox(boxes []*mcBox) int {
	worst := -1
	worstRange := -1
	for i, bx := range boxes {
		if len(bx.colors) < 2 {
			continue
		}
		if r := colorRange(bx.colors); r > worstRange {
			worst, worstRange = i, r
		}
	}
	return worst
}

func colorRange(cols []color.RGBA) int {
	dr, dg, db := colorRanges(cols)
	if dr >= dg && dr >= db {
		return dr
	}
	if dg >= db {
		return dg
	}
	return db
}

// splitBox splits a box along its widest channel at the channel median,
// returning two boxes (nil if it cannot be split further).
func splitBox(bx *mcBox) (*mcBox, *mcBox) {
	if len(bx.colors) < 2 {
		return nil, nil
	}

	dr, dg, db := colorRanges(bx.colors)
	channel := 0
	switch {
	case dr >= dg && dr >= db:
		channel = 0
	case dg >= db:
		channel = 1
	default:
		channel = 2
	}

	switch channel {
	case 0:
		sort.Slice(bx.colors, func(i, j int) bool { return bx.colors[i].R < bx.colors[j].R })
	case 1:
		sort.Slice(bx.colors, func(i, j int) bool { return bx.colors[i].G < bx.colors[j].G })
	default:
		sort.Slice(bx.colors, func(i, j int) bool { return bx.colors[i].B < bx.colors[j].B })
	}

	mid := len(bx.colors) / 2
	return &mcBox{colors: bx.colors[:mid]}, &mcBox{colors: bx.colors[mid:]}
}

func colorRanges(cols []color.RGBA) (dr, dg, db int) {
	var minR, minG, minB uint8 = 255, 255, 255
	var maxR, maxG, maxB uint8 = 0, 0, 0
	for _, c := range cols {
		if c.R < minR {
			minR = c.R
		}
		if c.G < minG {
			minG = c.G
		}
		if c.B < minB {
			minB = c.B
		}
		if c.R > maxR {
			maxR = c.R
		}
		if c.G > maxG {
			maxG = c.G
		}
		if c.B > maxB {
			maxB = c.B
		}
	}
	return int(maxR) - int(minR), int(maxG) - int(minG), int(maxB) - int(minB)
}

func averageColor(cols []color.RGBA) color.RGBA {
	var r, g, b uint64
	for _, c := range cols {
		r += uint64(c.R)
		g += uint64(c.G)
		b += uint64(c.B)
	}
	n := uint64(len(cols))
	if n == 0 {
		return color.RGBA{}
	}
	return color.RGBA{R: uint8(r / n), G: uint8(g / n), B: uint8(b / n), A: 255}
}
