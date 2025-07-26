package service

import (
	"fmt"
	"time"
	"url-shortener-go-backend/internal/model"
	"url-shortener-go-backend/internal/repository"
	"url-shortener-go-backend/internal/utils"
)

type URLService interface {
	CreateShortURL(rawURL string, userID *string, isPublic bool, codeLength int8) (*model.URL, error)
	GetURLByShortCode(code string) (*model.URL, error)
	GetUserUrls(user model.User) ([]model.URL, error)
	IncrementClickCount(shortcode string) error
}

type urlService struct {
	repo repository.URLRepository
}

func NewURLService(repo repository.URLRepository) URLService {
	return &urlService{repo: repo}
}

func (s *urlService) CreateShortURL(
    rawURL string,
    userID *string,
    isPublic bool,
    codeLength int8,
) (*model.URL, error) {
    parsed, err := utils.ValidateUrl(rawURL)
    if err != nil {
        return nil, fmt.Errorf("invalid URL: %w", err)
    }

    if codeLength < 6 || codeLength > 12 {
        codeLength = 8
    }

    if !isPublic && userID == nil {
        return nil, fmt.Errorf("userID must be provided for private URLs")
    }

    shortcode, err := utils.GenerateCode(parsed.Original, codeLength)
    if err != nil {
        return nil, fmt.Errorf("failed to generate shortcode: %w", err)
    }

    existing, err := s.repo.GetURLByShortCode(shortcode)
    if err != nil && err != utils.ErrNotFound {
        return nil, fmt.Errorf("failed to check shortcode uniqueness: %w", err)
    }
    if existing != nil {
        return nil, fmt.Errorf("generated shortcode already exists: %s", shortcode)
    }

    url := &model.URL{
        UserID:      userID,
        OriginalURL: rawURL,
        ShortCode:   shortcode,
        IsPublic:    isPublic,
        ClickCount:  0,
        CreatedAt:   time.Now().UTC(),
    }

    if err := s.repo.SaveURL(url); err != nil {
        return nil, fmt.Errorf("failed to save URL: %w", err)
    }

    return url, nil
}


func (s *urlService) GetURLByShortCode(shortcode string) (*model.URL, error) {
    return s.repo.GetURLByShortCode(shortcode)
}

func (s *urlService) IncrementClickCount(shortcode string) error {
    return s.repo.IncrementClickCount(shortcode)
}

func (s *urlService) GetUserUrls(user model.User) ([]model.URL, error) {
	return s.repo.GetUserUrls(user)
}
