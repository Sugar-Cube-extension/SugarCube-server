package middleware

import "github.com/labstack/echo/v4"

func GlobalHeaderMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		headers := c.Response().Header()

		// Security headers
		headers.Set("Strict-Transport-Security", "max-age=63072000; includeSubDomains; preload")
		headers.Set("X-Content-Type-Options", "nosniff")
		headers.Set("X-Frame-Options", "DENY")
		headers.Set("X-XSS-Protection", "1; mode=block")
		headers.Set("Referrer-Policy", "no-referrer")
		headers.Set("Content-Security-Policy", "default-src 'self'")
		headers.Set("Permissions-Policy", "interest-cohort=()")
		headers.Set("Cache-Control", "no-store")

		// App metadata
		headers.Set("X-App-Name", "Sugarcube")
		headers.Set("X-App-Version", "1.0.0")
		headers.Set("X-API-Version", "v1")

		return next(c)
	}
}
