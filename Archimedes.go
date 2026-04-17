package main

import (
	"fmt"
	"math/big"
	"time"
)

// @formatter:off

func ArchimedesBig(done chan bool, webPrint func(string)) {

	// 1. Send a padding comment to prime the SSE connection
	padding := ":"
	for k := 0; k < 1024; k++ {
		padding += " "
	}
	webPrint(padding)

	// 2. Brief network breath
	time.Sleep(50 * time.Millisecond)

	webPrint("You've selected a demonstration of Dick's improved version of Archimedes' method for aproximating the value of Pi : 3.14159...")
	webPrint("The goal is to accurately calculate over 2,700 correct digits of Pi. We'll need to use floating-point numbers with thousands of decimal places.")
	webPrint("This can be done using Dick's most-favoured language: go.lang, or simply Go. GoLand (by JetBrains) wil be our IDE.")
	webPrint("All of our variables must be big.Floats (as in the above code, this we now do, again.)")

	r := big.NewFloat(1)
	s1 := big.NewFloat(1)
	numberOfSides := big.NewFloat(6)

	a := new(big.Float)
	b := new(big.Float)
	p := new(big.Float)
	s2 := new(big.Float)
	p_d := new(big.Float)
	s1_2 := new(big.Float)

	webPrint("Go's precision is set to 55000 on all of our variables (as per the above code).")

	precision := 55000
	p_d.SetPrec(uint(precision))
	a.SetPrec(uint(precision))
	s1_2.SetPrec(uint(precision))
	s2.SetPrec(uint(precision))
	b.SetPrec(uint(precision))
	p.SetPrec(uint(precision))
	r.SetPrec(uint(precision))
	s1.SetPrec(uint(precision))
	numberOfSides.SetPrec(uint(precision))
	webPrint("Then, we do some initial assignments and calculations (as per above):")

	// Initial calculation
	numberOfSides.Mul(numberOfSides, big.NewFloat(2))
	s1_2.Quo(s1, big.NewFloat(2))
	a.Sqrt(new(big.Float).Sub(r, new(big.Float).Mul(s1_2, s1_2)))

	webPrint("  First we will determine the height (a) of a right triangle formed by bisecting a side of a polygon inscribed")
	webPrint("  in a unit circle (radius r = 1).")
	webPrint("  The polygon's side length (s1) is halved (s1_2 = s1 / 2), and this computation helps refine the polygon's ")
	webPrint("  perimeter to approximate π as the number of the sides of the polygon increases.")

	b.Sub(r, a)
	s2.Sqrt(new(big.Float).Add(new(big.Float).Mul(b, b), new(big.Float).Mul(s1_2, s1_2)))

	webPrint("     Here is some pseudo code for the algorithm:")
	webPrint("     Inputs:  b: short side from midpoint to circle edge (a big float)")
	webPrint("              s1_2: half the current side length (s1 / 2, a big float)")
	webPrint("     Output:")
	webPrint("              s2: new side length of the polygon (a big float)")
	webPrint("     Step 1.")
	webPrint("        Compute b^2")
	webPrint("        temp1 = b * b")
	webPrint("     Step 2. Compute (s1_2)^2")
	webPrint("        temp2 = s1_2 * s1_2")
	webPrint("     Step 3. Add the two squares")
	webPrint("        temp3 = temp1 + temp2")
	webPrint("     Step 4. Take the square root to get the new side length")
	webPrint("        s2 = square_root(temp3)")
	webPrint("     Now, we get to work!!")

	s1.Set(s2)
	p.Mul(numberOfSides, s1)
	p_d.Set(p)

	// sleepOrAbort sleeps for the given duration but returns false immediately
	// if the done channel is closed, allowing the loop to exit cleanly.
	sleepOrAbort := func(d time.Duration) bool {
		if d == 0 {
			select {
			case <-done:
				return false
			default:
				return true
			}
		}
		select {
		case <-done:
			return false
		case <-time.After(d):
			return true
		}
	}

	for i := 0; i < 5001; i++ {

		// Check for stop signal at the top of every iteration.
		select {
		case <-done:
			webPrint("  Archimedes calculation was stopped by user.")
			return
		default:
		}

		// 1. DO THE MATH
		numberOfSides.Mul(numberOfSides, big.NewFloat(2))
		s1_2.Quo(s1, big.NewFloat(2))
		a.Sqrt(new(big.Float).Sub(r, new(big.Float).Mul(s1_2, s1_2)))
		b.Sub(r, a)
		s2.Sqrt(new(big.Float).Add(new(big.Float).Mul(b, b), new(big.Float).Mul(s1_2, s1_2)))
		s1.Set(s2)
		p.Mul(numberOfSides, s1)
		p_d.Set(p)
		p_d.Quo(p_d, big.NewFloat(2))

		// 2. SEND THE RESULTS TO THE BROWSER
		if i == 24 {
			webPrint(fmt.Sprintf("------------------------------------------------------------------" +
				"  %d iterations were completed in order to yeild the following digits of π", i))
			webPrint(fmt.Sprintf("    %.20f is the big.Float of what we have calculated  ----- per Archimedes' at 24 iterations, formatted: 20f", p_d))
			webPrint("    3.141592653589793238  vs the value of π from the web")
			formattedNum := formatWithThousandSeparators(numberOfSides)
			webPrint(fmt.Sprintf("the above was estimated from a %s  --- sided polygon", formattedNum))
			_, lenOfPi := checkPiTo59766(p_d)
			webPrint(fmt.Sprintf("... And, it has been verified that we actually calculated pi correctly to %d digits!", lenOfPi))
			webPrint("... Mister A. would have wept!")
		}

		if i == 50 {
			webPrint(fmt.Sprintf("------------------------------------------------------------------" +
				"  %d iterations were completed in order to yeild the following digits of π", i))
			webPrint(fmt.Sprintf("    %.33f is the big.Float of what we have calculated  ----- per Archimedes' at 50 iters, formatted: 33f", p_d))
			webPrint("    3.141592653589793238462643383279502  ----- is the value of π from the web")
			_, lenOfPi := checkPiTo59766(p_d)
			webPrint(fmt.Sprintf("... And, it has been verified that we actually calculated pi correctly to %d digits!", lenOfPi))
			formattedNum := formatWithThousandSeparators(numberOfSides)
			webPrint(fmt.Sprintf(" the above was estimated from a %s  --- sided polygon", formattedNum))
		}

		if i == 150 {
			webPrint(fmt.Sprintf("------------------------------------------------------------------" +
				"  %d iterations were completed in order to yeild the following digits of π", i))
			webPrint(fmt.Sprintf("   %.95f   ----- per Rick's modified Archimedes' method, formatted 95f", p_d))
			webPrint("   3.14159265358979323846264338327950288419716939937510582097494459230781640628620899862803482534211  ----- is from web")
			_, lenOfPi := checkPiTo59766(p_d)
			webPrint(fmt.Sprintf("... And, it has been verified that we actually calculated pi correctly to %d digits!", lenOfPi))
			formattedNum := formatWithThousandSeparators(numberOfSides)
			webPrint(fmt.Sprintf(" the above was estimated from a %s  --- sided polygon", formattedNum))
		}

		if i == 200 {
			webPrint(fmt.Sprintf("------------------------------------------------------------------" +
				"  %d iterations were completed in order to yeild the following digits of π", i))
			webPrint(fmt.Sprintf("   %.122f   ---- ... Archimedes' method, formatted: 122f", p_d))
			webPrint("   3.14159265358979323846264338327950288419716939937510582097494459230781640628620899862803482534211706798214808651328230664709  ----- is from web")
			formattedNum := formatWithThousandSeparators(numberOfSides)
			webPrint(fmt.Sprintf("our figure was estimated from a %s  --- sided polygon", formattedNum))
			_, lenOfPi := checkPiTo59766(p_d)
			webPrint(fmt.Sprintf("... And, it has been verified that we actually calculated pi correctly to %d digits!", lenOfPi))
			webPrint("... working ...")
		}

		if i == 1200 || i == 2200 || i == 3200 || i == 4200 {
			webPrint(fmt.Sprintf("... still working, %d iterations completed ...", i))
		}

		if i == 4500 {
			webPrint("------------------------------------------------------------------")
			formattedNum := formatWithThousandSeparators(numberOfSides)
			webPrint(fmt.Sprintf("All Done! So, how many sides does our polygon have now? A lot: A staggering:%sSIDED POLYGON !!!", formattedNum))
			webPrint(fmt.Sprintf("%d iterations were completed to yeild well over 2,700 correct digits of π!!!", i))
			webPrint(fmt.Sprintf("Go's math/big objects were set to a precision value of: %d  --- here is your GIANT slice of pie: %.2800f", precision, p_d))
			_, lenOfPi := checkPiTo59766(p_d)
			webPrint(fmt.Sprintf("... And, it has been verified that we actually calculated pi correctly to %d digits! by Richard (Rick) H. Woolley", lenOfPi))
		}

		// Paced sleeps -- each one is interruptible by the done channel.
		var sleepDur time.Duration
		switch {
		case i < 24:
			if i == 2 { webPrint("\t\tSleeping each iteration for 135 milliseconds...") }
			sleepDur = 135 * time.Millisecond
		case i < 50:
			if i == 26 { webPrint("\t\tSleeping each iteration for 55 milliseconds...") }
			sleepDur = 55 * time.Millisecond
		case i < 150:
			if i == 52 { webPrint("\t\tSleeping each iteration for 35 milliseconds...") }
			sleepDur = 35 * time.Millisecond
		case i < 400:
			if i == 152 { webPrint("\t\tSleeping each iteration for 7 milliseconds...") }
			sleepDur = 7 * time.Millisecond
		case i < 1100:
			if i == 402 { webPrint("\t\tSleeping each iteration for 2 milliseconds...") }
			sleepDur = 2 * time.Millisecond
		case i < 2000:
			if i == 1102 { webPrint("\t\tSleeping each iteration for 1 millisecond...") }
			sleepDur = time.Millisecond
		default:
			if i == 2002 { webPrint("\t\tNo more sleeping!!!...") }
			sleepDur = 0
		}

		if !sleepOrAbort(sleepDur) {
			webPrint("  Archimedes calculation was stopped by user.")
			return
		}
	}
}
