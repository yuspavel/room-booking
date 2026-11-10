package models

import (
	"room-booking/internal/patch"
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `db:"id"`
	Name      string    `db:"name"`
	Password  string    `db:"pass_hash"`
	Email     string    `db:"email"`
	Phone     string    `db:"phone"`
	CreatedAt time.Time `db:"created_at"`
}

type UserUpdate struct {
	ID       uuid.UUID
	Name     patch.Field[string]
	Email    string
	Password patch.Field[string]
	Phone    patch.Field[string]
}

func (u *User) ApplyPatch(user UserUpdate) {
	if user.Name.HasValue() {
		u.Name = *user.Name.Val
	}
	u.Email = user.Email
	if user.Password.HasValue() {
		u.Password = *user.Password.Val
	}
	if user.Phone.HasValue() {
		u.Phone = *user.Phone.Val
	}
}
