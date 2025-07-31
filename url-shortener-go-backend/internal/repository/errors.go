package repository

import "errors"

var ErrUniqueViolation = errors.New("unique violation")
var (
	ErrAccountCreation = errors.New("error creating account")
	ErrEmailInUse      = errors.New("email is already in use")
	ErrEmailInvalid    = errors.New("email is not valid")
	ErrUserNotFound    = errors.New("user not found")
)
