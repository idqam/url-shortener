package service

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"url-shortener-go-backend/internal/cache"
	"url-shortener-go-backend/internal/model"
	"url-shortener-go-backend/internal/repository"
	"url-shortener-go-backend/internal/utils"
)

type URLService interface {
	CreateShortURL(ctx context.Context, rawURL string, userID *string, isPublic bool, codeLength int8) (*model.URL, error)
	GetURLByShortCode(ctx context.Context, code string) (*model.URL, error)
	GetUserUrls(user model.User) ([]model.URL, error)
	IncrementClickCount(shortcode string) error
}

type urlService struct {
	repo              repository.URLRepository
	cache             cache.Cache
	defaultCodeLength int8
	maxRetries        int
	cacheTTL          time.Duration
}

func NewURLService(repo repository.URLRepository, c cache.Cache) URLService {
	return &urlService{
		repo:              repo,
		cache:             c,
		defaultCodeLength: 8,
		maxRetries:        5,
		cacheTTL:          24 * time.Hour,
	}
}

func (s *urlService) CreateShortURL(ctx context.Context, rawURL string, userID *string, isPublic bool, codeLength int8) (*model.URL, error) {
	const logPrefix = "[CreateShortURL]"

	parsed, err := utils.ValidateUrl(rawURL)
	if err != nil {
		log.Printf("%s validate failed: %v", logPrefix, err)
		return nil, fmt.Errorf("invalid URL: %w", err)
	}

	if codeLength < 6 || codeLength > 12 {
		log.Printf("%s codeLength=%d out of bounds, using default=%d", logPrefix, codeLength, s.defaultCodeLength)
		codeLength = s.defaultCodeLength
	}

	if !isPublic && userID == nil {
		log.Printf("%s private URL without userID", logPrefix)
		return nil, fmt.Errorf("userID must be provided for private URLs")
	}

	if userID == nil {
		log.Printf("%s anonymous user path: trying reuse", logPrefix)
		if s.cache != nil {
			if sc, ok, _ := s.cache.Get(ctx, cache.KeyOriginalURL(parsed.Original)); ok {
				log.Printf("%s cache hit original->short_code %s", logPrefix, sc)
				if u, err := s.repo.GetURLByShortCode(sc); err == nil && u != nil {
					log.Printf("%s found URL in DB for short_code %s (id=%s)", logPrefix, sc, u.ID)
					return u, nil
				}
				log.Printf("%s cache had code but DB miss, continuing", logPrefix)
			} else {
				log.Printf("%s cache miss original->short_code", logPrefix)
			}
		}

		if existing, err := s.repo.GetURLByOriginalURL(parsed.Original); err == nil && existing != nil {
			log.Printf("%s found existing URL in DB (id=%s, short_code=%s)", logPrefix, existing.ID, existing.ShortCode)
			_ = s.setCaches(ctx, existing)
			return existing, nil
		} else if err != nil {
			log.Printf("%s repo.GetURLByOriginalURL error: %v", logPrefix, err)
		}
	}

	var url *model.URL

	for i := 0; i < s.maxRetries; i++ {
		shortcode, err := utils.GenerateCode(parsed.Original, codeLength)
		if err != nil {
			log.Printf("%s generate code failed on attempt %d: %v", logPrefix, i+1, err)
			return nil, fmt.Errorf("failed to generate shortcode: %w", err)
		}
		log.Printf("%s attempt %d generated shortcode=%s", logPrefix, i+1, shortcode)

		url = &model.URL{
			UserID:      userID,
			OriginalURL: parsed.Original,
			ShortCode:   shortcode,
			IsPublic:    isPublic,
			ClickCount:  0,
			CreatedAt:   time.Now().UTC(),
		}

		if err := s.repo.SaveURL(url); err != nil {
			log.Printf("%s SaveURL failed on attempt %d: %v", logPrefix, i+1, err)

			if errors.Is(err, repository.ErrUniqueViolation) || isUniqueViolation(err) {
				log.Printf("%s unique violation, retrying with a new code", logPrefix)
				continue
			}
			return nil, fmt.Errorf("failed to save URL (unknown reason)")
		}

		log.Printf("%s SaveURL success on attempt %d (id=%s, short_code=%s)", logPrefix, i+1, url.ID, url.ShortCode)
		break
	}

	if url == nil || url.ID == "" {
		log.Printf("%s url is nil or id is empty after %d attempts", logPrefix, s.maxRetries)
		return nil, fmt.Errorf("could not generate a unique shortcode after %d attempts", s.maxRetries)
	}

	if err := s.setCaches(ctx, url); err != nil {
		log.Printf("%s setCaches error: %v", logPrefix, err)
	}

	log.Printf("%s returning url id=%s short_code=%s", logPrefix, url.ID, url.ShortCode)
	return url, nil
}

