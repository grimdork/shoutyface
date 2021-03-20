package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/Urethramancer/signor/opt"
)

type CmdChannel struct {
	opt.DefaultHelp
	Add    CmdChannelAdd    `command:"add" help:"Add a channel."`
	Remove CmdChannelRemove `command:"remove" aliases:"rm" help:"Remove a channel."`
	List   CmdChannelList   `command:"list" aliases:"ls" help:"List channels."`
}

func (cmd *CmdChannel) Run(in []string) error {
	return opt.ErrUsage
}

//
// Add
//

// CmdUserAdd options.
type CmdChannelAdd struct {
	opt.DefaultHelp
	Name string `placeholder:"NAME" help:"Name of the new channel."`
}

// Run add.
func (cmd *CmdChannelAdd) Run(in []string) error {
	if cmd.Help || cmd.Name == "" {
		return opt.ErrUsage
	}

	headers := map[string]string{
		"channel": cmd.Name,
	}

	res, err := Request(http.MethodPost, "channel", headers)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error adding channel: %s\n", err.Error())
		return err
	}

	if res.StatusCode == http.StatusNotModified {
		pr("Channel '%s' already exists.", cmd.Name)
		return nil
	}

	if res.StatusCode != http.StatusOK {
		fmt.Fprintf(os.Stderr, "Error adding channel: %s\n", http.StatusText(res.StatusCode))
		return nil
	}

	println("Channel created.")
	return nil
}

//
// Remove
//

// CmdChannelRemove options.
type CmdChannelRemove struct {
	opt.DefaultHelp
	Name string `placeholder:"NAME" help:"Name of the channel to be removed."`
}

// Run remove.
func (cmd *CmdChannelRemove) Run(in []string) error {
	if cmd.Help {
		return opt.ErrUsage
	}

	headers := map[string]string{
		"channel": cmd.Name,
	}

	res, err := Request(http.MethodDelete, "channel", headers)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error removing channel: %s\n", err.Error())
		return err
	}

	switch res.StatusCode {
	case http.StatusOK:
		println("Channel removed.")
	case http.StatusNotFound:
		pr("Channel '%s' not found.", cmd.Name)
	default:
		pr("Unknown error removing '%s': %s", cmd.Name, http.StatusText(res.StatusCode))
	}

	return nil
}

//
// List
//

// CmdChannelList options.
type CmdChannelList struct{}

// Run list.
func (cmd *CmdChannelList) Run(in []string) error {
	var list []Channel
	headers := make(map[string]string)
	res, err := RequestJSON(http.MethodGet, "channels", headers, &list)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error listing channels: %s\nStatus: %d\n", err.Error(), res.StatusCode)
		return err
	}

	if res.StatusCode != http.StatusOK || len(list) == 0 {
		println("No channels.")
		return nil
	}

	println("Channel:")
	for _, c := range list {
		println(c.Name)
	}
	return nil
}
