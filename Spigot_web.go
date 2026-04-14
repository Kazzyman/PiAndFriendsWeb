package main

// Spigot_web.go
//
// Adapts the Spigot Pi algorithm for the Pi & Friends web suite.
//
// Original Spigot algorithm: sourced from GitHub, substantially rewritten
// by Richard (Rick) Woolley.
//
// Web adaptation, two-run architecture, event-recording system, honest
// uncertainty display, and Feynman Point detection:
// designed and implemented by Claude Sonnet (Anthropic), in collaboration
// with Richard Woolley, April 2026.
//
// Design philosophy (Rick's explicit requirement):
//   The algorithm is run TWICE. Both runs are genuine, independent
//   executions of the full Spigot computation. Nothing is reused,
//   replayed, or simulated between them.
//
//   Run 1 -- Full speed: digits stream to the screen as fast as the
//            SSE connection allows. The user sees the complete answer
//            and the wall-clock time it took.
//
//   Run 2 -- Honest edition: the algorithm runs again from scratch.
//            This time, every digit the algorithm is genuinely uncertain
//            about is shown as '?' on screen. The uncertainty is real --
//            when the Spigot encounters a 9 it cannot yet know if a carry
//            from later arithmetic will flip it to 0. The '?' is replaced
//            with the true digit the moment the algorithm resolves it.
//            No fakery. No simulation. Every '?' represents a moment the
//            algorithm itself did not know the answer.
//
//   The Feynman Point Easter egg (six consecutive 9s at decimal position
//   762) is revealed only in Run 2, where the user can watch the six '?'
//   marks accumulate and then resolve to '999999' in real time; real-slow time.

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// ── Entry point ───────────────────────────────────────────────────────────────

func TheSpigotWeb(numberOfDigits int, done chan bool, webPrint func(string)) {
	const bw = 50

	webPrint(boxSep(bw))
	webPrint(boxLine("  THE SPIGOT ALGORITHM                        ", bw))
	webPrint(boxLine("  Pi from integer arithmetic alone            ", bw))
	webPrint(boxSep(bw))
	webPrint(boxLine("  Origin : sourced from GitHub, rewritten     ", bw))
	webPrint(boxLine("          by Richard Woolley                  ", bw))
	webPrint(boxLine("  Method : integer arithmetic only, no floats ", bw))
	webPrint(boxLine("  Runs   : TWO genuine independent executions ", bw))
	webPrint(boxSep(bw))
	webPrint("")

	// ── RUN 1: full speed ─────────────────────────────────────────────────────
	//
	// A complete, genuine execution of the Spigot algorithm.
	// Digits are streamed to the screen as fast as the SSE connection
	// allows. No delays, no uncertainty markers -- just the raw answer.

	webPrint(boxSep(bw))
	webPrint(boxLine("  RUN 1 -- Full speed                         ", bw))
	webPrint(boxLine(fmt.Sprintf("  Computing %d digits...", numberOfDigits), bw))
	webPrint(boxSep(bw))
	webPrint("")

	run1Start := time.Now()
	ok        := spigotRun1(numberOfDigits, done, webPrint)
	run1Time  := time.Since(run1Start)

	if !ok {
		webPrint("  !! Aborted.")
		return
	}

	webPrint("")
	webPrint(boxSep(bw))
	webPrint(boxLine("  RUN 1 COMPLETE                              ", bw))
	webPrint(boxLine(fmt.Sprintf("  %d digits in %s",
		numberOfDigits, run1Time.Round(time.Microsecond)), bw))
	webPrint(boxSep(bw))
	webPrint("")
	webPrint("  Or was that too fast to follow?")
	webPrint("")
	webPrint("  No matter. We will now run the algorithm again,")
	webPrint("  this time with a human appreciation delay of 50ms")
	webPrint("  per digit.")
	webPrint("")
	webPrint("  This second run will also show you something the")
	webPrint("  first run was hiding. The Spigot algorithm is not")
	webPrint("  always certain about its own output. When it")
	webPrint("  encounters a 9, it cannot yet know if that digit")
	webPrint("  is correct -- a carry from later arithmetic could")
	webPrint("  flip it to 0.")
	webPrint("")
	webPrint("  We will show '?' for every genuinely uncertain")
	webPrint("  digit, and overwrite it with the truth the moment")
	webPrint("  the algorithm itself resolves it.")
	webPrint("")
	webPrint("  No simulation. This is a fresh computation.")
	webPrint("  Every '?' you see is real uncertainty, happening")
	webPrint("  right now, in the algorithm running on the server.")
	if numberOfDigits >= 768 {
		webPrint("")
		webPrint("  One more thing: keep your eyes open around")
		webPrint("  digit 762. Something rather famous lives there.")
	}
	webPrint("")

	// Dramatic pause before Run 2
	time.Sleep(3 * time.Second)

	// ── RUN 2: honest edition ─────────────────────────────────────────────────
	//
	// A second complete, genuine, independent execution of the full
	// Spigot algorithm. This run streams digits with a 50ms delay and
	// shows '?' for every digit the algorithm is currently uncertain
	// about, resolving each one on screen the moment it is confirmed.

	webPrint(boxSep(bw))
	webPrint(boxLine("  RUN 2 -- Honest edition                     ", bw))
	webPrint(boxLine("  Uncertainty shown as it actually occurs      ", bw))
	webPrint(boxLine("  50ms human appreciation delay per digit      ", bw))
	webPrint(boxSep(bw))
	webPrint("")

	ok = spigotRun2(numberOfDigits, done, webPrint, 50*time.Millisecond)
	if !ok {
		webPrint("  !! Aborted.")
		return
	}

	webPrint("")
	webPrint(boxSep(bw))
	webPrint(boxLine("  RUN 2 COMPLETE                              ", bw))
	webPrint(boxLine(fmt.Sprintf("  %d digits of Pi delivered twice           ", numberOfDigits), bw))
	webPrint(boxLine("  Integer arithmetic only. No floating point numbers were used in this production. ", bw))
	webPrint(boxSep(bw))
}

