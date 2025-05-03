package middleware

import (
	goctx "context"
	"errors"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"net/http"
	"time"
)

func CheckIPBanList(next echo.HandlerFunc) echo.HandlerFunc {
	db := DBMiddlewareClient.Database("sugarcube_admin")
	return func(ctx echo.Context) error {
		ip := ctx.RealIP()
		filter := bson.M{"ip": ip}
		ctxGO, cancel := goctx.WithTimeout(goctx.Background(), 30*time.Second)
		defer cancel()

		found := db.Collection("ip_bans").FindOne(ctxGO, filter).Err()
		if found == nil {
			log.Warn().
				Str("ip", ctx.RealIP()).
				Str("user_agent", ctx.Request().UserAgent()).
				Str("path", ctx.Request().URL.Path).
				Msg("Blocked request due to IP being on a blacklist")
			return ctx.String(http.StatusForbidden, "Forbidden")

		} else if !errors.Is(found, mongo.ErrNoDocuments) {
			log.Fatal().
				Str("ip", ctx.RealIP()).
				Str("user_agent", ctx.Request().UserAgent()).
				Str("path", ctx.Request().URL.Path).
				Err(found).
				Msg("Error on checking request against IP list")

		}

		return next(ctx)

	}
}
