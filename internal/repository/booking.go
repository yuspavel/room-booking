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
)

const (
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
	bookingPaidColumnName         = "paid"
	bookingCreatedAtColumnName    = "created_at"
	bookingUpdatedAtColumnName    = "updated_at"
)

func (p *pgxRepository) InsertBooking(ctx context.Context, roomID int, userID uuid.UUID, title string, description string, participants int, fromTime time.Time, toTime time.Time, status models.BookingStatus, paid bool) (*models.Booking, error) {
	query, args, err := psql.
		Insert(bookingTableName).
		Columns(bookingRoomIDColumnName, bookingUserIDColumnName, bookingTitleColumnName, bookingDescriptionColumnName, bookingParticipantsColumnName, bookingFromTimeColumnName, bookingToTimeColumnName, bookingStatusColumnName, bookingPaidColumnName).
		Values(roomID, userID, title, description, participants, fromTime, toTime, status, paid).
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
		Select(bookingIDColumnName, bookingRoomIDColumnName, bookingUserIDColumnName, bookingTitleColumnName, bookingDescriptionColumnName, bookingParticipantsColumnName, bookingFromTimeColumnName, bookingToTimeColumnName, bookingStatusColumnName, bookingPaidColumnName, bookingCreatedAtColumnName, bookingUpdatedAtColumnName).
		From(bookingTableName).
		Where(sq.Eq{bookingIDColumnName: id}).
		ToSql()

	if err != nil {
		return nil, err
	}

	var bkg models.Booking
	err = p.db.QueryRow(ctx, query, args...).Scan(&bkg.ID, &bkg.RoomID, &bkg.UserID, &bkg.Title, &bkg.Description, &bkg.ParticipantsCount, &bkg.FromTime, &bkg.ToTime, &bkg.Status, &bkg.Paid, &bkg.CreatedAt, &bkg.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errs.ErrNoRowsFound
		}
		return nil, err
	}

	return &bkg, nil
}

func (p *pgxRepository) UpdateBooking(ctx context.Context, id int, b models.BookingUpdate) (int, error) {
	updated := false

	builder := psql.Update(bookingTableName)

	if b.RoomID != nil {
		builder = builder.Set(bookingRoomIDColumnName, *b.RoomID)
		updated = true
	}
	if b.Title != nil {
		builder = builder.Set(bookingTitleColumnName, *b.Title)
		updated = true
	}
	if b.Descriptions != nil {
		builder = builder.Set(bookingDescriptionColumnName, *b.Descriptions)
		updated = true
	}
	if b.ParticipantsCount != nil {
		builder = builder.Set(bookingParticipantsColumnName, *b.ParticipantsCount)
		updated = true
	}
	if b.FromDate != nil {
		builder = builder.Set(bookingFromTimeColumnName, *b.FromDate)
		updated = true
	}
	if b.ToDate != nil {
		builder = builder.Set(bookingToTimeColumnName, *b.ToDate)
		updated = true
	}
	if b.Status != nil {
		builder = builder.Set(bookingStatusColumnName, *b.Status)
		updated = true
	}
	if b.Paid != nil {
		builder = builder.Set(bookingPaidColumnName, *b.Paid)
		updated = true
	}
	if !updated {
		return 0, errs.ErrInvalidInput
	}

	builder = builder.Set(bookingUpdatedAtColumnName, time.Now()) //-----------Обязательное обновляем поле updated_at
	query, args, err := builder.Where(sq.Eq{bookingIDColumnName: id}).
		ToSql()

	if err != nil {
		return 0, err
	}

	c, err := p.db.Exec(ctx, query, args...)
	if err != nil {
		return 0, err
	}

	if int(c.RowsAffected()) == 0 {
		return 0, errs.ErrNoRowsFound
	}
	return int(c.RowsAffected()), nil
}

func (p *pgxRepository) ListBookings(ctx context.Context) ([]models.Booking, error) {
	query, _, err := psql.
		Select(bookingIDColumnName, bookingRoomIDColumnName, bookingUserIDColumnName, bookingTitleColumnName, bookingDescriptionColumnName, bookingParticipantsColumnName, bookingFromTimeColumnName, bookingToTimeColumnName, bookingStatusColumnName, bookingPaidColumnName, bookingCreatedAtColumnName, bookingUpdatedAtColumnName).
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
		rows.Scan(&bkg.ID, &bkg.RoomID, &bkg.UserID, &bkg.Title, &bkg.Description, &bkg.ParticipantsCount, &bkg.FromTime, &bkg.ToTime, &bkg.Status, &bkg.Paid, &bkg.CreatedAt, &bkg.UpdatedAt)
		bookings = append(bookings, bkg)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	return bookings, nil
}
