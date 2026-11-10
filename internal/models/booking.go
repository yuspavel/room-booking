package models

import (
	"fmt"
	errs "room-booking/internal/errors"
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
	//Status            patch.Field[BookingStatus]
	Amount patch.Field[int64]
	Paid   patch.Field[bool]
}

type BookingStatus string

const (
	BookingStatusPending   = "pending"
	BookingStatusCancelled = "cancelled"
	BookingStatusConfirmed = "confirmed"
	BookingStatusCompleted = "completed"
)

func (b *Booking) ApplyPatch(bu BookingUpdate) {
	b.RoomID = bu.RoomID
	if bu.Title.HasValue() {
		b.Title = *bu.Title.Val
	}
	if bu.Description.HasValue() {
		b.Description = *bu.Description.Val
	}
	if bu.FromTime.HasValue() {
		b.FromTime = *bu.FromTime.Val
	}
	if bu.ToTime.HasValue() {
		b.ToTime = *bu.ToTime.Val
	}
	if bu.ParticipantsCount.HasValue() {
		b.ParticipantsCount = *bu.ParticipantsCount.Val
	}
	if bu.Amount.HasValue() {
		b.Amount = *bu.Amount.Val
	}
	if bu.Paid.HasValue() {
		b.Paid = *bu.Paid.Val
	}
}

// Можно использовать только при Insert брони
func (b *Booking) UpdateStatus(hourPrice int64) {
	//Проверяем оплачиваемую сумму для выставления первоначального статуса брони
	minutes := int64(time.Duration(b.ToTime.Sub(b.FromTime)).Minutes())
	total_price := (minutes * hourPrice) / 60

	if b.Amount < int64(total_price) {
		b.Status = BookingStatusPending
		return
	}
	b.Status = BookingStatusConfirmed
}

func (b *Booking) Completed() error {
	if b.Status != BookingStatusConfirmed {
		return errs.ErrInvalidStatusTransition
	}
	b.Status = BookingStatusCompleted
	return nil
}

func (b *Booking) Canceled() error {
	if b.Status != BookingStatusConfirmed && b.Status != BookingStatusPending {
		return errs.ErrInvalidStatusTransition
	}
	b.Status = BookingStatusCancelled
	return nil
}

func (b *Booking) ValidateBooking() error {
	var err error
	if err = b.validateTitle(); err != nil {
		return err
	}
	if err = b.validateDescription(); err != nil {
		return err
	}
	if err = b.validateTimeRange(); err != nil {
		return err
	}
	if err = b.validateMinMeetingTime(); err != nil {
		return err
	}
	if err = b.validateMaxMeetingTime(); err != nil {
		return err
	}
	if err = b.validateBackdate(); err != nil {
		return err
	}
	if err = b.validateParticipantsNumber(); err != nil {
		return err
	}
	if err = b.validateAmount(); err != nil {
		return err
	}
	return nil
}

func (b *Booking) validateTimeRange() error {
	//Проверка, что дата окончания брони больше его начала
	if !b.ToTime.After(b.FromTime) {
		return errs.ErrInvalidTimeRange
	}
	return nil
}

// Можно использовать при Create и Update брони
func (b *Booking) validateMinMeetingTime() error {
	//Проверка, что время встречи не может быть меньше 30 минут
	dur := b.ToTime.Sub(b.FromTime)
	if dur < 30*time.Minute {
		return errs.ErrMinTimeLack
	}
	return nil
}

// Можно использовать при Create и Update  брони
func (b *Booking) validateMaxMeetingTime() error {
	//Проверка на максимальную длительность забронированного совещания. Не более 8 часов (стандартный рабочий день)
	dur := b.ToTime.Sub(b.FromTime)
	if dur > 8*time.Hour {
		return errs.ErrMaxTimeExceed
	}
	return nil
}

// Можно использовать при Create и Update  брони
func (b *Booking) validateBackdate() error {
	//Проверка что на создание брони задним числом
	if b.FromTime.Before(time.Now().Add(1 * time.Minute)) {
		return errs.ErrInvalidBookingTime
	}
	return nil
}

// Можно использовать при Create и Update брони
func (b *Booking) validateParticipantsNumber() error {
	if b.ParticipantsCount <= 0 {
		return errs.ErrNoParticipantsInvolved
	}
	return nil
}

// Можно использовать при Create и Update брони
func (b *Booking) validateAmount() error {
	//Проверяем, что сумма оплаты не отрицательная
	if b.Amount < 0 {
		return errs.ErrNegativeAmount
	}
	return nil
}

func (b *Booking) validateTitle() error {
	if b.Title == "null" {
		return fmt.Errorf("title cannot be null")
	}
	return nil
}

func (b *Booking) validateDescription() error {
	if b.Description == "null" {
		return fmt.Errorf("description cannot be null")
	}
	return nil
}
