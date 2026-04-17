package main

// ScoreBoardNila_web.go
//
// Nilakantha Somayaji's alternating series for π (c. 1530):
//   π = 3 + 4/(2·3·4) − 4/(4·5·6) + 4/(6·7·8) − ···
//
// Runs in two phases to demonstrate the hard ceiling of float64
// and what happens when you push past it with big.Float:
//
// Phase 1 — float64, concurrent goroutines
//   Up to 10,000 terms, each computed in its own goroutine.
//   A live scoreboard updates in place as terms arrive.
//   Milestone messages celebrate each new correct decimal digit.
//   Converges to ~15 digits, then float64 has nothing left to give.
//
// Phase 2 — big.Float, 512-bit precision, sequential
//   Picks up the series where Phase 1 left off (same k counter).
//   Runs up to 50,000 additional terms.
//   The moment new digits appear beyond the float64 wall, they are
//   announced with fanfare. Uses piForGauss as reference (3000+ digits).
//
// Adapted by Richard Woolley with a lot of help from Claude.

import (
	"fmt"
	"math"
	"math/big"
	"strings"
	"time"
)

// ── Phase 1 reference (float64 era) ──────────────────────────────────────────

const piReference = "3.14159265358979323846"

// countCorrectDigits counts correct decimal digits in a float64 estimate.
// Used only during Phase 1 -- float64 can give at most 15.
func countCorrectDigits(f float64) int {
	s := fmt.Sprintf("%.20f", f)
	correct := 0
	for i := 0; i < len(piReference) && i < len(s); i++ {
		if s[i] == piReference[i] {
			correct++
		} else {
			break
		}
	}
	if correct < 2 {
		return 0
	}
	return correct - 2
}

// countCorrectDigitsBig counts correct decimal digits in a big.Float estimate.
// Used during Phase 2 -- references piForGauss for up to 3000+ digit comparison.
func countCorrectDigitsBig(pi *big.Float, prec uint) int {
	s   := pi.Text('f', int(prec/3)+10)
	ref := piForGauss
	correct := 0
	for i := 0; i < len(ref) && i < len(s); i++ {
		if s[i] == ref[i] {
			correct++
		} else {
			break
		}
	}
	if correct < 2 {
		return 0
	}
	return correct - 2
}

// ── Progress bar ──────────────────────────────────────────────────────────────

func progressBar(pct float64, width int) string {
	filled := int(math.Round(pct * float64(width)))
	if filled > width {
		filled = width
	}
	bar := strings.Repeat("█", filled) + strings.Repeat("░", width-filled)
	return fmt.Sprintf("[%s] %5.1f%%", bar, pct*100)
}

// ── Phase 1 milestone messages ────────────────────────────────────────────────

func dramaticMilestone(digits int) string {
	msgs := map[int]string{
		1:  "* First decimal digit locked -- 3.1",
		2:  "** 3.14 -- the famous digits begin!",
		3:  "*** 3.141 -- three digits confirmed!",
		4:  "**** 3.1415 -- four digits secured!",
		5:  "***** 3.14159 -- FIVE correct decimal digits!",
		6:  "****** 3.141592 -- six digits! Series converging hard.",
		7:  "******* 3.1415926 -- SEVEN digits. Nilakantha is flying!",
		8:  "******** 3.14159265 -- eight digits. Extraordinary precision.",
		9:  "********* 3.141592653 -- NINE decimal digits confirmed!",
		10: "********** 3.1415926535 -- TEN digits. You legend, Nilakantha.",
		11: "*********** 3.14159265358 -- ELEVEN digits. Deep Pi territory.",
		12: "************ 3.141592653589 -- TWELVE. Most calculators tap out here.",
		13: "************* 3.1415926535897 -- THIRTEEN digits. Rarefied air.",
		14: "************** 3.14159265358979 -- FOURTEEN. Almost float64 maximum!",
		15: "*************** 3.141592653589793 -- FIFTEEN digits. Float64 maxed out!",
	}
	if msg, ok := msgs[digits]; ok {
		return msg
	}
	return fmt.Sprintf("*** %d correct decimal digits!", digits)
}

// ── Phase 2 milestone messages ────────────────────────────────────────────────

