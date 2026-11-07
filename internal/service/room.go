package service

import (
	"context"
	"errors"
	errs "room-booking/internal/errors"
	"room-booking/internal/models"
)

type RoomRepo interface {
	CreateRoom(ctx context.Context, name string, capacity int, price int64, currency string) (*models.Room, error)
	GetRoom(ctx context.Context, id int) (models.Room, error)
	UpdateRoom(ctx context.Context, roomID int, room models.RoomUpdate) (int, error)
	ListRooms(ctx context.Context) ([]models.Room, error)
	Exist(ctx context.Context, roomID int) (bool, error)
}

type RoomService struct {
	rservice RoomRepo
}

func (s *RoomService) CreateRoom(ctx context.Context, r models.Room) (*models.Room, error) {
	room, err := s.rservice.CreateRoom(ctx, r.Name, r.Capacity, r.HourPrice, r.Currency)
	if err != nil {
		return nil, err
	}
	return room, nil
}

func (s *RoomService) GetRoom(ctx context.Context, id int) (models.Room, error) {
	room, err := s.rservice.GetRoom(ctx, id)
	if err != nil {
		if errors.Is(err, errs.ErrRoomNotFound) {
			return models.Room{}, errs.ErrRoomNotFound
		}
		return models.Room{}, err
	}
	return room, nil
}

func (s *RoomService) UpdateRoom(ctx context.Context, roomID int, room models.RoomUpdate) (int, error) {
	d, err := s.rservice.UpdateRoom(ctx, roomID, room)
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrInvalidInput):
			return 0, errs.ErrInvalidInput
		case errors.Is(err, errs.ErrRoomNotFound):
			return 0, errs.ErrRoomNotFound
		default:
			return 0, err
		}
	}
	return d, nil

}

func (s *RoomService) ListRooms(ctx context.Context) ([]models.Room, error) {
	rooms, err := s.rservice.ListRooms(ctx)
	if err != nil {
		return nil, err
	}
	return rooms, nil
}
