package simplegradient

import (
	"fmt"
	"math"
	"math/rand"
	"time"
)

var DEFAULT_NUM_TESTS = 10
var DEFAULT_MAX_MODEL_ITERATION = 100
var VERBOSE = false

type Config struct {
	// The parameter space for the model.
	Params [][]float64

	// The model to minimize or maximize.
	Model func([]float64) float64

	// The number of concurrent tests to run. The more tests
	// the higher liklihood of finding the global min or max.
	NumTests int

	// The maximum interations for each gradient descent test to run.
	// TThe higher the iterations the higher
	MaxIteration int

	// The length of each param array.
	// Prevents having to calculate this repeatedly throughout the model.
	paramLens []int
}

type solution struct {
	soln float64
	vals []float64
}

type empty struct{}
type done chan empty

func init() {
	rand.Seed(time.Now().Unix())
}

func Minimize(c Config) (float64, []float64) {
	solutions := c.solveN(true)

	bestSoln := solution{math.Inf(1), []float64{}}
	for _, solution := range solutions {
		if solution.soln < bestSoln.soln {
			bestSoln = solution
		}
	}

	return bestSoln.soln, bestSoln.vals
}

func Maximize(c Config) (float64, []float64) {
	solutions := c.solveN(false)

	bestSoln := solution{math.Inf(-1), []float64{}}
	for _, solution := range solutions {
		if solution.soln > bestSoln.soln {
			bestSoln = solution
		}
	}

	return bestSoln.soln, bestSoln.vals
}

/*
 * This runs each gradient descent in its own goroutine.
 *
 * Adapted from https://gist.github.com/danieldk/3810803
 * This could be improved to pass each solution back using a channel to
 * prevent the `solutions` array from getting copied into each routine.
 */
func (c *Config) solveN(minimize bool) []solution {
	c.setParamLens()

	if c.NumTests <= 0 {
		c.NumTests = DEFAULT_NUM_TESTS
	}

	if c.MaxIteration <= 0 {
		c.MaxIteration = DEFAULT_MAX_MODEL_ITERATION
	}

	solutions := make([]solution, c.NumTests)
	allDone := make(done, c.NumTests)

	for i := 0; i < c.NumTests; i++ {
		go func(idx int) {
			soln, vals := c.solveOne(minimize)
			solutions[idx] = solution{soln, vals}
			allDone <- empty{}
		}(i)
	}

	for i := 0; i < c.NumTests; i++ {
		<-allDone
	}

	return solutions
}

func (c *Config) solveOne(minimize bool) (float64, []float64) {
	idxs := c.getRandomIndexes()
	vals := c.valsFromIndexes(idxs)
	soln := c.Model(vals)
	bestIdxs := c.followGradient(idxs, soln, minimize, 0)
	bestVals := c.valsFromIndexes(bestIdxs)

	return c.Model(bestVals), bestVals
}

func (c *Config) followGradient(idxs []int, currentSoln float64, minimize bool, count int) []int {

	if count > c.MaxIteration {
		fmt.Println("Reached model max iteration")
		return idxs
	}

	nextIdxs := append([]int{}, idxs...)
	didUpdate := false

	for i, idx := range idxs {
		testIdxs := append([]int{}, idxs...)
		testVals := make([]float64, len(idxs))
		var upIdx int
		var downIdx int
		var upSoln float64
		var downSoln float64
		var nextIdx int

		// If we are at the ends of our arrays for this parameter, move inwards.
		if c.paramLens[i] == 1 {
			downIdx = 0
			upIdx = 0
			downSoln = currentSoln
			upSoln = currentSoln
			testIdxs[i] = 0
		} else if idx == 0 {
			downIdx = 0
			upIdx = 1
			downSoln = currentSoln
			testIdxs[i] = 1
			testVals = c.valsFromIndexes(testIdxs)
			upSoln = c.Model(testVals)
		} else if idx == c.paramLens[i]-1 {
			downIdx = idx - 1
			upIdx = idx
			upSoln = currentSoln
			testIdxs[i] = idx - 1
			testVals = c.valsFromIndexes(testIdxs)
			downSoln = c.Model(testVals)
		} else {
			downIdx = idx - 1

			testIdxs[i] = downIdx
			testVals = c.valsFromIndexes(testIdxs)
			downSoln = c.Model(testVals)

			upIdx = idx + 1
			testIdxs[i] = upIdx
			testVals = c.valsFromIndexes(testIdxs)
			upSoln = c.Model(testVals)
		}

		if minimize {
			// If we are minimizing we follow the gradient down
			if upSoln < currentSoln {
				nextIdx = upIdx
			} else if downSoln < currentSoln {
				nextIdx = downIdx
			} else {
				nextIdx = idx
			}
		} else {
			// Or else we are maximizing
			if upSoln > currentSoln {
				nextIdx = upIdx
			} else if downSoln > currentSoln {
				nextIdx = downIdx
			} else {
				nextIdx = idx
			}
		}

		if VERBOSE {
			fmt.Printf("i: %v, idx: %v, currentSoln: %v, downIdx: %v, downSoln: %v, upIdx: %v, upSoln: %v, nextIdx: %v\n",
				i, idx, currentSoln, downIdx, downSoln, upIdx, upSoln, nextIdx)
		}

		if nextIdx != idx {
			didUpdate = true
		}

		nextIdxs[i] = nextIdx
	}

	nextVals := c.valsFromIndexes(nextIdxs)
	nextSoln := c.Model(nextVals)

	if VERBOSE {
		fmt.Printf("Iteration: %v, lastIdxs: %v, nextIdxs: %v, current: %v, next %v\n",
			count, idxs, nextIdxs, currentSoln, nextSoln)
	}

	if !didUpdate || (minimize && currentSoln < nextSoln) || (!minimize && currentSoln > nextSoln) {
		return idxs
	}

	return c.followGradient(nextIdxs, nextSoln, minimize, count+1)
}

func (c *Config) getRandomIndexes() []int {
	idxs := make([]int, len(c.Params))

	for i, vals := range c.Params {
		idxs[i] = rand.Intn(len(vals))
	}

	return idxs
}

func (c *Config) valsFromIndexes(idxs []int) []float64 {
	vals := make([]float64, len(c.Params))
	for i, param := range c.Params {
		vals[i] = param[idxs[i]]
	}

	return vals
}

func (c *Config) setParamLens() {
	c.paramLens = make([]int, len(c.Params))
	for i, vals := range c.Params {
		c.paramLens[i] = len(vals)
	}
}
