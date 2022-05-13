package github

import (
	"context"
	"fmt"
	"time"

	"github.com/google/go-github/v44/github"
	"golang.org/x/oauth2"
)

const (
	labelMigration = "asana-to-github"
	labelDuplicate = "duplicate"
	stateClosed    = "closed"
	stateOpen      = "open"
)

type GithubInterface interface {
	updateIssue(ctx context.Context, issue *github.Issue, org string, repo string, updateRequest *github.IssueRequest) (*github.Issue, *github.Response, error)
	createIssueImplementation(ctx context.Context, org string, repo string, issue *github.IssueRequest) (*github.Issue, *github.Response, error)
}

type githubClient struct {
	client            *github.Client
	delayForRateLimit int
}

func New(ctx context.Context, token string, delayForRateLimit int) GithubInterface {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := githubClient{client: github.NewClient(tc), delayForRateLimit: delayForRateLimit}
	return client
}

func (c githubClient) createIssueImplementation(ctx context.Context, org string, repo string, issue *github.IssueRequest) (*github.Issue, *github.Response, error) {
	// Avoid rate limits by sleeping for a bit.
	duration := time.Duration(c.delayForRateLimit) * time.Second
	time.Sleep(duration)

	return c.client.Issues.Create(ctx, org, repo, issue)
}

func (c githubClient) updateIssue(ctx context.Context, issue *github.Issue, org string, repo string, updateRequest *github.IssueRequest) (*github.Issue, *github.Response, error) {
	// Avoid rate limits by sleeping for a bit.
	duration := time.Duration(c.delayForRateLimit) * time.Second
	time.Sleep(duration)

	return c.client.Issues.Edit(ctx, org, repo, *issue.Number, updateRequest)
}

func MigrateIssue(ctx context.Context, client GithubInterface, migratedIssues map[string]*github.Issue, org string, repo string, title string, body string, labels []string, completed bool, permalink string, forceUpdate bool) error {
	// Ignore issues with an empty title, since they may have been
	// created by mistake.
	if title == "" {
		return fmt.Errorf("empty title field")
	}

	state := stateOpen
	if completed {
		state = stateClosed
	}

	bodyWithPermalink := fmt.Sprintf("%s\n\n%s", permalink, body)
	if body == "" {
		bodyWithPermalink = permalink
	}

	// Update the issue if it was already migrated.
	// Use the title to match the asana issues with the github ones.
	if existingIssue, alreadyMigrated := migratedIssues[title]; alreadyMigrated {
		if forceUpdate || *existingIssue.State != state {
			updateRequest := &github.IssueRequest{
				State: &state,
				Body:  &bodyWithPermalink,
			}

			fmt.Printf("Updating state of %q. %s -> %s\n", title, *existingIssue.State, state)
			_, _, err := client.updateIssue(ctx, existingIssue, org, repo, updateRequest)
			return err
		}

		fmt.Printf("Skipping migration of %s\n", title)
		return nil
	}

	// Ensure migrated issues have the asana-to-github label.
	nonEmptyLabels := []string{}
	for _, label := range labels {
		if label != "" {
			nonEmptyLabels = append(nonEmptyLabels, label)
		}
	}
	nonEmptyLabels = append(nonEmptyLabels, labelMigration)

	// Create a new issue, since it doesn't exist in github yet.
	fmt.Printf("Creating %q, bodyWithPermalink: %s, labels: %v, state: %s\n", title, bodyWithPermalink, labels, state)
	req := github.IssueRequest{
		Title:  &title,
		Body:   &bodyWithPermalink,
		Labels: &nonEmptyLabels,
	}
	createdIssue, _, err := client.createIssueImplementation(ctx, org, repo, &req)
	if err != nil {
		return err
	}

	// Update the issue state (when necessary), since this isn't handled by the create request.
	if *createdIssue.State != state {
		fmt.Printf("Updating state (post-create) of %q. %s -> %s\n", title, *createdIssue.State, state)
		updateRequest := github.IssueRequest{
			State: &state,
		}
		_, _, err = client.updateIssue(ctx, createdIssue, org, repo, &updateRequest)
		return err
	}

	return nil
}

func ListMigratedIssues(ctx context.Context, write GithubInterface, token string, org string, repo string) (map[string]*github.Issue, error) {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)
	fetch := github.NewClient(tc)

	// Fetch all issues with the asana-to-github label.
	opts := &github.IssueListByRepoOptions{
		Labels: []string{labelMigration},
		State:  "all",
	}
	var issues []*github.Issue
	for {
		pageIssues, resp, err := fetch.Issues.ListByRepo(ctx, org, repo, opts)
		if err != nil {
			return nil, err
		}
		issues = append(issues, pageIssues...)
		if resp.NextPage == 0 {
			break
		}
		opts.ListOptions.Page = resp.NextPage
	}

	// Put all issues into a map using the title as key.
	migratedIssues := make(map[string]*github.Issue)
	numDups := 0
	for _, issue := range issues {
		title := *issue.Title
		if title == "" {
			continue
		}

		// Tag duplicates that occur during the migration and close them.
		if _, isDup := migratedIssues[title]; isDup {
			state := stateClosed
			updateRequest := &github.IssueRequest{
				State:  &state,
				Labels: &[]string{labelMigration, labelDuplicate},
			}

			numDups++
			fmt.Printf("closing duplicated %q\n", title)
			_, _, err := write.updateIssue(ctx, issue, org, repo, updateRequest)
			if err != nil {
				fmt.Printf("failed to update duplicated issue %q\n", title)
			}
			continue
		}

		migratedIssues[title] = issue
	}

	fmt.Printf("found %v duplicated issues out of %v %s issues\n", numDups, len(issues), labelMigration)
	fmt.Printf("found %v already migrated issues out of %v %s issues\n", len(migratedIssues), len(issues), labelMigration)
	return migratedIssues, nil
}
