package asm

import (
	"strings"

	log "github.com/willie68/tps_asm/internal/logging"
)

func (a *Assembler) generate() {
	a.Binary = make([]byte, 0)

	for x, cmd := range a.Code {
		a.prgCounter = x
		a.parts = strings.Split(cmd, " ")
		a.command = a.parts[0]
		mno, err := GetMnemonic(a.command)
		if err != nil {
			a.addError(err)
		}
		a.Binary = append(a.Binary, mno.Generate(a.parts[1:], x, a))
		log.Infof("process command: %s", cmd)
	}
}
