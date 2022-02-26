package tac

import "math"

// CathodeBypassCapacitor computes an approximation of the cathode bypass capacitor value
// according the given cathode resistor Rk in Ohm and the wanted frequency F in Hertz.
//
// For full baypass, use F = 5Hz
func CathodeBypassCapacitor(Rk, F float64) float64 {
	return 1 / (2 * math.Pi * Rk * F)
}

// SmoothingCapacitor computes the capicitor value to put after a voltage dropping resistor.
func SmoothingCapacitor(R float64) float64 {
	return 1 / (2 * math.Pi * R)
}
