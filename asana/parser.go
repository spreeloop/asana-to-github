package asana

import (
	"encoding/json"
)

type Assignee struct {
	GID          string
	Name         string
	ResourceType string `json:"resource_type"`
}

type Follower struct {
	GID          string
	Name         string
	ResourceType string `json:"resource_type"`
}

type User struct {
	GID          string
	Name         string
	ResourceType string `json:"resource_type"`
}

type Like struct {
	GID  string
	User User
}

type Project struct {
	GID          string
	Name         string
	ResourceType string `json:"resource_type"`
}

type Section struct {
	GID          string
	Name         string
	ResourceType string `json:"resource_type"`
}

type Membership struct {
	Project Project
	Section Section
}

type Tag struct {
	GID          string
	Name         string
	ResourceType string `json:"resource_type"`
}

type Workspace struct {
	GID          string
	Name         string
	ResourceType string `json:"resource_type"`
}

type Parent struct {
	GID          string
	Name         string
	ResourceType string `json:"resource_type"`
}

type Task struct {
	Gid             string `json:"gid"`
	Assignee        Assignee
	AssigneeStatus  string `json:"assignee_status"`
	Completed       bool
	CompletedAt     string `json:"completed_at"`
	CreatedAt       string `json:"created_at"`
	DueAt           string `json:"due_at"`
	Followers       []Follower
	Hearted         bool
	Hearts          []Like
	Liked           bool
	Likes           []Like
	Membership      []Membership
	ModifiedAt      string `json:"modified_at"`
	Name            string
	Notes           string
	NumHearts       int `json:"num_hearts"`
	NumLikes        int `json:"num_likes"`
	Parent          Parent
	PermalinkURL    string `json:"permalink_url"`
	Projects        []Project
	ResourceType    string `json:"resource_type"`
	StartAt         string `json:"start_at"`
	StartOn         string `json:"start_on"`
	Subtasks        []Task
	Tags            []Tag
	ResourceSubtype string `json:"resource_subtype"`
}

type root struct {
	Data []Task
}

// ParseJSON unmarshals asana tasks from bytes in JSON format.
func ParseJSON(data []byte) ([]Task, error) {
	var r root
	err := json.Unmarshal(data, &r)
	if err != nil {
		return []Task{}, err
	}

	tasks := r.Data
	subTasks := make([]Task, 0)
	for i := 0; i < len(tasks); i++ {
		for j := 0; j < len(tasks[i].Subtasks); j++ {
			subTasks = append(subTasks, tasks[i].Subtasks[j])
		}
	}

	return append(tasks, subTasks...), nil
}
