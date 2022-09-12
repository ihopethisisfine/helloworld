package user

type User struct {
	DateOfBirth string `json:"dateOfBirth"`
}

type Response struct {
	Message string `json:"message"`
}
