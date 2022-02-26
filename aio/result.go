package aio

import (
	"fmt"
	"strings"

	"github.com/mdouchement/tac"
)

type (
	Results struct {
		Uac1        float64
		Uac2        float64
		Diode       float64
		DC          float64
		IDC         float64
		WDC         float64
		Choke       float64
		ChokeScreen bool
		Rload       float64
		PowerWa     float64
		PowerIa     float64
		PowerIs     float64
		PreampWa    float64
		PreampIa    float64
		If          float64
		Cks         []CathodeBypassCapacitor
		LED         LED
	}

	CathodeBypassCapacitor struct {
		Rk float64
		Ck float64
	}

	LED struct {
		VAC float64
		R   float64
		W   float64
	}
)

func (r Results) String() string {
	var lines []string

	voltage := "0-" + tac.FormatUnit(r.Uac2, 0, "")
	if r.Uac1 > 0 {
		voltage = tac.FormatUnit(r.Uac2, 0, "") + "-" + voltage
	}

	lines = append(lines,
		fmt.Sprintf("Selected transformer voltage: %s", voltage),
	)

	if r.Diode > 0 {
		diode := tac.Tolerance(tac.Tolerance(r.Diode))
		lines = append(lines,
			fmt.Sprintf("Diode Reverse Repetative Maximum (Vrrm) rating: %s", tac.FormatUnit(diode, 1, "V")),
			"  -> We knock off 10% to allow for variation in mains voltage, and knock off another 10% to allow for the transformer voltage being high if loaded only lightly.",
			"  -> The popular 1N4007 is rated for 1000V. Two 1N4007 in series can handle 2000v.",
			"  -> http://www.valvewizard.co.uk/bridge.html",
		)
	}
	lines = append(lines,
		fmt.Sprintf("Calculated voltage at first capacitor (B+): %s", tac.FormatUnit(r.DC, 1, "V")),
	)

	if r.WDC > 0 {
		lines = append(lines,
			fmt.Sprintf("Calculated DC current at first capacitor: %s", tac.FormatUnit(r.IDC, 2, "A")),
			fmt.Sprintf("Calculated DC wattage at first capacitor: %s", tac.FormatUnit(r.WDC, 1, "W")),
		)
	}

	//
	//
	//

	if r.Choke > 0 {
		I := r.PowerIs
		if r.ChokeScreen {
			I += r.PreampIa
		}
		I = tac.Tolerance(I)
		lines = append(lines,
			fmt.Sprintf("Calculated choke: %s @ %s (10%% plus factored in)", tac.FormatUnit(tac.Tolerance(r.Choke), 2, "H"), tac.FormatUnit(I, 2, "A")),
			"  -> Resonant frequency is kept below 10Hz.",
		)
	}

	//
	//
	//

	I := tac.Tolerance(r.PowerIa + r.PowerIs + r.PreampIa)
	lines = append(lines,
		fmt.Sprintf("Calculated current: %s (10%% plus factored in) at %s load", tac.FormatUnit(I, 2, "A"), tac.FormatUnit(r.Rload, 1, "Ω")),
		"  -> Maximum preamp current consumption is used.",
		fmt.Sprintf("Calculated maximal valve(s) wattage: %s", tac.FormatUnit(r.PowerWa+r.PreampWa, 2, "W")),
		fmt.Sprintf("Calculated filament current (typically the 6.3v secondary): %s", tac.FormatUnit(r.If, 1, "A")),
	)

	//
	//
	//

	for _, bypass := range r.Cks {
		lines = append(lines,
			fmt.Sprintf("Calculated cathode bypass capacitor: Ck=%s for Rk=%s @ F=5Hz", tac.FormatUnit(tac.Tolerance(bypass.Ck), 2, "F"), tac.FormatUnit(bypass.Rk, 1, "Ω")),
		)
	}

	if r.LED.R > 0 {
		lines = append(lines,
			fmt.Sprintf("Calculated LED resistor: %s for %.3fW @ VAC=%s", tac.FormatUnit(r.LED.R, 2, "Ω"), tac.Tolerance(r.LED.W), tac.FormatUnit(r.LED.VAC, 2, "V")),
		)
	}

	//
	//
	//

	lines = append(lines,
		"",
		"Disclaimer:",
		"  Please, consult your chosen tube model manufacturer data sheet, especially, when you entered the Rload.",
		"  Then the calculation provided is an estimation based on the assumption that Rload is not relatively small!",
		"  Tubes are non-linear devices and current draw may change under a number of conditions!",
		"  This calculation is provided in the hope of being useful as a reference only.",
		"  Proceed with caution! This tool disclaims any responsibility from any intended or actual use and consequences of this calculation.",
		"",
		"Thanks:",
		"  This tool is fully inspired by https://thesubjectmatter.com/calcptcurrent.html",
		"  It use formula from http://www.valvewizard.co.uk",
	)

	return strings.Join(lines, "\n")
}
