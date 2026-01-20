package middleware

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func RoleMiddleware(role string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {

		return func(c echo.Context) error {
			userRole := c.Get("role")

			if userRole == nil {
				return c.JSON(http.StatusForbidden, map[string]string{
					"error": "role not found",
				})
			}

			if userRole.(string) != role {
				return c.JSON(http.StatusForbidden, map[string]string{
					"error": "access denied",
				})
			}

			return next(c)
		}
	}
}
