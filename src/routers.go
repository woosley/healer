//contain all routers
package src

import (
	"github.com/labstack/echo"
	"net/http"
)

func GetHealth(c echo.Context) error {
	var cha = TreeChan{
		TChan: make(chan map[string]Health),
		Sync:  make(chan bool),
	}
	cc := c.(*Healer)
	cc.readChan <- cha
	data := <-cha.TChan

	dt := cc.JSON(http.StatusOK, data)
	// tell channel I am done
	cha.Sync <- true
	return dt
}

func GetHealthFromHost(c echo.Context) error {
	cc := c.(*Healer)
	key := cc.Param("key")
	cha := KeyChan{
		Key:   key,
		KChan: make(chan Health),
		Sync:  make(chan bool),
	}

	cc.readChan <- cha
	data := <-cha.KChan
	dt := cc.JSON(http.StatusOK, data)
	cha.Sync <- true
	return dt
}
