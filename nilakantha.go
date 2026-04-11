package main

import (
	"fmt"
	"math/big"
	"os"
	"strings"
)

// @formatter:off

func NilakanthaBig(iters int, precision int, done chan bool, webPrint func(string2 string)) { // Changed signature ::: - -
var printThisThen string
var printThis []string
var lenOfPi int

	webPrint("... working ...")

	if iters > 36111222 {
		webPrint("... working ... Nilakantha using big floats")
	}
	if iters > 42000000 {
		webPrint("... werkin ...")
	}
	if iters > 55111222 {
		webPrint("... working for a while ...")
	}
	if iters > 69111222 {
		webPrint("... will be working for quite a while ...")
	}
	if iters > 80111222 {
		webPrint("... a very long while ... working ...")
	}

	// start := time.Now()

	var iterBig int

	// big.Float "constants":

		twoBig := big.NewFloat(2)
		threeBig := big.NewFloat(3)
		fourBig := big.NewFloat(4)

	// big.Float variables:

		digitoneBig := new(big.Float)
		*digitoneBig = *twoBig
	
		digittwoBig := new(big.Float)
		*digittwoBig = *threeBig
	
		digitthreeBig := new(big.Float)
		*digitthreeBig = *fourBig
	
		sumBig := new(big.Float)
		nexttermBig := new(big.Float)

	// set precision to a user-specified value
		sumBig.SetPrec(uint(precision))
		twoBig.SetPrec(uint(precision))
		threeBig.SetPrec(uint(precision))
		fourBig.SetPrec(uint(precision))
		digitoneBig.SetPrec(uint(precision))
		digittwoBig.SetPrec(uint(precision))
		digitthreeBig.SetPrec(uint(precision))
		nexttermBig.SetPrec(uint(precision))

	// ::: calculate initial value  	
	sumBig.Add(threeBig, new(big.Float).Quo(fourBig, new(big.Float).Mul(digitoneBig, new(big.Float).Mul(digittwoBig, digitthreeBig))))

	fileHandleBig, err1prslc2c := os.OpenFile("big_pie_is_in_here.txt", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600) // append to file
	check(err1prslc2c)           // ... gets a file handle to dataLog-From_calculate-pi-and-friends.txt // It’s idiomatic to defer a Close immediately after opening a file.
	defer func(fileHandleBig *os.File) {
		err := fileHandleBig.Close()
		if err != nil {}
	}(fileHandleBig) 

	iterBig = 1
	for iterBig < iters {
		select {
		case <-done: // ::: here an attempt is made to read from the channel (a closed channel can be read from successfully; but what is read will be the null/zero value of the type of chan (0, false, "", 0.0, etc.)
			// in the case of this particular channel (which is of type bool) we get the value false from having received from the channel when it is already closed. 
			// ::: if the channel known by the moniker "done" is already closed, that/it is to be interpreted as the abort signal by all listening processes. 
			fmt.Println("Goroutine Nilakantha for-loop (1 of 1) is being terminated by select case finding the done channel to be already closed")
			return // Exit the goroutine
		default:
		
		/*
		  -- Nilakantha Somayaji -- on Mac-mini.local
		was run on: Sun Mar 23 21:08:37 2025
		100000000 was total Iterations; 512 was precision setting for the big.Float types
		Total run was 1m4.656298791s   25 verified digits   3.141592653589793238462643
		*/
		
		iterBig++
		
			// ::: Calculate: 
				digitoneBig.Add(digitoneBig, twoBig)
				digittwoBig.Add(digittwoBig, twoBig)
				digitthreeBig.Add(digitthreeBig, twoBig)
		
				nexttermBig.Quo(fourBig, new(big.Float).Mul(digitoneBig, new(big.Float).Mul(digittwoBig, digitthreeBig)))

				if iterBig%2 == 0 { // % is modulus operator
					sumBig.Sub(sumBig, nexttermBig)
				} else {
					sumBig.Add(sumBig, nexttermBig)
				}

		if iterBig == 20111222 {
			webPrint(" ... doin some ... ") // Send to channel
		}
		if iterBig == 36111222 {
			webPrint(" ... werkin ... ")
		}
		if iterBig == 42000000 {
			webPrint("... still werkin ... Nilakantha Somayaji method using big.Float types   -- with some patience one can generate 31 correct digits of pi this way.")
		}
		if iterBig == 55111222 {
			webPrint("... been working for a while ...")
		}
		if iterBig == 69111222 {
			webPrint("... been working for quite a while ...")
		}
		if iterBig == 80111222 {
			webPrint("... it's been a very long while ... but still working ...")
		}
		if iterBig == 180111222 {
			webPrint("... it's been a very long while, 180,111,222 done, ... and still working ...")
		}
		if iterBig == 280111222 {
			webPrint("... it's been a very long while, 280,111,222 done, ... and still working ...")
		}
		if iterBig == 480111222 {
			webPrint("... it's been a very long while, 480,111,222 done, ... still working ...")
		}
		if iterBig == 680111222 {
			webPrint("... it's been a very long while, 680,111,222 done, ...  working ...")
		}
		if iterBig == 880111222 {
			webPrint("... it's been a very long while, done, 880,111,222, done ... still, working ...")
		}
		if iterBig == 977111222 {
			webPrint("... it's been a very long while, 977,111,222 already ... why am I still working? ...")
		}
		}
	} // end of the loop, the only calculating loop
	// t := time.Now()
	
		// ::: bug hammer = do this just once; KISS
		printThis, lenOfPi = checkPiTo100(sumBig) // all local variables defined at the top of this function 
		printThisThen = strings.Join(printThis, "")

	if lenOfPi > 55000 { // if length of pi is > 55,000 digits we have something really big
		// print to ::: screen
			webPrint("We have been tasked with making a lot of pie and it was sooo big it needed its own file ...")
			webPrint("  After allowing this process to finish (you may have to continue prodding this thing along for a while) ... ")
			webPrint("... Go have a look in /.big_pie_is_in_here.txt to find all the digits of π you had requested. ")

		// print (log) to a special ::: file
			_, err2prslc2c := fmt.Fprintf(fileHandleBig, "These are the %d verified digits we have calculated, dumped by rick  :: ", lenOfPi)
				check(err2prslc2c)
	
			_, err8prslc2c := fmt.Fprint(fileHandleBig, printThisThen) // to a file
				check(err8prslc2c)

		err := fileHandleBig.Close()
		if err != nil {
			return
		}
	} else { 

		webPrint(fmt.Sprintf("pi as calculated herein is: %s", printThisThen))

		floatIterBig := float64(iterBig)
		printableIterbigWithcommas := formatFloat64WithThousandSeparators(floatIterBig)

		webPrint(fmt.Sprintf(".... we have matched %d digits in %s iterations: ", lenOfPi, printableIterbigWithcommas))

		fileHandleNilakan, err1prslc2c := os.OpenFile("dataLog-From_Nilakantha_Method_lengthy_prints.txt", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
		check(err1prslc2c)         // ... gets a file handle to dataLog-From_calculate-pi-and-friends.txt  // It’s idiomatic to defer a Close immediately after opening a file.
		defer func(fileHandleNilakan *os.File) {
			err := fileHandleNilakan.Close()
				if err != nil {}
			}(fileHandleNilakan) 
		
		// print to ::: file
				_, err2prslc2c := fmt.Fprintf(fileHandleNilakan,
					"Below rick are the %d verified digits we have calculated via Nilakantha using precision of %d and iterations of %d: ", lenOfPi, precision, iterBig)
						check(err2prslc2c)

		webPrint(fmt.Sprintf("hey, rick, pi as calculated herein is: %s", printThisThen))
		
	}
	webPrint("")

	// Open a log file
		fileHandleDefault, err1 := os.OpenFile("dataLog-From_calculate-pi-and-friends.txt", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600) // append to file
			check(err1)                              // ... gets a file handle to dataLog-From_calculate-pi-and-friends.txt
				defer fileHandleDefault.Close()     // It’s idiomatic to defer a Close immediately after opening a file.

	// print to ::: screen
	webPrint(" via Nilakantha with big floats. Written entirely by Richard Woolley")

	// ::: Prepare to exit the Gregory Nilakantha method functions
	calculating = false // Allow another method to be selected.
	// written entirely by Richard Woolley
} 
