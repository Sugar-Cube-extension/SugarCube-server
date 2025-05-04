package api

import (
	"net/http"

	"github.com/MisterNorwood/SugarCube-Server/internal/database"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

var ApiClient *mongo.Database

// GET /api/coupons?site=<sitename>
func GetCouponsForPage(c echo.Context) error {
	site := c.QueryParam("site")
	if site == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Missing site parameter",
		})
	}

	coupons, err := database.GetSiteStruct(site, ApiClient)
	if err != nil {
		log.Fatal().
			Str("ip", c.RealIP()).
			Str("user_agent", c.Request().UserAgent()).
			Str("path", c.Request().URL.Path).
			Str("query_parm", site).
			Err(err).
			Msg("Error retriving data from database")

		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Database error",
		})

	}

	return c.JSON(http.StatusOK, coupons)
}

// POST /api/coupons?site=<sitename>
func AddCouponToSite(c echo.Context) error {
	site := c.QueryParam("site")
	if site == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Missing site parameter",
		})
	}

	var coupon database.CouponEntry
	if err := c.Bind(&coupon); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid JSON format",
		})
	}

	err := database.AddCouponToExistingSite(site, coupon, ApiClient)
	if err != nil {
		log.Error().
			Str("site", site).
			Str("ip", c.RealIP()).
			Err(err).
			Msg("Failed to insert coupon")
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	} else {
		log.Info().
			Str("site", site).
			Str("ip", c.RealIP()).
			Msg("Inserted coupon")
	}

	return c.JSON(http.StatusCreated, map[string]string{
		"status": "Coupon added",
	})
}

// POST /api/site?url=<sitename>
func RequestAddSite(c echo.Context) error {
	site := c.QueryParam("url")
	if site == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Missing site parameter",
		})
	}

	err := database.AddSite(site, ApiClient)
	if err != nil {
		log.Error().
			Str("site", site).
			Str("ip", c.RealIP()).
			Err(err).
			Msg("Failed to add site")
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	} else {
		log.Info().
			Str("site", site).
			Str("ip", c.RealIP()).
			Msg("Added site")
	}

	return c.JSON(http.StatusCreated, map[string]string{
		"status": "Site added",
	})
}
