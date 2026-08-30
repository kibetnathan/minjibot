package commands

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/kibetnathan/minjibot/internal/config"
)

const (
	factCheckAPI          = "https://factchecktools.googleapis.com/v1alpha1/claims:search"
	factCheckMinScore     = 0.55
	FactCheckColorTrue    = 0x57F287
	FactCheckColorMixed   = 0xFEE75C
	FactCheckColorFalse   = 0xED4245
	FactCheckColorUnknown = 0x95A5A6

	geminiAPIBase      = "https://generativelanguage.googleapis.com/v1beta"
	geminiDefaultModel = "gemini-3.6-flash"

	factCheckSystemPrompt = "You are a concise fact-checking assistant. Start directly with the verdict (True, Mostly True, Mixed, Mostly False, False, or Unverifiable) followed by a 1-2 sentence consensus summary. Do not output preamble or bullet markers."
)

type FactCheckClaimReview struct {
	Publisher struct {
		Name string `json:"name"`
	} `json:"publisher"`
	URL           string `json:"url"`
	Title         string `json:"title"`
	ReviewDate    string `json:"review_date"`
	TextualRating string `json:"textual_rating"`
}

type FactCheckClaim struct {
	Text        string                 `json:"text"`
	Claimant    string                 `json:"claimant"`
	ClaimDate   string                 `json:"claim_date"`
	ClaimReview []FactCheckClaimReview `json:"claimReview"`
}

type factCheckResponse struct {
	Claims []FactCheckClaim `json:"claims"`
}

type FactCheckSource struct {
	Title   string
	Snippet string
	URL     string
}

func factcheckMessageCommandHandler(s *discordgo.Session, m *discordgo.MessageCreate, args []string, cfg *config.Config) error {
	claim := strings.TrimSpace(strings.Join(args, " "))
	if claim == "" {
		claim = referencedMessageContent(s, m)
	}
	if claim == "" {
		_, err := s.ChannelMessageSend(m.ChannelID, "Usage: `!factcheck <claim>` — or reply to a message to fact-check it")
		return err
	}

	embed := runFactCheck(claim, cfg)
	_, err := s.ChannelMessageSendEmbed(m.ChannelID, embed)
	return err
}

func factcheckSlashCommandHandler(s *discordgo.Session, i *discordgo.InteractionCreate, cfg *config.Config) error {
	opts := OptionMap(i.ApplicationCommandData().Options)
	claim := strings.TrimSpace(OptString(opts, "claim"))
	if claim == "" {
		return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{Content: "Usage: `/factcheck claim:<claim>`"},
		})
	}

	embed := runFactCheck(claim, cfg)
	return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{Embeds: []*discordgo.MessageEmbed{embed}},
	})
}

// referencedMessageContent returns the text of the message the issuing message
// replied to, fetching it via REST if the gateway didn't include it.
func referencedMessageContent(s *discordgo.Session, m *discordgo.MessageCreate) string {
	target := m.ReferencedMessage
	if target == nil && m.MessageReference != nil {
		if ref := m.MessageReference; ref.MessageID != "" {
			fetched, err := s.ChannelMessage(ref.ChannelID, ref.MessageID)
			if err == nil {
				target = fetched
			}
		}
	}
	if target == nil {
		return ""
	}
	return strings.TrimSpace(target.Content)
}

