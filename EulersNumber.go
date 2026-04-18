package main

import (
	"fmt"
	"math"
)

func EulersNumber(done chan bool, webPrint func(string)) { 
	select {
	case <-done:
		return
	default:
	}
	var n float64
	var sum float64
	var Eulers float64
	// TODO: How do we make the following line print in a much larger font?
	webPrint("LARGE:Euler's Number \u2107 or \u2147 the natural logarithmic base")
	webPrint("LARGE:\u2147 = (1+1/n)^n")
	webPrint("LARGE:... the limit of an increasing value for n")
	webPrint("LARGE: ")
	n = 9
	sum = 1 + 1/n
	Eulers = math.Pow(sum, n)
	webPrint(fmt.Sprintf("LARGE:%0.45f", Eulers))
	webPrint(fmt.Sprintf("LARGE: was calculated with an exponent of %0.f ", n))
	webPrint("LARGE: ")
	n = 99
	sum = 1 + 1/n
	Eulers = math.Pow(sum, n)
	webPrint(fmt.Sprintf("LARGE:%0.45f", Eulers))
	webPrint(fmt.Sprintf("LARGE:  was calculated with an exponent of %0.f ", n))
	webPrint("LARGE: ")
	n = 999
	sum = 1 + 1/n
	Eulers = math.Pow(sum, n)
	webPrint(fmt.Sprintf("LARGE:%0.45f", Eulers))
	webPrint(fmt.Sprintf("LARGE: was calculated with an exponent of %0.f ", n))
	webPrint("LARGE: ")
	n = 9999
	sum = 1 + 1/n
	Eulers = math.Pow(sum, n)
	webPrint(fmt.Sprintf("LARGE:%0.45f", Eulers))
	webPrint(fmt.Sprintf("LARGE: was calculated with an exponent of %0.f ", n))
	webPrint("LARGE: ")
	n = 99999999999
	sum = 1 + 1/n
	Eulers = math.Pow(sum, n)
	webPrint(fmt.Sprintf("LARGE:%0.45f", Eulers))
	webPrint(fmt.Sprintf("LARGE: was calculated with an exponent of %0.f ", n))
	webPrint("LARGE: ")
	webPrint("LARGE: ")
	webPrint("LARGE:2.71828182845904523536028747135266249775724 is Euler's Number from the web")
	webPrint("LARGE:2.718281828 is the dollar value of $1 compounded continuously for one year.")
	webPrint("LARGE:2.714567 is from daily compound interest which is near-enough to continuous interest.")
	webPrint("LARGE: ")
	webPrint("LARGE:An account starts with $1.00 and pays 100 percent interest per year. If the interest is credited once,")
	webPrint("LARGE:at the end of the year, the value of the account at year-end will be $2.00. What happens if the interest")
	webPrint("LARGE:is computed and credited more frequently during the year?")
	webPrint("LARGE: ")
	webPrint("LARGE:If the interest is credited twice in the year, the interest rate for each 6 months will be 50%, so the ")
	webPrint("LARGE:initial $1 is multiplied by 1.5 twice, yielding $2.25 at the end of the year. Compounding quarterly")
	webPrint("LARGE:yields $2.44140625, and compounding monthly yields $2.613035 = $1.00 × (1 + 1/12)^12 Generally, if there")
	webPrint("LARGE:are n compounding intervals, the interest for each interval will be 100%/n and the value at the end of")
	webPrint("LARGE:the year will be $1.00 × (1 + 1/n)^n.")
	webPrint("LARGE: ")
	webPrint("LARGE:Bernoulli noticed that this sequence approaches a limit (the force of interest) with larger n and, thus, smaller")
	webPrint("LARGE:compounding intervals. Compounding weekly (n = 52) yields $2.692596..., while compounding daily (n = 365) yields")
	webPrint("LARGE:$2.714567... (approximately two cents more). The limit as n grows large is the number that came to be known as e.")
	webPrint("LARGE:That is, with continuous compounding, the account value will reach $2.718281828")

	// written entirely by Richard Woolley
} 
