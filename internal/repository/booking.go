package repository

import (
	"context"
	"errors"
	errs "room-booking/internal/errors"
	"room-booking/internal/models"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

const (
	constraint = "chk_booking_status"

	bookingTableName              = "booking"
	bookingIDColumnName           = "ID"
	bookingRoomIDColumnName       = "room_id"
	bookingUserIDColumnName       = "user_id"
	bookingTitleColumnName        = "title"
	bookingDescriptionColumnName  = "description"
	bookingParticipantsColumnName = "participants_count"
	bookingFromTimeColumnName     = "from_time"
	bookingToTimeColumnName       = "to_time"
	bookingStatusColumnName       = "status"
	bookingAmountTableColumn      = "amount"
	bookingPaidColumnName         = "paid"
	bookingCreatedAtColumnName    = "created_at"
	bookingUpdatedAtColumnName    = "updated_at"
)

func (p *pgxRepository) InsertBooking(ctx context.Context, roomID int, userID uuid.UUID, title string, description string, participants int, fromTime time.Time, toTime time.Time, status models.BookingStatus, amount float64, paid bool) (*models.Booking, error) {
	query, args, err := psql.
		Insert(bookingTableName).
		Columns(bookingRoomIDColumnName, bookingUserIDColumnName, bookingTitleColumnName, bookingDescriptionColumnName, bookingParticipantsColumnName, bookingFromTimeColumnName, bookingToTimeColumnName, bookingStatusColumnName, bookingAmountTableColumn, bookingPaidColumnName).
		Values(roomID, userID, title, description, participants, fromTime, toTime, status, amount, paid).
		Suffix("RETURNING " + bookingIDColumnName + ", " + bookingTitleColumnName).
		ToSql()

	if err != nil {
		return nil, err
	}

	var booking models.Booking

	err = p.db.QueryRow(ctx, query, args...).Scan(&booking.ID, &booking.Title)
	if err != nil {
		return nil, err
	}
	return &booking, nil
}

func (p *pgxRepository) GetBooking(ctx context.Context, id int) (*models.Booking, error) {
	query, args, err := psql.
		Select(bookingIDColumnName, bookingRoomIDColumnName, bookingUserIDColumnName, bookingTitleColumnName, bookingDescriptionColumnName, bookingParticipantsColumnName, bookingFromTimeColumnName, bookingToTimeColumnName, bookingStatusColumnName, bookingAmountTableColumn, bookingPaidColumnName, bookingCreatedAtColumnName, bookingUpdatedAtColumnName).
		From(bookingTableName).
		Where(sq.Eq{bookingIDColumnName: id}).
		ToSql()

	if err != nil {
		return nil, err
	}

	var bkg models.Booking
	err = p.db.QueryRow(ctx, query, args...).Scan(&bkg.ID, &bkg.RoomID, &bkg.UserID, &bkg.Title, &bkg.Description, &bkg.ParticipantsCount, &bkg.FromTime, &bkg.ToTime, &bkg.Status, &bkg.Amount, &bkg.Paid, &bkg.CreatedAt, &bkg.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errs.ErrBookingNotFound
		}
		return nil, err
	}

	return &bkg, nil
}

func (p *pgxRepository) UpdateBooking(ctx context.Context, b models.Booking) (*models.Booking, error) {
	query, args, err := psql.Update(bookingTableName).
		Set(bookingIDColumnName, b.ID).
		Set(bookingRoomIDColumnName, b.RoomID).
		Set(bookingTitleColumnName, b.Title).
		Set(bookingDescriptionColumnName, b.Description).
		Set(bookingFromTimeColumnName, b.FromTime).
		Set(bookingToTimeColumnName, b.ToTime).
		Set(bookingStatusColumnName, b.Status).
		Set(bookingAmountTableColumn, b.Amount).
		Set(bookingPaidColumnName, b.Paid).Where(sq.Eq{bookingIDColumnName: b.ID}).
		Suffix("RETURNING " + bookingIDColumnName + ", " + bookingTitleColumnName + ", " + bookingStatusColumnName).
		ToSql()

	if err != nil {
		return nil, err
	}

	var booking models.Booking
	err = p.db.QueryRow(ctx, query, args...).Scan(&booking.ID, &booking.Title, &booking.Status)
	if err != nil { //---------------------------проверяем на ошибки. Если err!=nil, то..
		var pgErr *pgconn.PgError
		if ok := errors.As(err, &pgErr); ok && pgErr.ConstraintName == constraint { //.. пытаемся привести ошибку к типу pgconn.PgErr.
			//---------------------------------Если удачно и в поле ConstraintName содержится имя ограничения (смотри миграцию таблицы booking,поле status), то возвращаем кастомную ошибку
			return nil, errs.ErrStatusNotAllowed
		}
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errs.ErrBookingNotFound
		}
		return nil, err
	}

	return &booking, nil
}

func (p *pgxRepository) ListBookings(ctx context.Context) ([]models.Booking, error) {
	query, _, err := psql.
		Select(bookingIDColumnName, bookingRoomIDColumnName, bookingUserIDColumnName, bookingTitleColumnName, bookingDescriptionColumnName, bookingParticipantsColumnName, bookingFromTimeColumnName, bookingToTimeColumnName, bookingStatusColumnName, bookingAmountTableColumn, bookingPaidColumnName, bookingCreatedAtColumnName, bookingUpdatedAtColumnName).
		From(bookingTableName).
		ToSql()

	if err != nil {
		return nil, err
	}

	rows, err := p.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var bookings []models.Booking

	for rows.Next() {
		var bkg models.Booking
		rows.Scan(&bkg.ID, &bkg.RoomID, &bkg.UserID, &bkg.Title, &bkg.Description, &bkg.ParticipantsCount, &bkg.FromTime, &bkg.ToTime, &bkg.Status, &bkg.Amount, &bkg.Paid, &bkg.CreatedAt, &bkg.UpdatedAt)
		bookings = append(bookings, bkg)
	}
	defer rows.Close()

	if err = rows.Err(); err != nil {
		return nil, err
	}
	return bookings, nil
}

func (p *pgxRepository) HasOverlaped(ctx context.Context, bookingID int, roomID int, FromTime, ToTime time.Time) (bool, error) {
	const query = "SELECT EXIST (SELECT 1 FROM booking WHERE room_id=$1 AND id <> $2 AND status in ('pending','confirmed') AND from_time < $4 AND to_time > $3);"

	var overlaped bool

	err := p.db.QueryRow(ctx, query, roomID, ToTime, FromTime).Scan(&overlaped)
	if err != nil {
		return false, err
	}
	return overlaped, nil
}
