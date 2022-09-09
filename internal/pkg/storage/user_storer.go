package storage

import (
	"context"
	"time"
)

type User struct {
	Username    string    `json:"username"`
	DateOfBirth time.Time `json:"dateOfBirth"`
}

type UserStorer interface {
	Insert(ctx context.Context, user User) error
	Find(ctx context.Context, username string) (User, error)
}
