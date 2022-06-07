package asm

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnknown(t *testing.T) {
	ast := assert.New(t)
	mno, err := GetMnemonic("muck")
	ast.NotNil(err)
	ast.Nil(mno)
}

func TestCtrl(t *testing.T) {
	ast := assert.New(t)

	var testDatas = []struct {
		name string
		code byte
	}{
		{
			name: "nop",
			code: 0x00,
		},
		{
			name: "rest",
			code: 0xEF,
		},
		{
			name: "pend",
			code: 0xFF,
		},
	}

	for x, td := range testDatas {
		mno, err := GetMnemonic(td.name)
		ast.Nil(err)
		ast.NotNil(mno)

		name := mno.Name
		name = strings.ToLower(name)
		ast.Equal(td.name, name)

		ast.Nil(mno.CheckParameter(""))

		ast.NotNil(mno.CheckParameter("#0x0e"))
		ast.NotNil(mno.CheckParameter("#12"))
		ast.NotNil(mno.CheckParameter("#0b0011"))
		ast.NotNil(mno.CheckParameter("#0b0211"))
		ast.NotNil(mno.CheckParameter("#16"))
		ast.NotNil(mno.CheckParameter("#0x3e"))
		ast.NotNil(mno.CheckParameter(":loop"))

		ast.Equal(td.code, mno.Code)

		if x == 0 {
			ast.Nil(mno.CheckHardware(Holtek))
			ast.Nil(mno.CheckHardware(ATMega8))
			ast.Nil(mno.CheckHardware(ArduinoSPS))
			ast.Nil(mno.CheckHardware(TinySPS))
		} else {
			ast.NotNil(mno.CheckHardware(Holtek))
			ast.NotNil(mno.CheckHardware(ATMega8))
			ast.Nil(mno.CheckHardware(ArduinoSPS))
			ast.Nil(mno.CheckHardware(TinySPS))
		}
	}
}

func TestPort(t *testing.T) {
	ast := assert.New(t)
	mno, err := GetMnemonic("PoRt")
	ast.Nil(err)
	ast.NotNil(mno)

	name := mno.Name
	name = strings.ToLower(name)
	ast.Equal("port", name)

	ast.NotNil(mno.CheckParameter(""))
	ast.Nil(mno.CheckParameter("#0x0e"))
	ast.Nil(mno.CheckParameter("#12"))
	ast.Nil(mno.CheckParameter("#0b0011"))
	ast.NotNil(mno.CheckParameter("#0x3e"))
	ast.NotNil(mno.CheckParameter("#0x3f"))
	ast.NotNil(mno.CheckParameter(":loop"))

	ast.Nil(mno.CheckHardware(Holtek))
	ast.Nil(mno.CheckHardware(ATMega8))
	ast.Nil(mno.CheckHardware(ArduinoSPS))
	ast.Nil(mno.CheckHardware(TinySPS))
}

func TestWait(t *testing.T) {
	ast := assert.New(t)
	mno, err := GetMnemonic("wait")
	ast.Nil(err)
	ast.NotNil(mno)

	name := mno.Name
	name = strings.ToLower(name)
	ast.Equal("wait", name)

	ast.NotNil(mno.CheckParameter(""))
	ast.Nil(mno.CheckParameter("#0x0e"))
	ast.Nil(mno.CheckParameter("200ms"))
	ast.NotNil(mno.CheckParameter(":loop"))

	ast.Nil(mno.CheckHardware(Holtek))
	ast.Nil(mno.CheckHardware(ATMega8))
	ast.Nil(mno.CheckHardware(ArduinoSPS))
	ast.Nil(mno.CheckHardware(TinySPS))
}

func TestPage(t *testing.T) {
	ast := assert.New(t)
	mno, err := GetMnemonic("Page")
	ast.Nil(err)
	ast.NotNil(mno)

	name := mno.Name
	name = strings.ToLower(name)
	ast.Equal("page", name)

	ast.NotNil(mno.CheckParameter(""))
	ast.Nil(mno.CheckParameter("#0x07"))
	ast.Nil(mno.CheckParameter(":?"))
	ast.NotNil(mno.CheckParameter(":loop"))

	ast.Nil(mno.CheckHardware(Holtek))
	ast.Nil(mno.CheckHardware(ATMega8))
	ast.Nil(mno.CheckHardware(ArduinoSPS))
	ast.Nil(mno.CheckHardware(TinySPS))
}

