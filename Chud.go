package main

import (
	"fmt"
	"math"
	"math/big"
	"strings"
	"time"
)

// @formatter:off
// ::: todo: finish printing the elapsed time 
// Chudnovsky method, based on https://arxiv.org/pdf/1809.00533.pdf
/*
    The Chudnovsky algorithm is an incredibly-fast algorithm for calculating the digits of pi. It was developed by Gregory Chudnovsky and his
	brother David Chudnovsky in the 1980s. It is more efficient than other algorithms and is based on the theory of modular equations. It has
	been used to calculate pi to over 62 trillion digits.
*/
//  Using this procedure, calculating 1,000,000 digits required 70,516 loops -- per the run on: Sun May 7 2023
//  Total run-time was 8h4m39.7847064s on an old i7 
//  AND, THAT CALCULATION WAS INDEPENDENTLY VERIFIED !!!!!!!!!!!

// This will be a little-bit tricky. We want to use callbacks etc. so that we can use the smoother-scrolling webPrint(fmt.Sprintf("")) way of doing prints ... 
// ... but this chudnovsky section is a cascade of functions: chudnovskyBig()-->calcPi()-->finishChudIfsAndPrint() 
func chudnovskyBig(webPrint func(string), digits int, done chan bool) { // ::: - -

	// ::: webPrint will use updateOutput[1-4] depending on from which window called -- so we pass webPrint to calcPi(webPrint, float64(digits), start, loops) thusly 
	webPrint(fmt.Sprintf("... working ..."))
	var loops int
	piAsBigFloat := new(big.Float).SetPrec(512).SetFloat64(0.0)
	start := time.Now() // ::: start will be passed, and then passed back, in order to be compared with end time t

				// ::: calcPi  <---- runs from here: v v v v v v v  
				loops, start, piAsBigFloat = calcPi(webPrint, float64(digits), start, done, piAsBigFloat) // ::: This is the call to calcPi
				// ::: calcPi --- - - - --- ^ ^ ^ 


	if loops < 100 {
		// .Text('f', 122) converts the big.Float to a string with 122 decimal places
		piString := piAsBigFloat.Text('f', 122) // Unresolved reference 'Text'

		webPrint(fmt.Sprintf("Less than 100 loops, so here is a peek at the prospective value of Pi: %s", piString))
		// prints: Less than 100 loops, so here is a peek at the prospective value of Pi: 3.14159265358979323846264338327950288419716939937510582097494459230781640628620899862803482534211706920799063933996827509056
	}
	

	// The following runs ::: after calcPi 
	webPrint(fmt.Sprintf(" loops were: %d, and digits requested was: %d ", loops, digits))

	webPrint(fmt.Sprintf(" 	The Chudnovsky algorithm is an incredibly-fast algorithm for calculating the digits of pi. It was developed by Gregory Chudnovsky and his "))
	webPrint("brother David Chudnovsky in the 1980s. It is more efficient than other algorithms and is based on the theory of modular equations. It has been ")
	webPrint(fmt.Sprintf("used to calculate pi to over 62 trillion digits."))

	// webPrint(fmt.Sprintf("Final Pi: %s", finalResult))
	done <- true
	}
