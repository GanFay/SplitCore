package domain

import "time"

type User struct {
	TgID      int64
	Username  string
	FirstName string
	CreatedAt time.Time
}
