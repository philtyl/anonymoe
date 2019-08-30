package main

import (
	"os"

	"github.com/philtyl/anonymoe/cmd"
	"github.com/philtyl/anonymoe/pkg/setting"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "Anonymoe"
	app.Usage = "An open source anonymous password-less email client"
	app.Version = setting.Config.AppVer
	app.Commands = []cli.Command{
		cmd.Server,
		cmd.Install,
	}
	app.Run(os.Args)
}
