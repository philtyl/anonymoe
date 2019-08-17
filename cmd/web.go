package cmd

import (
	"fmt"
	"net/http"
	"path"

	"anonymoe/pkg/context"
	"anonymoe/pkg/setting"
	"anonymoe/pkg/template"
	"anonymoe/routes"
	"github.com/urfave/cli"
	log "gopkg.in/clog.v1"
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
	m.Use(context.Contexter())
	return m
}

func runWeb(c *cli.Context) error {
	m := newMacaron()
	m.SetAutoHead(true)
	m.Get("/", routes.Home)
	m.Get("/inbox/:username", routes.Inbox)
	m.NotFound(routes.Home)

	//// Flag for port number in case first time run conflict.
	//if c.IsSet("port") {
	//	setting.AppURL = strings.Replace(setting.AppURL, setting.HTTPPort, c.String("port"), 1)
	//	setting.HTTPPort = c.String("port")
	//}

	listenAddr := fmt.Sprintf("%s:%s", setting.HTTPAddr, setting.HTTPPort)
	log.Info("Listen: %v://%s%s", setting.Protocol, listenAddr, setting.AppURL)

	var err error
	server := &http.Server{Addr: listenAddr, Handler: m}
	err = server.ListenAndServe()

	if err != nil {
		log.Fatal(4, "Failed to start server: %v", err)
	}

	return nil
}
