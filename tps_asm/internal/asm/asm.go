package asm

import (
	"bufio"
	"fmt"
	"os"
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
	Source   string
	Labels   map[string]label
	Macros   map[string]macro
	actMacro macro
	Code     []string
}

func (a *Assembler) Parse() (errs []error) {
	now := time.Now()
	errs = a.parseOne()
	log.Infof("time to parse: %d", time.Until(now).Milliseconds())
	return
}

func (a *Assembler) parseOne() (errs []error) {
	a.Labels = make(map[string]label)
	a.Macros = make(map[string]macro)
	a.Code = make([]string, 0)

	inMacro := false
	inComment := false

	file, err := os.Open(a.Source)
	//handle errors while opening
	if err != nil {
		errs = append(errs, fmt.Errorf("error when opening file: %v", err))
	}
	defer file.Close()

	fileScanner := bufio.NewScanner(file)
	lineNumber := 0
	prgCounter := 1
	// read line by line
	log.Info("----- start -----")
	for fileScanner.Scan() {
		line := fileScanner.Text()
		lineNumber++
		line = strings.TrimSpace(line)
		// ignoring empty lines
		if line == "" {
			continue
		}

		// Commant processing
		if strings.HasPrefix(line, "*/") {
			if !inComment {
				errs = append(errs, fmt.Errorf("line %d: missing starting /* comment directrive", lineNumber))
			}
			inComment = false
			continue
		}
		if inComment {
			continue
		}
		if strings.HasPrefix(line, "/*") {
			if inComment {
				errs = append(errs, fmt.Errorf("line %d: already in comment, nested comments are not supported", lineNumber))
			}
			inComment = true
			continue
		}

		// label processing
		line = a.removeComment(line)
		if strings.HasPrefix(line, ":") {
			parts := strings.Split(line, " ")
			labelName := strings.ToLower(parts[0])
			a.Labels[labelName] = label{
				Name:       labelName,
				PrgCounter: prgCounter,
			}
			log.Debugf("define label: %s", labelName)
			continue
		}

		// Splitting into parts
		parts := strings.Split(line, " ")
		command := parts[0]

		// include files
		if command == ".include" {
			errs = append(errs, fmt.Errorf("line %d: .include is not implemented, ignoring", lineNumber))
			continue
		}

		// macro definition
		if command == ".endmacro" {
			if !inMacro {
				errs = append(errs, fmt.Errorf("line %d: missing starting .macro directrive", lineNumber))
			}
			a.Macros[a.actMacro.Name] = a.actMacro
			inMacro = false
			continue
		}
		if inMacro {
			a.actMacro.Code = append(a.actMacro.Code, line)
			continue
		}
		if command == ".macro" {
			if inMacro {
				errs = append(errs, fmt.Errorf("line %d: already in macro definition, nested macros are not supported", lineNumber))
			}
			macroName := strings.ToLower(parts[1])
			a.actMacro = macro{
				Name:   macroName,
				Params: parts[2:],
			}
			inMacro = true
			continue
		}

		// macro processing
		if strings.HasPrefix(command, ".") {
			macroName := strings.ToLower(parts[0][1:])
			macro, ok := a.Macros[macroName]
			if !ok {
				errs = append(errs, fmt.Errorf("can't find macro: %s", macroName))
			} else {
				log.Infof("use macro: %s", macro.Name)
				for _, cmd := range macro.Code {
					prgCounter++
					a.Code = append(a.Code, cmd)
					log.Debugf("line %d: %s", lineNumber, cmd)
				}
				continue
			}
		}

		prgCounter++
		a.Code = append(a.Code, line)
		log.Debugf("line %d: %s", lineNumber, line)
	}
	log.Info("----- stop -----")
	// handle first encountered error while reading
	if err := fileScanner.Err(); err != nil {
		errs = append(errs, fmt.Errorf("line %d: error while reading file: %v", lineNumber, err))
	}

	return
}

func (a *Assembler) removeComment(line string) string {
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
