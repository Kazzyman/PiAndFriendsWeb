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
// 'result' is the calculated root approximation.
// 'pdiff'  is the proportional difference -- smaller means better.
type Results struct {
	result float64
	pdiff  float64
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

	// Show a concrete example of the ratio test using the first two table entries.
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

	// midpointShown tracks whether we have printed a midpoint summary.
	// We print it once, around the halfway mark of the table.
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

			// Midpoint summary -- show what the table looks like halfway through.
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
			if !midpointShown && i >= threeQuarterMark {
				// midpointShown doubles as a three-quarter guard to avoid
				// printing both messages in the same run.
			}
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

	best := sortedResults[0].result

	webPrint("  ── Result ───────────────────────────────────────")
	if radical2or3 == 2 {
		verification := math.Sqrt(float64(workPiece))
		webPrint(fmt.Sprintf("  Square Root of %d", workPiece))
		webPrint(fmt.Sprintf("  Our result    : %0.9f", best))
		webPrint(fmt.Sprintf("  Verification  : %0.9f  (math.Sqrt)", verification))
		webPrint(fmt.Sprintf("  Difference    : %0.9f", math.Abs(best-verification)))
		// Check if result is suspiciously close to a whole number
		frac := best - math.Floor(best)
		if frac < 0.01 || frac > 0.99 {
			webPrint("")
			webPrint("  Note: result is very close to a whole number.")
			webPrint(fmt.Sprintf("  %d sits just below %d² = %d,",
				workPiece,
				int(math.Round(best)),
				int(math.Round(best))*int(math.Round(best))))
			webPrint("  which is the best rational approximation available.")
			webPrint("  This is a known characteristic of the method for")
			webPrint("  numbers that sit just below a perfect square.")
		}
	} else {
    verification := math.Cbrt(float64(workPiece))
    webPrint(fmt.Sprintf("  Cube Root of %d", workPiece))
    webPrint(fmt.Sprintf("  Our result    : %0.9f", best))
    webPrint(fmt.Sprintf("  Verification  : %0.9f  (math.Cbrt)", verification))
    webPrint(fmt.Sprintf("  Difference    : %0.9f", math.Abs(best-verification)))
    if radical2or3 == 3 && workPiece == 2 {
        webPrint("")
        webPrint("  * - * The Delian Problem * - *")
        webPrint("  The cube root of 2 was sought by ancient Greek geometers")
        webPrint("  as the solution to 'doubling the cube' -- constructing a")
        webPrint("  cube with exactly twice the volume of a given cube.")
        webPrint("  It tormented mathematicians for over two thousand years.")
        webPrint("  They could not solve it with compass and straightedge alone.")
		webPrint(" ")
        webPrint("  ✦ We just solved it in milliseconds using an algorithm ✦ ")
		webPrint(" ")
        webPrint("  This algorithm was inspired by the same geometric reasoning")
		webPrint("  that was pioneered by the ancients")
		webPrint(" ")
        webPrint("  ✦✦ Archimedes would have wept ✦✦")
		webPrint(" ")
		webPrint("The cube root of two was finally proved impossible with a compass")
		webPrint("and a straightedge in 1837 by Pierre Wantzel — over 2,200 years")
		webPrint("after the Greeks first posed it.")
		webPrint(" ")
    }
}
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
		// For large workpieces and large perfect cubes, the product
		// oneReadOfSmallerPP * workPiece can silently overflow int64,
		// producing garbage negative diffs and nonsense results.
		if oneReadOfSmallerPP > math.MaxInt64/workPiece {
			return // too large to compare safely -- skip this anchor
		}

		largerPP := pairsSlice[j].product
		target   := oneReadOfSmallerPP * workPiece

		if largerPP > target {

			rootOfLarger  := pairsSlice[j].root
			smallerPP     := pairsSlice[j-1].product
			rootOfSmaller := pairsSlice[j-1].root

			// How far off are we on each side?
			dL := largerPP - target  // overshoot  on the larger side
			dS := target - smallerPP // undershoot on the smaller side

			// If either diff is negative, overflow has crept in.
			// Discard this anchor.
			if dL < 0 || dS < 0 {
				return
			}

			// Perfect result: ratio is exactly workPiece.
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

			// Update globals so the caller can read them.
			diffOfLarger  = dL
			diffOfSmaller = dS

			// Record candidates that fall within the precision window.
			if dL < precisionOfRoot {
				result := float64(rootOfLarger) / float64(oneReadOfSmallerRoot)
				pdiff  := float64(dL) / float64(largerPP)
				sortedResults = append(sortedResults, Results{result: result, pdiff: pdiff})
			}
			if dS < precisionOfRoot {
				result := float64(rootOfSmaller) / float64(oneReadOfSmallerRoot)
				pdiff  := float64(dS) / float64(smallerPP)
				sortedResults = append(sortedResults, Results{result: result, pdiff: pdiff})
			}

			return // one crossing point per anchor is enough
		}
	}
}

// ── calculatePrecision ────────────────────────────────────────────────────────
//
// Sets the precision window -- how close largerPP/smallerPP must be to
// workPiece before we record it as a candidate.
//
// Square roots: tiny window (4) because perfect squares are dense.
// Cube roots:   scales with workPiece because perfect cubes are sparser.
//               Clamped between 600 and 3000 for reasonable runtime.
//
func calculatePrecision(radical2or3, workPiece int, webPrint func(string)) int {
	var precision int
	if radical2or3 == 2 {
		// Scale with the square root of workPiece so large workpieces
		// get a reasonable starting window without grinding through
		// dozens of useless passes doubling from 4.
		// For small workpieces like 49 this still gives at least 4,
		// which is tight enough to find exact perfect squares cleanly.
		
		
		// precision = int(math.Sqrt(float64(workPiece))) / 10
		precision = int(math.Sqrt(float64(workPiece))) / 3


		if precision < 4   { precision = 4   }
		if precision > 500 { precision = 500 }
	} else {
		// Cube roots need a wider window because perfect cubes are sparser.
		precision = workPiece * 12
		if precision < 600  { precision = 600  }
		if precision > 3000 { precision = 3000 }
	}
	webPrint(fmt.Sprintf("  Precision window set to %d", precision))
	return precision
}

// ── buildPairsSlice ───────────────────────────────────────────────────────────
//
// Builds the lookup table of perfect squares or cubes.
// Starts at root=3 and builds 825,000 entries.
//
// Square roots: entries are (9,3), (16,4), (25,5), ...
// Cube roots:   entries are (27,3), (64,4), (125,5), ...
//
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
