package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func HomeHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.Render(http.StatusOK, "index", nil)
	}
}
