package main

import (
	"log"
	"net"
	"net/http"
	"os"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/oschwald/geoip2-golang"
)

type Error struct {
	Message string `json:"message"`
	Info    string `json:"info"`
}

type ResponseV1 struct {
	IP        string          `json:"ip"`
	City      string          `json:"city,omitempty"`
	Country   string          `json:"country,omitempty"`
	Latitude  float64         `json:"latitude,omitempty"`
	Longitude float64         `json:"longitude,omitempty"`
	Timezone  string          `json:"timezone,omitempty"`
	Astral    *AstralResponse `json:"astral,omitempty"`
}

type CodeResponse struct {
	Name string `json:"name"`
	Code string `json:"code"`
}

type LocationResponse struct {
	City      string       `json:"city,omitempty"`
	Country   CodeResponse `json:"country"`
	Latitude  float64      `json:"latitude,omitempty"`
	Longitude float64      `json:"longitude,omitempty"`
	Timezone  string       `json:"timezone,omitempty"`
}

type ASResponse struct {
	Number       uint   `json:"number,omitempty"`
	Organization string `json:"organization,omitempty"`
}

type ResponseV2 struct {
	IP       string            `json:"ip"`
	AS       *ASResponse       `json:"as,omitempty"`
	Location *LocationResponse `json:"location,omitempty"`
	Astral   *AstralResponse   `json:"astral,omitempty"`
}

func getIp(c echo.Context) string {
	ipFrom := c.QueryParam("ip")
	if ipFrom == "" {
		ipFrom = c.Request().Header.Get("X-Forwarded-For")
	}
	if ipFrom == "" {
		ipFrom = c.Request().Header.Get("X-Real-IP")
	}
	if ipFrom == "" {
		ipFrom = c.Request().RemoteAddr
	}
	return ipFrom
}

func main() {
	// Open database
	pathCity, exist := os.LookupEnv("GEOIP_PATH_CITY")
	if !exist {
		pathCity = "/data/GeoLite2-City.mmdb"
	}
	dbCity, err := geoip2.Open(pathCity)
	if err != nil {
		log.Fatal(err)
	}
	defer dbCity.Close()

	pathASN, exist := os.LookupEnv("GEOIP_PATH_ASN")
	if !exist {
		pathASN = "/data/GeoLite2-ASN.mmdb"
	}
	dbASN, err := geoip2.Open(pathASN)
	if err != nil {
		log.Fatal(err)
	}
	defer dbASN.Close()

	// Create router
	e := echo.New()
	e.HideBanner = true

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowMethods: []string{http.MethodGet},
	}))

	e.GET("/healthz", func(c echo.Context) error {
		return c.HTML(http.StatusOK, "OK")
	})

	e.GET("/plain", func(c echo.Context) error {
		ipFrom := getIp(c)
		if ipFrom == "" {
			return c.JSON(http.StatusBadRequest, "Cannot get an IP (?ip, X-Forwarded-For, X-Real-IP or Remote address)")
		}

		ip := net.ParseIP(ipFrom)

		record, err := dbCity.City(ip)
		if err != nil || record == nil {
			return c.String(http.StatusNotFound, "Cannot found IP `"+ipFrom+"` in our database")
		}

		as, err := dbASN.ASN(ip)
		if err != nil {
			log.Println(err)
		}

		return c.String(http.StatusOK, ip.String()+" from "+record.City.Names["en"]+", "+record.Country.Names["en"]+", using "+as.AutonomousSystemOrganization)
	})

	e.GET("/v1", func(c echo.Context) error {
		ipFrom := getIp(c)
		if ipFrom == "" {
			return c.JSON(http.StatusBadRequest, Error{
				Message: "Bad IP",
				Info:    "Cannot get an IP (X-Forwarded-For, X-Real-IP or Remote address)",
			})
		}

		ip := net.ParseIP(ipFrom)

		record, err := dbCity.City(ip)
		if err != nil || record == nil {
			return c.JSON(http.StatusNotFound, Error{
				Message: "Not found",
				Info:    "Cannot found IP `" + ipFrom + "` in our database",
			})
		}
		return c.JSON(http.StatusOK, ResponseV1{
			IP:        ip.String(),
			City:      record.City.Names["en"],
			Country:   record.Country.Names["en"],
			Latitude:  record.Location.Latitude,
			Longitude: record.Location.Longitude,
			Timezone:  record.Location.TimeZone,
		})
	})

	e.GET("/v2", func(c echo.Context) error {
		ipFrom := getIp(c)
		if ipFrom == "" {
			return c.JSON(http.StatusBadRequest, Error{
				Message: "Bad IP",
				Info:    "Cannot get an IP (?ip, X-Forwarded-For, X-Real-IP or Remote address)",
			})
		}

		ip := net.ParseIP(ipFrom)

		as, err := dbASN.ASN(ip)
		if err != nil {
			log.Println(err)
		}

		record, err := dbCity.City(ip)
		if err != nil || record == nil {
			return c.JSON(http.StatusNotFound, Error{
				Message: "Not found",
				Info:    "Cannot found IP `" + ipFrom + "` in our database",
			})
		}
		return c.JSON(http.StatusOK, ResponseV2{
			IP: ip.String(),
			AS: &ASResponse{
				Number:       as.AutonomousSystemNumber,
				Organization: as.AutonomousSystemOrganization,
			},
			Location: &LocationResponse{
				City: record.City.Names["en"],
				Country: CodeResponse{
					Code: record.Country.IsoCode,
					Name: record.Country.Names["en"],
				},
				Latitude:  record.Location.Latitude,
				Longitude: record.Location.Longitude,
				Timezone:  record.Location.TimeZone,
			},
			Astral: getAstral(record.Location.Latitude, record.Location.Longitude),
		})
	})

	e.Logger.Fatal(e.Start(":1323"))
}
