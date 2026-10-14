package repository

import (
	"context"
	"errors"
	"room-booking/internal/repository/models"
	"time"

	"github.com/Masterminds/squirrel"
	sq "github.com/Masterminds/squirrel"
)

var (
	psql = squirrel.StatementBuilder.PlaceholderFormat(sq.Dollar)
)

const (
	roomTableName        = "room"
	idTableColumn        = "id"
	nameTableColumn      = "name"
	capacityTableColumn  = "capacity"
	priceTableColumn     = "hour_price"
	currencyTableColumn  = "currency"
	createdAtTableColumn = "created_at"
)

var (
	ErrNoFieldsToUpdate = errors.New("no fields to update")
)

func (p *pgxRepository) CreateUser(ctx context.Context, name string, capacity int, price float64, currency string, created_at time.Time) (*models.Room, error) {
	query, args, err := psql.Insert(roomTableName).
		Columns(nameTableColumn, capacityTableColumn, priceTableColumn, currencyTableColumn, createdAtTableColumn).
		Values(name, capacity, price, currency, created_at).
		Suffix("RETURNING " + idTableColumn + ", " + createdAtTableColumn).
		ToSql()

	if err != nil {
		return &models.Room{}, err
	}

	var room models.Room
	err = p.db.QueryRow(ctx, query, args...).Scan(&room.ID, &room.CreatedAt)
	if err != nil {
		return &models.Room{}, err
	}
	return &room, nil
}

func (p *pgxRepository) GetRoom(ctx context.Context, id int) (*models.Room, error) {
	query, args, err := psql.
		Select(idTableColumn, nameTableColumn, capacityTableColumn, priceTableColumn, currencyTableColumn, createdAtTableColumn).
		From(roomTableName).
		Where(sq.Eq{idTableColumn: id}).
		ToSql()

	if err != nil {
		return &models.Room{}, err
	}
	var room models.Room

	err = p.db.QueryRow(ctx, query, args...).Scan(&room.ID, &room.Name, &room.Capacity, &room.HourPrice, &room.Currency, &room.CreatedAt)
	if err != nil {
		return &models.Room{}, err
	}

	return &room, nil
}

func (p *pgxRepository) UpdateRoom(ctx context.Context, roomID int, room models.UpdateRoom) (int, error) {
	updated := false
	builder := psql.Update(roomTableName)

	if room.Name != nil {
		builder = builder.Set(nameTableColumn, *room.Name)
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
		return 0, ErrNoFieldsToUpdate
	}
	query, args, err := builder.
		Where(sq.Eq{idTableColumn: roomID}).
		ToSql()

	c, err := p.db.Exec(ctx, query, args...)
	if err != nil {
		return 0, err
	}
	return int(c.RowsAffected()), nil
}
