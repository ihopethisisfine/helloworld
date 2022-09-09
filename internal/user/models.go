package user

import "time"

type User struct {
	DateOfBirth time.Time `json:"dateOfBirth"`
}
