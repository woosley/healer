package src
import (
	"github.com/labstack/echo"
)


type Healer struct {
	echo.Context
	Opts     Opt
	//Status   State
	//Contents *Content
}

type Opt struct {
	Listen      int
	Is_master   bool
	Help        bool
	Master_addr string
	Expire      int
	Key         string
	Version     bool
	Debug       bool
	ConfigFile 		string
}

type Host struct {
	Ip string
	Name string
	HealthUrl string
}
