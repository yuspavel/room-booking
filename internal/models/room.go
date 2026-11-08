package models

import (
	errs "room-booking/internal/errors"
	"room-booking/internal/patch"
	"time"
)

type Room struct {
	ID        int       `db:"id"`
	Name      string    `db:"name"`
	Capacity  int       `db:"capacity"`
	HourPrice int64     `db:"hour_price"`
	Currency  string    `db:"currency"`
	CreatedAt time.Time `db:"created_at"`
}

type RoomUpdate struct {
	ID        int
	Name      patch.Field[string]
	Capacity  patch.Field[int]
	HourPrice patch.Field[float32]
	Currency  patch.Field[string]
}

type Currency string

const (
	RUB = "RUB"
	KZT = "KZT"
	USD = "USD"
	EUR = "EUR"
)

func (r *Room) ApplyPatch(ru RoomUpdate) {
	if ru.Name.Set && ru.Name.Val != nil {
		r.Name = *ru.Name.Val
	}
	if ru.Capacity.Set && ru.Capacity.Val != nil {
		r.Capacity = *ru.Capacity.Val
	}
	if ru.HourPrice.Set && ru.Capacity.Val != nil {
		r.Capacity = *ru.Capacity.Val
	}
	if ru.Currency.Set && ru.Currency.Val != nil {
		r.Currency = *ru.Currency.Val
	}
}

func (r *Room) ValidateRoom() error {
	var err error
	if err = r.validateRoomName(); err != nil {
		return err
	}
	if err = r.validateCapacity(); err != nil {
		return err
	}
	if err = r.validatePrice(); err != nil {
		return err
	}
	if err = r.validateCurency(); err != nil {
		return err
	}
	return nil
}
func (r *Room) validateRoomName() error {
	if r.Name == "" || r.Name == "null" {
		return errs.ErrInvalidRoomName
	}
	return nil
}

func (r *Room) validateCapacity() error {
	if r.Capacity <= 0 {
		return errs.ErrInvalidRoomCapacity
	}
	return nil
}

func (r *Room) validatePrice() error {
	if r.HourPrice <= 0 {
		return errs.ErrInvalidRoomPrice
	}
	return nil
}

func (r *Room) validateCurency() error {
	isCurrency := r.Currency == RUB || r.Currency == USD || r.Currency == KZT || r.Currency == EUR
	if !isCurrency {
		return errs.ErrCurrencyNotAllowed
	}
	return nil
}
