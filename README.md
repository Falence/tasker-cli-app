## Task manager CLI application with MongoDB database
- This CLI application permits users create task, get list of tasks(uncompleted and/or completed), complete task and delete tasks.
- Uses MongoDB persistence.

## Commands
Pass the commands and arguements whe running the application.
Here are the availaible commands:
- **Add a task:** `go run main.go add "<task>"`
- **Display uncompleted tasks:** `go run main.go`
- **Display completed tasks:** `go run main.go finished`
- **Complete a task:** `go run main.go done "<task>"`
- **Delete a task:** `go run main.go rm "<task>"`