// runFactCheck performs the multi-step fact check and returns a formatted embed.
func runFactCheck(claim string, cfg *config.Config) *discordgo.MessageEmbed {
	footer := "Google FactCheck ClaimSearch"

	// Step 1: exact claim match via the Google FactCheck Claim Search API.
	if key := strings.TrimSpace(cfg.GoogleFactCheckKey); key != "" {
		if matched := queryClaimSearch(key, claim); matched != nil {
			return BuildClaimMatchEmbed(claim, matched)
		}
	} else {
		footer = "Google FactCheck ClaimSearch (set GOOGLE_FACTCHECK_API_KEY to enable)"
	}

	// Step 2: web search for supporting sources (broader query, no strict quotes).
	sources := factCheckWebSearch(claim)

	// Step 3: AI assessment — always runs as final fallback.
	// If sources exist, it summarizes them. If not, it uses its own knowledge.
	assessment := "No exact claim rating found. Here is what search sources say."
	if len(sources) == 0 {
		assessment = "No external sources were found. The AI is answering from its training knowledge."
	}
	if geminiKey := strings.TrimSpace(cfg.GeminiAPIKey); geminiKey != "" {
		if summary := llmAssessClaim(cfg, geminiKey, claim, sources); summary != "" {
			assessment = summary
			if len(sources) > 0 {
				footer = "AI assessment based on the sources above (treat as informal)"
			} else {
				footer = "AI assessment from model knowledge (no external sources found; treat as informal)"
			}
		}
	} else {
		if len(sources) > 0 {
			footer = "Web search consensus (no claim rating matched). Set GEMINI_API_KEY for AI summary."
		} else {
			footer = "No sources found. Set GEMINI_API_KEY for AI assessment."
		}
		assessment = fmt.Sprintf("No exact claim rating found. Here's what fact-check and news sources currently say. Review the sources yourself before drawing conclusions.")
	}

	return BuildConsensusEmbed(claim, assessment, sources, footer)
}

// queryClaimSearch hits the Google FactCheck API and returns the best
// matching claim review, or nil if none is similar enough.
func queryClaimSearch(apiKey, claim string) *FactCheckClaim {
	u, err := url.Parse(factCheckAPI)
	if err != nil {
		return nil
	}
	q := u.Query()
	q.Set("key", apiKey)
	q.Set("query", claim)
	q.Set("languageCode", "en")
	u.RawQuery = q.Encode()

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/126.0.0.0 Safari/537.36")

	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 4*1024))
		if len(body) > 0 {
			var apiErr struct {
				Error struct {
					Message string `json:"message"`
				} `json:"error"`
			}
			if json.Unmarshal(body, &apiErr) == nil && apiErr.Error.Message != "" {
				return &FactCheckClaim{Text: apiErr.Error.Message} // surfaced below as no-match
			}
		}
		return nil
	}

	var out factCheckResponse
	body, err := io.ReadAll(io.LimitReader(resp.Body, 2*1024*1024))
	if err != nil {
		return nil
	}
	if err := json.Unmarshal(body, &out); err != nil {
		return nil
	}

	return BestFactCheckClaim(out.Claims, claim)
}

func BestFactCheckClaim(claims []FactCheckClaim, claim string) *FactCheckClaim {
	var best *FactCheckClaim
	bestScore := 0.0
	for i := range claims {
		c := &claims[i]
		if len(c.ClaimReview) == 0 {
			continue
		}
		score := ClaimSimilarity(claim, c.Text)
		if score > bestScore {
			best, bestScore = c, score
		}
	}
	if best == nil || bestScore < factCheckMinScore {
		return nil
	}
	return best
}

// ClaimSimilarity is a word-bigram Dice score with exact-substring priority.
func ClaimSimilarity(a, b string) float64 {
	a, b = NormalizeClaim(a), NormalizeClaim(b)
	if a == "" || b == "" {
		return 0
	}
	if a == b || strings.Contains(a, b) || strings.Contains(b, a) {
		return 1
	}
	ga, gb := Bigrams(a), Bigrams(b)
	if len(ga) == 0 || len(gb) == 0 {
		return 0
	}
	overlap := 0
	for _, x := range ga {
		for _, y := range gb {
			if x == y {
				overlap++
				break
			}
		}
	}
	return 2.0 * float64(overlap) / float64(len(ga)+len(gb))
}

