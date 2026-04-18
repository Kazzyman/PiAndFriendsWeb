package main

// roots.go
//
// A brute-force demonstration of how square and cube roots can be
// calculated using nothing but whole-number (integer) arithmetic.
//
// Original concept and algorithm: Richard (Rick) Woolley.
//
// The core insight: if A and B are both perfect powers (squares or cubes),
// and A/B is close to our target number N, then the ratio of their roots
// is a close approximation to the Nth root of N.
//
// Example for square roots:
//   We want √11.
//   We find that 49/4 = 12.25  (just above 11)
//   and that    36/4 = 9.0     (just below 11)
//   So √11 ≈ 7/2 = 3.5  (rough first guess)
//   As we search larger perfect squares we find ever-closer ratios,
//   giving us ever more accurate approximations.
//
// This is the same insight ancient Greek mathematicians used to find
// that √2 ≈ 7/5, because 49/25 = 1.96 ≈ 2.
//
// Rewritten cleanly: April 2026.

import (
	"fmt"
	"math"
	"sort"
	"time"
)

// @formatter:off

// ── Package-level variables ───────────────────────────────────────────────────

var (
	pairsSlice    []Pairs // The lookup table of perfect squares or cubes
	mathSqrtCheat float64 // Used only to print a verification for perfect squares
	mathCbrtCheat float64 // Used only to print a verification for perfect cubes
)

// ── Structs ───────────────────────────────────────────────────────────────────

// Pairs holds one entry in our lookup table.
// 'product' is a perfect square or cube (e.g. 49, 125).
// 'root'    is the whole number whose power produced it (e.g. 7, 5).
type Pairs struct {
	product int
	root    int
}

// Results holds one candidate approximation found during the search.
// We collect many and sort by pdiff to find the best one.
// 'result'        is the calculated root approximation (rootOfLarger / rootOfSmaller).
// 'pdiff'         is the proportional difference -- smaller means better.
// 'largerPP'      is the larger perfect power that bracketed the target on the high side.
// 'smallerPP'     is the anchor (smaller) perfect power.
// 'rootOfLarger'  is the root of largerPP.
// 'rootOfSmaller' is the root of smallerPP (the anchor entry).
type Results struct {
	result        float64
	pdiff         float64
	largerPP      int
	smallerPP     int
	rootOfLarger  int
	rootOfSmaller int
}

// sortedResults accumulates all candidate answers during one run.
var sortedResults []Results

// ── Entry point ───────────────────────────────────────────────────────────────

