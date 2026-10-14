package errs

import "errors"

var (
	ErrInvalidInput      = errors.New("invalid input")
	ErrNoRowsFound       = errors.New("no rows found")
	ErrEmailAlreadyExist = errors.New("email already exist")
)
