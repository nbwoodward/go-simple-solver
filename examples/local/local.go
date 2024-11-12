// This demonstrates finding local minima in an 8-dimensional function with 2 solutions in each
// dimension, 1 and -1.5. The minimize function runs one single test, and thus will return a non-optimal
// solution. The maximize function runs 100 tests and is likely to find the optimal solution.
package main

import (
	"fmt"
	"math"
	"time"

	simplegradient "github.com/nbwoodward/simplegradient"
)

var paramSpace = [][]float64{
	{-3, -2.5, -2, -1.5, -1, -0.5, 0, 0.5, 1, 1.5, 2, 2.5, 3},
	{-3, -2.5, -2, -1.5, -1, -0.5, 0, 0.5, 1, 1.5, 2, 2.5, 3},
	{-3, -2.5, -2, -1.5, -1, -0.5, 0, 0.5, 1, 1.5, 2, 2.5, 3},
	{-3, -2.5, -2, -1.5, -1, -0.5, 0, 0.5, 1, 1.5, 2, 2.5, 3},
	{-3, -2.5, -2, -1.5, -1, -0.5, 0, 0.5, 1, 1.5, 2, 2.5, 3},
	{-3, -2.5, -2, -1.5, -1, -0.5, 0, 0.5, 1, 1.5, 2, 2.5, 3},
	{-3, -2.5, -2, -1.5, -1, -0.5, 0, 0.5, 1, 1.5, 2, 2.5, 3},
	{-3, -2.5, -2, -1.5, -1, -0.5, 0, 0.5, 1, 1.5, 2, 2.5, 3},
}

// x^4 - 3x^2 + 2x
// This has local minimums at x=1 and x=-1.5. the global minimum is -1.5.
var poly = func(x float64) float64 {
	return math.Pow(x, 4) - 3*math.Pow(x, 2) + 2*x
}

func main() {
	minimize()
	maximize()
}

func minimize() {
	model := func(vals []float64) float64 {
		f := 0.0
		for _, val := range vals {
			f += poly(val)
		}

		return f
	}

	cfg := simplegradient.Config{
		Params:   paramSpace,
		Model:    model,
		NumTests: 1,
	}

	fmt.Println("Minimize Model:")
	start := time.Now()
	bestVal, bestParameters := simplegradient.Minimize(cfg)
	duration := time.Since(start)
	fmt.Println("Duration:", duration)
	fmt.Println("Best val:", bestVal)
	fmt.Println("Best params:", bestParameters)
}

func maximize() {
	model := func(vals []float64) float64 {
		f := 0.0
		for _, val := range vals {
			f -= poly(val)
		}

		return f
	}

	cfg := simplegradient.Config{
		Params:   paramSpace,
		Model:    model,
		NumTests: 100,
	}

	fmt.Println("")
	fmt.Println("Maximize Model:")
	start := time.Now()
	bestVal, bestParameters := simplegradient.Maximize(cfg)
	duration := time.Since(start)
	fmt.Println("Duration:", duration)
	fmt.Println("Best val:", bestVal)
	fmt.Println("Best params:", bestParameters)
}
