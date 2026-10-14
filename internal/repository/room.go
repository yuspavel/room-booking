package repository

import (
	"context"
	"errors"
	errs "room-booking/internal/errors"
	"room-booking/internal/models"

	"github.com/Masterminds/squirrel"
	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

var (
	psql = squirrel.StatementBuilder.PlaceholderFormat(sq.Dollar)
)

const (
	roomTableName            = "room"
	idRoomTableColumn        = "id"
	nameRoomTableColumn      = "name"
	capacityTableColumn      = "capacity"
	priceTableColumn         = "hour_price"
	currencyTableColumn      = "currency"
	createdAtRoomTableColumn = "created_at"
)

func (p *pgxRepository) CreateRoom(ctx context.Context, name string, capacity int, price float64, currency string) (*models.Room, error) {
	query, args, err := psql.Insert(roomTableName).
		Columns(nameRoomTableColumn, capacityTableColumn, priceTableColumn, currencyTableColumn).
		Values(name, capacity, price, currency).
		Suffix("RETURNING " + idRoomTableColumn + ", " + createdAtRoomTableColumn).
		ToSql()

	if err != nil {
		return nil, err
	}

	var room models.Room
	err = p.db.QueryRow(ctx, query, args...).Scan(&room.ID, &room.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &room, nil
}

func (p *pgxRepository) GetRoom(ctx context.Context, id uuid.UUID) (*models.Room, error) {
	query, args, err := psql.
		Select(idRoomTableColumn, nameRoomTableColumn, capacityTableColumn, priceTableColumn, currencyTableColumn, createdAtRoomTableColumn).
		From(roomTableName).
		Where(sq.Eq{idRoomTableColumn: id}).
		ToSql()

	if err != nil {
		return nil, err
	}
	var room models.Room

	err = p.db.QueryRow(ctx, query, args...).Scan(&room.ID, &room.Name, &room.Capacity, &room.HourPrice, &room.Currency, &room.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errs.ErrNoRowsFound
		}
		return nil, err
	}

	return &room, nil
}

func (p *pgxRepository) UpdateRoom(ctx context.Context, roomID int, room models.RoomUpdate) (int, error) {
	updated := false
	builder := psql.Update(roomTableName)

	if room.Name != nil {
		builder = builder.Set(nameRoomTableColumn, *room.Name)
		updated = true
	}
	if room.Capacity != nil {
		builder = builder.Set(capacityTableColumn, *room.Capacity)
		updated = true
	}
	if room.HourPrice != nil {
		builder = builder.Set(priceTableColumn, *room.HourPrice)
		updated = true
	}
	if room.Currency != nil {
		builder = builder.Set(currencyTableColumn, *room.Currency)
		updated = true
	}

	if !updated {
		return 0, errs.ErrInvalidInput
	}
	query, args, err := builder.
		Where(sq.Eq{idRoomTableColumn: roomID}).
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

func (p *pgxRepository) ListRooms(ctx context.Context) ([]models.Room, error) {
	query, _, err := psql.
		Select(roomTableName).
		Column(idRoomTableColumn, nameRoomTableColumn, capacityTableColumn, priceTableColumn, currencyTableColumn, createdAtRoomTableColumn).
		ToSql()

	if err != nil {
		return nil, err
	}

	rows, err := p.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rooms []models.Room
	for rows.Next() {
		var room models.Room
		rows.Scan(&room.ID,
			&room.Name,
			&room.Capacity,
			&room.HourPrice,
			&room.Currency,
			&room.CreatedAt)

		rooms = append(rooms, room)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	return rooms, nil

}
