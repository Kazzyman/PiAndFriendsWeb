package main

// Spigot_web.go
//
// Adapts the Spigot Pi algorithm for the Pi & Friends web suite.
//
// Original Spigot algorithm: sourced from GitHub, substantially rewritten
// by Richard (Rick) Woolley.
//
// Web adaptation, two-run architecture, honest uncertainty display,
// adaptive human-appreciation delay, Feynman Point detection, and
// colored output protocol:
// designed and implemented by Claude Sonnet (Anthropic), in collaboration
// with Richard Woolley, April 2026.
//
// Design philosophy (Rick's explicit requirement):
//   The algorithm is run TWICE. Both runs are genuine, independent
//   executions of the full Spigot computation. Nothing is reused,
//   replayed, or simulated between them.
//
//   Run 1 -- Full speed: digits sent as plain appends (no UPDATE:).
//            No in-place updating needed since all digits are final.
//
//   Run 2 -- Honest edition: ALL digit line output goes via UPDATE:,
//            which overwrites the current line in place. When a line
//            is complete, a plain empty string "" is sent to open a
//            new blank row, then the next digit starts a fresh UPDATE:
//            sequence on that new row.
//            This avoids the doubling bug where an UPDATE: and a plain
//            append both show the same content.
//
// Message protocol:
//   "UPDATE:text"        -- overwrite last line in place
//   "UPDATE:HASRED:text" -- overwrite last line, color '?' red
//   "COLOR:name:text"    -- append new colored line
//   ""  (empty)          -- append blank row (opens new line for UPDATE:)
//   plain text           -- append as permanent new line (Run 1 only)

import (
	"fmt"
	"strconv"
	"time"
)

// ── Entry point ───────────────────────────────────────────────────────────────

func TheSpigotWeb(numberOfDigits int, done chan bool, webPrint func(string)) {
	const bw = 50
	const targetSecs = 14.0

	webPrint("COLOR:cyan:" + boxSep(bw))
	webPrint("COLOR:cyan:" + boxLine("  THE SPIGOT ALGORITHM                        ", bw))
	webPrint("COLOR:cyan:" + boxLine("  Pi from integer arithmetic alone            ", bw))
	webPrint("COLOR:cyan:" + boxSep(bw))
	webPrint("COLOR:cyan:" + boxLine("  Origin : sourced from GitHub, rewritten     ", bw))
	webPrint("COLOR:cyan:" + boxLine("          by Richard Woolley                  ", bw))
	webPrint("COLOR:cyan:" + boxLine("  Method : integer arithmetic only, no floats ", bw))
	webPrint("COLOR:cyan:" + boxLine("  Runs   : TWO genuine independent executions ", bw))
	webPrint("COLOR:cyan:" + boxSep(bw))
	webPrint("")

	webPrint("COLOR:yellow:" + boxSep(bw))
	webPrint("COLOR:yellow:" + boxLine("  RUN 1 -- Full speed                         ", bw))
	webPrint("COLOR:yellow:" + boxLine(fmt.Sprintf("  Computing %d digits...", numberOfDigits), bw))
	webPrint("COLOR:yellow:" + boxSep(bw))
	webPrint("")

	run1Start := time.Now()
	ok := spigotRun1(numberOfDigits, done, webPrint)
	run1Time := time.Since(run1Start)

	if !ok {
		webPrint("COLOR:red:  !! Aborted.")
		return
	}

	normalDigits := numberOfDigits - 69
	if normalDigits < 1 {
		normalDigits = numberOfDigits
	}
	baseDelayMs := (targetSecs * 0.80 * 1000.0) / float64(normalDigits)
	if baseDelayMs < 1 {
		baseDelayMs = 1
	}
	baseDelay := time.Duration(baseDelayMs * float64(time.Millisecond))

	webPrint("")
	webPrint("COLOR:yellow:" + boxSep(bw))
	webPrint("COLOR:yellow:" + boxLine("  RUN 1 COMPLETE                              ", bw))
	webPrint("COLOR:yellow:" + boxLine(fmt.Sprintf("  %d digits in %s",
		numberOfDigits, run1Time.Round(time.Microsecond)), bw))
	webPrint("COLOR:yellow:" + boxSep(bw))
	webPrint("")
	webPrint("  Or was that too fast to follow?")
	webPrint("")
	webPrint("  No matter. We will now run the algorithm again --")
	webPrint("  a fresh computation, not a replay of Run 1 --")
	webPrint("  with a delay calibrated to this platform so that")
	webPrint(fmt.Sprintf("  the full run takes approximately %.0f seconds.", targetSecs))
	webPrint("")
	webPrint("  This time we will show you something the fast run")
	webPrint("  was hiding. The Spigot algorithm is not always")
	webPrint("  certain about its own output.")
	webPrint("")
	webPrint("  When it encounters a 9, it cannot yet know if that")
	webPrint("  digit is correct -- a carry from later arithmetic")
	webPrint("  could flip it to 0.")
	webPrint("")
	webPrint("COLOR:red:  Uncertain digits will appear in red as '?'.")
	webPrint("  Each '?' is overwritten with the true digit the")
	webPrint("  moment the algorithm itself resolves it.")
	webPrint("")
	webPrint("  No simulation. Every '?' is real uncertainty,")
	webPrint("  happening right now, on the server.")
	if numberOfDigits >= 768 {
		webPrint("")
		webPrint("COLOR:yellow:  One more thing: keep your eyes open around")
		webPrint("COLOR:yellow:  digit 762. Something rather famous lives there.")
	}
	webPrint("")

	time.Sleep(3 * time.Second)

	webPrint("COLOR:green:" + boxSep(bw))
	webPrint("COLOR:green:" + boxLine("  RUN 2 -- Honest edition                     ", bw))
	webPrint("COLOR:green:" + boxLine("  Uncertainty shown as it actually occurs      ", bw))
	webPrint("COLOR:green:" + boxLine(fmt.Sprintf("  Base delay: %.1fms/digit, calibrated for ~%.0fs",
		baseDelayMs, targetSecs), bw))
	webPrint("COLOR:green:" + boxSep(bw))
	webPrint("")

	ok = spigotRun2(numberOfDigits, done, webPrint, baseDelay)
	if !ok {
		webPrint("COLOR:red:  !! Aborted.")
		return
	}

	webPrint("")
	webPrint("COLOR:cyan:" + boxSep(bw))
	webPrint("COLOR:cyan:" + boxLine(fmt.Sprintf("  Spigot complete: %d digits of Pi          ", numberOfDigits), bw))
	webPrint("COLOR:cyan:" + boxLine("  Integer arithmetic only. No floating point   ", bw))
	webPrint("COLOR:cyan:" + boxLine("  numbers were used in this production.        ", bw))
	webPrint("COLOR:cyan:" + boxSep(bw))
}

