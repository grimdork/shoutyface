package main

import (
	"fmt"
	"os"

	"github.com/Urethramancer/signor/opt"
)

// Options for the app.
var Options struct {
	opt.DefaultHelp
	Serve   CmdServe `command:"serve" help:"Start server."`
	Env     CmdEnv   `command:"env" help:"List available environment variables for the server."`
	Init    CmdInit  `command:"init" help:"Configure client settings."`
	User    CmdEnv   `command:"user" help:"User management."`
	Channel CmdEnv   `command:"channel" help:"Channel management."`
	Sub     CmdEnv   `command:"sub" help:"Subscription management."`
}

func main() {
	a := opt.Parse(&Options)
	if Options.Help {
		a.Usage()
		return
	}

	err := a.RunCommand(false)
	if err != nil {
		if err.Error() == opt.ErrorNoCommand {
			a.Usage()
			return
		}

		fmt.Printf("Error running command: %s\n", err.Error())
		os.Exit(2)
	}
}
