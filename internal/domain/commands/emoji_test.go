package commands

import (
	"strings"
	"testing"

	"github.com/bwmarrin/discordgo"
)

func TestParseEmoji(t *testing.T) {
	cases := []struct {
		name     string
		in       string
		wantName string
		wantID   string
		wantAnim bool
		wantOK   bool
	}{
		{"static full", "<:wave:1234>", "wave", "1234", false, true},
		{"animated full", "<a:wave:1234>", "wave", "1234", true, true},
		{"bare mention", ":wave:1234", "wave", "1234", false, true},
		{"plain id", "1234", "", "1234", false, true},
		{"whitespace padded", "  <:wave:1234>  ", "wave", "1234", false, true},
		{"missing id", "<:wave:>", "wave", "", false, false},
		{"bare word is an id candidate", "abc", "", "abc", false, false},
		{"mixed digits", "12a4", "", "12a4", false, false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			name, id, anim, ok := parseEmoji(tc.in)
			if name != tc.wantName || id != tc.wantID || anim != tc.wantAnim || ok != tc.wantOK {
				t.Errorf("parseEmoji(%q) = (%q,%q,%v,%v), want (%q,%q,%v,%v)",
					tc.in, name, id, anim, ok, tc.wantName, tc.wantID, tc.wantAnim, tc.wantOK)
			}
		})
	}
}

func TestEmojiTarget(t *testing.T) {
	reply := &discordgo.MessageCreate{Message: &discordgo.Message{
		ReferencedMessage: &discordgo.Message{Content: "look at this <a:spin:777>"},
	}}

	name, id, anim, ok := emojiTarget(reply, nil)
	if !ok || name != "spin" || id != "777" || !anim {
		t.Errorf("reply target = (%q,%q,%v,%v)", name, id, anim, ok)
	}

	// Explicit arg should win over the reply.
	name, id, anim, ok = emojiTarget(reply, []string{"<:star:1>"})
	if !ok || name != "star" || id != "1" || anim {
		t.Errorf("arg target = (%q,%q,%v,%v)", name, id, anim, ok)
	}

	// Unparseable arg -> not ok even with a reply present.
	if _, _, _, ok := emojiTarget(reply, []string{"junk"}); ok {
		t.Error("expected unparseable arg to fail")
	}

	// No args, no reply.
	if _, _, _, ok := emojiTarget(&discordgo.MessageCreate{Message: &discordgo.Message{}}, nil); ok {
		t.Error("expected no match without args or reply")
	}

	// Reply without an emoji.
	plain := &discordgo.MessageCreate{Message: &discordgo.Message{
		ReferencedMessage: &discordgo.Message{Content: "just text"},
	}}
	if _, _, _, ok := emojiTarget(plain, nil); ok {
		t.Error("expected no match for emoji-free reply")
	}
}

func TestDigitsOnly(t *testing.T) {
	for _, s := range []string{"", "a", "12a", " 1", "1.2", "-1"} {
		if digitsOnly(s) {
			t.Errorf("digitsOnly(%q) = true, want false", s)
		}
	}
	for _, s := range []string{"0", "1", "1234567890"} {
		if !digitsOnly(s) {
			t.Errorf("digitsOnly(%q) = false, want true", s)
		}
	}
}

func TestEmojiImageURL(t *testing.T) {
	if got := emojiImageURL("123", false, 256); got != "https://cdn.discordapp.com/emojis/123.png?size=256" {
		t.Errorf("static URL = %q", got)
	}
	if got := emojiImageURL("123", true, 128); got != "https://cdn.discordapp.com/emojis/123.gif?size=128" {
		t.Errorf("animated URL = %q", got)
	}
}

