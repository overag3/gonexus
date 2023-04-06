package nexusrm

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const (
	restTasks = "service/rest/v1/tasks"
)

type listTasksResponse struct {
	Items             []Task `json:"items"`
	ContinuationToken string `json:"continuationToken"`
}

type Task struct {
	ID            string    `json:"id"`
	Name          string    `json:"name"`
	Type          string    `json:"type"`
	Message       string    `json:"message"`
	CurrentState  string    `json:"currentState"`
	LastRunResult string    `json:"lastRunResult"`
	NextRun       time.Time `json:"nextRun"`
	LastRun       time.Time `json:"lastRun"`
}

// TasksList gets a list of all tasks in the Repository Manager
func TasksList(rm RM) ([]Task, error) {
	continuation := ""

	getTasks := func() (listResp listTasksResponse, err error) {
		url := fmt.Sprintf(restTasks)

		if continuation != "" {
			url += "&continuationToken=" + continuation
		}

		body, resp, err := rm.Get(url)
		if err != nil || resp.StatusCode != http.StatusOK {
			return
		}

		err = json.Unmarshal(body, &listResp)

		return
	}

	var items []Task
	for {
		resp, err := getTasks()
		if err != nil {
			return items, fmt.Errorf("could not get tasks: %v", err)
		}

		items = append(items, resp.Items...)

		if resp.ContinuationToken == "" {
			break
		}

		continuation = resp.ContinuationToken
	}
	return items, nil
}

// TasksGetByID gets a single task by ID
func TasksGetByID(rm RM, id string) (Task, error) {
	var task Task

	url := fmt.Sprintf("%s/%s", restTasks, id)
	body, resp, err := rm.Get(url)
	if err != nil || resp.StatusCode != http.StatusOK {
		return task, fmt.Errorf("could not fetch task: %w", err)
	}

	if err := json.Unmarshal(body, &task); err != nil {
		return task, fmt.Errorf("could not unmarshal task: %w", err)
	}

	return task, nil
}

// TasksRun triggers the task with the specified ID
func TasksRun(rm RM, id string) error {
	url := fmt.Sprintf("%s/%s/run", restTasks, id)
	_, resp, err := rm.Post(url, nil)
	if err != nil && resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("could not run task: %w", err)
	}
	return nil
}

// TasksStop stops the task with the specified ID
func TasksStop(rm RM, id string) error {
	url := fmt.Sprintf("%s/%s/stop", restTasks, id)
	_, resp, err := rm.Post(url, nil)
	if err != nil && resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("could not stop task: %w", err)
	}
	return nil
}
