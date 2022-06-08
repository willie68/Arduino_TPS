:loop
LDA #15
STA Dout
WAIT 200ms
PRG0
LDA #00
STA Dout
WAIT 500ms
RJMP :loop
