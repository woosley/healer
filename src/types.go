package src

import (
	"github.com/labstack/echo"
)

type Healer struct {
	echo.Context
	Opts Opt
	//Status   State
	//Contents *Content
	readChan ReadChan
}

type Opt struct {
	Listen  int    `short:"l" long:"listen" description:"set listen port"`
	Version bool   `short:"v" long:"version" description:"show current version"`
	Debug   bool   `short:"d" long:"debug" description:"debug mode"`
	Config  string `short:"c" long:"config" description:"config file location"`
}

type Host map[string]string

type Health map[string]string

// ReadChan passes data between web request goroutine and main goroutine
type ReadChan chan interface{}

// DataChan passes updated health status in health goroutine to main goroutine
type DataChan chan map[string]Health

// TreeChan passes all health tree to web request goroutine
type TreeChan struct {
	TChan chan map[string]Health
	Sync  chan bool
}

// KeyChan passes health status for a host to web request goroutine
type KeyChan struct {
	KChan chan Health
	Key   string
	Sync  chan bool
}
