package main

import (
	"fmt"
	"math"
	"math/big"
	"time"
)

func bbpFast44(webPrint func(string), digits int, done chan bool) { // case 42: // -- AMFbbp_formulaA
	webPrint(fmt.Sprintf("BBP executed with %d digits", digits))

	usingBigFloats = true
	start := time.Now()
	// numCPU := runtime.NumCPU()
	// runtime.GOMAXPROCS(numCPU)

	n := digits
	p := uint((int(math.Log2(10)))*n + 3)

	result := make(chan *big.Float, n)
	worker := workers(p, done)

	pi := new(big.Float).SetPrec(p).SetInt64(0)

	for i := 0; i < n; i++ {
		select {
		case <-done: // ::: here an attempt is made to read from the channel (a closed channel can be read from successfully; but what is read will be the null/zero value of the type of chan (0, false, "", 0.0, etc.)
			// in the case of this particular channel (which is of type bool) we get the value false from having received from the channel when it is already closed. 
			// ::: if the channel known by the moniker "done" is already closed, that/it is to be interpreted as the abort signal by all listening processes. 
			webPrint("Goroutine BBPfast44 for-loop (1 of 3) is being terminated by select case finding the done channel to be already closed")
			return // Exit the goroutine
		default:
			go worker(i, result)
			iters_bbp = i
		}
	}

	for i := 0; i < n; i++ {
		select {
		case <-done: // ::: here an attempt is made to read from the channel (a closed channel can be read from successfully; but what is read will be the null/zero value of the type of chan (0, false, "", 0.0, etc.)
			// in the case of this particular channel (which is of type bool) we get the value false from having received from the channel when it is already closed. 
			// ::: if the channel known by the moniker "done" is already closed, that/it is to be interpreted as the abort signal by all listening processes. 
			webPrint("Goroutine BBPfast44 for-loop (2 of 3) is being terminated by select case finding the done channel to be already closed")
			return // Exit the goroutine
		default:
			pi.Add(pi, <-result)
			iters_bbp = i
		}
	}

	dur := time.Since(start)
	// fyneFunc(fmt.Sprintf("took %v to calculate %d digits of pi \n", dur, n)) // original, prior to grok

	// output := fmt.Sprintf("%s\nIt only took BBP %v to calculate the following %d digits of pi\n", codeSnippet, dur, n)
	output := fmt.Sprintf("\nIt only took BBP %v to calculate the following %d digits of pi\n", dur, n)

	// Display in the GUI
	webPrint(output)

	// fmt.Printf("%[1]*.[2]*[3]f \n", 1, n, pi) // original from CLI version

	// updateChan <- updateData{text:"%[1]*.[2]*[3]f \n", 1, n, pi} // does not work, even with the correct signature for updateChan <- updateData{text:"
	webPrint(fmt.Sprintf("%[1]*.[2]*[3]f \n", 1, n, pi)) // n is the number of digits of pi to calculate

}

func workers(p uint, done chan bool) func(id int, result chan *big.Float) {
	B1 := new(big.Float).SetPrec(p).SetInt64(1)
	B2 := new(big.Float).SetPrec(p).SetInt64(2)
	B4 := new(big.Float).SetPrec(p).SetInt64(4)
	B5 := new(big.Float).SetPrec(p).SetInt64(5)
	B6 := new(big.Float).SetPrec(p).SetInt64(6)
	B8 := new(big.Float).SetPrec(p).SetInt64(8)
	B16 := new(big.Float).SetPrec(p).SetInt64(16)

	return func(id int, result chan *big.Float) {
		Bn := new(big.Float).SetPrec(p).SetInt64(int64(id))

		C1 := new(big.Float).SetPrec(p).SetInt64(1)
		for i := 0; i < id; i++ {
			select {
			case <-done: // ::: here an attempt is made to read from the channel (a closed channel can be read from successfully; but what is read will be the null/zero value of the type of chan (0, false, "", 0.0, etc.)
				// in the case of this particular channel (which is of type bool) we get the value false from having received from the channel when it is already closed. 
				// ::: if the channel known by the moniker "done" is already closed, that/it is to be interpreted as the abort signal by all listening processes. 
				// webPrint("Goroutine BBPfast44-func-workers for-loop (3 of 3) is being terminated by select case finding the done channel to be already closed")
				return // Exit the goroutine
			default:
				C1.Mul(C1, B16)
			}
		}

		C2 := new(big.Float).SetPrec(p)
		C2.Mul(B8, Bn)

		T1 := new(big.Float).SetPrec(p)
		T1.Add(C2, B1)
		T1.Quo(B4, T1)

		T2 := new(big.Float).SetPrec(p)
		T2.Add(C2, B4)
		T2.Quo(B2, T2)

		T3 := new(big.Float).SetPrec(p)
		T3.Add(C2, B5)
		T3.Quo(B1, T3)

		T4 := new(big.Float).SetPrec(p)
		T4.Add(C2, B6)
		T4.Quo(B1, T4)

		R := new(big.Float).SetPrec(p)
		R.Sub(T1, T2)
		R.Sub(R, T3)
		R.Sub(R, T4)
		R.Quo(R, C1)

		result <- R
	}
	// adapted by Richard Woolley
} // end of bbp_formula()
