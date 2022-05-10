package asm

import (
	"fmt"
	"strings"
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

	lineNumber int
	prgCounter int
	actMacro   macro
	Code       []string
	inMacro    bool
	inComment  bool
	parts      []string
	command    string
	errs       []error
	line       string
}

func (a *Assembler) Parse() []error {
	start := time.Now()
	a.parseOne()
	log.Infof("time to parse: %d ms", -time.Until(start).Milliseconds())
	return a.errs
}

func (a *Assembler) parseOne() {
	a.Labels = make(map[string]label)
	a.Macros = make(map[string]macro)
	a.Code = make([]string, 0)

	a.inMacro = false
	a.inComment = false

	a.lineNumber = 0
	a.prgCounter = 1
	// read line by line
	log.Info("----- start -----")
	for x, line := range a.Source {
		a.lineNumber = x
		line = strings.TrimSpace(line)

		// ignoring empty lines
		if line == "" {
			continue
		}

		a.line = line

		// Commant processing
		if a.processComment() {
			continue
		}

		a.line = a.removeLineComment(a.line)
		a.parts = strings.Split(line, " ")
		a.command = a.parts[0]

		// label processing
		if a.processLabelDefinition() {
			continue
		}

		// Splitting into parts

		command := a.command

		// include files
		if command == ".include" {
			a.addErrorS(".include is not implemented, ignoring")
			continue
		}

		// macro definition
		if a.processMacroDefinition() {
			continue
		}

		// macro processing
		if a.processMacro() {
			continue
		}

		a.prgCounter++
		a.Code = append(a.Code, line)
		log.Debugf("line %d: %s", a.lineNumber, line)
	}
	log.Info("----- stop -----")
}

func (a *Assembler) removeLineComment(line string) string {
	if strings.ContainsAny(line, ";") {
		pos := strings.Index(line, ";")
		return substr(line, 0, pos-1)
	}
	return line
}

// NOTE: this isn't multi-Unicode-codepoint aware, like specifying skintone or
//       gender of an emoji: https://unicode.org/emoji/charts/full-emoji-modifiers.html
func substr(input string, start int, length int) string {
	asRunes := []rune(input)

	if start >= len(asRunes) {
		return ""
	}

	if start+length > len(asRunes) {
		length = len(asRunes) - start
	}

	return string(asRunes[start : start+length])
}

func (a *Assembler) processMacroDefinition() bool {
	if a.command == ".endmacro" {
		if !a.inMacro {
			a.addErrorS("missing starting .macro directrive")
		}
		a.Macros[a.actMacro.Name] = a.actMacro
		a.inMacro = false
		return true
	}
	if a.inMacro {
		a.actMacro.Code = append(a.actMacro.Code, a.line)
		return true
	}
	if a.command == ".macro" {
		if a.inMacro {
			a.addErrorS("already in macro definition, nested macros are not supported")
		}
		macroName := strings.ToLower(a.parts[1])
		a.actMacro = macro{
			Name:   macroName,
			Params: a.parts[2:],
		}
		a.inMacro = true
		return true
	}
	return false
}

func (a *Assembler) processLabelDefinition() bool {
	if strings.HasPrefix(a.line, ":") {
		labelName := strings.ToLower(a.parts[0])
		a.Labels[labelName] = label{
			Name:       labelName,
			PrgCounter: a.prgCounter,
		}
		log.Debugf("define label: %s", labelName)
		return true
	}
	return false
}

func (a *Assembler) processMacro() bool {
	if strings.HasPrefix(a.command, ".") {
		macroName := strings.ToLower(a.parts[0][1:])
		macro, ok := a.Macros[macroName]
		if !ok {
			a.addErrorS(fmt.Sprintf("can't find macro: %s", macroName))
		} else {
			log.Infof("use macro: %s", macro.Name)
			for _, cmd := range macro.Code {
				a.prgCounter++
				for x, mac := range macro.Params {
					cmd = strings.ReplaceAll(cmd, mac, a.parts[x+1])
				}
				a.Code = append(a.Code, cmd)
				log.Debugf("line %d: %s", a.lineNumber, cmd)
			}
			return true
		}
	}
	return false
}

func (a *Assembler) processComment() bool {
	if strings.HasPrefix(a.line, "*/") {
		if !a.inComment {
			a.addErrorS("missing starting /* comment directrive")
		}
		a.inComment = false
		return true
	}
	if a.inComment {
		return true
	}
	if strings.HasPrefix(a.line, "/*") {
		if a.inComment {
			a.addErrorS("already in comment, nested comments are not supported")
		}
		a.inComment = true
		return true
	}
	return false
}

func (a *Assembler) addErrorS(msg string) {
	a.errs = append(a.errs, fmt.Errorf("line %d: %s", a.lineNumber, msg))
}

func (a *Assembler) addError(err error) {
	a.errs = append(a.errs, err)
}
