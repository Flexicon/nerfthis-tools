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
	GeoIP *GeoIPResult
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

	res, err := http.Get("http://ip-api.com/json/" + ip)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var lookup *GeoIPLookupResponse
	if err := json.NewDecoder(res.Body).Decode(&lookup); err != nil {
		return nil, err
	}

	result := lookup.ToResult()
	s.cache.Set(cacheKey, &result)

	return &result, nil
}

type GeoIPLookupResponse struct {
	Query       string  `json:"query"`
	Country     string  `json:"country"`
	CountryCode string  `json:"countryCode"`
	RegionName  string  `json:"regionName"`
	City        string  `json:"city"`
	Zip         string  `json:"zip"`
	Timezone    string  `json:"timezone"`
	Lat         float64 `json:"lat"`
	Lon         float64 `json:"lon"`
}

func (r GeoIPLookupResponse) ToResult() GeoIPResult {
	return mapGeoIPResponseToResult(&r)
}

func mapGeoIPResponseToResult(result *GeoIPLookupResponse) GeoIPResult {
	return GeoIPResult{
		IP:          result.Query,
		CountryCode: result.CountryCode,
		CountryName: result.Country,
		RegionName:  result.RegionName,
		City:        result.City,
		ZipCode:     result.Zip,
		TimeZone:    result.Timezone,
		Latitude:    result.Lat,
		Longitude:   result.Lon,
	}
}
