package models

import (
	"time"
)

type Room struct {
	ID        int       `db:"id" json:"id"`
	Name      string    `db:"name" json:"name"`
	Capacity  int       `db:"capacity" json:"capacity"`
	HourPrice int64     `db:"hour_price" json:"price"`
	Currency  string    `db:"currency" json:"currency"`
	CreatedAt time.Time `db:"created_at" json:"createdat"`
}

type RoomUpdate struct {
	Name      *string
	Capacity  *int
	HourPrice *float32
	Currency  *string
}
