package middleware

import (
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"net/http"
	"strings"
)

var USER_AGENT string = "SugarCube/"

func CheckUserAgent(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		ua := ctx.Request().UserAgent()
		if !strings.HasPrefix(ua, USER_AGENT) {
			log.Warn().
				Str("ip", ctx.RealIP()).
				Str("user_agent", ua).
				Str("path", ctx.Request().URL.Path).
				Msg("Blocked request due to invalid User-Agent")

			return ctx.String(http.StatusForbidden, "Forbidden")

		}
		return next(ctx)
	}
}
