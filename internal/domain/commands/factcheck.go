package commands

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

const (
	factCheckAPI      = "https://factchecktools.googleapis.com/v1alpha1/claims:search"
	factCheckMinScore = 0.55
	factCheckColorTrue   = 0x57F287
	factCheckColorMixed  = 0xFEE75C
	factCheckColorFalse  = 0xED4245
	factCheckColorUnknown = 0x95A5A6
)

type factCheckClaimReview struct {
	Publisher      struct {
		Name string `json:"name"`
	} `json:"publisher"`
	URL           string `json:"url"`
	Title         string `json:"title"`
	ReviewDate    string `json:"review_date"`
	TextualRating string `json:"textual_rating"`
}

type factCheckClaim struct {
	Text        string                 `json:"text"`
	Claimant    string                 `json:"claimant"`
	ClaimDate   string                 `json:"claim_date"`
	ClaimReview []factCheckClaimReview `json:"claimReview"`
}

type factCheckResponse struct {
	Claims []factCheckClaim `json:"claims"`
}

type factCheckSource struct {
	Title   string
	Snippet string
	URL     string
}

func factcheckMessageCommandHandler(s *discordgo.Session, channelID string, args []string) error {
	claim := strings.TrimSpace(strings.Join(args, " "))
	if claim == "" {
		_, err := s.ChannelMessageSend(channelID, "Usage: `!factcheck <claim>`")
		return err
	}

	embed := runFactCheck(claim)
	_, err := s.ChannelMessageSendEmbed(channelID, embed)
	return err
}

func factcheckSlashCommandHandler(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	opts := optionMap(i.ApplicationCommandData().Options)
	claim := strings.TrimSpace(optString(opts, "claim"))
	if claim == "" {
		return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{Content: "Usage: `/factcheck claim:<claim>`"},
		})
	}

	embed := runFactCheck(claim)
	return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{Embeds: []*discordgo.MessageEmbed{embed}},
	})
}

// runFactCheck performs the multi-step fact check and returns a formatted embed.
func runFactCheck(claim string) *discordgo.MessageEmbed {
	footer := "Google FactCheck ClaimSearch"

	// Step 1: exact claim match via the Google FactCheck Claim Search API.
	if key := strings.TrimSpace(os.Getenv("GOOGLE_FACTCHECK_API_KEY")); key != "" {
		if matched := queryClaimSearch(key, claim); matched != nil {
			return buildClaimMatchEmbed(claim, matched)
		}
	} else {
		footer = "Google FactCheck ClaimSearch (set GOOGLE_FACTCHECK_API_KEY to enable)"
	}

	// Step 2: fallback to web search consensus.
	sources := factCheckWebSearch(claim)
	if len(sources) == 0 {
		embed := &discordgo.MessageEmbed{
			Color:       factCheckColorUnknown,
			Title:       "Fact check: unable to verify",
			Description: fmt.Sprintf("No claim rating or search consensus was found for:\n\n> %q", claim),
			Footer:      &discordgo.MessageEmbedFooter{Text: footer},
		}
		return embed
	}

	assessment := fmt.Sprintf("No exact claim rating found. Here's what fact-check and news sources currently say. Review the sources yourself before drawing conclusions.")
	if openaiKey := strings.TrimSpace(os.Getenv("OPENAI_API_KEY")); openaiKey != "" {
		if summary := llmAssessClaim(openaiKey, claim, sources); summary != "" {
			assessment = summary
			footer = "AI assessment based on the sources above (treat as informal)"
		}
	} else {
		footer = "Web search consensus (no claim rating matched). Set OPENAI_API_KEY for AI summary."
	}

	return buildConsensusEmbed(claim, assessment, sources, footer)
}

// queryClaimSearch hits the Google FactCheck API and returns the best
// matching claim review, or nil if none is similar enough.
func queryClaimSearch(apiKey, claim string) *factCheckClaim {
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
				return &factCheckClaim{Text: apiErr.Error.Message} // surfaced below as no-match
			}
		}
		return nil
	}

	var out factCheckResponse
	if err := json.Unmarshal(resp.Body, &out); err != nil {
		return nil
	}

	return bestFactCheckClaim(out.Claims, claim)
}

func bestFactCheckClaim(claims []factCheckClaim, claim string) *factCheckClaim {
	var best *factCheckClaim
	bestScore := 0.0
	for i := range claims {
		c := &claims[i]
		if len(c.ClaimReview) == 0 {
			continue
		}
		score := claimSimilarity(claim, c.Text)
		if score > bestScore {
			best, bestScore = c, score
		}
	}
	if best == nil || bestScore < factCheckMinScore {
		return nil
	}
	return best
}

