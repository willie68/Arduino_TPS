package asm

import log "github.com/willie68/tps_asm/internal/logging"

func (a *Assembler) generate() {
	a.Binary = make([]byte, 0)

	for _, cmd := range a.Code {
		log.Debugf("process command: %s", cmd)
	}
}
