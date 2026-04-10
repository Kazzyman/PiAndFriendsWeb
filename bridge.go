package main

import (
	"fmt"
	"net/url"
	"strconv"
)

func runRootsWeb(params url.Values, fyneFunc func(string)) {
	radStr := params.Get("radical")
	workStr := params.Get("workpiece")

	rad, _ := strconv.Atoi(radStr)
	work, _ := strconv.Atoi(workStr)

	// Bridging to your original roots.go variables
	radical_index = rad

	fyneFunc(fmt.Sprintf("Starting Roots Demo: Radical %d, Workpiece %d\n", rad, work))

	// This calls your original build logic from roots.go
	buildPairsSlice(rad)

	// Note: You can add the specific root-finding call here 
	// based on how you normally trigger it in window2.go
	fyneFunc("\nCalculation complete.\n")
}
