package asm

import (
	"fmt"
	"strings"
)

type Hardware int64

const (
	Holtek     Hardware = 0
	ATMega8    Hardware = 1
	ArduinoTPS Hardware = 2
	TinyTPS    Hardware = 3
	RPI2040    Hardware = 4

	// private parts
	sHoltek     string = "HOLTEK"
	sATMega8    string = "ATMEGA8"
	sArduinoTPS string = "ARDUINOTPS"
	sTinyTPS    string = "TINYTPS"
	sRPI2040    string = "RPI2040"
)

func (h Hardware) String() string {
	switch h {
	case Holtek:
		return sHoltek
	case ATMega8:
		return sATMega8
	case ArduinoTPS:
		return sArduinoTPS
	case TinyTPS:
		return sTinyTPS
	case RPI2040:
		return sRPI2040
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
	case sArduinoTPS:
		return ArduinoTPS
	case sTinyTPS:
		return TinyTPS
	case sRPI2040:
		return RPI2040
	}
	return Holtek
}

func checkSize(h Hardware, size int) error {
	switch h {
	case Holtek:
		if size > 128 {
			return fmt.Errorf("program exceed size limit. Max size for %s is %d", h.String(), 128)
		}
	case ATMega8:
		if size > 256 {
			return fmt.Errorf("program exceed size limit. Max size for %s is %d", h.String(), 256)
		}
	case ArduinoTPS:
		if size > 1024 {
			return fmt.Errorf("program exceed size limit. Max size for %s is %d", h.String(), 1024)
		}
	case RPI2040:
		if size > 1024 {
			return fmt.Errorf("program exceed size limit. Max size for %s is %d", h.String(), 1024)
		}
	case TinyTPS:
		if size > 512 {
			return fmt.Errorf("program exceed size limit. Max size for %s is %d", h.String(), 512)
		}
	default:
		if size > 128 {
			return fmt.Errorf("program exceed size limit. Max size for %s is %d", "Default", 128)
		}
	}
	return nil
}