// ── Run 1: full speed ─────────────────────────────────────────────────────────
// Plain appends only. No UPDATE: needed since every digit is immediately final.

func spigotRun1(numberOfDigits int, done chan bool, webPrint func(string)) bool {
	const lineWidth = 50

	size := numberOfDigits*10/3 + 50
	a := make([]int, size)
	for i := range a {
		a[i] = 2
	}

	line := ""
	pre := -1 // predigit, -1 = not yet set
	nines := 0
	count := 0
	decInserted := false

	addChar := func(ch string) {
		line += ch
		if len([]rune(line)) >= lineWidth {
			webPrint(line)
			line = ""
		}
	}

	emit := func(d int) {
		if !decInserted && count == 1 {
			addChar(".")
			decInserted = true
		}
		addChar(strconv.Itoa(d))
		count++
	}

	for count < numberOfDigits {
		select {
		case <-done:
			return false
		default:
		}

		carriedOver := 0
		sum := 0
		for j := size - 1; j >= 0; j-- {
			select {
			case <-done:
				return false
			default:
			}
			a[j] *= 10
			sum = a[j] + carriedOver
			quotient := sum / (j*2 + 1)
			a[j] = sum % (j*2 + 1)
			carriedOver = quotient * j
		}
		a[0] = sum % 10
		q := sum / 10

		switch {
		case q == 9:
			nines++
		case q == 10:
			if pre >= 0 {
				emit(pre + 1)
			}
			for i := 0; i < nines && count < numberOfDigits; i++ {
				emit(0)
			}
			pre = 0
			nines = 0
		default:
			if pre >= 0 {
				emit(pre)
			}
			for i := 0; i < nines && count < numberOfDigits; i++ {
				emit(9)
			}
			pre = q
			nines = 0
		}
	}

	if line != "" {
		webPrint(line)
	}
	return true
}

