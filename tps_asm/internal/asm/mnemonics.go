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
	Name   string
	Params []([]string)
	Enums  map[string]int
	Code   byte
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
		Name:   "NOP",
		Params: [][]string{},
		Code:   0x00,
	},
	{
		Name:   "WAIT",
		Params: [][]string{{int4, enum}},
		Enums: map[string]int{
			"1ms":   0x00,
			"2ms":   0x01,
			"5ms":   0x02,
			"10ms":  0x03,
			"20ms":  0x04,
			"50ms":  0x05,
			"100ms": 0x06,
			"200ms": 0x07,
			"500ms": 0x08,
			"1s":    0x09,
			"2s":    0x0A,
			"5s":    0x0B,
			"10s":   0x0C,
			"20s":   0x0D,
			"30s":   0x0E,
			"60s":   0x0F,
		},
		Code: 0x20,
	},
	{
		Name:   "RJMP",
		Params: [][]string{{int4, lbl}},
		Code:   0x30,
	},
	{
		Name:   "PAGE",
		Params: [][]string{{int4, enum}},
		Enums: map[string]int{
			":?": 0x10,
		},
		Code: 0x80,
	},
	{
		Name:   "JMP",
		Params: [][]string{{int4, lbl}},
		Code:   0x90,
	},
	{
		Name:   "LOOPC",
		Params: [][]string{{int4, lbl}},
		Code:   0xA0,
	},
	{
		Name:   "LOOPD",
		Params: [][]string{{int4, lbl}},
		Code:   0xB0,
	},
	{
		Name:   "SKIP0",
		Params: [][]string{},
		Code:   0xC0,
	},
	{
		Name:   "AGTB",
		Params: [][]string{},
		Code:   0xC1,
	},
	{
		Name:   "ALTB",
		Params: [][]string{},
		Code:   0xC2,
	},
	{
		Name:   "AEQB",
		Params: [][]string{},
		Code:   0xC3,
	},
	{
		Name:   "DEQ0",
		Params: [][]string{{enum}},
		Enums: map[string]int{
			"1": 0x08,
			"2": 0x09,
			"3": 0x0A,
			"4": 0x0B,
		},
		Code: 0xC0,
	},
	{
		Name:   "DEQ1",
		Params: [][]string{{enum}},
		Enums: map[string]int{
			"1": 0x04,
			"2": 0x05,
			"3": 0x06,
			"4": 0x07,
		},
		Code: 0xC0,
	},
	{
		Name:   "PRG0",
		Params: [][]string{},
		Code:   0xCC,
	},
	{
		Name:   "SEL0",
		Params: [][]string{},
		Code:   0xCD,
	},
	{
		Name:   "PRG1",
		Params: [][]string{},
		Code:   0xCE,
	},
	{
		Name:   "SEL1",
		Params: [][]string{},
		Code:   0xCF,
	},
	{
		Name:   "CALL",
		Params: [][]string{{int4, lbl}},
		Code:   0xD0,
	},
	{
		Name:   "RTR",
		Params: [][]string{},
		Code:   0xE0,
	},
	{
		Name:   "CASB",
		Params: [][]string{{int4, lbl}},
		Code:   0xE0,
	},
	{
		Name:   "DFSB",
		Params: [][]string{{int4, lbl}},
		Code:   0xE7,
	},
	{
		Name:   "REST",
		Params: [][]string{},
		Code:   0xEF,
	},
	{
		Name:   "PEND",
		Params: [][]string{},
		Code:   0xFF,
	},

	// Load and save
	{
		Name:   "LDA",
		Params: [][]string{{int4, enum}},
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
	},
	{
		Name:   "SWAP",
		Params: [][]string{},
		Code:   0x50,
	},
	{
		Name:   "MOV",
		Params: [][]string{{enum}},
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
	},
	{
		Name:   "PUSH",
		Params: [][]string{},
		Code:   0x5F,
	},
	{
		Name:   "POP",
		Params: [][]string{},
		Code:   0x6F,
	},
	// Math
	{
		Name:   "INC",
		Params: [][]string{},
		Code:   0x71,
	},
	{
		Name:   "DEC",
		Params: [][]string{},
		Code:   0x72,
	},
	{
		Name:   "ADD",
		Params: [][]string{},
		Code:   0x73,
	},
	{
		Name:   "SUB",
		Params: [][]string{},
		Code:   0x74,
	},
	{
		Name:   "MUL",
		Params: [][]string{},
		Code:   0x75,
	},
	{
		Name:   "DIV",
		Params: [][]string{},
		Code:   0x76,
	},
	{
		Name:   "AND",
		Params: [][]string{},
		Code:   0x77,
	},
	{
		Name:   "OR",
		Params: [][]string{},
		Code:   0x78,
	},
	{
		Name:   "XOR",
		Params: [][]string{},
		Code:   0x79,
	},
	{
		Name:   "NOT",
		Params: [][]string{},
		Code:   0x7A,
	},
	{
		Name:   "MOD",
		Params: [][]string{},
		Code:   0x7B,
	},
	{
		Name:   "BYTE",
		Params: [][]string{},
		Code:   0x7C,
	},
	{
		Name:   "BSUBA",
		Params: [][]string{},
		Code:   0x7D,
	},
	{
		Name:   "SHR",
		Params: [][]string{},
		Code:   0x7E,
	},
	{
		Name:   "SHL",
		Params: [][]string{},
		Code:   0x7F,
	},
	// Input/Output
	{
		Name:   "PORT",
		Params: [][]string{{int4}},
		Code:   0x10,
	},
	{
		Name:   "STA",
		Params: [][]string{{enum}},
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
	},
	// Byte mnemonics
	{
		Name:   "BLDA",
		Params: [][]string{{enum}},
		Enums: map[string]int{
			"ADC1": 0xF0,
			"ADC2": 0xF1,
			"RC1":  0xF2,
			"RC2":  0xF3,
		},
		Code: 0x00,
	},
	{
		Name:   "BSTA",
		Params: [][]string{{enum}},
		Enums: map[string]int{
			"PWM1": 0xF4,
			"PWM2": 0xF5,
			"SRV1": 0xF6,
			"SRV2": 0xF7,
		},
		Code: 0x00,
	},
	{
		Name:   "TONE",
		Params: [][]string{},
		Code:   0xF8,
	},
}

