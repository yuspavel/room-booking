package service

import (
	"context"
	"errors"
	errs "room-booking/internal/errors"
	"room-booking/internal/models"
	"time"

	"github.com/google/uuid"
)

type finalbook struct {
	Title              string
	Description        string
	From_time, To_Time time.Time
	ParticipantsCount  int
	Status             models.BookingStatus
	Amount             int64
	Paid               bool
}

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

	completed := applyPatch(bkg, *booking) //Применяем patch (составляем финальную структуру для обновления данных)

	room_exist, err = b.roomRepo.Exist(ctx, bkg.RoomID) //Проверяем, существует ли совещательная комната
	if err != nil {
		return nil, err
	} else if !room_exist {
		return nil, errs.ErrRoomNotFound
	}

	//Проверка на пересечение обновляемых диапазонов дат бронирования
	needOverLapCheck := booking.RoomID != bkg.RoomID || bkg.FromTime.Set || bkg.ToTime.Set //Если клиент прислал необходимые значения, то только в этом случае вызываем "дорогую" проверку с запросом к БД.
	if needOverLapCheck {
		overlapped, err := b.bookingRepo.HasOverlapped(ctx, bkg.ID, bkg.RoomID, completed.From_time, completed.To_Time)
		if err != nil {
			return nil, err
		}
		if overlapped {
			return nil, errs.ErrBookingOverlap
		}
	}

	room, err := b.roomRepo.GetRoom(ctx, bkg.RoomID)
	if err != nil {
		if errors.Is(err, errs.ErrRoomNotFound) {
			return nil, errs.ErrRoomNotFound
		}
		return nil, err
	}

	for _, v := range buildValidators(bkg, room) {
		if err := v(); err != nil {
			return nil, err
		}
	}

	if !validateStatusTransition(booking.Status, completed.Status) {
		return nil, errs.ErrStatusNotAllowed
	}
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

func applyPatch(bu models.BookingUpdate, b models.Booking) finalbook {
	var (
		title, description string
		from_time, to_time time.Time
		pc                 int
		status             models.BookingStatus
		amount             int64
		paid               bool
	)
	from_time = b.FromTime
	to_time = b.ToTime
	title = b.Title
	description = b.Description
	pc = b.ParticipantsCount
	status = b.Status
	amount = b.Amount
	paid = b.Paid

	if bu.FromTime.Set && bu.FromTime.Val != nil {
		from_time = *bu.FromTime.Val
	}
	if bu.ToTime.Set && bu.ToTime.Val != nil {
		to_time = *bu.ToTime.Val
	}
	if bu.Title.Set && bu.Title.Val != nil {
		title = *bu.Title.Val
	}
	if bu.Description.Set && bu.Description.Val != nil {
		description = *bu.Description.Val
	}
	if bu.ParticipantsCount.Set && bu.ParticipantsCount.Val != nil {
		pc = *bu.ParticipantsCount.Val
	}
	if bu.Status.Set && bu.Status.Val != nil {
		status = *bu.Status.Val
	}
	if bu.Amount.Set && bu.Amount.Val != nil {
		amount = *bu.Amount.Val
	}
	if bu.Paid.Set && bu.Paid.Val != nil {
		paid = *bu.Paid.Val
	}

	result := finalbook{
		Title:             title,
		Description:       description,
		From_time:         from_time,
		To_Time:           to_time,
		ParticipantsCount: pc,
		Status:            status,
		Amount:            amount,
		Paid:              paid,
	}
	return result
}

func buildValidators(bkg models.BookingUpdate, room models.Room) []func() error {
	var Validators []func() error

	if bkg.Title.HasValue() {
		Validators = append(Validators, func() error {
			return validateTitle(*bkg.Title.Val)
		})
	}

	if bkg.Description.HasValue() {
		Validators = append(Validators, func() error {
			return validateDescription(*bkg.Description.Val)
		})
	}

	if bkg.FromTime.HasValue() {
		Validators = append(Validators, func() error {
			return validateBackdate(*bkg.FromTime.Val)
		})
	}

	if bkg.FromTime.HasValue() && bkg.ToTime.HasValue() {
		Validators = append(Validators, func() error {
			return validateTimeRange(*bkg.FromTime.Val, *bkg.ToTime.Val)
		})
		Validators = append(Validators, func() error {
			return validateMinMeetingTime(*bkg.FromTime.Val, *bkg.ToTime.Val)
		})
		Validators = append(Validators, func() error {
			return validateMaxMeetingTime(*bkg.FromTime.Val, *bkg.ToTime.Val)
		})
	}

	if bkg.ParticipantsCount.HasValue() {
		Validators = append(Validators, func() error {
			return validateParticipantsNumber(*bkg.ParticipantsCount.Val)
		})
		Validators = append(Validators, func() error {
			return validateMeetingCapacity(room.Capacity, *bkg.ParticipantsCount.Val)
		})
	}

	if bkg.Amount.HasValue() {
		Validators = append(Validators, func() error {
			return validateAmount(*bkg.Amount.Val)
		})
	}
	return nil
}
