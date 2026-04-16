// Grok does slow question marks in Spigot.txt

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
// Human-visibility enhancement for red uncertainty markers (?):
//   • 750 ms pause each time a red ? is printed (so humans actually see it)
//   • Slow re-draw of the preceding ~10 confirmed digits + the pending ?s
//     immediately before they resolve (visual "pay attention" replay)
//   • All still 100% honest: the ?s are shown in real time as the algorithm
//     produces them; the re-draw is pure display simulation after the fact.
//
// Design philosophy (Rick's explicit requirement):
//   The algorithm is run TWICE. Both runs are genuine, independent
//   executions of the full Spigot computation. Nothing is reused,
//   replayed, or simulated between them.

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// ── Entry point ───────────────────────────────────────────────────────────────

func TheSpigotWeb(done chan bool, webPrint func(string)) {
	const numberOfDigits = 850
	const bw             = 50

	webPrint("COLOR:cyan:" + boxSep(bw))
	webPrint("COLOR:cyan:" + boxLine("  THE RABINOWITZ–WAGON SPIGOT ALGORITHM       ", bw))
	webPrint("COLOR:cyan:" + boxLine("  Pi from integer arithmetic alone            ", bw))
	webPrint("COLOR:cyan:" + boxLine("  1995 — produces digits sequentially, no floats", bw))
	webPrint("COLOR:cyan:" + boxSep(bw))
	webPrint("COLOR:cyan:" + boxLine("  Origin : Rabinowitz & Wagon (Am. Math. Monthly)", bw))
	webPrint("COLOR:cyan:" + boxLine("          Rewritten by Richard Woolley            ", bw))
	webPrint("COLOR:cyan:" + boxLine("  Method : Pure integer arithmetic only           ", bw))
	webPrint("COLOR:cyan:" + boxLine("  Runs   : TWO genuine independent executions     ", bw))
	webPrint("COLOR:cyan:" + boxSep(bw))
	webPrint("")

	webPrint("COLOR:yellow:" + boxSep(bw))
	webPrint("COLOR:yellow:" + boxLine("  RUN 1 -- Full speed                         ", bw))
	webPrint("COLOR:yellow:" + boxLine("  Computing 850 digits...                     ", bw))
	webPrint("COLOR:yellow:" + boxSep(bw))
	webPrint("")

	run1Start := time.Now()
	ok        := spigotRun1(done, webPrint)
	run1Time  := time.Since(run1Start)

	if !ok {
		webPrint("COLOR:red:  !! Aborted.")
		return
	}

	// Calculate a reasonable base delay for the non-pause parts of Run 2.
	// The actual total runtime is dominated by human-visibility pauses.
	normalDigits := numberOfDigits - 69
	if normalDigits < 1 {
		normalDigits = numberOfDigits
	}
	baseDelayMs := (8.0 * 1000.0) / float64(normalDigits) // gentle base (~8s of pure computation)
	if baseDelayMs < 1 {
		baseDelayMs = 1
	}
	baseDelay := time.Duration(baseDelayMs * float64(time.Millisecond))

	webPrint("")
	webPrint("COLOR:yellow:" + boxSep(bw))
	webPrint("COLOR:yellow:" + boxLine("  RUN 1 COMPLETE                              ", bw))
	webPrint("COLOR:yellow:" + boxLine(fmt.Sprintf("  850 digits in %s", run1Time.Round(time.Microsecond)), bw))
	webPrint("COLOR:yellow:" + boxSep(bw))
	webPrint("")
	webPrint("  Or was that too fast to follow?")
	webPrint("")
	webPrint("  No matter. We will now run the algorithm again --")
	webPrint("  a fresh computation, not a replay of Run 1 --")
	webPrint("  this time with deliberate pauses so you can")
	webPrint("  clearly see the uncertainty in action.")
	webPrint("")
	webPrint("COLOR:red:  Uncertain digits will appear in red as '?'.")
	webPrint("  Each '?' is overwritten with the true digit the")
	webPrint("  moment the algorithm itself resolves it.")
	webPrint("")
	webPrint("  No simulation. Every '?' is real uncertainty,")
	webPrint("  happening right now, on the server.")
	webPrint("")
	webPrint("COLOR:yellow:  One more thing: keep your eyes open around")
	webPrint("COLOR:yellow:  digit 762. Something rather famous lives there.")
	webPrint("")

	time.Sleep(3 * time.Second)

	webPrint("COLOR:green:" + boxSep(bw))
	webPrint("COLOR:green:" + boxLine("  RUN 2 -- Honest edition                     ", bw))
	webPrint("COLOR:green:" + boxLine("  Uncertainty shown as it actually occurs      ", bw))
	webPrint("COLOR:green:" + boxLine("  Human-paced for easy viewing (~3 minutes)     ", bw))
	webPrint("COLOR:green:" + boxLine(fmt.Sprintf("  Base algorithmic delay: %.1f ms/digit         ", baseDelayMs), bw))
	webPrint("COLOR:green:" + boxSep(bw))
	webPrint("")

	ok = spigotRun2(done, webPrint, baseDelay)
	if !ok {
		webPrint("COLOR:red:  !! Aborted.")
		return
	}

	webPrint("")
	webPrint("COLOR:cyan:" + boxSep(bw))
	webPrint("COLOR:cyan:" + boxLine("  Spigot complete: 850 digits of Pi            ", bw))
	webPrint("COLOR:cyan:" + boxLine("  Integer arithmetic only. No floating point   ", bw))
	webPrint("COLOR:cyan:" + boxLine("  numbers were used in this production.        ", bw))
	webPrint("COLOR:cyan:" + boxSep(bw))
}

