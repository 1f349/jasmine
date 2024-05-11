package main

import (
	"context"
	"encoding/json"
	"flag"
	"github.com/1f349/jasmine"
	"github.com/1f349/violet/utils"
	"github.com/charmbracelet/log"
	"github.com/google/subcommands"
	"github.com/mrmelon54/exit-reload"
	"os"
	"path/filepath"
)

type serveCmd struct {
	configPath string
	debugLog   bool
}

func (s *serveCmd) Name() string { return "serve" }

func (s *serveCmd) Synopsis() string { return "Serve calendar service" }

func (s *serveCmd) SetFlags(f *flag.FlagSet) {
	f.StringVar(&s.configPath, "conf", "", "/path/to/config.json : path to the config file")
	f.BoolVar(&s.debugLog, "debug", false, "enable debug logging")
}

func (s *serveCmd) Usage() string {
	return `serve [-conf <config file>]
  Serve calendar service using information from the config file
`
}

func (s *serveCmd) Execute(_ context.Context, _ *flag.FlagSet, _ ...any) subcommands.ExitStatus {
	if s.debugLog {
		jasmine.Logger.SetLevel(log.DebugLevel)
	}
	jasmine.Logger.Info("Starting...")

	if s.configPath == "" {
		jasmine.Logger.Error("Config flag is missing")
		return subcommands.ExitUsageError
	}

	openConf, err := os.Open(s.configPath)
	if err != nil {
		if os.IsNotExist(err) {
			jasmine.Logger.Error("Missing config file")
		} else {
			jasmine.Logger.Error("Open config file", "err", err)
		}
		return subcommands.ExitFailure
	}

	var config jasmine.Conf
	err = json.NewDecoder(openConf).Decode(&config)
	if err != nil {
		jasmine.Logger.Error("Invalid config file", "err", err)
		return subcommands.ExitFailure
	}

	configPathAbs, err := filepath.Abs(s.configPath)
	if err != nil {
		jasmine.Logger.Error("Failed to get absolute config path")
		return subcommands.ExitFailure
	}
	wd := filepath.Dir(configPathAbs)
	normalLoad(config, wd)
	return subcommands.ExitSuccess
}

func normalLoad(startUp jasmine.Conf, wd string) {
	srv := jasmine.NewHttpServer(startUp, wd)
	jasmine.Logger.Infof("Starting HTTP server on '%s'", srv.Addr)
	go utils.RunBackgroundHttp("HTTP", srv)

	exit_reload.ExitReload("Jasmine", func() {}, func() {
		// stop http server
		_ = srv.Close()
	})
}