func NormalizeClaim(s string) string {
	s = strings.ToLower(s)
	var b strings.Builder
	for _, r := range s {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == ' ' {
			b.WriteRune(r)
		}
	}
	return strings.Join(strings.Fields(b.String()), " ")
}

func Bigrams(s string) []string {
	words := strings.Split(s, " ")
	if len(words) == 1 {
		return words
	}
	out := make([]string, 0, len(words)-1)
	for i := 0; i < len(words)-1; i++ {
		out = append(out, words[i]+" "+words[i+1])
	}
	return out
}

func BuildClaimMatchEmbed(claim string, c *FactCheckClaim) *discordgo.MessageEmbed {
	review := c.ClaimReview[0]

	verdict := strings.TrimSpace(review.TextualRating)
	if verdict == "" {
		verdict = "Unrated"
	}
	title := strings.TrimSpace(review.Title)
	if title == "" {
		title = "Claim review by " + strings.TrimSpace(review.Publisher.Name)
	}

	embed := &discordgo.MessageEmbed{
		Color:       FactCheckColorFor(verdict),
		Title:       "Fact check",
		Description: fmt.Sprintf("> %s", TruncateForEmbed(claim, 1024)),
		Fields: []*discordgo.MessageEmbedField{
			{Name: "Publisher", Value: OrEmptyText(strings.TrimSpace(review.Publisher.Name)), Inline: true},
			{Name: "Verdict", Value: TruncateForEmbed(verdict, 256), Inline: true},
		},
	}

	if date := strings.TrimSpace(review.ReviewDate); date != "" {
		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{Name: "Reviewed", Value: date, Inline: true})
	}
	if title != "" {
		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{Name: "Article", Value: TruncateForEmbed(title, 1024)})
	}
	if review.URL != "" {
		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{Name: "Source", Value: review.URL})
	}

	embed.Footer = &discordgo.MessageEmbedFooter{Text: "Rating from Google FactCheck ClaimSearch"}
	return embed
}

func BuildConsensusEmbed(claim, assessment string, sources []FactCheckSource, footer string) *discordgo.MessageEmbed {
	embed := &discordgo.MessageEmbed{
		Color:       FactCheckColorUnknown,
		Title:       "Fact check (search consensus)",
		Description: fmt.Sprintf("> %s\n\n%s", TruncateForEmbed(claim, 1024), TruncateForEmbed(assessment, 1800)),
		Footer:      &discordgo.MessageEmbedFooter{Text: footer},
	}

	limit := 5
	if len(sources) < limit {
		limit = len(sources)
	}
	var sourceText strings.Builder
	for _, src := range sources[:limit] {
		sourceText.WriteString("- ")
		if src.Title != "" {
			sourceText.WriteString(TruncateForEmbed(src.Title, 150))
			sourceText.WriteString(" — ")
		}
		sourceText.WriteString(src.URL)
		sourceText.WriteString("\n")
	}
	embed.Fields = []*discordgo.MessageEmbedField{
		{Name: fmt.Sprintf("Sources (%d)", len(sources)), Value: sourceText.String()},
	}
	return embed
}

// factCheckWebSearch runs a DuckDuckGo search for the claim and turns the
// relevant hits into a small source list. The HTML endpoint is preferred
// because it returns real web results; the Instant Answer API usually omits
// them for fact-check queries.
func factCheckWebSearch(claim string) []FactCheckSource {
	query := fmt.Sprintf("%s fact check", claim)

	if sources := newDDGClient().htmlSearch(query, 8); len(sources) > 0 {
		return sources
	}

	res, err := newDDGClient().Search(query)
	if err != nil {
		return nil
	}

	topics := res.FlattenedTopics()
	if len(topics) == 0 {
		topics = res.Results
	}

	sources := make([]FactCheckSource, 0, len(topics))
	for _, t := range topics {
		if t.FirstURL == "" {
			continue
		}
		title := t.Text
		if idx := strings.Index(title, " - "); idx != -1 {
			title = title[:idx]
		}
		sources = append(sources, FactCheckSource{Title: strings.TrimSpace(title), URL: t.FirstURL})
	}
	return sources
}

