package main

import (
	"github.com/Urethramancer/signor/opt"
)

type CmdSubs struct {
	opt.DefaultHelp
}

func (cmd *CmdSubs) Run(in []string) error {
	if cmd.Help {
		return opt.ErrUsage
	}

	return nil
}
