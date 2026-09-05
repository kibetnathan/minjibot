package api

import "testing"

func TestNewAuthorizer_IgnoresBlankEntries(t *testing.T) {
	a := newAuthorizer([]string{"111", "  ", "", "222 "})
	if !a.isAdmin("111") {
		t.Errorf("expected 111 to be admin")
	}
	// Trailing whitespace is trimmed on both sides.
	if !a.isAdmin("222") {
		t.Errorf("expected 222 (whitespace-trimmed) to be admin")
	}
	// The blank/whitespace entries must not create an empty-string admin.
	if a.isAdmin("") {
		t.Errorf("empty user id must never be admin")
	}
	if len(a.adminIDs) != 2 {
		t.Errorf("expected 2 admins, got %d", len(a.adminIDs))
	}
}

func TestAuthorizer_IsAdmin(t *testing.T) {
	a := newAuthorizer([]string{"admin1", "admin2"})
	if !a.isAdmin("admin1") {
		t.Errorf("admin1 should be admin")
	}
	if a.isAdmin("nobody") {
		t.Errorf("nobody should not be admin")
	}
}

func TestAuthorizer_EmptyAllowlistFailsClosed(t *testing.T) {
	a := newAuthorizer(nil)
	if !a.empty() {
		t.Errorf("nil allowlist should report empty")
	}
	if a.isAdmin("anyone") {
		t.Errorf("empty allowlist must authorize nobody")
	}

	// An allowlist of only blank entries is also effectively empty.
	blank := newAuthorizer([]string{"", "   "})
	if !blank.empty() {
		t.Errorf("all-blank allowlist should report empty")
	}
}
