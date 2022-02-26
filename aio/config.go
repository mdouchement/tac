package aio

type Config struct {
	Voltage   string  `yaml:"voltage"`
	Current   float64 `yaml:"current"`
	Rectifier string  `yaml:"rectifier"`
	Choke     struct {
		IncludeScreen  bool    `yaml:"include_screen"`
		FirstCapacitor float64 `yaml:"first_capacitor"`
	} `yaml:"choke"`
	PowerValves struct {
		Model          string  `yaml:"model"`
		NumberOfValves int     `yaml:"number_of_valves"`
		Topology       string  `yaml:"topology"`
		Rload          float64 `yaml:"rload"`
	} `yaml:"power_valves"`
	PreampValves []struct {
		Model           string  `yaml:"model"`
		Voltage         float64 `yaml:"voltage"`
		NumberOfTriodes int     `yaml:"number_of_triodes"`
	} `yaml:"preamp_valves"`
	CathodeBypassCapacitors []float64 `yaml:"cathode_bypass_capacitors"`
	LEDResistor             struct {
		VAC              float64 `yaml:"vac"`
		DCForwardVoltage float64 `yaml:"dc_forward_voltage"`
		Current          float64 `yaml:"current"`
	} `yaml:"led_resistor"`
}
