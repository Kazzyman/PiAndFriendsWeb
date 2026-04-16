package main

import (
	"fmt"
	"math"
	"math/big"
	"strconv"
	"time"
)

// @formatter:off

func MonteCarloWeb(gridSizeAsString string, webPrint func(string)) {
	
	// Produce an alternate string suitable for printing, with commas every three digits from the right
	withCommas := ""
	for i, char := range gridSizeAsString {
		if i > 0 && (len(gridSizeAsString)-i)%3 == 0 {
			withCommas += ","
		}
		withCommas += string(char)
	}
	// ::: screen
	webPrint(fmt.Sprintf("Size of the grid has been set to: %s", withCommas))

	// convert gridSize to an int
	gridSize, err := strconv.Atoi(gridSizeAsString)
	if err != nil {
		webPrint("Invalid input, please enter an integer in string form.")
		return
	}
		// ::: screen
		if gridSize < 5 {
			webPrint(" A grid that small makes me puke! ")
			return
		}
		webPrint(" ... working ... ")
		if gridSize > 3000 && gridSize <= 4000 {
			webPrint(" ... working ... expect 7s")
		} else if gridSize > 4000 && gridSize <= 5000 {
			webPrint(" ... working ... expect 10s")
		} else if gridSize > 6000 && gridSize <= 8500 {
			webPrint(" ... working ... expect 15s")
		} else if gridSize > 8500 && gridSize <= 11000 {
			webPrint(" ... really working ... expect 25s")
		} else if gridSize > 11000 && gridSize <= 12000 {
			webPrint(" ... I will be working really hard ...expect 30s")
		} else if gridSize > 12000 && gridSize <= 13000 {
			webPrint(" ... working really really hard...expect 40s")
		} else if gridSize > 13000 && gridSize <= 15000 {
			webPrint(" ... for so very long, I'll be working ...expect 50s")
		} else if gridSize > 15000 && gridSize <= 18000 {
			webPrint(" ... Yikes, I'll be working, for too long ...expect 1m5s")
		} else if gridSize > 18000 && gridSize <= 24000 {
			webPrint(" ... while you take a nap, I'll still be working ... expect 1m25s")
		} else if gridSize > 24000 && gridSize <= 34000 {
			webPrint(" ... Brace yourself for how long I'll be working ... expect 4min")
		} else if gridSize > 34000 && gridSize <= 100000 {
			webPrint(" ... Expect 5–15 minutes for ~4–5 digits ...")
			webPrint(" ... and be advised that 120k, or more, will make me puke! ...")
		} else if gridSize > 100000 && gridSize <= 119999 {
			webPrint(" ... Working insanely hard, expect 15–30 minutes for ~5 digits ...")
		} else if gridSize > 119999 {
			webPrint(" ... I have puked! ")
			return
		}
		
	piApprox := GridPi(gridSize, webPrint) // ::: run GridPi < - - - - - - - - - - < -

		// ::: screen
		webPrint(fmt.Sprintf("Size of the grid was set at: %s", withCommas))
		webPrint(fmt.Sprintf("Approximated Pi as big float: %s", piApprox.Text('f', 30)))
			piApproxFloat64, _ := piApprox.Float64()
		webPrint(fmt.Sprintf("Approximated Pi as float64:   %f", piApproxFloat64))
			piFromMathLib := math.Pi
			piFromMathLibBF := big.NewFloat(piFromMathLib)
		webPrint(fmt.Sprintf("Pi from Math Library:         %s", piFromMathLibBF.Text('f', 30)))
		webPrint(fmt.Sprintf("Difference: %f", math.Abs(piApproxFloat64-math.Pi)))
			_, digitCount := checkPiTo100(piApprox)
		webPrint(fmt.Sprintf("We verified Pi to %d digits", digitCount))
}
/*
.
 */
func GridPi(gridSize int, webPrint func(string)) *big.Float {
	start := time.Now()
		insideCircle := big.NewInt(0)
		totalPoints := big.NewInt(int64(gridSize * gridSize))
		increment := big.NewFloat(1.0 / float64(gridSize)).SetPrec(256)
		halfIncrement := new(big.Float).Quo(increment, big.NewFloat(2.0)).SetPrec(256)
	for i := 0; i < gridSize; i++ {
		/*
		case <-done: // ::: here an attempt is made to read from the channel (a closed channel can be read from successfully; but what is read will be the null/zero value of the type of chan (0, false, "", 0.0, etc.)
			// in the case of this particular channel (which is of type bool) we get the value false from having received from the channel when it is already closed. 
			// ::: if the channel known by the moniker "done" is already closed, that/it is to be interpreted as the abort signal by all listening processes. 
			fmt.Println("Goroutine Monty for-loop (2 of 2) is being terminated by select case finding the done channel to be already closed")
			return increment // Exit the goroutine ::: We had to return some kind of a big float ... 
			
		 */
		for j := 0; j < gridSize; j++ {
			// ::: x = (i * increment) + halfIncrement
				x := new(big.Float).SetPrec(256)
				x.Mul(increment, big.NewFloat(float64(i)))
				x.Add(x, halfIncrement)
			// ::: y = (j * increment) + halfIncrement
				y := new(big.Float).SetPrec(256)
				y.Mul(increment, big.NewFloat(float64(j)))
				y.Add(y, halfIncrement)
			xSquared := new(big.Float).Mul(x, x)
			ySquared := new(big.Float).Mul(y, y)
			sum := new(big.Float).Add(xSquared, ySquared)
				if sum.Cmp(big.NewFloat(1.0)) <= 0 {
					insideCircle.Add(insideCircle, big.NewInt(1))
				}
			iterationsForMonte16j = j
		}
		iterationsForMonte16i = i
		
	}
	iterationsForMonteTotal = iterationsForMonte16j * iterationsForMonte16i
		four := big.NewFloat(4.0).SetPrec(256)
		insideCircleF := new(big.Float).SetPrec(256).SetInt(insideCircle)
		totalPointsF := new(big.Float).SetPrec(256).SetInt(totalPoints)
		piApprox := new(big.Float).SetPrec(256)
		piApprox.Quo(insideCircleF, totalPointsF)
		piApprox.Mul(piApprox, four)
	t := time.Now()
	elapsed := t.Sub(start)
	TotalRun := elapsed.String()
		// ::: put commas into the Total-iterations var
		numStr := strconv.FormatInt(int64(iterationsForMonteTotal), 10)
		result := ""
		for i, char := range numStr {
			if i > 0 && (len(numStr)-i)%3 == 0 {
				result += ","
			}
			result += string(char)
		}
	// ::: screen	
	webPrint(fmt.Sprintf("Total iterations: %s Elapsed time: %s ", result, TotalRun))
	return piApprox
}
