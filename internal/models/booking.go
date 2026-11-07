package models

import (
	"room-booking/internal/patch"
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
	Amount            int64         `db:"amount"`
	Paid              bool          `db:"paid"`
	CreatedAt         time.Time     `db:"created_at"`
	UpdatedAt         time.Time     `db:"updated_at"`
}

type BookingUpdate struct {
	ID                int
	RoomID            int
	Title             patch.Field[string]
	Description       patch.Field[string]
	ParticipantsCount patch.Field[int]
	FromTime          patch.Field[time.Time]
	ToTime            patch.Field[time.Time]
	Status            patch.Field[BookingStatus]
	Amount            patch.Field[int64]
	Paid              patch.Field[bool]
}

type BookingStatus string

const (
	BookingStatusPending   = "pending"
	BookingStatusCancelled = "cancelled"
	BookingStatusConfirmed = "confirmed"
	BookingStatusCompleted = "completed"
)

func (b *Booking) ApplyPatch(bu BookingUpdate) {
	if bu.Title.Set && bu.Title.Val != nil {
		b.Title = *bu.Title.Val
	}
	if bu.Description.Set && bu.Description.Val != nil {
		b.Description = *bu.Description.Val
	}
	if bu.FromTime.Set && bu.FromTime.Val != nil {
		b.FromTime = *bu.FromTime.Val
	}
	if bu.ToTime.Set && bu.ToTime.Val != nil {
		b.ToTime = *bu.ToTime.Val
	}
	if bu.ParticipantsCount.Set && bu.ParticipantsCount.Val != nil {
		b.ParticipantsCount = *bu.ParticipantsCount.Val
	}
	if bu.Status.Set && bu.Status.Val != nil {
		b.Status = *bu.Status.Val
	}
	if bu.Amount.Set && bu.Amount.Val != nil {
		b.Amount = *bu.Amount.Val
	}
	if bu.Paid.Set && bu.Paid.Val != nil {
		b.Paid = *bu.Paid.Val
	}
}
