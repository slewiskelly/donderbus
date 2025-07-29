// Package serve provides functionality for serving webhooks.
package serve

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"

	"github.com/google/go-github/v73/github"

	_assign "github.com/slewiskelly/donderbus/internal/pkg/assign"
)

// Server is a HTTP server which, in receipt of specific webhook events, prompts
// the assignment of pull requests.
//
// After calling [ListenAndServe], and upon receipt of pull_request events, in
// which the action is ready_for_review, invididuals will be assigned who are
// members of the assigned team.
//
// By default, webhook deliveries are validated against the value of
// GITHUB_WEBHOOK_SECRET.
type Server struct {
	opts *options

	srv *http.Server
}

// New initializes a [Server].
func New(opts ...Option) (*Server, error) {
	s := &Server{
		opts: &options{
			host:          "0.0.0.0",
			insecure:      false,
			port:          8080,
			webhookSecret: []byte(os.Getenv("GITHUB_WEBHOOK_SECRET")),
		},

		srv: &http.Server{},
	}

	for _, opt := range opts {
		opt.apply(s.opts)
	}

	if len(s.opts.webhookSecret) < 1 && !s.opts.insecure {
		return nil, errors.New("GITHUB_WEBHOOK_SECRET is required to validate webhook deliveries")
	}

	http.Handle("/assign", s)

	return s, nil
}

// ListenAndServe listens on the server's HOST:PORT and begins sering HTTP
// requests.
func (s *Server) ListenAndServe() error {
	slog.Info("Serving HTTP", "host", s.opts.host, "port", s.opts.port)

	ln, err := net.Listen("tcp", fmt.Sprintf("%s:%d", s.opts.host, s.opts.port))
	if err != nil {
		return fmt.Errorf("listening: %w", err)
	}

	if err := s.srv.Serve(ln); !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("serving: %w", err)
	}

	return nil
}

// ServeHTTP serves HTTP requests.
//
// Specifically, ServeHTTP handles webhook deliveries from GitHub, handling
// only pull request events.
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	payload, err := github.ValidatePayload(r, s.opts.webhookSecret)
	if err != nil {
		slog.Error("Failed to validate payload", "error", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	event, err := github.ParseWebHook(github.WebHookType(r), payload)
	if err != nil {
		slog.Error("Failed to parse payload", "error", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	switch event := event.(type) {
	case *github.PullRequestEvent:
		err = s.handlePullRequestEvent(r.Context(), event)
	default:
		slog.Debug("Untargeted event")
		return
	}

	if err != nil {
		slog.Error("Failed to handle event", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

// Shutdown gracefully shuts down the server.
func (s *Server) Shutdown(ctx context.Context) {
	slog.Info("Gracefully shutting down HTTP")
	defer slog.Info("Shutdown HTTP")

	s.srv.Shutdown(ctx)
}

func (s *Server) handlePullRequestEvent(ctx context.Context, e *github.PullRequestEvent) error {
	if a := e.GetAction(); a != "ready_for_review" {
		slog.Debug("Untargeted action", "action", a)
		return nil
	}

	if err := _assign.PullRequest(ctx, e.GetRepo().GetOwner().GetLogin(), e.GetRepo().GetName(), e.GetPullRequest().GetNumber()); err != nil {
		return fmt.Errorf("handling pull request: %w", err)
	}

	return nil
}
