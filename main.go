package main

import (
	"fmt"
	"os"

	"github.com/Urethramancer/daemon"
	"github.com/Urethramancer/signor/opt"
)

// Options for the app.
var Options struct {
	opt.DefaultHelp
}

func main() {
	a := opt.Parse(&Options)
	if Options.Help {
		a.Usage()
		return
	}

	srv, err := NewServer()
	if err != nil {
		fmt.Printf("ERror starting server: %s\n", err.Error())
		os.Exit(2)
	}

	srv.Start()
	<-daemon.BreakChannel()
	srv.Stop()
}
