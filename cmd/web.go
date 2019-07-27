package cmd

import (
	"github.com/urfave/cli"
)

var Web = cli.Command{
	Name:  "web",
	Usage: "Start web server",
	Description: `Anonymoe web server starts all necessary components`,
	Action: runWeb,
	Flags: []cli.Flag{
		stringFlag("port, p", "3000", "Temporary port number to prevent conflict"),
		stringFlag("config, c", "custom/conf/app.ini", "Custom configuration file path"),
	},
}

func runWeb(c *cli.Context) error {
	return nil
}