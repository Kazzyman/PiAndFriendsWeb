package main

// Gauss.go
//
// The Gauss-Legendre algorithm for computing π.
// Developed by Carl Friedrich Gauss, refined by Adrien-Marie Legendre.
//
// Quadratic convergence: correct digits DOUBLE every iteration.
//   Iteration  1 →    1 digit
//   Iteration  5 →   16 digits
//   Iteration 10 →  512 digits
//   Iteration 12 → ~2048 digits  (practical web demo ceiling)
//
// The algorithm maintains four values:
//   a  = arithmetic mean  (starts at 1)
//   b  = geometric mean   (starts at 1/√2)
//   t  = correction term  (starts at 1/4)
//   p  = power of 2       (starts at 1)
//
// Each iteration:
//   a_next = (a + b) / 2
//   b_next = √(a × b)
//   t_next = t - p × (a - a_next)²
//   p_next = 2 × p
//
// Then: π ≈ (a + b)² / (4 × t)
//
// Original float64 version: Richard (Rick) Woolley.
// Rewritten with big.Float and verified-digit output: April 2026.

import (
	"fmt"
	"math/big"
	"strings"
	"time"
)

// gaussMaxIters is the hard ceiling on iterations.
// At iteration 12 you get ~2048 theoretical digits, which takes a few
// seconds even on slow hardware. Beyond 12 the runtime grows rapidly.
const gaussMaxIters = 12

// gaussVerifiedDigits walks the computed π string character by character
// against piForGauss (our ~3000-digit reference) and returns the count
// of digits after "3." that are correct.
// Only verified digits are ever shown to the user -- no estimated output.
func gaussVerifiedDigits(computed string) int {
	computed = strings.TrimSpace(computed)
	ref      := piForGauss

	// Walk both strings together. Both start with "3." so we skip
	// the first two characters when counting, but include them in
	// the comparison to keep the indices aligned.
	count := 0
	for i := 0; i < len(ref) && i < len(computed); i++ {
		if computed[i] != ref[i] {
			break
		}
		if i >= 2 { // past the "3." prefix
			count++
		}
	}
	return count
}

