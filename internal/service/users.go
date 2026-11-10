package service

import (
	"context"
	"errors"
	errs "room-booking/internal/errors"
	"room-booking/internal/models"

	"golang.org/x/crypto/bcrypt"
)

type UserRepo interface {
	CreateUser(ctx context.Context, user models.User) (models.User, error)
	GetUser(ctx context.Context, email string) (*models.User, error)
	UpdateUser(ctx context.Context, user models.UserUpdate) (*models.User, error)
	ListUsers(ctx context.Context) ([]models.User, error)
}

type UserService struct {
	urepo UserRepo
}

// ОБРАТИТЬ ВНИМАНИЕ: на передачу возвращаемой структуры models.User. В этом методе она передается по значению
func (s *UserService) Register(ctx context.Context, u models.User) (*models.User, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	u.Password = string(hash)
	user, err := s.urepo.CreateUser(ctx, u)
	if err != nil {
		if errors.Is(err, errs.ErrEmailAlreadyExist) {
			return nil, errs.ErrEmailAlreadyExist
		}
		return nil, err
	}
	return &user, nil
}

func (s *UserService) Authorize(ctx context.Context, email, password string) (*models.User, error) {
	user, err := s.urepo.GetUser(ctx, email)
	if err != nil {
		if errors.Is(err, errs.ErrInvalidCredentials) {
			return nil, errs.ErrInvalidCredentials
		}
		return nil, err
	}

	if err = checkPassword(user.Password, password); err != nil {
		return nil, err
	}
	return user, nil
}

func (s *UserService) UpdateUser(ctx context.Context, u models.UserUpdate) (*models.User, error) {
	user, err := s.urepo.UpdateUser(ctx, u)
	if err != nil {
		if errors.Is(err, errs.ErrEmailAlreadyExist) {
			return nil, errs.ErrEmailAlreadyExist
		}
		return nil, err
	}
	return user, nil
}

func checkPassword(hashedPass, password string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(hashedPass), []byte(password)); err != nil {
		return errs.ErrInvalidCredentials
	}
	return nil
}
