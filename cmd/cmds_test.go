package cmd

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"testing"

	"github.com/fmravic96/godo/internal"
)

var staticMockTasks = []internal.Task{
	{ID: 1, Description: "Mock Task 1", Status: internal.StatusTodo},
	{ID: 2, Description: "Mock Task 2", Status: internal.StatusDone},
}

func getMockTasks() []internal.Task {
	tasks := make([]internal.Task, len(staticMockTasks))
	copy(tasks, staticMockTasks)
	return tasks
}

func writeMockTasksToFile(filePath string, tasks []internal.Task) error {
	return internal.WriteTasks(filePath, tasks)
}

func TestAddTask(t *testing.T) {
	tmpFile := "test_tasks.json"
	taskName := "Test CLI task"

	if err := writeMockTasksToFile(tmpFile, []internal.Task{}); err != nil {
		t.Fatalf("Failed to write empty mock tasks: %v", err)
	}

	cmd := RootCmd
	cmd.SetArgs([]string{"add", "Test CLI task", "--file", tmpFile})
	err := cmd.Execute()

	if err != nil {
		t.Fatalf("Command failed: %v", err)
	}

	tasks, err := internal.ReadTasks(tmpFile)
	if err != nil {
		t.Fatalf("ReadTasks failed: %v", err)
	}
	if len(tasks) != 1 {
		t.Fatalf("Expected 1 task, got %d", len(tasks))
	}
	if tasks[0].Description != taskName {
		t.Errorf("Expected description '%s', got '%s'", taskName, tasks[0].Description)
	}
	if tasks[0].Status != internal.StatusTodo {
		t.Errorf("Expected status '%s', got '%s'", internal.StatusTodo, tasks[0].Status)
	}
}

func TestUpdateTask(t *testing.T) {
	tmpFile := "test_tasks.json"
	tasks := getMockTasks()
	if err := writeMockTasksToFile(tmpFile, tasks); err != nil {
		t.Fatalf("Failed to write mock tasks: %v", err)
	}
	taskID := tasks[0].ID

	// Update the task
	newDesc := "Updated Task Description"
	cmd := RootCmd
	cmd.SetArgs([]string{"update", fmt.Sprintf("%d", taskID), newDesc, "--file", tmpFile})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("UpdateCmd failed: %v", err)
	}

	// Verify update
	updatedTasks, err := internal.ReadTasks(tmpFile)
	if err != nil {
		t.Fatalf("ReadTasks failed: %v", err)
	}
	if updatedTasks[0].Description != newDesc {
		t.Errorf("Expected description '%s', got '%s'", newDesc, updatedTasks[0].Description)
	}
}

func TestDeleteTask(t *testing.T) {
	tmpFile := "test_tasks.json"
	tasks := getMockTasks()
	if err := writeMockTasksToFile(tmpFile, tasks); err != nil {
		t.Fatalf("Failed to write mock tasks: %v", err)
	}
	taskID := tasks[0].ID

	cmd := RootCmd
	cmd.SetArgs([]string{"delete", fmt.Sprintf("%d", taskID), "--file", tmpFile})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("DeleteCmd failed: %v", err)
	}

	updatedTasks, err := internal.ReadTasks(tmpFile)
	if err != nil {
		t.Fatalf("ReadTasks failed: %v", err)
	}
	for _, task := range updatedTasks {
		if task.ID == taskID {
			t.Errorf("Task with ID %d was not deleted", taskID)
		}
	}
	if len(updatedTasks) != len(tasks)-1 {
		t.Errorf("Expected %d tasks after delete, got %d", len(tasks)-1, len(updatedTasks))
	}
}

