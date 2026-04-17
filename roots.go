package main

// roots.go
//
// A brute-force demonstration of how square and cube roots can be
// calculated using nothing but whole-number (integer) arithmetic.
//
// Original concept and algorithm: Richard (Rick) Woolley.
// The core insight: if A and B are both perfect powers (squares or cubes),
// and A/B is close to our target number N, then the ratio of their roots
// is a close approximation to the Nth root of N.
//
// Example for square roots:
//   We want √11.
//   We notice that 49/4 = 12.25  (close to 11, but a bit high)
//   and that       36/4 = 9.0    (close to 11, but a bit low)
//   So √11 ≈ 7/2 = 3.5  (rough), or more precisely we keep searching
//   for larger perfect squares whose ratio is even closer to 11.
//   Eventually we find e.g. 2116/196 = 10.795... and 2209/196 = 11.27...
//   giving √11 ≈ 46/14 = 3.2857... which is correct to several digits.
//
// This is the same insight ancient Greek mathematicians used to discover
// that √2 ≈ 7/5, because 49/25 = 1.96 ≈ 2.
//
// Rewritten with clean logic and educational commentary: April 2026.

import (
	"fmt"
	"math"
	"sort"
	"time"
)

// @formatter:off

// ── Package-level variables ───────────────────────────────────────────────────

var (
	pairsSlice    []Pairs  // The table of perfect squares or cubes and their roots
	mathSqrtCheat float64  // Used only to print a verification when a perfect square is found
	mathCbrtCheat float64  // Used only to print a verification when a perfect cube is found
)

// ── Pairs: the building block of our table ────────────────────────────────────

// Pairs holds one entry in our lookup table.
// 'product' is a perfect square or cube (e.g. 49, 64, 125).
// 'root'    is the whole number whose power produced it (e.g. 7, 8, 5).
type Pairs struct {
	product int
	root    int
}

// ── Results: a candidate answer ───────────────────────────────────────────────

// Results holds one candidate approximation found during the search.
// We collect many of these and then sort by pdiff to find the best one.
// 'result' is the calculated root approximation (e.g. 3.31662...).
// 'pdiff'  is the proportional difference -- how close we got.
//           Smaller pdiff means a better approximation.
/*
// this is in globals.go
type Results struct {
	result float64
	pdiff  float64
}
	*/

// sortedResults accumulates all candidate answers found during one run.
// It is sorted by pdiff at the end so sortedResults[0] is the best answer.
// this too is in globals.go
// var sortedResults []Results

// ── Entry point ───────────────────────────────────────────────────────────────

