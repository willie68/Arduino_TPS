package asm

import "strings"

type Hardware int64

const (
	Holtek     Hardware = 0
	ATMega8    Hardware = 1
	ArduinoSPS Hardware = 2
	TinySPS    Hardware = 3
)

func (h Hardware) String() string {
	switch h {
	case Holtek:
		return "HOLTEK"
	case ATMega8:
		return "ATMEGA8"
	case ArduinoSPS:
		return "ARDUINOSPS"
	case TinySPS:
		return "TINYSPS"
	default:
		return "HOLTEK"
	}
}

func ParseHardware(dest string) Hardware {
	dest = strings.ToUpper(dest)
	switch dest {
	case "HOLTEK":
		return Holtek
	case "ATMEGA8":
		return ATMega8
	case "ARDUINOSPS":
		return ArduinoSPS
	case "TINYSPS":
		return TinySPS
	}
	return Holtek
}
