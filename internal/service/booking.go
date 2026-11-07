package service

import (
	"context"
	"errors"
	errs "room-booking/internal/errors"
	"room-booking/internal/models"
	"time"

	"github.com/google/uuid"
)

type BookingRepo interface {
	InsertBooking(ctx context.Context, roomID int, userID uuid.UUID, title string, description string, participants int, fromTime time.Time, toTime time.Time, status models.BookingStatus, paid bool) (*models.Booking, error)
	GetBooking(ctx context.Context, id int) (*models.Booking, error)
	UpdateBooking(ctx context.Context, id int, b models.BookingUpdate) (*models.Booking, error)
	ListBookings(ctx context.Context) ([]models.Booking, error)
	HasOverlapped(ctx context.Context, id int, roomID int, FromTime, toTime time.Time) (bool, error)
}

type BookingService struct {
	bookingRepo BookingRepo
	roomRepo    RoomRepo
}

func (b *BookingService) InsertBooking(ctx context.Context, bkg *models.Booking) (*models.Booking, error) {
	var (
		err        error
		room_exist bool
	)

	if err = validateTimeRange(bkg.FromTime, bkg.ToTime); err != nil {
		return nil, errs.ErrInvalidTimeRange
	}

	if err = validateMinMeetingTime(bkg.FromTime, bkg.ToTime); err != nil {
		return nil, errs.ErrMinTimeLack
	}

	if err = validateMaxMeetingTime(bkg.FromTime, bkg.ToTime); err != nil {
		return nil, errs.ErrMaxTimeExceed
	}

	if err = validateBackdate(bkg.FromTime); err != nil {
		return nil, errs.ErrInvalidBookingTime
	}

	if err = validateParticipantsNumber(bkg.ParticipantsCount); err != nil {
		return nil, errs.ErrNoParticipantsInvolved
	}

	room_exist, err = b.roomRepo.Exist(ctx, bkg.RoomID)
	if err != nil { //-------------------Ошибка БД, провверка не пройдена
		return nil, err
	} else if !room_exist { //------------------------Комната не существует, проверка не пройдена
		return nil, errs.ErrRoomNotFound
	} //-------------------------Если два условия выше не выполнены, проверка пройдена

	overlapped, err := b.bookingRepo.HasOverlapped(ctx, bkg.ID, bkg.RoomID, bkg.FromTime, bkg.ToTime)
	if err != nil {
		return nil, err
	}
	if overlapped {
		return nil, errs.ErrBookingOverlap
	}

	room, err := b.roomRepo.GetRoom(ctx, bkg.RoomID)
	if err != nil {
		if errors.Is(err, errs.ErrRoomNotFound) {
			return nil, errs.ErrRoomNotFound
		}
		return nil, err
	}

	if err = validateMeetingCapacity(room.Capacity, bkg.ParticipantsCount); err != nil {
		return nil, errs.ErrMaxPaticipantsExceeded
	}

	if err = validateAmount(bkg.Amount); err != nil {
		return nil, errs.ErrNegativeAmount
	}

	//Вычисляем статус в зависимости от поступившей оплаты. Если денег 0 или меньше полной оплаты, то pending. Если оплата полная, то confirmed
	status := validateInsertStatus(bkg.FromTime, bkg.ToTime, bkg.Amount, room.HourPrice)

	booking, err := b.bookingRepo.InsertBooking(ctx, bkg.RoomID, bkg.UserID, bkg.Title, bkg.Description, bkg.ParticipantsCount, bkg.FromTime, bkg.ToTime, status, bkg.Paid)
	if err != nil {
		return nil, err
	}
	return booking, nil
}

func (b *BookingService) UpdateBooking(ctx context.Context, bkg models.BookingUpdate) (*models.Booking, error) {
	var (
		err        error
		room_exist bool
	)
	booking, err := b.bookingRepo.GetBooking(ctx, bkg.ID) //Выбираем текущий booking
	if err != nil {
		if errors.Is(err, errs.ErrBookingNotFound) {
			return nil, errs.ErrBookingNotFound
		}
		return nil, err
	}

	booking.ApplyPatch(bkg) //Применяем patch (составляем финальную структуру для обновления данных)

	room_exist, err = b.roomRepo.Exist(ctx, booking.RoomID) //Проверяем, существует ли совещательная комната
	if err != nil {
		return nil, err
	} else if !room_exist {
		return nil, errs.ErrRoomNotFound
	}

	//Проверка на пересечение обновляемых диапазонов дат бронирования
	needOverLapCheck := booking.RoomID != bkg.RoomID || bkg.FromTime.Set || bkg.ToTime.Set //Если клиент прислал необходимые значения, то только в этом случае вызываем "дорогую" проверку с запросом к БД.
	if needOverLapCheck {
		overlapped, err := b.bookingRepo.HasOverlapped(ctx, bkg.ID, booking.RoomID, booking.FromTime, booking.ToTime)
		if err != nil {
			return nil, err
		}
		if overlapped {
			return nil, errs.ErrBookingOverlap
		}
	}

	room, err := b.roomRepo.GetRoom(ctx, booking.RoomID)
	if err != nil {
		if errors.Is(err, errs.ErrRoomNotFound) {
			return nil, errs.ErrRoomNotFound
		}
		return nil, err
	}

	if err = validateBooking(*booking, room); err != nil {
		return nil, err
	}

	if !validateStatusTransition(booking.Status, booking.Status) {
		return nil, errs.ErrStatusNotAllowed
	}

	//fBooking, err := b.bookingRepo.UpdateBooking(ctx, bkg.ID, completed)
	return nil, nil
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

func validateBooking(completed models.Booking, room models.Room) error {
	var err error

	if err = validateTitle(completed.Title); err != nil {
		return err
	}
	if err = validateDescription(completed.Description); err != nil {
		return err
	}

	if err = validateBackdate(completed.FromTime); err != nil {
		return err
	}

	if err = validateTimeRange(completed.FromTime, completed.ToTime); err != nil {
		return err
	}

	if err = validateMinMeetingTime(completed.FromTime, completed.ToTime); err != nil {
		return err
	}

	if err = validateMaxMeetingTime(completed.FromTime, completed.ToTime); err != nil {
		return err
	}

	if err = validateParticipantsNumber(completed.ParticipantsCount); err != nil {
		return err
	}

	if err = validateMeetingCapacity(room.Capacity, completed.ParticipantsCount); err != nil {
		return err
	}

	if err = validateAmount(completed.Amount); err != nil {
		return err
	}

	return nil
}
