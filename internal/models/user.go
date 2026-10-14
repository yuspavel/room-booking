package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `db:"id"`
	Name      string    `db:"name"`
	PassHash  string    `db:"pass_hash"`
	Email     string    `db:"email"`
	Phone     string    `db:"phone"`
	CreatedAt time.Time `db:"created_at"`
}

type UserUpdate struct {
	Name     *string
	PassHash *string
	Email    *string
	Phone    *string
}
