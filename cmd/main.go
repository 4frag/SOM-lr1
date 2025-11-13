package main

import (
	"context"
	"fmt"

	"github.com/4frag/SOM-lr1/internal/tasks"
	"github.com/yarlson/tap"
)

func main() {
	ctx := context.Background()

    tap.Intro(fmt.Sprintf("%sWelcome! 👋%s", tap.Green, tap.Green))

	tasks_options := make([]tap.SelectOption[string], 0, len(tasks.Registry))
    for _, task := range tasks.Registry {
        tasks_options = append(tasks_options, tap.SelectOption[string]{
            Value: task.ID,
            Label: task.Name, 
            Hint:  task.Description,
        })
    }
    selected := tap.Select(ctx, tap.SelectOptions[string]{
		Message: "Choose task for execution:",
		Options: tasks_options,
	})

	task := tasks.GetByID(selected)
	task.Handler(ctx)
}