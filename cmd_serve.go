package main

import (
	"github.com/Urethramancer/daemon"
)

type CmdServe struct{}

func (cmd *CmdServe) Run(in []string) error {
	srv, err := NewServer()
	if err != nil {
		return err
	}

	srv.Start()
	<-daemon.BreakChannel()
	srv.Stop()
	return nil
}