func TestBase64DataURI(t *testing.T) {
	cases := []struct {
		name string
		md   *mediaData
		want string
	}{
		{"jpg to jpeg", &mediaData{Data: []byte{0xff}, Ext: "jpg"}, "data:image/jpeg;base64,/w=="},
		{"png stays", &mediaData{Data: []byte{0xff}, Ext: "png"}, "data:image/png;base64,/w=="},
		{"gif stays", &mediaData{Data: []byte{0xff}, Ext: "gif"}, "data:image/gif;base64,/w=="},
		{"webm to png", &mediaData{Data: []byte{0xff}, Ext: "webm"}, "data:image/png;base64,/w=="},
		{"mp4 to png", &mediaData{Data: []byte{0xff}, Ext: "mp4"}, "data:image/png;base64,/w=="},
		{"bin to png", &mediaData{Data: []byte{0xff}, Ext: "bin"}, "data:image/png;base64,/w=="},
		{"empty ext to png", &mediaData{Data: []byte{0xff}, Ext: ""}, "data:image/png;base64,/w=="},
		{"no data", &mediaData{Data: nil, Ext: "png"}, "data:image/png;base64,"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := base64DataURI(tc.md); got != tc.want {
				t.Errorf("got %q, want %q", got, tc.want)
			}
		})
	}
}

func TestSanitizeEmojiName(t *testing.T) {
	cases := []struct {
		in   string
		want string
	}{
		{"WEIRD_Name!", "weird_name"},
		{"with spaces", "withspaces"},
		{"a.b-c", "abc"},
		{"!!!", ""},
		{"", ""},
		{strings.Repeat("x", 50), strings.Repeat("x", emojiNameMaxLen)},
	}
	for _, tc := range cases {
		if got := sanitizeEmojiName(tc.in); got != tc.want {
			t.Errorf("sanitizeEmojiName(%q) = %q, want %q", tc.in, got, tc.want)
		}
	}
}

func TestFilepathExt(t *testing.T) {
	cases := []struct{ in, want string }{
		{"a.png", ".png"},
		{"a.tar.gz", ".gz"},
		{"noext", ""},
		{".hidden", ".hidden"},
	}
	for _, tc := range cases {
		if got := filepathExt(tc.in); got != tc.want {
			t.Errorf("filepathExt(%q) = %q, want %q", tc.in, got, tc.want)
		}
	}
}

func TestAttachmentEmojiName(t *testing.T) {
	att := &discordgo.MessageAttachment{Filename: "happy_face.PNG"}
	if got := attachmentEmojiName(att); got != "happy_face" {
		t.Errorf("got %q", got)
	}
}

func TestChunkEmojiList(t *testing.T) {
	emojis := []*discordgo.Emoji{
		{Name: "a", ID: "1"},
		{Name: "b", ID: "2"},
		{Name: "c", ID: "3"},
	}
	one := chunkEmojiList(emojis, 10000)
	if len(one) != 1 {
		t.Fatalf("expected 1 field with big maxLen, got %d", len(one))
	}
	single := chunkEmojiList(emojis, 1)
	if len(single) != len(emojis) {
		t.Fatalf("expected 1 field per emoji with small maxLen, got %d", len(single))
	}
	for i, f := range single {
		if f.Name != "—" || !strings.Contains(f.Value, emojis[i].Name) {
			t.Errorf("field %d = %+v", i, f)
		}
	}
	if got := chunkEmojiList(nil, 1000); len(got) != 0 {
		t.Errorf("expected empty for nil input, got %d fields", len(got))
	}
}

func TestOptionMapAndOptString(t *testing.T) {
	opts := []*discordgo.ApplicationCommandInteractionDataOption{
		{Name: "text", Value: "hello"},
		{Name: "count", Value: 3},
	}
	m := optionMap(opts)
	if got := optString(m, "text"); got != "hello" {
		t.Errorf("optString text = %q", got)
	}
	if got := optString(m, "count"); got != "" {
		t.Errorf("optString count should be empty for non-string, got %q", got)
	}
	if got := optString(m, "missing"); got != "" {
		t.Errorf("optString missing = %q", got)
	}
	if optString(nil, "x") != "" {
		t.Error("optString on nil should be empty")
	}
}
