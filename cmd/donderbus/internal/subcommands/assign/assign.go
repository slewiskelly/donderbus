// Package assign implements the "assign" subcommand.
package assign

import (
	"context"
	"flag"
	"fmt"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/google/subcommands"

	_assign "github.com/slewiskelly/donderbus/internal/pkg/assign"
)

// Assign implements the "assign" subcommand.
type Assign struct{}

// Name returns the name of the subcommand.
func (*Assign) Name() string {
	return "assign"
}

// Synopsis returns a one-line summary of the subcommand.
func (*Assign) Synopsis() string {
	return "assigns individuals to a pull request"
}

// Usage returns a longer explanation and/or usage example(s) of the subcommand.
func (*Assign) Usage() string {
	return `Assigns individuals to a pull request.

The pull request must already be assigned to one or more teams.

Usage:
	donderbus assign [flags] <url>

Examples:
	donderbus assign https://github.com/acme/foo/pull/123
`
}

// SetFlags sets the flags specific to the subcommand.
func (a *Assign) SetFlags(f *flag.FlagSet) {}

// Execute executes the subcommand.
func (a *Assign) Execute(ctx context.Context, fs *flag.FlagSet, args ...any) subcommands.ExitStatus {
	if fs.NArg() < 1 {
		fmt.Fprintf(os.Stderr, "No URL provided.\n\n")
		fs.Usage()
		return subcommands.ExitUsageError
	}

	if err := a.execute(ctx, fs, args...); err != nil {
		fmt.Fprintln(os.Stderr, err)

		return subcommands.ExitFailure
	}

	return subcommands.ExitSuccess
}

func (a *Assign) execute(ctx context.Context, fs *flag.FlagSet, _ ...any) error {
	owner, repo, pr, err := fromURL(fs.Arg(0))
	if err != nil {
		return fmt.Errorf("parsing URL: %w", err)
	}

	if err := _assign.PullRequest(ctx, owner, repo, pr); err != nil {
		return fmt.Errorf("assigning pull request: %w", err)
	}

	fmt.Printf("https://github.com/%s/%s/pull/%d\n", owner, repo, pr)

	return nil
}

func fromURL(u string) (string, string, int, error) {
	parsed, err := url.Parse(u)
	if err != nil {
		return "", "", 0, fmt.Errorf("invalid URL: %w", err)
	}

	if parsed.Host != "github.com" {
		return "", "", 0, fmt.Errorf("unsupported host: %s", parsed.Host)
	}

	parts := strings.Split(strings.Trim(parsed.Path, "/"), "/")
	if len(parts) < 4 || parts[2] != "pull" {
		return "", "", 0, fmt.Errorf("unsupported URL path: %q", parsed.Path)
	}

	owner, repo := parts[0], parts[1]

	pr, err := strconv.Atoi(parts[3])
	if err != nil {
		return "", "", 0, fmt.Errorf("invalid PR number: %q", parts[3])
	}

	return owner, repo, pr, nil
}
