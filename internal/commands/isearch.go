package commands

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

const iSearchMaxResults = 4

// yandexMatch is one visually-similar result returned by Yandex's reverse
// image search. PageURL is the page the image appears on; OriginalImage is a
// direct link to the image itself.
type yandexMatch struct {
	Title   string `json:"title"`
	PageURL string `json:"url"`
	Domain  string `json:"domain"`
	Thumb   struct {
		URL string `json:"url"`
	} `json:"thumb"`
	OriginalImage struct {
		URL string `json:"url"`
	} `json:"originalImage"`
}

// yandexReverseSearch asks Yandex to find pages containing (or visually
// similar to) imageURL, which must be publicly reachable. Results are parsed
// out of the server-rendered results HTML; no API key is required.
func yandexReverseSearch(imageURL string) ([]yandexMatch, error) {
	u := url.URL{Scheme: "https", Host: "yandex.com", Path: "/images/search"}
	q := u.Query()
	q.Set("rpt", "imageview")
	q.Set("url", imageURL)
	u.RawQuery = q.Encode()

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/126.0.0.0 Safari/537.36")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")

	client := &http.Client{Timeout: 20 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("reverse image search returned status %d", resp.StatusCode)
	}
	body, err := io.ReadAll(io.LimitReader(resp.Body, 6*1024*1024))
	if err != nil {
		return nil, err
	}
	if bytes.Contains(bytes.ToLower(body), []byte("captcha")) {
		return nil, fmt.Errorf("temporarily blocked by the image service, try again in a bit")
	}

	decoded := html.UnescapeString(string(body))
	var matches []yandexMatch
	seen := map[string]bool{}
	for _, raw := range extractJSONObjects(decoded) {
		var m yandexMatch
		if err := json.Unmarshal([]byte(raw), &m); err != nil {
			continue
		}
		if m.PageURL == "" || m.OriginalImage.URL == "" {
			continue
		}
		if seen[m.OriginalImage.URL] {
			continue
		}
		seen[m.OriginalImage.URL] = true
		matches = append(matches, m)
		if len(matches) >= iSearchMaxResults*2 {
			break
		}
	}
	return matches, nil
}

// extractJSONObjects pulls brace-balanced `{"title":...}` JSON objects out of
// the unescaped page, tracking string state so braces inside quoted values are
// ignored. It's a targeted parser for Yandex's server-rendered result cards.
func extractJSONObjects(s string) []string {
	var objs []string
	i := 0
	for i < len(s) {
		p := strings.Index(s[i:], `{"title":"`)
		if p == -1 {
			break
		}
		p += i

		depth := 0
		inString := false
		escaped := false
		end := -1
		for j := p; j < len(s); j++ {
			c := s[j]
			if inString {
				switch {
				case escaped:
					escaped = false
				case c == '\\':
					escaped = true
				case c == '"':
					inString = false
				}
				continue
			}
			switch c {
			case '"':
				inString = true
			case '{':
				depth++
			case '}':
				depth--
				if depth == 0 {
					end = j
				}
			}
			if end != -1 {
				break
			}
		}
		if end == -1 {
			break
		}
		objs = append(objs, s[p:end+1])
		i = end + 1
	}
	return objs
}

// stripYandexTrackers removes the utm_* query params Yandex appends to result
// page URLs so the links we show are clean and shareable.
func stripYandexTrackers(raw string) string {
	u, err := url.Parse(raw)
	if err != nil {
		return raw
	}
	q := u.Query()
	if q.Get("utm_medium") == "" && q.Get("utm_source") == "" {
		return raw
	}
	q.Del("utm_medium")
	q.Del("utm_source")
	u.RawQuery = q.Encode()
	return u.String()
}

// isearchMessage resolves the image to search from an explicit URL arg first,
// then the replied-to message's attachment, then the message's own attachment.
func isearchMessage(s *discordgo.Session, m *discordgo.MessageCreate, args []string) (string, error) {
	for _, a := range args {
		a = strings.Trim(a, "<>")
		if u, err := url.Parse(a); err == nil && (u.Scheme == "http" || u.Scheme == "https") {
			return a, nil
		}
	}
	if att := FirstAttachment(m.ReferencedMessage); att != nil {
		return att.URL, nil
	}
	if att := FirstAttachment(m.Message); att != nil {
		return att.URL, nil
	}
	return "", fmt.Errorf("usage: `-isearch <image-url>` or reply to / attach an image")
}

func isearchMessageCommandHandler(s *discordgo.Session, m *discordgo.MessageCreate, args []string) error {
	imageURL, err := isearchMessage(s, m, args)
	if err != nil {
		_, serr := s.ChannelMessageSend(m.ChannelID, err.Error())
		return serr
	}

	results, err := yandexReverseSearch(imageURL)
	if err != nil {
		return err
	}
	if len(results) == 0 {
		_, err := s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("No matches found for that image: %s", imageURL))
		return err
	}

	_, err = s.ChannelMessageSendComplex(m.ChannelID, buildISearchSend(imageURL, results))
	return err
}

func isearchSlashCommandHandler(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	opts := OptionMap(i.ApplicationCommandData().Options)
	imageURL := strings.TrimSpace(OptString(opts, "image_url"))
	if imageURL == "" {
		return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{Content: "Usage: `/isearch image_url:<url>`"},
		})
	}

	results, err := yandexReverseSearch(imageURL)
	if err != nil {
		return err
	}
	if len(results) == 0 {
		return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{Content: fmt.Sprintf("No matches found for that image: %s", imageURL)},
		})
	}

	return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: buildISearchContent(imageURL, results),
			Files:   downloadISearchImages(results),
		},
	})
}

// buildISearchSend assembles the prefix-command reply: matched images as
// attachments plus numbered source links. Fails softly if none of the matched
// images can be downloaded (still posts the sources).
func buildISearchSend(imageURL string, results []yandexMatch) *discordgo.MessageSend {
	return &discordgo.MessageSend{
		Content: buildISearchContent(imageURL, results),
		Files:   downloadISearchImages(results),
	}
}

func buildISearchContent(imageURL string, results []yandexMatch) string {
	var b strings.Builder
	fmt.Fprintf(&b, "Reverse image matches for **%s**\n", imageURL)
	for i, r := range results {
		if i >= iSearchMaxResults {
			break
		}
		fmt.Fprintf(&b, "%d. %s\n", i+1, stripYandexTrackers(r.PageURL))
	}
	return strings.TrimSpace(b.String())
}

// downloadISearchImages fetches up to iSearchMaxResults matched images, trying
// the high-res original first and falling back to Yandex's hosted thumbnail.
func downloadISearchImages(results []yandexMatch) []*discordgo.File {
	var files []*discordgo.File
	for _, r := range results {
		if len(files) >= iSearchMaxResults {
			break
		}
		md, err := fetchURL(r.OriginalImage.URL)
		if err != nil {
			thumb := r.Thumb.URL
			if !strings.HasPrefix(thumb, "http") {
				thumb = "https:" + thumb
			}
			md, err = fetchURL(thumb)
		}
		if err != nil || len(md.Data) == 0 {
			continue
		}
		ext := "png"
		switch md.Ext {
		case "jpg", "jpeg", "gif", "webp":
			ext = md.Ext
		case "bin":
			ext = "png"
		}
		files = append(files, &discordgo.File{
			Name:        fmt.Sprintf("isearch_%d.%s", len(files), ext),
			ContentType: "image/" + ext,
			Reader:      bytes.NewReader(md.Data),
		})
	}
	return files
}
