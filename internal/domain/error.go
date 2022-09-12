package domain

import "errors"

var (
	ErrInternal         = errors.New("internal")
	ErrNotFound         = errors.New("not found")
	ErrInvalidBirthDate = errors.New("birthdate must be before today")
	ErrInvalidDate      = errors.New("date is invalid (must be a valid date and format must be YYYY-MM-DD)")
	ErrInvalidUsername  = errors.New("username should contain only letters")
)
