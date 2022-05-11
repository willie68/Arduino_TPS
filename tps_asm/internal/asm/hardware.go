package asm

import "strings"

type Hardware int64

const (
	Holtek     Hardware = 0
	ATMega8    Hardware = 1
	ArduinoSPS Hardware = 2
	TinySPS    Hardware = 3

	// private parts
	sHoltek     string = "HOLTEK"
	sATMega8    string = "ATMEGA8"
	sArduinoSPS string = "ARDUINOSPS"
	sTinySPS    string = "TINYSPS"
)

func (h Hardware) String() string {
	switch h {
	case Holtek:
		return sHoltek
	case ATMega8:
		return sATMega8
	case ArduinoSPS:
		return sArduinoSPS
	case TinySPS:
		return sTinySPS
	default:
		return sHoltek
	}
}

func ParseHardware(dest string) Hardware {
	dest = strings.ToUpper(dest)
	switch dest {
	case sHoltek:
		return Holtek
	case sATMega8:
		return ATMega8
	case sArduinoSPS:
		return ArduinoSPS
	case sTinySPS:
		return TinySPS
	}
	return Holtek
}
