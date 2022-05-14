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
		Binary: []byte{0x4c},
	},
	errCount: 0,
}, {
	name: "testTrimming1",
	asm: Assembler{
		Source: []string{"LDA #12", " LDA #12", "LDA #12 ", "\t LDA #12", "\t LDA #12  \t"},
		Labels: map[string]label{},
		Code:   []string{"LDA #12", "LDA #12", "LDA #12", "LDA #12", "LDA #12"},
		Binary: []byte{0x4c, 0x4c, 0x4c, 0x4c, 0x4c},
	},
	errCount: 0,
}, {
	name: "testComment1",
	asm: Assembler{
		Source: []string{"/* dies ist ein Test", "1", "2", "3", "4", "*/"},
		Labels: map[string]label{},
		Code:   []string{},
		Binary: []byte{},
	},
	errCount: 0,
}, {
	name: "testComment2",
	asm: Assembler{
		Source: []string{"/* dies ist ein Test", "1", "2", "*/", "LDA #12", "/*", "3", "4", "*/"},
		Labels: map[string]label{},
		Code:   []string{"LDA #12"},
		Binary: []byte{0x4c},
	},
	errCount: 0,
}, {
	name: "testErrorComment1",
	asm: Assembler{
		Source: []string{"*/"},
		Labels: map[string]label{},
		Code:   []string{},
		Binary: []byte{},
	},
	errCount: 1,
}, {
	name: "testErrorComment2",
	asm: Assembler{
		Source: []string{"/*", "/*", "*/"},
		Labels: map[string]label{},
		Code:   []string{},
		Binary: []byte{},
	},
	errCount: 1,
}, {
	name: "testLineComment1",
	asm: Assembler{
		Source: []string{"; dies ist ein Test"},
		Labels: map[string]label{},
		Code:   []string{},
		Binary: []byte{},
	},
	errCount: 0,
}, {
	name: "testLineComment2",
	asm: Assembler{
		Source: []string{" ; dies ist ein Test"},
		Labels: map[string]label{},
		Code:   []string{},
		Binary: []byte{},
	},
	errCount: 0,
}, {
	name: "testInlineComment1",
	asm: Assembler{
		Source: []string{"LDA #12; dies ist ein Test"},
		Labels: map[string]label{},
		Code:   []string{"LDA #12"},
		Binary: []byte{0x4C},
	},
	errCount: 0,
}, {
	name: "testInlineComment2",
	asm: Assembler{
		Source: []string{"LDA #12 ; dies ist ein Test"},
		Labels: map[string]label{},
		Code:   []string{"LDA #12"},
		Binary: []byte{0x4c},
	},
	errCount: 0,
}, {
	name: "testlabel1",
	asm: Assembler{
		Source: []string{":label1", ":label2"},
		Labels: map[string]label{
			"label1": {
				Name:       "label1",
				PrgCounter: 0,
				Line:       0,
			},
			"label2": {
				Name:       "label2",
				PrgCounter: 0,
				Line:       1,
			},
		},
		Code:   []string{},
		Binary: []byte{},
	},
	errCount: 0,
}, {
	name: "testMacro1",
	asm: Assembler{
		Source: []string{":label1", ":label2", ".macro define", "LDA #12", ".endmacro"},
		Labels: map[string]label{
			"label1": {
				Name:       "label1",
				PrgCounter: 0,
				Line:       0,
			},
			"label2": {
				Name:       "label2",
				PrgCounter: 0,
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
		Code:   []string{},
		Binary: []byte{},
	},
	errCount: 0,
}, {
	name: "testMacroWithParameter",
	asm: Assembler{
		Source: []string{":label1", ":label2", ".macro define muck puck", "LDA #12", ".endmacro"},
		Labels: map[string]label{
			"label1": {
				Name:       "label1",
				PrgCounter: 0,
				Line:       0,
			},
			"label2": {
				Name:       "label2",
				PrgCounter: 0,
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
		Code:   []string{},
		Binary: []byte{},
	},
	errCount: 0,
}, {
	name: "testMacroError1",
	asm: Assembler{
		Source: []string{":label1", ":label2", ".macro", ".endmacro"},
		Labels: map[string]label{
			"label1": {
				Name:       "label1",
				PrgCounter: 0,
				Line:       0,
			},
			"label2": {
				Name:       "label2",
				PrgCounter: 0,
				Line:       1,
			},
		},
		Code:   []string{},
		Binary: []byte{},
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
		Code:   []string{},
		Binary: []byte{},
	},
	errCount: 1,
}, {
	name: "testMacroProcessing1",
	asm: Assembler{
		Source: []string{".macro define1", "LDA #12", "PORT #12", ".endmacro", ".define1"},
		Labels: map[string]label{},
		Macros: map[string]macro{
			"define1": {
				Name: "define1",
				Code: []string{"LDA #12", "PORT #12"},
			},
		},
		Code:   []string{"LDA #12", "PORT #12"},
		Binary: []byte{0x4c, 0x1c},
	},
	errCount: 0,
}, {
	name: "testMacroProcessing2",
	asm: Assembler{
		Source: []string{".macro define1 time name", "LDA time", "PORT name", ".endmacro", ".define1 #12 #13"},
		Labels: map[string]label{},
		Macros: map[string]macro{
			"define1": {
				Name:   "define1",
				Params: []string{"time", "name"},
				Code:   []string{"LDA time", "PORT name"},
			},
		},
		Code:   []string{"LDA #12", "PORT #13"},
		Binary: []byte{0x4c, 0x1d},
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
		Code:   []string{},
		Binary: []byte{},
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
}, {
	name: "testProgramm6",
	asm: Assembler{
		Source: []string{
			"MOV A,B",
			"MOV A,C",
			"MOV A,D",
			"LDA DIN",
			"LDA DIN1",
			"LDA DIN2",
			"LDA DIN3",
			"LDA DIN4",
			"LDA ADC1",
			"LDA ADC2",
			"LDA RC1",
			"LDA RC2",
			"MOV A,E",
			"MOV A,F",
			"POP",
		},
		Labels: map[string]label{},
		Macros: map[string]macro{},
		Code: []string{
			"MOV A,B",
			"MOV A,C",
			"MOV A,D",
			"LDA DIN",
			"LDA DIN1",
			"LDA DIN2",
			"LDA DIN3",
			"LDA DIN4",
			"LDA ADC1",
			"LDA ADC2",
			"LDA RC1",
			"LDA RC2",
			"MOV A,E",
			"MOV A,F",
			"POP",
		},
		Binary: []byte{
			0x61,
			0x62,
			0x63,
			0x64,
			0x65,
			0x66,
			0x67,
			0x68,
			0x69,
			0x6A,
			0x6B,
			0x6C,
			0x6D,
			0x6E,
			0x6F,
		},
	},
	errCount: 0,
}, {
	name: "testProgramm7",
	asm: Assembler{
		Source: []string{
			"INC",
			"DEC",
			"ADD",
			"SUB",
			"MUL",
			"DIV",
			"AND",
			"OR",
			"XOR",
			"NOT",
			"MOD",
			"BYTE",
			"BSUBA",
			"SHR",
			"SHL",
		},
		Labels: map[string]label{},
		Macros: map[string]macro{},
		Code: []string{
			"INC",
			"DEC",
			"ADD",
			"SUB",
			"MUL",
			"DIV",
			"AND",
			"OR",
			"XOR",
			"NOT",
			"MOD",
			"BYTE",
			"BSUBA",
			"SHR",
			"SHL",
		},
		Binary: []byte{
			0x71,
			0x72,
			0x73,
			0x74,
			0x75,
			0x76,
			0x77,
			0x78,
			0x79,
			0x7A,
			0x7B,
			0x7C,
			0x7D,
			0x7E,
			0x7F,
		},
	},
	errCount: 0,
}, {
	name: "testProgramm8.BD",
	asm: Assembler{
		Source: []string{
			"PAGE #0",
			"PAGE #4",
			"PAGE #7",
			"PAGE #8",
			"PAGE #15",
			"JMP #0",
			"JMP #4",
			"JMP #7",
			"JMP #8",
			"JMP #0x0f",
			"LOOPC #0",
			"LOOPC #4",
			"LOOPC #7",
			"LOOPC #8",
			"LOOPC #0x0f",
			"LOOPD #0",
			"LOOPD #4",
			"LOOPD #7",
			"LOOPD #8",
			"LOOPD #0x0f",
			"CALL #0",
			"CALL #4",
			"CALL #7",
			"CALL #8",
			"CALL #0x0f",
		},
		Labels: map[string]label{},
		Macros: map[string]macro{},
		Code: []string{
			"PAGE #0",
			"PAGE #4",
			"PAGE #7",
			"PAGE #8",
			"PAGE #15",
			"JMP #0",
			"JMP #4",
			"JMP #7",
			"JMP #8",
			"JMP #0x0f",
			"LOOPC #0",
			"LOOPC #4",
			"LOOPC #7",
			"LOOPC #8",
			"LOOPC #0x0f",
			"LOOPD #0",
			"LOOPD #4",
			"LOOPD #7",
			"LOOPD #8",
			"LOOPD #0x0f",
			"CALL #0",
			"CALL #4",
			"CALL #7",
			"CALL #8",
			"CALL #0x0f",
		},
		Binary: []byte{
			0x80,
			0x84,
			0x87,
			0x88,
			0x8F,
			0x90,
			0x94,
			0x97,
			0x98,
			0x9F,
			0xA0,
			0xA4,
			0xA7,
			0xA8,
			0xAF,
			0xB0,
			0xB4,
			0xB7,
			0xB8,
			0xBF,
			0xD0,
			0xD4,
			0xD7,
			0xD8,
			0xDF,
		},
	},
	errCount: 0,
}, {
	name: "testProgrammCEF",
	asm: Assembler{
		Source: []string{
			"SKIP0",
			"AGTB",
			"ALTB",
			"AEQB",
			"DEQ1 1",
			"DEQ1 2",
			"DEQ1 3",
			"DEQ1 4",
			"DEQ0 1",
			"DEQ0 2",
			"DEQ0 3",
			"DEQ0 4",
			"PRG0",
			"SEL0",
			"PRG1",
			"SEL1",
			"RTR",
			"CASB #1",
			"CASB #2",
			"CASB #3",
			"CASB #4",
			"CASB #5",
			"CASB #6",
			"DFSB #1",
			"DFSB #2",
			"DFSB #3",
			"DFSB #4",
			"DFSB #5",
			"DFSB #6",
			"REST",
			"BLDA ADC1",
			"BLDA ADC2",
			"BLDA RC1",
			"BLDA RC2",
			"BSTA PWM1",
			"BSTA PWM2",
			"BSTA SRV1",
			"BSTA SRV2",
			"TONE",
			"PEND",
		},
		Labels: map[string]label{},
		Macros: map[string]macro{},
		Code: []string{
			"SKIP0",
			"AGTB",
			"ALTB",
			"AEQB",
			"DEQ1 1",
			"DEQ1 2",
			"DEQ1 3",
			"DEQ1 4",
			"DEQ0 1",
			"DEQ0 2",
			"DEQ0 3",
			"DEQ0 4",
			"PRG0",
			"SEL0",
			"PRG1",
			"SEL1",
			"RTR",
			"CASB #1",
			"CASB #2",
			"CASB #3",
			"CASB #4",
			"CASB #5",
			"CASB #6",
			"DFSB #1",
			"DFSB #2",
			"DFSB #3",
			"DFSB #4",
			"DFSB #5",
			"DFSB #6",
			"REST",
			"BLDA ADC1",
			"BLDA ADC2",
			"BLDA RC1",
			"BLDA RC2",
			"BSTA PWM1",
			"BSTA PWM2",
			"BSTA SRV1",
			"BSTA SRV2",
			"TONE",
			"PEND"},
		Binary: []byte{
			0xC0,
			0xC1,
			0xC2,
			0xC3,
			0xC4,
			0xC5,
			0xC6,
			0xC7,
			0xC8,
			0xC9,
			0xCA,
			0xCB,
			0xCC,
			0xCD,
			0xCE,
			0xCF,
			0xE0,
			0xE1,
			0xE2,
			0xE3,
			0xE4,
			0xE5,
			0xE6,
			0xE8,
			0xE9,
			0xEA,
			0xEB,
			0xEC,
			0xED,
			0xEF,
			0xF0,
			0xF1,
			0xF2,
			0xF3,
			0xF4,
			0xF5,
			0xF6,
			0xF7,
			0xF8,
			0xFF,
		},
	},
	errCount: 0,
},
}

func TestOne(t *testing.T) {
	log.Logger.SetLevel(log.LvError)
	oneTest(t, "testProgrammCEF")
}

func TestAll(t *testing.T) {
	log.Logger.SetLevel(log.LvError)
	for _, test := range testdatas {
		oneTest(t, test.name)
	}
}

func oneTest(t *testing.T, name string) {
	ast := assert.New(t)
	found := false
	for _, test := range testdatas {
		if name == test.name {
			found = true
			fmt.Printf("testing %s\r\n", test.name)

			tasm := Assembler{
				Hardware: Holtek,
				Source:   test.asm.Source,
			}
			tasm.Parse()
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
	if !found {
		t.Logf("test %s not found", name)
		t.Fail()
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
