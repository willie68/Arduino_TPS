package asm

import (
	"testing"

	"github.com/stretchr/testify/assert"
	log "github.com/willie68/tps_asm/internal/logging"
)

func TestPrgSize(t *testing.T) {
	var td = []struct {
		h  Hardware
		sc int
		ec int
	}{
		{
			h:  Holtek,
			sc: 127,
			ec: 0,
		}, {
			h:  Holtek,
			sc: 129,
			ec: 1,
		}, {
			h:  ATMega8,
			sc: 256,
			ec: 0,
		}, {
			h:  ATMega8,
			sc: 257,
			ec: 1,
		}, {
			h:  ArduinoSPS,
			sc: 1024,
			ec: 0,
		}, {
			h:  ArduinoSPS,
			sc: 1025,
			ec: 1,
		}, {
			h:  TinySPS,
			sc: 512,
			ec: 0,
		}, {
			h:  TinySPS,
			sc: 513,
			ec: 1,
		},
	}
	log.Logger.SetLevel(log.LvError)
	ast := assert.New(t)
	for _, t := range td {
		tasm := Assembler{
			Hardware: t.h,
			Source:   generateSrc(t.sc),
		}
		tasm.Parse()
		ast.Equal(t.ec, len(tasm.errs), "error count not equal")
	}
}

func generateSrc(d int) []string {
	src := make([]string, 0)
	for i := 0; i < d; i++ {
		src = append(src, "NOP")
	}
	return src
}