func xRootOfy(radical2or3 int, workPiece int, webPrint func(string)) {

	if radical2or3 != 2 && radical2or3 != 3 {
		webPrint(fmt.Sprintf("  Error: radical %d is not supported. Please enter 2 or 3.", radical2or3))
		return
	}

	usingBigFloats = false

	// Reset ALL state from any previous run.
	sortedResults = nil
	diffOfLarger  = 1
	diffOfSmaller = 1
	mathSqrtCheat = 0
	mathCbrtCheat = 0

	startTime := time.Now()

	precisionOfRoot = calculatePrecision(radical2or3, workPiece, webPrint)

	webPrint("")
	webPrint("  ── Root Calculation ─────────────────────────────")
	if radical2or3 == 2 {
		webPrint(fmt.Sprintf("  Finding the Square Root of %d", workPiece))
	} else {
		webPrint(fmt.Sprintf("  Finding the Cube Root of %d", workPiece))
	}
	webPrint("  Method: integer arithmetic only, no floating point")
	webPrint(fmt.Sprintf("  Precision window: %d", precisionOfRoot))
	webPrint("  ─────────────────────────────────────────────────")
	webPrint("")

	if radical2or3 == 2 {
		webPrint("  Building table of perfect squares...")
	} else {
		webPrint("  Building table of perfect cubes...")
	}
	buildPairsSlice(radical2or3)
	if radical2or3 == 2 {
		webPrint(fmt.Sprintf("  Table built: %d perfect squares.", len(pairsSlice)))
	} else {
		webPrint(fmt.Sprintf("  Table built: %d perfect cubes.", len(pairsSlice)))
	}

	// Show the first few table entries so the user can see what we built.
	webPrint("")
	webPrint("  First few entries in the table:")
	for k := 0; k < 5 && k < len(pairsSlice); k++ {
		webPrint(fmt.Sprintf("    root=%d  →  %d^%d = %d",
			pairsSlice[k].root, pairsSlice[k].root, radical2or3, pairsSlice[k].product))
	}
	webPrint("    ...")
	webPrint("")

	webPrint("  The method: for each perfect power in the table,")
	webPrint("  search forward for another whose ratio brackets")
	webPrint(fmt.Sprintf("  our target number %d.", workPiece))
	webPrint("  When largerPP / smallerPP ≈ workPiece,")
	webPrint("  then rootOfLarger / rootOfSmaller ≈ answer.")
	webPrint("")

	// Show a concrete example using the first two table entries.
	if len(pairsSlice) >= 2 {
		a := pairsSlice[0]
		b := pairsSlice[1]
		webPrint("  Example with first two entries:")
		webPrint(fmt.Sprintf("    entry[0]: root=%d  product=%d", a.root, a.product))
		webPrint(fmt.Sprintf("    entry[1]: root=%d  product=%d", b.root, b.product))
		webPrint(fmt.Sprintf("    ratio = %d / %d = %.4f  (target: %d)",
			b.product, a.product,
			float64(b.product)/float64(a.product),
			workPiece))
		webPrint(fmt.Sprintf("    root ratio = %d / %d = %.6f",
			b.root, a.root,
			float64(b.root)/float64(a.root)))
		webPrint("    That ratio is our first rough approximation.")
		webPrint("    We keep searching for ratios ever closer to the target.")
	}
	webPrint("")
	webPrint("  Starting search...")
	webPrint("")

	midpointShown    := false
	midpoint         := len(pairsSlice) / 2
	threeQuarterMark := (len(pairsSlice) * 3) / 4

	runSearch := func() {
		sortedResults  = nil
		diffOfLarger   = 1
		diffOfSmaller  = 1
		firstHitShown := false
		midpointShown  = false

		for i := 0; i < len(pairsSlice)-1; i += 2 {

			prevCount := len(sortedResults)

			readPairsSlice(i, workPiece)

			if diffOfLarger == 0 || diffOfSmaller == 0 {
				return
			}

			// Show the first candidate in detail.
			if !firstHitShown && len(sortedResults) > prevCount {
				firstHitShown = true
				best := sortedResults[len(sortedResults)-1]
				if radical2or3 == 2 {
					verification := math.Sqrt(float64(workPiece))
					webPrint(fmt.Sprintf("  First candidate found at table index %d:", i))
					webPrint(fmt.Sprintf("    smaller PP : %d  (root %d)",
						pairsSlice[i].product, pairsSlice[i].root))
					webPrint(fmt.Sprintf("    larger PP  : product near %d",
						pairsSlice[i].product*workPiece))
					webPrint(fmt.Sprintf("    root ratio = %0.9f", best.result))
					webPrint(fmt.Sprintf("    math.Sqrt  = %0.9f", verification))
					webPrint(fmt.Sprintf("    difference = %0.9f", math.Abs(best.result-verification)))
				} else {
					verification := math.Cbrt(float64(workPiece))
					webPrint(fmt.Sprintf("  First candidate found at table index %d:", i))
					webPrint(fmt.Sprintf("    smaller PP : %d  (root %d)",
						pairsSlice[i].product, pairsSlice[i].root))
					webPrint(fmt.Sprintf("    root ratio = %0.9f", best.result))
					webPrint(fmt.Sprintf("    math.Cbrt  = %0.9f", verification))
					webPrint(fmt.Sprintf("    difference = %0.9f", math.Abs(best.result-verification)))
				}
				webPrint("")
				webPrint("  Continuing search for a better approximation...")
				webPrint("")
			}

			// Progress every 80,000 iterations.
			if i%80000 == 0 && i > 0 {
				stringI := formatInt64WithThousandSeparators(int64(i))
				webPrint(fmt.Sprintf("  %s iterations completed...  elapsed: %s",
					stringI, time.Since(startTime).Round(time.Millisecond)))
				if len(sortedResults) > 0 {
					webPrint(fmt.Sprintf("  Best so far: %0.9f  (candidates found: %d)",
						sortedResults[len(sortedResults)-1].result, len(sortedResults)))
				}
			}

			// Midpoint summary.
			if !midpointShown && i >= midpoint {
				midpointShown = true
				mid := pairsSlice[i]
				webPrint("")
				webPrint(fmt.Sprintf("  ── Midpoint check (index %d of %d) ──",
					i, len(pairsSlice)))
				webPrint(fmt.Sprintf("  We are now testing ratios involving root=%d", mid.root))
				webPrint(fmt.Sprintf("  whose perfect power is %d", mid.product))
				if len(sortedResults) > 0 {
					best := sortedResults[len(sortedResults)-1]
					webPrint(fmt.Sprintf("  Best approximation so far: %0.9f", best.result))
					if radical2or3 == 2 {
						webPrint(fmt.Sprintf("  True value (math.Sqrt):    %0.9f",
							math.Sqrt(float64(workPiece))))
					} else {
						webPrint(fmt.Sprintf("  True value (math.Cbrt):    %0.9f",
							math.Cbrt(float64(workPiece))))
					}
				} else {
					webPrint("  No candidates found yet -- precision window may need widening.")
				}
				webPrint("")
			}

			// Three-quarter mark summary.
			if midpointShown && i == threeQuarterMark {
				tq := pairsSlice[i]
				webPrint("")
				webPrint(fmt.Sprintf("  ── Three-quarter mark (index %d) ──", i))
				webPrint(fmt.Sprintf("  Now testing root=%d  (perfect power: %d)",
					tq.root, tq.product))
				if len(sortedResults) > 0 {
					webPrint(fmt.Sprintf("  Best so far: %0.9f  candidates: %d",
						sortedResults[len(sortedResults)-1].result, len(sortedResults)))
				} else {
					webPrint("  Still no candidates -- will widen precision if needed.")
				}
				webPrint("")
			}

			// Gap-based early exit.
			if i+2 < len(pairsSlice) {
				gap := pairsSlice[i+2].product - pairsSlice[i].product
				if gap > precisionOfRoot && len(sortedResults) > 0 {
					webPrint("")
					webPrint(fmt.Sprintf("  Gap between consecutive perfect powers (%d)", gap))
					webPrint(fmt.Sprintf("  now exceeds precision window (%d).", precisionOfRoot))
					webPrint("  No better result is mathematically possible.")
					webPrint("  Searching 2000 more iterations just to be sure...")
					for extra := i + 2; extra < i+2002 && extra < len(pairsSlice)-1; extra += 2 {
						readPairsSlice(extra, workPiece)
					}
					webPrint("  Done. Stopping search.")
					return
				}
			}
		}
	}

	// First attempt.
	runSearch()

	// Auto-widen and retry if nothing found.
	for len(sortedResults) == 0 && diffOfLarger != 0 && diffOfSmaller != 0 {
		if precisionOfRoot >= 500000 {
			webPrint("")
			webPrint("  Could not find a result even at maximum precision window.")
			webPrint("  This workPiece may be too large for the current table size.")
			return
		}
		precisionOfRoot *= 2
		webPrint("")
		webPrint(fmt.Sprintf("  No results found. Widening precision window to %d and retrying...  (elapsed: %s)",
			precisionOfRoot, time.Since(startTime).Round(time.Millisecond)))
		webPrint("")
		runSearch()
	}

	// Perfect result.
	if diffOfLarger == 0 || diffOfSmaller == 0 {
		elapsed := time.Since(startTime)
		webPrint("")
		webPrint("  ── Perfect Result ───────────────────────────────")
		if radical2or3 == 2 {
			webPrint(fmt.Sprintf("  %d is a perfect square.", workPiece))
			webPrint(fmt.Sprintf("  Its square root is exactly %0.0f", mathSqrtCheat))
		} else {
			webPrint(fmt.Sprintf("  %d is a perfect cube.", workPiece))
			webPrint(fmt.Sprintf("  Its cube root is exactly %0.0f", mathCbrtCheat))
		}
		webPrint(fmt.Sprintf("  Completed in: %s", elapsed.Round(time.Millisecond)))
		webPrint("  ─────────────────────────────────────────────────")
		return
	}

	// Final result.
	elapsed := time.Since(startTime)

	sort.Slice(sortedResults, func(i, j int) bool {
		return sortedResults[i].pdiff < sortedResults[j].pdiff
	})

	b := sortedResults[0] // best candidate

	webPrint("  ── Result ───────────────────────────────────────")
	if radical2or3 == 2 {
		verification := math.Sqrt(float64(workPiece))
		webPrint(fmt.Sprintf("  Square Root of %d", workPiece))
		webPrint(fmt.Sprintf("  Our result    : %0.9f", b.result))
		webPrint(fmt.Sprintf("  Verification  : %0.9f  (math.Sqrt)", verification))
		webPrint(fmt.Sprintf("  Difference    : %0.9f", math.Abs(b.result-verification)))

		// Check if result is suspiciously close to a whole number.
		frac := b.result - math.Floor(b.result)
		if frac < 0.01 || frac > 0.99 {
			webPrint("")
			webPrint("  Note: result is very close to a whole number.")
			webPrint(fmt.Sprintf("  %d sits just below %d² = %d,",
				workPiece,
				int(math.Round(b.result)),
				int(math.Round(b.result))*int(math.Round(b.result))))
			webPrint("  which is the best rational approximation available.")
			webPrint("  This is a known characteristic of the method for")
			webPrint("  numbers that sit just below a perfect square.")
		}
	} else {
		verification := math.Cbrt(float64(workPiece))
		webPrint(fmt.Sprintf("  Cube Root of %d", workPiece))
		webPrint(fmt.Sprintf("  Our result    : %0.9f", b.result))
		webPrint(fmt.Sprintf("  Verification  : %0.9f  (math.Cbrt)", verification))
		webPrint(fmt.Sprintf("  Difference    : %0.9f", math.Abs(b.result-verification)))
		if workPiece == 2 {
			webPrint("")
			webPrint("  * - * The Delian Problem * - *")
			webPrint("  The cube root of 2 was sought by ancient Greek geometers")
			webPrint("  as the solution to 'doubling the cube' -- constructing a")
			webPrint("  cube with exactly twice the volume of a given cube.")
			webPrint("  It tormented mathematicians for over two thousand years.")
			webPrint("  They could not solve it with compass and straightedge alone.")
			webPrint(" ")
			webPrint("  ✦ We just solved it in milliseconds using an algorithm ✦")
			webPrint(" ")
			webPrint("  This algorithm was inspired by the same geometric reasoning")
			webPrint("  that was pioneered by the ancients")
			webPrint(" ")
			webPrint("  ✦✦ Archimedes would have wept ✦✦")
			webPrint(" ")
			webPrint("  The cube root of two was finally proved impossible with a")
			webPrint("  compass and a straightedge in 1837 by Pierre Wantzel —")
			webPrint("  over 2,200 years after the Greeks first posed it.")
			webPrint(" ")
		}
	}

	// ── Show the actual perfect power pair that produced the best result ──────
	//
	// This is the educational heart of the output -- the user can verify the
	// arithmetic themselves: divide the two perfect powers, take the ratio of
	// their roots, and average high and low to get an even better approximation.
	webPrint("")
	webPrint("  ── The perfect power pair behind this result ─────")
	webPrint(fmt.Sprintf("  Anchor (smaller) PP : %d  (root %d)",
		b.smallerPP, b.rootOfSmaller))
	webPrint(fmt.Sprintf("  Bracket PP          : %d  (root %d)",
		b.largerPP, b.rootOfLarger))
	webPrint(fmt.Sprintf("  PP ratio            : %d / %d = %.6f  (target: %d)",
		b.largerPP, b.smallerPP,
		float64(b.largerPP)/float64(b.smallerPP),
		workPiece))
	webPrint("")

	// The high and low root ratios that bracket the true answer.
	highRatio := float64(b.rootOfLarger) / float64(b.rootOfSmaller)
	lowRatio  := float64(b.rootOfLarger-1) / float64(b.rootOfSmaller)
	average   := (highRatio + lowRatio) / 2.0

	webPrint(fmt.Sprintf("  Root ratio (high)   : %d / %d = %.9f",
		b.rootOfLarger, b.rootOfSmaller, highRatio))
	webPrint(fmt.Sprintf("  Root ratio (low)    : %d / %d = %.9f",
		b.rootOfLarger-1, b.rootOfSmaller, lowRatio))
	webPrint(fmt.Sprintf("  Average of high+low : %.9f", average))
	if radical2or3 == 2 {
		webPrint(fmt.Sprintf("  True value          : %.9f  (math.Sqrt)", math.Sqrt(float64(workPiece))))
	} else {
		webPrint(fmt.Sprintf("  True value          : %.9f  (math.Cbrt)", math.Cbrt(float64(workPiece))))
	}
	webPrint("")
	webPrint("  You can verify this yourself: divide the two perfect")
	if radical2or3 == 2 {
		webPrint("  squares above, then take the ratio of their roots.")
	} else {
		webPrint("  cubes above, then take the ratio of their roots.")
	}
	webPrint("  Averaging the high and low ratios gives an even")
	webPrint("  better approximation than either alone.")

	webPrint(fmt.Sprintf("  Candidates found: %d (best pdiff: %0.8f)",
		len(sortedResults), sortedResults[0].pdiff))
	webPrint(fmt.Sprintf("  Completed in    : %s", elapsed.Round(time.Millisecond)))
	webPrint("  ─────────────────────────────────────────────────")
	webPrint("")
	webPrint("  (Verification values from math.Sqrt/Cbrt are shown")
	webPrint("   only to confirm accuracy. The calculation above")
	webPrint("   used integer arithmetic exclusively.)")
}