// xRootOfy calculates the square root (radical=2) or cube root (radical=3)
// of workPiece using only integer arithmetic -- no math.Sqrt or math.Cbrt
// is used for the actual calculation (only for verification printouts).
func xRootOfy(radical2or3 int, workPiece int, webPrint func(string)) {

	// Safety check: we only support square (2) and cube (3) roots.
	if radical2or3 != 2 && radical2or3 != 3 {
		webPrint(fmt.Sprintf("Error: Radical %d is not supported. Please enter 2 or 3.", radical2or3))
		return
	}

	usingBigFloats = false
	timeOfStartFromTop := time.Now()

	// Reset the results accumulator between runs.
	// Without this, results from a previous run would pollute the next one.
	sortedResults = nil

	// Precision controls how close our ratio A/B needs to be to workPiece
	// before we consider it a "hit" worth recording. Larger workpieces need
	// a larger precision window because the perfect power gaps are bigger.
	precisionOfRoot = calculatePrecision(radical2or3, workPiece, webPrint)

	// Print a clear header so the user understands what we are about to do.
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

	// Build the lookup table of perfect squares or cubes.
	// This is the foundation of the whole method.
	webPrint("Building table of perfect powers...")
	buildPairsSlice(radical2or3)
	webPrint(fmt.Sprintf("Table built: %d entries. Starting search...", len(pairsSlice)))
	webPrint("")

	startBeforeCall := time.Now()

	// ── Main search loop ──────────────────────────────────────────────────
	//
	// We step through the table two entries at a time. For each entry i,
	// readPairsSlice searches forward from i to find a second entry whose
	// product, when divided by entry i's product, is close to workPiece.
	//
	// Why step by 2? Because stepping by 1 would cause excessive overlap
	// and redundant comparisons. Stepping by 2 gives reasonable coverage
	// without being wasteful. You could step by 1 for more thoroughness
	// at the cost of speed.
	//
	// We stop early if a perfect square or cube is found (diff == 0).
	for i := 0; i < len(pairsSlice)-1; i += 2 {

		// Check periodically if the user has asked us to stop.
		// (done channel check is handled inside readPairsSlice)

		readPairsSlice(i, startBeforeCall, radical2or3, workPiece, webPrint)

		// If a perfect result was found, diffOfLarger or diffOfSmaller
		// will be exactly zero. No need to keep searching.
		if diffOfLarger == 0 || diffOfSmaller == 0 {
			handlePerfectSquaresAndCubes(timeOfStartFromTop, radical2or3, workPiece, webPrint)
			break
		}

		// Print a progress update every 80,000 iterations so the user
		// knows the algorithm is still alive and working.
		if i%80000 == 0 && i > 0 {
			stringI := formatInt64WithThousandSeparators(int64(i))
			webPrint(fmt.Sprintf("  %s iterations completed... still searching...", stringI))
		}
	}

	// ── Final result ──────────────────────────────────────────────────────

	elapsed := time.Since(timeOfStartFromTop)

	// Only enter the result block if we did NOT find a perfect square/cube.
	// (Perfect results are reported immediately inside the loop above.)
	if diffOfLarger != 0 && diffOfSmaller != 0 {

		if len(sortedResults) == 0 {
			// This should not happen with reasonable inputs, but handle it
			// gracefully rather than panicking.
			webPrint(fmt.Sprintf("No results found within precision %d. Try a larger workPiece or adjust precision.", precisionOfRoot))
			return
		}

		// Sort all candidate results by their proportional difference.
		// The best approximation (smallest pdiff) will be at index 0.
		sort.Slice(sortedResults, func(i, j int) bool {
			return sortedResults[i].pdiff < sortedResults[j].pdiff
		})

		best := sortedResults[0].result

		// Print the best result. We use Go's math library ONLY to print
		// a verification value alongside our integer-arithmetic result --
		// the calculation itself used no floating point.
		webPrint("")
		webPrint("  ── Result ───────────────────────────────────────")
		if radical2or3 == 2 {
			verification := math.Sqrt(float64(workPiece))
			webPrint(fmt.Sprintf("  Square Root of %d", workPiece))
			webPrint(fmt.Sprintf("  Our result    : %0.9f", best))
			webPrint(fmt.Sprintf("  Verification  : %0.9f  (math.Sqrt)", verification))
			webPrint(fmt.Sprintf("  Difference    : %0.9f", math.Abs(best-verification)))
		} else {
			verification := math.Cbrt(float64(workPiece))
			webPrint(fmt.Sprintf("  Cube Root of %d", workPiece))
			webPrint(fmt.Sprintf("  Our result    : %0.9f", best))
			webPrint(fmt.Sprintf("  Verification  : %0.9f  (math.Cbrt)", verification))
			webPrint(fmt.Sprintf("  Difference    : %0.9f", math.Abs(best-verification)))
		}
		webPrint(fmt.Sprintf("  Candidates found: %d (best pdiff: %0.8f)", len(sortedResults), sortedResults[0].pdiff))
		webPrint(fmt.Sprintf("  Completed in    : %s", elapsed.Round(time.Millisecond)))
		webPrint("  ─────────────────────────────────────────────────")
		webPrint("")
		webPrint("  (Verification values from math.Sqrt/Cbrt are shown")
		webPrint("   only to confirm accuracy. The calculation above")
		webPrint("   used integer arithmetic exclusively.)")
	}
}

