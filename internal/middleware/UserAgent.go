package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

var HEADER = "SC-Api-version"
var API_VER = "v1"

func CheckUserAgent(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		version := strings.TrimSpace(ctx.Request().Header.Get(HEADER))

		if version != API_VER {
			log.Warn().
				Str("ip", ctx.RealIP()).
				Str("received_version", version).
				Str("expected_version", API_VER).
				Str("header_dump", HeaderToString(ctx.Request().Header)).
				Str("path", ctx.Request().URL.Path).
				Msg("Blocked request due to invalid API version header")

			return ctx.String(http.StatusForbidden, "Forbidden")
		}
		return next(ctx)
	}
}

func HeaderToString(header http.Header) string {
	var b strings.Builder
	for k, v := range header {
		for _, val := range v {
			b.WriteString(fmt.Sprintf("%s: %s\n", k, val))
		}
	}
	return b.String()
}
