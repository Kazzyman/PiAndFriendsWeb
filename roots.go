package main

import (
	"fmt"
	"math"
	"os"
	"sort"
	"time"
)

// @formatter:off

// TODO share with Claude. Ask, is this cool or what. 

var (
	pairsSlice []Pairs
	mathSqrtCheat            float64
	mathCbrtCheat            float64
)

func xRootOfy(radical2or3 int, workPiece int, webPrint func(string)) {
	// This says: "I don't know HOW to print to the web, but I expect to be passed a function (webPrint) that does."

	// 1. Safety Check: Only allow Square (2) or Cube (3) roots
	if radical2or3 != 2 && radical2or3 != 3 {
		webPrint(fmt.Sprintf("Error: Radical %d is not supported. Please enter 2 or 3.", radical2or3))
		return // Gracefully quit the function
	}
	
	usingBigFloats = false
	TimeOfStartFromTop := time.Now()

	// radical2or3 := mgr.GetRadical() // How do we now get this? Who will be calling xRootOfy()
	// workPiece := mgr.GetWorkPiece()

	radical2or3, workPiece = setPrecisionForSquareOrCubeRoot(radical2or3, workPiece, webPrint) // sets precision only, no actual need to digest and pass our inputs
	// mgr.SetRadical(radical2or3) // no need for these
	// mgr.SetWorkPiece(workPiece)

	webPrint("Building table...")
	buildPairsSlice(radical2or3)
	webPrint("Table built, starting calculation...")
	startBeforeCall := time.Now()

	var indx int
	for i := 0; i < 400000; i += 2 { // this is meant to be a pretty big loop 825,000 is the number of 
		/*
		if mgr.ShouldStop() {
			webPrint("Calculation of a root aborted")
			webPrint(fmt.Sprintf("Calculation of a root aborted"))
			return
		}
		// Need to implement the alternative to this. 
		 */
		readPairsSlice(i, startBeforeCall, radical2or3, workPiece, webPrint)
		handlePerfectSquaresAndCubes(TimeOfStartFromTop, radical2or3, workPiece, webPrint)
		if diffOfLarger == 0 || diffOfSmaller == 0 {
			break // because we have a perfect square or cube
		}
		if i%80000 == 0 && i > 0 { // if remainder of div is 0 (every 80,000 iterations) conditional progress updates print
			stringVindx := formatInt64WithThousandSeparators(int64(indx))
			webPrint(fmt.Sprintf("%s iterations completed... of 400,000", stringVindx))
			webPrint("... still working ...") // ok

			webPrint(fmt.Sprintf("%s iterations completed... of 400,000", stringVindx))
			fmt.Println(i, "... still working ...")
		}
		indx = i // save/copy to a wider scope for later use outside this loop
	}
	fmt.Println("Loop completed at index:", indx) // ::: debug to console.

	// ::: Show the final result
	webPrint(fmt.Sprintf("Entering result block, mathSqrtCheat 'square': %v, mathCbrtCheat 'cube': %v", mathSqrtCheat, mathCbrtCheat)) 
	// ::: "Entering result block ... "

	t_s2 := time.Now()
	elapsed_s2 := t_s2.Sub(TimeOfStartFromTop)
	if diffOfLarger != 0 || diffOfSmaller != 0 { // if not a perfect square or cube do this else skip due to detection of perfect result
		defer func() {
			if r := recover(); r != nil {
				fmt.Println("Panic in result block:", r)
				webPrint("Error calculating result")
			}
		}()
		fileHandle, err31 := os.OpenFile("dataLog-From_calculate-pi-and-friends.txt", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
		check(err31)
		defer fileHandle.Close()

		Hostname, _ := os.Hostname()
		fmt.Fprintf(fileHandle, "  -- %d root of %d by a ratio of perfect Products -- on %s ", radical2or3, workPiece, Hostname)
		fmt.Fprint(fileHandle, "was run on: ", time.Now().Format(time.ANSIC), "")
		fmt.Fprintf(fileHandle, "%d was total Iterations ", indx)

		fmt.Println("Sorting results...") // ::: debug to console 
		sort.Slice(sortedResults, func(i, j int) bool { return sortedResults[i].pdiff < sortedResults[j].pdiff })
		fmt.Println("Sorted results, length:", len(sortedResults)) // ::: debug to console

		if len(sortedResults) > 0 {
			if radical2or3 == 2 {
				result := fmt.Sprintf("%0.9f, it's the best approximation for the Square Root of %d", sortedResults[0].result, workPiece)
				fmt.Println("Updating GUI with:", result) // ::: debug to console
				webPrint(result)
				fmt.Println("GUI updated, printing to console...") // ::: debug to console
				webPrint(fmt.Sprintf("%s", result))
				// webPrint(fmt.Sprintf("Square-root Result is: %s", result))
				fmt.Println("Writing to file...") // ::: debug to console
				fmt.Fprintf(fileHandle, "%s ", result)
				fmt.Println("File written") // ::: debug to console
			}
			if radical2or3 == 3 {
				result := fmt.Sprintf("%0.9f, it's the best approximation for the Cube Root of %d", sortedResults[0].result, workPiece) //  ::: todo, This seems to run/print twice?
				fmt.Println("Updating GUI with:", result) // ::: debug to console
				webPrint(result)
				fmt.Println("GUI updated, printing to console...") // ::: debug to console
				webPrint(fmt.Sprintf("%s", result))
				// webPrint(fmt.Sprintf("Cube-root Result is: %s", result))
				fmt.Println("Writing to file...") // ::: debug to console
				fmt.Fprintf(fileHandle, "%s ", result)
				fmt.Println("File written") // ::: debug to console
			}
		} else {
			webPrint(fmt.Sprintf("No results found within precision %d after %d iterations", precisionOfRoot, indx))
			fmt.Printf("No results found within precision %d after %d iterations", precisionOfRoot, indx)
		}

		// TotalRun := elapsed_s2.String()
		// fmt.Fprintf(fileHandle, "Total run was %s  ", TotalRun)
		webPrint(fmt.Sprintf("Calculation completed in %s", elapsed_s2))
	} else {
		fmt.Println("Skipped result block due to perfect result detection") // ::: debug to console
	}
}

func readPairsSlice(i int, startBeforeCall time.Time, radical2or3, workPiece int, webPrint func(string)) {
	// The "webPrint func(string)" part says: "I don't know HOW to print to the web, but I expect to be passed a function (webPrint) that does."
	
	oneReadOfSmallerRoot := pairsSlice[i].root // Read a smaller PP and its root (just once) for each time readPairsSlice is called
	oneReadOfSmallerPP := pairsSlice[i].product

	for iter := 0; iter < 410000 && i < len(pairsSlice); iter++ { // go big, but not so big that you would read past the end of the pairsSlice
		i++
		largerPerfectProduct := pairsSlice[i].product // i has been incremented since the initial one-time read of oneReadOfSmallerPP

		// ... and, keep incrementing the i until largerPerfectProduct is greater than (oneReadOfSmallerPP * workPiece)
		if largerPerfectProduct > oneReadOfSmallerPP*workPiece { // For example: workPiece may be 11, 3.32*3.32.   Larger PP may be 49, 7*7.   Smaller oneReadPP may be 4, 2*2. ::: oneRead is 4

			ProspectivePHitOnLargeSide := largerPerfectProduct // rename it, badly;
			rootOfProspectivePHitOnLargeSide := pairsSlice[i].root // grab larger side's root

			ProspectivePHitOnSmallerSide := pairsSlice[i-1].product
			rootOfProspectivePHitOnSmallerSide := pairsSlice[i-1].root


			// we next look at two roots (PHs) of two PPs. 
			diffOfLarger = ProspectivePHitOnLargeSide - workPiece*oneReadOfSmallerPP // ::: PH_larger - (WP * _once)     7 - (11 * 4)
			// What does it tell us if we find that the sum of one of the larger roots from the table : ProspectivePHitOnLargeSide
			// and/plus the negative of another smaller root from the table (times our WP) turns out to be zero?


			diffOfSmaller = workPiece*oneReadOfSmallerPP - ProspectivePHitOnSmallerSide // ::: (WP * _once) - PH_smaller    (11 * 4) - 

			if diffOfLarger == 0 {
				fmt.Println(colorCyan, " The", radical2or3, "root of", workPiece, "is", colorGreen,
					float64(rootOfProspectivePHitOnLargeSide)/float64(oneReadOfSmallerRoot), colorReset, "")
				webPrint(fmt.Sprintf(" The %d root of %d is %0.33f", radical2or3, workPiece, float64(rootOfProspectivePHitOnLargeSide)/float64(oneReadOfSmallerRoot)))

				mathCbrtCheat = math.Cbrt(float64(workPiece))
				break // TODO is this supposed to be here, see the remaining ifs? 
			}
			if diffOfSmaller == 0 {
				fmt.Println(colorCyan, " The", radical2or3, "root of", workPiece, "is", colorGreen,
					float64(rootOfProspectivePHitOnSmallerSide)/float64(oneReadOfSmallerRoot), colorReset, "")
				webPrint(fmt.Sprintf(" The %d root of %d is %0.33f", radical2or3, workPiece, float64(rootOfProspectivePHitOnSmallerSide)/float64(oneReadOfSmallerRoot)))

				mathSqrtCheat = math.Sqrt(float64(workPiece)) // ::: I cheated? Yea, a bit. But only in order to generate verbiage to print re a perfect root having been found 
				mathCbrtCheat = math.Cbrt(float64(workPiece))
				break // TODO is this supposed to be here, see the remaining ifs? 
			}

			if diffOfLarger < precisionOfRoot {
				result := float64(rootOfProspectivePHitOnLargeSide) / float64(oneReadOfSmallerRoot)
				pdiff := float64(diffOfLarger) / float64(ProspectivePHitOnLargeSide)

				sortedResults = append(sortedResults, Results{result: result, pdiff: pdiff})

				webPrint(fmt.Sprintf("Found large prospect at index %d: result=%f, diff=%d", i, result, diffOfLarger)) 
				// break // was apparenly here, but is now on the next line? // TODO ?
				if diffOfLarger < 2 {break}
			}
			if diffOfSmaller < precisionOfRoot {
				result := float64(rootOfProspectivePHitOnSmallerSide) / float64(oneReadOfSmallerRoot)
				pdiff := float64(diffOfSmaller) / float64(ProspectivePHitOnSmallerSide)

				sortedResults = append(sortedResults, Results{result: result, pdiff: pdiff})

				webPrint(fmt.Sprintf("Found small prospect at index %d: result=%f, diff=%d", i, result, diffOfSmaller)) 
				webPrint(fmt.Sprintf("Found small prospect at index %d: result=%f, diff=%d", i, result, diffOfSmaller)) 
				// break // was apparenly here, but is now on the next line? // TODO ?
				if diffOfSmaller < 2 {break}
			}

			// ::: we will be potentially duplicating Results struct -> slice 
			// larger side section: ----------------------------------------------------------------------------------------------------------------------------------------
			// ---------------------------------------------------------------------------------------------------------------------------------------------------------------
			if diffOfLarger < precisionOfRoot { // report the prospects, their differences, and the calculated result for the Sqrt or Cbrt
				fmt.Println("small PP is", colorCyan, oneReadOfSmallerPP, colorReset, "and, slightly on the higher side of", workPiece,
					"* that we found a PP of", colorCyan, ProspectivePHitOnLargeSide, colorReset, "a difference of", diffOfLarger) // ::: debug to console
				webPrint(fmt.Sprintf("small PP is %d and, slightly on the higher side of %d * that we found a PP of %d a difference of %d", oneReadOfSmallerPP, workPiece, ProspectivePHitOnLargeSide, diffOfLarger))

				fmt.Println("the ", radical2or3, " root of ", workPiece, " is calculated as ", colorGreen,
					float64(rootOfProspectivePHitOnLargeSide)/float64(oneReadOfSmallerRoot), colorReset)
				webPrint(fmt.Sprintf("the %d root of %d is calculated as %0.6f ", radical2or3, workPiece, float64(rootOfProspectivePHitOnLargeSide)/float64(oneReadOfSmallerRoot)))

				webPrint(fmt.Sprintf("with pdiff of %0.4f ", (float64(diffOfLarger)/float64(ProspectivePHitOnLargeSide))*100000))
				webPrint(fmt.Sprintf("with pdiff of %0.4f ", (float64(diffOfLarger)/float64(ProspectivePHitOnLargeSide))*100000))
				// save the result to an accumulator array so we can Fprint all such hits at the very end
				// List_of_2_results_case18 = append(List_of_2_results_case18, float64(rootOfProspectivePHitOnLargeSide) / float64(oneReadOfSmallerRoot) )
				// corresponding_diffs = append(corresponding_diffs, diffOfLarger)
				// diffs_as_percent = append(diffs_as_percent, float64(diffOfLarger)/float64(ProspectivePHitOnLargeSide))

				// in the next five lines we load (append) a record into/to the file (array) of Results
				Result1 := Results{
					result: float64(rootOfProspectivePHitOnLargeSide) / float64(oneReadOfSmallerRoot),
					pdiff:  float64(diffOfLarger) / float64(ProspectivePHitOnLargeSide),
				}
				sortedResults = append(sortedResults, Result1)

				t2 := time.Now()
				elapsed2 := t2.Sub(startBeforeCall)
				// if needed, notify the user that we are still working
				Tim_win = 0.178
				if radical2or3 == 3 {
					if workPiece > 13 {
						Tim_win = 0.0012
					} else {
						Tim_win = 0.003
					}
				}
				if elapsed2.Seconds() > Tim_win {
					fmt.Println(elapsed2.Seconds(), "Seconds have elapsed ... working ...")
					webPrint(fmt.Sprintf("%0.4f Seconds have elapsed ... working ...", elapsed2.Seconds()))
				}
				// TODO no break here?
			}

			// smaller side section: ----------------------------------------------------------------------------------------------------------------------------------------
			// ---------------------------------------------------------------------------------------------------------------------------------------------------------------
			if diffOfSmaller < precisionOfRoot { // report the prospects, their differences, and the calculated result for the Sqrt or Cbrt
				fmt.Println("small PP is", colorCyan, oneReadOfSmallerPP, colorReset, "and, slightly on the lesser side of", workPiece,
					"* that we found a PP of", colorCyan, ProspectivePHitOnSmallerSide, colorReset, "a difference of", diffOfSmaller)
				webPrint(fmt.Sprintf("small PP is %d and, slightly on the higher side of %d * that we found a PP of %d a difference of %d", oneReadOfSmallerPP, workPiece, ProspectivePHitOnSmallerSide, diffOfSmaller))

				fmt.Println("the ", radical2or3, " root of ", workPiece, " is calculated as ", colorGreen,
					float64(rootOfProspectivePHitOnSmallerSide)/float64(oneReadOfSmallerRoot), colorReset)
				webPrint(fmt.Sprintf("the %d root of %d is calculated as %0.6f ", radical2or3, workPiece, float64(rootOfProspectivePHitOnSmallerSide)/float64(oneReadOfSmallerRoot)))

				webPrint(fmt.Sprintf("with pdiff of %0.4f ", (float64(diffOfSmaller)/float64(ProspectivePHitOnSmallerSide))*100000))
				webPrint(fmt.Sprintf("with pdiff of %0.4f ", (float64(diffOfSmaller)/float64(ProspectivePHitOnSmallerSide))*100000))

				// save the result to three accumulator arrays so we can Fprint all such hits, diffs, and p-diffs, at the very end of run
				// List_of_2_results_case18 = append(List_of_2_results_case18, float64(rootOfProspectivePHitOnSmallerSide) / float64(oneReadOfSmallerRoot) )
				// corresponding_diffs = append(corresponding_diffs, diffOfSmaller)
				// diffs_as_percent = append(diffs_as_percent, float64(diffOfSmaller)/float64(ProspectivePHitOnSmallerSide))
				// ***** ^^^^ ****** the preceeding was replaced with the following five lines *******************************************

				// in the next five lines we load (append) a record into/to the file (array) of Results
				Result1 := Results{
					result: float64(rootOfProspectivePHitOnSmallerSide) / float64(oneReadOfSmallerRoot),
					pdiff:  float64(diffOfSmaller) / float64(ProspectivePHitOnSmallerSide),
				}
				sortedResults = append(sortedResults, Result1)

				t2 := time.Now()
				elapsed2 := t2.Sub(startBeforeCall)
				// if needed, notify the user that we are still working
				Tim_win = 0.178
				if radical2or3 == 3 {
					if workPiece > 13 {
						Tim_win = 0.0012
					} else {
						Tim_win = 0.003
					}
				}
				if elapsed2.Seconds() > Tim_win {
					fmt.Println(elapsed2.Seconds(), "Seconds have elapsed ... working ...")
					webPrint(fmt.Sprintf("%0.4f Seconds have elapsed ... working ...", elapsed2.Seconds()))
				}
			}
			break // TODO resonsider this, and the preceding if
		}
	}
}

// handlePerfectSquaresAndCubes reports perfect squares/cubes to file and UI
func handlePerfectSquaresAndCubes(TimeOfStartFromTop time.Time, radical2or3, workPiece int, webPrint func(string)) {
	if diffOfLarger == 0 || diffOfSmaller == 0 {
		// t_s1 := time.Now()
		// elapsed_s1 := t_s1.Sub(TimeOfStartFromTop)

		if radical2or3 == 2 {
			result := fmt.Sprintf("Perfect square: %v is the %d root of %d", mathSqrtCheat, radical2or3, workPiece)
			webPrint(result)
		}
		if radical2or3 == 3 {
			result := fmt.Sprintf("Perfect cube: %0.0f is the %d root of %d", mathCbrtCheat, radical2or3, workPiece)
			webPrint(result)
		}
	}
}


// setPrecisionForSquareOrCubeRoot adjusts precision based on radical and workPiece
func setPrecisionForSquareOrCubeRoot(radical2or3, workPiece int, webPrint func(string)) (int, int) {
	if radical2or3 == 3 { // ::: setting the optimal precision this way is a crude kluge
		if workPiece > 4 {
			precisionOfRoot = 1700
			fmt.Println(" Default precision is 1700 ")
			webPrint(" Default precision is 1700 ")
		}
		if workPiece == 2 || workPiece == 11 || workPiece == 17 {
			precisionOfRoot = 600
			fmt.Println(" resetting precision to 600 ")
			webPrint(" resetting precision to 600 ")
		}
		if workPiece == 3 || workPiece == 4 || workPiece == 14 {
			precisionOfRoot = 900
			fmt.Println(" resetting precision to 900 ")
			webPrint(" resetting precision to 900 ")
		}
	}
	if radical2or3 == 2 {
		precisionOfRoot = 4
	}
	return radical2or3, workPiece
}

// Pairs A struct to contain two related whole numbers: an identity product (perfect square or cube), e.g. 49; and its root, which in that case would be 7 
type Pairs struct {
	product int
	root int
}

// build a table of ::: perfect squares or cubes
func buildPairsSlice(radical2or3 int) {
// func buildPairsSlice(radical2or3 int) { // ::: - -
	var identityProduct int
	pairsSlice = nil // Clear/reset the slice between runs
	root := 2 // Because; 2 is the smallest possible whole-number root, i.e., it's the square root of 4 and the cube root of 8 // I used to have this as root := 10 but I do not recall why : (how I had decided on 10?)
	for i := 0; i < 825000; i++ {
		root++
		if radical2or3 == 3 {                   // ::: depending on passed radical 
			identityProduct = root * root * root
		}
		if radical2or3 == 2 {
			identityProduct = root * root
		}
		pairsSlice = append(pairsSlice, Pairs{
			product: identityProduct,
			root:  root,
		})
	}
}