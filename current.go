package tac

import (
	"math"
)

// TransformerCurrentFWC computes an approximation of the
// transformer's DC current after Full Wave rectifier with a Capacitor input load.
// https://www.hammfg.com/electronics/transformers/rectifier
func TransformerCurrentFWC(Uac1, Uac2, Iac float64) float64 {
	if Uac1 > 0 && Uac2 > 0 {
		// Center tap transformer.
		return Iac
	}

	// Non-center tap transformer.
	return Iac * 0.62
}

// PreampCurrent computes the maximum current consumption for the given n triode.
func PreampCurrent(v PreampValve, U float64, n int) float64 {
	return v.Wa / U * float64(n)
}

// PowerCurrent computes the maximum current consumption for the given n triode.
//
// http://www.valvewizard.co.uk/pp.html
// http://www.valvewizard.co.uk/smoothing.html (Design Example)
func PowerCurrent(v PowerValve, topology string, n int, DC, Rload float64) (Ia, Is float64) {

	Ia = v.Power / DC * float64(n)
	ratio := v.SingleEnded.Ia / v.SingleEnded.Is

	if topology == "PP" {
		if Rload <= 0 {
			Rload = Load(DC, v.Power)
		}
		DC50 := DC - 50

		P := float64(n) * math.Pow(DC50, 2) / Rload
		Ia = P / DC50

		ratio = v.PushPull.Ia / v.PushPull.Is
	}

	return Ia, Ia / ratio
}
