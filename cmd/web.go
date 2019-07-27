package cmd

import (
	"path"

	"github.com/urfave/cli"
	"gopkg.in/macaron.v1"
)

var Web = cli.Command{
	Name:        "web",
	Usage:       "Start web server",
	Description: `Anonymoe web server starts all necessary components`,
	Action:      runWeb,
	Flags: []cli.Flag{
		stringFlag("port, p", "3000", "Temporary port number to prevent conflict"),
		stringFlag("config, c", "custom/conf/app.ini", "Custom configuration file path"),
	},
}

func newMacaron() *macaron.Macaron {
	m := macaron.New()
	m.Use(macaron.Logger())
	m.Use(macaron.Static(
		path.Join(setting.StaticRootPath, "public"),
	))

	funcMap := template.NewFuncMap()
	m.Use(macaron.Renderer(macaron.RenderOptions{
		Directory:  path.Join(setting.StaticRootPath, "templates"),
		Funcs:      funcMap,
		IndentJSON: macaron.Env != macaron.PROD,
	}))
}

func runWeb(c *cli.Context) error {
	return nil
}
