package main

import (
	"log"
	"net"
	"net/http"

	"github.com/labstack/echo"
	"github.com/oschwald/geoip2-golang"
)

type Error struct {
	Message string `json:"message"`
	Info    string `json:"info"`
}

type Response struct {
	IP        string  `json:"ip"`
	City      string  `json:"city,omitempty"`
	Country   string  `json:"country,omitempty"`
	Latitude  float64 `json:"latitude,omitempty"`
	Longitude float64 `json:"longitude,omitempty"`
	Timezone  string  `json:"timezone,omitempty"`
}

func main() {
	// Open database
	db, err := geoip2.Open("/data/GeoIP2-City.mmdb")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Create router
	e := echo.New()
	e.HideBanner = true

	e.GET("/healthz", func(c echo.Context) error {
		return c.HTML(http.StatusOK, "OK")
	})

	e.GET("/v1", func(c echo.Context) error {
		ipFrom := c.Request().Header.Get("X-Forwarded-For")
		if ipFrom == "" {
			ipFrom = c.Request().Header.Get("X-Real-IP")
		}
		if ipFrom == "" {
			ipFrom = c.Request().RemoteAddr
		}
		if ipFrom == "" {
			return c.JSON(http.StatusBadRequest, Error{
				Message: "Bad IP",
				Info:    "Cannot get an IP (X-Forwarded-For, X-Real-IP or Remote address)",
			})
		}

		ip := net.ParseIP(ipFrom)

		record, err := db.City(ip)
		if err != nil || record == nil {
			return c.JSON(http.StatusNotFound, Error{
				Message: "Not found",
				Info:    "Cannot found IP `" + ipFrom + "` in our database",
			})
		}

		return c.JSON(http.StatusOK, Response{
			IP:        ip.String(),
			City:      record.City.Names["en"],
			Country:   record.Country.Names["en"],
			Latitude:  record.Location.Latitude,
			Longitude: record.Location.Longitude,
			Timezone:  record.Location.TimeZone,
		})
	})

	e.Logger.Fatal(e.Start(":1323"))
}