func (m mnemonic) CheckParameter(params []string) error {
	// Check parameter count
	if len(params) != len(m.Params) {
		return ErrParamCount
	}
	if len(m.Params) == 0 {
		return nil
	}
	var err error
	found := false
	for x, pts := range m.Params {
		p := params[x]
	ptsloop:
		for _, pt := range pts {
			switch pt {
			case int4:
				_, found, err = convertNumber(p)
				if err != nil {
					return err
				}
				if !found {
					continue
				}
				break ptsloop
			case enum:
				_, ok := m.Enums[p]
				if ok {
					found = true
				}
				break ptsloop
			case lbl:
				if strings.HasPrefix(p, ":") {
					found = true
				}
				break ptsloop
			}
		}
	}
	if !found {
		return ErrIllegalValue
	}
	return nil
}

func (m mnemonic) Generate(params []string, prgCounter int, a *Assembler) byte {
	found := false
	for x, pts := range m.Params {
		p := params[x]
	ptsloop:
		for _, pt := range pts {
			var v byte
			switch pt {
			case int4:
				v, found, _ = convertNumber(p)
				if !found {
					continue
				}
				return byte(m.Code + v)
			case enum:
				v, ok := m.Enums[p]
				if ok {
					found = true
				}
				if v > 0x0f {
					return byte(byte(v))
				}
				return byte(m.Code + byte(v))
			case lbl:
				switch m.Name {
				case "RJMP":
					p = strings.TrimPrefix(p, ":")
					lbl, ok := a.Labels[p]
					if ok {
						df := prgCounter - lbl.PrgCounter
						return byte(m.Code + byte(df))
					}
				case "CASB", "DFSB":
					id := a.subNumber(p)
					return byte(m.Code + byte(id+1))
				}
				break ptsloop
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
