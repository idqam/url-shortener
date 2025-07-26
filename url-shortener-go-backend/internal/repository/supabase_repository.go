package repository

import (
	"fmt"

	"github.com/supabase-community/supabase-go"
)


type SupabaseRepository struct {
	Client *supabase.Client
}

func NewSupabaseRepository(apiURL string, apiKey string) (*SupabaseRepository, error) {
    client, err := supabase.NewClient(apiURL, apiKey, &supabase.ClientOptions{})
    if err != nil {
        return nil, fmt.Errorf("failed to initialize Supabase client: %w", err)
    }

    return &SupabaseRepository{Client: client}, nil
}
