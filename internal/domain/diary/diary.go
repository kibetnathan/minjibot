package diary

import "time"

// Entry is a single private journal entry owned by a user (diary add/view/
// delete). Each row belongs to one user and is only visible to them.
type Entry struct {
	ID        int64
	UserID    string
	Content   string
	CreatedAt time.Time
}
