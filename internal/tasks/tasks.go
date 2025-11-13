package tasks

import "context"

type Task struct {
    ID          string
    Name        string
    Description string
    Handler     func(ctx context.Context) error
}

var Registry = []Task{}

func Register(task Task) {
    Registry = append(Registry, task)
}

func GetByID(id string) *Task {
    for _, task := range Registry {
        if task.ID == id {
            return &task
        }
    }
    return nil
}