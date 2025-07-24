// Package assign implements the "assign" subcommand.
package assign

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log/slog"
	"math/rand/v2"
	"net/url"
	"os"
	"slices"
	"strconv"
	"strings"

	"github.com/google/go-github/v73/github"
	"github.com/google/subcommands"
	"golang.org/x/oauth2"
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
	return `Assigns individuals to a pull request

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
		fmt.Fprintf(os.Stderr, "No path provided.\n\n")
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
	if fs.NArg() < 1 {
		fmt.Fprintln(os.Stderr, a.Usage())

		return errors.New("no URL provided")
	}

	owner, repo, pr, err := fromURL(fs.Arg(0))
	if err != nil {
		return fmt.Errorf("parsing URL: %w", err)
	}

	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		return errors.New("GITHUB_TOKEN is required for for authentication")
	}

	client := github.NewClient(oauth2.NewClient(ctx, oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})))

	teams, err := teamReviewers(ctx, client, owner, repo, pr)
	if err != nil {
		return fmt.Errorf("retrieving current reviewers: %w", err)
	}

	if len(teams) == 0 {
		fmt.Println("No teams are currently assigned as reviewers!")
		return nil
	}

	var individuals []string

	for _, team := range teams {
		members, err := teamMembers(ctx, client, owner, team)
		if err != nil {
			return fmt.Errorf("retrieving users in team %q: %w", team, err)
		}

		individuals = append(individuals, members...)
	}

	slices.Sort(individuals)
	individuals = slices.Compact(individuals)

	if err := assignReviewers(ctx, client, owner, repo, pr, individuals); err != nil {
		return fmt.Errorf("assigning reviewers: %w", err)
	}

	fmt.Printf("https://github.com/%s/%s/pull/%d\n", owner, repo, pr)

	return nil
}

func assignReviewers(ctx context.Context, client *github.Client, owner, repo string, pr int, users []string) error {
	slog.Debug("Assigning individuals.")

	reviewers := sample(users, proportion)

	slog.Debug("Individuals being assigned", "reviewers", strings.Join(reviewers, ","), "sample", proportion)

	_, _, err := client.PullRequests.RequestReviewers(ctx, owner, repo, pr, github.ReviewersRequest{
		Reviewers: reviewers,
	})

	slog.Debug("Individuals assigned.")

	return err
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

func teamMembers(ctx context.Context, client *github.Client, owner, team string) ([]string, error) {
	slog.Debug("Retrieving team members.", "team", team)

	var users []string

	opts := &github.TeamListTeamMembersOptions{
		ListOptions: github.ListOptions{PerPage: 100},
		Role:        "all",
	}

	for {
		members, resp, err := client.Teams.ListTeamMembersBySlug(ctx, owner, team, opts)
		if err != nil {
			return nil, fmt.Errorf("listing team members: %w", err)
		}

		for _, m := range members {
			users = append(users, m.GetLogin())
		}

		if resp.NextPage == 0 {
			break
		}

		opts.Page = resp.NextPage
	}

	slog.Debug("Retrieved team members.", "team", team, "members", strings.Join(users, ","))

	return users, nil
}

func sample[S ~[]T, T any](s S, proportion float64) S {
	if len(s) == 0 || proportion <= 0 {
		return nil
	}

	x := slices.Clone(s)

	k := int((float64(len(s)) * proportion) + rand.Float64())
	k = min(max(k, 1), len(s))

	rand.Shuffle(len(x), func(i, j int) { x[i], x[j] = x[j], x[i] })

	return x[:k]
}

func teamReviewers(ctx context.Context, client *github.Client, owner, repo string, pr int) ([]string, error) {
	slog.Debug("Retrieving team reviewers.")

	var teams []string

	reviewers, _, err := client.PullRequests.ListReviewers(ctx, owner, repo, pr, nil)
	if err != nil {
		return nil, err
	}

	for _, t := range reviewers.Teams {
		teams = append(teams, t.GetSlug())
	}

	slog.Debug("Retrieved team reviewers.", "teams", strings.Join(teams, ","))

	return teams, nil
}

const (
	proportion = 0.3
)

var (
	errUsage = errors.New("usage error")
)
