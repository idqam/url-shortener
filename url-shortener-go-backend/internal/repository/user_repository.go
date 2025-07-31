package repository

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"url-shortener-go-backend/internal/model"
	"url-shortener-go-backend/internal/utils"
)

type UserRepositoryImpl struct {
	*SupabaseRepository
}

func NewUserRepository(baseRepo *SupabaseRepository) UserRepository {
	return &UserRepositoryImpl{SupabaseRepository: baseRepo}
}

func (u *UserRepositoryImpl) CreateUser(email string) error {
	if !utils.ValidateEmail(email) {
		return ErrEmailInvalid
	}

	data := map[string]interface{}{
		"email": strings.TrimSpace(strings.ToLower(email)),
	}

	resp, _, err := u.Client.From("users").Insert(data, false, "", "", "").Execute()
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
			return ErrEmailInUse
		}
		return fmt.Errorf("%w: %v", ErrAccountCreation, err)
	}

	var inserted []model.User
	if err := json.Unmarshal(resp, &inserted); err != nil || len(inserted) == 0 {
		log.Printf("[CreateUser] Failed to decode insert response: %v", err)

		var errResp struct {
			Message string `json:"message"`
			Code    string `json:"code"`
			Details string `json:"details"`
			Hint    string `json:"hint"`
		}
		if json.Unmarshal(resp, &errResp) == nil && errResp.Message != "" {
			return fmt.Errorf("%w: %s (code: %s, details: %s, hint: %s)",
				ErrAccountCreation, errResp.Message, errResp.Code, errResp.Details, errResp.Hint)
		}

		return ErrAccountCreation
	}

	return nil
}
func (u *UserRepositoryImpl) GetUserByEmail(email string) (*model.User, error) {
	if !utils.ValidateEmail(email) {
		log.Printf("[GetUserByEmail] Invalid email: %s", email)
		return nil, ErrEmailInvalid
	}

	resp, _, err := u.Client.
		From("users").
		Select("*", "exact", false).
		Eq("email", email).
		Single().
		Execute()
	if err != nil {
		if strings.Contains(err.Error(), "No rows") {
			log.Printf("[GetUserByEmail] No user found for email: %s", email)
			return nil, ErrUserNotFound
		}
		log.Printf("[GetUserByEmail] Supabase query failed for email %s: %v", email, err)
		return nil, fmt.Errorf("failed to retrieve user: %w", err)
	}

	var user model.User
	if err := json.Unmarshal(resp, &user); err != nil {
		log.Printf("[GetUserByEmail] Failed to decode Supabase response: %v", err)

		var errResp struct {
			Message string `json:"message"`
			Code    string `json:"code"`
			Details string `json:"details"`
			Hint    string `json:"hint"`
		}
		if json.Unmarshal(resp, &errResp) == nil && errResp.Message != "" {
			log.Printf("[GetUserByEmail] Supabase error response: %+v", errResp)
			return nil, fmt.Errorf("supabase error: %s (code: %s, details: %s, hint: %s)",
				errResp.Message, errResp.Code, errResp.Details, errResp.Hint)
		}

		return nil, fmt.Errorf("failed to parse user response: %w", err)
	}

	return &user, nil
}
