package main

// ScoreBoardNila_web.go
//
// Nilakantha Somayaji's alternating series for π (c. 1530):
//   π = 3 + 4/(2·3·4) − 4/(4·5·6) + 4/(6·7·8) − ···
//
// Convergence rate: error after N terms ≈ 1/(2N³)
// This means each factor of 10 in terms buys 3 more correct digits:
//   10,000 terms  →  ~12 digits
//   100,000 terms →  ~15 digits   (float64 wall)
//   1,000,000     →  ~18 digits   (wall broken!)
//   5,000,000     →  ~21 digits
//
// Phase 1 — float64, concurrent goroutines
//   Live scoreboard, milestone fanfare as each digit locks in.
//   Caps at ~15 digits no matter how long you run it.
//   That cap IS the demonstration.
//
// Phase 2 — big.Float, 512-bit, sequential
//   Runs from term 1 through n1+n2, accumulating in big.Float.
//   At ~100,000-200,000 terms the wall breaks and new digits appear.
//   Progress updates every 10,000 terms with live digit count.
//   The wall-break moment is announced with fanfare.
//
// Adapted by Richard Woolley with a lot of help from Claude.

import (
	"fmt"
	"math"
	"math/big"
	"strings"
	"time"
)

// ── Phase 1 reference (float64 era, 20 digits is plenty) ─────────────────────

const piReference = "3.14159265358979323846"

