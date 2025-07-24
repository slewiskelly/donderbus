// Package version implements the "version" subcommand.
package version

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"
	"runtime/debug"

	"github.com/google/subcommands"
)

// Version implements the "version" subcommand.
type Version struct{}

// Name returns the name of the subcommand.
func (*Version) Name() string {
	return "version"
}

// Synopsis returns a one-line summary of the subcommand.
func (*Version) Synopsis() string {
	return "displays the version of donderbus"
}

// Usage returns a longer explanation and/or usage example(s) of the subcommand.
func (*Version) Usage() string {
	return `donderbus version
`
}

// SetFlags sets the flags specific to the subcommand.
func (v *Version) SetFlags(f *flag.FlagSet) {}

// Execute executes the subcommand.
func (v *Version) Execute(ctx context.Context, fs *flag.FlagSet, args ...interface{}) subcommands.ExitStatus {
	if err := v.execute(ctx, fs, args...); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return subcommands.ExitFailure
	}

	return subcommands.ExitSuccess
}

func (v *Version) execute(_ context.Context, _ *flag.FlagSet, _ ...interface{}) error {
	i, ok := debug.ReadBuildInfo()
	if !ok { // Shouldn't happen.
		return errors.New("failed to read build info")
	}

	var os, arch string

	for _, s := range i.Settings {
		switch s.Key {
		case "GOOS":
			os = s.Value
		case "GOARCH":
			arch = s.Value
		}
	}

	fmt.Printf("%s %s/%s\n", i.Main.Version, os, arch)

	return nil
}
