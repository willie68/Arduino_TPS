package asm

import (
	"fmt"
	"io"
	"strings"

	"github.com/marcinbor85/gohex"
)

/* writing the file in tps format:
#TPS:Willies SPS Arduino
0x00,1,5,""
0x01,2,8,""
0x02,1,A,""
0x03,2,8,""
0x04,3,4,""
0x05,0,0,""
*/
func TpsFile(writer io.Writer, tpsasm Assembler) error {
	switch strings.ToLower(tpsasm.Outputformat) {
	case "bin":
		_, err := writer.Write(tpsasm.Binary)
		if err != nil {
			return err
		}
	case "intelhex":
		mem := gohex.NewMemory()
		err := mem.AddBinary(0, tpsasm.Binary)
		if err != nil {
			return err
		}
		err = mem.DumpIntelHex(writer, 8)
		if err != nil {
			return err
		}
	case "tps":
		_, err := io.WriteString(writer, fmt.Sprintf("#TPS: asm generated file: %s \n", tpsasm.Filename))
		if err != nil {
			return err
		}
		for x, v := range tpsasm.Binary {
			_, err = io.WriteString(writer, fmt.Sprintf("0x%04x,%x,%x,\"%s\"\n", x, (v&0xf0)>>4, (v&0x0f), tpsasm.Code[x]))
			if err != nil {
				return err
			}
		}
	}
	return nil
}
