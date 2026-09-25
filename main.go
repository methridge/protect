package main

import (
	"os"

	"github.com/methridge/protect/cmd"
	"github.com/methridge/protect/internal/logger"
)

func main() {
	err := cmd.Execute()

	// Fetch the logger after Execute: SetLevel replaces the global logger
	// once the configured log level is known.
	log := logger.Get()

	if err != nil {
		log.Errorw("Failed to execute command", "error", err)
		log.Sync()
		os.Exit(1)
	}
	log.Sync()
}
