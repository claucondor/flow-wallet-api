package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/flow-hydraulics/flow-wallet-api/configs"
	"github.com/flow-hydraulics/flow-wallet-api/internal/app"
	log "github.com/sirupsen/logrus"
)

const version = "2.0.0"

var (
	sha1ver   string // sha1 revision used to build the program
	buildTime string // when the executable was built
)

func main() {
	var (
		printVersion bool
		envFilePath  string // LEGACY: now used to check if user still is using envFilePath
	)

	// If we should just print the version number and exit
	flag.BoolVar(&printVersion, "version", false, "if true, print version and exit")
	flag.StringVar(&envFilePath, "envfile", "", "deprecated")
	flag.Parse()

	if envFilePath != "" {
		panic("'-envfile' is no longer supported, see readme")
	}

	if printVersion {
		fmt.Printf("v%s build on %s from sha1 %s\n", version, buildTime, sha1ver)
		os.Exit(0)
	}

	cfg, err := configs.Parse()
	if err != nil {
		panic(err)
	}

	// Example: Configure V2-specific settings
	cfg.Port = 3001 // Different port for V2
	log.Info("Starting API v2.0 with enhanced features")

	// Future: You could create app.NewAPIV2() for different business logic
	apiApp := app.NewAPI(cfg, sha1ver, buildTime)
	if err := apiApp.Run(); err != nil {
		panic(err)
	}

	os.Exit(0)
}