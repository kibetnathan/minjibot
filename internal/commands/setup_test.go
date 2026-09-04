package commands

import (
	"strings"
	"testing"
)

func TestParseToggle(t *testing.T) {
	on := []string{"on", "ON", "enable", "enabled", "true", "yes", "1", "  on  "}
	off := []string{"off", "OFF", "disable", "disabled", "false", "no", "0"}
	bad := []string{"", "maybe", "onn", "2", "toggle"}

	for _, v := range on {
		got, ok := parseToggle(v)
		if !ok || !got {
			t.Errorf("parseToggle(%q) = (%v, %v), want (true, true)", v, got, ok)
		}
	}
	for _, v := range off {
		got, ok := parseToggle(v)
		if !ok || got {
			t.Errorf("parseToggle(%q) = (%v, %v), want (false, true)", v, got, ok)
		}
	}
	for _, v := range bad {
		if _, ok := parseToggle(v); ok {
			t.Errorf("parseToggle(%q) should be unrecognised", v)
		}
	}
}

func TestSetupStatusText(t *testing.T) {
	// No channel, logging off.
	s := setupStatusText("", false)
	if !strings.Contains(s, "Not configured") || !strings.Contains(s, "Message-content logging: off") {
		t.Errorf("unexpected status text: %q", s)
	}
	// Channel set, logging on.
	s = setupStatusText("123", true)
	if !strings.Contains(s, "<#123>") || !strings.Contains(s, "Message-content logging: on") {
		t.Errorf("unexpected status text: %q", s)
	}
}

func TestMessageLoggingConfirm(t *testing.T) {
	if on := messageLoggingConfirm(true); !strings.Contains(on, "on") {
		t.Errorf("enabled confirm should say on: %q", on)
	}
	if off := messageLoggingConfirm(false); !strings.Contains(off, "off") {
		t.Errorf("disabled confirm should say off: %q", off)
	}
}
