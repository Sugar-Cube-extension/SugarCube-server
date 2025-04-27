package api

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

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
	coupons := []string{"SAVE10", "SAVE20"}

	return c.JSON(http.StatusOK, coupons)
}
