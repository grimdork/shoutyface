package main

import (
	"github.com/Urethramancer/signor/opt"
)

type CmdChannel struct {
	opt.DefaultHelp
}

func (cmd *CmdChannel) Run(in []string) error {
	if cmd.Help {
		return opt.ErrUsage
	}

	return nil
}
