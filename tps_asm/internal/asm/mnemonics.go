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
)

var Mnos = []mnemonic{
	{
		Name:   "NOP",
		Params: [][]string{}, // int4: full 4 bit variable, label : goto to label
		Code:   0x00,
	},
	{
		Name:   "PORT",
		Params: [][]string{{int4}}, // int4: full 4 bit variable, label : goto to label
		Code:   0x10,
	},
	{
		Name:   "WAIT",
		Params: [][]string{{int4, enum}}, // int4: full 4 bit variable, label : goto to label
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
		Params: [][]string{{int4, label}}, // int4: full 4 bit variable, label : goto to label
		Code:   0x30,
	},
	{
		Name:   "LDA",
		Params: [][]string{{int4}}, // int4: full 4 bit variable, label : goto to label
		Code:   0x40,
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
				found, err = convertNumber(p)
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
			}
		}
	}
	if !found {
		return ErrIllegalValue
	}
	return nil
}

func convertNumber(p string) (bool, error) {
	if strings.HasPrefix(p, "#0x") {
		_, err := strconv.ParseUint(p[3:], 16, 4)
		if err != nil {
			return false, ErrIllegalValue
		}
		return true, nil
	}
	if strings.HasPrefix(p, "#0b") {
		_, err := strconv.ParseUint(p[3:], 2, 4)
		if err != nil {
			return false, ErrIllegalValue
		}
		return true, nil
	}
	if strings.HasPrefix(p, "#") {
		_, err := strconv.ParseUint(p[1:], 10, 4)
		if err != nil {
			return false, ErrIllegalValue
		}
		return true, nil
	}
	return false, nil
}
