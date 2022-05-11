package asm

import (
	"errors"
	"fmt"
	"strings"
)

var ErrParamCount = errors.New("parameter count is wrong")

type mnemonic interface {
	Name() string
	CheckParameter(params []string) error
}

func GetMnemonic(name string) (mnemonic, error) {
	name = strings.ToUpper(name)
	mno, ok := Mnos[name]
	if !ok {
		return nil, fmt.Errorf("unknown mnemonics: %s ", name)
	}
	return mno, nil
}

var Mnos = map[string]mnemonic{
	"PORT": Port{},
}

type Port struct{}

func (p Port) Name() string {
	return "PORT"
}

func (p Port) CheckParameter(params []string) error {
	if len(params) != 1 {
		return ErrParamCount
	}
	return nil
}
