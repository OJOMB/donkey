package main

import (
	"flag"
	"log/slog"
	"os"

	"github.com/OJOMB/donkey/internal/repl"
	"github.com/OJOMB/donkey/pkg/logs"
)

var (
	// flagLogOn is a flag for setting log on or off, default is off
	flagLogOn bool
)

func main() {
	// flag for setting log on or off, default is off
	flag.BoolVar(&flagLogOn, "log", false, "set logging on or off")
	flag.Parse()

	// create a new logger that writes to a file if the flag is set, otherwise use a null logger
	var logger logs.Logger
	if flagLogOn {
		logger = logs.NewMultiSlogger(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	}

	r := repl.New(os.Stdin, os.Stdout, logger)
	r.Start()
}
