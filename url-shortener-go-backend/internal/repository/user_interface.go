package repository

import "url-shortener-go-backend/internal/model"

type UserRepository interface {
	CreateUser(email string) error
	GetUserByEmail(email string) (*model.User, error)
}
