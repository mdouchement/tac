package tac

import "math"

// PowerRI computes the power for the given R and I.
func PowerRI(R, I float64) float64 {
	return R * math.Pow(I, 2)
}
