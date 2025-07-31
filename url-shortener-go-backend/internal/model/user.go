package model

import "time"

type User struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

type CreateUserInput struct {
	Email string `json:"email"`
}
