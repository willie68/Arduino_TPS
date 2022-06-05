package asm

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

var (
	ErrParamCount   = errors.New("parameter count is wrong")
	ErrIllegalValue = errors.New("error on value")
)

type mnemonic struct {
	Name  string
	Param []string
	Enums map[string]int
	Code  byte
	H     []Hardware
}

func GetMnemonic(name string) (*mnemonic, error) {
	name = strings.ToUpper(name)
	for _, mno := range Mnos {
		if name == mno.Name {
			return &mno, nil
		}
	}
	return nil, fmt.Errorf("unknown mnemonics: %s ", name)
}

const (
	int4 string = "int4"
	lbl  string = "label"
	enum string = "enum"
	nme  string = "name"
)

var Mnos = []mnemonic{
	// program control
	{
		Name:  "NOP",
		Param: []string{},
		Code:  0x00,
	},
	{
		Name:  "WAIT",
		Param: []string{int4, enum},
		Enums: map[string]int{
			"1MS":   0x00,
			"2MS":   0x01,
			"5MS":   0x02,
			"10MS":  0x03,
			"20MS":  0x04,
			"50MS":  0x05,
			"100MS": 0x06,
			"200MS": 0x07,
			"500MS": 0x08,
			"1S":    0x09,
			"2S":    0x0A,
			"5S":    0x0B,
			"10S":   0x0C,
			"20S":   0x0D,
			"30S":   0x0E,
			"60S":   0x0F,
		},
		Code: 0x20,
	},
	{
		Name:  "RJMP",
		Param: []string{int4, lbl},
		Code:  0x30,
	},
	{
		Name:  "PAGE",
		Param: []string{int4, enum},
		Enums: map[string]int{
			":?": 0x10,
		},
		Code: 0x80,
		H: []Hardware{
			ArduinoSPS,
			TinySPS,
		},
	},
	{
		Name:  "JMP",
		Param: []string{int4, lbl},
		Code:  0x90,
	},
	{
		Name:  "LOOPC",
		Param: []string{int4, lbl},
		Code:  0xA0,
	},
	{
		Name:  "LOOPD",
		Param: []string{int4, lbl},
		Code:  0xB0,
	},
	{
		Name:  "SKIP0",
		Param: []string{},
		Code:  0xC0,
		H: []Hardware{
			ArduinoSPS,
			TinySPS,
		},
	},
	{
		Name:  "AGTB",
		Param: []string{},
		Code:  0xC1,
	},
	{
		Name:  "ALTB",
		Param: []string{},
		Code:  0xC2,
	},
	{
		Name:  "AEQB",
		Param: []string{},
		Code:  0xC3,
	},
	{
		Name:  "DEQ0",
		Param: []string{enum},
		Enums: map[string]int{
			"1": 0x08,
			"2": 0x09,
			"3": 0x0A,
			"4": 0x0B,
		},
		Code: 0xC0,
	},
	{
		Name:  "DEQ1",
		Param: []string{enum},
		Enums: map[string]int{
			"1": 0x04,
			"2": 0x05,
			"3": 0x06,
			"4": 0x07,
		},
		Code: 0xC0,
	},
	{
		Name:  "PRG0",
		Param: []string{},
		Code:  0xCC,
	},
	{
		Name:  "SEL0",
		Param: []string{},
		Code:  0xCD,
	},
	{
		Name:  "PRG1",
		Param: []string{},
		Code:  0xCE,
	},
	{
		Name:  "SEL1",
		Param: []string{},
		Code:  0xCF,
	},
	{
		Name:  "CALL",
		Param: []string{int4, lbl},
		Code:  0xD0,
	},
	{
		Name:  "RTR",
		Param: []string{},
		Code:  0xE0,
	},
	{
		Name:  "CASB",
		Param: []string{int4, lbl},
		Code:  0xE0,
		H: []Hardware{
			ArduinoSPS,
			TinySPS,
		},
	},
	{
		Name:  "DFSB",
		Param: []string{int4, lbl},
		Code:  0xE7,
		H: []Hardware{
			ArduinoSPS,
			TinySPS,
		},
	},
	{
		Name:  "REST",
		Param: []string{},
		Code:  0xEF,
		H: []Hardware{
			ArduinoSPS,
			TinySPS,
		},
	},
	{
		Name:  "PEND",
		Param: []string{},
		Code:  0xFF,
		H: []Hardware{
			ArduinoSPS,
			TinySPS,
		},
	},

	// Load and save
	{
		Name:  "LDA",
		Param: []string{int4, enum},
		Enums: map[string]int{
			"DIN":  0x64,
			"DIN1": 0x65,
			"DIN2": 0x66,
			"DIN3": 0x67,
			"DIN4": 0x68,
			"ADC1": 0x69,
			"ADC2": 0x6A,
			"RC1":  0x6B,
			"RC2":  0x6C,
		},
		Code: 0x40,
		H: []Hardware{
			ArduinoSPS,
			TinySPS,
		},
	},
	{
		Name:  "SWAP",
		Param: []string{},
		Code:  0x50,
		H: []Hardware{
			ArduinoSPS,
			TinySPS,
		},
	},
	{
		Name:  "MOV",
		Param: []string{enum},
		Enums: map[string]int{
			"B,A": 0x51,
			"C,A": 0x52,
			"D,A": 0x53,
			"A,B": 0x61,
			"A,C": 0x62,
			"A,D": 0x63,
			"E,A": 0x5D,
			"F,A": 0x5E,
			"A,E": 0x6D,
			"A,F": 0x6E,
		},
		Code: 0x00,
		H: []Hardware{
			ArduinoSPS,
			TinySPS,
		},
	},
	{
		Name:  "PUSH",
		Param: []string{},
		Code:  0x5F,
		H: []Hardware{
			ArduinoSPS,
			TinySPS,
		},
	},
	{
		Name:  "POP",
		Param: []string{},
		Code:  0x6F,
		H: []Hardware{
			ArduinoSPS,
			TinySPS,
		},
	},
	// Math
	{
		Name:  "INC",
		Param: []string{},
		Code:  0x71,
	},
	{
		Name:  "DEC",
		Param: []string{},
		Code:  0x72,
	},
	{
		Name:  "ADD",
		Param: []string{},
		Code:  0x73,
	},
	{
		Name:  "SUB",
		Param: []string{},
		Code:  0x74,
	},
	{
		Name:  "MUL",
		Param: []string{},
		Code:  0x75,
	},
	{
		Name:  "DIV",
		Param: []string{},
		Code:  0x76,
	},
	{
		Name:  "AND",
		Param: []string{},
		Code:  0x77,
	},
	{
		Name:  "OR",
		Param: []string{},
		Code:  0x78,
	},
	{
		Name:  "XOR",
		Param: []string{},
		Code:  0x79,
	},
	{
		Name:  "NOT",
		Param: []string{},
		Code:  0x7A,
	},
	{
		Name:  "MOD",
		Param: []string{},
		Code:  0x7B,
		H: []Hardware{
			ArduinoSPS,
			TinySPS,
		},
	},
	{
		Name:  "BYTE",
		Param: []string{},
		Code:  0x7C,
		H: []Hardware{
			ArduinoSPS,
			TinySPS,
		},
	},
	{
		Name:  "BSUBA",
		Param: []string{},
		Code:  0x7D,
		H: []Hardware{
			ArduinoSPS,
			TinySPS,
		},
	},
	{
		Name:  "SHR",
		Param: []string{},
		Code:  0x7E,
		H: []Hardware{
			ArduinoSPS,
			TinySPS,
		},
	},
	{
		Name:  "SHL",
		Param: []string{},
		Code:  0x7F,
		H: []Hardware{
			ArduinoSPS,
			TinySPS,
		},
	},
	// Input/Output
	{
		Name:  "PORT",
		Param: []string{int4},
		Code:  0x10,
	},
	{
		Name:  "STA",
		Param: []string{enum},
		Enums: map[string]int{
			"DOUT":  0x54,
			"DOUT1": 0x55,
			"DOUT2": 0x56,
			"DOUT3": 0x57,
			"DOUT4": 0x58,
			"PWM1":  0x59,
			"PWM2":  0x5A,
			"SRV1":  0x5B,
			"SRV2":  0x5C,
		},
		Code: 0x00,
		H: []Hardware{
			ArduinoSPS,
			TinySPS,
		},
	},
	// Byte mnemonics
	{
		Name:  "BLDA",
		Param: []string{enum},
		Enums: map[string]int{
			"ADC1": 0xF0,
			"ADC2": 0xF1,
			"RC1":  0xF2,
			"RC2":  0xF3,
		},
		Code: 0x00,
		H: []Hardware{
			ArduinoSPS,
			TinySPS,
		},
	},
	{
		Name:  "BSTA",
		Param: []string{enum},
		Enums: map[string]int{
			"PWM1": 0xF4,
			"PWM2": 0xF5,
			"SRV1": 0xF6,
			"SRV2": 0xF7,
		},
		Code: 0x00,
		H: []Hardware{
			ArduinoSPS,
			TinySPS,
		},
	},
	{
		Name:  "TONE",
		Param: []string{},
		Code:  0xF8,
		H: []Hardware{
			ArduinoSPS,
			TinySPS,
		},
	},
}

