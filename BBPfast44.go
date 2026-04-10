package main

import (
	"fmt"
	"math"
	"math/big"
	"os"
	"time"
)

func bbpFast44(digits int) { // case 42: // -- AMFbbp_formulaA
	fmt.Sprintf("bbpFast46 executed with %d digits", digits)

	usingBigFloats = true
	iters_bbp := 1
	start := time.Now()
	// numCPU := runtime.NumCPU()
	// runtime.GOMAXPROCS(numCPU)

	n := digits
	p := uint((int(math.Log2(10)))*n + 3)

	result := make(chan *big.Float, n)

	pi := new(big.Float).SetPrec(p).SetInt64(0)

	for i := 0; i < n; i++ {
		select {

		default:
			iters_bbp = i
		}
	}

	for i := 0; i < n; i++ {
		select {

		default:
			pi.Add(pi, <-result)
			iters_bbp = i
		}
	}

	fmt.Sprintf("%[1]*.[2]*[3]f \n", 1, n, pi) // n is the number of digits of pi to calculate

	// log run stats to a log file
	t := time.Now()
	elapsed := t.Sub(start)
	fileHandle, err1 := os.OpenFile("dataLog-From_calculate-pi-and-friends.txt", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600) // append to file
	check(err1)                                                                                                             // ... gets a file handle to dataLog-From_calculate-pi-and-friends.txt
	defer fileHandle.Close()                                                                                                // It’s idiomatic to defer a Close immediately after opening a file.
	Hostname, _ := os.Hostname()
	_, err0 := fmt.Fprintf(fileHandle, "\n  -- calculate pi using the bbp formula -- on %s \n", Hostname)
	check(err0)
	current_time := time.Now()
	_, err6 := fmt.Fprint(fileHandle, "was run on: ", current_time.Format(time.ANSIC), "\n")
	check(err6)
	_, err4 := fmt.Fprintf(fileHandle, "%.02f was Iterations/Seconds \n", float64(iters_bbp)/elapsed.Seconds())
	check(err4)
	_, err5 := fmt.Fprintf(fileHandle, "%d was total Iterations \n", iters_bbp)
	check(err5)
	TotalRun := elapsed.String() // cast time durations to a String type for Fprintf "formatted print"
	_, err7 := fmt.Fprintf(fileHandle, "Total run was %s \n ", TotalRun)
	check(err7)

	// ::: Prepare to exit the BBP fast 44 method functions
	calculating = false // Allow another method to be selected.
	/*
		for _, btn := range buttons2 { // ok to only Enable buttons1, because I expect to only ever execute this from window2
			btn.Enable() // ::: Enable
		}

	*/
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
				fmt.Println("Goroutine BBPfast44-func-workers for-loop (3 of 3) is being terminated by select case finding the done channel to be already closed")
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
