package userpermission

import "time"

type UserPermission struct {
	ID              int64
	UserID          string
	GuildID         string
	Role            string
	PermissionsJSON []byte
	CreatedAt       time.Time
}