package main

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const (
	TestIP       = "49.129.61.133"
	ExpectedJSON = "{\"ip\":\"49.129.61.133\",\"country_code\":\"de-DE\",\"country_name\":\"Germany\",\"region_name\":\"Idk\",\"city\":\"Stuttgart\",\"zip_code\":\"70173\",\"time_zone\":\"GMT+1\",\"latitude\":48.78,\"longitude\":9.18}"
)

var (
	TestGeoIP = GeoIPResult{
		Query:       TestIP,
		CountryCode: "de-DE",
		Country:     "Germany",
		RegionName:  "Idk",
		City:        "Stuttgart",
		Zip:         "70173",
		Timezone:    "GMT+1",
		Lat:         48.78,
		Lon:         9.18,
	}
)

func TestIPHandler(t *testing.T) {
	e := echo.New()

	t.Run("when requesting JSON", func(t *testing.T) {
		t.Run("should handle any ip lookup errors", func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			req.RemoteAddr = TestIP
			rec := httptest.NewRecorder()

			c := e.NewContext(req, rec)
			c.SetPath("/ip.json")

			expectedErr := errors.New("something went wrong")

			ipServiceMock := new(TestIPService)
			ipServiceMock.On("Lookup", TestIP).Return(nil, expectedErr)

			handler := IPHandler(ipServiceMock)

			if err := handler(c); assert.Error(t, err) {
				assert.Contains(t, err.Error(), expectedErr.Error())
			}
		})

		t.Run("should return expected JSON", func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			req.RemoteAddr = TestIP
			rec := httptest.NewRecorder()

			c := e.NewContext(req, rec)
			c.SetPath("/ip.json")

			ipServiceMock := new(TestIPService)
			ipServiceMock.On("Lookup", TestIP).Return(&TestGeoIP, nil)

			handler := IPHandler(ipServiceMock)

			if assert.NoError(t, handler(c)) {
				assert.Equal(t, http.StatusOK, rec.Code)
				assert.Equal(t, ExpectedJSON, strings.TrimSpace(rec.Body.String()))
			}
		})
	})
}

type TestIPService struct {
	mock.Mock
}

func (m *TestIPService) Lookup(ip string) (*GeoIPResult, error) {
	args := m.Called(ip)

	if v := args.Get(0); v != nil {
		return args.Get(0).(*GeoIPResult), args.Error(1)
	}
	return nil, args.Error(1)
}