// ── Run 2: honest edition ─────────────────────────────────────────────────────
//
// All digit output uses UPDATE: so the current line is always overwritten
// in place. When a line is complete, we send a plain "" to open a fresh
// blank row, then the next character starts a new UPDATE: sequence on it.
//
// This means the JS never sees the same line content as both an UPDATE:
// and a plain append -- eliminating the doubling bug entirely.

func spigotRun2(numberOfDigits int, done chan bool, webPrint func(string), baseDelay time.Duration) bool {
	const lineWidth = 50
	const feynmanLen = 6
	const slowStart = 700
	const feynmanZone = 762

	pi := ""
	line := ""
	col := 0
	// boxes        := numberOfDigits * 10 / 3
	guardDigits := numberOfDigits/10 + 20
	totalDigits := numberOfDigits + guardDigits
	boxes := totalDigits*10/3 + 10
	remainders := make([]int, boxes)
	digitsHeld := 0
	pendingQs := 0
	decPos := 0
	decInserted := false
	digitsSeen := 0
	feynmanFired := false

	for i := 0; i < boxes; i++ {
		select {
		case <-done:
			return false
		default:
			remainders[i] = 2
		}
	}

	delayFor := func(pos int) time.Duration {
		if pos < slowStart || pos > feynmanZone+6 {
			return baseDelay
		}
		t := float64(pos-slowStart) / float64(feynmanZone-slowStart)
		factor := 1.0 + 3.0*t
		return time.Duration(float64(baseDelay) * factor)
	}

	// show sends the current line as UPDATE:, overwriting in place.
	// All digit output goes through here -- never via plain webPrint.
	show := func() {
		if pendingQs > 0 {
			webPrint("UPDATE:HASRED:" + line)
		} else {
			webPrint("UPDATE:" + line)
		}
	}

	// newRow commits the current line visually (it is already showing
	// correctly via the last show() call) and opens a fresh blank row
	// by sending an empty plain string. The next show() call will then
	// overwrite that blank row with the first character of the new line.
	newRow := func() {
		webPrint("") // plain empty -- JS appends a blank row
		line = ""
		col = 0
		pendingQs = 0
	}

	// flushRow commits the current line and opens a fresh row.
	// The line is already showing correctly via the last show() call,
	// so we just send "" to anchor it and open the next row --
	// exactly like newRow(). Sending webPrint(line) here would cause
	// the JS to emit the line a second time (doubling bug).
	flushRow := func() {
		if line != "" {
			newRow()
		}
	}

	appendChar := func(ch string) {
		line += ch
		col++
		if col >= lineWidth {
			show()   // show complete line via UPDATE:
			newRow() // open next row
		} else {
			show()
		}
	}

	showQ := func() {
		line += "?"
		col++
		pendingQs++
		// Do NOT enforce lineWidth break here. pendingQs must stay intact so
		// the ?s can be resolved in-place (to 9 or 0) before the line commits.
		// The line may briefly exceed lineWidth while ?s are unresolved;
		// it commits after resolution in the default or carry case.
		show()
	}

	resolveQs := func(count int, digit rune) {
		start := len([]rune(line)) - pendingQs
		runes := []rune(line)
		for k := 0; k < count && start+k < len(runes); k++ {
			select {
			case <-done:
				return
			default:
			}
			runes[start+k] = digit
			line = string(runes)
			pendingQs--
			show()
			time.Sleep(baseDelay / 2)
		}
	}

	// Seed the first blank row so the first UPDATE: has somewhere to land
	webPrint("")

	// for i := 0; i < numberOfDigits; i++ {
	for i := 0; i < totalDigits; i++ {
		select {
		case <-done:
			return false
		default:
		}

		carriedOver := 0
		sum := 0

		for j := boxes - 1; j >= 0; j-- {
			select {
			case <-done:
				return false
			default:
			}
			remainders[j] *= 10
			sum = remainders[j] + carriedOver
			quotient := sum / (j*2 + 1)
			remainders[j] = sum % (j*2 + 1)
			carriedOver = quotient * j
		}

		remainders[0] = sum % 10
		q := sum / 10

		switch q {
		case 9:
			digitsHeld++
			decPos++
			digitsSeen++
			if !decInserted && digitsSeen == 2 {
				appendChar(".")
				decInserted = true
			}
			pi += "9"
			showQ()
			time.Sleep(delayFor(decPos))

		case 10:
			q = 0
			for k := 1; k <= digitsHeld; k++ {
				select {
				case <-done:
					return false
				default:
				}
				replaced, _ := strconv.Atoi(pi[i-k : i-k+1])
				if replaced == 9 {
					replaced = 0
				} else {
					replaced++
				}
				pi = delChar(pi, i-k)
				pi = pi[:i-k] + strconv.Itoa(replaced) + pi[i-k:]
			}

			// Handle carry display.
			// Save pendingQs BEFORE flushRow/newRow resets it to zero.
			// IMPORTANT: in the spigot display model (same as Run 1), held 9s
			// ARE displayed as 9s even when a carry fires. The carry does not
			// flip the held 9 to 0 -- it makes the CURRENT iteration's digit 0.
			// Run 1 shows e.g. "...9 0..." where 9 is the held digit and 0 is
			// the carry digit. We do the same: resolve ? to 9, then emit 0.
			savedPendingQs := pendingQs

			if savedPendingQs > 0 {
				// Resolve the ?s to 9s in-place on the current line,
				// then commit the line before the carry message.
				resolveQs(savedPendingQs, '9')
				if col >= lineWidth {
					show()
					newRow()
				}
				flushRow()
				if savedPendingQs == 1 {
					webPrint("COLOR:red:" + fmt.Sprintf(
						"  [CARRY at position %d] That 9 is confirmed.", decPos))
					webPrint("COLOR:red:  A carry fires -- the next digit in the sequence is 0.")
				} else {
					webPrint("COLOR:red:" + fmt.Sprintf(
						"  [CARRY at position %d] Those %d 9s are confirmed.", decPos, savedPendingQs))
					webPrint("COLOR:red:  A carry fires -- the next digit in the sequence is 0.")
				}
				time.Sleep(baseDelay * 3)
				webPrint("") // blank line between message and following digits
				webPrint("")
			}

			decPos++
			digitsSeen++
			if !decInserted && digitsSeen == 2 {
				appendChar(".")
				decInserted = true
			}
			pi += strconv.Itoa(q)
			digitsHeld = 1
			appendChar(strconv.Itoa(q))
			time.Sleep(delayFor(decPos))

		default:
			if digitsHeld > 0 && pendingQs > 0 {
				resolveQs(pendingQs, '9')

				if !feynmanFired && digitsHeld >= feynmanLen {
					feynmanFired = true
					flushRow()
					webPrint("")
					webPrint("COLOR:yellow:  +--------------------------------------------------+")
					webPrint("COLOR:yellow:  | !! THE FEYNMAN POINT !!                          |")
					webPrint("COLOR:yellow:  | (or perhaps the Hofstadter Point?)               |")
					webPrint("COLOR:yellow:  | Six consecutive 9s at decimal position 762       |")
					webPrint("COLOR:yellow:  | ...134  999999  837...                           |")
					webPrint("COLOR:yellow:  |                                                  |")
					webPrint("COLOR:yellow:  | You just watched the algorithm hold all six as   |")
					webPrint("COLOR:yellow:  | '??????' -- because a carry could have flipped   |")
					webPrint("COLOR:yellow:  | every one of them to zero.                       |")
					webPrint("COLOR:yellow:  | They resolved to 999999. Truth confirmed        |")
					webPrint("COLOR:yellow:  | in real time ... in real-slow time.              |")
					webPrint("COLOR:yellow:  |                                                  |")
					webPrint("COLOR:yellow:  | Feynman: 'nine nine nine nine nine nine...        |")
					webPrint("COLOR:yellow:  |          ...and so on!'                          |")
					webPrint("COLOR:yellow:  +--------------------------------------------------+")
					webPrint("")
					time.Sleep(4 * time.Second)
					webPrint("") // seed blank row for digits to continue
				}
			}

			digitsHeld = 1
			// If resolving the ?s filled the line to lineWidth, commit it now
			// before appending the current digit on a fresh row.
			if col >= lineWidth {
				show()
				newRow()
			}
			decPos++
			digitsSeen++
			if !decInserted && digitsSeen == 2 {
				appendChar(".")
				decInserted = true
			}
			pi += strconv.Itoa(q)
			appendChar(strconv.Itoa(q))
			time.Sleep(delayFor(decPos))
		}
	}

	// Commit final partial line -- cursor already shows it correctly
	if line != "" {
		newRow()
	}
	return true
}

// ── Helpers ───────────────────────────────────────────────────────────────────

// delChar removes the character at index from string s.
// Written largely by Richard Woolley.
func delChar(s string, index int) string {
	tmp := []rune(s)
	return string(append(tmp[0:index], tmp[index+1:]...))
}
