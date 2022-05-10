.macro blink
PORT #0B0101
WAIT 200ms
PORT #0B1010
WAIT 200ms
.endmacro

:loop
.blink
RJMP :loop
/* 
Kommentar über mehrere Zeilen
/*
*/

.macro macro1 output time
PORT output
WAIT time
PORT #0x00
WAIT time
.endmacro

.include macro_blink
:loop1
.macro1 #0x0f 100ms

PORT #0x0F ;Zeilenkommentar
WAIT 200ms
PORT #0x00
WAIT 200ms
RJMP :loop1

DFSB 1
PORT #0x0F ;Zeilenkommentar
WAIT 200ms
PORT #0x00
WAIT 200ms
RTR
