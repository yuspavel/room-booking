package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `db:"id" json:"id"`
	Name      string    `db:"name" json:"name"`
	Password  string    `db:"pass_hash" json:"password"`
	Email     string    `db:"email" json:"email"`
	Phone     string    `db:"phone" json:"phone"`
	CreatedAt time.Time `db:"created_at" json:"createdat"`
}

type UserUpdate struct {
	Name     *string
	PassHash *string
	Phone    *string
}
