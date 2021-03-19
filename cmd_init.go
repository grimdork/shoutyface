package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type CmdInit struct{}

func (cmd *CmdInit) Run(in []string) error {
	cfg, err := loadConfig()
	if err != nil {
		return err
	}

	r := bufio.NewReader(os.Stdin)
	fmt.Printf("Server URL [%s]: ", cfg.ServerURL)
	t, err := r.ReadString('\n')
	if err != nil {
		return err
	}

	cfg.ServerURL = strings.ToLower(strings.TrimSpace(t))
	fmt.Printf("Admin token [%s]: ", cfg.AdminToken)
	t, err = r.ReadString('\n')
	if err != nil {
		return err
	}

	cfg.AdminToken = strings.TrimSpace(t)
	return saveConfig(cfg)
}
