package repository

import (
	"context"
	"errors"
	errs "room-booking/internal/errors"
	"room-booking/internal/models"

	"github.com/google/uuid"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

const (
	userTableName            = "users"
	idUserTableColumn        = "id"
	nameUserTableColumn      = "name"
	passHashTableColumn      = "pass_hash"
	emailTableColumn         = "email"
	phoneTableColumn         = "phone"
	createdAtUserTableColumn = "created_at"
)

var ()

func (p *pgxRepository) CreareUser(ctx context.Context, name string, passhash string, email string, phone string) (*models.User, error) {
	query, args, err := psql.
		Insert(userTableName).
		Columns(nameUserTableColumn, passHashTableColumn, emailTableColumn, phoneTableColumn).
		Values(name, passhash, email, phone).
		Suffix("RETURNING " + idUserTableColumn + ", " + nameUserTableColumn).
		ToSql()

	if err != nil {
		return nil, err
	}

	var user models.User

	err = p.db.QueryRow(ctx, query, args...).Scan(&user.ID, &user.Name)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) { //--------------------------В таблице users имеется уникальное поле email.
			if pgErr.Code == "23505" { //-----------------------Если произойдет вставка не уникального значения...
				return nil, errs.ErrEmailAlreadyExist //--возвращаем кастомную ошибку
			}
		}
		return nil, err
	}
	return &user, nil
}

func (p *pgxRepository) GetUser(ctx context.Context, id uuid.UUID) (*models.User, error) {
	query, args, err := psql.
		Select(idUserTableColumn, nameUserTableColumn, emailTableColumn, phoneTableColumn, createdAtUserTableColumn).
		From(userTableName).
		Where(sq.Eq{idUserTableColumn: id}).
		ToSql()

	if err != nil {
		return &models.User{}, err
	}

	var user models.User

	err = p.db.QueryRow(ctx, query, args...).Scan(&user.ID, &user.Name, &user.Email, &user.Phone, &user.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return &models.User{}, errs.ErrNoRowsFound
		}
		return &models.User{}, err
	}
	return &user, nil
}

func (p *pgxRepository) UpdateUser(ctx context.Context, id uuid.UUID, user models.UserUpdate) (int, error) {
	updated := false
	builder := psql.
		Update(userTableName)

	if user.Name != nil {
		builder = builder.Set(nameUserTableColumn, *user.Name)
		updated = true
	}
	if user.PassHash != nil {
		builder = builder.Set(passHashTableColumn, *user.PassHash)
		updated = true
	}
	if user.Email != nil {
		builder = builder.Set(emailTableColumn, *user.Email)
		updated = true
	}
	if user.Phone != nil {
		builder = builder.Set(phoneTableColumn, *user.Phone)
		updated = true
	}

	if !updated {
		return 0, errs.ErrInvalidInput
	}

	query, args, err := builder.
		Where(sq.Eq{idUserTableColumn: id}).
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

func (p *pgxRepository) ListUsers(ctx context.Context) ([]models.User, error) {
	query, _, err := psql.
		Select(idUserTableColumn, nameUserTableColumn, emailTableColumn, phoneTableColumn, createdAtUserTableColumn).
		From(userTableName).
		OrderBy(nameUserTableColumn).
		ToSql()

	if err != nil {
		return nil, err
	}

	rows, err := p.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		rows.Scan(
			&user.ID,
			&user.Name,
			&user.Email,
			&user.Phone,
			&user.CreatedAt,
		)
		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	return users, nil
}
