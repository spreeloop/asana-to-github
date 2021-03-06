package asana_test

import (
	"context"
	"testing"

	"spreeloop.com/asana-to-github/asana"
)

func TestParseEmptyJson(t *testing.T) {
	source := "testdata/empty.json"
	client := asana.NewFake(source)
	tasks, err := asana.FetchTasks(context.Background(), client, "dummy-project-id")
	if err != nil {
		t.Fatal("ParseJSON returned an error:", err)
	}

	want := 0
	if got := len(tasks); got != want {
		t.Errorf("ParseJSON(%v) len = %d; want %d", source, got, want)
	}
}

func TestParseTaskWithoutSubtasks(t *testing.T) {
	source := "testdata/tasks_without_subtasks.json"
	client := asana.NewFake(source)
	tasks, err := asana.FetchTasks(context.Background(), client, "dummy-project-id")
	if err != nil {
		t.Fatal("ParseJSON returned an error:", err)
	}

	wantLength := 1
	if gotLength := len(tasks); gotLength != wantLength {
		t.Errorf("ParseJSON(%v) len = %d; want %d", source, gotLength, wantLength)
	}

	wantName := "Continuously update courier location"
	if gotName := tasks[0].Name; gotName != wantName {
		t.Errorf("ParseJSON(%v) name = %s; want %s", source, gotName, wantName)
	}
}

func TestParseTaskWithSubtasks(t *testing.T) {
	source := "testdata/tasks_with_subtasks.json"
	client := asana.NewFake(source)
	tasks, err := asana.FetchTasks(context.Background(), client, "dummy-project-id")
	if err != nil {
		t.Fatal("ParseJSON returned an error:", err)
	}

	want := 3
	if got := len(tasks); got != want {
		t.Errorf("ParseJSON(%v) len = %d; want %d", source, got, want)
	}
}

func TestParseBrokenJson(t *testing.T) {
	source := "testdata/bad_data.json"
	client := asana.NewFake(source)
	_, err := asana.FetchTasks(context.Background(), client, "dummy-project-id")
	if err == nil {
		t.Errorf("got nil error; want non-nil error")
	}
}

func TestParseNonExistingFile(t *testing.T) {
	source := "testdata/non_existing_file.json"
	client := asana.NewFake(source)
	_, err := asana.FetchTasks(context.Background(), client, "dummy-project-id")
	if err == nil {
		t.Errorf("got nil error; want non-nil error")
	}
}
