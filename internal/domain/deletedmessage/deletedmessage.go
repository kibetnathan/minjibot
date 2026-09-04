// Package deletedmessage provides the domain model for messages that were
// deleted in a guild, captured for auditing and review on the dashboard.
package deletedmessage

import "time"

// DeletedMessage is a single message that was removed from a guild, together
// with what we knew about its author and content at the time of deletion.
type DeletedMessage struct {
	ID            int64
	GuildID       string
	ChannelID     string
	MessageID     string
	AuthorID      string
	AuthorName    string
	Content       string
	Attachments   []byte
	DeletedBy     string
	DeletedByName string
	CreatedAt     time.Time
}
