package repository

import (
	"url-shortener-go-backend/internal/model"
)


type URLRepository interface {
	SaveURL(url *model.URL) error
	GetURLByShortCode(shortcode string) (*model.URL, error)
	IncrementClickCount(shortcode string) error
	GetUserUrls(user model.User) ([]model.URL, error) 
	GetURLByOriginalURL(original string) (*model.URL, error)
}