// llmAssessClaim asks Google Gemini (free tier) to summarise the consensus
// from the given sources. Returns "" on any failure so the caller falls back
// to its own assessment text. The model defaults to gemini-2.5-flash and can
// be overridden with GEMINI_MODEL.
func llmAssessClaim(cfg *config.Config, apiKey, claim string, sources []FactCheckSource) string {
	model := strings.TrimSpace(cfg.GeminiModel)
	if model == "" {
		model = geminiDefaultModel
	}

	payload := map[string]any{
		"system_instruction": map[string]any{
			"parts": []map[string]string{{"text": factCheckSystemPrompt}},
		},
		"contents": []map[string]any{
			{"role": "user", "parts": []map[string]string{{"text": BuildLLMSourcesText(claim, sources)}}},
		},
		"generationConfig": map[string]any{
			"maxOutputTokens": 800,
			"temperature":     0.2,
		},
	}
	raw, err := json.Marshal(payload)
	if err != nil {
		return ""
	}

	u := fmt.Sprintf("%s/models/%s:generateContent", geminiAPIBase, url.PathEscape(model))
	req, err := http.NewRequest(http.MethodPost, u, strings.NewReader(string(raw)))
	if err != nil {
		return ""
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-goog-api-key", apiKey)

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyErr, _ := io.ReadAll(io.LimitReader(resp.Body, 4*1024))
		fmt.Printf("[Gemini API Error] Status: %d, Response: %s\n", resp.StatusCode, string(bodyErr))
		return ""
	}
	bodyResp, err := io.ReadAll(io.LimitReader(resp.Body, 64*1024))
	if err != nil {
		return ""
	}

	var out struct {
		Candidates []struct {
			FinishReason string `json:"finishReason"`
			Content      struct {
				Parts []struct {
					Text string `json:"text"`
				} `json:"parts"`
			} `json:"content"`
		} `json:"candidates"`
	}
	if err := json.Unmarshal(bodyResp, &out); err != nil {
		return ""
	}
	if len(out.Candidates) == 0 {
		fmt.Printf("[Gemini API Error] No candidates returned\n")
		return ""
	}
	if out.Candidates[0].FinishReason != "STOP" {
		fmt.Printf("[Gemini API Error] Non-STOP finish reason: %s\n", out.Candidates[0].FinishReason)
		return ""
	}

	var b strings.Builder
	for _, part := range out.Candidates[0].Content.Parts {
		b.WriteString(part.Text)
	}
	return strings.TrimSpace(b.String())
}

func BuildLLMSourcesText(claim string, sources []FactCheckSource) string {
	var b strings.Builder
	fmt.Fprintf(&b, "Claim: %s\n\nSearch results:\n", claim)
	for i, src := range sources {
		fmt.Fprintf(&b, "%d. %s — %s\n", i+1, OrEmptyText(src.Title), src.URL)
		if src.Snippet != "" {
			fmt.Fprintf(&b, "   %s\n", src.Snippet)
		}
	}
	return b.String()
}

func OrEmptyText(s string) string {
	if s == "" {
		return "—"
	}
	return s
}

func FactCheckColorFor(verdict string) int {
	v := strings.ToLower(verdict)
	switch {
	case strings.Contains(v, "true") && !strings.Contains(v, "mostly") && !strings.Contains(v, "half"):
		return FactCheckColorTrue
	case strings.Contains(v, "mostly"), strings.Contains(v, "half"), strings.Contains(v, "mixed"):
		return FactCheckColorMixed
	case strings.Contains(v, "false"), strings.Contains(v, "pants"), strings.Contains(v, "no evidence"):
		return FactCheckColorFalse
	default:
		return FactCheckColorUnknown
	}
}
