package main

import (
	"fmt"
	"math"
	"strings"
	"time"
)

const piReference = "3.14159265358979323846"

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

func progressBar(pct float64, width int) string {
	filled := int(math.Round(pct * float64(width)))
	if filled > width {
		filled = width
	}
	bar := strings.Repeat("█", filled) + strings.Repeat("░", width-filled)
	return fmt.Sprintf("[%s] %5.1f%%", bar, pct*100)
}

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

// boxSep returns a separator line: +----...----+  (inner width w dashes)
func boxSep(w int) string {
	return "+" + strings.Repeat("-", w) + "+"
}

// boxLine returns a content line padded to exactly w inner chars: |content...|
func boxLine(content string, w int) string {
	// Truncate if somehow over (shouldn't happen with correct literals)
	if len(content) > w {
		content = content[:w]
	}
	return "|" + content + strings.Repeat(" ", w-len(content)) + "|"
}

// sleepForTerm paces the summation loop in three acts.
// Total sleep across 5,000 terms: ~3.6 seconds.
func sleepForTerm(k, n int) {
	pct := float64(k) / float64(n)
	switch {
	case pct < 0.05:
		return
	case pct < 0.70:
		t := (pct - 0.05) / 0.65
		us := 200.0 + 800.0*(math.Log1p(t*9)/math.Log1p(9))
		time.Sleep(time.Duration(us) * time.Microsecond)
	default:
		wobble := 50.0 * math.Sin(float64(k)*0.05)
		time.Sleep(time.Duration(800+wobble) * time.Microsecond)
	}
}

// nifty_scoreBoardWeb runs the concurrent Nilakantha Pi calculation
// with a dramatic live scoreboard streamed to the web client.
//
// Message protocol:
//   Normal lines  --> appended as new rows
//   "UPDATE:..."  --> JavaScript overwrites the last line in place
func nifty_scoreBoardWeb(done chan bool, webPrint func(string)) float64 {

	displayChan     := make(chan float64, 1)
	resultChan      := make(chan float64, 1)
	computationDone := make(chan bool, 1)
	termsCount      := 0
	const n          = 5000
	const bw         = 50 // inner box width -- sep line has exactly bw dashes

	startTime  := time.Now()
	ticker     := time.NewTicker(time.Millisecond * 108)
	lastPi     := 3.0
	bestDigits := 0

	webPrint(boxSep(bw))
	webPrint(boxLine("  NILAKANTHA SERIES -- CONCURRENT Pi ENGINE  ", bw))
	webPrint(boxLine("  pi = 3 + 4/(2*3*4) - 4/(4*5*6) + ...      ", bw))
	webPrint(boxSep(bw))
	webPrint(boxLine(fmt.Sprintf("  Terms planned : %d goroutines", n), bw))
	webPrint(boxLine("  Reference     : 3.14159265358979323846     ", bw))
	webPrint(boxLine("  Pacing        : burst -> ramp -> plateau   ", bw))
	webPrint(boxSep(bw))
	webPrint("  >> Launching goroutines...")

	go pi_nf_web(n, done, webPrint, displayChan, resultChan, computationDone, &termsCount)

	go func() {
		for range ticker.C {
			select {
			case <-done:
				ticker.Stop()
				return
			case piValue := <-displayChan:
				elapsed := time.Since(startTime)
				pct     := float64(termsCount) / float64(n)
				bar     := progressBar(pct, 20)
				delta   := math.Abs(piValue - lastPi)
				digits  := countCorrectDigits(piValue)
				lastPi   = piValue

				if digits > bestDigits {
					bestDigits = digits
					webPrint(dramaticMilestone(digits))
				}

				webPrint(fmt.Sprintf(
					"UPDATE:  pi ~ %.15f  %s  terms:%d  delta:%.2e  time:%s",
					piValue, bar, termsCount, delta,
					elapsed.Round(time.Millisecond),
				))
			default:
				elapsed := time.Since(startTime)
				pct     := float64(termsCount) / float64(n)
				bar     := progressBar(pct, 20)
				webPrint(fmt.Sprintf(
					"UPDATE:  pi ~ %.15f  %s  terms:%d  time:%s",
					lastPi, bar, termsCount,
					elapsed.Round(time.Millisecond),
				))
			}
		}
	}()

	for {
		select {
		case <-computationDone:
			ticker.Stop()
			final   := <-resultChan
			elapsed := time.Since(startTime)
			digits  := countCorrectDigits(final)

			fs := fmt.Sprintf("%.15f", final)
			match := ""
			for i := 0; i < len(piReference) && i < len(fs); i++ {
				if fs[i] == piReference[i] {
					match += string(fs[i])
				} else {
					match += "~" + fs[i:]
					break
				}
			}

			webPrint(boxSep(bw))
			webPrint(boxLine(fmt.Sprintf("  DONE in %s", elapsed.Round(time.Millisecond)), bw))
			webPrint(boxLine(fmt.Sprintf("  Terms : %d / %d (100%%)", termsCount, n), bw))
			webPrint(boxLine(fmt.Sprintf("  Digits: %d correct decimal places", digits), bw))
			webPrint(boxLine(fmt.Sprintf("  Final : %.15f", final), bw))
			webPrint(boxLine(fmt.Sprintf("  Ref   : %.15f", 3.14159265358979323846), bw))
			webPrint(boxLine(fmt.Sprintf("  Match : %s", match), bw))
			webPrint(boxSep(bw))
			webPrint(boxLine("         Pi  COMPUTATION  DONE               ", bw))
			webPrint(boxSep(bw))
			return final

		case <-done:
			ticker.Stop()
			webPrint("  !! Aborted by user.")
			return 0.0
		}
	}
}

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

func nilakanthaTermWeb(ch chan float64, k float64) {
	j := 2 * k
	if int64(k)%2 == 1 {
		ch <- 4.0 / (j * (j + 1) * (j + 2))
	} else {
		ch <- -4.0 / (j * (j + 1) * (j + 2))
	}
}