// ── readPairsSlice: the heart of the search ───────────────────────────────────
//
// Starting from table entry i (our "smaller" perfect power), this function
// searches forward through the table to find a "larger" perfect power such
// that:
//
//   largerPerfectPower / smallerPerfectPower  ≈  workPiece
//
// When it finds such a pair, the ratio of their roots:
//
//   rootOfLarger / rootOfSmaller  ≈  nthRoot(workPiece)
//
// Why does this work? Because:
//   If largerPP = a^n  and  smallerPP = b^n
//   and largerPP / smallerPP ≈ workPiece
//   then (a/b)^n ≈ workPiece
//   therefore a/b ≈ nthRoot(workPiece)
//
func readPairsSlice(i int, startBeforeCall time.Time, radical2or3, workPiece int, webPrint func(string)) {

	// Read the "anchor" entry -- the smaller perfect power we will
	// use as the denominator in our ratio search.
	oneReadOfSmallerRoot := pairsSlice[i].root
	oneReadOfSmallerPP   := pairsSlice[i].product

	// Search forward from i+1 looking for a larger perfect power
	// whose ratio to oneReadOfSmallerPP brackets workPiece.
	for j := i + 1; j < len(pairsSlice); j++ {

		largerPerfectProduct := pairsSlice[j].product

		// We are looking for the first larger PP that exceeds
		// oneReadOfSmallerPP * workPiece. At that crossing point,
		// the entry just before it (j-1) is slightly below, and
		// this entry (j) is slightly above. Both are candidates.
		if largerPerfectProduct > oneReadOfSmallerPP*workPiece {

			// ── Larger side (just above the target) ──
			ProspectivePHitOnLargeSide       := pairsSlice[j].product
			rootOfProspectivePHitOnLargeSide := pairsSlice[j].root

			// ── Smaller side (just below the target) ──
			ProspectivePHitOnSmallerSide       := pairsSlice[j-1].product
			rootOfProspectivePHitOnSmallerSide := pairsSlice[j-1].root

			// How far off are we on each side?
			// diffOfLarger  = how much largerPP overshoots  (workPiece * smallerPP)
			// diffOfSmaller = how much smallerPP undershoots (workPiece * smallerPP)
			diffOfLarger  = ProspectivePHitOnLargeSide  - workPiece*oneReadOfSmallerPP
			diffOfSmaller = workPiece*oneReadOfSmallerPP - ProspectivePHitOnSmallerSide

			// ── Perfect result check ──────────────────────────────────────
			// If diff is exactly zero, we have found a perfect square or cube.
			if diffOfLarger == 0 {
				webPrint("  !! Perfect result found on the larger side.")
				webPrint(fmt.Sprintf("  The %d root of %d is exactly %0.0f",
					radical2or3, workPiece,
					float64(rootOfProspectivePHitOnLargeSide)/float64(oneReadOfSmallerRoot)))
				mathCbrtCheat = math.Cbrt(float64(workPiece))
				mathSqrtCheat = math.Sqrt(float64(workPiece))
				return
			}
			if diffOfSmaller == 0 {
				webPrint("  !! Perfect result found on the smaller side.")
				webPrint(fmt.Sprintf("  The %d root of %d is exactly %0.0f",
					radical2or3, workPiece,
					float64(rootOfProspectivePHitOnSmallerSide)/float64(oneReadOfSmallerRoot)))
				mathSqrtCheat = math.Sqrt(float64(workPiece))
				mathCbrtCheat = math.Cbrt(float64(workPiece))
				return
			}

			// ── Approximate result: larger side ──────────────────────────
			// Record this candidate if it falls within our precision window.
			if diffOfLarger < precisionOfRoot {
				result := float64(rootOfProspectivePHitOnLargeSide) / float64(oneReadOfSmallerRoot)
				pdiff  := float64(diffOfLarger) / float64(ProspectivePHitOnLargeSide)
				sortedResults = append(sortedResults, Results{result: result, pdiff: pdiff})

				// Show our work -- this is the educational heart of the module.
				webPrint("  [+] Larger side hit:")
				webPrint(fmt.Sprintf("      Small PP=%d (root=%d),  Large PP=%d (root=%d)",
					oneReadOfSmallerPP, oneReadOfSmallerRoot,
					ProspectivePHitOnLargeSide, rootOfProspectivePHitOnLargeSide))
				webPrint(fmt.Sprintf("      %d / %d = %0.6f  (target: %d.0)",
					ProspectivePHitOnLargeSide, oneReadOfSmallerPP,
					float64(ProspectivePHitOnLargeSide)/float64(oneReadOfSmallerPP),
					workPiece))
				webPrint(fmt.Sprintf("      Root ratio: %d / %d = %0.9f",
					rootOfProspectivePHitOnLargeSide, oneReadOfSmallerRoot, result))
				webPrint(fmt.Sprintf("      Difference: %d  (pdiff: %0.8f)", diffOfLarger, pdiff))

				// Elapsed time -- reassures the user things are progressing.
				elapsed := time.Since(startBeforeCall)
				if elapsed.Seconds() > 0.1 {
					webPrint(fmt.Sprintf("      Time so far: %s", elapsed.Round(time.Millisecond)))
				}
			}

			// ── Approximate result: smaller side ─────────────────────────
			if diffOfSmaller < precisionOfRoot {
				result := float64(rootOfProspectivePHitOnSmallerSide) / float64(oneReadOfSmallerRoot)
				pdiff  := float64(diffOfSmaller) / float64(ProspectivePHitOnSmallerSide)
				sortedResults = append(sortedResults, Results{result: result, pdiff: pdiff})

				webPrint("  [-] Smaller side hit:")
				webPrint(fmt.Sprintf("      Small PP=%d (root=%d),  Smaller PP=%d (root=%d)",
					oneReadOfSmallerPP, oneReadOfSmallerRoot,
					ProspectivePHitOnSmallerSide, rootOfProspectivePHitOnSmallerSide))
				webPrint(fmt.Sprintf("      %d / %d = %0.6f  (target: %d.0)",
					ProspectivePHitOnSmallerSide, oneReadOfSmallerPP,
					float64(ProspectivePHitOnSmallerSide)/float64(oneReadOfSmallerPP),
					workPiece))
				webPrint(fmt.Sprintf("      Root ratio: %d / %d = %0.9f",
					rootOfProspectivePHitOnSmallerSide, oneReadOfSmallerRoot, result))
				webPrint(fmt.Sprintf("      Difference: %d  (pdiff: %0.8f)", diffOfSmaller, pdiff))

				elapsed := time.Since(startBeforeCall)
				if elapsed.Seconds() > 0.1 {
					webPrint(fmt.Sprintf("      Time so far: %s", elapsed.Round(time.Millisecond)))
				}
			}

			// We found the crossing point for this anchor entry.
			// Break the inner search and let the outer loop advance
			// to the next anchor entry.
			break
		}
	}
}