func bigMilestone(digits int) string {
	switch {
	case digits == 16:
		return "COLOR:cyan:  !! 16 digits -- we just broke through the float64 wall !!"
	case digits == 17:
		return "COLOR:cyan:  !! 17 digits -- territory float64 cannot even see"
	case digits == 18:
		return "COLOR:cyan:  !! 18 digits -- Nilakantha in the 18th decimal place"
	case digits == 20:
		return "COLOR:cyan:  !! 20 digits -- two full decades of π"
	case digits%5 == 0:
		return fmt.Sprintf("COLOR:cyan:  !! %d digits -- big.Float climbing steadily", digits)
	default:
		return fmt.Sprintf("COLOR:green:  >> %d correct digits", digits)
	}
}

// ── Pacing for Phase 1 ────────────────────────────────────────────────────────

// sleepForTerm paces the Phase 1 summation loop in three acts:
// burst at the start, gentle ramp through the middle, plateau at the end.
func sleepForTerm(k, n int) {
	pct := float64(k) / float64(n)
	switch {
	case pct < 0.05:
		return
	case pct < 0.70:
		t  := (pct - 0.05) / 0.65
		us := 200.0 + 800.0*(math.Log1p(t*9)/math.Log1p(9))
		time.Sleep(time.Duration(us) * time.Microsecond)
	default:
		wobble := 50.0 * math.Sin(float64(k)*0.05)
		time.Sleep(time.Duration(800+wobble) * time.Microsecond)
	}
}

// ── Entry point ───────────────────────────────────────────────────────────────

