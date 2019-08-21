package main

import (
	"os"

	"anonymoe/cmd"
	"anonymoe/pkg/bindata"
	"anonymoe/pkg/setting"
	"github.com/urfave/cli"
)

var APP_VER = string(bindata.MustAsset("conf/VERSION"))

func init() {
	setting.AppVer = APP_VER
}

func main() {
	app := cli.NewApp()
	app.Name = "Anonymoe"
	app.Usage = "An open source anonymous password-less email client"
	app.Version = APP_VER
	app.Commands = []cli.Command{
		cmd.Server,
		cmd.Install,
	}
	app.Run(os.Args)
}
