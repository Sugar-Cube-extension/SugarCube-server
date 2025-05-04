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

	//coupons, err := FindCouponsForSite(site)
	// if err != nil {
	//     return c.JSON(http.StatusInternalServerError, map[string]string{
	//         "error": "Database error",
	//     })
	// }
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