func TestJump(t *testing.T) {
	ast := assert.New(t)
	mno, err := GetMnemonic("JMP")
	ast.Nil(err)
	ast.NotNil(mno)

	name := mno.Name
	name = strings.ToLower(name)
	ast.Equal("jmp", name)

	ast.NotNil(mno.CheckParameter(""))
	ast.Nil(mno.CheckParameter("#0x0e"))
	ast.Nil(mno.CheckParameter(":loop"))
	ast.NotNil(mno.CheckParameter("#0x1e"))
	ast.NotNil(mno.CheckParameter("muck"))
}

func TestLoop(t *testing.T) {
	ast := assert.New(t)
	var testdatas = []struct {
		name string
		code byte
	}{
		{
			name: "loopc",
			code: 0xa0,
		},
		{
			name: "loopc",
			code: 0xa0,
		},
	}
	for _, td := range testdatas {
		mno, err := GetMnemonic(td.name)
		ast.Nil(err)
		ast.NotNil(mno)

		name := mno.Name
		name = strings.ToLower(name)
		ast.Equal(td.name, name)

		ast.NotNil(mno.CheckParameter(""))
		ast.Nil(mno.CheckParameter("#0x0e"))
		ast.Nil(mno.CheckParameter(":loop"))
		ast.NotNil(mno.CheckParameter("#0x1e"))
		ast.NotNil(mno.CheckParameter("muck"))
	}
}

func TestSkip(t *testing.T) {
	ast := assert.New(t)

	var testDatas = []struct {
		name string
		code byte
	}{
		{
			name: "skip0",
			code: 0xC0,
		},
		{
			name: "agtb",
			code: 0xC1,
		},
		{
			name: "altb",
			code: 0xC2,
		},
		{
			name: "aeqb",
			code: 0xC3,
		},
		{
			name: "prg0",
			code: 0xCC,
		},
		{
			name: "sel0",
			code: 0xCD,
		},
		{
			name: "prg1",
			code: 0xCE,
		},
		{
			name: "sel1",
			code: 0xCF,
		},
	}

	for _, td := range testDatas {
		mno, err := GetMnemonic(td.name)
		ast.Nil(err)
		ast.NotNil(mno)

		name := mno.Name
		name = strings.ToLower(name)
		ast.Equal(td.name, name)

		ast.Nil(mno.CheckParameter(""))

		ast.NotNil(mno.CheckParameter("#0x0e"))
		ast.NotNil(mno.CheckParameter("#12"))
		ast.NotNil(mno.CheckParameter("#0b0011"))
		ast.NotNil(mno.CheckParameter("#0b0211"))
		ast.NotNil(mno.CheckParameter("#16"))
		ast.NotNil(mno.CheckParameter("#0x3e"))
		ast.NotNil(mno.CheckParameter(":loop"))

		ast.Equal(td.code, mno.Code)
	}
}

func TestDEQ(t *testing.T) {
	ast := assert.New(t)

	var testDatas = []struct {
		name string
		code byte
	}{
		{
			name: "deq0",
			code: 0xC0,
		},
		{
			name: "deq1",
			code: 0xC0,
		},
	}

	for _, td := range testDatas {
		mno, err := GetMnemonic(td.name)
		ast.Nil(err)
		ast.NotNil(mno)

		name := mno.Name
		name = strings.ToLower(name)
		ast.Equal(td.name, name)

		ast.Nil(mno.CheckParameter("1"))
		ast.Nil(mno.CheckParameter("2"))
		ast.Nil(mno.CheckParameter("3"))
		ast.Nil(mno.CheckParameter("4"))

		ast.NotNil(mno.CheckParameter(""))

		ast.NotNil(mno.CheckParameter("#0x0e"))
		ast.NotNil(mno.CheckParameter("#12"))
		ast.NotNil(mno.CheckParameter("#0b0011"))
		ast.NotNil(mno.CheckParameter("#0b0211"))
		ast.NotNil(mno.CheckParameter("#16"))
		ast.NotNil(mno.CheckParameter("#0x3e"))
		ast.NotNil(mno.CheckParameter(":loop"))

		ast.Equal(td.code, mno.Code)
	}
}

