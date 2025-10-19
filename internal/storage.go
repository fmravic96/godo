package internal

import (
	"encoding/json"
	"os"
)

const DefaultTasksFile = "tasks.json"
const DefaultTestsFile = "test_tasks.json"

// WriteTasks writes the list of tasks to the specified JSON file
func WriteTasks(filePath string, tasks []Task) error {
	data, err := json.MarshalIndent(tasks, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filePath, data, 0644)
}

// ReadTasks reads the list of tasks from the specified JSON file
func ReadTasks(filePath string) ([]Task, error) {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return []Task{}, nil // No file yet, return empty slice
	}
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	var tasks []Task
	if err := json.Unmarshal(data, &tasks); err != nil {
		return nil, err
	}
	return tasks, nil
}
