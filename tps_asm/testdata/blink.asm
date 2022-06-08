.arduinotps ;damit wird die Hardware auf Arduino_TPS festgelegt. Diese Directrive sollte vor dem eigentlichen Code erscheinen.
;.tinytps ; legt die Hardware auf die Tiny_TPS fest.
;.atmega8 ; legt die Hardware auf die ATMega8 fest.
;.holtek ; legt die Hardware auf Holtek fest.

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
*/

.macro macro1 output time
PORT output
WAIT time
PORT #0x00
WAIT time
.endmacro

:loop1
PORT #0x0F ;Zeilenkommentar
WAIT 200ms
PORT #0x00
WAIT 200ms
RJMP :loop1

DFSB #1
PORT #0x0F ;Zeilenkommentar
WAIT 200ms
PORT #0x00
WAIT 200ms
RTR
