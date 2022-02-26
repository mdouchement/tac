package tac

import "math"

// Vpp computes the peak-to-peak voltage of an AC voltage.
//
// e.g. 230VAC is about 325Vpp.
func Vpp(U float64) float64 {
	return U * math.Sqrt(2)
}
