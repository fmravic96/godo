package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var taskFile string

var AddCmd = &cobra.Command{
	Use:   "add [task name]",
	Short: "Add a new task to your todo list",
	Run:   addTask,
	Args:  cobra.MinimumNArgs(1),
}

var UpdateCmd = &cobra.Command{
	Use:   "update [task ID] [new description]",
	Short: "Update an existing task's description",
	Run:   updateTask,
	Args:  cobra.MinimumNArgs(2),
}

var DeleteCmd = &cobra.Command{
	Use:   "delete [task ID]",
	Short: "Delete a task from your todo list",
	Run:   deleteTask,
	Args:  cobra.MinimumNArgs(1),
}

var MarkInProgressCmd = &cobra.Command{
	Use:   "mark-in-progress [task ID]",
	Short: "Mark a task as in-progress",
	Args:  cobra.MinimumNArgs(1),
	Run:   markInProgressTask,
}

var MarkDoneCmd = &cobra.Command{
	Use:   "mark-done [task ID]",
	Short: "Mark a task as done",
	Args:  cobra.MinimumNArgs(1),
	Run:   markDoneTask,
}

var ListCmd = &cobra.Command{
	Use:   "list [status]",
	Short: "List all tasks or filter by status",
	Args:  cobra.MaximumNArgs(1),
	Run:   listTasks,
}

// rootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "godo",
	Short: "Simple CLI todo application",
	Long:  `godo is a simple command-line todo application to help you manage your tasks efficiently.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := RootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	RootCmd.PersistentFlags().StringVar(&taskFile, "file", "tasks.json", "Path to the tasks storage file")
	RootCmd.AddCommand(AddCmd, UpdateCmd, DeleteCmd, MarkInProgressCmd, MarkDoneCmd, ListCmd)
}
