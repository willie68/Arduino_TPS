package asm

import (
	"fmt"
	"io"
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
func TpsFile(writer io.Writer, filename string, tpsasm Assembler) error {
	io.WriteString(writer, fmt.Sprintf("#TPS: asm generated file: %s \n", filename))
	for x, v := range tpsasm.Binary {
		io.WriteString(writer, fmt.Sprintf("0x%04x,%x,%x,\"%s\"\n", x, (v&0xf0)>>4, (v&0x0f), tpsasm.Code[x]))
	}
	return nil
}
