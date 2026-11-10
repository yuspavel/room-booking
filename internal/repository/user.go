package repository

import (
	"context"
	"errors"
	errs "room-booking/internal/errors"
	"room-booking/internal/models"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

const (
	userTableName                = "users"
	userIDTableColumn            = "id"
	userNameTableColumn          = "name"
	userPassHashTableColumn      = "pass_hash"
	userEmailTableColumn         = "email"
	userPhoneTableColumn         = "phone"
	userCreatedAtUserTableColumn = "created_at"
)

func (p *pgxRepository) CreareUser(ctx context.Context, u models.User) (*models.User, error) {
	query, args, err := psql.
		Insert(userTableName).
		Columns(userNameTableColumn, userPassHashTableColumn, userEmailTableColumn, userPhoneTableColumn).
		Values(u.Name, u.Password, u.Email, u.Phone).
		Suffix("RETURNING " + userIDTableColumn + ", " + userNameTableColumn + ", " + userCreatedAtUserTableColumn).
		ToSql()

	if err != nil {
		return nil, err
	}

	var user models.User

	err = p.db.QueryRow(ctx, query, args...).Scan(&user.ID, &user.Name, &user.CreatedAt)
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

func (p *pgxRepository) GetUser(ctx context.Context, email string) (*models.User, error) {
	query, args, err := psql.
		Select(userIDTableColumn, userNameTableColumn, userEmailTableColumn, userPhoneTableColumn, userCreatedAtUserTableColumn).
		From(userTableName).
		Where(sq.Eq{userEmailTableColumn: email}).
		ToSql()

	if err != nil {
		return &models.User{}, err
	}

	var user models.User

	err = p.db.QueryRow(ctx, query, args...).Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.Phone, &user.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errs.ErrInvalidCredentials
		}
		return nil, err
	}
	return &user, nil
}

func (p *pgxRepository) UpdateUser(ctx context.Context, user models.User) (*models.User, error) {
	query, args, err := psql.
		Update(userTableName).
		Set(userNameTableColumn, user.Name).
		Set(userPassHashTableColumn, user.Password).
		Set(userPhoneTableColumn, user.Phone).
		Where(sq.Eq{userEmailTableColumn: user.Email}).
		ToSql()

	if err != nil {
		return nil, err
	}

	var usr models.User
	err = p.db.QueryRow(ctx, query, args...).Scan(&usr.ID, &usr.Name, &usr.Email)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) { //--------------------------В таблице users имеется уникальное поле email.
			if pgErr.Code == "23505" { //-----------------------Если произойдет вставка не уникального значения...
				return nil, errs.ErrEmailAlreadyExist //--------возвращаем кастомную ошибку
			}
		}
	}

	return &usr, nil
}

func (p *pgxRepository) ListUsers(ctx context.Context) ([]models.User, error) {
	query, _, err := psql.
		Select(userIDTableColumn, userNameTableColumn, userEmailTableColumn, userPhoneTableColumn, userCreatedAtUserTableColumn).
		From(userTableName).
		OrderBy(userNameTableColumn).
		ToSql()

	if err != nil {
		return nil, err
	}

	rows, err := p.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users, err := pgx.CollectRows(rows, pgx.RowToStructByNameLax[models.User]) //Выбираем все строки из rows в срез структур Users
	if err != nil {
		return nil, err
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	return users, nil
}
