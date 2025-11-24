package tasks

import (
	"context"
	"fmt"
	"strconv"

	"github.com/yarlson/tap"
)

func init() {
    Register(Task{
        ID:          "1",
        Name:        "Runge-Kutta Method",
        Description: "4th order numerical integration",
        Handler:     RunRungeKutta,
    })
}

func inputFloat(ctx context.Context, prompt string, placeholder string, defaultValue float64) float64 {
    for {
        input := tap.Text(ctx, tap.TextOptions{
			Message: prompt,
            Placeholder: placeholder,
        })

        if input == "" {
            return defaultValue
        }

        value, err := strconv.ParseFloat(input, 64)
        if err == nil {
            return value
        }

        tap.Message(fmt.Sprintf("%s❌ Ошибка: введите корректное число%s", tap.Red, tap.Red))
    }
}

func RK4(function func(float64, float64) float64, x0 float64, y0 float64, x_end float64, h float64) float64 {
	n := int((x_end - x0)/h)

	for range n {
		k1 := function(x0, y0)
		k2 := function(x0 + h/2, y0 + (h/2) * k1)
		k3 := function(x0 + h/2, y0 + (h/2) * k2)
		k4 := function(x0 + h, y0 + h * k3)

		x0 += h
		y0 = y0 + (h/6) * (k1 + 2 * k2 + 2 * k3 + k4)
	}

	return y0
}

func RunRungeKutta(ctx context.Context) error {
	expression := tap.Text(ctx, tap.TextOptions{
		Message: "Input f(x, y) for different equation dy/dx = f(x, y)",
		Placeholder: "Input f(x, y)...",
	})
	function, err := CreateFunction(expression)
	if err != nil {
		panic(err)
	}

	x0 := inputFloat(ctx, "Input initial value for x0", "Input x0...", 0.0)
	y0 := inputFloat(ctx, "Input initial value for y0", "Input y0...", 0.0)
	x_end := inputFloat(ctx, "Input right border of the interval for integration (x1)", "Input x1...", 1.0)
	h := inputFloat(ctx, "Input step (h)", "Input h...", 0.0)

	result := RK4(function, x0, y0, x_end, h)
	
	fmt.Printf("Result: %f", result)

    return nil
}