package cmd

import (
	"github.com/urfave/cli"
)

var Mail = cli.Command{
	Name:        "mail",
	Usage:       "This command controls the SMTP Daemon",
	Description: `The SMTP Daemon receives new emails directed at domain and fires events`,
	Action:      runServ,
	Flags: []cli.Flag{
		stringFlag("config, c", "custom/conf/app.ini", "Custom configuration file path"),
	},
}

func runServ(c *cli.Context) error {
	return nil
}