func TestCall(t *testing.T) {
	ast := assert.New(t)
	mno, err := GetMnemonic("CALL")
	ast.Nil(err)
	ast.NotNil(mno)

	name := mno.Name
	name = strings.ToLower(name)
	ast.Equal("call", name)

	ast.NotNil(mno.CheckParameter(""))
	ast.Nil(mno.CheckParameter("#0x0e"))
	ast.Nil(mno.CheckParameter(":loop"))
	ast.NotNil(mno.CheckParameter("#0x1e"))
	ast.NotNil(mno.CheckParameter("muck"))

	mno, err = GetMnemonic("RTR")
	ast.Nil(err)
	ast.NotNil(mno)

	name = mno.Name
	name = strings.ToLower(name)
	ast.Equal("rtr", name)

	ast.Nil(mno.CheckParameter(""))
	ast.NotNil(mno.CheckParameter("#0x0e"))
	ast.NotNil(mno.CheckParameter(":loop"))
	ast.NotNil(mno.CheckParameter("#0x1e"))
	ast.NotNil(mno.CheckParameter("muck"))
}

func TestSub(t *testing.T) {
	ast := assert.New(t)
	mno, err := GetMnemonic("DFSB")
	ast.Nil(err)
	ast.NotNil(mno)

	name := mno.Name
	name = strings.ToLower(name)
	ast.Equal("dfsb", name)

	ast.NotNil(mno.CheckParameter(""))
	ast.Nil(mno.CheckParameter("#0x0e"))
	ast.Nil(mno.CheckParameter(":loop"))
	ast.NotNil(mno.CheckParameter("#0x1e"))
	ast.NotNil(mno.CheckParameter("muck"))

	mno, err = GetMnemonic("CASB")
	ast.Nil(err)
	ast.NotNil(mno)

	name = mno.Name
	name = strings.ToLower(name)
	ast.Equal("casb", name)

	ast.NotNil(mno.CheckParameter(""))
	ast.Nil(mno.CheckParameter("#0x0e"))
	ast.Nil(mno.CheckParameter(":loop"))
	ast.NotNil(mno.CheckParameter("#0x1e"))
	ast.NotNil(mno.CheckParameter("muck"))
}

func TestLDA(t *testing.T) {
	ast := assert.New(t)
	mno, err := GetMnemonic("LDA")
	ast.Nil(err)
	ast.NotNil(mno)

	name := mno.Name
	name = strings.ToLower(name)
	ast.Equal("lda", name)

	ast.NotNil(mno.CheckParameter(""))
	ast.Nil(mno.CheckParameter("#0x0e"))
	ast.Nil(mno.CheckParameter("#12"))
	ast.Nil(mno.CheckParameter("#0b0011"))
	ast.NotNil(mno.CheckParameter("#0x3e"))
	ast.NotNil(mno.CheckParameter("#0x3f"))
	ast.NotNil(mno.CheckParameter(":loop"))
	ast.Equal(uint8(0x40), mno.Code)

	ast.Nil(mno.CheckParameter("DIN"))
	ast.Nil(mno.CheckParameter("DIN1"))
	ast.Nil(mno.CheckParameter("DIN2"))
	ast.Nil(mno.CheckParameter("DIN3"))
	ast.Nil(mno.CheckParameter("DIN4"))
	ast.Nil(mno.CheckParameter("ADC1"))
	ast.Nil(mno.CheckParameter("ADC2"))
	ast.Nil(mno.CheckParameter("RC1"))
	ast.Nil(mno.CheckParameter("RC2"))
}

