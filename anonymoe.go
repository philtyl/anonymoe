package main

import (
	"os"

	"github.com/urfave/cli"

	"github.com/philtyl/anonymoe/cmd"
	"github.com/philtyl/anonymoe/pkg/setting"
)

const APP_VER = "0.0.1"

func init() {
	setting.AppVer = APP_VER
}

func main() {
	app := cli.NewApp()
	app.Name = "Anonymoe"
	app.Usage = "An open source anonymous password-less email client"
	app.Version = APP_VER
	app.Commands = []cli.Command{
		cmd.Web,
		cmd.Mail,
	}
	app.Run(os.Args)
}
