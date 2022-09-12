package storage

import (
	"context"
	"regexp"
	"time"

	"github.com/ihopethisisfine/helloworld/internal/domain"
)

type User struct {
	Username    string `json:"username"`
	DateOfBirth string `json:"dateOfBirth"`
}

func (u User) Validate() error {
	var IsLetter = regexp.MustCompile(`^[a-zA-Z]+$`).MatchString

	if !IsLetter(u.Username) {
		return domain.ErrInvalidUsername
	}

	isDateBeforeToday, err := isDateBeforeToday(u.DateOfBirth)

	if err != nil {
		return domain.ErrInvalidDate
	}

	if !isDateBeforeToday {
		return domain.ErrInvalidBirthDate
	}

	return nil
}

func isDateBeforeToday(dateString string) (bool, error) {
	date, err := time.Parse("2006-01-02", dateString)
	if err != nil {
		return false, err
	}
	return date.Before(time.Now()), nil
}

type UserStorer interface {
	Put(ctx context.Context, user User) error
	Find(ctx context.Context, username string) (User, error)
}
