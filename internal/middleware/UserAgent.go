package middleware

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

var USER_AGENT string = "SugarCube/" // For now we wont enforce API versions, dont care

func CheckUserAgent(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		ua := ctx.Request().UserAgent()
		if !strings.HasPrefix(ua, USER_AGENT) {
			log.Warn().
				Str("ip", ctx.RealIP()).
				Str("user_agent", ua).
				Str("path", ctx.Request().URL.Path).
				Msg("Blocked request due to invalid User-Agent")

			return next(ctx)

		}
		return ctx.String(http.StatusForbidden, "Forbidden")
	}
}
