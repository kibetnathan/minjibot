package commands

import (
	"testing"
	"time"
)

func TestParseDuration(t *testing.T) {
	cases := []struct {
		in   string
		want time.Duration
	}{
		{"30s", 30 * time.Second},
		{"5m", 5 * time.Minute},
		{"2h", 2 * time.Hour},
		{"1d", 24 * time.Hour},
		{"1h30m", 90 * time.Minute},
		{"2d12h", 60 * time.Hour},
		{"0s", 0},
		{"", 0},
		{"10", 0},
		{"abc", 0},
		{"1x", 0},
		{"-5m", 0},
		{"1h 30m", 0},
	}
	for _, tc := range cases {
		if got := parseDuration(tc.in); got != tc.want {
			t.Errorf("parseDuration(%q) = %v, want %v", tc.in, got, tc.want)
		}
	}
}

func TestIsClock(t *testing.T) {
	for _, s := range []string{"14:00", "09:30", "23:59", "00:00"} {
		if !isClock(s) {
			t.Errorf("isClock(%q) = false, want true", s)
		}
	}
	for _, s := range []string{"9:00", "14:0", "1a:00", "12:3a", "12345", "", "14:00:00"} {
		if isClock(s) {
			t.Errorf("isClock(%q) = true, want false", s)
		}
	}
}

func TestHumanizeDuration(t *testing.T) {
	cases := []struct {
		in   time.Duration
		want string
	}{
		{45 * time.Second, "45s"},
		{59 * time.Second, "59s"},
		{90 * time.Minute, "1h30m"},
		{2 * time.Hour, "2h"},
		{120 * time.Minute, "2h"},
		{36 * time.Hour, "1d12h"},
		{48 * time.Hour, "2d"},
		{25 * time.Hour, "1d1h"},
	}
	for _, tc := range cases {
		if got := humanizeDuration(tc.in); got != tc.want {
			t.Errorf("humanizeDuration(%v) = %q, want %q", tc.in, got, tc.want)
		}
	}
}

func TestParseReminderArgs(t *testing.T) {
	cases := []struct {
		name     string
		args     []string
		wantText string
		wantDur  time.Duration
	}{
		{"in form", []string{"in", "30m", "buy", "milk"}, "buy milk", 30 * time.Minute},
		{"bare duration", []string{"2h", "take", "a", "break"}, "take a break", 2 * time.Hour},
		{"compact", []string{"1h30m", "code"}, "code", 90 * time.Minute},
		{"from now", []string{"1h", "from", "now", "stretch"}, "from now stretch", time.Hour},
		{"no args", nil, "", 0},
		{"no duration", []string{"buy", "milk"}, "", 0},
		{"zero duration", []string{"0m", "x"}, "", 0},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			text, d := parseReminderArgs(tc.args)
			if text != tc.wantText || d != tc.wantDur {
				t.Errorf("got (%q,%v), want (%q,%v)", text, d, tc.wantText, tc.wantDur)
			}
		})
	}
}

func TestParseReminderArgsClock(t *testing.T) {
	// Absolute time resolves relative to the local clock, so only assert that
	// the text is split out and the delay is positive.
	text, d := parseReminderArgs([]string{"in", "14:00", "lunch"})
	if text != "lunch" {
		t.Errorf("text = %q, want %q", text, "lunch")
	}
	if d <= 0 {
		t.Errorf("expected a positive delay for an absolute clock time, got %v", d)
	}
}
