.rpi2040
.define DELAYTIME 500ms
:loop
LDA #5
LED1
LDA ADC1
MOV B,A
LDA ADC2
BYTE
BSTA SRV1
BSTA SRV2
BSTA SRV3
BSTA SRV4
WAIT DELAYTIME
RJMP :loop
