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
	var config []Host
	if !isUrl {
		config = loadConfigFromFile(options.Config)
	} else {
		config = loadConfigFromURL(options.Config)
	}
	config = []Host{{"url": "hello"}}
	var index int
	for idx, host := range config {
		go getHealthForHost(host, syncChan)
		index = idx
	}
	for i := 0; i <= index; i++ {
		h := <-syncChan
		data[h["key"]] = h
	}
	if !initial {
		time.Sleep(10 * time.Second)
	} else {
		initial = false
	}
	dataChan <- data
}

func getHealthForHost(h Host, syncChan chan Health) {
	fmt.Println("I am here")
	healthy := make(Health)
	healthy["status"] = "running"
	syncChan <- healthy
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
	e.Logger.Info(fmt.Sprintf("Starting gogate on %v", 1000))
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", 1000)))
}