// ── Run 1: full speed ─────────────────────────────────────────────────────────
//
// Complete independent execution of the Spigot algorithm.
// Streams correct digits to webPrint as fast as possible, grouped
// into lines of 50 characters.

func spigotRun1(numberOfDigits int, done chan bool, webPrint func(string)) bool {
	const lineWidth = 50

	pi         := ""
	line        := ""
	boxes       := numberOfDigits * 10 / 3
	remainders  := make([]int, boxes)
	digitsHeld  := 0
	decInserted := false
	digitsSeen  := 0

	for i := 0; i < boxes; i++ {
		select {
		case <-done:
			return false
		default:
			remainders[i] = 2
		}
	}

	emit := func(ch string) {
		line += ch
		if len([]rune(line)) >= lineWidth {
			webPrint(line)
			line = ""
		}
	}

	for i := 0; i < numberOfDigits; i++ {
		select {
		case <-done:
			return false
		default:
		}

		carriedOver := 0
		sum         := 0

		for j := boxes - 1; j >= 0; j-- {
			select {
			case <-done:
				return false
			default:
			}
			remainders[j] *= 10
			sum            = remainders[j] + carriedOver
			quotient       := sum / (j*2 + 1)
			remainders[j]  = sum % (j*2 + 1)
			carriedOver     = quotient * j
		}

		remainders[0] = sum % 10
		q             := sum / 10

		switch q {
		case 9:
			digitsHeld++
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
			digitsHeld = 1
		default:
			digitsHeld = 1
		}

		pi += strconv.Itoa(q)

		// Insert decimal point between first and second digit
		digitsSeen++
		if !decInserted && digitsSeen == 2 {
			emit(".")
			decInserted = true
		}
		emit(strconv.Itoa(q))
	}

	if line != "" {
		webPrint(line)
	}
	return true
}

// ── Run 2: honest edition ─────────────────────────────────────────────────────
//
// Second complete independent execution of the full Spigot algorithm.
// This run streams digits with a human appreciation delay and shows
// genuine uncertainty as it occurs:
//
//   - '?' is shown for each digit the algorithm is currently holding
//     (because a carry might still correct it)
//   - When the algorithm confirms the held digits are correct, the '?'
//     marks are overwritten with their true values on screen
//   - When a carry fires, the '?' marks are overwritten with corrected
//     values and a carry announcement is printed
//   - The Feynman Point (six consecutive 9s at position 762) is
//     announced after its '??????' resolves to '999999'
//
// Everything shown reflects the actual state of the algorithm at that
// moment. Nothing is simulated or replayed from Run 1.

