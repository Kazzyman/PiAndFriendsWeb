package main

import (
	"fmt"
	"time"
)

// @formatter:off

func CustomSeries(webPrint func(string)) {
	usingBigFloats = false
	
	webPrint("You selected a Custom Recursive π Approximation ... this will be quick!")
	/*
	The code’s sequence (4, 3.5556, 3.4134, 3.3438, …) doesn’t match Leibniz (4, 2.6667, 3.4667, 2.8952, …), but it still converges toward π over many iterations.
	*/
	webPrint("Three-hundred-million iterations will be executed ... working ...")

	var nextOdd float64
	var tally float64
	
			start := time.Now()
			
	iterFloat64 = 0
	nextOdd = 1
	four = 4
	tally = (four / nextOdd)
	iterInt64 = 0
	/*
	Pi = 4/1 + 4/3 + 4/5
	 */
	
	for iterInt64 < 300000000 {

		iterInt64++
		iterFloat64++
		nextOdd = nextOdd + 2
		tally = tally - (tally / nextOdd)
		tally = tally + (tally / nextOdd) // pi (tally) is set equl to the sum of a subtraction and an addition, alternatively

		if iterInt64 == 10000000 {
			webPrint("... 10,000,000 of three hundred million iterations already completed. still working, but ...")
			webPrint(fmt.Sprintf("   %0.6f was calculated thus far via the Gregory-Leibniz series", tally))
				t := time.Now()
				elapsed := t.Sub(start)
			webPrint(fmt.Sprintf("  10,000,000 iterations in %s yields 7 digits of π", elapsed))
		}
		// 7 digits of Pi were found per the above code
		// the next two ifs give eight digits of Pi
		if iterInt64 == 50000000 {
			webPrint("... 50,000,000 of three hundred million completed. still working, but ...")
			webPrint(fmt.Sprintf("      %0.8f was calculated by the Gregory-Leibniz series, so far", tally))
				t := time.Now()
				elapsed := t.Sub(start)
			webPrint(fmt.Sprintf("  50,000,000 iterations in %s yields 8 digits of π", elapsed))
		}
		if iterInt64 == 100000000 {
			webPrint("... 100,000,000 of three hundred million completed. still working, and ...")
			webPrint(fmt.Sprintf("      %0.9f was calculated by the Gregory-Leibniz series", tally))
				t := time.Now()
				elapsed := t.Sub(start)
			webPrint(fmt.Sprintf("  100,000,000 iterations in %s yields 8 digits of π", elapsed))
		}
		// 9 digits of Pi are found below
		if iterInt64 == 200000000 {
			webPrint("... 200,000,000 of three hundred million now completed. still working, but ...")
			webPrint(fmt.Sprintf("      %0.10f was calculated thus far by the Gregory-Leibniz series", tally))
				t := time.Now()
				elapsed := t.Sub(start)
			webPrint(fmt.Sprintf("  200,000,000 iterations in %s yields 9 digits of π", elapsed))
		}
		if iterInt64 == 300000000 { // last one, still 9 digits
			webPrint(fmt.Sprintf("       %0.11f was calculated by the Gregory-Leibniz series ", tally))
				t := time.Now()
				elapsed := t.Sub(start)
			webPrint("  300 million iterations have finally finished; still yielding only 9 digits of pi, ") // no Println here
			webPrint(fmt.Sprintf("in %s", elapsed))
			webPrint(" per the Gregory-Leibniz series, circa 1676")

				LinesPerIter = 11 // an estimate of the number of lines per iteration
				linePerApp := LinesPerIter * 300000000
				stringOfTotal := formatFloat64WithThousandSeparators(linePerApp)
				LinesPerSecond = (LinesPerIter * iterFloat64) / elapsed.Seconds() // .Seconds() returns a float64
			webPrint(fmt.Sprintf("at aprox %0.0f lines of code per iteration ... SLOC executed was aprox. %s ", LinesPerIter, stringOfTotal))
			webPrint(fmt.Sprintf("       %.0f lines of code were executed per second ", LinesPerSecond))

					TotalRun := elapsed.String() // cast time durations to a String type for Fprintf "formatted print"
			
			webPrint(" That was the Gregory-Leibniz series:")
			webPrint("π = (4/1) - (4/3) + (4/5) - (4/7) + (4/9) - (4/11) + (4/13) - (4/15) ...")
			webPrint(fmt.Sprintf("Runtime was: %s", TotalRun))
			webPrint("Three-hundred-million iterations were executed. This section was written entirely by Richard Woolley")

		}
	
	}
	// ::: Prepare to exit the Gregory Leibniz method function
	calculating = false // Allow another method to be selected.
	/*
	for _, btn := range buttons1 { // ok to only Enable buttons1, because I expect to only ever execute this from window1
		btn.Enable() // ::: Enable
	}
	
	 */
	// written entirely by Richard Woolley
}