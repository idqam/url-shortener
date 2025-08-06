package middleware

import (
	"net/http"

	"github.com/unrolled/secure"
)

func SecurityHeaders() func(http.Handler) http.Handler {
	secureMiddleware := secure.New(secure.Options{
		FrameDeny:             true,
		ContentTypeNosniff:    true,
		BrowserXssFilter:      true,
		ContentSecurityPolicy: "default-src 'self'; script-src 'self' 'unsafe-inline'; style-src 'self' 'unsafe-inline'",
		ReferrerPolicy:        "strict-origin-when-cross-origin",

		STSSeconds:           31536000,
		STSIncludeSubdomains: true,
		STSPreload:           true,

		CustomFrameOptionsValue: "DENY",

		IsDevelopment: true, // Set based on environment
	})

	return func(next http.Handler) http.Handler {
		return secureMiddleware.Handler(next)
	}
}
