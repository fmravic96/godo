# godo

A simple CLI todo application written in Go.

## Inspiration

https://roadmap.sh/projects/task-tracker

## Features
- Add, update, delete tasks
- Mark tasks as in-progress or done
- List tasks (all, or filtered by status)
- Persistent storage in a JSON file

## Usage

### Build
```
go build -o godo
```

### Run
```
./godo [command] [arguments] [flags]
```

### Examples
- Add a task:
  ```
  ./godo add "Buy groceries"
  ```
- Update a task's description:
  ```
  ./godo update 1 "Buy groceries and cook dinner"
  ```
- Delete a task:
  ```
  ./godo delete 1
  ```
- Mark a task as in-progress:
  ```
  ./godo mark-in-progress 2
  ```
- Mark a task as done:
  ```
  ./godo mark-done 2
  ```
- List all tasks:
  ```
  ./godo list
  ```
- List only done tasks:
  ```
  ./godo list done
  ```

### Custom storage file
You can specify a custom JSON file for storage:
```
./godo add "Read a book" --file=mytasks.json
```

## Testing
Run all tests:
```
go test ./...
```