// nifty_scoreBoardWeb runs both phases sequentially.
// n1 = Phase 1 terms (float64, concurrent goroutines)
// n2 = Phase 2 terms (big.Float, sequential, picks up where Phase 1 left off)
func nifty_scoreBoardWeb(n1, n2 int, done chan bool, webPrint func(string)) float64 {

	const bw   = 50
	const prec = uint(512) // big.Float precision in bits for Phase 2

	// ── Phase 1 header ────────────────────────────────────────────────────
	webPrint(boxSep(bw))
	webPrint(boxLine("  NILAKANTHA TWO-PHASE Pi ENGINE              ", bw))
	webPrint(boxLine("  π = 3 + 4/(2·3·4) − 4/(4·5·6) + ···       ", bw))
	webPrint(boxSep(bw))
	webPrint(boxLine("  PHASE 1: float64 · concurrent goroutines    ", bw))
	webPrint(boxLine(fmt.Sprintf("  Terms   : %d", n1), bw))
	webPrint(boxLine("  Ceiling : ~15 correct decimal digits        ", bw))
	webPrint(boxSep(bw))
	webPrint(boxLine("  PHASE 2: big.Float · 512-bit · sequential   ", bw))
	webPrint(boxLine(fmt.Sprintf("  Terms   : %d  (continues from Phase 1)", n2), bw))
	webPrint(boxLine("  Ceiling : genuine arbitrary precision        ", bw))
	webPrint(boxSep(bw))
	webPrint("")

	// ── Phase 1 ───────────────────────────────────────────────────────────

	webPrint("COLOR:yellow:  ── PHASE 1 BEGIN ──────────────────────────────────────")
	webPrint("")

	displayChan     := make(chan float64, 1)
	resultChan      := make(chan float64, 1)
	computationDone := make(chan bool, 1)
	termsCount      := 0

	phase1Start := time.Now()
	ticker      := time.NewTicker(time.Millisecond * 108)
	lastPi      := 3.0
	bestDigits  := 0

	webPrint(fmt.Sprintf("  Launching %d goroutines...", n1))

	go pi_nf_web(n1, done, webPrint, displayChan, resultChan, computationDone, &termsCount)

	// Live scoreboard ticker goroutine
	go func() {
		for range ticker.C {
			select {
			case <-done:
				ticker.Stop()
				return
			case piValue := <-displayChan:
				elapsed := time.Since(phase1Start)
				pct     := float64(termsCount) / float64(n1)
				bar     := progressBar(pct, 20)
				delta   := math.Abs(piValue - lastPi)
				digits  := countCorrectDigits(piValue)
				lastPi   = piValue

				if digits > bestDigits {
					bestDigits = digits
					webPrint(dramaticMilestone(digits))
				}
				webPrint(fmt.Sprintf(
					"UPDATE:  π ~ %.15f  %s  terms:%d  delta:%.2e  %s",
					piValue, bar, termsCount, delta,
					elapsed.Round(time.Millisecond),
				))
			default:
				elapsed := time.Since(phase1Start)
				pct     := float64(termsCount) / float64(n1)
				bar     := progressBar(pct, 20)
				webPrint(fmt.Sprintf(
					"UPDATE:  π ~ %.15f  %s  terms:%d  %s",
					lastPi, bar, termsCount,
					elapsed.Round(time.Millisecond),
				))
			}
		}
	}()

	// Wait for Phase 1 to complete
	var phase1Result float64
	for {
		select {
		case <-computationDone:
			ticker.Stop()
			phase1Result = <-resultChan
			phase1Elapsed := time.Since(phase1Start)
			finalDigits   := countCorrectDigits(phase1Result)

			// Build the match string for the summary
			fs    := fmt.Sprintf("%.15f", phase1Result)
			match := ""
			for i := 0; i < len(piReference) && i < len(fs); i++ {
				if fs[i] == piReference[i] {
					match += string(fs[i])
				} else {
					match += "~" + fs[i:]
					break
				}
			}

			webPrint("")
			webPrint(boxSep(bw))
			webPrint(boxLine("  PHASE 1 COMPLETE                            ", bw))
			webPrint(boxLine(fmt.Sprintf("  Time  : %s", phase1Elapsed.Round(time.Millisecond)), bw))
			webPrint(boxLine(fmt.Sprintf("  Terms : %d", termsCount), bw))
			webPrint(boxLine(fmt.Sprintf("  Digits: %d correct decimal places", finalDigits), bw))
			webPrint(boxLine(fmt.Sprintf("  Final : %.15f", phase1Result), bw))
			webPrint(boxLine(fmt.Sprintf("  Match : %s", match), bw))
			webPrint(boxSep(bw))
			webPrint("")
			webPrint("  float64 has given everything it has.")
			webPrint("  Fifteen digits is the wall.")
			webPrint("  Phase 2 will now pick up the series at the same k")
			webPrint("  and continue with big.Float at 512-bit precision.")
			webPrint("  Watch for new digits to appear beyond position 15.")
			webPrint("")

			goto phase2
		case <-done:
			ticker.Stop()
			webPrint("  !! Aborted by user during Phase 1.")
			return 0.0
		}
	}

phase2:
	// ── Phase 2 ───────────────────────────────────────────────────────────
	//
	// We restart the Nilakantha series from scratch using big.Float at
	// 512-bit precision. The series is not resumable from a float64 value
	// (the accumulated rounding error would pollute the big.Float result)
	// so we recompute all n1 terms quickly in big.Float, then continue
	// with n2 more terms, displaying progress as we go.
	//
	// The "wall break" moment -- when digits beyond 15 first appear --
	// is announced with special fanfare.

	webPrint("COLOR:yellow:  ── PHASE 2 BEGIN ──────────────────────────────────────")
	webPrint("")
	webPrint(fmt.Sprintf("  Recomputing %d Phase-1 terms in big.Float (fast)...", n1))
	webPrint("")

	phase2Start  := time.Now()
	bestBigDigits := 0
	wallBroken   := false

	// big.Float constants
	three := new(big.Float).SetPrec(prec).SetFloat64(3.0)
	four2 := new(big.Float).SetPrec(prec).SetFloat64(4.0)
	one   := new(big.Float).SetPrec(prec).SetFloat64(1.0)

	// Accumulator starts at 3
	piB := new(big.Float).SetPrec(prec).Set(three)

	// nilakanthaTerm returns the kth term as a big.Float
	// term_k = ±4 / (2k · (2k+1) · (2k+2))
	nilakanthaBigTerm := func(k int) *big.Float {
		j   := float64(2 * k)
		den := j * (j + 1) * (j + 2)
		t   := new(big.Float).SetPrec(prec).Quo(four2, new(big.Float).SetPrec(prec).SetFloat64(den))
		if k%2 == 0 {
			t.Neg(t)
		}
		return t
	}

	// Fast recompute of Phase 1 terms (no display, no sleep)
	for k := 1; k <= n1; k++ {
		select {
		case <-done:
			webPrint("  !! Aborted by user during Phase 2 recompute.")
			return 0.0
		default:
		}
		piB.Add(piB, nilakanthaBigTerm(k))
	}

	webPrint(fmt.Sprintf("  Phase-1 recompute complete. Continuing with %d new terms...", n2))
	webPrint("")

	// Seed blank row for first UPDATE:
	webPrint("")

	// Now run Phase 2 terms, displaying live progress
	updateInterval := n2 / 200 // update display ~200 times
	if updateInterval < 1 {
		updateInterval = 1
	}

	for k := n1 + 1; k <= n1+n2; k++ {
		select {
		case <-done:
			webPrint("  !! Aborted by user during Phase 2.")
			return 0.0
		default:
		}

		piB.Add(piB, nilakanthaBigTerm(k))

		// Check for new correct digits periodically
		if k%updateInterval == 0 || k == n1+n2 {
			digits  := countCorrectDigitsBig(piB, prec)
			pct     := float64(k-n1) / float64(n2)
			bar     := progressBar(pct, 20)
			elapsed := time.Since(phase2Start)
			piStr   := piB.Text('f', 20)

			// Announce wall break the first time we exceed 15 digits
			if !wallBroken && digits > 15 {
				wallBroken = true
				webPrint("")
				webPrint("COLOR:cyan:  ╔══════════════════════════════════════════════════╗")
				webPrint("COLOR:cyan:  ║  !! THE FLOAT64 WALL HAS BEEN BROKEN !!         ║")
				webPrint("COLOR:cyan:  ║  big.Float is now showing digits that           ║")
				webPrint("COLOR:cyan:  ║  float64 arithmetic cannot even represent.      ║")
				webPrint("COLOR:cyan:  ╚══════════════════════════════════════════════════╝")
				webPrint("")
			}

			if digits > bestBigDigits {
				bestBigDigits = digits
				webPrint(bigMilestone(digits))
			}

			webPrint(fmt.Sprintf(
				"UPDATE:  π ~ %s  %s  term:%d  %s",
				piStr, bar, k, elapsed.Round(time.Millisecond),
			))
		}
	}

	// ── Phase 2 final summary ─────────────────────────────────────────────

	phase2Elapsed := time.Since(phase2Start)
	finalDigits   := countCorrectDigitsBig(piB, prec)

	// Display string: verified digits only
	verifiedLen := finalDigits + 2 // +2 for "3."
	piStr       := piB.Text('f', finalDigits+5)
	if len(piStr) > verifiedLen {
		piStr = piStr[:verifiedLen]
	}

	webPrint("")
	webPrint(boxSep(bw))
	webPrint(boxLine("  PHASE 2 COMPLETE                            ", bw))
	webPrint(boxLine(fmt.Sprintf("  Time  : %s", phase2Elapsed.Round(time.Millisecond)), bw))
	webPrint(boxLine(fmt.Sprintf("  Terms : %d  (Phase 1) + %d  (Phase 2)", n1, n2), bw))
	webPrint(boxLine(fmt.Sprintf("  Digits: %d correct decimal places", finalDigits), bw))
	webPrint(boxSep(bw))
	webPrint("")
	webPrint(fmt.Sprintf("  π = %s", piStr))
	webPrint("")
	webPrint(boxSep(bw))
	webPrint(boxLine("  TWO-PHASE NILAKANTHA COMPLETE                ", bw))
	webPrint(boxLine(fmt.Sprintf("  Phase 1 ceiling  : 15 digits (float64)      ", ), bw))
	webPrint(boxLine(fmt.Sprintf("  Phase 2 achieved : %d digits (big.Float)  ", finalDigits), bw))
	webPrint(boxLine("  Kerala school · c. 1530 · still climbing     ", bw))
	webPrint(boxSep(bw))

	_ = phase1Result // acknowledged, not used -- Phase 2 recomputes cleanly
	_ = one
	return 0.0
}

