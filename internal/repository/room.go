package repository

import (
	"context"
	"errors"
	errs "room-booking/internal/errors"
	"room-booking/internal/models"

	"github.com/Masterminds/squirrel"
	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

var (
	psql = squirrel.StatementBuilder.PlaceholderFormat(sq.Dollar)
)

const (
	room_constraint = "chk_currency"

	roomTableName                = "room"
	roomIDTableColumn            = "id"
	roomNameTableColumn          = "name"
	roomCapacityTableColumn      = "capacity"
	roomPriceTableColumn         = "hour_price"
	roomCurrencyTableColumn      = "currency"
	roomCreatedAtRoomTableColumn = "created_at"
)

func (p *pgxRepository) CreateRoom(ctx context.Context, room models.Room) (*models.Room, error) {
	query, args, err := psql.Insert(roomTableName).
		Columns(roomNameTableColumn, roomCapacityTableColumn, roomPriceTableColumn, roomCurrencyTableColumn).
		Values(room.Name, room.Capacity, room.HourPrice, room.Currency).
		Suffix("RETURNING " + roomIDTableColumn + ", " + roomNameTableColumn + ", " + roomCreatedAtRoomTableColumn).
		ToSql()

	if err != nil {
		return nil, err
	}

	var r models.Room
	err = p.db.QueryRow(ctx, query, args...).Scan(&r.ID, &r.Name, &r.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &r, nil
}

func (p *pgxRepository) GetRoom(ctx context.Context, id int) (*models.Room, error) {
	query, args, err := psql.
		Select(roomIDTableColumn, roomNameTableColumn, roomCapacityTableColumn, roomPriceTableColumn, roomCurrencyTableColumn, roomCreatedAtRoomTableColumn).
		From(roomTableName).
		Where(sq.Eq{roomIDTableColumn: id}).
		ToSql()

	if err != nil {
		return nil, err
	}
	var room models.Room

	err = p.db.QueryRow(ctx, query, args...).Scan(&room.ID, &room.Name, &room.Capacity, &room.HourPrice, &room.Currency, &room.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errs.ErrRoomNotFound
		}
		return nil, err
	}

	return &room, nil
}

func (p *pgxRepository) UpdateRoom(ctx context.Context, r models.Room) (*models.Room, error) {
	query, args, err := psql.Update(roomTableName).
		Set(roomNameTableColumn, r.Name).
		Set(roomCapacityTableColumn, r.Capacity).
		Set(roomPriceTableColumn, r.HourPrice).
		Set(roomCurrencyTableColumn, r.Currency).
		Where(sq.Eq{roomIDTableColumn: r.ID}).
		Suffix("RETURNING " + roomNameTableColumn + ", " + roomCapacityTableColumn + ", " + roomPriceTableColumn + ", " + roomCurrencyTableColumn).
		ToSql()
	if err != nil {
		return nil, err
	}

	var room models.Room
	err = p.db.QueryRow(ctx, query, args...).Scan(&room)
	if err != nil {
		var pgErr *pgconn.PgError
		if ok := errors.As(err, &pgErr); ok && pgErr.ConstraintName == room_constraint {
			return nil, errs.ErrCurrencyNotAllowed
		}
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errs.ErrBookingNotFound
		}
		return nil, err
	}

	return &room, nil
}

func (p *pgxRepository) ListRooms(ctx context.Context) ([]models.Room, error) {
	query, _, err := psql.
		Select(roomTableName).
		Column(roomIDTableColumn, roomNameTableColumn, roomCapacityTableColumn, roomPriceTableColumn, roomCurrencyTableColumn, roomCreatedAtRoomTableColumn).
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

func (p *pgxRepository) Exist(ctx context.Context, roomID int) (bool, error) {
	query := "SELECT EXISTS(SELECT 1 FROM room WHERE id = $1)"

	var exist bool

	err := p.db.QueryRow(ctx, query, roomID).Scan(&exist)

	if err != nil {
		return false, err
	}
	return exist, nil
}
