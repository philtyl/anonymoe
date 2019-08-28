package cmd

import (
	"os"
	"path"

	"github.com/philtyl/anonymoe/models"
	"github.com/philtyl/anonymoe/pkg/bindata"
	"github.com/philtyl/anonymoe/pkg/setting"
	"github.com/urfave/cli"
	log "gopkg.in/clog.v1"
)

var Install = cli.Command{
	Name:   "install",
	Usage:  "Install or upgrade base system configuration and database",
	Action: InstallServer,
	Flags: []cli.Flag{
		boolFlag("init, i", "Initial database/config install"),
		boolFlag("migrate, m", "Migrate existing database"),

		boolFlag("force", "Force initial database/config install despite existing files"),
	},
}

func InstallServer(c *cli.Context) (err error) {
	SetupLogger("install.log")

	if c.Bool("init") {
		installPath := setting.InstallDir()
		installConfigFile := path.Join(installPath, "app.ini")
		if fileExists(installConfigFile) {
			if !c.Bool("force") {
				log.Fatal(2, "Configuration file [%s] already exists, please run with --force flag to overwrite", installConfigFile)
			}
		}
		if err = copyConf("conf/app.ini", installConfigFile); err != nil {
			log.Fatal(2, "Failed to initialize 'app.ini' to '%s': %v", installConfigFile, err)
		}
		if err = setting.NewContext(); err != nil {
			log.Fatal(2, "Failed to initialize settings engine: %v", err)
		}
		if err = models.NewEngine(); err != nil {
			log.Fatal(2, "Failed to initialize ORM engine: %v", err)
		}
	} else if c.Bool("migrate") {
		if err = setting.NewContext(); err != nil {
			log.Fatal(2, "Failed to initialize settings engine: %v", err)
		}
		if err = models.NewEngine(); err != nil {
			log.Fatal(2, "Failed to initialize ORM engine: %v", err)
		}
	}

	return err
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func copyConf(asset, dest string) (err error) {
	if err = os.MkdirAll(path.Dir(dest), 0755); err != nil {
		return err
	}

	if fileExists(dest) {
		if err = os.Remove(dest); err != nil {
			return err
		}
	}

	destination, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer destination.Close()
	if _, err := destination.Write(bindata.MustAsset(asset)); err != nil {
		return err
	}
	return err
}
