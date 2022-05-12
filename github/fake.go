package github

import (
	"context"
	"fmt"

	"github.com/google/go-github/v44/github"
)

type fakeGithubClient struct{}

func NewFake() GithubInterface {
	return fakeGithubClient{}
}

func (c fakeGithubClient) createIssueImplementation(ctx context.Context, org string, repo string, issue *github.IssueRequest) (*github.Issue, *github.Response, error) {
	fmt.Printf("%v/%v title: %s, body: %s, labels: %v, state: %s\n", org, repo, *issue.Title, *issue.Body, *issue.Labels, *issue.State)
	return nil, nil, nil
}
