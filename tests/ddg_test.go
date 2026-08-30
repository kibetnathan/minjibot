package tests

import (
	"testing"

	"github.com/kibetnathan/minjibot/internal/commands"
)

func TestParseDDGHTMLResults(t *testing.T) {
	body := `
<div class="result results_links results_links_deep web_result">
	<h2 class="result__title">
		<a rel="nofollow" class="result__a" href="//duckduckgo.com/l/?uddg=https%3A%2F%2Fexample.com%2Ffactcheck%3Fa%3D1%26b%3D2&amp;rut=abc">Fact check: example claim is <b>true</b></a>
	</h2>
	<a class="result__snippet">The claim checks out, says <b>Example News</b>. It&#x27;s correct.</a>
</div>
<div class="result">
	<h2 class="result__title">
		<a rel="nofollow" class="result__a" href="https://example.org/article">Second result title</a>
	</h2>
</div>
`

	got := commands.ParseDDGHTMLResults(body, 8)
	if len(got) != 2 {
		t.Fatalf("expected 2 results, got %d: %+v", len(got), got)
	}

	if got[0].URL != "https://example.com/factcheck?a=1&b=2" {
		t.Errorf("url = %q", got[0].URL)
	}
	if got[0].Title != "Fact check: example claim is true" {
		t.Errorf("title = %q", got[0].Title)
	}
	if got[0].Snippet != "The claim checks out, says Example News. It's correct." {
		t.Errorf("snippet = %q", got[0].Snippet)
	}

	if got[1].URL != "https://example.org/article" {
		t.Errorf("second url = %q", got[1].URL)
	}

	if trimmed := commands.ParseDDGHTMLResults(body, 1); len(trimmed) != 1 {
		t.Errorf("limit not applied, got %d", len(trimmed))
	}
}
