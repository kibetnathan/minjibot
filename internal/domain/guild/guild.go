package guild

import "time"

type Guild struct {
	ID          string
	Name        string
	PremiumTier int32
	CreatedAt   time.Time
}