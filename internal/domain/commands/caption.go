package commands

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

// memegenText converts user text into memegen.link's path-safe form
// (spaces -> dashes, underscores preserved via ~).
func memegenText(s string) string {
	s = strings.Join(strings.Fields(strings.TrimSpace(s)), " ")
	s = strings.ReplaceAll(s, " ", "-")
	s = strings.ReplaceAll(s, "_", "~")
	for strings.Contains(s, "--") {
		s = strings.ReplaceAll(s, "--", "-")
	}
	return strings.Trim(s, "-")
}

func parseCaptionArgs(args []string) (top, bottom, imageURL string) {
	var rest []string
	for _, arg := range args {
		switch {
		case strings.HasPrefix(strings.ToLower(arg), "top:"):
			top = strings.TrimSpace(arg[len("top:"):])
		case strings.HasPrefix(strings.ToLower(arg), "bottom:"):
			bottom = strings.TrimSpace(arg[len("bottom:"):])
		case strings.HasPrefix(strings.ToLower(arg), "url:"):
			imageURL = strings.TrimSpace(arg[len("url:"):])
		default:
			rest = append(rest, arg)
		}
	}

	// Fallback positional form: "top text / bottom text [url]"
	if top == "" && bottom == "" && imageURL == "" && len(rest) > 0 {
		joined := strings.Join(rest, " ")
		if parts := strings.SplitN(joined, "/", 2); len(parts) == 2 {
			split := strings.Fields(parts[1])
			if len(split) >= 2 {
				top = strings.TrimSpace(parts[0])
				bottom = strings.TrimSpace(strings.Join(split[:len(split)-1], " "))
				imageURL = split[len(split)-1]
			}
		}
	}
	return top, bottom, imageURL
}

func captionMessageCommandHandler(s *discordgo.Session, m *discordgo.MessageCreate, args []string) error {
	top, bottom, imageURL := parseCaptionArgs(args)
	if top == "" && bottom == "" {
		_, err := s.ChannelMessageSend(m.ChannelID, "Usage: `!caption top:<text> bottom:<text> url:<image-url>`")
		return err
	}

	md, err := resolveMedia(s, m, imageURL)
	if err != nil {
		return err
	}

	// Host the source image so memegen can fetch it as a background.
	background := ""
	if imageURL != "" {
		background = imageURL
	} else {
		background, err = uploadToCatbox(md)
		if err != nil {
			return fmt.Errorf("uploading image for caption: %w", err)
		}
	}

	memeURL := fmt.Sprintf(
		"https://api.memegen.link/templates/custom/%s/%s.png?background=%s&font=impact",
		url.PathEscape(memegenText(top)),
		url.PathEscape(memegenText(bottom)),
		url.QueryEscape(background),
	)

	_, err = s.ChannelMessageSendEmbed(m.ChannelID, &discordgo.MessageEmbed{
		Color: 0xFEE75C,
		Image: &discordgo.MessageEmbedImage{URL: memeURL},
	})
	return err
}

func captionSlashCommandHandler(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	opts := optionMap(i.ApplicationCommandData().Options)
	top := optString(opts, "top")
	bottom := optString(opts, "bottom")
	imageURL := optString(opts, "image_url")

	if top == "" && bottom == "" {
		return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{Content: "Usage: `/caption top:<text> bottom:<text> image_url:<image>`"},
		})
	}

	memeURL := fmt.Sprintf(
		"https://api.memegen.link/templates/custom/%s/%s.png?background=%s&font=impact",
		url.PathEscape(memegenText(top)),
		url.PathEscape(memegenText(bottom)),
		url.QueryEscape(imageURL),
	)

	return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{{Color: 0xFEE75C, Image: &discordgo.MessageEmbedImage{URL: memeURL}}},
		},
	})
}

// uploadToCatbox hosts an image anonymously and returns its public URL so
// memegen can fetch it as a background.
func uploadToCatbox(md *mediaData) (string, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	if err := writer.WriteField("reqtype", "fileupload"); err != nil {
		return "", err
	}
	part, err := writer.CreateFormFile("fileToUpload", "image."+md.Ext)
	if err != nil {
		return "", err
	}
	if _, err := part.Write(md.Data); err != nil {
		return "", err
	}
	if err := writer.Close(); err != nil {
		return "", err
	}

	req, err := http.NewRequest(http.MethodPost, "https://catbox.moe/user/api.php", body)
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/126.0.0.0 Safari/537.36")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("catbox returned status %d", resp.StatusCode)
	}
	href, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	link := strings.TrimSpace(string(href))
	if !strings.HasPrefix(link, "https://") {
		return "", fmt.Errorf("catbox upload failed: %s", link)
	}
	return link, nil
}
