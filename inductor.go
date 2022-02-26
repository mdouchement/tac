package tac

import (
	"math"
)

// Choke computes the choke value you should use in your power supply unit.
//
// C is the capacitor value after the rectifier and before the Choke, usually ~47µF.
//
// F is the resonant frequency that we want to keep below 10Hz.
//
// http://www.valvewizard.co.uk/smoothing.html
// -> Don't forget that the choke must handle the poweramp screen current and preamp current.
func Choke(C, F float64) float64 {
	return 1 / (C * math.Pow(2*math.Pi*F, 2))
}

// Load computes output transformer load.
func Load(DC, power float64) float64 {
	return math.Pow(DC, 2) / power
}
