package tac

const PreampValveTypicalVoltage float64 = 250

type (
	PreampValve struct {
		// Maximum anode/plate voltage.
		Ua float64
		// Maximum power dissipation.
		Wa float64
		// Typical power dissipation (Ua*Ia from typical characteristics)
		Wat float64
		// Heater current.
		If float64
	}

	PowerValve struct {
		Power       float64
		SingleEnded Currents
		PushPull    Currents
		// Heater current.
		If float64
	}

	Currents struct {
		// Anode/Plate current in Ampère.
		Ia float64
		// Screen current in Ampère. Ig2 in JJ datasheets.
		Is float64
	}
)

// Toloreance adds a 10% of error factor.
func Tolerance(v float64) float64 {
	return 1.10 * v
}

// Rectifiers' coef.
// SQRT(2) = ~1.4 is the ideal coef.
// Based on https://thesubjectmatter.com/calcptcurrent.html
var Rectifiers = map[string]float64{
	"DIODE": 1.37,
	"GZ34":  1.36,
	"EZ81":  1.30,
	"5U4B":  1.28,
	"5Y3":   1.25,
}

// PreampValves characteristics.
// Based on https://www.jj-electronic.com/en/preamplifying-tubes
var PreampValves = map[string]PreampValve{
	"ECC81": {
		Ua:  300,
		Wa:  2.5,
		Wat: PreampValveTypicalVoltage * 10 * UnitMilli,
		If:  0.3, // At 6.3V

	},
	"ECC83S": {
		Ua:  300,
		Wa:  1,
		Wat: PreampValveTypicalVoltage * 1.2 * UnitMilli,
		If:  0.3, // At 6.3V
	},
	"ECC803S": {
		Ua:  300,
		Wa:  1.2,
		Wat: PreampValveTypicalVoltage * 1.2 * UnitMilli,
		If:  0.3, // At 6.3V
	},
	"6SL7": {
		Ua:  300,
		Wa:  1,
		Wat: PreampValveTypicalVoltage * 2.3 * UnitMilli,
		If:  0.3, // At 6.3V
	},
}

// PowerValves characteristics.
// Based on https://www.jj-electronic.com/en/power-tubes
var PowerValves = map[string]PowerValve{
	"EL34": {
		Power: 25,
		SingleEnded: Currents{
			Ia: 100 * UnitMilli,
			Is: 14.9 * UnitMilli,
		},
		PushPull: Currents{
			Ia: 100 * UnitMilli,
			Is: 14.9 * UnitMilli,
		},
		If: 1.5,
	},
	"6L6GC": {
		Power: 30,
		SingleEnded: Currents{
			Ia: 72 * UnitMilli,
			Is: 5 * UnitMilli,
		},
		PushPull: Currents{
			Ia: 134 * UnitMilli,
			Is: 11 * UnitMilli,
		},
		If: 0.9,
	},
	"6V6S": {
		Power: 14,
		SingleEnded: Currents{
			Ia: 45 * UnitMilli,
			Is: 5 * UnitMilli,
		},
		PushPull: Currents{
			Ia: 70 * UnitMilli,
			Is: 13 * UnitMilli,
		},
		If: 0.5,
	},
}
