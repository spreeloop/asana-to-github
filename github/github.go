package github

import (
	"context"
	"fmt"

	"github.com/google/go-github/v44/github"
	"golang.org/x/oauth2"
)

type githubClientInterface interface {
	createIssueInternal(ctx context.Context, org string, repo string, issue *github.IssueRequest) (*github.Issue, *github.Response, error)
}

type githubClient struct {
	client *github.Client
}

type fakeGithubClient struct{}

func New(ctx context.Context, token string) githubClientInterface {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := githubClient{client: github.NewClient(tc)}
	return client
}

func NewFake() githubClientInterface {
	return fakeGithubClient{}
}

func (c githubClient) createIssueInternal(ctx context.Context, org string, repo string, issue *github.IssueRequest) (*github.Issue, *github.Response, error) {
	return c.client.Issues.Create(ctx, org, repo, issue)
}

func (c fakeGithubClient) createIssueInternal(ctx context.Context, org string, repo string, issue *github.IssueRequest) (*github.Issue, *github.Response, error) {
	fmt.Printf("%v/%v title: %s, body: %s, labels: %v, state: %s\n", org, repo, *issue.Title, *issue.Body, *issue.Labels, *issue.State)
	return nil, nil, nil
}

func CreateIssue(ctx context.Context, client githubClientInterface, org string, repo string, title string, body string, labels []string, completed bool) error {
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

	_, _, err := client.createIssueInternal(ctx, org, repo, &req)
	if err != nil {
		return err
	}

	return nil
}
