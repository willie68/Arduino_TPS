package asm

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	log "github.com/willie68/tps_asm/internal/logging"
)

var testdatas = []struct {
	name     string
	source   []string
	asm      Assembler
	errCount int
}{{
	name: "testEmptyLine",
	asm: Assembler{
		Source: []string{" ", " ", " ", "LDA #12", " "},
		Labels: map[string]label{},
		Macros: map[string]macro{},
		Code:   []string{"LDA #12"},
	},
	errCount: 0,
}, {
	name: "testTrimming1",
	asm: Assembler{
		Source: []string{"LDA #12", " LDA #12", "LDA #12 ", "\t LDA #12", "\t LDA #12  \t"},
		Labels: map[string]label{},
		Code:   []string{"LDA #12", "LDA #12", "LDA #12", "LDA #12", "LDA #12"},
	},
	errCount: 0,
}, {
	name: "testComment1",
	asm: Assembler{
		Source: []string{"/* dies ist ein Test", "1", "2", "3", "4", "*/"},
		Labels: map[string]label{},
		Code:   []string{},
	},
	errCount: 0,
}, {
	name: "testComment2",
	asm: Assembler{
		Source: []string{"/* dies ist ein Test", "1", "2", "*/", "LDA #12", "/*", "3", "4", "*/"},
		Labels: map[string]label{},
		Code:   []string{"LDA #12"},
	},
	errCount: 0,
}, {
	name: "testErrorComment1",
	asm: Assembler{
		Source: []string{"*/"},
		Labels: map[string]label{},
		Code:   []string{},
	},
	errCount: 1,
}, {
	name: "testErrorComment2",
	asm: Assembler{
		Source: []string{"/*", "/*", "*/"},
		Labels: map[string]label{},
		Code:   []string{},
	},
	errCount: 1,
}, {
	name: "testLineComment1",
	asm: Assembler{
		Source: []string{"; dies ist ein Test"},
		Labels: map[string]label{},
		Code:   []string{},
	},
	errCount: 0,
}, {
	name: "testLineComment2",
	asm: Assembler{
		Source: []string{" ; dies ist ein Test"},
		Labels: map[string]label{},
		Code:   []string{},
	},
	errCount: 0,
}, {
	name: "testInlineComment1",
	asm: Assembler{
		Source: []string{"LDA #12; dies ist ein Test"},
		Labels: map[string]label{},
		Code:   []string{"LDA #12"},
	},
	errCount: 0,
}, {
	name: "testInlineComment2",
	asm: Assembler{
		Source: []string{"LDA #12 ; dies ist ein Test"},
		Labels: map[string]label{},
		Code:   []string{"LDA #12"},
	},
	errCount: 0,
}, {
	name: "testlabel1",
	asm: Assembler{
		Source: []string{":label1", ":label2"},
		Labels: map[string]label{
			"label1": label{
				Name:       "label1",
				PrgCounter: 1,
			},
			"label2": label{
				Name:       "label2",
				PrgCounter: 1,
			},
		},
		Code: []string{},
	},
	errCount: 0,
}, {
	name: "testMacro1",
	asm: Assembler{
		Source: []string{":label1", ":label2", ".macro define", "LDA #12", ".endmacro"},
		Labels: map[string]label{
			"label1": label{
				Name:       "label1",
				PrgCounter: 1,
			},
			"label2": label{
				Name:       "label2",
				PrgCounter: 1,
			},
		},
		Macros: map[string]macro{
			"define": {
				Name:   "define",
				Params: []string{},
				Code:   []string{"LDA #12"},
			},
		},
		Code: []string{},
	},
	errCount: 0,
}, {
	name: "testMacroWithParameter",
	asm: Assembler{
		Source: []string{":label1", ":label2", ".macro define muck puck", "LDA #12", ".endmacro"},
		Labels: map[string]label{
			"label1": label{
				Name:       "label1",
				PrgCounter: 1,
			},
			"label2": label{
				Name:       "label2",
				PrgCounter: 1,
			},
		},
		Macros: map[string]macro{
			"define": {
				Name:   "define",
				Params: []string{"muck", "puck"},
				Code:   []string{"LDA #12"},
			},
		},
		Code: []string{},
	},
	errCount: 0,
}, {
	name: "testMacroError1",
	asm: Assembler{
		Source: []string{":label1", ":label2", ".macro", ".endmacro"},
		Labels: map[string]label{
			"label1": label{
				Name:       "label1",
				PrgCounter: 1,
			},
			"label2": label{
				Name:       "label2",
				PrgCounter: 1,
			},
		},
		Code: []string{},
	},
	errCount: 2,
}, {
	name: "testMacroError2",
	asm: Assembler{
		Source: []string{".macro define1", ".macro define2", ".endmacro"},
		Labels: map[string]label{},
		Macros: map[string]macro{
			"define1": {
				Name: "define1",
			},
		},
		Code: []string{},
	},
	errCount: 1,
}, {
	name: "testMacroProcessing1",
	asm: Assembler{
		Source: []string{".macro define1", "LDA #12", "STA #12", ".endmacro", ".define1"},
		Labels: map[string]label{},
		Macros: map[string]macro{
			"define1": {
				Name: "define1",
				Code: []string{"LDA #12", "STA #12"},
			},
		},
		Code: []string{"LDA #12", "STA #12"},
	},
	errCount: 0,
}, {
	name: "testMacroProcessing2",
	asm: Assembler{
		Source: []string{".macro define1 time name", "LDA time", "STA name", ".endmacro", ".define1 #12 #23"},
		Labels: map[string]label{},
		Macros: map[string]macro{
			"define1": {
				Name:   "define1",
				Params: []string{"time", "name"},
				Code:   []string{"LDA time", "STA name"},
			},
		},
		Code: []string{"LDA #12", "STA #23"},
	},
	errCount: 0,
}, {
	name: "testMacroProcessingError1",
	asm: Assembler{
		Source: []string{".macro define1 time name", "LDA time", "STA name", ".endmacro", ".define2 #12 #23"},
		Labels: map[string]label{},
		Macros: map[string]macro{
			"define1": {
				Name:   "define1",
				Params: []string{"time", "name"},
				Code:   []string{"LDA time", "STA name"},
			},
		},
		Code: []string{},
	},
	errCount: 1,
},
}

