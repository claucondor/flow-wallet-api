package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/flow-hydraulics/flow-wallet-api/configs"
	log "github.com/sirupsen/logrus"
)

const version = "0.9.0"

var (
	sha1ver   string // sha1 revision used to build the program
	buildTime string // when the executable was built
)

func main() {
	var (
		printVersion bool
		adminCommand string
	)

	// Parse flags
	flag.BoolVar(&printVersion, "version", false, "if true, print version and exit")
	flag.StringVar(&adminCommand, "command", "", "admin command to run (migrate, seed, backup)")
	flag.Parse()

	if printVersion {
		fmt.Printf("v%s build on %s from sha1 %s\n", version, buildTime, sha1ver)
		os.Exit(0)
	}

	cfg, err := configs.Parse()
	if err != nil {
		panic(err)
	}

	configs.ConfigureLogger(cfg.LogLevel)
	log.Info("Starting admin tool")

	// Example admin commands
	switch adminCommand {
	case "migrate":
		log.Info("Running database migrations...")
		// Future: Add migration logic
	case "seed":
		log.Info("Seeding database...")
		// Future: Add seeding logic  
	case "backup":
		log.Info("Creating database backup...")
		// Future: Add backup logic
	default:
		fmt.Println("Available commands: migrate, seed, backup")
		fmt.Println("Usage: ./admin -command=migrate")
	}

	os.Exit(0)
}