func TestSTA(t *testing.T) {
	ast := assert.New(t)
	mno, err := GetMnemonic("STA")
	ast.Nil(err)
	ast.NotNil(mno)

	name := mno.Name
	name = strings.ToLower(name)
	ast.Equal("sta", name)

	ast.NotNil(mno.CheckParameter(""))
	ast.NotNil(mno.CheckParameter("#0x0e"))
	ast.NotNil(mno.CheckParameter("#12"))
	ast.NotNil(mno.CheckParameter("#0b0011"))
	ast.NotNil(mno.CheckParameter("#0x3e"))
	ast.NotNil(mno.CheckParameter("#0x3f"))
	ast.NotNil(mno.CheckParameter(":loop"))
	ast.Equal(uint8(0x00), mno.Code)

	ast.Nil(mno.CheckParameter("DOUT"))
	ast.Nil(mno.CheckParameter("DOUT1"))
	ast.Nil(mno.CheckParameter("DOUT2"))
	ast.Nil(mno.CheckParameter("DOUT3"))
	ast.Nil(mno.CheckParameter("DOUT4"))
	ast.Nil(mno.CheckParameter("PWM1"))
	ast.Nil(mno.CheckParameter("PWM2"))
	ast.Nil(mno.CheckParameter("SRV1"))
	ast.Nil(mno.CheckParameter("SRV2"))
}

func TestMov(t *testing.T) {
	ast := assert.New(t)
	mno, err := GetMnemonic("MOV")
	ast.Nil(err)
	ast.NotNil(mno)

	name := mno.Name
	name = strings.ToLower(name)
	ast.Equal("mov", name)

	ast.NotNil(mno.CheckParameter(""))
	ast.NotNil(mno.CheckParameter("#0x3e"))
	ast.NotNil(mno.CheckParameter("#0x3f"))
	ast.NotNil(mno.CheckParameter(":loop"))

	ast.Nil(mno.CheckParameter("A,B"))
	ast.Nil(mno.CheckParameter("A,C"))
	ast.Nil(mno.CheckParameter("A,D"))
	ast.Nil(mno.CheckParameter("A,E"))
	ast.Nil(mno.CheckParameter("A,F"))
	ast.Nil(mno.CheckParameter("B,A"))
	ast.Nil(mno.CheckParameter("C,A"))
	ast.Nil(mno.CheckParameter("D,A"))
	ast.Nil(mno.CheckParameter("E,A"))
	ast.Nil(mno.CheckParameter("F,A"))
}

func TestLoadStore(t *testing.T) {
	ast := assert.New(t)

	var testDatas = []struct {
		Name string
		Code byte
	}{
		{
			Name: "SWAP",
			Code: 0x50,
		},
		{
			Name: "PUSH",
			Code: 0x5F,
		},
		{
			Name: "POP",
			Code: 0x6F,
		},
	}

	for _, testdata := range testDatas {
		mno, err := GetMnemonic(testdata.Name)
		ast.Nil(err)
		ast.NotNil(mno)
		ast.Nil(mno.CheckParameter(""))

		ast.NotNil(mno.CheckParameter("#0x0e"))
		ast.NotNil(mno.CheckParameter("#12"))
		ast.NotNil(mno.CheckParameter("#0b0011"))
		ast.NotNil(mno.CheckParameter("#0b0211"))
		ast.NotNil(mno.CheckParameter("#16"))
		ast.NotNil(mno.CheckParameter("#0x3e"))
		ast.NotNil(mno.CheckParameter(":loop"))

		ast.Equal(testdata.Code, mno.Code)
	}
}

func TestNumberFormat(t *testing.T) {
	ast := assert.New(t)
	mno, err := GetMnemonic("LDA")
	ast.Nil(err)
	ast.NotNil(mno)

	ast.Nil(mno.CheckParameter("#0x0e"))
	ast.Nil(mno.CheckParameter("#12"))
	ast.Nil(mno.CheckParameter("#0b0011"))
	ast.NotNil(mno.CheckParameter("#0b0211"))
	ast.NotNil(mno.CheckParameter("#16"))
	ast.NotNil(mno.CheckParameter("#0x3e"))
	ast.NotNil(mno.CheckParameter(":loop"))
}

