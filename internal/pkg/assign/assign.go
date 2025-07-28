// Package assign provides functionality for assigning individuals to a pull
// request.
package assign

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"math/rand/v2"
	"os"
	"slices"

	"github.com/google/go-github/v73/github"
	"golang.org/x/oauth2"
)

func init() {
	githubToken := os.Getenv("GITHUB_TOKEN")

	if githubToken == "" {
		fmt.Fprintf(os.Stderr, "GITHUB_TOKEN is required for for authentication.\n")
		os.Exit(2)
	}

	client = github.NewClient(oauth2.NewClient(context.Background(), oauth2.StaticTokenSource(&oauth2.Token{AccessToken: githubToken})))
}

var (
	ErrNoTeamsAssigned = errors.New("no teams are currently assigned as reviewers")
)

// PullRequest assigns the pull request for the given owner, repo, and PR number
// to a random set of individuals from a GitHub team.
func PullRequest(ctx context.Context, owner, repo string, pr int) error {
	if owner == "" || repo == "" || pr == 0 {
		return errors.New("owner, repo, and PR number are required")
	}

	teams, err := teamReviewers(ctx, client, owner, repo, pr)
	if err != nil {
		return fmt.Errorf("retrieving current reviewers: %w", err)
	}

	if len(teams) == 0 {
		return ErrNoTeamsAssigned
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

	return nil
}

func assignReviewers(ctx context.Context, client *github.Client, owner, repo string, pr int, users []string) error {
	slog.Debug("Assigning individuals.")

	reviewers := sample(users, proportion)

	slog.Debug("Individuals being assigned", "reviewers", reviewers, "sample", proportion)

	_, _, err := client.PullRequests.RequestReviewers(ctx, owner, repo, pr, github.ReviewersRequest{
		Reviewers: reviewers,
	})

	slog.Debug("Individuals assigned.")

	return err
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

	slog.Debug("Retrieved team members.", "team", team, "members", users)

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

	slog.Debug("Retrieved team reviewers.", "teams", teams)

	return teams, nil
}

const (
	proportion = 0.3
)

var (
	client *github.Client
)
