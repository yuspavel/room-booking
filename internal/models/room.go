package models

import (
	"time"
)

type Room struct {
	ID        int       `db:"id"`
	Name      string    `db:"name"`
	Capacity  int       `db:"capacity"`
	HourPrice float64   `db:"hour_price"`
	Currency  string    `db:"currency"`
	CreatedAt time.Time `db:"created_at"`
}

type RoomUpdate struct {
	Name      *string
	Capacity  *int
	HourPrice *float32
	Currency  *string
}
