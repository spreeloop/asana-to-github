package asana

import (
	"context"
	"io/ioutil"
	"os"
)

type fake struct {
	source string
}

func NewAsanaClient(source string) AsanaClient {
	return fake{source: source}
}

func (a fake) fetchTasksJSON(ctx context.Context, projectId string) ([]byte, error) {
	f, err := os.Open(a.source)
	if err != nil {
		return []byte{}, err
	}
	defer f.Close()

	return ioutil.ReadAll(f)
}
