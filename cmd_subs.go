package main

import (
	"github.com/Urethramancer/signor/opt"
)

type CmdSubs struct {
	opt.DefaultHelp
	Add    CmdChannelAdd    `command:"add" help:"Add a channel."`
	Remove CmdChannelRemove `command:"remove" aliases:"rm" help:"Remove a channel."`
	List   CmdChannelList   `command:"list" aliases:"ls" help:"List channels."`
}

func (cmd *CmdSubs) Run(in []string) error {
	return opt.ErrUsage
}
