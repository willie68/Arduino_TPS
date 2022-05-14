package asm

import (
	"fmt"
	"time"

	log "github.com/willie68/tps_asm/internal/logging"
)

type label struct {
	Name       string
	PrgCounter int
	Line       int
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
	Subs     []string
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
	pageLabel  int
}

func (a *Assembler) Parse() []error {
	start := time.Now()
	a.init()
	a.parse()
	log.Infof("time to parse: %d ms", -time.Until(start).Milliseconds())
	start = time.Now()
	a.generate()
	log.Infof("time to generate: %d ms", -time.Until(start).Milliseconds())
	return a.errs
}

func (a *Assembler) init() {
	a.Labels = make(map[string]label)
	a.Subs = make([]string, 0)
	a.Macros = make(map[string]macro)
	a.Code = make([]string, 0)

	a.inMacro = false
	a.inComment = false

	a.lineNumber = 0
	a.prgCounter = 0
	a.pageLabel = -1
}

func (a *Assembler) addErrorS(msg string) {
	a.errs = append(a.errs, fmt.Errorf("line %d: %s", a.lineNumber, msg))
}

func (a *Assembler) addError(err error) {
	a.errs = append(a.errs, err)
}

func (a *Assembler) subNumber(subname string) int {
	for x, v := range a.Subs {
		if subname == v {
			return x
		}
	}
	return -1 //not found.
}