func (s *urlService) GetURLByShortCode(ctx context.Context, code string) (*model.URL, error) {
	const logPrefix = "[GetURLByShortCode]"

	if s.cache != nil {
		if orig, ok, err := s.cache.Get(ctx, cache.KeyShortCode(code)); err == nil && ok {
			log.Printf("%s cache hit for code=%s, original URL=%s", logPrefix, code, orig)

			if u, err := s.repo.GetURLByShortCode(code); err == nil && u != nil {
				log.Printf("%s found full URL in DB for code=%s (id=%s)", logPrefix, code, u.ID)
				return u, nil
			}

			log.Printf("%s cache had original URL but DB lookup failed or not found for code=%s", logPrefix, code)
			return &model.URL{OriginalURL: orig, ShortCode: code}, nil
		} else if err != nil {
			log.Printf("%s cache lookup error for code=%s: %v", logPrefix, code, err)
		} else {
			log.Printf("%s cache miss for code=%s", logPrefix, code)
		}
	}

	u, err := s.repo.GetURLByShortCode(code)
	if err != nil {
		log.Printf("%s DB error for code=%s: %v", logPrefix, code, err)
		return nil, err
	}
	if u == nil {
		log.Printf("%s no URL found in DB for code=%s", logPrefix, code)
		return nil, nil
	}

	log.Printf("%s found in DB, caching and returning URL (id=%s, short_code=%s)", logPrefix, u.ID, u.ShortCode)
	if err := s.setCaches(ctx, u); err != nil {
		log.Printf("%s setCaches failed: %v", logPrefix, err)
	}
	return u, nil
}


func (s *urlService) GetUserUrls(user model.User) ([]model.URL, error) {
	log.Printf("[GetUserUrls] Fetching URLs for userID=%s", user.ID)
	return s.repo.GetUserUrls(user)
}

func (s *urlService) IncrementClickCount(shortcode string) error {
	log.Printf("[IncrementClickCount] Incrementing click count for short_code=%s", shortcode)
	return s.repo.IncrementClickCount(shortcode)
}

func (s *urlService) setCaches(ctx context.Context, url *model.URL) error {
	const logPrefix = "[setCaches]"

	if s.cache == nil {
		log.Printf("%s skipping cache — no cache implementation", logPrefix)
		return nil
	}

	if err := s.cache.Set(ctx, cache.KeyShortCode(url.ShortCode), url.OriginalURL, s.cacheTTL); err != nil {
		log.Printf("%s failed to cache short_code -> original: %v", logPrefix, err)
		return err
	}

	if url.UserID == nil {
		if err := s.cache.Set(ctx, cache.KeyOriginalURL(url.OriginalURL), url.ShortCode, s.cacheTTL); err != nil {
			log.Printf("%s failed to cache original -> short_code: %v", logPrefix, err)
			return err
		}
	}

	log.Printf("%s caching successful for short_code=%s", logPrefix, url.ShortCode)
	return nil
}


func isUniqueViolation(err error) bool {
	msg := strings.ToLower(err.Error())
	return strings.Contains(msg, "duplicate key") || strings.Contains(msg, "unique constraint")
}
