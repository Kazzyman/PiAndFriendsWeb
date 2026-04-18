package main

import (
	"fmt"
	"math"
)

func ErdosBorwein(done chan bool, webPrint func(string)) { 
	select {
	case <-done:
		return
	default:
	}

	webPrint("")

	webPrint("We calculate E as E = the sum of 1/((2^n)-1) as n grows from 1 to 'infinity'")

	var Erdos_Borwein float64
	Erdos_Borwein = 1
	var iter float64
	iter = 1
	for iter < 100 {
		iter++ // iter will therefore begin at 2
		Erdos_Borwein = Erdos_Borwein + (1 / ((math.Pow(2, iter)) - 1))
		if iter == 10 || iter == 20 || iter == 30 || iter == 40 || iter == 50 || iter == 60 || iter == 70 || iter == 100 || iter == 101 {
			webPrint(fmt.Sprintf("%0.25f", Erdos_Borwein))
		}
		// 100 and 101 prove that we ended on 100 as the final exponent
		// ... we only get 8 results, not nine
	}
	webPrint("for 10, 20, 30, 40, 50, 60, 70, and 100 iterations respectively\n")
	webPrint("Our calculated Erdos-Borwein constant is ")
	webPrint(fmt.Sprintf("Erdos_Borwein, after, %0.9f, iterations, i.e., with a final exponent of", iter))
	webPrint("1.606695152415291763 is what we get from the web\n")
	/*
	// TODO: the following is what I had in my old CLI version, what "...with a final exponent of", iter)" was doing I cannot now figure?
	// TODO ... and why iter is a float64 is also a mystery ?? Possibly I was just that ignorant back then?
	fmt.Println("")

	fmt.Println(colorCyan, rune, "\n", colorReset)
	fmt.Println("We calculate E as E = the sum of 1/((2^n)-1) as n grows from 1 to 'infinity'")

	var Erdos_Borwein float64
	Erdos_Borwein = 1
	var iter float64
	iter = 1
	for iter < 100 {
		iter++ // iter will therefore begin at 2
		Erdos_Borwein = Erdos_Borwein + (1 / ((math.Pow(2, iter)) - 1))
		if iter == 10 || iter == 20 || iter == 30 || iter == 40 || iter == 50 || iter == 60 || iter == 70 || iter == 100 || iter == 101 {
			fmt.Println(Erdos_Borwein)
		}
		// 100 and 101 prove that we ended on 100 as the final exponent
		// ... we only get 8 results, not nine
	}
	fmt.Println("for 10, 20, 30, 40, 50, 60, 70, and 100 iterations respectively\n")
	fmt.Println("Our calculated Erdos-Borwein constant is ")
	fmt.Println(Erdos_Borwein, "after", iter, "iterations, i.e., with a final exponent of", iter)
	fmt.Println("1.606695152415291763 is what we get from the web\n")
	*/
	// written entirely by Richard Woolley
}
