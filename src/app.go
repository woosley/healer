package src

import (
	"fmt"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/labstack/gommon/log"
)

func GetHealthState(url string) map[string]string{
	health := make(map[string]string)
	health["state"] = "running"
	return health
}


func runHealthCheck(options Opt)[]Host{
	config := LoadConfig(options.ConfigFile)
	return config
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
	}else{
		e.Logger.SetLevel(log.INFO)
	}

	go runHealthCheck(options)
	e.HideBanner = true
	e.GET("/", GetHealth)
	e.Logger.Info(fmt.Sprintf("Starting gogate on %v", 1000))
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", 1000)))

}
