package main

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

// @formatter:off

func JohnWallis(done chan bool, webPrint func(string)) float64 { // case 8: // -- AMFJohnWallisA
webPrint("I am here in JW")
// ::: it makes it to here before hanging

	webPrint("The forgoing is the entire code for this method.")

	usingBigFloats = false
	webPrint("   You selected A Go language exercize which can be used to test the speed of your hardware.")
	webPrint("   We will calculate π to a maximum of ten digits of accuracy using an infinite series by John Wallis circa 1655")
	webPrint("   Up to 40 Billion iterations of the following formula will be executed ")
	webPrint("   π = 2 * ((2/1)*(2/3)) * ((4/3)*(4/5)) * ((6/5)*(6/7)) ...")
	start := time.Now()
	iterFloat64 = 0
	var numerators float64
	numerators = 2
	var firstDenom float64
	firstDenom = 1
	var secondDenom float64
	secondDenom = 3
	var cumulativeProduct float64
	cumulativeProduct = (numerators / firstDenom) * (numerators / secondDenom)
	iterInt64 = 0
	// Wallis one:
	for iterInt64 < 1000000000 { // was 1000000000
		select {
		case <-done: // ::: here an attempt is made to read from the channel (a closed channel can be read from successfully; but what is read will be the null/zero value of the type of chan (0, false, "", 0.0, etc.)
			// in the case of this particular channel (which is of type bool) we get the value false from having received from the channel when it is already closed. 
			// ::: if the channel known by the moniker "done" is already closed, that/it is to be interpreted as the abort signal by all listening processes. 
			webPrint("Goroutine Wallis for-loop (1 of 2) is being terminated by select case finding the done channel to be already closed")
			return π// Exit the goroutine
		default:
			iterInt64++
			iterFloat64++
			numerators = numerators + 2
			firstDenom = firstDenom + 2
			secondDenom = secondDenom + 2
			cumulativeProduct = cumulativeProduct * (numerators / firstDenom) * (numerators / secondDenom)
			π = cumulativeProduct * 2

			if iterInt64 == 2000 {
				webPrint(fmt.Sprintf("%0.5f calculated using an infinite series by John Wallis circa 1655", π))
				webPrint("3.14159265358  is, again, the value of π from the web")
				t := time.Now()
				elapsed := t.Sub(start)
				RunTimeAsString := elapsed.String()
				webPrint(fmt.Sprintf("%d iterations were completed in %s yielding 4 digits of π", iterInt64, RunTimeAsString))
			}
			if iterInt64 == 10000 {
				webPrint(fmt.Sprintf("%0.6f calculated using an infinite series by John Wallis circa 1655", π))
				webPrint("3.14159265358  is, again, the value of π from the web")
				t := time.Now()
				elapsed := t.Sub(start)
				RunTimeAsString := elapsed.String()
				webPrint(fmt.Sprintf("10,000 iterations were completed in %s yielding 5 digits of π", RunTimeAsString))
			}
			if iterInt64 == 50000 { // 50,000
				webPrint(fmt.Sprintf("%0.7f calculated using an infinite series by John Wallis circa 1655", π))
				webPrint("3.14159265358  is, again, the value of π from the web")
				t := time.Now()
				elapsed := t.Sub(start)
				RunTimeAsString := elapsed.String()
				webPrint(fmt.Sprintf("50,000 iterations were completed in %s yielding 5 digits of π", RunTimeAsString))
			}
			if iterInt64 == 500000 { // 500,000 done
				webPrint(fmt.Sprintf("%0.8f calculated using an infinite series by John Wallis circa 1655", π))
				webPrint("3.14159265358  is, again, the value of π from the web")
				t := time.Now()
				elapsed := t.Sub(start)
				RunTimeAsString := elapsed.String()
				webPrint(fmt.Sprintf("500,000 iterations were completed in %s yielding 6 digits of π", RunTimeAsString))
			}
			if iterInt64 == 2000000 { // 2M done
				webPrint(fmt.Sprintf("%0.9f calculated using an infinite series by John Wallis circa 1655", π))
				webPrint("3.14159265358  is, again, the value of π from the web")
				t := time.Now()
				elapsed := t.Sub(start)
				RunTimeAsString := elapsed.String()
				webPrint(fmt.Sprintf("2,000,000 iterations were completed in %s yielding 7 digits of π", RunTimeAsString))
			}
			if iterInt64 == 40000000 { // 40M done
				webPrint(fmt.Sprintf("%0.10f calculated using an infinite series by John Wallis circa 1655", π))
				webPrint("3.14159265358  is, again, the value of π from the web")
				t := time.Now()
				elapsed := t.Sub(start)
				RunTimeAsString := elapsed.String()
	
				piAsAString := strconv.FormatFloat(π, 'g', -1, 64)
				checkPiUpTo255chars(piAsAString)
				webPrint(fmt.Sprintf("40,000,000 iterations were completed in %s yielding %d confirmed digits of π", RunTimeAsString, copyOfLastPosition))
				webPrint("  .. working .. on another factor-of-ten iterations")
			}
			if iterInt64 == 400000000 { // 400M done
				webPrint(fmt.Sprintf("%0.11f calculated using an infinite series by John Wallis circa 1655", π))
				webPrint("3.14159265358  is, again, the value of π from the web")
				t := time.Now()
				elapsed := t.Sub(start)
				RunTimeAsString := elapsed.String()
	
				piAsAString := strconv.FormatFloat(π, 'g', -1, 64)
				checkPiUpTo255chars(piAsAString)
	
				webPrint(fmt.Sprintf("400,000,000 iterations were completed in %s yielding %d confirmed digits of π", RunTimeAsString, copyOfLastPosition))
	
				LinesPerIter = 36 // an estimate
				webPrint(fmt.Sprintf("at aprox %0.1f lines of code per iteration ...", LinesPerIter))
				LinesPerSecond = (LinesPerIter * iterFloat64) / elapsed.Seconds()
				formattedLinesPerSecond := formatInt64WithThousandSeparators(int64(LinesPerSecond)) // .Seconds() returns a float64
				webPrint(fmt.Sprintf("Aprox %s lines of code were executed per second ", formattedLinesPerSecond))
				// a brief Red notification follows :
				webPrint(" ... will be working on doing Billions more iterations ...")
			}
			//
			if iterInt64 == 600000000 { // 600M done
				webPrint("  600M done, still working on another Two-Hundred-Thousand iterations ... working ...")
				t := time.Now()
				elapsed := t.Sub(start)
				RunTimeAsString := elapsed.String()
				webPrint(fmt.Sprintf("%s ", RunTimeAsString))
				webPrint("Calculating the next digit of pi may require 40B iterations, which takes a few minutes ")
				LinesPerIter = 36 // an estimate
				webPrint(fmt.Sprintf("at aprox %0.1f lines of code per iteration ...", LinesPerIter))
				LinesPerSecond = (LinesPerIter * iterFloat64) / elapsed.Seconds()
				formattedLinesPerSecond := formatInt64WithThousandSeparators(int64(LinesPerSecond)) // .Seconds() returns a float64
				webPrint(fmt.Sprintf("Aprox %s lines of code were executed per second ", formattedLinesPerSecond))
				webPrint(" ... still working ...")
			}
			if iterInt64 == 800000000 { // 800M done
				webPrint("  800M done, still working on yet another Two Hundred Thousand iterations ... working ...")
				t := time.Now()
				elapsed := t.Sub(start)
				RunTimeAsString := elapsed.String()
				webPrint(fmt.Sprintf("%s ", RunTimeAsString))
			}
			if iterInt64 == 1000000000 { // 1B done
				webPrint(fmt.Sprintf("%0.11f calculated using an infinite series by John Wallis circa 1655", π))
				webPrint("3.14159265358  is the value of π from the web")
				t := time.Now()
				elapsed := t.Sub(start)
				RunTimeAsString := elapsed.String()
	
				piAsAString := strconv.FormatFloat(π, 'g', -1, 64)
				checkPiUpTo255chars(piAsAString)
				webPrint(fmt.Sprintf("One Billion iterations were completed in %s still only yielding π to %d confirmed digits", RunTimeAsString, copyOfLastPosition))
				webPrint(" per --  an infinite series by John Wallis circa 1655") // ----------------------
	
				LinesPerIter = 36 // an estimate
				webPrint(fmt.Sprintf("at aprox %0.1f lines of code per iteration ...", LinesPerIter))
				LinesPerSecond = (LinesPerIter * iterFloat64) / elapsed.Seconds()
				formattedLinesPerSecond := formatInt64WithThousandSeparators(int64(LinesPerSecond)) // .Seconds() returns a float64
				webPrint(fmt.Sprintf("Aprox %s lines of code were executed per second ", formattedLinesPerSecond))
	
				// store reults in a log file which can be displayed from within the program by selecting option #12
				fileHandle, err1 := os.OpenFile("dataLog-From_calculate-pi-and-friends.txt", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600) // append to file
				check(err1)                                                                                                             // ... gets a file handle to dataLog-From_calculate-pi-and-friends.txt
				defer fileHandle.Close()                                                                                                // It’s idiomatic to defer a Close immediately after opening a file.
				Hostname, _ := os.Hostname()
				_, err0 := fmt.Fprintf(fileHandle, "  -- John Wallis -- on %s ", Hostname)
				check(err0)
				current_time := time.Now()
				_, err6 := fmt.Fprint(fileHandle, "was run on: ", current_time.Format(time.ANSIC), "")
				check(err6)
				_, err2 := fmt.Fprintf(fileHandle, "%s was Lines/Second  ", formattedLinesPerSecond)
				check(err2)
				_, err4 := fmt.Fprintf(fileHandle, "%.02f was Iterations/Seconds  ", iterFloat64/elapsed.Seconds())
				check(err4)
				_, err5 := fmt.Fprintf(fileHandle, "%e was total Iterations  ", iterFloat64)
				check(err5)
				TotalRun := elapsed.String()                                         // cast time durations to a String type for Fprintf "formatted print"
				_, err7 := fmt.Fprintf(fileHandle, "Total run was %s  ", TotalRun) // add total runtime of this calculation
				check(err7)
			} // ifs
		} // select
	} // end of first for loop

// :::webPrint(fmt.Sprintf("Enter any positive digit to continue with an additional 39 billion iterations, 0 to exit")

	webPrint("You elected to continue the infinite series by John Wallis")
	webPrint("    an additionl 39 billion iterations will be executed    ... working ...")

	webPrint(" ... still working ... on Billions of iterations, 39 to go ...")

	webPrint(" ... 39 Billion additional loops now ensue, just to get maybe one additional digit of pi")

	start = time.Now()

	// Wallis two:
	for iterInt64 < 40000000000 { // was 40000000000
		select {
		case <-done: // ::: here an attempt is made to read from the channel (a closed channel can be read from successfully; but what is read will be the null/zero value of the type of chan (0, false, "", 0.0, etc.)
			// in the case of this particular channel (which is of type bool) we get the value false from having received from the channel when it is already closed. 
			// ::: if the channel known by the moniker "done" is already closed, that/it is to be interpreted as the abort signal by all listening processes. 
			webPrint("Goroutine Wallis for-loop (2 of 2) is being terminated by select case finding the done channel to be already closed")
			return π// Exit the goroutine
		default:
				iterInt64++
				iterFloat64++
				numerators = numerators + 2
				firstDenom = firstDenom + 2
				secondDenom = secondDenom + 2
				cumulativeProduct = cumulativeProduct * (numerators / firstDenom) * (numerators / secondDenom)
				π = cumulativeProduct * 2
	
				if iterInt64 == 2000000000 { // 2B completed
					webPrint("  2B done, still working ... on another Billion iterations ... working ...")
						t := time.Now()
						elapsed := t.Sub(start)
					RunTimeAsString := elapsed.String()
					webPrint(fmt.Sprintf("%s", RunTimeAsString))
				}
				if iterInt64 == 3000000000 { // 3B completed
					webPrint("  3B done, still working ... on another Billion iterations ... working ... ")
						t := time.Now()
						elapsed := t.Sub(start)
					RunTimeAsString := elapsed.String()
					webPrint(fmt.Sprintf("%s", RunTimeAsString))
				}
				if iterInt64 == 4000000000 { // 4B completed
					webPrint("  4B done, still working ... on another Billion iterations ... working ...")
						t := time.Now()
						elapsed := t.Sub(start)
					RunTimeAsString := elapsed.String()
					webPrint(fmt.Sprintf("%s", RunTimeAsString))
				}
				if iterInt64 == 5000000000 { // 5B completed
					webPrint("  5B done, still working ... on another Billion iterations ... working ...")
						t := time.Now()
						elapsed := t.Sub(start)
					RunTimeAsString := elapsed.String()
					webPrint(fmt.Sprintf("%s", RunTimeAsString))
				}
				if iterInt64 == 6000000000 { // 6B completed
					webPrint("  6B done, still working ... on another Billion iterations ... working ...")
						t := time.Now()
						elapsed := t.Sub(start)
					RunTimeAsString := elapsed.String()
					webPrint(fmt.Sprintf("%s", RunTimeAsString))
				}
				if iterInt64 == 7000000000 { // 7B completed
					webPrint("  7B done, still working ... on another Billion iterations ... working ...")
						t := time.Now()
						elapsed := t.Sub(start)
					RunTimeAsString := elapsed.String()
					webPrint(fmt.Sprintf("%s", RunTimeAsString))
				}
				if iterInt64 == 8000000000 { // 8B completed
					webPrint("  8B done, still working ... on another Billion iterations ... working ...")
						t := time.Now()
						elapsed := t.Sub(start)
					RunTimeAsString := elapsed.String()
					webPrint(fmt.Sprintf("%s", RunTimeAsString))
				}
				if iterInt64 == 9000000000 { // 9B completed
					webPrint("  9B done, still working ... on another five Billion iterations ... working ...")
						t := time.Now()
						elapsed := t.Sub(start)
					RunTimeAsString := elapsed.String()
					webPrint(fmt.Sprintf("%s", RunTimeAsString))
				}
				if iterInt64 == 14000000000 { // 14B completed
					webPrint("  14B done, still working ... on another five Billion iterations ... working ...")
						t := time.Now()
						elapsed := t.Sub(start)
					RunTimeAsString := elapsed.String()
					webPrint(fmt.Sprintf("%s", RunTimeAsString))
				}
				if iterInt64 == 19000000000 { // 19B completed
					webPrint("  19B done, still working ... on another five Billion iterations ... working ...")
						t := time.Now()
						elapsed := t.Sub(start)
					RunTimeAsString := elapsed.String()
					webPrint(fmt.Sprintf("%s", RunTimeAsString))
				}
				if iterInt64 == 24000000000 { // 24B completed
					webPrint("  24B done, still working ... on another five Billion iterations ... working ...")
						t := time.Now()
						elapsed := t.Sub(start)
					RunTimeAsString := elapsed.String()
					webPrint(fmt.Sprintf("%s", RunTimeAsString))
				}
				if iterInt64 == 29000000000 { // 29B completed
					webPrint("  29B done, still working ... on another five Billion iterations ... working ...")
						t := time.Now()
						elapsed := t.Sub(start)
					RunTimeAsString := elapsed.String()
					webPrint(fmt.Sprintf("%s", RunTimeAsString))
				}
				if iterInt64 == 34000000000 { // 34B completed
					webPrint("  34B done, still working ... just another six Billion iterations to go! ... ")
						t := time.Now()
						elapsed := t.Sub(start)
					RunTimeAsString := elapsed.String()
					webPrint(fmt.Sprintf("%s", RunTimeAsString))
				}
				if iterInt64 == 40000000000 { // 40B completed
					webPrint(fmt.Sprintf("%0.12f is our Pi calculated using an infinite series by John Wallis circa 1655", π))
					webPrint("3.14159265358  is the value of π from the web")
	
					t := time.Now()
					elapsed := t.Sub(start)
					RunTimeAsString := elapsed.String()
	
					piAsAString := strconv.FormatFloat(π, 'g', -1, 64)
					checkPiUpTo255chars(piAsAString)
					webPrint(fmt.Sprintf("Forty Billion iterations were completed in %s yielding π to %d confirmed digits", RunTimeAsString, copyOfLastPosition))
					webPrint(" per --  an infinite series by John Wallis circa 1655") // ----------------------
					LinesPerIter = 36                                                                // an estimate
					webPrint(fmt.Sprintf("at aprox %0.1f lines of code per iteration ...", LinesPerIter))
					LinesPerSecond = (LinesPerIter * iterFloat64) / elapsed.Seconds()
					formattedLinesPerSecond := formatInt64WithThousandSeparators(int64(LinesPerSecond)) // .Seconds() returns a float64
					webPrint(fmt.Sprintf("Aprox %s lines of code were executed per second ", formattedLinesPerSecond))
	
					// store reults in a log file which can be displayed from within the program by selecting option #12
					fileHandle, err1 := os.OpenFile("dataLog-From_calculate-pi-and-friends.txt", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600) // append to file
					check(err1)                                                                                                             // ... gets a file handle to dataLog-From_calculate-pi-and-friends.txt
					defer fileHandle.Close()                                                                                                // It’s idiomatic to defer a Close immediately after opening a file.
					Hostname, _ := os.Hostname()
					_, err0 := fmt.Fprintf(fileHandle, "  -- John Wallis (cont.) -- on %s ", Hostname)
					check(err0)
					current_time := time.Now()
					_, err6 := fmt.Fprint(fileHandle, "was run on: ", current_time.Format(time.ANSIC), "")
					check(err6)
					_, err2 := fmt.Fprintf(fileHandle, "%s was Lines/Second  ", formattedLinesPerSecond)
					check(err2)
					_, err4 := fmt.Fprintf(fileHandle, "%.02f was Iterations/Seconds  ", iterFloat64/elapsed.Seconds())
					check(err4)
					_, err5 := fmt.Fprintf(fileHandle, "%e was total Iterations  ", iterFloat64)
					check(err5)
					TotalRun := elapsed.String()                                         // cast time durations to a String type for Fprintf "formatted print"
					_, err7 := fmt.Fprintf(fileHandle, "Total run was %s  ", TotalRun) // add total runtime of this calculation
					check(err7)
				}
		} // end of select
	} // end of for interInt64 < 40B
// written entirely by Richard Woolley
calculating = false
return π

} // end of JohnWallis()

func JohnWallis_grok(updateOutput func(string), done chan bool, webPrint func(string)) float64 {
	result := 2.0
	iterations := 1_000_000 // Reduced for testing; revert to 40_000_000_000 for full run
	for n := 1; n <= iterations; n++ {
		select {
		case <-done:
			updateOutput("Wallis stopped early")
			return result
		default:
			term := float64(2*n) / float64(2*n-1) * float64(2*n) / float64(2*n+1)
			result *= term
			if n%100_000 == 0 {
				// updateOutput(webPrint(fmt.Sprintf("Iteration %d: π ≈ %.10f", n, result))) // ::: todo: CHECK THIS OUT ---------------------------
				webPrint(fmt.Sprintf("Iteration %d: π ≈ %.10f", n, result)) // ::: todo: CHECK THIS OUT ---------------------------
			}
		}
	}
	return result
}