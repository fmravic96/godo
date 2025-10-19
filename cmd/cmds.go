package cmd

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/fmravic96/godo/internal"
	"github.com/spf13/cobra"
)

func addTask(cmd *cobra.Command, args []string) {
	taskName := strings.Join(args, " ")
	fmt.Printf("Adding task: %s\n", taskName)

	tasks, err := internal.ReadTasks(taskFile)
	if err != nil {
		fmt.Println("Error reading tasks:", err)
		os.Exit(1)
	}

	// Generate new ID
	newID := 1
	if len(tasks) > 0 {
		newID = tasks[len(tasks)-1].ID + 1
	}

	now := time.Now()
	task := internal.Task{
		ID:          newID,
		Description: taskName,
		Status:      internal.StatusTodo,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	tasks = append(tasks, task)

	if err := internal.WriteTasks(taskFile, tasks); err != nil {
		fmt.Println("Error writing tasks file:", err)
		os.Exit(1)
	}

	fmt.Println("Task added successfully!")
}

func updateTask(cmd *cobra.Command, args []string) {

	// Parse taskID
	taskID := args[0]
	newDescription := strings.Join(args[1:], " ")

	tasks, err := internal.ReadTasks(taskFile)
	if err != nil {
		fmt.Println("Error reading tasks:", err)
		os.Exit(1)
	}

	found := false
	for i, t := range tasks {
		if fmt.Sprintf("%d", t.ID) == taskID {
			tasks[i].Description = newDescription
			tasks[i].UpdatedAt = time.Now()
			found = true
			break
		}
	}

	if !found {
		fmt.Printf("Task with ID %s not found\n", taskID)
		os.Exit(1)
	}

	if err := internal.WriteTasks(taskFile, tasks); err != nil {
		fmt.Println("Error writing tasks file:", err)
		os.Exit(1)
	}

	fmt.Println("Task updated successfully!")
}

func deleteTask(cmd *cobra.Command, args []string) {
	taskID := args[0]
	tasks, err := internal.ReadTasks(taskFile)
	if err != nil {
		fmt.Println("Error reading tasks:", err)
		os.Exit(1)
	}
	found := false
	newTasks := make([]internal.Task, 0, len(tasks))
	for _, t := range tasks {
		if fmt.Sprintf("%d", t.ID) == taskID {
			found = true
			continue // skip this task (delete)
		}
		newTasks = append(newTasks, t)
	}
	if !found {
		fmt.Printf("Task with ID %s not found\n", taskID)
		os.Exit(1)
	}
	if err := internal.WriteTasks(taskFile, newTasks); err != nil {
		fmt.Println("Error writing tasks file:", err)
		os.Exit(1)
	}
	fmt.Println("Task deleted successfully!")
}

func markInProgressTask(cmd *cobra.Command, args []string) {
	if len(args) < 1 {
		fmt.Println("Usage: godo mark-in-progress <taskID>")
		os.Exit(1)
	}
	taskID := args[0]
	tasks, err := internal.ReadTasks(taskFile)
	if err != nil {
		fmt.Println("Error reading tasks:", err)
		os.Exit(1)
	}
	found := false
	for i, t := range tasks {
		if fmt.Sprintf("%d", t.ID) == taskID {
			tasks[i].Status = internal.StatusInProgress
			tasks[i].UpdatedAt = time.Now()
			found = true
			break
		}
	}
	if !found {
		fmt.Printf("Task with ID %s not found\n", taskID)
		os.Exit(1)
	}
	if err := internal.WriteTasks(taskFile, tasks); err != nil {
		fmt.Println("Error writing tasks file:", err)
		os.Exit(1)
	}
	fmt.Println("Task marked as in-progress!")
}

func markDoneTask(cmd *cobra.Command, args []string) {
	if len(args) < 1 {
		fmt.Println("Usage: godo mark-done <taskID>")
		os.Exit(1)
	}
	taskID := args[0]
	tasks, err := internal.ReadTasks(taskFile)
	if err != nil {
		fmt.Println("Error reading tasks:", err)
		os.Exit(1)
	}
	found := false
	for i, t := range tasks {
		if fmt.Sprintf("%d", t.ID) == taskID {
			tasks[i].Status = internal.StatusDone
			tasks[i].UpdatedAt = time.Now()
			found = true
			break
		}
	}
	if !found {
		fmt.Printf("Task with ID %s not found\n", taskID)
		os.Exit(1)
	}
	if err := internal.WriteTasks(taskFile, tasks); err != nil {
		fmt.Println("Error writing tasks file:", err)
		os.Exit(1)
	}
	fmt.Println("Task marked as done!")
}

func listTasks(cmd *cobra.Command, args []string) {
	tasks, err := internal.ReadTasks(taskFile)
	if err != nil {
		fmt.Println("Error reading tasks:", err)
		os.Exit(1)
	}
	var filterStatus string
	if len(args) > 0 {
		filterStatus = args[0]
	}
	count := 0
	for _, t := range tasks {
		if filterStatus == "" || string(t.Status) == filterStatus {
			fmt.Printf("[%d] %s | %s | %s\n", t.ID, t.Description, t.Status, t.UpdatedAt.Format("2006-01-02 15:04"))
			count++
		}
	}
	if count == 0 {
		if filterStatus == "" {
			fmt.Println("No tasks found.")
		} else {
			fmt.Printf("No tasks found with status '%s'.\n", filterStatus)
		}
	}
}
