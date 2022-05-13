package asana

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"

	"golang.org/x/oauth2"
)

type AsanaClient interface {
	fetchTasksJSON(ctx context.Context, projectId string) ([]byte, error)
}

type asana struct {
	client *http.Client
}

func New(ctx context.Context, token string) AsanaClient {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)
	return asana{client: tc}
}

func (a asana) fetchTasksJSON(ctx context.Context, projectId string) ([]byte, error) {
	resp, err := a.client.Get(fmt.Sprintf("https://app.asana.com/api/1.0/tasks?project=%s&opt_fields=name,tags,completed,notes,permalink_url", projectId))
	if err != nil {
		return []byte{}, err
	}

	return ioutil.ReadAll(resp.Body)
}

func FetchTasks(ctx context.Context, client AsanaClient, projectId string) ([]Task, error) {
	response, err := client.fetchTasksJSON(ctx, projectId)
	if err != nil {
		return []Task{}, err
	}

	return ParseJSON(response)

}
