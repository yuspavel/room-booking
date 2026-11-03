package service

import (
	"context"
	"errors"
	errs "room-booking/internal/errors"
	"room-booking/internal/models"

	"golang.org/x/crypto/bcrypt"
)

type UserRepo interface {
	CreateUser(ctx context.Context, name, passhash, email, phone string) (models.User, error)
	GetUser(ctx context.Context, email string) (*models.User, error)
	UpdateUser(ctx context.Context, email string, user models.UserUpdate) (*models.User, error)
}

type UserService struct {
	urepo UserRepo
}

// ОБРАТИТЬ ВНИМАНИЕ: на передачу возвращаемой структуры models.User. В этом методе она передается по значению
func (s *UserService) Register(ctx context.Context, u models.User) (models.User, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return models.User{}, err
	}

	user, err := s.urepo.CreateUser(ctx, u.Name, string(hash), u.Email, u.Phone)
	if err != nil {
		if errors.Is(err, errs.ErrEmailAlreadyExist) {
			return models.User{}, errs.ErrEmailAlreadyExist
		}
		return models.User{}, err
	}
	return user, nil
}

// ОБРАТИТЬ ВНИМАНИЕ: на передачу возвращаемой структуры models.User. В этом методе она передается по ссылке
func (s *UserService) Authorize(ctx context.Context, email, password string) (*models.User, error) {
	user, err := s.urepo.GetUser(ctx, email)
	if err != nil {
		if errors.Is(err, errs.ErrUserNotFound) {
			return nil, errs.ErrUserNotFound
		}
		return nil, err
	}

	hashedPassword := []byte(user.Password)
	if err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(password)); err != nil {
		return nil, err
	}
	return user, nil
}

func (s *UserService) UpdateUserData(ctx context.Context, email string, u models.UserUpdate) (*models.User, error) {
	user, err := s.urepo.UpdateUser(ctx, email, u)
	if err != nil {
		if errors.Is(err, errs.ErrEmailAlreadyExist) {
			return nil, errs.ErrEmailAlreadyExist
		}
		return nil, err
	}
	return user, nil
}
