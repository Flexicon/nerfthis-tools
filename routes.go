package main

import (
	"time"

	"github.com/labstack/echo/v4"
)

func SetupRoutes(g *echo.Group) {
	g.GET("", HomeHandler())

	ipService := NewIPService(AppCache)
	g.GET("/ip", IPHandler(ipService), NoCacheMiddleware())
	g.GET("/ip.json", IPHandler(ipService), NoCacheMiddleware())
}

// NoCacheMiddleware prevents routers and clients from caching the response of a handler.
func NoCacheMiddleware() echo.MiddlewareFunc {
	noCacheHeaders := map[string]string{
		"Expires":         time.Unix(0, 0).Format(time.RFC1123),
		"Cache-Control":   "no-cache, no-store, must-revalidate, max-age=0",
		"Pragma":          "no-cache",
		"X-Accel-Expires": "0",
	}
	etagHeaders := []string{
		"ETag",
		"If-Modified-Since",
		"If-Match",
		"If-None-Match",
		"If-Range",
		"If-Unmodified-Since",
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			req := c.Request()
			// Delete any ETag headers that may have been set
			for _, v := range etagHeaders {
				if req.Header.Get(v) != "" {
					req.Header.Del(v)
				}
			}

			// Set NoCache headers
			res := c.Response()
			for k, v := range noCacheHeaders {
				res.Header().Set(k, v)
			}

			return next(c)
		}
	}
}
