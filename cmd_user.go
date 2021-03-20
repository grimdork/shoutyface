package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/Urethramancer/signor/opt"
)

// CmdUser options.
type CmdUser struct {
	opt.DefaultHelp
	Add    CmdUserAdd    `command:"add" help:"Add a user."`
	Remove CmdUserRemove `command:"remove" aliases:"rm" help:"Remove a user."`
	List   CmdUserList   `command:"list" aliases:"ls" help:"List users."`
}

// Run user help.
func (cmd *CmdUser) Run(in []string) error {
	return opt.ErrUsage
}

//
// Add
//

// CmdUserAdd options.
type CmdUserAdd struct {
	opt.DefaultHelp
	Name  string `placeholder:"NAME" help:"Username of the new user."`
	Email string `placeholder:"EMAIL" help:"E-mail address ofr alerts to the user."`
}

// Run add.
func (cmd *CmdUserAdd) Run(in []string) error {
	if cmd.Help || cmd.Email == "" {
		return opt.ErrUsage
	}

	headers := map[string]string{
		"username": cmd.Name,
		"email":    cmd.Email,
	}

	res, err := Request(http.MethodPost, "user", headers)
	if err != nil {
		pr("Error adding user: %s\n", err.Error())
		return err
	}

	if res.StatusCode != http.StatusCreated {
		pr("Error adding user: %s\n", http.StatusText(res.StatusCode))
		return nil
	}

	println("User created.")
	return nil
}

//
// Remove
//

// CmdUserRemove options.
type CmdUserRemove struct {
	opt.DefaultHelp
	Name string `placeholder:"NAME" help:"Username of the user to be removed."`
}

// Run remove.
func (cmd *CmdUserRemove) Run(in []string) error {
	if cmd.Help {
		return opt.ErrUsage
	}

	headers := map[string]string{
		"username": cmd.Name,
	}

	res, err := Request(http.MethodDelete, "user", headers)
	if err != nil {
		pr("Error removing user: %s\n", err.Error())
		return err
	}

	if res.StatusCode != http.StatusOK {
		pr("Error removing user: %s\n", http.StatusText(res.StatusCode))
		return nil
	}

	println("User removed.")
	return nil
}

//
// List
//

// CmdUserList options.
type CmdUserList struct {
	opt.DefaultHelp
}

// Run list.
func (cmd *CmdUserList) Run(in []string) error {
	if cmd.Help {
		return opt.ErrUsage
	}

	var list []User
	headers := make(map[string]string)
	res, err := RequestJSON(http.MethodGet, "users", headers, &list)
	if err != nil {
		pr("Error fetching data: %s\nStatus: %d\n", err.Error(), res.StatusCode)
		return err
	}

	if res.StatusCode != http.StatusOK || len(list) == 0 {
		println("No users.")
		return nil
	}

	buf := strings.Builder{}
	buf.WriteString("Username:\tE-mail:\n")
	for _, u := range list {
		buf.WriteString(u.Name)
		buf.WriteByte('\t')
		buf.WriteString(u.Email)
		buf.WriteByte('\n')
	}
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 4, ' ', 0)
	fmt.Fprint(w, buf.String())
	w.Flush()
	return nil
}
