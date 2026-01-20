package utils

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func Success(c echo.Context, message string, data interface{}) error {
    response := map[string]interface{}{
        "status":  "success",
        "message": message,
    }

    // hanya tampilkan data jika tidak nil
    if data != nil {
        response["data"] = data
    }

    return c.JSON(http.StatusOK, response)
}

func Error(c echo.Context, status int, message string) error {
	return c.JSON(status, map[string]interface{}{
		"status":  "error",
		"message": message,
	})
}
