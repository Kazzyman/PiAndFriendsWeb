package main

import (
	"time"
)

// @formatter:off

var calculating bool
var Radical_index int
var copyOfLastPosition int
var iters_bbp int

// convenience globals:
var usingBigFloats = false // a variable of type bool which is passed by many funcs to print Result Stats Long()

var iterationsForMonte16i int
var iterationsForMonte16j int
var iterationsForMonteTotal int
var four float64 // is initialized to 4 where needed
var π float64    // a var can be any character, as in this Pi symbol/character
var LinesPerSecond float64
var LinesPerIter float64
var iterInt64 int64     // to be used primarily in selections which require modulus calculations
var iterFloat64 float64 // to be used in selections which do not require modulus calculations
var t2 time.Time

// The following globals, are used in multiple funcs of case 18: calculate either square or cube root of any integer
var Tim_win float64             // Time Window
// var sortedResults = []Results{} // sortedResults is an array of type Results as defined at the top of this file
var diffOfLarger int
var diffOfSmaller int
var precisionOfRoot int    // this being global means we do not need to pass it in to the read func

const colorReset = "\033[0m"
const colorGreen = "\033[32m"
const colorCyan = "\033[36m"
