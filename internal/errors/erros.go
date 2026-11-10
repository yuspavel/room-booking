package errs

import (
	"errors"
)

var (
	//common errors
	ErrInvalidInput = errors.New("invalid input")
	//room errors
	ErrRoomNotFound        = errors.New("room not found")
	ErrInvalidRoomName     = errors.New("invalid room name")
	ErrInvalidRoomPrice    = errors.New("invalid room price")
	ErrCurrencyNotAllowed  = errors.New("invalid currency")
	ErrInvalidRoomCapacity = errors.New("invalid room capacity")
	//user errors
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrEmailAlreadyExist  = errors.New("email already exist")
	//booking errors
	ErrBookingNotFound         = errors.New("booking not found")
	ErrInvalidTimeRange        = errors.New("invalid bookig time range")
	ErrInvalidBookingTime      = errors.New("booking time lower than now")
	ErrBookingOverlap          = errors.New("booking overlap detected")
	ErrInvalidStatusTransition = errors.New("invalid status")
	ErrMaxTimeExceed           = errors.New("maximum time exceeded")
	ErrMinTimeLack             = errors.New("minimum meeting time 30 min allowed")
	ErrMaxPaticipantsExceeded  = errors.New("too much participants")
	ErrNoParticipantsInvolved  = errors.New("participants can be less or equal 0")
	ErrNegativeAmount          = errors.New("amount cannot be negative")
)
