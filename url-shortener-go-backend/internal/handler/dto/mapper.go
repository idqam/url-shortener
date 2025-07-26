package dto

import (
	"time"
	"url-shortener-go-backend/internal/model"
)

const timeLayout = time.RFC3339 

func ToURLResponse(u model.URL) URLResponse {
	return URLResponse{
		ID:          u.ID,
		UserID:      u.UserID,
		OriginalURL: u.OriginalURL,
		ShortCode:   u.ShortCode,
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

func ToCreateURLResponse(u *model.URL) CreateURLResponse {
	return CreateURLResponse{
		ID:          u.ID,
		OriginalURL: u.OriginalURL,
		ShortCode:   u.ShortCode,
	}
}
