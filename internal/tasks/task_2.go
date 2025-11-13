package tasks

import (
	"context"
	"fmt"

	"github.com/yarlson/tap"
)

func init() {
	Register(Task{
		ID: "2",
		Name: "Spacecraft movement model",
		Description: "Calculate spacecraft movement model for central gravitational force field or for gravity field including the second zonal harmonic J20",
		Handler: task2,
	})
}

func task2(ctx context.Context) error {
	choices := []tap.SelectOption[string]{
		{Value: "1", Label: "Central gravitational force field"},
		{Value: "2", Label: "Gravity field including the second zonal harmonic J20"},
	}

	selected := tap.Select(ctx, tap.SelectOptions[string]{
		Message: "Select type",
		Options: choices,
	})

	fmt.Println(selected)
	return nil
}