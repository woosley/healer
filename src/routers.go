//contain all routers
package src

import (
	"github.com/labstack/echo"
	"net/http"
)

func GetHealth(c echo.Context) error {
	cc := c.(*Healer)
	return cc.JSON(http.StatusOK, map[string]string{"me": "cool"})
}

func GetHealthFromHost(c echo.Context) error {
	cc := c.(*Healer)
	return cc.JSON(http.StatusOK, map[string]string{"me": "cool"})
}
