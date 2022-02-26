package tac

import (
	"fmt"
	"math"
	"strconv"
)

// Unit multiplier.
var (
	UnitPico  = math.Pow(10, -12)
	UnitNano  = math.Pow(10, -9)
	UnitMicro = math.Pow(10, -6)
	UnitMilli = math.Pow(10, -3)
	UnitNone  = math.Pow(10, 0)
	UnitKillo = math.Pow(10, 3)
	UnitMega  = math.Pow(10, 6)
)

// FormatUnit returns a pretty print version of the given value.
func FormatUnit(v float64, precision int, unit string) string {
	var p float64
	var letter string

	switch {
	case v >= UnitMega:
		letter = "M"
		p = math.Pow(10, -6)
	case v >= UnitKillo:
		letter = "k"
		p = math.Pow(10, -3)
	case v >= UnitNone:
		letter = ""
		p = math.Pow(10, 0)
	case v >= UnitMilli:
		letter = "m"
		p = math.Pow(10, 3)
	case v >= UnitMicro:
		letter = "µ"
		p = math.Pow(10, 6)
	case v >= UnitNano:
		letter = "n"
		p = math.Pow(10, 9)
	case v >= UnitPico:
		letter = "p"
		p = math.Pow(10, 12)
	}

	v *= p
	if v == math.Trunc(v) {
		precision = 0
	}

	format := "%." + strconv.Itoa(precision) + "f%s%s"
	return fmt.Sprintf(format, v, letter, unit)
}