// claimSimilarity is a word-bigram Dice score with exact-substring priority.
func claimSimilarity(a, b string) float64 {
	a, b = normalizeClaim(a), normalizeClaim(b)
	if a == "" || b == "" {
		return 0
	}
	if a == b || strings.Contains(a, b) || strings.Contains(b, a) {
		return 1
	}
	ga, gb := bigrams(a), bigrams(b)
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

func normalizeClaim(s string) string {
	s = strings.ToLower(s)
	var b strings.Builder
	for _, r := range s {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == ' ' {
			b.WriteRune(r)
		}
	}
	return strings.Join(strings.Fields(b.String()), " ")
}

func bigrams(s string) []string {
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

func buildClaimMatchEmbed(claim string, c *factCheckClaim) *discordgo.MessageEmbed {
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
		Color:       factCheckColorFor(verdict),
		Title:       "Fact check",
		Description: fmt.Sprintf("> %s", truncateForEmbed(claim, 1024)),
		Fields: []*discordgo.MessageEmbedField{
			{Name: "Publisher", Value: orEmptyText(strings.TrimSpace(review.Publisher.Name)), Inline: true},
			{Name: "Verdict", Value: truncateForEmbed(verdict, 256), Inline: true},
		},
	}

	if date := strings.TrimSpace(review.ReviewDate); date != "" {
		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{Name: "Reviewed", Value: date, Inline: true})
	}
	if title != "" {
		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{Name: "Article", Value: truncateForEmbed(title, 1024)})
	}
	if review.URL != "" {
		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{Name: "Source", Value: review.URL})
	}

	embed.Footer = &discordgo.MessageEmbedFooter{Text: "Rating from Google FactCheck ClaimSearch"}
	return embed
}

func buildConsensusEmbed(claim, assessment string, sources []factCheckSource, footer string) *discordgo.MessageEmbed {
	embed := &discordgo.MessageEmbed{
		Color:       factCheckColorUnknown,
		Title:       "Fact check (search consensus)",
		Description: fmt.Sprintf("> %s\n\n%s", truncateForEmbed(claim, 1024), truncateForEmbed(assessment, 1800)),
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
			sourceText.WriteString(truncateForEmbed(src.Title, 150))
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
// relevant hits into a small source list.
func factCheckWebSearch(claim string) []factCheckSource {
	query := fmt.Sprintf("%q fact check", claim)
	res, err := newDDGClient().Search(query)
	if err != nil {
		return nil
	}

	topics := res.FlattenedTopics()
	if len(topics) == 0 {
		topics = res.Results
	}

	sources := make([]factCheckSource, 0, len(topics))
	for _, t := range topics {
		if t.FirstURL == "" {
			continue
		}
		title := t.Text
		if idx := strings.Index(title, " - "); idx != -1 {
			title = title[:idx]
		}
		sources = append(sources, factCheckSource{Title: strings.TrimSpace(title), URL: t.FirstURL})
	}
	return sources
}

// llmAssessClaim asks an OpenAI-compatible chat model to summarise the
// consensus from the given sources. Returns "" on any failure.
func llmAssessClaim(apiKey, claim string, sources []factCheckSource) string {
	body := fmt.Sprintf(`{
  "model": "gpt-4o-mini",
  "max_tokens": 300,
  "temperature": 0.2,
  "messages": [
    {"role": "system", "content": "You are a careful fact-checking assistant. Based only on the web search results provided by the user, give: 1) a verdict — True, Mostly True, Mixed, Mostly False, False, or Unverifiable; and 2) a 1-2 sentence summary of the current consensus that names the most relevant source. If the results are inconclusive, say Unverifiable and explain why. Do not invent facts or sources."},
    {"role": "user", "content": %q}
  ]
}`, buildLLMSourcesText(claim, sources))

	data := []byte(strings.ReplaceAll(body, "%q", ""))
	_ = data

	payload := map[string]any{
		"model":       "gpt-4o-mini",
		"max_tokens":  300,
		"temperature": 0.2,
		"messages": []map[string]string{
			{"role": "system", "content": "You are a careful fact-checking assistant. Based only on the web search results provided by the user, give: 1) a verdict — True, Mostly True, Mixed, Mostly False, False, or Unverifiable; and 2) a 1-2 sentence summary of the current consensus that names the most relevant source. If the results are inconclusive, say Unverifiable and explain why. Do not invent facts or sources."},
			{"role": "user", "content": buildLLMSourcesText(claim, sources)},
		},
	}
	raw, _ := json.Marshal(payload)

	req, err := http.NewRequest(http.MethodPost, "https://api.openai.com/v1/chat/completions", strings.NewReader(string(raw)))
	if err != nil {
		return ""
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return ""
	}
	bodyResp, err := io.ReadAll(io.LimitReader(resp.Body, 64*1024))
	if err != nil {
		return ""
	}

	var out struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}
	if err := json.Unmarshal(bodyResp, &out); err != nil || len(out.Choices) == 0 {
		return ""
	}
	return strings.TrimSpace(out.Choices[0].Message.Content)
}

func buildLLMSourcesText(claim string, sources []factCheckSource) string {
	var b strings.Builder
	fmt.Fprintf(&b, "Claim: %s\n\nSearch results:\n", claim)
	for i, src := range sources {
		fmt.Fprintf(&b, "%d. %s — %s\n", i+1, orEmptyText(src.Title), src.URL)
		if src.Snippet != "" {
			fmt.Fprintf(&b, "   %s\n", src.Snippet)
		}
	}
	return b.String()
}

func orEmptyText(s string) string {
	if s == "" {
		return "—"
	}
	return s
}

func factCheckColorFor(verdict string) int {
	v := strings.ToLower(verdict)
	switch {
	case strings.Contains(v, "true") && !strings.Contains(v, "mostly"):
		return factCheckColorTrue
	case strings.Contains(v, "mostly"), strings.Contains(v, "half"), strings.Contains(v, "mixed"):
		return factCheckColorMixed
	case strings.Contains(v, "false"), strings.Contains(v, "pants"), strings.Contains(v, "no evidence"):
		return factCheckColorFalse
	default:
		return factCheckColorUnknown
	}
}