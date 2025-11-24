package tasks

import (
	"context"
	"encoding/csv"
	"fmt"
	"log"
	"math"
	"os"
	"os/exec"

	"github.com/yarlson/tap"
)

func init() {
	Register(Task{
		ID:          "3",
		Name:        "Numerical integration for 30 days",
		Description: "Simulate orbit for given initial conditions",
		Handler:     task3,
	})
}

func task3(ctx context.Context) error {

	// Начальные условия (в километрах и километрах в секунду)
	state := []float64{
		-9665.932022,  // x (km)
		21247.590022,  // y
		10246.812448,  // z
		-2.269448,     // vx (km/s)
		0.507755,      // vy
		-3.199603,     // vz
	}

	// Интервал интегрирования: 30 суток
	duration := 30 * 86400.0 // seconds
	dt := 60.0               // шаг 60 секунд — оптимально для J2

	tap.Message("Simulating 30 days with J2 gravity model...")

	states := simulateOrbit(state, dt, duration, j2Gravity)

	// Вывод ключевых точек
	displayResults(states, dt)

	// Генерация CSV-файлов
	err := exportOrbitCSV(states, dt)
	if err != nil {
		return err
	}

	tap.Message("CSV files generated: orbit_xyz.csv, radius.csv, speed.csv")

	cmd := exec.Command("pwd")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}

	cmd = exec.Command("./plot.py", "orbit_xyz.csv", "radius.csv", "speed.csv")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}

	return nil
}


func exportOrbitCSV(states [][]float64, dt float64) error {

	// 1) Орбита x,y,z
	f1, _ := os.Create("orbit_xyz.csv")
	w1 := csv.NewWriter(f1)
	w1.Write([]string{"t", "x", "y", "z"})
	for i, s := range states {
		t := float64(i) * dt
		w1.Write([]string{
			fmt.Sprintf("%.1f", t),
			fmt.Sprintf("%.6f", s[0]),
			fmt.Sprintf("%.6f", s[1]),
			fmt.Sprintf("%.6f", s[2]),
		})
	}
	w1.Flush()
	f1.Close()

	// 2) Радиус-вектор r(t)
	f2, _ := os.Create("radius.csv")
	w2 := csv.NewWriter(f2)
	w2.Write([]string{"t", "r"})
	for i, s := range states {
		t := float64(i) * dt
		r := math.Sqrt(s[0]*s[0] + s[1]*s[1] + s[2]*s[2])
		w2.Write([]string{
			fmt.Sprintf("%.1f", t),
			fmt.Sprintf("%.6f", r),
		})
	}
	w2.Flush()
	f2.Close()

	// 3) Модуль скорости |v|
	f3, _ := os.Create("speed.csv")
	w3 := csv.NewWriter(f3)
	w3.Write([]string{"t", "speed"})
	for i, s := range states {
		t := float64(i) * dt
		v := math.Sqrt(s[3]*s[3] + s[4]*s[4] + s[5]*s[5])
		w3.Write([]string{
			fmt.Sprintf("%.1f", t),
			fmt.Sprintf("%.6f", v),
		})
	}
	w3.Flush()
	f3.Close()

	return nil
}
