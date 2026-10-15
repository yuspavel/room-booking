package service

import (
	"context"
	"errors"
	"fmt"
	errs "room-booking/internal/errors"
	"room-booking/internal/models"

	"golang.org/x/crypto/bcrypt"
)

type UserRepo interface {
	CreateUser(ctx context.Context, name, passhash, email, phone string) (models.User, error)
	GetUser(ctx context.Context, email string) (*models.User, error)
	UpdateUser(ctx context.Context, email string, user models.UserUpdate) (int, error)
	//ListUsers(ctx context.Context) ([]models.User, error)
}

type Service struct {
	urepo UserRepo
}

// ОБРАТИТЬ ВНИМАНИЕ: на передачу возвращаемой структуры models.User. В этом методе она передается по значению
func (s *Service) Register(ctx context.Context, u models.User) (models.User, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(u.PassHash), bcrypt.DefaultCost)
	if err != nil {
		return models.User{}, err
	}

	user, err := s.urepo.CreateUser(ctx, u.Name, string(hash), u.Email, u.Phone)
	if err != nil {
		if errors.Is(err, errs.ErrEmailAlreadyExist) {
			return models.User{}, fmt.Errorf("user already exist")
		}
		return models.User{}, err
	}
	return user, nil
}

// ОБРАТИТЬ ВНИМАНИЕ: на передачу возвращаемой структуры models.User. В этом методе она передается по ссылке
func (s *Service) Authorize(ctx context.Context, email, password string) (*models.User, error) {
	user, err := s.urepo.GetUser(ctx, email)
	if err != nil {
		if errors.Is(err, errs.ErrNoRowsFound) {
			return nil, fmt.Errorf("invalid credentials")
		}
		return nil, err
	}

	hashedPassword := []byte(user.PassHash)
	if err = bcrypt.CompareHashAndPassword(hashedPassword, hashedPassword); err != nil {
		return nil, err
	}
	return user, nil
}

func (s *Service) UpdateUserData(ctx context.Context, email string, user models.UserUpdate) (int, error) {
	c, err := s.urepo.UpdateUser(ctx, email, user)
	if err != nil {
		if errors.Is(err, errs.ErrEmailAlreadyExist) {
			return 0, fmt.Errorf("email already exist")
		}
		return 0, err
	}
	return c, nil
}
