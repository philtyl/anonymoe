package cmd

import (
	"log"

	"anonymoe/models"
	"anonymoe/pkg/setting"
	"github.com/urfave/cli"
)

var Install = cli.Command{
	Name:        "install",
	Usage:       "[-i|-init] [-m|-migrate]",
	Description: "Install necessary configuration and database models",
	Action:      InstallServer,
	Flags: []cli.Flag{
		boolFlag("init, i", "Initial database install"),
		boolFlag("migrate, m", "Migrate existing database"),
	},
}

func InstallServer(c *cli.Context) (err error) {
	setting.NewContext()
	models.LoadConfigs()
	if err = models.NewEngine(); err != nil {
		log.Fatalf("Fail to initialize ORM engine: %v", err)
	}
	return err
}