func TestMarkInProgressTask(t *testing.T) {
	tmpFile := "test_tasks.json"
	tasks := getMockTasks()
	if err := writeMockTasksToFile(tmpFile, tasks); err != nil {
		t.Fatalf("Failed to write mock tasks: %v", err)
	}
	taskID := tasks[0].ID

	cmd := RootCmd
	cmd.SetArgs([]string{"mark-in-progress", fmt.Sprintf("%d", taskID), "--file", tmpFile})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("MarkInProgressCmd failed: %v", err)
	}

	updatedTasks, err := internal.ReadTasks(tmpFile)
	if err != nil {
		t.Fatalf("ReadTasks failed: %v", err)
	}
	if updatedTasks[0].Status != internal.StatusInProgress {
		t.Errorf("Expected status '%s', got '%s'", internal.StatusInProgress, updatedTasks[0].Status)
	}
}

func TestMarkDoneTask(t *testing.T) {
	tmpFile := "test_tasks.json"
	tasks := getMockTasks()
	if err := writeMockTasksToFile(tmpFile, tasks); err != nil {
		t.Fatalf("Failed to write mock tasks: %v", err)
	}
	taskID := tasks[0].ID

	cmd := RootCmd
	cmd.SetArgs([]string{"mark-done", fmt.Sprintf("%d", taskID), "--file", tmpFile})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("MarkDoneCmd failed: %v", err)
	}

	updatedTasks, err := internal.ReadTasks(tmpFile)
	if err != nil {
		t.Fatalf("ReadTasks failed: %v", err)
	}
	if updatedTasks[0].Status != internal.StatusDone {
		t.Errorf("Expected status '%s', got '%s'", internal.StatusDone, updatedTasks[0].Status)
	}
}

func TestListAllTasks(t *testing.T) {
	tmpFile := "test_tasks.json"
	tasks := getMockTasks()
	if err := writeMockTasksToFile(tmpFile, tasks); err != nil {
		t.Fatalf("Failed to write mock tasks: %v", err)
	}

	cmd := RootCmd
	cmd.SetArgs([]string{"list", "--file", tmpFile})
	// Capture output
	output := captureOutput(func() {
		_ = cmd.Execute()
	})
	for _, task := range tasks {
		if !contains(output, task.Description) {
			t.Errorf("Expected to find task '%s' in output", task.Description)
		}
	}
}

func TestListDoneTasks(t *testing.T) {
	tmpFile := "test_tasks.json"
	tasks := getMockTasks()
	if err := writeMockTasksToFile(tmpFile, tasks); err != nil {
		t.Fatalf("Failed to write mock tasks: %v", err)
	}

	cmd := RootCmd
	cmd.SetArgs([]string{"list", "done", "--file", tmpFile})
	output := captureOutput(func() {
		_ = cmd.Execute()
	})
	for _, task := range tasks {
		if string(task.Status) == "done" {
			if !contains(output, task.Description) {
				t.Errorf("Expected to find done task '%s' in output", task.Description)
			}
		} else {
			if contains(output, task.Description) {
				t.Errorf("Did not expect to find non-done task '%s' in output", task.Description)
			}
		}
	}
}

func TestListInProgressTasks(t *testing.T) {
	tmpFile := "test_tasks.json"
	tasks := []internal.Task{
		{ID: 1, Description: "Task 1", Status: internal.StatusInProgress},
		{ID: 2, Description: "Task 2", Status: internal.StatusTodo},
	}
	if err := writeMockTasksToFile(tmpFile, tasks); err != nil {
		t.Fatalf("Failed to write mock tasks: %v", err)
	}

	cmd := RootCmd
	cmd.SetArgs([]string{"list", "in-progress", "--file", tmpFile})
	output := captureOutput(func() {
		_ = cmd.Execute()
	})
	if !contains(output, "Task 1") {
		t.Errorf("Expected to find 'Task 1' in output")
	}
	if contains(output, "Task 2") {
		t.Errorf("Did not expect to find 'Task 2' in output")
	}
}

// Helper to capture stdout
func captureOutput(f func()) string {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	f()
	w.Close()
	os.Stdout = old
	var buf bytes.Buffer
	io.Copy(&buf, r)
	return buf.String()
}

func contains(s, substr string) bool {
	return bytes.Contains([]byte(s), []byte(substr))
}