func TestOne(t *testing.T) {
	log.Logger.SetLevel(log.LvError)
	oneTest(t, "testInlineComment1")
}

func TestAll(t *testing.T) {
	log.Logger.SetLevel(log.LvError)
	for _, test := range testdatas {
		oneTest(t, test.name)
	}
}

func oneTest(t *testing.T, name string) {
	ast := assert.New(t)
	for _, test := range testdatas {
		if name == test.name {

			fmt.Printf("testing %s\r\n", test.name)

			tasm := Assembler{
				Hardware: Holtek,
				Source:   test.asm.Source,
			}
			tasm.parse()
			ast.Equal(test.errCount, len(tasm.errs), "error count not equal")
			if len(tasm.errs) != test.errCount {
				for _, err := range tasm.errs {
					t.Logf("error: %v", err)
				}
				t.Fail()
			}
			checkAssembler(ast, test.asm, tasm)
		}
	}
}

func checkAssembler(ast *assert.Assertions, expected, actual Assembler) {
	// Testing hardware
	ast.Equal(expected.Hardware, actual.Hardware, "hardware not equal, exp: %v , act: %v", expected.Hardware, actual.Hardware)
	// Testing source
	ast.Equal(expected.Source, actual.Source, "generated code is not equal")
	// Testing labels
	ast.Equal(fmt.Sprint(expected.Labels), fmt.Sprint(actual.Labels), "Labels are not equal")
	// Testing macros
	ast.Equal(fmt.Sprint(expected.Macros), fmt.Sprint(actual.Macros), "Macros are not equal")
	// Testing code
	ast.Equal(expected.Code, actual.Code, "Code is not equal")
}
