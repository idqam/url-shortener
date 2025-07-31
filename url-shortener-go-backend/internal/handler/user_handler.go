package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"url-shortener-go-backend/internal/handler/dto"
	"url-shortener-go-backend/internal/repository"
	"url-shortener-go-backend/internal/service"
)

type UserHandler struct {
	svc service.UserService

}

func NewUserHandler(svc service.UserService) *UserHandler {
	return &UserHandler{svc: svc}
}

//POST /api/users
func (h *UserHandler) UserCreationHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req dto.GetUserByEmailRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		email := strings.TrimSpace(strings.ToLower(req.Email))
		user, err := h.svc.RegisterUser(email)
		if err != nil {
			status := http.StatusInternalServerError
			if errors.Is(err, repository.ErrEmailInvalid) {
				status = http.StatusBadRequest
			} else if errors.Is(err, repository.ErrEmailInUse) {
				status = http.StatusConflict
			}
			http.Error(w, err.Error(), status)
			return
		}

		resp := dto.GetUserByEmailResponse{
			ID:        user.ID,
			Email:     user.Email,
			CreatedAt: user.CreatedAt,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}
}
