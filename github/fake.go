package github

import (
	"context"

	"github.com/google/go-github/v44/github"
)

type fakeGithubClient struct{}

func NewFake() GithubInterface {
	return fakeGithubClient{}
}

func (c fakeGithubClient) updateIssue(ctx context.Context, issue *github.Issue, org string, repo string, updateRequest *github.IssueRequest) (*github.Issue, *github.Response, error) {
	return issue, &github.Response{}, nil
}

func (c fakeGithubClient) createIssueImplementation(ctx context.Context, org string, repo string, issue *github.IssueRequest) (*github.Issue, *github.Response, error) {
	state := "open"
	return &github.Issue{State: &state}, &github.Response{}, nil
}
