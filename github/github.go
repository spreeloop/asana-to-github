package github

import (
	"context"
	"fmt"

	"github.com/google/go-github/v44/github"
	"golang.org/x/oauth2"
)

type GithubInterface interface {
	createIssueImplementation(ctx context.Context, org string, repo string, issue *github.IssueRequest) (*github.Issue, *github.Response, error)
}

type githubClient struct {
	client *github.Client
}

func New(ctx context.Context, token string) GithubInterface {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := githubClient{client: github.NewClient(tc)}
	return client
}

func (c githubClient) createIssueImplementation(ctx context.Context, org string, repo string, issue *github.IssueRequest) (*github.Issue, *github.Response, error) {
	return c.client.Issues.Create(ctx, org, repo, issue)
}

func CreateIssue(ctx context.Context, client GithubInterface, org string, repo string, title string, body string, labels []string, completed bool) error {
	if title == "" {
		return fmt.Errorf("empty title field")
	}

	state := "open"
	if completed {
		state = "closed"
	}
	labels = append(labels, "asana-to-github")

	req := github.IssueRequest{
		Title:  &title,
		Body:   &body,
		Labels: &labels,
		State:  &state,
	}

	_, _, err := client.createIssueImplementation(ctx, org, repo, &req)
	if err != nil {
		return err
	}

	return nil
}
