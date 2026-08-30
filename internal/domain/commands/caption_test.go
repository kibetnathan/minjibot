package commands

import (
	"testing"
)

func TestMemegenText(t *testing.T) {
	cases := []struct {
		in   string
		want string
	}{
		{"hello world", "hello-world"},
		{"  hello  world  ", "hello-world"},
		{"under_score", "under~score"},
		{"under_score x", "under~score-x"},
		{"a - b", "a-b"},
		{" -hello ", "hello"},
		{"", ""},
	}
	for _, tc := range cases {
		if got := memegenText(tc.in); got != tc.want {
			t.Errorf("memegenText(%q) = %q, want %q", tc.in, got, tc.want)
		}
	}
}

func TestParseCaptionArgs(t *testing.T) {
	cases := []struct {
		name       string
		args       []string
		wantTop    string
		wantBottom string
		wantImage  string
	}{
		{"flags", []string{"top:Hello", "bottom:World", "url:https://x/y.png"}, "Hello", "World", "https://x/y.png"},
		{"uppercase flags", []string{"TOP:Hi", "BOTTOM:Lo", "URL:https://z/q.png"}, "Hi", "Lo", "https://z/q.png"},
		{"top only", []string{"top:Hi"}, "Hi", "", ""},
		{"positional with url", []string{"hello world", "/", "bottom text", "https://x/y.png"}, "hello world", "bottom text", "https://x/y.png"},
		{"mixed flags beat positional", []string{"top:x", "anything", "else"}, "x", "", ""},
		{"empty", nil, "", "", ""},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			top, bottom, image := parseCaptionArgs(tc.args)
			if top != tc.wantTop || bottom != tc.wantBottom || image != tc.wantImage {
				t.Errorf("got (%q,%q,%q), want (%q,%q,%q)", top, bottom, image, tc.wantTop, tc.wantBottom, tc.wantImage)
			}
		})
	}
}
