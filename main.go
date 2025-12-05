package main

import (
	"os"

	"github.com/methridge/protect/cmd"
	"github.com/methridge/protect/internal/logger"
)

func main() {
	log := logger.New()
	defer log.Sync()

	if err := cmd.Execute(); err != nil {
		log.Error("Failed to execute command", "error", err)
		os.Exit(1)
	}
}
