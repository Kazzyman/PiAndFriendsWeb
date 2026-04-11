package main

import (
	"fmt"
	"os"
	"time"
)

// @formatter:off

func GregoryLeibniz(webPrint func(string), done chan bool) {
	// Open a log file 
	fileHandle, err1 := os.OpenFile("dataLog-From_calculate-pi-and-friends.txt", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600) // append to file
	check(err1)                              // ... gets a file handle to dataLog-From_calculate-pi-and-friends.txt
		defer func(fileHandle *os.File) {   // It’s idiomatic to defer a Close immediately after opening a file.
			err := fileHandle.Close()
			if err != nil {}
		}(fileHandle)

	usingBigFloats = false // π = (4/1) - (4/3) + (4/5) - (4/7) + (4/9) - (4/11) + (4/13) - (4/15) ...
		webPrint("You selected Gregory-Leibniz formula  :  π = 4 * ( 1 - 1/3 + 1/5 - 1/7 + 1/9 ...) ")
		webPrint("   Infinitesimal calculus was developed independently in the late 17th century by Isaac Newton")
		webPrint("    ... James Gregory, and Gottfried Wilhelm Leibniz")
		webPrint("   4 Billion iterations will (initially) be executed ... ")
		webPrint(" ... working ...")

	start := time.Now()

	var denom float64
	var sum float64
	denom = 3
	sum = 1 - (1 / denom)
	
	iterInt64 = 1   // global 
	iterFloat64 = 0 // global
	
	for iterInt64 < 4000000000 {
		select {
		case <-done: // ::: here an attempt is made to read from the channel (a closed channel can be read from successfully; but what is read will be the null/zero value of the type of chan (0, false, "", 0.0, etc.)
			// in the case of this particular channel (which is of type bool) we get the value false from having received from the channel when it is already closed. 
			// ::: if the channel known by the moniker "done" is already closed, that/it is to be interpreted as the abort signal by all listening processes. 
			fmt.Println("Goroutine Gregory-Leibniz for-loop (1 of 2) is being terminated by select case finding the done channel to be already closed")
			return // Exit the goroutine
		default:
		iterFloat64++
		iterInt64++
		
		denom = denom + 2
		
		if iterInt64%2 == 0 {
			sum = sum + 1/denom
		} else {
			sum = sum - 1/denom
		}
		
		π = 4 * sum // calculate ::: pi : π
		
			if iterInt64 == 100000000 {
			webPrint("... 100,000,000 completed iterations ...")
			webPrint(fmt.Sprintf("   %0.7f was calculated by the Gottfried Wilhelm Leibniz formula", π))
			t := time.Now()
			elapsed := t.Sub(start)
			webPrint(fmt.Sprintf("  100,000,000 iterations in %s yields 8 digits of π", elapsed))
			}
			if iterInt64 == 200000000 {
			webPrint("... 200,000,000 gets another digit ...")
			webPrint(fmt.Sprintf("   %0.9f was calculated by the Gottfried Wilhelm Leibniz formula", π))
			t := time.Now()
			elapsed := t.Sub(start)
			webPrint(fmt.Sprintf("  200,000,000 iterations in %s yields 9 digits of π", elapsed))
			}
			if iterInt64 == 400000000 {
			webPrint("... 400,000,000 iterations completed, still at nine ...")
			webPrint(fmt.Sprintf("   %0.10f was calculated by the Gottfried Wilhelm Leibniz formula", π))
			t := time.Now()
			elapsed := t.Sub(start)
			webPrint(fmt.Sprintf("  400,000,000 iterations in %s yields 9 digits of π", elapsed))
			}
			if iterInt64 == 600000000 {
			webPrint("... 600,000,000 iterations, still at nine ...")
			webPrint(fmt.Sprintf("   %0.5f was calculated by the Gottfried Wilhelm Leibniz formula", π))
			t := time.Now()
			elapsed := t.Sub(start)
			webPrint(fmt.Sprintf("  600,000,000 iterations in %s yields 9 digits of π", elapsed))
			}
			if iterInt64 == 1000000000 {
			webPrint("... 1 Billion iterations completed, still nine ...")
			webPrint(fmt.Sprintf("   %0.5f was calculated by the Gottfried Wilhelm Leibniz formula", π))
			t := time.Now()
			elapsed := t.Sub(start)
			webPrint(fmt.Sprintf("  1,000,000,000 iterations in %s yields 9 digits of π", elapsed))
			}
			if iterInt64 == 2000000000 {
			webPrint("... 2 Billion, and still just nine ...")
			webPrint(fmt.Sprintf("   %0.5f was calculated by the Gottfried Wilhelm Leibniz formula", π))
			t := time.Now()
			elapsed := t.Sub(start)
			webPrint(fmt.Sprintf("  2,000,000,000 iterations in %s yields 9 digits of π", elapsed))
			}
			if iterInt64 == 4000000000 { // ::: last one
				webPrint("... 4 Billion, gets us ten digits  ...")
				webPrint(fmt.Sprintf("   %0.5f was calculated by the Gottfried Wilhelm Leibniz formula", π))
				t := time.Now()
				elapsed := t.Sub(start)
				webPrint(fmt.Sprintf("  4,000,000,000 iterations in %s yields 10 digits of π", elapsed))
				webPrint(" per the Gottfried Wilhelm Leibniz formula")
		
				LinesPerIter = 14
				webPrint(fmt.Sprintf("at aprox %0.2f lines of code per iteration ...", LinesPerIter))
				LinesPerSecond = (LinesPerIter * iterFloat64) / elapsed.Seconds() // .Seconds() returns a float64
				webPrint(fmt.Sprintf("Aprox %.0f lines of code were executed per second ", LinesPerSecond))
				
				// store results in a log file 
					Hostname, _ := os.Hostname()
					current_time := time.Now()
					TotalRun := elapsed.String()   // cast time duration to a String type for Fprintf "formatted print"
				
					// to ::: file
						_, err0 := fmt.Fprintf(fileHandle, "  -- Gottfried Wilhelm Leibniz --  on %s ", Hostname)
							check(err0)
						_, err6 := fmt.Fprint(fileHandle, "was run on: ", current_time.Format(time.ANSIC), "")
							check(err6)
						_, err2 := fmt.Fprintf(fileHandle, "%.0f was Lines/Second  ", LinesPerSecond)
							check(err2)
						_, err4 := fmt.Fprintf(fileHandle, "%e was Iterations/Seconds  ", iterFloat64/elapsed.Seconds())
							check(err4)
						_, err5 := fmt.Fprintf(fileHandle, "%e was total Iterations  ", iterFloat64)
							check(err5)
						_, err7 := fmt.Fprintf(fileHandle, "Total runTime was %s ", TotalRun) // add total runtime of this calculation
							check(err7)
			} // end of last if
			}
	} // end of first for loop
	

		// print to ::: file
			webPrint( "We continue the Gottfried Wilhelm Leibniz formula  :  π = 4 * ( 1 - 1/3 + 1/5 - 1/7 + 1/9 ... ")
			webPrint("    π = 3 + 4/(2*3*4) - 4/(4*5*6) + 4/(6*7*8) - 4/(8*9*10) + 4/(10*11*12) ...")
			
			webPrint("   Infinitesimal calculus was developed independently in the late 17th century by Isaac Newton")
			webPrint("   and Gottfried Wilhelm Leibniz")
			webPrint("   9 billion iterations will be executed    ... working ...")

		start = time.Now()

		for iterInt64 < 9000000000 {
			select {
			case <-done: // ::: here an attempt is made to read from the channel (a closed channel can be read from successfully; but what is read will be the null/zero value of the type of chan (0, false, "", 0.0, etc.)
				// in the case of this particular channel (which is of type bool) we get the value false from having received from the channel when it is already closed. 
				// ::: if the channel known by the moniker "done" is already closed, that/it is to be interpreted as the abort signal by all listening processes. 
				fmt.Println("Goroutine Gregory-Leibniz for-loop (2 of 2) is being terminated by select case finding the done channel to be already closed")
				return // Exit the goroutine
			default:
			iterFloat64++
			iterInt64++
			denom = denom + 2
			if iterInt64%2 == 0 {
			sum = sum + 1/denom
			} else {
			sum = sum - 1/denom
			}
			π = 4 * sum
			
			if iterInt64 == 6000000000 {
			webPrint("... 6 Billion completed ... ")
			webPrint(fmt.Sprintf("   %0.13f was calculated by the Gottfried Wilhelm Leibniz formula", π))
			t := time.Now()
			elapsed := t.Sub(start)
			webPrint(fmt.Sprintf("  6,000,000,000 iterations in %s still yields 10 digits of π", elapsed))
			webPrint( "  ... working ...")
			}
			if iterInt64 == 8000000000 {
			webPrint("... 8 Billion completed. still ten ...")
			webPrint(fmt.Sprintf("   %0.13f was calculated by the Gottfried Wilhelm Leibniz formula", π))
			t := time.Now()
			elapsed := t.Sub(start)
			webPrint(fmt.Sprintf("  8,000,000,000 iterations in %s still yields 10 digits of π", elapsed))
			webPrint( " ... working ...")
			}
			if iterInt64 == 9000000000 {
			webPrint(fmt.Sprintf("   %0.13f was calculated by the Gottfried Wilhelm Leibniz formula", π))
			// webPrint(fmt.Sprintf("   ", iter)
			t := time.Now()
			elapsed := t.Sub(start)
			webPrint(fmt.Sprintf("... 9B iterations in %s, but to get 10 digits we only needed 4B iterations", elapsed))
			webPrint(" per  --  the Gottfried Wilhelm Leibniz formula")
			
				t = time.Now()
				elapsed = t.Sub(start)
				TotalRun := elapsed.String()          // cast time duration to a String type for Fprintf "formatted print"
			
			LinesPerIter = 14 // estimate
			// print to ::: screen
				webPrint(fmt.Sprintf("at aprox %0.2f lines of code per iteration ...", LinesPerIter))
			
					webPrint(fmt.Sprintf("%e was Iterations/Seconds", iterFloat64/elapsed.Seconds()))
					webPrint("   Infinitesimal calculus was developed independently in the late 17th century by Isaac Newton")
					webPrint("   and Gottfried Wilhelm Leibniz. This implementaion was done entirely by Richard Woolley")
			
				LinesPerSecond = (LinesPerIter * iterFloat64) / elapsed.Seconds() // .Seconds() returns a float64
			// to ::: screen
				webPrint(fmt.Sprintf("Aprox %.0f lines of code were executed per second ", LinesPerSecond))
				webPrint(fmt.Sprintf("Total runTime was %s ", TotalRun)) // add total runtime of this calculation
			}
			}
		}
		
	// ::: Prepare to exit the Gottfried method function
	webPrint("Gregory-Leibniz calculation complete.")
	done <- true
	
} // written entirely by Richard Woolley
