package main

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

func runRootsWeb(params url.Values, webPrint func(string)) {
	/* ^ ^ this function is called like so: ... what is that referred to as? A closure?
	runRootsWeb(r.URL.Query(), func(s string) {
					outputChan <- s
				})
	// In any case what is ultimately obtained via params via Get is a string.
	*/
	radicalEntryAsString := params.Get("radical")
	workPieceAsString := params.Get("workpiece")

	radicalEntryAsInt, _ := strconv.Atoi(radicalEntryAsString)
	workPieceAsInt, _ := strconv.Atoi(workPieceAsString)

	// Bridging to my original roots.go variables
	Radical_index = radicalEntryAsInt

	webPrint(fmt.Sprintf("Starting Roots Demo: Radical %d, Workpiece %d\n", radicalEntryAsInt, workPieceAsInt))

	// This calls my original build logic from roots.go
	buildPairsSlice(radicalEntryAsInt, webPrint)

	// Note: You can add the specific root-finding call here 
	// Gemini: based on how you normally trigger it in window2.go
	// SetupRootsDemo(mgr, radicalEntryAsString, workEntry, updateOutput2) // That is how it was called in window2.go 
	SetupRootsDemo(radicalEntryAsString, workPieceAsString, webPrint) // Properly named variables are passed. 
	// webPrint("Calculation complete, in bridge.go.") // This seemed to be duplicated in roots.go
	xRootOfy(radicalEntryAsInt, workPieceAsInt, webPrint)
}

/*
// said to need to look like this:
func runRootsWeb(params url.Values, webPrint func(string)) {
    // 1. Get the strings from the URL
    radicalEntryAsString := params.Get("radical")
    workPieceAsString := params.Get("workpiece")

    // 2. Convert to Integers
    radical, _ := strconv.Atoi(radicalEntryAsString)
    workPiece, _ := strconv.Atoi(workPieceAsString)

    // 3. Set the globals so the rest of roots.go can see them
    Radical_index = radical // This satisfies your global in globals.go

    // 4. Call the math function directly, passing our webPrint baton
    // We pass the actual integers we just made.
    xRootOfyWeb(radical, workPiece, webPrint)
}
*/

func SetupRootsDemo(radicalEntryAsString, workPieceAsString string, webPrint func(string)) {
	trimmedRadicalString := strings.TrimRight(radicalEntryAsString, " ")
	radical, err := strconv.Atoi(trimmedRadicalString)
	if err != nil || (radical != 2 && radical != 3) {
		webPrint("Invalid radical: enter 2 or 3\n")
		return
	}
	trimmedWorkPieceString := strings.TrimRight(workPieceAsString, " ")
	workPiece, err := strconv.Atoi(trimmedWorkPieceString)
	if err != nil || workPiece < 0 {
		webPrint("Invalid number: enter a non-negative integer\n")
		webPrint(fmt.Sprintf("Invalid number: enter a non-negative integer\n")) // This does not assign the string to any var?
		return
	}
	webPrint(fmt.Sprintf(" ::: - Radical is set to: %d\n", radical)) // These were Printf 
	webPrint(fmt.Sprintf(" ::: - Work Piece is set to: %d\n", workPiece))
}
