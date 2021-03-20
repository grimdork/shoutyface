package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/Urethramancer/signor/opt"
)

type CmdSubs struct {
	opt.DefaultHelp
	Add    CmdSubsAdd    `command:"add" help:"Add a channel."`
	Remove CmdSubsRemove `command:"remove" aliases:"rm" help:"Remove a channel."`
	List   CmdSubsList   `command:"list" aliases:"ls" help:"List channels."`
}

func (cmd *CmdSubs) Run(in []string) error {
	return opt.ErrUsage
}

//
// Add
//

type CmdSubsAdd struct {
	opt.DefaultHelp
	User    string `placeholder:"USERNAME" help:"Name of the user to subscribe to a channel."`
	Channel string `placeholder:"CHANNEL" help:"Name of the channel to subscribe the user to."`
}

func (cmd *CmdSubsAdd) Run(in []string) error {
	if cmd.Help || cmd.Channel == "" {
		return opt.ErrUsage
	}

	headers := make(map[string]string)
	headers["username"] = cmd.User
	headers["channel"] = cmd.Channel
	res, err := Request(http.MethodPost, "subscribe", headers)
	if err != nil {
		pr("Error subscribing %s to %s: %s", cmd.User, cmd.Channel, err.Error())
		return err
	}

	if res.StatusCode != http.StatusCreated {
		pr("Couldn't subscribe %s to %s: %s", cmd.User, cmd.Channel, http.StatusText(res.StatusCode))
		return nil
	}

	pr("Subscribed %s to %s.", cmd.User, cmd.Channel)
	return nil
}

//
// Remove
//

type CmdSubsRemove struct {
	opt.DefaultHelp
	User    string `placeholder:"USERNAME" help:"Name of the user to unsubscribe from a channel."`
	Channel string `placeholder:"CHANNEL" help:"Name of the channel to unsubscribe the user from."`
}

func (cmd *CmdSubsRemove) Run(in []string) error {
	if cmd.Help || cmd.Channel == "" {
		return opt.ErrUsage
	}

	headers := make(map[string]string)
	headers["username"] = cmd.User
	headers["channel"] = cmd.Channel
	res, err := Request(http.MethodDelete, "subscribe", headers)
	if err != nil {
		pr("Error unsubscribing %s from %s: %s", cmd.User, cmd.Channel, err.Error())
		return err
	}

	if res.StatusCode != http.StatusOK {
		pr("Error unsubscribing %s from %s: %s", cmd.User, cmd.Channel, http.StatusText(res.StatusCode))
		return nil
	}

	pr("Unsubscribed %s from %s.", cmd.User, cmd.Channel)
	return nil
}

//
// List
//

type CmdSubsList struct {
	opt.DefaultHelp
	User string `placeholder:"FILTER" help:"User to list subscriptions for. All subscriptions will be listed if not specified."`
}

func (cmd *CmdSubsList) Run(in []string) error {
	if cmd.Help {
		return opt.ErrUsage
	}

	var list []Sub
	headers := make(map[string]string)
	headers["username"] = cmd.User
	res, err := RequestJSON(http.MethodGet, "subscriptions", headers, &list)
	if err != nil {
		pr("Error listing subscriptions: %s (%s)\n", err.Error(), http.StatusText(res.StatusCode))
		return err
	}

	if res.StatusCode != http.StatusOK || len(list) == 0 {
		println("No channels.")
		return nil
	}

	buf := strings.Builder{}
	buf.WriteString("Username:\tChannel:\n")

	for _, s := range list {
		buf.WriteString(s.User)
		buf.WriteByte('\t')
		buf.WriteString(s.Channel)
		buf.WriteByte('\n')
	}
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 4, ' ', 0)
	fmt.Fprint(w, buf.String())
	w.Flush()
	return nil
}