// ── Phase 1 goroutine machinery ───────────────────────────────────────────────

func pi_nf_web(
	n int,
	done chan bool,
	webPrint func(string),
	displayChan chan float64,
	resultChan chan float64,
	computationDone chan bool,
	termsCount *int,
) float64 {

	ch := make(chan float64, n)
	f  := 3.0

	// Launch all n goroutines immediately -- this is the concurrent burst
	for k := 1; k <= n; k++ {
		select {
		case <-done:
			return f
		default:
			go nilakanthaTermWeb(ch, float64(k))
		}
	}

	// Collect results as they arrive, pacing with sleepForTerm
	for k := 1; k <= n; k++ {
		select {
		case <-done:
			return f
		case term := <-ch:
			*termsCount++
			f += term
			select {
			case displayChan <- f:
			default:
			}
			sleepForTerm(k, n)
		}
	}

	resultChan      <- f
	computationDone <- true
	return f
}

// nilakanthaTermWeb computes one term of the Nilakantha series
// and sends it to the channel. Runs as a goroutine in Phase 1.
func nilakanthaTermWeb(ch chan float64, k float64) {
	j := 2 * k
	if int64(k)%2 == 1 {
		ch <- 4.0 / (j * (j + 1) * (j + 2))
	} else {
		ch <- -4.0 / (j * (j + 1) * (j + 2))
	}
}
