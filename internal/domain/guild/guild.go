// Package guild defines the domain model for Discord guilds (servers)
// tracked by MinjiBot.
package guild

import "time"

type Guild struct {
	ID          string
	Name        string
	PremiumTier int32
	CreatedAt   time.Time
}
