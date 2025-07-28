package dto

import (
	"time"
	"url-shortener-go-backend/internal/model"
)

const timeLayout = time.RFC3339

func ToURLResponse(u model.URL) URLResponse {
	u.PopulateShortURL()
	return URLResponse{
		ID:          u.ID,
		UserID:      u.UserID,
		OriginalURL: u.OriginalURL,
		ShortCode:   u.ShortCode,
		ShortURL:    u.ShortURL,
		IsPublic:    u.IsPublic,
		ClickCount:  u.ClickCount,
		CreatedAt:   u.CreatedAt.UTC().Format(timeLayout),
	}
}

func ToGetUrlsResponse(urls []model.URL) GetUrlsResponse {
	resp := GetUrlsResponse{
		URLs: make([]URLResponse, len(urls)),
	}
	for i, u := range urls {
		resp.URLs[i] = ToURLResponse(u)
	}
	return resp
}

func ToShortenResponse(u *model.URL) ShortenResponse {
	u.PopulateShortURL()
	return ShortenResponse{
		ID:          u.ID,
		ShortCode:   u.ShortCode,
		OriginalURL: u.OriginalURL,
		ShortURL:    u.ShortURL,
	}
}