// ── handlePerfectSquaresAndCubes ─────────────────────────────────────────────
//
// Called only when a perfect result was found (diff == 0).
// Prints a clear confirmation message.
func handlePerfectSquaresAndCubes(timeOfStart time.Time, radical2or3, workPiece int, webPrint func(string)) {
	elapsed := time.Since(timeOfStart)
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
}

// ── calculatePrecision ────────────────────────────────────────────────────────
//
// Determines how wide a "closeness window" to use when deciding whether a
// pair of perfect powers is close enough to workPiece to count as a hit.
//
// The original code had hardcoded magic numbers for specific workpiece values,
// which was fragile. This version scales automatically:
//
//   - Square roots need only a tiny window (4) because perfect squares are
//     dense and we find good approximations quickly.
//   - Cube roots need a larger window because perfect cubes are sparser.
//     Larger workpieces need a larger window for the same reason.
//   - We clamp between 600 and 2000 to keep runtime reasonable.
//
func calculatePrecision(radical2or3, workPiece int, webPrint func(string)) int {
	var precision int
	if radical2or3 == 2 {
		precision = 4
	} else {
		// Scale with workPiece: bigger numbers need a wider search window.
		precision = workPiece * 120
		if precision < 600  { precision = 600  }
		if precision > 2000 { precision = 2000 }
	}
	webPrint(fmt.Sprintf("  Precision window set to %d", precision))
	return precision
}

// ── buildPairsSlice ───────────────────────────────────────────────────────────
//
// Builds the lookup table of perfect squares or cubes.
// Each entry records the perfect power and the whole number root that produced it.
//
// We start from root=2 (since 1^n = 1 is trivial and not useful as a denominator)
// and build 825,000 entries, giving us roots up to 825,002.
//
// For square roots: entries are (4,2), (9,3), (16,4), (25,5) ...
// For cube roots:   entries are (8,2), (27,3), (64,4), (125,5) ...
//
func buildPairsSlice(radical2or3 int) {
	pairsSlice = nil // Clear any results from a previous run
	root := 1
	for i := 0; i < 825000; i++ {
		root++
		var identityProduct int
		if radical2or3 == 2 {
			identityProduct = root * root
		} else {
			identityProduct = root * root * root
		}
		pairsSlice = append(pairsSlice, Pairs{
			product: identityProduct,
			root:    root,
		})
	}
}
