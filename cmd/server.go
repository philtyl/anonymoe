package cmd

import (
	"fmt"
	"log"
	"net/http"
	"path"

	"anonymoe/models"
	"anonymoe/pkg/context"
	"anonymoe/pkg/mail"
	"anonymoe/pkg/setting"
	"anonymoe/pkg/template"
	"anonymoe/routes"
	"github.com/go-macaron/session"
	"github.com/urfave/cli"
	"gopkg.in/macaron.v1"
)

var Server = cli.Command{
	Name:        "server",
	Usage:       "[start] [-p|-port 3000] [-c|-config custom/conf/app.ini]",
	Description: "Anonymoe web server starts all necessary components",
	Action:      StartServer,
	Flags: []cli.Flag{
		stringFlag("port, p", "3000", "Temporary port number to prevent conflict"),
		stringFlag("config, c", "custom/conf/app.ini", "Custom configuration file path"),
	},
}

func NewMacaron() *macaron.Macaron {
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
	m.Use(session.Sessioner(setting.SessionConfig))
	m.Use(context.Contexter())
	return m
}

func StartServer(c *cli.Context) error {
	go mail.NewSMTPServer()
	return RunWeb(c)
}

func RunWeb(c *cli.Context) error {
	routes.GlobalInit()

	models.LoadConfigs()
	if err := models.SetEngine(); err != nil {
		log.Fatalf("Internal error: SetEngine: %v", err)
	}

	m := NewMacaron()
	m.SetAutoHead(true)
	m.Get("/", routes.Home)
	m.Get("/inbox", routes.NewInbox)
	m.Get("/inbox/:user", routes.Inbox)
	m.NotFound(routes.Home)

	listenAddr := fmt.Sprintf("%s:%s", setting.HTTPAddr, setting.HTTPPort)
	log.Printf("Listen: %v://%s%s", setting.Protocol, listenAddr, setting.AppURL)

	var err error
	server := &http.Server{Addr: listenAddr, Handler: m}
	err = server.ListenAndServe()

	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}

	return nil
}
