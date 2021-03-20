package main

import (
	"fmt"
	"os"

	"github.com/Urethramancer/signor/opt"
)

var o struct {
	opt.DefaultHelp
	Serve   CmdServe   `command:"serve" help:"Start server."`
	Env     CmdEnv     `command:"env" help:"List available environment variables for the server."`
	Init    CmdInit    `command:"init" help:"Configure client settings."`
	User    CmdUser    `command:"user" aliases:"u" help:"User management."`
	Channel CmdChannel `command:"channel" aliases:"c" help:"Channel management."`
	Sub     CmdSubs    `command:"sub" aliases:"s" help:"Subscription management."`
}

func main() {
	a := opt.Parse(&o)
	if o.Help {
		a.Usage()
		return
	}

	err := a.RunCommand(false)
	if err != nil {
		if err == opt.ErrNoCommand {
			a.Usage()
			return
		}

		fmt.Printf("Error running command: %s\n", err.Error())
		os.Exit(2)
	}
}
