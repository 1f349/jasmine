package main

import (
	"context"
	"encoding/json"
	"flag"
	"github.com/1f349/jasmine"
	"github.com/1f349/violet/utils"
	"github.com/google/subcommands"
	"github.com/mrmelon54/exit-reload"
	"log"
	"os"
	"path/filepath"
)

type serveCmd struct{ configPath string }

func (s *serveCmd) Name() string { return "serve" }

func (s *serveCmd) Synopsis() string { return "Serve calendar service" }

func (s *serveCmd) SetFlags(f *flag.FlagSet) {
	f.StringVar(&s.configPath, "conf", "", "/path/to/config.json : path to the config file")
}

func (s *serveCmd) Usage() string {
	return `serve [-conf <config file>]
  Serve calendar service using information from the config file
`
}

func (s *serveCmd) Execute(_ context.Context, _ *flag.FlagSet, _ ...any) subcommands.ExitStatus {
	log.Println("[Jasmine] Starting...")

	if s.configPath == "" {
		log.Println("[Jasmine] Error: config flag is missing")
		return subcommands.ExitUsageError
	}

	openConf, err := os.Open(s.configPath)
	if err != nil {
		if os.IsNotExist(err) {
			log.Println("[Jasmine] Error: missing config file")
		} else {
			log.Println("[Jasmine] Error: open config file: ", err)
		}
		return subcommands.ExitFailure
	}

	var config jasmine.Conf
	err = json.NewDecoder(openConf).Decode(&config)
	if err != nil {
		log.Println("[Jasmine] Error: invalid config file: ", err)
		return subcommands.ExitFailure
	}

	configPathAbs, err := filepath.Abs(s.configPath)
	if err != nil {
		log.Fatal("[Jasmine] Failed to get absolute config path")
	}
	wd := filepath.Dir(configPathAbs)
	normalLoad(config, wd)
	return subcommands.ExitSuccess
}

func normalLoad(startUp jasmine.Conf, wd string) {
	srv := jasmine.NewHttpServer(startUp, wd)
	log.Printf("[Jasmine] Starting HTTP server on '%s'\n", srv.Addr)
	go utils.RunBackgroundHttp("HTTP", srv)

	exit_reload.ExitReload("Jasmine", func() {}, func() {
		// stop http server
		_ = srv.Close()
	})
}