/*
.
.
.
*/
// calculate Pi for n number of digits
func calcPi(webPrint func(string), digits float64, start time.Time, done chan bool, float *big.Float) (int, time.Time, *big.Float) {
	webPrint("This is an implementation for https://en.wikipedia.org/wiki/Chudnovsky_algorithm")
	webPrint("It can be improved using binary splitting http://numbers.computation.free.fr/Constants/Algorithms/splitting.html")
	webPrint("if we were to split it into two independent parts and simplify the formula. For more details, visit:")
	webPrint("https://www.craig-wood.com/nick/articles/pi-chudnovsky")

	var i int

	// ::: n ...
	// ... apparently, n, is the expected number of loops we may need to produce digits number of digits
	n := int64(2 + int(float64(digits)/14.181647462))
	// comments re: n := int64(2 + int(float64(digits)/12))  // I tried this, and may try something like it again someday?? like /14.0 ?

	// set precision 
	// comments re: precision := uint(int(math.Ceil(math.Log2(10)*digits)) + int(math.Ceil(math.Log10(digits))) + 2) // the original
	// comments re: precision := uint(digits) // not good, not large enough, so ...
	digitsPlus := digits + digits*0.10 // because we needed a little more than the orriginal programmer had figured on :)
	precision := uint(int(math.Ceil(math.Log2(10)*digitsPlus)) + int(math.Ceil(math.Log10(digitsPlus))) + 2)

	c := new(big.Float).Mul(
		big.NewFloat(float64(426880)),
		new(big.Float).SetPrec(precision).Sqrt(big.NewFloat(float64(10005))),
	)

	k := big.NewInt(int64(6))
	k12 := big.NewInt(int64(12))
	l := big.NewFloat(float64(13591409))
	lc := big.NewFloat(float64(545140134))
	x := big.NewFloat(float64(1))
	xc := big.NewFloat(float64(-262537412640768000))
	m := big.NewFloat(float64(1))
	sum := big.NewFloat(float64(13591409))

	pi := big.NewFloat(0)

	x.SetPrec(precision)
	m.SetPrec(precision)
	sum.SetPrec(precision)
	pi.SetPrec(precision)

	bigI := big.NewInt(0)
	bigOne := big.NewInt(1)

	// this is a flag; if it is set to zero we exit
	queryIfTimeToDie := 1
	i = 1 // a secondary dedicated loop counter


	if n > 8998 {
		webPrint(" Well, this is going to take a while, because you asked for too much pie (> 8990)")
	}


	for ; n > 0; n-- {
		select {
		case <-done: // ::: here an attempt is made to read from the channel (a closed channel can be read from successfully; but what is read will be the null/zero value of the type of chan (0, false, "", 0.0, etc.)
			// in the case of this particular channel (which is of type bool) we get the value false from having received from the channel when it is already closed. 
			// ::: if the channel known by the moniker "done" is already closed, that/it is to be interpreted as the abort signal by all listening processes. 
			webPrint("Goroutine chud-func-calcPi for-loop (1 of 1) is being terminated by select case finding the done channel to be already closed")
			return i, start, pi // Exit the goroutine
		default:
			i++

			// L calculation
			l.Add(l, lc)

			// X calculation
			x.Mul(x, xc)

			// M calculation
			kpower3 := big.NewInt(0)
			kpower3.Exp(k, big.NewInt(3), nil)
			ktimes16 := new(big.Int).Mul(k, big.NewInt(16))
			mtop := big.NewFloat(0).SetPrec(precision)
			mtop.SetInt(new(big.Int).Sub(kpower3, ktimes16))
			mbot := big.NewFloat(0).SetPrec(precision)
			mbot.SetInt(new(big.Int).Exp(new(big.Int).Add(bigI, bigOne), big.NewInt(3), nil))
			mtmp := big.NewFloat(0).SetPrec(precision)
			mtmp.Quo(mtop, mbot)
			m.Mul(m, mtmp)

			// Sum calculation
			t := big.NewFloat(0).SetPrec(precision)
			t.Mul(m, l)
			t.Quo(t, x)
			sum.Add(sum, t)

			// Pi calculation
			pi.Quo(c, sum)
			k.Add(k, k12)
			bigI.Add(bigI, bigOne)

			if i == 2 {
				finishChudIfsAndPrint(webPrint, pi, "no")
			}
			if i == 4 {
				finishChudIfsAndPrint(webPrint, pi, "no")
			}
			if i == 8 {
				finishChudIfsAndPrint(webPrint, pi, "no")
			}
			if i == 16 {
				finishChudIfsAndPrint(webPrint, pi, "no")
			}
			if i == 32 {
				finishChudIfsAndPrint(webPrint, pi, "no")
			}
			if i == 44 {
				finishChudIfsAndPrint(webPrint, pi, "no")
			}
			if i == 52 {
				finishChudIfsAndPrint(webPrint, pi, "no")
			}
			if i == 62 {
				finishChudIfsAndPrint(webPrint, pi, "no")
			}
			if i == 72 {
				finishChudIfsAndPrint(webPrint, pi, "no")
			}
			if i == 82 {
				finishChudIfsAndPrint(webPrint, pi, "no")
			}
			if i == 92 {
				finishChudIfsAndPrint(webPrint, pi, "no")
			}

			useAlternateFile := "no" // no means to use the standard log file rather than some special one
			// the compiler is not happy unless it sees this created outside of an if
			// But, wait. Why is the compiler allowing me to violate the no new var left of the := assignment ??? This IS in a loop !!!!
			if i == 100 {
				// useAlternateFile := "no" // the compiler is not happy unless it sees this created outside of an if
				webPrint(fmt.Sprintf(" we are at %d loops", i))
			}
			if i == 200 {
				// useAlternateFile = "no" // still no
				webPrint(fmt.Sprintf(" we are at %d loops: ", i))
				finishChudIfsAndPrint(webPrint, pi, useAlternateFile)
			}
			if i == 400 {
				// useAlternateFile = "no" // still no ::: based on this flag ...
				webPrint(fmt.Sprintf(" we are at %d loops: ", i))
				finishChudIfsAndPrint(webPrint, pi, useAlternateFile)
			}
			// ::: ... up to this point the user will be shown the verified pi message
			//
			// note below the: useAlternateFile = "chudDid800orMoreLoops"
			if i == 800 {
				useAlternateFile = "chudDid800orMoreLoops"
				webPrint(fmt.Sprintf(" we are at %d loops: ", i))
				finishChudIfsAndPrint(webPrint, pi, useAlternateFile)
			}
			if i == 1600 {
				useAlternateFile = "chudDid800orMoreLoops"
				webPrint(fmt.Sprintf(" we are at %d loops: ", i))
				finishChudIfsAndPrint(webPrint, pi, useAlternateFile)
			}
			if i == 2000 {
				useAlternateFile = "chudDid800orMoreLoops"
				webPrint(fmt.Sprintf(" we are at %d loops: ", i))
				finishChudIfsAndPrint(webPrint, pi, useAlternateFile)
			}
			if i == 2400 {
				useAlternateFile = "chudDid800orMoreLoops"
				webPrint(fmt.Sprintf(" we are at %d loops: ", i))
				finishChudIfsAndPrint(webPrint, pi, useAlternateFile)
			}
			if i == 2800 {
				useAlternateFile = "chudDid800orMoreLoops"
				webPrint(fmt.Sprintf(" we are at %d loops: ", i))
				finishChudIfsAndPrint(webPrint, pi, useAlternateFile)
			}
			if i == 3200 {
				useAlternateFile = "chudDid800orMoreLoops"
				webPrint(fmt.Sprintf(" we are at %d loops: ", i))
				finishChudIfsAndPrint(webPrint, pi, useAlternateFile)
			}
			if i == 4000 {
				useAlternateFile = "chudDid800orMoreLoops"
				webPrint(fmt.Sprintf(" we are at %d loops: ", i))
				finishChudIfsAndPrint(webPrint, pi, useAlternateFile)
			}
			if i == 6000 {
				useAlternateFile = "chudDid800orMoreLoops"
				webPrint(fmt.Sprintf(" we are at %d loops: ", i))
				finishChudIfsAndPrint(webPrint, pi, useAlternateFile)
			}
			if i == 8000 {
				useAlternateFile = "chudDid800orMoreLoops"
				webPrint(fmt.Sprintf(" we are at %d loops: ", i))
				finishChudIfsAndPrint(webPrint, pi, useAlternateFile)
			}
			if queryIfTimeToDie == 0 {
				webPrint(fmt.Sprintf("if queryIfTimeToDie is 0, time to die"))
				webPrint(fmt.Sprintf("precisionision was: %d ", precision))
				break
			}
			// 1,000,000 digits requires 70516 loops, per the run on May 7 2023 at 10:30
			//  was run on: Sun May  7 08:50:23 2023
			//  Total run was 8h4m39.7847064s
			// AND THE CALCULATION WAS INDEPENDANTLY VERIFIED !!!!!!!!!!!
		} // end of select
	} // end of for loop way up thar :: it prompts periodically to continue or die

	/*
	// ::: we are out of the loop, so we do the following just once:

	// obtain file handle
	fileHandleBig, err1prslc2c := os.OpenFile("big_pie_is_in_here.txt", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600) // append to file
	check(err1prslc2c)                                                                                             // ... gets a file handle to dataLog-From_calculate-pi-and-friends.txt
	defer fileHandleBig.Close()                                                                                    // It’s idiomatic to defer a Close immediately after opening a file.

	// to ::: file		
	_, err9bigpie := fmt.Fprint(fileHandleBig, pi)                               // dump this big-assed pie to a special log file
	check(err9bigpie)
	_, err9bigpie = fmt.Fprint(fileHandleBig, "was pi as a big.Float")  // add a suffix 
	check(err9bigpie)

	_, errGoesHere := fmt.Fprint(fileHandleBig, "")
	check(errGoesHere)

	fileHandleBig.Close()
	
	 */

	return i, start, pi // assigning i to loops in caller
}
/*
.
.
.
.
.
*/
// a helper func   
func finishChudIfsAndPrint(webPrint func(string), pi *big.Float, useAlternateFile string) { // ::: - -

	// ::: Check pi and convert to []string -- and, set lenOfPi
	stringVerOfOurCorrectDigits, lenOfPi := checkPiTo59766(pi)

	
	/*
	if lenOfPi < 600 {
		
		//	print to ::: screen
		webPrint(fmt.Sprintf("lenOfPi < 600, so, Here are %d calculated digits that we have verified by reference (one at a time): ", lenOfPi))

		// print via range to ::: screen	
		for _, oneChar := range stringVerOfOurCorrectDigits { // pi is finally ::: printed here via ranging 
			// to screen:
			webPrint(fmt.Sprintf("%s", oneChar)) // ::: to screen
		}

		// print to ::: screen	
		// webPrint(fmt.Sprintf("")) // ::: ---- this may be needed ------------------
	}
	
	 */



	if lenOfPi < 600 {
		webPrint(fmt.Sprintf("Here are %d verified digits (one at a time): ", lenOfPi))

		for _, oneChar := range stringVerOfOurCorrectDigits {
			webPrint(string(oneChar)) // Send one character

			// Optional: Add a tiny delay to simulate the old desktop 'crawl'
			// 10-20 milliseconds is usually perfect
			time.Sleep(20 * time.Millisecond)
		}
		webPrint("\n")
	}
	
	
	

	if lenOfPi > 46000 { // if length of pi is > 48,000 digits we have something really big
		// print to ::: screen
		webPrint("We have been tasked with making a lot of pie and it was sooo big it needed its own file ...")
		webPrint("... Go have a look in /.big_pie_is_in_here.txt to find all the digits of π you had requested. ")
	} else {

		// } else { continues below: (in other words, the following if-else conditions are only checked if length of pi was < 55,000 digits)
		if useAlternateFile == "ChudDidLessThanOneHundredLoops" {
			// print to ::: screen
			webPrint(fmt.Sprintf(" Here are %d calculated digits that we have verified by reference: ", lenOfPi))

			for _, oneChar := range stringVerOfOurCorrectDigits {
				// ::: to screen
				webPrint(oneChar)
			}
			// ::: this final else handles any instances of useAlternateFile not caught above
		} else {
			// to ::: screen
			webPrint(fmt.Sprintf(" Here are %d calculated digits that we have verified by reference:", lenOfPi))

			asString := strings.Join(stringVerOfOurCorrectDigits, "")
			webPrint(fmt.Sprintf(" catch-all, asString: %s", asString))
		}
	} // end of if's else, way up thar "if lenOfPi > 46000 {} else {"   so, this has been the instance where pi is shorter than 55,000
}