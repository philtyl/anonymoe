package cmd

import (
	"fmt"
	"net/http"
	"path"

	"github.com/go-macaron/session"
	"github.com/philtyl/anonymoe/models"
	"github.com/philtyl/anonymoe/pkg/context"
	"github.com/philtyl/anonymoe/pkg/mail"
	"github.com/philtyl/anonymoe/pkg/setting"
	"github.com/philtyl/anonymoe/pkg/template"
	"github.com/philtyl/anonymoe/routes"
	"github.com/urfave/cli"
	log "gopkg.in/clog.v1"
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
	SetupLogger("server.log")

	if err := setting.NewContext(); err != nil {
		log.Fatal(2, "Unable to initialize settings context: %v", err)
	}
	if err := models.SetEngine(); err != nil {
		log.Fatal(2, "Internal error: SetEngine: %v", err)
	}

	if c.Bool("start") {
		go mail.NewSMTPServer()
		return runWeb(c)
	} else if c.Bool("info") {
		if err := setting.NewContext(); err != nil {
			log.Fatal(2, "Unable to initialize settings context: %v", err)
		}

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
	log.Info("Listen: %v://%s%s", setting.Protocol, listenAddr, setting.AppURL)

	var err error
	server := &http.Server{Addr: listenAddr, Handler: m}
	err = server.ListenAndServe()

	if err != nil {
		log.Fatal(2, "Failed to start server: %v", err)
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
