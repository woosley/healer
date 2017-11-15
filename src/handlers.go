//contain all routers
package src

import (
	"fmt"
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
	if len(data) == 0 {
		return cc.JSON(http.StatusNotFound,
			map[string]string{"code": fmt.Sprintf("%v", http.StatusNotFound), "message": "Not found"})
	}
	dt := cc.JSON(http.StatusOK, data)
	cha.Sync <- true
	return dt
}

func GetSelfHealth(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{"state": "running"})

}
