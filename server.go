package main

import (
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/AlekSi/zabbix"
	"github.com/labstack/echo"
)

const (
	Version string = "0.1.0"
)

var (
	port           string = os.Getenv("PORT")
	zabbixUsername string = os.Getenv("ZABBIX_USERNAME")
	zabbixPassword string = os.Getenv("ZABBIX_PASSWORD")
	zabbixUrl      string = os.Getenv("ZABBIX_URL")
)

func main() {
	e := echo.New()

	e.GET("/", func(c echo.Context) error { return remoteAddress(c) })
	e.GET("/stop/:eventid", func(c echo.Context) error { return stopAlert(c) })
	e.GET("/version", func(c echo.Context) error { return version(c) })

	e.Logger.Fatal(e.Start(":" + port))
}

func remoteAddress(c echo.Context) error {
	req := c.Request()
	ip := strings.Split(req.RemoteAddr, ":")[0]
	return c.String(http.StatusOK, ip)
}

func shouldBeBlocked(userAgent string) bool {
	if userAgent == "Slackbot-LinkExpanding" {
		return true
	}
	return false
}

func stopAlert(c echo.Context) error {
	userAgent := c.Request().UserAgent()
	if shouldBeBlocked(userAgent) {
		return c.String(http.StatusForbidden, "Forbidden")
	}
	eventid := c.Param("eventid")
	api := zabbix.NewAPI(zabbixUrl)
	responseAuth, err := api.Call("user.login", zabbix.Params{
		"user":     zabbixUsername,
		"password": zabbixPassword,
	})
	if err != nil {
		return c.String(http.StatusBadRequest, "Failed to authorization.")
	}
	api.Auth = responseAuth.Result.(string)
	if _, err = api.Call("event.acknowledge", zabbix.Params{
		"eventids": eventid,
		"message":  "Alert-Stop",
	}); err != nil {
		log.Println(err)
		return c.String(http.StatusBadRequest, "Failed to acknowledge event.")
	}
	return c.String(http.StatusOK, "Alert-Stop")
}

func version(c echo.Context) error {
	return c.String(http.StatusOK, Version)
}
