package commands

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

const mediaTimeout = 30 * time.Second

// mediaData is a chunk of downloaded bytes with an inferred extension.
type mediaData struct {
	Data []byte
	Ext  string
}

// fetchURL downloads an image/video file from a URL.
func fetchURL(rawURL string) (*mediaData, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return nil, err
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		return nil, fmt.Errorf("invalid URL: %s", rawURL)
	}

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/126.0.0.0 Safari/537.36")

	client := &http.Client{Timeout: mediaTimeout}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("fetch failed: %s returned status %d", rawURL, resp.StatusCode)
	}

	data, err := io.ReadAll(io.LimitReader(resp.Body, 25*1024*1024))
	if err != nil {
		return nil, err
	}

	ext := strings.ToLower(strings.TrimPrefix(path.Ext(u.Path), "."))
	if ext == "" {
		ext = extensionFromContentType(resp.Header.Get("Content-Type"))
	}
	return &mediaData{Data: data, Ext: ext}, nil
}

func extensionFromContentType(ct string) string {
	ct = strings.ToLower(strings.Split(ct, ";")[0])
	switch {
	case strings.Contains(ct, "image/png"):
		return "png"
	case strings.Contains(ct, "image/jpeg"), strings.Contains(ct, "image/jpg"):
		return "jpg"
	case strings.Contains(ct, "image/gif"):
		return "gif"
	case strings.Contains(ct, "image/webp"):
		return "webp"
	case strings.Contains(ct, "video/mp4"):
		return "mp4"
	case strings.Contains(ct, "video/webm"):
		return "webm"
	case strings.Contains(ct, "image/avif"):
		return "avif"
	}
	return "bin"
}

// firstAttachment returns the first attachment of a message.
func firstAttachment(m *discordgo.Message) *discordgo.MessageAttachment {
	if m == nil {
		return nil
	}
	for _, a := range m.Attachments {
		if a.ID != "" && a.URL != "" {
			return a
		}
	}
	return nil
}

// resolveMedia returns the first media source available from: an explicit URL,
// the replied-to message's attachment, or the issuing message's attachment.
// If the gateway didn't include the referenced message, it is fetched via REST
// so replies to images still work.
func resolveMedia(s *discordgo.Session, m *discordgo.MessageCreate, explicitURL string) (*mediaData, error) {
	if explicitURL != "" {
		return fetchURL(explicitURL)
	}

	target := m.ReferencedMessage
	if target == nil && m.MessageReference != nil {
		if ref := m.MessageReference; ref.MessageID != "" {
			fetched, err := s.ChannelMessage(ref.ChannelID, ref.MessageID)
			if err == nil {
				target = fetched
			}
		}
	}

	att := firstAttachment(target)
	if att == nil {
		att = firstAttachment(m.Message)
	}
	if att == nil {
		return nil, fmt.Errorf("no image found — reply to a message with an attachment, attach an image, or pass a URL")
	}
	return fetchURL(att.URL)
}
