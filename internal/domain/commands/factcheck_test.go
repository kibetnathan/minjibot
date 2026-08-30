package commands

import (
	"math"
	"strings"
	"testing"
)

func TestNormalizeClaim(t *testing.T) {
	if got := normalizeClaim("Hello,  World !! 123"); got != "hello world 123" {
		t.Errorf("got %q", got)
	}
	if got := normalizeClaim(""); got != "" {
		t.Errorf("got %q", got)
	}
}

func TestBigrams(t *testing.T) {
	cases := []struct {
		in   string
		want []string
	}{
		{"a b c", []string{"a b", "b c"}},
		{"single", []string{"single"}},
		{"", []string{""}},
	}
	for _, tc := range cases {
		got := bigrams(tc.in)
		if len(got) != len(tc.want) {
			t.Fatalf("bigrams(%q) length = %d, want %d", tc.in, len(got), len(tc.want))
		}
		for i := range got {
			if got[i] != tc.want[i] {
				t.Errorf("bigrams(%q)[%d] = %q, want %q", tc.in, i, got[i], tc.want[i])
			}
		}
	}
}

func TestClaimSimilarity(t *testing.T) {
	cases := []struct {
		name string
		a, b string
		want float64
	}{
		{"equal", "the earth is flat", "the earth is flat", 1},
		{"normalized equal", "Hello, World!", "hello world", 1},
		{"substring", "the earth is definitely flat", "flat", 1},
		{"similar bigrams", "the earth is flat", "the earth is round", 0}, // shared "the earth", see assert below
		{"empty", "", "something", 0},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := claimSimilarity(tc.a, tc.b)
			if tc.name == "similar bigrams" {
				// bigrams: the-earth, earth-is, is-flat | the-earth, earth-is, is-round
				// overlap (per entry in a) = 2 of 3 -> 2*2/(3+3) = 0.667
				if math.Abs(got-2.0/3.0) > 0.001 {
					t.Errorf("got %v, want ~0.667", got)
				}
				return
			}
			if math.Abs(got-tc.want) > 0.001 {
				t.Errorf("got %v, want %v", got, tc.want)
			}
		})
	}
}

func TestBestFactCheckClaim(t *testing.T) {
	mk := func(text, rating string) factCheckClaim {
		return factCheckClaim{
			Text: text,
			ClaimReview: []factCheckClaimReview{{TextualRating: rating, Publisher: struct {
				Name string `json:"name"`
			}{Name: "Pub"}}},
		}
	}

	claims := []factCheckClaim{
		mk("water boils at 100 celsius", "True"),
		mk("completely unrelated thing", "False"),
	}
	best := bestFactCheckClaim(claims, "water boils at 100 celsius")
	if best == nil {
		t.Fatal("expected a match")
	}
	if best.ClaimReview[0].TextualRating != "True" {
		t.Errorf("matched wrong claim: %+v", *best)
	}

	// No claim close enough -> nil.
	if best := bestFactCheckClaim(claims, "the moon is made of cheese and jam"); best != nil {
		t.Errorf("expected nil for unrelated claim, got %+v", *best)
	}

	// Claims without reviews are skipped.
	noReview := []factCheckClaim{{Text: "the earth is flat"}}
	if best := bestFactCheckClaim(noReview, "the earth is flat"); best != nil {
		t.Error("expected nil when all claims lack reviews")
	}
}

func TestFactCheckColorFor(t *testing.T) {
	cases := []struct {
		in   string
		want int
	}{
		{"True", factCheckColorTrue},
		{"Mostly True", factCheckColorMixed},
		{"Half true", factCheckColorMixed},
		{"Mixed", factCheckColorMixed},
		{"False", factCheckColorFalse},
		{"Pants on fire", factCheckColorFalse},
		{"No evidence", factCheckColorFalse},
		{"Unrated", factCheckColorUnknown},
		{"", factCheckColorUnknown},
	}
	for _, tc := range cases {
		if got := factCheckColorFor(tc.in); got != tc.want {
			t.Errorf("factCheckColorFor(%q) = %d, want %d", tc.in, got, tc.want)
		}
	}
}

func TestOrEmptyText(t *testing.T) {
	if orEmptyText("") != "—" {
		t.Error("orEmptyText(\"\") should return em dash")
	}
	if orEmptyText("x") != "x" {
		t.Error("orEmptyText should pass through")
	}
}

func TestBuildClaimMatchEmbed(t *testing.T) {
	c := &factCheckClaim{
		Text: "claim text",
		ClaimReview: []factCheckClaimReview{{
			Publisher: struct {
				Name string `json:"name"`
			}{Name: "Snopes"},
			URL:           "https://snopes.com/x",
			Title:         "Snopes checks",
			TextualRating: "True",
		}},
	}
	embed := buildClaimMatchEmbed("claim text", c)
	if embed.Color != factCheckColorTrue {
		t.Errorf("color = %d", embed.Color)
	}
	if !strings.Contains(embed.Description, "claim text") {
		t.Errorf("description = %q", embed.Description)
	}
	found := map[string]string{}
	for _, f := range embed.Fields {
		found[f.Name] = f.Value
	}
	if found["Publisher"] != "Snopes" || found["Verdict"] != "True" || found["Source"] != "https://snopes.com/x" {
		t.Errorf("fields = %+v", found)
	}
}

func TestBuildConsensusEmbed(t *testing.T) {
	sources := make([]factCheckSource, 6)
	for i := range sources {
		sources[i] = factCheckSource{Title: "Src", URL: "https://x"}
	}
	embed := buildConsensusEmbed("claim", "assessment", sources, "my footer")
	if !strings.Contains(embed.Description, "claim") || !strings.Contains(embed.Description, "assessment") {
		t.Errorf("description = %q", embed.Description)
	}
	if len(embed.Fields) != 1 || !strings.Contains(embed.Fields[0].Name, "(6)") {
		t.Errorf("sources field = %+v", embed.Fields)
	}
	if strings.Count(embed.Fields[0].Value, "https://x") != 5 {
		t.Errorf("expected at most 5 sources listed, got %q", embed.Fields[0].Value)
	}
	if embed.Footer == nil || embed.Footer.Text != "my footer" {
		t.Errorf("footer = %+v", embed.Footer)
	}
}

func TestBuildLLMSourcesText(t *testing.T) {
	text := buildLLMSourcesText("the claim", []factCheckSource{
		{Title: "A", URL: "https://a", Snippet: "snip"},
	})
	if !strings.HasPrefix(text, "Claim: the claim") {
		t.Errorf("got %q", text)
	}
	if !strings.Contains(text, "1. A — https://a\n   snip") {
		t.Errorf("got %q", text)
	}
	empty := buildLLMSourcesText("c", []factCheckSource{{URL: "https://u"}})
	if !strings.Contains(empty, "— https://u") {
		t.Errorf("expected em dash for blank title, got %q", empty)
	}
}

func BenchmarkClaimSimilarity(b *testing.B) {
	for i := 0; i < b.N; i++ {
		claimSimilarity("the earth is flat and round at once", "the earth is round")
	}
}
