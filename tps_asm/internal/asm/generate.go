package asm

import (
	"fmt"
	"strings"

	log "github.com/willie68/tps_asm/internal/logging"
)

func (a *Assembler) generate() {
	a.Binary = make([]byte, 0)

	for x, cmd := range a.Code {
		a.prgCounter = x
		a.parts = strings.Split(cmd, " ")
		if len(a.parts) < 2 {
			a.parts = append(a.parts, "")
		}
		a.command = a.parts[0]
		mno, err := GetMnemonic(a.command)
		if err != nil {
			a.addError(err)
			continue
		}
		if mno == nil {
			a.addError(fmt.Errorf("unknown mnemonic: %s", cmd))
			continue
		}
		a.Binary = append(a.Binary, mno.Generate(a.parts[1], x, a))
		log.Infof("process command: %s", cmd)
	}
}

func (a *Assembler) SetPage(p int) {
	//TODO now every system has a 4 bit page register
	if (a.pageLabel >= 0) && (p < 16) {
		a.Binary[a.pageLabel] = a.Binary[a.pageLabel] + byte(p)
		a.pageLabel = -1
	}
}
