package main

import (
	"context"
	"flag"
	"log/slog"
	"os"

	"github.com/google/subcommands"

	"github.com/slewiskelly/donderbus/cmd/donderbus/internal/subcommands/assign"
	"github.com/slewiskelly/donderbus/cmd/donderbus/internal/subcommands/version"
)

var (
	debug = flag.Bool("debug", false, "enable debug logging")
)

func init() {
	flag.Parse()

	subcommands.Register(&assign.Assign{}, "")
	subcommands.Register(&version.Version{}, "")
	subcommands.Register(subcommands.HelpCommand(), "")

	if *debug {
		slog.SetDefault(slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelDebug})))
	}
}

func main() {
	os.Exit(int(subcommands.Execute(context.Background())))
}