// Gauss_Legendre runs the Gauss-Legendre algorithm for the requested
// number of iterations, printing only verified correct digits after each.
func Gauss_Legendre(iters int, webPrint func(string)) {

	// ── Validate user input ───────────────────────────────────────────────
	if iters > gaussMaxIters {
		webPrint(fmt.Sprintf("  The maximum number of iterations is %d.", gaussMaxIters))
		webPrint("  Beyond that, runtime becomes impractical for a web demo.")
		webPrint(fmt.Sprintf("  Please choose a number from 1 to %d.", gaussMaxIters))
		return
	}
	if iters < 1 {
		iters = 1
	}

	// ── Set big.Float precision dynamically ───────────────────────────────
	// Each iteration doubles the correct digits, so we need 2^iters decimal
	// digits of backing precision. Since 3.32 bits ≈ 1 decimal digit, we
	// set bits = 2^iters * 4, with a floor of 64 to keep early iters sane.
	precBits := uint(1<<uint(iters)) * 4
	if precBits < 64 {
		precBits = 64
	}

	usingBigFloats = true
	start := time.Now()

	webPrint("  ── Gauss-Legendre Algorithm ───────────────────────────────────")
	webPrint("  C F Gauss / Adrien-Marie Legendre")
	webPrint("  Convergence: QUADRATIC — correct digits double each iteration")
	webPrint(fmt.Sprintf("  Iterations requested : %d  (max %d)", iters, gaussMaxIters))
	webPrint(fmt.Sprintf("  Precision            : %d bits (~%d decimal digits)",
		precBits, precBits/4))
	webPrint("  ───────────────────────────────────────────────────────────────")
	webPrint("")
	webPrint("  Starting values:")
	webPrint("    a₀ = 1")
	webPrint("    b₀ = 1 / √2")
	webPrint("    t₀ = 1/4")
	webPrint("    p₀ = 1")
	webPrint("")

	// ── Initialize the four values ────────────────────────────────────────

	two  := new(big.Float).SetPrec(precBits).SetFloat64(2.0)
	four := new(big.Float).SetPrec(precBits).SetFloat64(4.0)

	// a₀ = 1
	a := new(big.Float).SetPrec(precBits).SetFloat64(1.0)

	// b₀ = 1 / √2
	b := new(big.Float).SetPrec(precBits).Sqrt(two)
	b.Quo(new(big.Float).SetPrec(precBits).SetFloat64(1.0), b)

	// t₀ = 0.25
	t := new(big.Float).SetPrec(precBits).SetFloat64(0.25)

	// p₀ = 1
	p := new(big.Float).SetPrec(precBits).SetFloat64(1.0)

	// ── Iterate ───────────────────────────────────────────────────────────

	refDigits := len(piForGauss) - 2 // total digits in our reference (minus "3.")

	for i := 1; i <= iters; i++ {

		// a_next = (a + b) / 2   (arithmetic mean)
		aNext := new(big.Float).SetPrec(precBits).Add(a, b)
		aNext.Quo(aNext, two)

		// b_next = √(a × b)   (geometric mean)
		bNext := new(big.Float).SetPrec(precBits).Mul(a, b)
		bNext.Sqrt(bNext)

		// t_next = t - p × (a - a_next)²
		// The term (a - a_next) is the key: it measures how far a moved
		// this iteration. As a and b converge toward their AGM, this
		// difference shrinks quadratically -- hence the speed.
		diff := new(big.Float).SetPrec(precBits).Sub(a, aNext)
		diff.Mul(diff, diff)  // square it
		diff.Mul(p, diff)     // multiply by p
		tNext := new(big.Float).SetPrec(precBits).Sub(t, diff)

		// p_next = 2 × p
		pNext := new(big.Float).SetPrec(precBits).Mul(two, p)

		// Advance all four values for the next iteration
		a, b, t, p = aNext, bNext, tNext, pNext

		// ── Compute current π estimate ────────────────────────────────
		// π ≈ (a + b)² / (4 × t)
		sumAB := new(big.Float).SetPrec(precBits).Add(a, b)
		pi    := new(big.Float).SetPrec(precBits).Mul(sumAB, sumAB)
		pi.Quo(pi, new(big.Float).SetPrec(precBits).Mul(four, t))

		// Convert to string with enough decimal places to saturate our
		// reference. We then count only what is actually correct.
		piStr    := pi.Text('f', refDigits+10)
		verified := gaussVerifiedDigits(piStr)

		// Cap at our reference length -- we cannot claim more than we
		// can verify.
		if verified > refDigits {
			verified = refDigits
		}

		// Build a display string containing exactly the verified digits.
		// "3." prefix + verified digits.
		displayLen := verified + 2
		displayStr := piStr
		if len(displayStr) > displayLen {
			displayStr = displayStr[:displayLen]
		}

		// The gap between a and b shrinks quadratically -- showing it
		// makes the convergence speed viscerally obvious.
		gap := new(big.Float).SetPrec(precBits).Sub(a, b)

		webPrint(fmt.Sprintf("  Iteration %2d  →  %d verified correct digits:", i, verified))
		webPrint(fmt.Sprintf("  π = %s", displayStr))
		webPrint(fmt.Sprintf("  gap (a-b) = %s  ← approaches zero quadratically",
			gap.Text('e', 4)))
		webPrint("")
	}

	// ── Final summary ─────────────────────────────────────────────────────

	elapsed := time.Since(start)

	// Recompute final π for the summary block
	sumAB := new(big.Float).SetPrec(precBits).Add(a, b)
	pi    := new(big.Float).SetPrec(precBits).Mul(sumAB, sumAB)
	pi.Quo(pi, new(big.Float).SetPrec(precBits).Mul(four, t))
	piStr    := pi.Text('f', refDigits+10)
	verified := gaussVerifiedDigits(piStr)
	if verified > refDigits {
		verified = refDigits
	}
	displayLen := verified + 2
	displayStr := piStr
	if len(displayStr) > displayLen {
		displayStr = displayStr[:displayLen]
	}

	webPrint("  ── Final Result ────────────────────────────────────────────────")
	webPrint(fmt.Sprintf("  %d iterations completed in %s",
		iters, elapsed.Round(time.Millisecond)))
	webPrint(fmt.Sprintf("  %d verified correct digits of π", verified))
	webPrint("")
	webPrint(fmt.Sprintf("  π = %s", displayStr))
	webPrint("  ───────────────────────────────────────────────────────────────")
}
