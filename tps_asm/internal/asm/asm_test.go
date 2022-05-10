package asm

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/imdario/mergo"
)

var testdatas = []struct {
	name   string
	source []string
	asm    Assembler
}{{
	name: "testlabel1",
	asm: Assembler{
		Source: []string{":label1", ":label2"},
		Labels: map[string]label{
			".label1": label{
				Name:       ".label1",
				PrgCounter: 1,
			},
			".label2": label{
				Name:       ".label2",
				PrgCounter: 1,
			},
		},
	},
},
}

func TestLabels(t *testing.T) {
	for _, test := range testdatas {
		fmt.Printf("testing %s\r\n", test.name)

		tasm := Assembler{
			Hardware: Holtek,
			Source:   test.asm.Source,
		}
		tasm.parseOne()
		if len(tasm.errs) > 0 {
			t.Fail()
		}

		src_jstr, err := json.Marshal(tasm)
		if err != nil {
			t.Fail()
		}
		tst_jstr, err := json.Marshal(test.asm)
		if err != nil {
			t.Fail()
		}
		t.Log("org:", string(src_jstr))
		t.Log("tst:", string(tst_jstr))

		mergo.Map(&tasm, test.asm, mergo.WithOverride)
		dsc_jstr, _ := json.Marshal(tasm)

		t.Log("dst:", string(dsc_jstr))

		if string(src_jstr) != string(dsc_jstr) {
			t.Fail()
		}
	}
}
