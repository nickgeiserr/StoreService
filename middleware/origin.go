package middleware

import "github.com/labstack/echo/v4"

func RequestOrigin() (echo.MiddlewareFunc, error) {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			return nil
		}
	}, nil
}
