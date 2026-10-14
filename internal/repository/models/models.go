package models

import (
	"time"

	"github.com/google/uuid"
)

type Room struct {
	ID        int       `db:"id"`
	Name      string    `db:"name"`
	Capacity  int       `db:"capacity"`
	HourPrice float64   `db:"hour_price"`
	Currency  string    `db:"currency"`
	CreatedAt time.Time `db:"created_at"`
}

type User struct {
	ID        uuid.UUID `db:"id"`
	Name      string    `db:"name"`
	PassHash  string    `db:"pass_hash"`
	Email     string    `db:"email"`
	Phone     string    `db:"phone"`
	CreatedAt time.Time `db:"created_at"`
}

type Booking struct {
	ID                int       `db:"id"`
	RoomID            int       `db:"room_id"`
	UserID            uuid.UUID `db:"user_id"`
	Title             string    `db:"title"`
	Description       string    `db:"description"`
	ParticipantsCount int       `db:"participants_count"`
	FromTime          time.Time `db:"from_time"`
	ToTime            time.Time `db:"to_time"`
	status            string    `db:"status"`
	paid              bool      `db:"paid"`
	CreatedAt         time.Time `db:"created_at"`
	UpdatedAt         time.Time `db:"updated_at"`
}

type UpdateRoom struct {
	Name      *string
	Capacity  *int
	HourPrice *float32
	Currency  *string
}
