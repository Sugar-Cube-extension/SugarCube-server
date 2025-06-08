package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

var HEADER string = "SC-Api-version"
var API_VER string = "v1"

func CheckUserAgent(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		header := ctx.Request().Header

		if strings.Compare(header.Get(HEADER), API_VER) == 0 {
			log.Warn().
				Str("ip", ctx.RealIP()).
				Str("header", HeaderToString(header)).
				Str("path", ctx.Request().URL.Path).
				Msg("Blocked request due to invalid header")

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
