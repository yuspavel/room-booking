package models

import (
	"time"

	"github.com/google/uuid"
)

type Booking struct {
	ID                int           `db:"id"`
	RoomID            int           `db:"room_id"`
	UserID            uuid.UUID     `db:"user_id"`
	Title             string        `db:"title"`
	Description       string        `db:"description"`
	ParticipantsCount int           `db:"participants_count"`
	FromTime          time.Time     `db:"from_time"`
	ToTime            time.Time     `db:"to_time"`
	Status            BookingStatus `db:"status"`
	Paid              bool          `db:"paid"`
	CreatedAt         time.Time     `db:"created_at"`
	UpdatedAt         time.Time     `db:"updated_at"`
}

type BookingUpdate struct {
	RoomID            *int
	Title             *string
	Descriptions      *string
	ParticipantsCount *int
	FromDate          *time.Time
	ToDate            *time.Time
	Status            *string
	Paid              *bool
}

type BookingStatus string

const (
	BookingStatusPending   = "pending"
	BookingStatusCancelled = "cancelled"
	BookingStatusConfirmed = "confirmed"
	BookingStatusCompleted = "completed"
)
