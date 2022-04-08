package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"github.com/ripexz/rip"
)

type IPTplArgs struct {
	GeoIP   *GeoIPResult
	Headers map[string]string
}

type GeoIPResult struct {
	IP          string  `json:"ip"`
	CountryCode string  `json:"country_code"`
	CountryName string  `json:"country_name"`
	RegionName  string  `json:"region_name"`
	City        string  `json:"city"`
	ZipCode     string  `json:"zip_code"`
	TimeZone    string  `json:"time_zone"`
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
}

type IPLookuper interface {
	Lookup(ip string) (*GeoIPResult, error)
}

func IPHandler(ipService IPLookuper) echo.HandlerFunc {
	return func(c echo.Context) error {
		geoIP, err := ipService.Lookup(rip.FromRequest(c.Request(), nil))
		if err != nil {
			err := errors.WithMessage(err, "failed to lookup ip geolocation")
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		if strings.Contains(c.Path(), "json") {
			return c.JSON(http.StatusOK, geoIP)
		}
		return c.Render(http.StatusOK, "ip", IPTplArgs{GeoIP: geoIP})
	}
}

type IPService struct {
	cache *Cache
}

func NewIPService(c *Cache) *IPService {
	return &IPService{
		cache: c,
	}
}

func (s IPService) Lookup(ip string) (*GeoIPResult, error) {
	cacheKey := fmt.Sprintf("ip-geolocation-%s", ip)
	if cached, hit := s.cache.Get(cacheKey); hit {
		return cached.(*GeoIPResult), nil
	}

	res, err := http.Get("https://freegeoip.app/json/" + ip)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var result *GeoIPResult
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, err
	}

	s.cache.Set(cacheKey, result)
	return result, nil
}
