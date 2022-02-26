package aio

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/mdouchement/tac"
)

type Controller struct {
	Logger  *Logger
	Config  Config
	Results Results
}

func (c *Controller) Execute() error {
	if err := c.heaters(); err != nil {
		return err
	}

	if err := c.voltage(); err != nil {
		return err
	}

	c.transformer()
	c.diode()
	c.choke()
	c.power()
	c.preamp()
	c.bypass()
	c.led()

	c.Results.Rload = c.Config.PowerValves.Rload
	if c.Results.Rload <= 0 {
		c.Results.Rload = tac.Load(c.Results.DC, tac.PowerValves[c.Config.PowerValves.Model].Power)
	}

	return nil
}

func (c *Controller) voltage() (err error) {
	voltages := strings.Split(c.Config.Voltage, "-")

	switch {
	case len(voltages) == 2:
		c.Results.Uac1 = 0
		c.Results.Uac2, err = strconv.ParseFloat(voltages[1], 64)
		if err != nil {
			return err
		}
	case len(voltages) == 3:
		c.Results.Uac1, err = strconv.ParseFloat(voltages[0], 64)
		if err != nil {
			return err
		}

		c.Results.Uac2, err = strconv.ParseFloat(voltages[2], 64)
		if err != nil {
			return err
		}
	default:
		return errors.New("invalid voltage format")
	}

	//
	//
	//

	coef, ok := tac.Rectifiers[c.Config.Rectifier]
	if !ok {
		return fmt.Errorf("unsupported rectifier: %s", c.Config.Rectifier)
	}
	c.Results.DC = c.Results.Uac2 * coef

	return nil
}

func (c *Controller) transformer() {
	c.Results.IDC = tac.TransformerCurrentFWC(c.Results.Uac1, c.Results.Uac2, c.Config.Current)
	c.Results.WDC = c.Results.IDC * c.Results.DC
}

func (c *Controller) heaters() error {
	c.Logger.Verbose("Heaters:")

	config := c.Config.PowerValves
	valve, ok := tac.PowerValves[config.Model]
	if !ok {
		return fmt.Errorf("unsupported power valve: %s", config.Model)
	}

	c.Results.If = valve.If * float64(config.NumberOfValves)
	c.Logger.Verbosef("- %-9s%d * %s = %s\n",
		config.Model+":",
		config.NumberOfValves,
		tac.FormatUnit(valve.If, 1, "A"),
		tac.FormatUnit(c.Results.If, 1, "A"),
	)

	//

	type preamp struct {
		triodes int
		current float64
	}
	m := map[string]*preamp{}

	//

	for _, config := range c.Config.PreampValves {
		valve, ok := tac.PreampValves[config.Model]
		if !ok {
			return fmt.Errorf("unsupported preamp valve: %s", config.Model)
		}

		if _, ok := m[config.Model]; !ok {
			m[config.Model] = &preamp{
				current: valve.If,
			}
		}
		m[config.Model].triodes += config.NumberOfTriodes
	}

	for model, p := range m {
		i := p.current * float64(p.triodes/2+p.triodes%2)
		c.Logger.Verbosef("- %-9s%d * %s = %s\n",
			model+":",
			p.triodes/2+p.triodes%2,
			tac.FormatUnit(p.current, 1, "A"),
			tac.FormatUnit(i, 1, "A"),
		)
		c.Results.If += i
	}

	return nil
}

func (c *Controller) diode() {
	if c.Config.Rectifier == "DIODE" {
		c.Results.Diode = tac.RectifierDiodeRating(c.Results.Uac1, c.Results.Uac2)
	}
}

func (c *Controller) choke() {
	if c.Config.Choke.FirstCapacitor > 0 {
		c.Results.Choke = tac.Choke(c.Config.Choke.FirstCapacitor, 10)
		c.Logger.Verbosef("Choke: %s\n", tac.FormatUnit(c.Results.Choke, 2, "H"))

		c.Results.ChokeScreen = c.Config.Choke.IncludeScreen
	}
}

func (c *Controller) preamp() {
	c.Logger.Verbose("Preamp consumed:")
	for _, v := range c.Config.PreampValves {
		valve := tac.PreampValves[v.Model]

		voltage := valve.Ua
		if v.Voltage > 0 {
			c.Logger.Println(`WARN: preamp: expiremental current computation ("voltage: 0" to disable)`)
			voltage = v.Voltage
		}

		Ia := tac.PreampCurrent(valve, voltage, v.NumberOfTriodes)
		Iat := float64(v.NumberOfTriodes) * (valve.Wat / tac.PreampValveTypicalVoltage)
		Wa := float64(v.NumberOfTriodes) * valve.Wa

		c.Logger.Verbosef("- %-9stypical: %d * %s = %s\n           maximum: %d * %s = %s (%s)\n",
			v.Model+":",
			v.NumberOfTriodes, tac.FormatUnit(valve.Wat/tac.PreampValveTypicalVoltage, 1, "A"), tac.FormatUnit(Iat, 1, "A"),
			v.NumberOfTriodes, tac.FormatUnit(tac.PreampCurrent(valve, voltage, 1), 1, "A"), tac.FormatUnit(Ia, 1, "A"),
			tac.FormatUnit(Wa, 2, "W"),
		)

		c.Results.PreampIa += Ia
		c.Results.PreampWa += Wa
	}
}

func (c *Controller) power() {
	c.Logger.Verbose("Power consumed:")

	valve := tac.PowerValves[c.Config.PowerValves.Model]
	c.Results.PowerIa, c.Results.PowerIs = tac.PowerCurrent(
		valve,
		c.Config.PowerValves.Topology,
		c.Config.PowerValves.NumberOfValves,
		c.Results.DC,
		c.Config.PowerValves.Rload,
	)
	c.Results.PowerWa = valve.Power * float64(c.Config.PowerValves.NumberOfValves)

	c.Logger.Verbosef("- Wa: %-9s- for %d valve(s) in %s\n", tac.FormatUnit(c.Results.PowerWa, 1, "W"), c.Config.PowerValves.NumberOfValves, c.Config.PowerValves.Topology)
	c.Logger.Verbosef("- Ia: %-9s- at plate for %d valve(s) in %s\n", tac.FormatUnit(c.Results.PowerIa, 1, "A"), c.Config.PowerValves.NumberOfValves, c.Config.PowerValves.Topology)
	c.Logger.Verbosef("- Is: %-9s- at screen for %d valve(s) in %s\n", tac.FormatUnit(c.Results.PowerIs, 1, "A"), c.Config.PowerValves.NumberOfValves, c.Config.PowerValves.Topology)
}

func (c *Controller) bypass() {
	for _, Rk := range c.Config.CathodeBypassCapacitors {
		c.Results.Cks = append(c.Results.Cks, CathodeBypassCapacitor{
			Rk: Rk,
			Ck: tac.CathodeBypassCapacitor(Rk, 5),
		})
	}
}

func (c *Controller) led() {
	c.Results.LED.VAC = c.Config.LEDResistor.VAC
	c.Results.LED.R = tac.LEDResistor(
		c.Config.LEDResistor.VAC,
		c.Config.LEDResistor.DCForwardVoltage,
		c.Config.LEDResistor.Current,
	)
	c.Results.LED.W = tac.PowerRI(c.Results.LED.R, c.Config.LEDResistor.Current)
}
