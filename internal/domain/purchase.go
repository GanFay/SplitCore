package domain

import "time"

type Purchase struct {
	ID          int
	FundID      int
	PayerID     int64
	Amount      float64
	Description string
	CreatedAt   time.Time
}
