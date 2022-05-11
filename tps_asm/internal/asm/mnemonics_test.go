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

func TestNop(t *testing.T) {
	ast := assert.New(t)
	mno, err := GetMnemonic("nop")
	ast.Nil(err)
	ast.NotNil(mno)

	name := mno.Name
	name = strings.ToLower(name)
	ast.Equal("nop", name)

	ast.Nil(mno.CheckParameter([]string{}))
	ast.NotNil(mno.CheckParameter([]string{"#0x0e"}))
	ast.NotNil(mno.CheckParameter([]string{":loop"}))
}

func TestPort(t *testing.T) {
	ast := assert.New(t)
	mno, err := GetMnemonic("PoRt")
	ast.Nil(err)
	ast.NotNil(mno)

	name := mno.Name
	name = strings.ToLower(name)
	ast.Equal("port", name)

	ast.NotNil(mno.CheckParameter([]string{}))
	ast.Nil(mno.CheckParameter([]string{"#0x0e"}))
	ast.Nil(mno.CheckParameter([]string{"#12"}))
	ast.Nil(mno.CheckParameter([]string{"#0b0011"}))
	ast.NotNil(mno.CheckParameter([]string{"#0x3e"}))
	ast.NotNil(mno.CheckParameter([]string{"#0x3f", "#0xde"}))
	ast.NotNil(mno.CheckParameter([]string{":loop"}))
}

func TestWait(t *testing.T) {
	ast := assert.New(t)
	mno, err := GetMnemonic("wait")
	ast.Nil(err)
	ast.NotNil(mno)

	name := mno.Name
	name = strings.ToLower(name)
	ast.Equal("wait", name)

	ast.NotNil(mno.CheckParameter([]string{}))
	ast.Nil(mno.CheckParameter([]string{"#0x0e"}))
	ast.Nil(mno.CheckParameter([]string{"200ms"}))
	ast.NotNil(mno.CheckParameter([]string{":loop"}))
}

func TestLDA(t *testing.T) {
	ast := assert.New(t)
	mno, err := GetMnemonic("LDA")
	ast.Nil(err)
	ast.NotNil(mno)

	name := mno.Name
	name = strings.ToLower(name)
	ast.Equal("lda", name)

	ast.NotNil(mno.CheckParameter([]string{}))
	ast.Nil(mno.CheckParameter([]string{"#0x0e"}))
	ast.Nil(mno.CheckParameter([]string{"#12"}))
	ast.Nil(mno.CheckParameter([]string{"#0b0011"}))
	ast.NotNil(mno.CheckParameter([]string{"#0x3e"}))
	ast.NotNil(mno.CheckParameter([]string{"#0x3f", "#0xde"}))
	ast.NotNil(mno.CheckParameter([]string{":loop"}))
}

func TestNumberFormat(t *testing.T) {
	ast := assert.New(t)
	mno, err := GetMnemonic("LDA")
	ast.Nil(err)
	ast.NotNil(mno)

	ast.Nil(mno.CheckParameter([]string{"#0x0e"}))
	ast.Nil(mno.CheckParameter([]string{"#12"}))
	ast.Nil(mno.CheckParameter([]string{"#0b0011"}))
	ast.NotNil(mno.CheckParameter([]string{"#0b0211"}))
	ast.NotNil(mno.CheckParameter([]string{"#16"}))
	ast.NotNil(mno.CheckParameter([]string{"#0x3e"}))
	ast.NotNil(mno.CheckParameter([]string{":loop"}))
}
