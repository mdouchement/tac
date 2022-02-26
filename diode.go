package tac

// RectifierDiodeRating computes the Reverse Repetative Maximum (Vrrm) rating
// that excess the peak AC voltage.
// For a tap-center transformer like 350-0-350, U1=350 and U2=350.
// For a normal transformer like 0-350, U1=0, U2=350.
//
// http://www.valvewizard.co.uk/bridge.html
//
// The popular 1N4007 is rated for 1000V.
// This corresponds to an AC voltage of 1000V/1.4 = 714Vrms.
// However, we should knock off 10% to allow for variation in mains voltage,
// and knock off another 10% to allow for the transformer voltage being high
// if loaded only lightly. Therefore, we can't use the 1N4007
// if the (advertised) transformer voltage is greater than 580Vrms.
func RectifierDiodeRating(U1, U2 float64) float64 {
	return Vpp(U1 + U2)
}

// LEDResistor computes the resistor value to put in series with the LED.
func LEDResistor(VAC, DCForwardVoltage, I float64) float64 {
	return (Vpp(VAC) - DCForwardVoltage) / I
}