// ── readPairsSlice: the heart of the search ───────────────────────────────────
//
// Starting from table entry i (the "anchor" perfect power), this function
// searches forward through the table to find the first entry j where:
//
//   pairsSlice[j].product > pairsSlice[i].product * workPiece
//
// At that crossing point, entries j and j-1 bracket the target.
// The ratio of their roots gives our approximation.
//
// This function makes NO webPrint calls -- completely silent.
// All output is handled by xRootOfy.
//
func readPairsSlice(i int, workPiece int) {

	oneReadOfSmallerRoot := pairsSlice[i].root
	oneReadOfSmallerPP   := pairsSlice[i].product

	for j := i + 1; j < len(pairsSlice); j++ {

		// Guard against integer overflow before multiplying.
		if oneReadOfSmallerPP > math.MaxInt64/workPiece {
			return
		}

		largerPP := pairsSlice[j].product
		target   := oneReadOfSmallerPP * workPiece

		if largerPP > target {

			rootOfLarger  := pairsSlice[j].root
			smallerPP     := pairsSlice[j-1].product
			rootOfSmaller := pairsSlice[j-1].root

			dL := largerPP - target
			dS := target - smallerPP

			if dL < 0 || dS < 0 {
				return
			}

			if dL == 0 {
				diffOfLarger  = 0
				mathSqrtCheat = math.Sqrt(float64(workPiece))
				mathCbrtCheat = math.Cbrt(float64(workPiece))
				return
			}
			if dS == 0 {
				diffOfSmaller = 0
				mathSqrtCheat = math.Sqrt(float64(workPiece))
				mathCbrtCheat = math.Cbrt(float64(workPiece))
				return
			}

			diffOfLarger  = dL
			diffOfSmaller = dS

			// Record candidates within the precision window,
			// now storing the full pair information for educational output.
			if dL < precisionOfRoot {
				sortedResults = append(sortedResults, Results{
					result:        float64(rootOfLarger) / float64(oneReadOfSmallerRoot),
					pdiff:         float64(dL) / float64(largerPP),
					largerPP:      largerPP,
					smallerPP:     oneReadOfSmallerPP,
					rootOfLarger:  rootOfLarger,
					rootOfSmaller: oneReadOfSmallerRoot,
				})
			}
			if dS < precisionOfRoot {
				sortedResults = append(sortedResults, Results{
					result:        float64(rootOfSmaller) / float64(oneReadOfSmallerRoot),
					pdiff:         float64(dS) / float64(smallerPP),
					largerPP:      smallerPP,
					smallerPP:     oneReadOfSmallerPP,
					rootOfLarger:  rootOfSmaller,
					rootOfSmaller: oneReadOfSmallerRoot,
				})
			}

			return
		}
	}
}

// ── calculatePrecision ────────────────────────────────────────────────────────

func calculatePrecision(radical2or3, workPiece int, webPrint func(string)) int {
	var precision int
	if radical2or3 == 2 {
		precision = int(math.Sqrt(float64(workPiece))) / 3
		if precision < 4   { precision = 4   }
		if precision > 500 { precision = 500 }
	} else {
		precision = workPiece * 12
		if precision < 600  { precision = 600  }
		if precision > 3000 { precision = 3000 }
	}
	webPrint(fmt.Sprintf("  Precision window set to %d", precision))
	return precision
}

// ── buildPairsSlice ───────────────────────────────────────────────────────────

func buildPairsSlice(radical2or3 int) {
	pairsSlice = nil
	root := 2
	for i := 0; i < 825000; i++ {
		root++
		var product int
		if radical2or3 == 2 {
			product = root * root
		} else {
			product = root * root * root
		}
		pairsSlice = append(pairsSlice, Pairs{
			product: product,
			root:    root,
		})
	}
}
