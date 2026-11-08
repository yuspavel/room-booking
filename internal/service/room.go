package service

import (
	"context"
	"errors"
	errs "room-booking/internal/errors"
	"room-booking/internal/models"
)

type RoomRepo interface {
	CreateRoom(ctx context.Context, room models.Room) (*models.Room, error)
	GetRoom(ctx context.Context, id int) (models.Room, error)
	UpdateRoom(ctx context.Context, room models.Room) (*models.Room, error)
	ListRooms(ctx context.Context) ([]models.Room, error)
}

type RoomService struct {
	roomRepo RoomRepo
}

func (r *RoomService) CreateRoom(ctx context.Context, rm models.Room) (*models.Room, error) {
	if err := rm.ValidateRoom(); err != nil {
		return nil, err
	}

	room, err := r.roomRepo.CreateRoom(ctx, rm)
	if err != nil {
		return nil, err
	}
	return room, nil
}

func (r *RoomService) GetRoom(ctx context.Context, id int) (models.Room, error) {
	room, err := r.roomRepo.GetRoom(ctx, id)
	if err != nil {
		if errors.Is(err, errs.ErrRoomNotFound) {
			return models.Room{}, errs.ErrRoomNotFound
		}
		return models.Room{}, err
	}
	return room, nil
}

func (r *RoomService) UpdateRoom(ctx context.Context, rUpd models.RoomUpdate) (*models.Room, error) {

	room, err := r.GetRoom(ctx, rUpd.ID)
	if err != nil {
		if errors.Is(err, errs.ErrRoomNotFound) {
			return nil, errs.ErrRoomNotFound
		}
		return nil, err
	}

	room.ApplyPatch(rUpd)
	if err = room.ValidateRoom(); err != nil {
		return nil, err
	}

	d, err := r.roomRepo.UpdateRoom(ctx, room)
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrInvalidInput):
			return nil, errs.ErrInvalidInput
		case errors.Is(err, errs.ErrRoomNotFound):
			return nil, errs.ErrRoomNotFound
		default:
			return nil, err
		}
	}
	return d, nil

}

func (r *RoomService) ListRooms(ctx context.Context) ([]models.Room, error) {
	rooms, err := r.roomRepo.ListRooms(ctx)
	if err != nil {
		return nil, err
	}
	return rooms, nil
}
