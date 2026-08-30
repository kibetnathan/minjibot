package commands

import (
	"strings"
	"testing"
)

func TestGiphyResultGIFURL(t *testing.T) {
	cases := []struct {
		name string
		r    giphyResult
		want string
	}{
		{
			"downsized preferred",
			giphyResult{Images: giphyImages{Downsized: giphyMedia{URL: "https://d"}, Original: giphyMedia{URL: "https://o"}}},
			"https://d",
		},
		{
			"original fallback",
			giphyResult{Images: giphyImages{Original: giphyMedia{URL: "https://o"}}},
			"https://o",
		},
		{
			"both empty",
			giphyResult{},
			"",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := tc.r.GIFURL(); got != tc.want {
				t.Errorf("GIFURL() = %q, want %q", got, tc.want)
			}
		})
	}
}

func TestGiphyResultDescriptionAndCreator(t *testing.T) {
	r := giphyResult{Title: "  Funny cat video  ", User: giphyUser{Username: "us", DisplayName: "Display"}}
	if d := r.Description(); d != "Funny cat video" {
		t.Errorf("Description = %q", d)
	}
	if c := r.CreatorName(); c != "Display" {
		t.Errorf("CreatorName = %q, want display name", c)
	}
	r.User.DisplayName = ""
	if c := r.CreatorName(); c != "us" {
		t.Errorf("CreatorName = %q, want username fallback", c)
	}
}

func TestFilterByCreator(t *testing.T) {
	results := []giphyResult{
		{ID: "1", User: giphyUser{Username: "bob_the_gif", DisplayName: "bob show"}, Title: "nice gif"},
		{ID: "2", User: giphyUser{Username: "zoe"}, Title: "Alice is cool"},
		{ID: "3", User: giphyUser{Username: "pete"}, Title: "plain"},
	}
	cases := []struct {
		name    string
		creator string
		wantIDs []string
	}{
		{"empty passthrough", "", []string{"1", "2", "3"}},
		{"by username", "bob_the_gif", []string{"1"}},
		{"by display name", "bob show", []string{"1"}},
		{"by title case-insensitive", "is cool", []string{"2"}},
		{"no match", "nobody", nil},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := filterByCreator(results, tc.creator)
			if len(got) != len(tc.wantIDs) {
				t.Fatalf("got %d results, want %d: %+v", len(got), len(tc.wantIDs), got)
			}
			for i, id := range tc.wantIDs {
				if got[i].ID != id {
					t.Errorf("result[%d].ID = %q, want %q", i, got[i].ID, id)
				}
			}
		})
	}
}

func TestParseGifSearchArgs(t *testing.T) {
	q, c := parseGifSearchArgs([]string{"funny", "cat", "creator:Giphy", "extra"})
	if q != "funny cat extra" || c != "Giphy" {
		t.Errorf("got query=%q creator=%q", q, c)
	}
	q, c = parseGifSearchArgs([]string{"by:nasa", "launch"})
	if q != "launch" || c != "nasa" {
		t.Errorf("got query=%q creator=%q", q, c)
	}
	q, c = parseGifSearchArgs([]string{"CREATOR:somethin", "a"})
	if q != "a" || c != "somethin" {
		t.Errorf("got query=%q creator=%q", q, c)
	}
	q, c = parseGifSearchArgs([]string{"just", "query"})
	if q != "just query" || c != "" {
		t.Errorf("got query=%q creator=%q", q, c)
	}
}

func TestBuildGifEmbed(t *testing.T) {
	r := giphyResult{
		Title:  "An amazing gif",
		Images: giphyImages{Original: giphyMedia{URL: "https://o.gif"}},
		User:   giphyUser{Username: "creator", DisplayName: "Creator"},
	}
	embed := buildGifEmbed("fallback", "", r)
	if embed.Title != "An amazing gif" {
		t.Errorf("title = %q", embed.Title)
	}
	if embed.Image == nil || embed.Image.URL != "https://o.gif" {
		t.Errorf("image = %+v", embed.Image)
	}
	if !strings.Contains(embed.Footer.Text, "https://giphy.com/@creator") {
		t.Errorf("footer = %q", embed.Footer.Text)
	}

	blank := giphyResult{}
	embed = buildGifEmbed("fallback title", "", blank)
	if embed.Title != "fallback title" {
		t.Errorf("fallback title = %q", embed.Title)
	}

	embed = buildGifEmbed("q", "someone", r)
	if !strings.Contains(embed.Footer.Text, `creator filter: "someone"`) {
		t.Errorf("footer should mention creator filter: %q", embed.Footer.Text)
	}
}

func TestBuildGifEmbedTruncatesDescription(t *testing.T) {
	r := giphyResult{Title: strings.Repeat("y", gifSearchMaxDesc*2)}
	embed := buildGifEmbed("q", "", r)
	if len(embed.Title) > gifSearchMaxDesc {
		t.Errorf("title too long: %d", len(embed.Title))
	}
	if !strings.HasSuffix(embed.Title, "...") {
		t.Errorf("title should be truncated: %q", embed.Title)
	}
}
