# Simple Gradient Descent

Do you have a mostly continuous numerical model with:
- Too big of a parameter space to do an exhaustive search?
- Small enough of a parameter space that you don't need big ol' GPUs?

You're in luck! To maximize or minimize your model all you need is:
 - A 2D array of `float64`s that covers your parameter space.
 - A function that outputs a `float64`.
 - Run `Maximize` or `Minimize`
 - Enjoy sweet freedom

Example Usage:
```go
package main

import (
	"fmt"
	"math"

	simplegradient "github.com/nbwoodward/simplegradient"
)

// All the params we think might be optimal for our function
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

	// A simple model that
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

	bestParams, bestAnswer := simplegradient.Minimize(cfg)
	fmt.Println(bestParams, bestAnswer)
	// Prints: 6, [0, 0, 0]

	// Solve again with verbose mode on
	cfg.NumTests = 1
	simplegradient.Verbose = true

	bestParams, bestAnswer = simplegradient.Minimize(cfg)
	// Prints more information about the process
}
```

## The Strategy
The solver takes the standard numerical model approach of:
- Picking N random starting locations in the param space (from the config value `NumTests`)
- Testing the index above and below each parameter to see if it performs better
- Setting the best parameters as the next starting point, trying again, and continuing to move up or down hill.
- Taking the best value and parameter set from each of the N tests.

Each of the N tests runs on its own goroutine for efficent use of system resources.

## Improvements
- The solution matrix is copied to each goroutine which makes it less memory efficient. Channels
should be used to pass the information around more cleanly.
- Parameters that have already been tested should have their values cached so we aren't duplicating work.