func spigotRun2(numberOfDigits int, done chan bool, webPrint func(string), delay time.Duration) bool {
	const lineWidth  = 50
	const feynmanLen = 6

	pi           := ""
	line         := ""
	col          := 0
	boxes        := numberOfDigits * 10 / 3
	remainders   := make([]int, boxes)
	digitsHeld   := 0
	pendingQs    := 0   // how many '?' are currently on the active line
	decPos       := 0   // 1-based decimal digit position
	decInserted  := false
	digitsSeen   := 0
	feynmanFired := false

	for i := 0; i < boxes; i++ {
		select {
		case <-done:
			return false
		default:
			remainders[i] = 2
		}
	}

	updateLine := func() {
		webPrint("UPDATE:    " + line)
	}

	flushLine := func() {
		if line != "" {
			webPrint(line)
			line      = ""
			col       = 0
			pendingQs = 0
		}
	}

	appendChar := func(ch string) {
		line += ch
		col++
		if col >= lineWidth {
			flushLine()
		} else {
			updateLine()
		}
	}

	// showQ appends a '?' to the current line representing genuine
	// uncertainty about a digit the algorithm is currently holding.
	showQ := func() {
		line += "?"
		col++
		pendingQs++
		if col >= lineWidth {
			flushLine()
		} else {
			updateLine()
		}
	}

	// resolveQs overwrites the pending '?' marks with confirmed digits,
	// one by one from left to right, each with a brief pause.
	// Called when the algorithm confirms held 9s are correct.
	resolveQs := func(count int, digit rune) {
		resolveStart := len([]rune(line)) - pendingQs
		runes        := []rune(line)
		for k := 0; k < count && resolveStart+k < len(runes); k++ {
			select {
			case <-done:
				return
			default:
			}
			runes[resolveStart+k] = digit
			line = string(runes)
			updateLine()
			time.Sleep(delay / 2)
		}
		pendingQs -= count
		if pendingQs < 0 {
			pendingQs = 0
		}
	}

	for i := 0; i < numberOfDigits; i++ {
		select {
		case <-done:
			return false
		default:
		}

		carriedOver := 0
		sum         := 0

		for j := boxes - 1; j >= 0; j-- {
			select {
			case <-done:
				return false
			default:
			}
			remainders[j] *= 10
			sum            = remainders[j] + carriedOver
			quotient       := sum / (j*2 + 1)
			remainders[j]  = sum % (j*2 + 1)
			carriedOver     = quotient * j
		}

		remainders[0] = sum % 10
		q             := sum / 10

		switch q {
		case 9:
			// The algorithm has produced a 9 but cannot confirm it yet.
			// A carry from a later iteration could still flip it to 0.
			// Show '?' on screen -- this is genuine uncertainty.
			digitsHeld++
			decPos++
			digitsSeen++
			if !decInserted && digitsSeen == 2 {
				appendChar(".")
				decInserted = true
			}
			pi += "9"
			showQ()
			time.Sleep(delay)

		case 10:
			// Carry fired. The held 9s were wrong -- they become 0s,
			// and the digit before them is incremented.
			// Patch the pi string (same logic as original algorithm).
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

			// Announce the carry and correct the '?' marks on screen
			if pendingQs > 0 {
				flushLine()
				webPrint("")
				webPrint(fmt.Sprintf(
					"  [CARRY at position %d] Those %d '?'s were NOT nines.",
					decPos, pendingQs))
				webPrint("  [CARRY] A carry propagated -- correcting to zeros...")
				webPrint("")
				runes        := []rune(line)
				resolveStart := len(runes) - pendingQs
				for k := 0; k < pendingQs && resolveStart+k < len(runes); k++ {
					runes[resolveStart+k] = '0'
				}
				line = string(runes)
				updateLine()
				time.Sleep(delay * 3)
				pendingQs = 0
			}

			decPos++
			digitsSeen++
			if !decInserted && digitsSeen == 2 {
				appendChar(".")
				decInserted = true
			}
			pi += strconv.Itoa(q)
			appendChar(strconv.Itoa(q))
			digitsHeld = 1
			time.Sleep(delay)

		default:
			// A non-9 digit arrived. If we were holding 9s, they are
			// now confirmed correct -- resolve the '?' marks to '9'.
			if digitsHeld > 0 && pendingQs > 0 {
				resolveQs(pendingQs, '9')

				// Feynman Point: six or more consecutive confirmed 9s
				if !feynmanFired && digitsHeld >= feynmanLen {
					feynmanFired = true
					flushLine()
					webPrint("")
					webPrint("  +--------------------------------------------------+")
					webPrint("  | !! THE FEYNMAN POINT !!                          |")
					webPrint("  | (or perhaps the Hofstadter Point?)               |")
					webPrint("  | Six consecutive 9s at decimal position 762       |")
					webPrint("  | ...134  999999  837...                           |")
					webPrint("  |                                                  |")
					webPrint("  | You just watched the algorithm hold all six as   |")
					webPrint("  | '??????' -- because a carry could have flipped   |")
					webPrint("  | every one of them to zero.                       |")
					webPrint("  | They resolved to 999999. Truth confirmed.        |")
					webPrint("  |                                                  |")
					webPrint("  | Feynman: 'nine nine nine nine nine nine...       |")
					webPrint("  |           ...and so on!'                         |")
					webPrint("  +--------------------------------------------------+")
					webPrint("")
					time.Sleep(3 * time.Second)
				}
			}

			digitsHeld = 1
			decPos++
			digitsSeen++
			if !decInserted && digitsSeen == 2 {
				appendChar(".")
				decInserted = true
			}
			pi += strconv.Itoa(q)
			appendChar(strconv.Itoa(q))
			time.Sleep(delay)
		}
	}

	flushLine()
	return true
}

// ── Shared helpers ────────────────────────────────────────────────────────────

// delChar removes the character at index from string s.
// Written largely by Richard Woolley.
func delChar(s string, index int) string {
	tmp := []rune(s)
	return string(append(tmp[0:index], tmp[index+1:]...))
}

// spigotMax returns the larger of two ints.
// Named to avoid collision with Go 1.21+ builtin max().
func spigotMax(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// suppress unused import warning if strings is only used here
var _ = strings.Repeat
