package tasks

import (
	"context"
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/yarlson/tap"
)

const (
	GM = 398600.4415 // км³/с²
	J2 = 0.0010826
	Re = 6378.137    // км
)

func init() {
	Register(Task{
		ID:          "2",
		Name:        "Spacecraft movement model",
		Description: "Calculate spacecraft movement model for central gravitational force field or for gravity field including the second zonal harmonic J20",
		Handler:     task2,
	})
}

func task2(ctx context.Context) error {
	choices := []tap.SelectOption[int]{
		{Value: 1, Label: "Central gravitational force field"},
		{Value: 2, Label: "Gravity field including the second zonal harmonic J20"},
	}

	selected := tap.Select(ctx, tap.SelectOptions[int]{
		Message: "Select type",
		Options: choices,
	})

	// Получаем высоту
	altitudeStr := tap.Text(ctx, tap.TextOptions{
		Message: "Input altitude (km)",
	})
	altitude, _ := strconv.ParseFloat(altitudeStr, 64)

	// Получаем начальные условия для скорости
	velocityRaw := strings.Split(tap.Text(ctx, tap.TextOptions{
		Message: "Input initial velocity components vx vy vz (km/s) separated by space",
	}), " ")
	
	if len(velocityRaw) != 3 {
		return fmt.Errorf("need exactly 3 velocity components")
	}

	// Парсим компоненты скорости
	vx, _ := strconv.ParseFloat(velocityRaw[0], 64)
	vy, _ := strconv.ParseFloat(velocityRaw[1], 64)
	vz, _ := strconv.ParseFloat(velocityRaw[2], 64)

	// Начальное состояние (позиция на высоте над экватором)
	initialState := []float64{
		Re + altitude, // x
		0,              // y  
		0,              // z
		vx,            // vx
		vy,            // vy
		vz,            // vz
	}

	// Получаем параметры симуляции
	dt := getFloatInput(ctx, "Input simulation time step (seconds)")
	duration := getFloatInput(ctx, "Input simulation duration (seconds)")

	// Выбираем модель гравитации
	var gravityFunc func(float64, []float64) []float64
	if selected == 1 {
		gravityFunc = centralGravity
	} else {
		gravityFunc = j2Gravity
	}

	// Запускаем моделирование
	states := simulateOrbit(initialState, dt, duration, gravityFunc)

	// Выводим результаты
	displayResults(states, dt)

	return nil
}

func getFloatInput(ctx context.Context, message string) float64 {
	for {
		input := tap.Text(ctx, tap.TextOptions{Message: message})
		if value, err := strconv.ParseFloat(input, 64); err == nil {
			return value
		}
		tap.Message("Invalid input, please enter a number")
	}
}

func simulateOrbit(initial []float64, dt, duration float64, gravityFunc func(float64, []float64) []float64) [][]float64 {
	steps := int(duration / dt)
	states := make([][]float64, steps+1)
	
	// Начальное состояние
	currentState := make([]float64, len(initial))
	copy(currentState, initial)
	states[0] = make([]float64, len(initial))
	copy(states[0], initial)

	// Интегрирование по времени
	for step := 1; step <= steps; step++ {
		t := float64(step) * dt
		
		// Один шаг RK4 для системы ОДУ
		newState := rk4StepSystem(gravityFunc, t-dt, currentState, dt)
		
		states[step] = make([]float64, len(newState))
		copy(states[step], newState)
		copy(currentState, newState)
	}

	return states
}

// RK4 для систем ОДУ
func rk4StepSystem(f func(float64, []float64) []float64, t float64, y []float64, h float64) []float64 {
	n := len(y)
	
	// k1 = f(t, y)
	k1 := f(t, y)
	
	// k2 = f(t + h/2, y + h/2 * k1)
	y2 := make([]float64, n)
	for i := 0; i < n; i++ {
		y2[i] = y[i] + (h/2)*k1[i]
	}
	k2 := f(t+h/2, y2)
	
	// k3 = f(t + h/2, y + h/2 * k2)
	y3 := make([]float64, n)
	for i := 0; i < n; i++ {
		y3[i] = y[i] + (h/2)*k2[i]
	}
	k3 := f(t+h/2, y3)
	
	// k4 = f(t + h, y + h * k3)
	y4 := make([]float64, n)
	for i := 0; i < n; i++ {
		y4[i] = y[i] + h*k3[i]
	}
	k4 := f(t+h, y4)
	
	// y_new = y + h/6 * (k1 + 2*k2 + 2*k3 + k4)
	result := make([]float64, n)
	for i := 0; i < n; i++ {
		result[i] = y[i] + (h/6)*(k1[i]+2*k2[i]+2*k3[i]+k4[i])
	}
	
	return result
}

func displayResults(states [][]float64, dt float64) {
	if len(states) == 0 {
		fmt.Println("No simulation data")
		return
	}
	
	fmt.Println("\n=== Simulation Results ===")
	for step, state := range states {
		time := float64(step) * dt
		r := math.Sqrt(state[0]*state[0] + state[1]*state[1] + state[2]*state[2])
		altitude := r - Re
		
		fmt.Printf("Time: %.1f s | Altitude: %.2f km | Position: (%.1f, %.1f, %.1f) km\n",
			time, altitude, state[0], state[1], state[2])
	}
	
	// Итоговое состояние
	final := states[len(states)-1]
	rFinal := math.Sqrt(final[0]*final[0] + final[1]*final[1] + final[2]*final[2])
	fmt.Printf("\nFinal altitude: %.2f km\n", rFinal-Re)
}

func centralGravity(t float64, state []float64) []float64 {
	x, y, z := state[0], state[1], state[2]
	vx, vy, vz := state[3], state[4], state[5]
	
	r := math.Sqrt(x*x + y*y + z*z)
	r3 := r * r * r
	
	dxdt := vx
	dydt := vy
	dzdt := vz
	dvxdt := -GM * x / r3
	dvydt := -GM * y / r3
	dvzdt := -GM * z / r3
	
	return []float64{dxdt, dydt, dzdt, dvxdt, dvydt, dvzdt}
}

func j2Gravity(t float64, state []float64) []float64 {
	x, y, z := state[0], state[1], state[2]
	vx, vy, vz := state[3], state[4], state[5]

	r := math.Sqrt(x*x + y*y + z*z)
	r2 := r * r
	r3 := r2 * r

	z2 := z * z

	// Общий множитель J2
	factor := -1.5 * J2 * (Re*Re) / (r2 * r2)

	// Поправки
	Cxy := factor * (5.0*z2/r2 - 1.0)
	Cz  := factor * (5.0*z2/r2 - 3.0)

	ax := -GM * x / r3 * (1 + Cxy)
	ay := -GM * y / r3 * (1 + Cxy)
	az := -GM * z / r3 * (1 + Cz)

	return []float64{vx, vy, vz, ax, ay, az}
}
