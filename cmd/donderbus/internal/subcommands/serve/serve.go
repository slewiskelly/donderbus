// Package serve implements the "serve" subcommand.
package serve

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/google/subcommands"
	"golang.org/x/sync/errgroup"

	_serve "github.com/slewiskelly/donderbus/internal/pkg/serve"
)

// Serve implements the "assign" subcommand.
type Serve struct {
	host          string
	insecure      bool
	port          int
	webhookSecret []byte

	srv *http.Server
}

// Name returns the name of the subcommand.
func (*Serve) Name() string {
	return "serve"
}

// Synopsis returns a one-line summary of the subcommand.
func (*Serve) Synopsis() string {
	return "starts a server in which to receive webhooks that will prompt assignment"
}

// Usage returns a longer explanation and/or usage example(s) of the subcommand.
func (*Serve) Usage() string {
	return `Starts a HTTP server in which to receive webhooks that will prompt assignment.

Upon receipt of pull_request events, in which the action is ready_for_review, invididuals will be assigned who are members of the assigned team.

By default, webhook deliveries are validated against the value of GITHUB_WEBHOOK_SECRET.

Usage:
  donderbus serve [flags]

Examples:
  donderbus serve

Flags:
`
}

// SetFlags sets the flags specific to the subcommand.
func (s *Serve) SetFlags(f *flag.FlagSet) {
	f.StringVar(&s.host, "host", "0.0.0.0", "host that the server is to serve from")
	f.IntVar(&s.port, "port", 8080, "port that the server is to listen on")
	f.BoolVar(&s.insecure, "insecure", false, "skips webhook validation")
}

// Execute executes the subcommand.
func (s *Serve) Execute(ctx context.Context, fs *flag.FlagSet, args ...any) subcommands.ExitStatus {
	if err := s.execute(ctx, fs, args...); err != nil {
		fmt.Fprintln(os.Stderr, err)

		return subcommands.ExitFailure
	}

	return subcommands.ExitSuccess
}

func (s *Serve) execute(ctx context.Context, _ *flag.FlagSet, _ ...any) error {
	ctx, cancel := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	srv, err := _serve.New(_serve.Host(s.host), _serve.Insecure(s.insecure), _serve.Port(s.port))
	if err != nil {
		return fmt.Errorf("initializing server: %w", err)
	}

	grp, _ := errgroup.WithContext(ctx)

	grp.Go(func() error {
		return srv.ListenAndServe()
	})

	go func() {
		<-ctx.Done()
		srv.Shutdown(ctx)
	}()

	return grp.Wait()
}
