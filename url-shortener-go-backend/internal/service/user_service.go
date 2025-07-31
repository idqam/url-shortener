package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"url-shortener-go-backend/internal/cache"
	"url-shortener-go-backend/internal/model"
	"url-shortener-go-backend/internal/repository"
)

type UserService interface {
	RegisterUser(email string) (*model.User, error)
	GetUserByEmail(email string) (*model.User, error)
}

type userServiceImpl struct {
	userRepo repository.UserRepository
	cache    cache.Cache
	ttl      time.Duration
}

func NewUserService(repo repository.UserRepository, c cache.Cache) UserService {
	return &userServiceImpl{
		userRepo: repo,
		cache:    c,
		ttl:      24 * time.Hour,
	}
}

func (s *userServiceImpl) RegisterUser(email string) (*model.User, error) {
	cleaned := strings.ToLower(strings.TrimSpace(email))

	existingUser, err := s.GetUserByEmail(cleaned)
	if err != nil && err != repository.ErrUserNotFound {
		return nil, err
	}
	if existingUser != nil {
		return nil, repository.ErrEmailInUse
	}

	if err := s.userRepo.CreateUser(cleaned); err != nil {
		return nil, err
	}

	user, err := s.userRepo.GetUserByEmail(cleaned)
	if err != nil {
		return nil, err
	}

	// Set cache
	if s.cache != nil {
		if err := s.setUserCache(user); err != nil {
			log.Printf("[UserService] Error setting user cache: %v", err)
		}
	}

	return user, nil
}

func (s *userServiceImpl) GetUserByEmail(email string) (*model.User, error) {
	cleaned := strings.ToLower(strings.TrimSpace(email))
	ctx := context.Background()

	cacheKey := fmt.Sprintf("user:%s", cleaned)

	if s.cache != nil {
		val, found, err := s.cache.Get(ctx, cacheKey)
		if err != nil {
			log.Printf("[UserService] Cache error for key %s: %v", cacheKey, err)
		} else if found {
			var user model.User
			if err := json.Unmarshal([]byte(val), &user); err == nil {
				log.Printf("[UserService] Cache HIT for user: %s", email)
				return &user, nil
			}
			log.Printf("[UserService] Cache unmarshal error: %v", err)
		}
	}

	log.Printf("[UserService] Cache MISS for user: %s", email)
	user, err := s.userRepo.GetUserByEmail(cleaned)
	if err != nil {
		return nil, err
	}

	if s.cache != nil {
		if err := s.setUserCache(user); err != nil {
			log.Printf("[UserService] Error caching user: %v", err)
		}
	}

	return user, nil
}

func (s *userServiceImpl) setUserCache(user *model.User) error {
	if user == nil {
		return nil
	}
	cacheKey := fmt.Sprintf("user:%s", strings.ToLower(user.Email))

	data, err := json.Marshal(user)
	if err != nil {
		return err
	}
	return s.cache.Set(context.Background(), cacheKey, string(data), s.ttl)
}