func TestMathMno(t *testing.T) {
	ast := assert.New(t)

	var testDatas = []struct {
		Name string
		Code byte
	}{
		{
			Name: "INC",
			Code: 0x71,
		},
		{
			Name: "DEC",
			Code: 0x72,
		},
		{
			Name: "ADD",
			Code: 0x73,
		},
		{
			Name: "SUB",
			Code: 0x74,
		},
		{
			Name: "MUL",
			Code: 0x75,
		},
		{
			Name: "DIV",
			Code: 0x76,
		},
		{
			Name: "AND",
			Code: 0x77,
		},
		{
			Name: "OR",
			Code: 0x78,
		},
		{
			Name: "XOR",
			Code: 0x79,
		},
		{
			Name: "NOT",
			Code: 0x7A,
		},
		{
			Name: "MOD",
			Code: 0x7B,
		},
		{
			Name: "SHR",
			Code: 0x7E,
		},
		{
			Name: "SHL",
			Code: 0x7F,
		},
		{
			Name: "BYTE",
			Code: 0x7C,
		},
		{
			Name: "BSUBA",
			Code: 0x7D,
		},
	}

	for _, testdata := range testDatas {
		mno, err := GetMnemonic(testdata.Name)
		ast.Nil(err)
		ast.NotNil(mno)
		ast.Nil(mno.CheckParameter(""))

		ast.NotNil(mno.CheckParameter("#0x0e"))
		ast.NotNil(mno.CheckParameter("#12"))
		ast.NotNil(mno.CheckParameter("#0b0011"))
		ast.NotNil(mno.CheckParameter("#0b0211"))
		ast.NotNil(mno.CheckParameter("#16"))
		ast.NotNil(mno.CheckParameter("#0x3e"))
		ast.NotNil(mno.CheckParameter(":loop"))

		ast.Equal(testdata.Code, mno.Code)
	}
}

func TestBLDA(t *testing.T) {
	ast := assert.New(t)
	mno, err := GetMnemonic("BLDA")
	ast.Nil(err)
	ast.NotNil(mno)

	name := mno.Name
	name = strings.ToLower(name)
	ast.Equal("blda", name)

	ast.NotNil(mno.CheckParameter(""))
	ast.NotNil(mno.CheckParameter("#0x0e"))
	ast.NotNil(mno.CheckParameter("#12"))
	ast.NotNil(mno.CheckParameter("#0b0011"))
	ast.NotNil(mno.CheckParameter("#0x3e"))
	ast.NotNil(mno.CheckParameter("#0x3f"))
	ast.NotNil(mno.CheckParameter(":loop"))
	ast.Equal(uint8(0x00), mno.Code)

	ast.Nil(mno.CheckParameter("ADC1"))
	ast.Nil(mno.CheckParameter("ADC2"))
	ast.Nil(mno.CheckParameter("RC1"))
	ast.Nil(mno.CheckParameter("RC2"))
}

func TestBSTA(t *testing.T) {
	ast := assert.New(t)
	mno, err := GetMnemonic("BSTA")
	ast.Nil(err)
	ast.NotNil(mno)

	name := mno.Name
	name = strings.ToLower(name)
	ast.Equal("bsta", name)

	ast.NotNil(mno.CheckParameter(""))
	ast.NotNil(mno.CheckParameter("#0x0e"))
	ast.NotNil(mno.CheckParameter("#12"))
	ast.NotNil(mno.CheckParameter("#0b0011"))
	ast.NotNil(mno.CheckParameter("#0x3e"))
	ast.NotNil(mno.CheckParameter("#0x3f"))
	ast.NotNil(mno.CheckParameter(":loop"))
	ast.Equal(uint8(0x00), mno.Code)

	ast.Nil(mno.CheckParameter("PWM1"))
	ast.Nil(mno.CheckParameter("PWM2"))
	ast.Nil(mno.CheckParameter("SRV1"))
	ast.Nil(mno.CheckParameter("SRV2"))
}

func TestByte(t *testing.T) {
	ast := assert.New(t)
	mno, err := GetMnemonic("TONE")
	ast.Nil(err)
	ast.NotNil(mno)

	name := mno.Name
	name = strings.ToLower(name)
	ast.Equal("tone", name)

	ast.Nil(mno.CheckParameter(""))
	ast.NotNil(mno.CheckParameter("#0x0e"))
	ast.NotNil(mno.CheckParameter("#12"))
	ast.NotNil(mno.CheckParameter("#0b0011"))
	ast.NotNil(mno.CheckParameter("#0x3e"))
	ast.NotNil(mno.CheckParameter("#0x3f"))
	ast.NotNil(mno.CheckParameter(":loop"))
	ast.Equal(uint8(0xF8), mno.Code)
}
