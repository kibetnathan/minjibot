// Package auditlog defines the domain model for moderation action audit
// logs, recording who did what to whom in a guild.
package auditlog

import "time"

type AuditLog struct {
	ID         int64
	GuildID    string
	Action     string
	ActorID    string
	ActorName  string
	TargetID   string
	TargetName string
	Metadata   []byte
	CreatedAt  time.Time
}
