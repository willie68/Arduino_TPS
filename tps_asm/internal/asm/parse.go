package asm

import (
	"fmt"
	"strings"

	log "github.com/willie68/tps_asm/internal/logging"
	"github.com/willie68/tps_asm/internal/utils"
)

func (a *Assembler) parse() {
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

		a.line = a.removeInlineComment(a.line)
		if a.line == "" {
			continue
		}
		a.parts = strings.Split(a.line, " ")
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

		a.processSubroutines()

		a.parts = a.checkSyntax()

		a.prgCounter++
		a.Code = append(a.Code, strings.Join(a.parts, " "))
		log.Debugf("line %d: %s", a.lineNumber, line)
	}
	log.Info("----- stop -----")
}

func (a *Assembler) removeInlineComment(line string) string {
	if strings.ContainsAny(line, ";") {
		pos := strings.Index(line, ";")
		newline := utils.Substr(line, 0, pos)
		return strings.TrimSpace(newline)
	}
	return line
}

func (a *Assembler) processMacroDefinition() bool {
	if a.command == ".endmacro" {
		if !a.inMacro {
			a.addErrorS("missing starting .macro directrive")
			return true
		}
		a.Macros[a.actMacro.Name] = a.actMacro
		a.inMacro = false
		return true
	}
	if a.command == ".macro" {
		if a.inMacro {
			a.addErrorS("already in macro definition, nested macros are not supported")
			return true
		}
		if len(a.parts) < 2 {
			a.addErrorS("missing macroname")
			return true
		}
		macroName := strings.ToLower(a.parts[1])
		a.actMacro = macro{
			Name:   macroName,
			Params: a.parts[2:],
		}
		a.inMacro = true
		return true
	}
	if a.inMacro {
		a.actMacro.Code = append(a.actMacro.Code, a.line)
		return true
	}
	return false
}

func (a *Assembler) processLabelDefinition() bool {
	if strings.HasPrefix(a.line, ":") {
		labelName := strings.TrimPrefix(strings.ToLower(a.parts[0]), ":")
		a.Labels[labelName] = label{
			Name:       labelName,
			PrgCounter: a.prgCounter,
			Line:       a.lineNumber,
		}
		log.Debugf("define label: %s", labelName)
		return true
	}
	return false
}

func (a *Assembler) processSubroutines() bool {
	if strings.ToUpper(a.command) == "DFSB" {
		subName := a.parts[1]
		a.Subs = append(a.Subs, subName)
		log.Debugf("define subroutine: %s", subName)
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
			return true
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
	if strings.HasPrefix(a.line, ";") {
		return true
	}
	if strings.HasPrefix(a.line, "*/") {
		if !a.inComment {
			a.addErrorS("missing starting /* comment directrive")
		}
		a.inComment = false
		return true
	}
	if strings.HasPrefix(a.line, "/*") {
		if a.inComment {
			a.addErrorS("already in comment, nested comments are not supported")
		}
		a.inComment = true
		return true
	}
	if a.inComment {
		return true
	}
	return false
}

func (a *Assembler) checkSyntax() []string {
	mno, err := GetMnemonic(a.command)
	if err != nil {
		a.addError(err)
		return a.parts
	}
	var part string
	for _, p := range a.parts[1:] {
		part += p
	}
	err = mno.CheckParameter(part)
	if err != nil {
		a.addError(err)
		return a.parts
	}
	if len(mno.Param) == 0 {
		return []string{a.command}
	}
	return []string{a.command, part}
}
