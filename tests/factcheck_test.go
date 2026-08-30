package tests

import (
	"math"
	"strings"
	"testing"

	"github.com/kibetnathan/minjibot/internal/commands"
)

func TestNormalizeClaim(t *testing.T) {
	if got := commands.NormalizeClaim("Hello,  World !! 123"); got != "hello world 123" {
		t.Errorf("got %q", got)
	}
	if got := commands.NormalizeClaim(""); got != "" {
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
		got := commands.Bigrams(tc.in)
		if len(got) != len(tc.want) {
			t.Fatalf("commands.Bigrams(%q) length = %d, want %d", tc.in, len(got), len(tc.want))
		}
		for i := range got {
			if got[i] != tc.want[i] {
				t.Errorf("commands.Bigrams(%q)[%d] = %q, want %q", tc.in, i, got[i], tc.want[i])
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
		{"similar commands.Bigrams", "the earth is flat", "the earth is round", 0}, // shared "the earth", see assert below
		{"empty", "", "something", 0},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := commands.ClaimSimilarity(tc.a, tc.b)
			if tc.name == "similar commands.Bigrams" {
				// commands.Bigrams: the-earth, earth-is, is-flat | the-earth, earth-is, is-round
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
	mk := func(text, rating string) commands.FactCheckClaim {
		return commands.FactCheckClaim{
			Text: text,
			ClaimReview: []commands.FactCheckClaimReview{{TextualRating: rating, Publisher: struct {
				Name string `json:"name"`
			}{Name: "Pub"}}},
		}
	}

	claims := []commands.FactCheckClaim{
		mk("water boils at 100 celsius", "True"),
		mk("completely unrelated thing", "False"),
	}
	best := commands.BestFactCheckClaim(claims, "water boils at 100 celsius")
	if best == nil {
		t.Fatal("expected a match")
	}
	if best.ClaimReview[0].TextualRating != "True" {
		t.Errorf("matched wrong claim: %+v", *best)
	}

	// No claim close enough -> nil.
	if best := commands.BestFactCheckClaim(claims, "the moon is made of cheese and jam"); best != nil {
		t.Errorf("expected nil for unrelated claim, got %+v", *best)
	}

	// Claims without reviews are skipped.
	noReview := []commands.FactCheckClaim{{Text: "the earth is flat"}}
	if best := commands.BestFactCheckClaim(noReview, "the earth is flat"); best != nil {
		t.Error("expected nil when all claims lack reviews")
	}
}

func TestFactCheckColorFor(t *testing.T) {
	cases := []struct {
		in   string
		want int
	}{
		{"True", commands.FactCheckColorTrue},
		{"Mostly True", commands.FactCheckColorMixed},
		{"Half true", commands.FactCheckColorMixed},
		{"Mixed", commands.FactCheckColorMixed},
		{"False", commands.FactCheckColorFalse},
		{"Pants on fire", commands.FactCheckColorFalse},
		{"No evidence", commands.FactCheckColorFalse},
		{"Unrated", commands.FactCheckColorUnknown},
		{"", commands.FactCheckColorUnknown},
	}
	for _, tc := range cases {
		if got := commands.FactCheckColorFor(tc.in); got != tc.want {
			t.Errorf("commands.FactCheckColorFor(%q) = %d, want %d", tc.in, got, tc.want)
		}
	}
}

func TestOrEmptyText(t *testing.T) {
	if commands.OrEmptyText("") != "—" {
		t.Error("commands.OrEmptyText(\"\") should return em dash")
	}
	if commands.OrEmptyText("x") != "x" {
		t.Error("commands.OrEmptyText should pass through")
	}
}

func TestBuildClaimMatchEmbed(t *testing.T) {
	c := &commands.FactCheckClaim{
		Text: "claim text",
		ClaimReview: []commands.FactCheckClaimReview{{
			Publisher: struct {
				Name string `json:"name"`
			}{Name: "Snopes"},
			URL:           "https://snopes.com/x",
			Title:         "Snopes checks",
			TextualRating: "True",
		}},
	}
	embed := commands.BuildClaimMatchEmbed("claim text", c)
	if embed.Color != commands.FactCheckColorTrue {
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
	sources := make([]commands.FactCheckSource, 6)
	for i := range sources {
		sources[i] = commands.FactCheckSource{Title: "Src", URL: "https://x"}
	}
	embed := commands.BuildConsensusEmbed("claim", "assessment", sources, "my footer")
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
	text := commands.BuildLLMSourcesText("the claim", []commands.FactCheckSource{
		{Title: "A", URL: "https://a", Snippet: "snip"},
	})
	if !strings.HasPrefix(text, "Claim: the claim") {
		t.Errorf("got %q", text)
	}
	if !strings.Contains(text, "1. A — https://a\n   snip") {
		t.Errorf("got %q", text)
	}
	empty := commands.BuildLLMSourcesText("c", []commands.FactCheckSource{{URL: "https://u"}})
	if !strings.Contains(empty, "— https://u") {
		t.Errorf("expected em dash for blank title, got %q", empty)
	}
}

func BenchmarkClaimSimilarity(b *testing.B) {
	for i := 0; i < b.N; i++ {
		commands.ClaimSimilarity("the earth is flat and round at once", "the earth is round")
	}
}
