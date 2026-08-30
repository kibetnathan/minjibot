package auditlog

import "time"

type AuditLog struct {
	ID        int64
	GuildID   string
	Action    string
	ActorID   string
	TargetID  string
	Metadata  []byte
	CreatedAt time.Time
}
