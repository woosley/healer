package src

import (
	"fmt"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/labstack/gommon/log"
	"net/url"
	"time"
)

func getHealthState(h Host, cha chan map[string]string) {
	health := make(map[string]string)
	health["state"] = "running"
	cha <- health
}

func looksLikeUrl(s string) bool {
	if _, err := url.ParseRequestURI(s); err != nil {
		return false
	} else {
		return true
	}
}

func runHealthCheck(options Opt) {

	isUrl := looksLikeUrl(options.Config)
	var config []Host
	if !isUrl {
		config = loadConfigFromFile(options.Config)
	}
	for {
		if isUrl {
			config = loadConfigFromURL(options.Config)
		}
		GetHealthFor(config)
		time.Sleep(5000 * time.Millisecond)
	}
}

func GetHealthFor(c []Host) interface{} {
	cha := make(chan map[string]string)
	for _, v := range c {
		go getHealthState(v, cha)
	}
	return nil
}

func App(options Opt) {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(func(h echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			healer := &Healer{c, options}
			return h(healer)
		}
	})

	if options.Debug {
		e.Logger.SetLevel(log.DEBUG)
	} else {
		e.Logger.SetLevel(log.INFO)
	}

	go runHealthCheck(options)
	e.HideBanner = true
	e.GET("/", GetHealth)
	e.Logger.Info(fmt.Sprintf("Starting gogate on %v", 1000))
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", 1000)))
}
