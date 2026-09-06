package api

import "strings"

// authorizer decides whether an authenticated Discord user is allowed to view
// the dashboard's guild data. It holds an allowlist of Discord user IDs; a user
// is authorized only if their ID is in the list. An empty allowlist authorizes
// nobody (fail closed), so a misconfigured deployment never leaks guild data.
type authorizer struct {
	adminIDs map[string]struct{}
}

// newAuthorizer builds an authorizer from a list of Discord user IDs. Blank and
// whitespace-only entries are ignored so a trailing comma in the env var does
// not silently create an empty-string "admin".
func newAuthorizer(ids []string) *authorizer {
	m := make(map[string]struct{}, len(ids))
	for _, id := range ids {
		id = strings.TrimSpace(id)
		if id == "" {
			continue
		}
		m[id] = struct{}{}
	}
	return &authorizer{adminIDs: m}
}

// isAdmin reports whether the given Discord user ID is on the allowlist.
func (a *authorizer) isAdmin(userID string) bool {
	if userID == "" {
		return false
	}
	_, ok := a.adminIDs[userID]
	return ok
}

// empty reports whether the allowlist has no entries. Used at startup to warn
// that every dashboard data endpoint will reject all users until configured.
func (a *authorizer) empty() bool {
	return len(a.adminIDs) == 0
}