func (m mnemonic) CheckHardware(h Hardware) error {
	if len(m.H) > 0 {
		switch m.Name {
		case "PAGE":
			return nil
		case "LDA":
			return nil
		case "MOV":
			return nil
		case "STA":
			return nil
		default:
			f := false
			for _, mh := range m.H {
				if mh.String() == h.String() {
					f = true
				}
			}
			if !f {
				return fmt.Errorf("%s Mnemonic isn't support on %s Hardware", m.Name, h.String())
			}
		}
	}
	return nil
}

func (m mnemonic) CheckParameter(param string) error {
	if (param == "") && (len(m.Param) > 0) {
		return ErrParamCount
	}
	if (param != "") && (len(m.Param) == 0) {
		return ErrParamCount
	}
	if len(m.Param) == 0 {
		return nil
	}
	var err error
	found := false
	for _, pt := range m.Param {
		switch pt {
		case int4:
			_, found, err = convertNumber(param)
			if err != nil {
				return err
			}
			if !found {
				continue
			}
		case enum:
			param = strings.ToUpper(param)
			_, ok := m.Enums[param]
			if ok {
				found = true
			}
		case lbl:
			if strings.HasPrefix(param, ":") {
				found = true
			}
		}
	}
	if !found {
		return ErrIllegalValue
	}
	return nil
}

