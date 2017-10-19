package src

import (
	"fmt"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/labstack/gommon/log"
	"net/url"
	"time"
)

var cha chan bool = make(chan bool)
var initial bool = true

func looksLikeUrl(s string) bool {
	if _, err := url.ParseRequestURI(s); err != nil {
		return false
	} else {
		return true
	}
}

func runHealthCheck(options Opt, dataChan DataChan) {

	isUrl := looksLikeUrl(options.Config)

	data := make(map[string]Health)
	syncChan := make(chan Health)
	var config map[string]Health
	if !isUrl {
		config = loadConfigFromFile(options.Config)
	} else {
		config = loadConfigFromURL(options.Config)
	}
	var index int
	for _, host := range config {
		go getHealthForHost(host, syncChan)
		index += 1
	}
	for i := 0; i < index; i++ {
		h := <-syncChan
		data[guessKey(h)] = h
	}
	if !initial {
		time.Sleep(10 * time.Second)
	} else {
		initial = false
	}
	dataChan <- data
}

func getHealthForHost(h Health, syncChan chan Health) {
	url := h["health_url"]
	state, reason, code := getHealthFromURL(url)
	h["state"] = state
	h["reason"] = reason
	h["code"] = code
	syncChan <- h
}

func run(options Opt, readChan ReadChan, dataChan DataChan) {
	var health map[string]Health
	//trigger dataChan
	go runHealthCheck(options, dataChan)
	for {
		select {
		// take health from health check
		case h := <-dataChan:
			//re run health check
			health = h
			go runHealthCheck(options, dataChan)
		// read channel from http request
		case cha := <-readChan:
			// put data into channel
			switch t := cha.(type) {
			case KeyChan:
				doKeyChan(t, health)
			case TreeChan:
				//send to TreeChan
				t.TChan <- health
				<-t.Sync
			}
		}
	}
}

func App(options Opt) {
	e := echo.New()

	readChan := make(ReadChan)
	dataChan := make(DataChan)
	e.Use(middleware.Logger())
	e.Use(func(h echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			healer := &Healer{c, options, readChan}
			return h(healer)
		}
	})

	if options.Debug {
		e.Logger.SetLevel(log.DEBUG)
	} else {
		e.Logger.SetLevel(log.INFO)
	}

	go run(options, readChan, dataChan)
	e.HideBanner = true
	e.GET("/", GetHealth)
	e.GET("/:key", GetHealthFromHost)
	e.Logger.Info(fmt.Sprintf("Starting healer on %v", options.Listen))
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", options.Listen)))
}