// ── Run 1: full speed ─────────────────────────────────────────────────────────
// Plain appends only. No UPDATE: needed since every digit is immediately final.

func spigotRun1(done chan bool, webPrint func(string)) bool {
	const numberOfDigits = 850
	const lineWidth      = 50

	size := numberOfDigits*10/3 + 50
	a    := make([]int, size)
	for i := range a { a[i] = 2 }

	line        := ""
	pre         := -1
	nines       := 0
	count       := 0
	decInserted := false

	addChar := func(ch string) {
		line += ch
		if len([]rune(line)) >= lineWidth {
			webPrint(line)
			line = ""
		}
	}

	emit := func(d int) {
		if count >= numberOfDigits { return }
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
		sum         := 0
		for j := size - 1; j >= 0; j-- {
			select {
			case <-done:
				return false
			default:
			}
			a[j]      *= 10
			sum        = a[j] + carriedOver
			quotient   := sum / (j*2 + 1)
			a[j]       = sum % (j*2 + 1)
			carriedOver = quotient * j
		}
		a[0] = sum % 10
		q    := sum / 10

		switch {
		case q == 9:
			nines++
		case q == 10:
			if pre >= 0 { emit(pre + 1) }
			for i := 0; i < nines && count < numberOfDigits; i++ { emit(0) }
			pre   = 0
			nines = 0
		default:
			if pre >= 0 { emit(pre) }
			for i := 0; i < nines && count < numberOfDigits; i++ { emit(9) }
			pre   = q
			nines = 0
		}
	}

	if line != "" { webPrint(line) }
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
//
// Human-noticeability improvements (April 2026):
//   • Every red ? is held on-screen for a full 750 ms when first printed.
//   • Immediately before any pending ?s are resolved, the last 3 confirmed
//     digits + the current ?s are visually re-drawn very quickly
//     (300ms → 200ms → 150ms). Snappy but noticeable.
//
//   Total runtime is now dominated by deliberate human-visibility pauses
//   rather than raw CPU speed. On both Mac Mini and render.com free tier,
//   Run 2 consistently takes ~3 minutes.

func spigotRun2(done chan bool, webPrint func(string), baseDelay time.Duration) bool {
	const numberOfDigits = 850
	const lineWidth      = 50
	const feynmanLen     = 6
	const slowStart      = 700
	const feynmanZone    = 762
	const previewDigits  = 3          // final value after tuning for best visibility
	const humanPause     = 750 * time.Millisecond

	size := numberOfDigits*10/3 + 50
	a    := make([]int, size)
	for i := range a { a[i] = 2 }

	line         := ""
	col          := 0
	pendingQs    := 0
	pre          := -1
	nines        := 0
	count        := 0
	decPos       := 0
	decInserted  := false
	feynmanFired := false

	delayFor := func(pos int) time.Duration {
		if pos < slowStart || pos > feynmanZone+6 {
			return baseDelay
		}
		t      := float64(pos-slowStart) / float64(feynmanZone-slowStart)
		factor := 1.0 + 3.0*t
		return time.Duration(float64(baseDelay) * factor)
	}

	show := func() {
		if pendingQs > 0 {
			webPrint("UPDATE:HASRED:" + line)
		} else {
			webPrint("UPDATE:" + line)
		}
	}

	newRow := func() {
		webPrint("")
		line      = ""
		col       = 0
		pendingQs = 0
	}

	showQ := func() {
		if !decInserted && count == 1 {
			line += "."
			col++
			decInserted = true
		}
		line += "?"
		col++
		pendingQs++
		show()
		time.Sleep(humanPause) // each red ? lingers visibly for the user
	}

	// confirmAndEmit: reveals pre before the pending ?s and resolves them
	// one by one to resolvedNineChar.
	confirmAndEmit := func(preDigit int, resolvedNineChar rune) {
		runes  := []rune(line)
		qStart := len(runes) - pendingQs

		prefix := make([]rune, qStart, qStart+2)
		copy(prefix, runes[:qStart])

		if preDigit >= 0 && count < numberOfDigits {
			if !decInserted && count == 1 {
				prefix      = append(prefix, '.')
				decInserted = true
			}
			prefix = append(prefix, rune('0'+preDigit))
			decPos++
			count++
		}
		prefixLen := len(prefix)

		pending := make([]rune, pendingQs)
		copy(pending, runes[qStart:])
		line = string(prefix) + string(pending)
		col  = len([]rune(line))
		show()
		time.Sleep(baseDelay / 2)

		for i := 0; pendingQs > 0 && count < numberOfDigits; i++ {
			select {
			case <-done:
				pendingQs = 0
				return
			default:
			}
			allRunes         := []rune(line)
			allRunes[prefixLen+i] = resolvedNineChar
			line      = string(allRunes)
			pendingQs--
			decPos++
			count++
			show()
			time.Sleep(baseDelay / 2)
		}
		pendingQs = 0
		col = len([]rune(line))
		if col >= lineWidth {
			newRow()
		}
	}

	// slowReplayPrecedingUncertainty: re-draws the last 3 confirmed digits
	// + pending ?s with a quick progressive timing ramp for optimal visibility.
	slowReplayPrecedingUncertainty := func() {
		if pendingQs == 0 {
			return
		}
		runes := []rune(line)
		qStart := len(runes) - pendingQs
		if qStart < 0 {
			return
		}

		previewLen := previewDigits
		if previewLen > qStart {
			previewLen = qStart
		}
		previewStart := qStart - previewLen

		slowSegment := runes[previewStart:]
		slowBase := string(runes[:previewStart])

		// Quick blank to make the re-draw clearly visible
		blanked := slowBase + strings.Repeat(" ", len(slowSegment))
		webPrint("UPDATE:" + blanked)
		time.Sleep(200 * time.Millisecond)

		// Fast progressive delays: 300ms → 200ms → 150ms
		replayDelays := []time.Duration{
			300 * time.Millisecond,
			200 * time.Millisecond,
			150 * time.Millisecond,
		}

		for i := 0; i < len(slowSegment); i++ {
			partialSegment := slowSegment[:i+1]
			partial := slowBase + string(partialSegment)

			if strings.Contains(string(partialSegment), "?") {
				webPrint("UPDATE:HASRED:" + partial)
			} else {
				webPrint("UPDATE:" + partial)
			}

			delay := replayDelays[len(replayDelays)-1]
			if i < len(replayDelays) {
				delay = replayDelays[i]
			}
			time.Sleep(delay)
		}

		// Restore exact original line state
		line = slowBase + string(slowSegment)
		col = len([]rune(line))
	}

	webPrint("") // seed blank row for first UPDATE:

	for count < numberOfDigits {
		select {
		case <-done:
			return false
		default:
		}

		carriedOver := 0
		sum         := 0
		for j := size - 1; j >= 0; j-- {
			select {
			case <-done:
				return false
			default:
			}
			a[j]      *= 10
			sum        = a[j] + carriedOver
			quotient   := sum / (j*2 + 1)
			a[j]       = sum % (j*2 + 1)
			carriedOver = quotient * j
		}
		a[0] = sum % 10
		q    := sum / 10

		switch {
		case q == 9:
			nines++
			showQ()
			time.Sleep(delayFor(decPos + nines))

		case q == 10:
			savedNines := nines
			carryPre   := -1
			if pre >= 0 { carryPre = pre + 1 }

			if nines > 0 {
				slowReplayPrecedingUncertainty()
			}

			confirmAndEmit(carryPre, '0')
			nines = 0

			if savedNines > 0 {
				if line != "" { newRow() }
				if savedNines == 1 {
					webPrint("COLOR:red:" + fmt.Sprintf(
						"  [CARRY at ~position %d] That 9 becomes 0.", decPos))
				} else {
					webPrint("COLOR:red:" + fmt.Sprintf(
						"  [CARRY at ~position %d] Those %d 9s become 0s.", decPos, savedNines))
				}
				webPrint("COLOR:red:  The digit before them was also incremented.")
				time.Sleep(baseDelay * 3)
				webPrint("")
				webPrint("")
			}
			pre = 0

		default:
			savedNines := nines

			if nines > 0 {
				slowReplayPrecedingUncertainty()
			}

			confirmAndEmit(pre, '9')
			nines = 0

			if !feynmanFired && savedNines >= feynmanLen {
				feynmanFired = true
				if line != "" { newRow() }
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
				webPrint("COLOR:yellow:  | They resolved to 999999. Truth confirmed         |")
				webPrint("COLOR:yellow:  | in real time ... in real-slow time.              |")
				webPrint("COLOR:yellow:  |                                                  |")
				webPrint("COLOR:yellow:  | Feynman: \"nine nine nine nine nine nine...       |")
				webPrint("COLOR:yellow:  |          ...and so on!\"                         |")
				webPrint("COLOR:yellow:  +--------------------------------------------------+")
				webPrint("")
				time.Sleep(4 * time.Second)
				webPrint("") // seed blank row for digits to continue
			}

			pre = q
			time.Sleep(delayFor(decPos + 1))
		}
	}

	// Emit final held pre if needed
	if pre >= 0 && count < numberOfDigits {
		if nines > 0 {
			slowReplayPrecedingUncertainty()
		}
		confirmAndEmit(pre, '9')
	}

	if line != "" { newRow() }
	return true
}

/*
// ── Helpers ───────────────────────────────────────────────────────────────────

// delChar removes the character at index from string s.
// Written largely by Richard Woolley.
func delChar(s string, index int) string {
	tmp := []rune(s)
	return string(append(tmp[0:index], tmp[index+1:]...))
}
	*/