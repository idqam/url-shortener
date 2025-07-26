package utils

import (
	"fmt"
	"net/url"
)

type Url struct {
	Original string
	Parsed   *url.URL
}

func ValidateUrl(rawUrl string) (*Url, error) {
	parsedUrl, err := url.ParseRequestURI(rawUrl)
	if err != nil {
		return nil, fmt.Errorf("parse error for url %s: %w", rawUrl, err)
	}
	return &Url{Original: rawUrl, Parsed: parsedUrl}, nil
}