// countCorrectDigits counts correct decimal digits in a float64 estimate.
func countCorrectDigits(f float64) int {
	s       := fmt.Sprintf("%.20f", f)
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
// Uses piForGauss which has 3000+ verified digits.
func countCorrectDigitsBig(pi *big.Float) int {
	s   := pi.Text('f', 35)
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
	return fmt.Sprintf("*** %d correct decimal digits -- beyond float64 resolution!", digits)
}

// ── Phase 2 milestone messages ────────────────────────────────────────────────

func bigMilestone(digits int) string {
	switch digits {
	case 16:
		return "COLOR:cyan:  !! 16 digits -- we just broke through the float64 wall !!"
	case 17:
		return "COLOR:cyan:  !! 17 digits -- territory float64 cannot even represent"
	case 18:
		return "COLOR:cyan:  !! 18 digits -- Nilakantha in the 18th decimal place"
	case 19:
		return "COLOR:cyan:  !! 19 digits -- Kerala school, c.1530, still climbing"
	case 20:
		return "COLOR:cyan:  !! 20 digits -- two full decades of π confirmed"
	default:
		if digits > 20 {
			return fmt.Sprintf("COLOR:cyan:  !! %d digits -- big.Float soaring", digits)
		}
		return fmt.Sprintf("COLOR:green:  >> %d correct digits", digits)
	}
}

// ── Pacing for Phase 1 ────────────────────────────────────────────────────────

// sleepForTerm paces the Phase 1 loop: burst → ramp → plateau.
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

// nifty_scoreBoardWeb runs both phases.
//   n1 = Phase 1 terms (float64, concurrent goroutines, default 5000)
//   n2 = Phase 2 terms (big.Float, sequential, default 1,000,000)
func nifty_scoreBoardWeb(n1, n2 int, done chan bool, webPrint func(string)) float64 {

	const bw   = 50
	const prec = uint(512)

	totalTerms := n1 + n2

	// ── Header ────────────────────────────────────────────────────────────
	webPrint(boxSep(bw))
	webPrint(boxLine("  NILAKANTHA TWO-PHASE Pi ENGINE              ", bw))
	webPrint(boxLine("  π = 3 + 4/(2·3·4) − 4/(4·5·6) + ···       ", bw))
	webPrint(boxSep(bw))
	webPrint(boxLine("  PHASE 1: float64 · concurrent goroutines    ", bw))
	webPrint(boxLine(fmt.Sprintf("  Terms   : %d", n1), bw))
	webPrint(boxLine("  Ceiling : ~15 correct decimal digits        ", bw))
	webPrint(boxSep(bw))
	webPrint(boxLine("  PHASE 2: big.Float · 512-bit · sequential   ", bw))
	webPrint(boxLine(fmt.Sprintf("  Terms   : %d additional", n2), bw))
	webPrint(boxLine(fmt.Sprintf("  Total   : %d terms", totalTerms), bw))
	webPrint(boxLine("  Ceiling : genuine arbitrary precision        ", bw))
	webPrint(boxLine("  Error   : ~1/(2N³): 10x terms = 3 new digits", bw))
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

	go pi_nf_web(n1, done, displayChan, resultChan, computationDone, &termsCount)

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

	// Wait for Phase 1 to finish
	for {
		select {
		case <-computationDone:
			ticker.Stop()
			phase1Final   := <-resultChan
			phase1Elapsed := time.Since(phase1Start)
			finalDigits   := countCorrectDigits(phase1Final)

			fs    := fmt.Sprintf("%.15f", phase1Final)
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
			webPrint(boxLine(fmt.Sprintf("  Final : %.15f", phase1Final), bw))
			webPrint(boxLine(fmt.Sprintf("  Match : %s", match), bw))
			webPrint(boxSep(bw))
			webPrint("")
			webPrint("  float64 has given everything it has.")
			webPrint("  The wall is real. ~15 digits is the ceiling.")
			webPrint("")
			webPrint(fmt.Sprintf("  Phase 2 will now run %d terms in big.Float", n2))
			webPrint("  at 512-bit precision (~154 decimal digits of headroom).")
			webPrint("  The wall should break somewhere around term 150,000.")
			webPrint("  This will take a few minutes. Watch the digit count.")
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
	// We recompute from term 1 entirely in big.Float. We cannot resume
	// from the float64 accumulator -- its accumulated rounding error sits
	// right at the digits we are trying to reveal. Starting fresh keeps
	// the big.Float accumulator clean for all n1+n2 terms.

	webPrint("COLOR:yellow:  ── PHASE 2 BEGIN ──────────────────────────────────────")
	webPrint("")
	webPrint(fmt.Sprintf("  Recomputing Phase-1 terms (%d) silently in big.Float...", n1))

	phase2Start   := time.Now()
	bestBigDigits := 0
	wallBroken    := false

	two  := new(big.Float).SetPrec(prec).SetFloat64(2.0)
	four := new(big.Float).SetPrec(prec).SetFloat64(4.0)

	piB      := new(big.Float).SetPrec(prec).SetFloat64(3.0)
	digitOne := new(big.Float).SetPrec(prec).SetFloat64(2.0)
	digitTwo := new(big.Float).SetPrec(prec).SetFloat64(3.0)
	digitThr := new(big.Float).SetPrec(prec).SetFloat64(4.0)

	// Apply k=1 term: 3 + 4/(2*3*4)
	firstTerm := new(big.Float).SetPrec(prec).Quo(four,
		new(big.Float).SetPrec(prec).Mul(digitOne,
			new(big.Float).SetPrec(prec).Mul(digitTwo, digitThr)))
	piB.Add(piB, firstTerm)

	// Silent recompute: k=2 through n1
	for k := 2; k <= n1; k++ {
		select {
		case <-done:
			webPrint("  !! Aborted during Phase 2 silent recompute.")
			return 0.0
		default:
		}
		digitOne.Add(digitOne, two)
		digitTwo.Add(digitTwo, two)
		digitThr.Add(digitThr, two)

		term := new(big.Float).SetPrec(prec).Quo(four,
			new(big.Float).SetPrec(prec).Mul(digitOne,
				new(big.Float).SetPrec(prec).Mul(digitTwo, digitThr)))

		if k%2 == 0 {
			piB.Sub(piB, term)
		} else {
			piB.Add(piB, term)
		}
	}

	silentElapsed := time.Since(phase2Start)
	webPrint(fmt.Sprintf("  Silent recompute done in %s.", silentElapsed.Round(time.Millisecond)))
	webPrint(fmt.Sprintf("  Now running %d new terms with live display...", n2))
	webPrint("")
	webPrint("") // seed blank row for first UPDATE:

	const updateEvery = 10000

	for k := n1 + 1; k <= totalTerms; k++ {
		select {
		case <-done:
			webPrint("  !! Aborted by user during Phase 2.")
			return 0.0
		default:
		}

		digitOne.Add(digitOne, two)
		digitTwo.Add(digitTwo, two)
		digitThr.Add(digitThr, two)

		term := new(big.Float).SetPrec(prec).Quo(four,
			new(big.Float).SetPrec(prec).Mul(digitOne,
				new(big.Float).SetPrec(prec).Mul(digitTwo, digitThr)))

		if k%2 == 0 {
			piB.Sub(piB, term)
		} else {
			piB.Add(piB, term)
		}

		if k%updateEvery == 0 || k == totalTerms {
			digits  := countCorrectDigitsBig(piB)
			pct     := float64(k-n1) / float64(n2)
			bar     := progressBar(pct, 20)
			elapsed := time.Since(phase2Start)
			piStr   := piB.Text('f', 22)

			if !wallBroken && digits > 15 {
				wallBroken = true
				webPrint("")
				webPrint("COLOR:cyan:  ╔══════════════════════════════════════════════════╗")
				webPrint("COLOR:cyan:  ║  !! THE FLOAT64 WALL HAS BEEN BROKEN !!         ║")
				webPrint("COLOR:cyan:  ║  big.Float is now showing digits that           ║")
				webPrint("COLOR:cyan:  ║  float64 cannot even represent.                 ║")
				webPrint("COLOR:cyan:  ║  Nilakantha Somayaji, c.1530, still delivering. ║")
				webPrint("COLOR:cyan:  ╚══════════════════════════════════════════════════╝")
				webPrint("")
				webPrint("") // re-seed for UPDATE:
			}

			if digits > bestBigDigits {
				bestBigDigits = digits
				webPrint(bigMilestone(digits))
				webPrint("") // re-seed after milestone
			}

			webPrint(fmt.Sprintf(
				"UPDATE:  π ~ %s  %s  term:%d  %s",
				piStr, bar, k, elapsed.Round(time.Millisecond),
			))
		}
	}

	// ── Final summary ─────────────────────────────────────────────────────

	phase2Elapsed := time.Since(phase2Start)
	finalDigits   := countCorrectDigitsBig(piB)

	displayStr := piB.Text('f', finalDigits+2)
	if len(displayStr) > finalDigits+2 {
		displayStr = displayStr[:finalDigits+2]
	}

	webPrint("")
	webPrint(boxSep(bw))
	webPrint(boxLine("  PHASE 2 COMPLETE                            ", bw))
	webPrint(boxLine(fmt.Sprintf("  Time     : %s", phase2Elapsed.Round(time.Millisecond)), bw))
	webPrint(boxLine(fmt.Sprintf("  New terms: %d", n2), bw))
	webPrint(boxLine(fmt.Sprintf("  Total    : %d terms", totalTerms), bw))
	webPrint(boxLine(fmt.Sprintf("  Digits   : %d correct decimal places", finalDigits), bw))
	webPrint(boxSep(bw))
	webPrint("")
	webPrint(fmt.Sprintf("  π = %s", displayStr))
	webPrint("")
	webPrint(boxSep(bw))
	webPrint(boxLine("  TWO-PHASE NILAKANTHA COMPLETE                ", bw))
	webPrint(boxLine("  Phase 1 ceiling  : 15 digits  (float64)     ", bw))
	webPrint(boxLine(fmt.Sprintf("  Phase 2 achieved : %d digits  (big.Float)  ", finalDigits), bw))
	webPrint(boxLine("  Kerala school · c.1530 · predates Newton     ", bw))
	webPrint(boxLine("  by 150 years · still climbing                ", bw))
	webPrint(boxSep(bw))

	return 0.0
}

// ── Phase 1 goroutine machinery ───────────────────────────────────────────────

func pi_nf_web(
    n int,
    done chan bool,
    displayChan chan float64,
	resultChan chan float64,
	computationDone chan bool,
	termsCount *int,
) float64 {

	ch := make(chan float64, n)
	f  := 3.0

	for k := 1; k <= n; k++ {
		select {
		case <-done:
			return f
		default:
			go nilakanthaTermWeb(ch, float64(k))
		}
	}

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

// nilakanthaTermWeb computes one term and sends it to the channel.
func nilakanthaTermWeb(ch chan float64, k float64) {
	j := 2 * k
	if int64(k)%2 == 1 {
		ch <- 4.0 / (j * (j + 1) * (j + 2))
	} else {
		ch <- -4.0 / (j * (j + 1) * (j + 2))
	}
}
