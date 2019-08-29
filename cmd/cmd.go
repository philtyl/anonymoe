package cmd

import (
	"path/filepath"

	"github.com/philtyl/anonymoe/pkg/setting"
	"github.com/urfave/cli"
	log "gopkg.in/clog.v1"
)

func SetupLogger(logName string) {
	level := log.TRACE
	if setting.ProdMode {
		level = log.ERROR
	}
	err := log.New(log.FILE, log.FileConfig{
		Level:    level,
		Filename: filepath.Join(setting.InstallDir(), "logs", logName),
		FileRotationConfig: log.FileRotationConfig{
			Rotate:  true,
			Daily:   true,
			MaxDays: 3,
		},
	})
	if err != nil {
		log.Warn("Unable to start file logger, defaulting to STDOUT: %v", err)
	} else {
		log.Delete(log.CONSOLE) // Remove primary logger
	}
}

func stringFlag(name, value, usage string) cli.StringFlag {
	return cli.StringFlag{
		Name:  name,
		Value: value,
		Usage: usage,
	}
}

func boolFlag(name, usage string) cli.BoolFlag {
	return cli.BoolFlag{
		Name:  name,
		Usage: usage,
	}
}
