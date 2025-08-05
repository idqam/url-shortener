package mapper

import (
	"url-shortener-go-backend/internal/handler/dto"
	"url-shortener-go-backend/internal/model"
)

func ToShortenURLResponse(url model.URL) dto.ShortenURLResponse {
	return dto.ShortenURLResponse{
		ID:         url.ID,
		ShortCode:  url.ShortCode,
		ShortURL:   url.ShortURL,
		CreatedAt:  url.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		IsPublic:   url.IsPublic,
		ClickCount: url.ClickCount,
	}
}
func ToShortenURLResponses(urls []model.URL) []dto.ShortenURLResponse {
	var responses []dto.ShortenURLResponse
	for _, u := range urls {
		responses = append(responses, ToShortenURLResponse(u))
	}
	return responses
}

func ToGetUrlsResponse(urls []model.URL) dto.GetUserURLsResponse {
	return dto.GetUserURLsResponse{
		URLs: ToShortenURLResponses(urls),
	}
}

func ToGetURLByShortCodeResponse(url model.URL) dto.GetURLByShortCodeResponse {
	return dto.GetURLByShortCodeResponse{
		OriginalURL: url.OriginalURL,
		ClickCount:  url.ClickCount,
	}
}
