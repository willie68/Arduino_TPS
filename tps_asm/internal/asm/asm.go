package asm

import (
	"fmt"
	"time"

	log "github.com/willie68/tps_asm/internal/logging"
)

type label struct {
	Name       string
	PrgCounter int
}

type macro struct {
	Name   string
	Params []string
	Code   []string
}

type Assembler struct {
	Hardware Hardware
	Source   []string
	Includes string
	Labels   map[string]label
	Macros   map[string]macro
	Code     []string
	Binary   []byte

	lineNumber int
	prgCounter int
	actMacro   macro
	inMacro    bool
	inComment  bool
	parts      []string
	command    string
	errs       []error
	line       string
}

func (a *Assembler) Parse() []error {
	start := time.Now()
	a.parse()
	log.Infof("time to parse: %d ms", -time.Until(start).Milliseconds())
	start = time.Now()
	a.generate()
	log.Infof("time to generate: %d ms", -time.Until(start).Milliseconds())
	return a.errs
}

func (a *Assembler) addErrorS(msg string) {
	a.errs = append(a.errs, fmt.Errorf("line %d: %s", a.lineNumber, msg))
}

func (a *Assembler) addError(err error) {
	a.errs = append(a.errs, err)
}
