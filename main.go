package main

import (
	"fmt"
	"github.com/jessevdk/go-flags"
	"github.com/woosley/healer/src"
	"os"
)

const version string = "0.1"

var options src.Opt

func main() {
	_, err := flags.Parse(&options)
	if err != nil {
		panic(err)
	}

	if options.Version {
		fmt.Println("healer version", version)
		os.Exit(0)
	}
	src.App(options)
}
