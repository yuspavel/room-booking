package service

import (
	"context"
	errs "room-booking/internal/errors"
	"room-booking/internal/models"
	"time"
)

/*Валидацию статуса я перенес в БД и сделал на основе ENUM перечисления и CONSTRAINT поля status*/

// Можно использовать при Create и Update  брони
func validateTimeRange(from, to time.Time) error {
	//Проверка на что дата окончания брони больше его начала
	if !to.After(from) {
		return errs.ErrInvalidTimeRange
	}
	return nil
}

// Можно использовать при Create и Update брони
func validateMinMeetingTime(from, to time.Time) error {
	dur := to.Sub(from)
	if dur < 30*time.Minute {
		return errs.ErrMinTimeLack
	}
	return nil
}

// Можно использовать при Create и Update  брони
func validateMaxMeetingTime(from, to time.Time) error {
	//Проверка на максимальную длительность забронированного совещания. Не более 8 часов (стандартный рабочий день)
	dur := to.Sub(from)
	if dur > 8*time.Hour {
		return errs.ErrMaxTimeExceed
	}
	return nil
}

// Можно использовать при Create и Update  брони
func validateBackdate(from time.Time) error {
	//Проверка что на создание брони задним числом
	if from.Before(time.Now().Add(1 * time.Minute)) {
		return errs.ErrInvalidBookingTime
	}
	return nil
}

// Можно использовать при Create и Update брони
func (s *BookingService) validateIsOverlapped(ctx context.Context, roomID int, from, to time.Time) error {
	//Проверка на пересечение дат бронирования одной и той же комнаты
	overlapped, err := s.bookingRepo.HasOverlapped(ctx, roomID, from, to)

	if err != nil { //---------------------Ошибка БД, проверка не пройдена
		return err
	}
	if overlapped { //----------------------Overlap найден, проверка не пройдена
		return errs.ErrBookingOverlap
	}
	return nil //--------------------------Если два условия выше не выполнены, проверка пройдена
}

// Можно использовать при Create и Update брони
func validateParticipantsNumber(pc int) error {
	if pc <= 0 {
		return errs.ErrNoParticipantsInvolved
	}
	return nil
}

// Можно использовать при Create и Update брони
func validateMeetingCapacity(cap int, pc int) error {
	if pc > cap {
		return errs.ErrMaxPaticipantsExceeded
	}
	return nil
}

// Можно использовать только при Update брони
func validateStatusTransition(bookingStatus models.BookingStatus, requestStatus models.BookingStatus) bool {
	switch bookingStatus {
	case models.BookingStatusPending:
		return requestStatus == models.BookingStatusConfirmed || requestStatus == models.BookingStatusCancelled || requestStatus == models.BookingStatusPending
	case models.BookingStatusConfirmed:
		return requestStatus == models.BookingStatusCancelled || requestStatus == models.BookingStatusCompleted || requestStatus == models.BookingStatusConfirmed
	case models.BookingStatusCancelled, models.BookingStatusCompleted:
		return false
	default:
		return false
	}
}

// Можно использовать при Create и Update брони
func validateAmount(amount int64) error {
	//Проверяем, что сумма оплаты не отрицательная
	if amount < 0 {
		return errs.ErrNegativeAmount
	}
	return nil
}

// Можно использовать только при Insert брони
func validateInsertStatus(from, to time.Time, amount int64, hourPrice int64) models.BookingStatus {
	//Проверяем оплачиваемую сумму для выставления первоначального статуса брони
	minutes := int64(time.Duration(to.Sub(from)).Minutes())
	total_price := (minutes * hourPrice) / 60

	if amount < int64(total_price) {
		return models.BookingStatusPending
	}
	return models.BookingStatusConfirmed
}
