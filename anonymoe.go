package main

import (
	"os"

	"anonymoe/cmd"
	"anonymoe/pkg/setting"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "Anonymoe"
	app.Usage = "An open source anonymous password-less email client"
	app.Version = setting.AppVer
	app.Commands = []cli.Command{
		cmd.Server,
		cmd.Install,
	}
	app.Run(os.Args)
}
