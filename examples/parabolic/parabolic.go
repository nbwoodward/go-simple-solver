// This demonstrates finding the minimum of a simple 3-dimensional paraboloid.
package main

import (
	"fmt"
	"math"

	simplegradient "github.com/nbwoodward/simplegradient"
)

var paramSpace = [][]float64{
	{-3, -2.5, -2, -1.5, -1, -0.5, 0, 0.5, 1, 1.5, 2, 2.5, 3}, // X
	{-3, -2.5, -2, -1.5, -1, -0.5, 0, 0.5, 1, 1.5, 2, 2.5, 3}, // Y
	{-3, -2.5, -2, -1.5, -1, -0.5, 0, 0.5, 1, 1.5, 2, 2.5, 3}, // Z
}

// x^2 + 2
var parabola = func(x float64) float64 {
	return math.Pow(x, 2) + 2
}

func main() {
	model := func(xyz []float64) float64 {
		x := xyz[0]
		y := xyz[1]
		z := xyz[2]

		return parabola(x) + parabola(y) + parabola(z)
	}

	cfg := simplegradient.Config{
		Model:    model,
		Params:   paramSpace,
		NumTests: 10,
	}

	bestAnswer, bestParams := simplegradient.Minimize(cfg)
	fmt.Println("Best Parameters:", bestParams)
	fmt.Println("Minimum Solution:", bestAnswer)
}
