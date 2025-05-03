package middleware

import (
	"time"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

func ZeroLogMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		start := time.Now()
		err := next(c)
		log.Info().
			Str("method", c.Request().Method).
			Str("url", c.Request().URL.String()).
			Int("status", c.Response().Status).
			Dur("duration", time.Since(start)).
			Msg("Request processed")
		return err
	}
}
