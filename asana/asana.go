package asana

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"golang.org/x/oauth2"
)

type asanaClientInterface interface {
	fetchTasksJSON(ctx context.Context, projectId string) ([]byte, error)
}

type asanaClient struct {
	client *http.Client
}

type fakeAsanaClient struct {
	source string
}

func New(ctx context.Context, token string) asanaClientInterface {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)
	return asanaClient{client: tc}
}

func NewFake(source string) asanaClientInterface {
	return fakeAsanaClient{source: source}
}

func (a asanaClient) fetchTasksJSON(ctx context.Context, projectId string) ([]byte, error) {
	resp, err := a.client.Get(fmt.Sprintf("https://app.asana.com/api/1.0/tasks?project=%s&opt_fields=name,tags,completed,notes", projectId))
	if err != nil {
		return []byte{}, err
	}

	return ioutil.ReadAll(resp.Body)
}

func (a fakeAsanaClient) fetchTasksJSON(ctx context.Context, projectId string) ([]byte, error) {
	f, err := os.Open(a.source)
	if err != nil {
		return []byte{}, err
	}
	defer f.Close()

	return ioutil.ReadAll(f)
}

func FetchTasks(ctx context.Context, client asanaClientInterface, projectId string) ([]Task, error) {
	response, err := client.fetchTasksJSON(ctx, projectId)
	if err != nil {
		return []Task{}, err
	}

	return ParseJSON(response)

}