func (m mnemonic) Generate(param string, prgCounter int, a *Assembler) byte {
	found := false
	for _, pt := range m.Param {
		var v byte
		switch pt {
		case int4:
			v, found, _ = convertNumber(param)
			if !found {
				continue
			}
			return byte(m.Code + v)
		case enum:
			param = strings.ToUpper(param)
			v, ok := m.Enums[param]
			if ok {
				found = true
			}
			switch param {
			case ":?":
				a.pageLabel = a.prgCounter
				return m.Code
			}
			if v > 0x0f {
				return byte(byte(v))
			}
			return byte(m.Code + byte(v))
		case lbl:
			switch m.Name {
			case "RJMP":
				param = strings.TrimPrefix(param, ":")
				lbl, ok := a.Labels[param]
				if ok {
					df := prgCounter - lbl.PrgCounter
					return byte(m.Code + byte(df))
				}
			case "CASB", "DFSB":
				id := a.subNumber(param)
				return byte(m.Code + byte(id+1))
			case "JMP", "LOOPC", "LOOPD", "CALL":
				param = strings.TrimPrefix(param, ":")
				lbl, ok := a.Labels[param]
				if ok {
					c := lbl.PrgCounter % 16
					p := lbl.PrgCounter / 16
					a.SetPage(p)
					return byte(m.Code + byte(c))
				}
			}
		}
	}
	return m.Code
}

func convertNumber(p string) (byte, bool, error) {
	if strings.HasPrefix(strings.ToLower(p), "#0x") {
		v, err := strconv.ParseUint(p[3:], 16, 4)
		if err != nil {
			return 0, false, ErrIllegalValue
		}
		return byte(v), true, nil
	}
	if strings.HasPrefix(strings.ToLower(p), "#0b") {
		v, err := strconv.ParseUint(p[3:], 2, 4)
		if err != nil {
			return 0, false, ErrIllegalValue
		}
		return byte(v), true, nil
	}
	if strings.HasPrefix(p, "#") {
		v, err := strconv.ParseUint(p[1:], 10, 4)
		if err != nil {
			return 0, false, ErrIllegalValue
		}
		return byte(v), true, nil
	}
	return 0, false, nil
}
