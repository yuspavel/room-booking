package repository

import (
	"github.com/jackc/pgx/v5/pgxpool"
)

type pgxRepository struct {
	db *pgxpool.Pool
}

func NewPgxRepository(db *pgxpool.Pool) *pgxRepository {
	return &pgxRepository{db: db}
}
