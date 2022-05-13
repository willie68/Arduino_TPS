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
				Line:       0,
			},
			"label2": label{
				Name:       "label2",
				PrgCounter: 1,
				Line:       1,
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
				Line:       0,
			},
			"label2": label{
				Name:       "label2",
				PrgCounter: 1,
				Line:       1,
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
				Line:       0,
			},
			"label2": label{
				Name:       "label2",
				PrgCounter: 1,
				Line:       1,
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
				Line:       0,
			},
			"label2": label{
				Name:       "label2",
				PrgCounter: 1,
				Line:       1,
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
}, {
	name: "testProgramm1",
	asm: Assembler{
		Source: []string{
			".macro blink",
			"PORT #0B0101",
			"WAIT 200ms",
			"PORT #0B1010",
			"WAIT 200ms",
			".endmacro",
			" ",
			":loop",
			".blink",
			"RJMP :loop",
			"/*",
			"Kommentar über mehrere Zeilen",
			"*/",
			" ",
			".macro macro1 output time",
			"PORT output",
			"WAIT time",
			"PORT #0x00",
			"WAIT time",
			".endmacro",
			"",
			":loop1",
			".macro1 #0x0f 100ms",
			" ",
			"PORT #0x0F ;Zeilenkommentar",
			"WAIT 200ms",
			"PORT #0x00",
			"WAIT 200ms",
			"CASB :sub1",
			"RJMP :loop1",
			" ",
			"DFSB :sub1",
			"PORT #0x0F ;Zeilenkommentar",
			"WAIT 200ms",
			"PORT #0x00",
			"WAIT 200ms",
			"RTR",
		},
		Labels: map[string]label{
			"loop": {
				Name:       "loop",
				PrgCounter: 0,
				Line:       7,
			},
			"loop1": {
				Name:       "loop1",
				PrgCounter: 5,
				Line:       21,
			},
		},
		Macros: map[string]macro{
			"blink": {
				Name:   "blink",
				Params: []string{},
				Code: []string{
					"PORT #0B0101",
					"WAIT 200ms",
					"PORT #0B1010",
					"WAIT 200ms",
				},
			},
			"macro1": {
				Name:   "macro1",
				Params: []string{"output", "time"},
				Code: []string{
					"PORT output",
					"WAIT time",
					"PORT #0x00",
					"WAIT time",
				},
			},
		},
		Code: []string{
			"PORT #0B0101",
			"WAIT 200ms",
			"PORT #0B1010",
			"WAIT 200ms",
			"RJMP :loop",
			"PORT #0x0f",
			"WAIT 100ms",
			"PORT #0x00",
			"WAIT 100ms",
			"PORT #0x0F",
			"WAIT 200ms",
			"PORT #0x00",
			"WAIT 200ms",
			"CASB :sub1",
			"RJMP :loop1",
			"DFSB :sub1",
			"PORT #0x0F",
			"WAIT 200ms",
			"PORT #0x00",
			"WAIT 200ms",
			"RTR",
		},
		Binary: []byte{
			0x15,
			0x27,
			0x1A,
			0x27,
			0x34,
			0x1F,
			0x26,
			0x10,
			0x26,
			0x1F,
			0x27,
			0x10,
			0x27,
			0xE1,
			0x39,
			0xE8,
			0x1F,
			0x27,
			0x10,
			0x27,
			0xE0,
		},
	},
	errCount: 0,
}, {
	name: "testProgramm2",
	asm: Assembler{
		Source: []string{
			"NOP",
			"PORT #0x0f",
			"WAIT 200ms",
			"RJMP #0x03",
			"LDA #0x09",
			"SWAP",
			"MOV B,A",
			"MOV C,A",
			"MOV D,A",
			"STA DOUT",
			"STA DOUT1",
			"STA DOUT2",
			"STA DOUT3",
			"STA DOUT4",
			"STA PWM1",
			"STA PWM2",
			"STA SRV1",
			"STA SRV2",
			"MOV E,A",
			"MOV F,A",
			"PUSH",
		},
		Labels: map[string]label{},
		Macros: map[string]macro{},
		Code: []string{
			"NOP",
			"PORT #0x0f",
			"WAIT 200ms",
			"RJMP #0x03",
			"LDA #0x09",
			"SWAP",
			"MOV B,A",
			"MOV C,A",
			"MOV D,A",
			"STA DOUT",
			"STA DOUT1",
			"STA DOUT2",
			"STA DOUT3",
			"STA DOUT4",
			"STA PWM1",
			"STA PWM2",
			"STA SRV1",
			"STA SRV2",
			"MOV E,A",
			"MOV F,A",
			"PUSH",
		},
		Binary: []byte{
			0x00,
			0x1f,
			0x27,
			0x33,
			0x49,
			0x50,
			0x51,
			0x52,
			0x53,
			0x54,
			0x55,
			0x56,
			0x57,
			0x58,
			0x59,
			0x5A,
			0x5B,
			0x5C,
			0x5D,
			0x5E,
			0x5F,
		},
	},
	errCount: 0,
},
}

func TestOne(t *testing.T) {
	log.Logger.SetLevel(log.LvError)
	oneTest(t, "testProgramm2")
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
			tasm.generate()
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
	// Testing code
	ast.Equal(expected.Binary, actual.Binary, "Code is not equal")
}
