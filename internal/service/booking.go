package service

import (
	"context"
	"errors"
	errs "room-booking/internal/errors"
	"room-booking/internal/models"
	"time"
)

type BookingRepo interface {
	InsertBooking(ctx context.Context, b models.Booking) (*models.Booking, error)
	GetBooking(ctx context.Context, id int) (*models.Booking, error)
	UpdateBooking(ctx context.Context, b models.Booking) (*models.Booking, error)
	ListBookings(ctx context.Context) ([]models.Booking, error)
	HasOverlapped(ctx context.Context, id int, roomID int, FromTime, toTime time.Time) (bool, error)
}

type BookingService struct {
	bookingRepo BookingRepo
	roomRepo    RoomRepo
}

func (b *BookingService) InsertBooking(ctx context.Context, bk models.Booking) (*models.Booking, error) {
	var err error
	if err = bk.ValidateBooking(); err != nil {
		return nil, err
	}

	overlapped, err := b.bookingRepo.HasOverlapped(ctx, bk.ID, bk.RoomID, bk.FromTime, bk.ToTime)
	if err != nil {
		return nil, err
	}
	if overlapped {
		return nil, errs.ErrBookingOverlap
	}

	room, err := b.roomRepo.GetRoom(ctx, bk.RoomID)
	if err != nil {
		if errors.Is(err, errs.ErrRoomNotFound) {
			return nil, errs.ErrRoomNotFound
		}
		return nil, err
	}

	if err = validateMeetingCapacity(room.Capacity, bk.ParticipantsCount); err != nil {
		return nil, errs.ErrMaxPaticipantsExceeded
	}

	//Вычисляем статус в зависимости от поступившей оплаты. Если денег 0 или меньше полной оплаты, то pending. Если оплата полная, то confirmed
	bk.UpdateStatus(room.HourPrice)

	return b.bookingRepo.InsertBooking(ctx, bk)
}

func (b *BookingService) UpdateBooking(ctx context.Context, bkg models.BookingUpdate) (*models.Booking, error) {
	var (
		err error
	)
	//Выбираем обновляемое бронирование их БД
	booking, err := b.bookingRepo.GetBooking(ctx, bkg.ID) //Выбираем текущий booking
	if err != nil {
		if errors.Is(err, errs.ErrBookingNotFound) {
			return nil, errs.ErrBookingNotFound
		}
		return nil, err
	}

	booking.ApplyPatch(bkg) //Применяем patch (составляем финальную структуру для обновления данных)

	if err = booking.ValidateBooking(); err != nil {
		return nil, err
	}

	//Проверка на пересечение обновляемых диапазонов дат бронирования. Условия: комната изменилась, даты from/to установлены
	needOverLapCheck := bkg.FromTime.Set || bkg.ToTime.Set //Если клиент прислал новые даты бронирования, то только в этом случае вызываем "дорогую" проверку с запросом к БД.
	if needOverLapCheck {
		overlapped, err := b.bookingRepo.HasOverlapped(ctx, bkg.ID, booking.RoomID, booking.FromTime, booking.ToTime)
		if err != nil {
			return nil, err
		}
		if overlapped {
			return nil, errs.ErrBookingOverlap
		}
	}

	//Находим комнату, она понадобится для вычисления вместимости количества участников в validateBooking функции5
	room, err := b.roomRepo.GetRoom(ctx, booking.RoomID)
	if err != nil {
		if errors.Is(err, errs.ErrRoomNotFound) {
			return nil, errs.ErrRoomNotFound
		}
		return nil, err
	}

	if err = validateMeetingCapacity(room.Capacity, booking.ParticipantsCount); err != nil {
		return nil, err
	}

	booking.UpdateStatus(room.HourPrice) //Обновляем статус брони

	return b.bookingRepo.UpdateBooking(ctx, *booking)
}

func (b *BookingService) GetBooking(ctx context.Context, id int) (*models.Booking, error) {
	bkg, err := b.bookingRepo.GetBooking(ctx, id)
	if err != nil {
		if errors.Is(err, errs.ErrBookingNotFound) {
			return nil, errs.ErrBookingNotFound
		}
		return nil, err
	}
	return bkg, nil
}

func (b *BookingService) ListBookings(ctx context.Context) ([]models.Booking, error) {
	bkg, err := b.bookingRepo.ListBookings(ctx)
	if err != nil {
		return nil, err
	}
	return bkg, nil
}

// Можно использовать при Create и Update брони
func validateMeetingCapacity(cap int, pc int) error {
	if pc > cap {
		return errs.ErrMaxPaticipantsExceeded
	}
	return nil
}
