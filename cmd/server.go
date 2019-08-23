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
	Name:   "server",
	Usage:  "Control or query information from server",
	Action: StartServer,
	Flags: []cli.Flag{
		boolFlag("info", "Dump settings information for webserver"),
		boolFlag("start", "Start webserver and necessary auxiliary components"),
	},
}

func StartServer(c *cli.Context) error {
	if err := setting.NewContext(); err != nil {
		log.Fatalf("Unable to initialize settings context: %v", err)
	}
	if err := models.SetEngine(); err != nil {
		log.Fatalf("Internal error: SetEngine: %v", err)
	}

	if c.Bool("start") {
		go mail.NewSMTPServer()
		return runWeb(c)
	} else if c.Bool("info") {
		if err := setting.NewContext(); err != nil {
			log.Fatalf("Unable to initialize settings context: %v", err)
		}

	} else {
		log.Fatalf("")
	}
	return nil
}

func runWeb(c *cli.Context) error {
	m := newMacaron()
	m.SetAutoHead(true)
	m.Get("/", routes.Home)
	m.Get("/inbox", routes.NewInbox)
	m.Get("/inbox/:user", routes.Inbox)
	m.Get("/inbox/:user/node", routes.InboxNode)
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
	m.Use(session.Sessioner(setting.SessionConfig))
	m.Use(context.Contexter())
	return m
}
