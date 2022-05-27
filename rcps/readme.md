# RCPC, Remote Controlled Programmable Controller

Die RCPC ist eine Steuerung von Modellbauern für Modellbauer. Sie soll einfache Steuerungsaufgaben ermöglichen, in einer einfachen Steuersprache, auch am Teich, der Rennstrecke oder dem Modellflughafen. 

Die Sprache ist einfach und systematisch, reduziert auf das nötigste. 

Als Hardware werden je nach Bedarf verschiedene Controller verwendet. Die Firmware beinhaltet den notwendigen Übersetzer. 

Im Web gibt es eine entsprechende Seite, mit deren Hilfe man die Steuerprogramm in den Steuercode übersetzen kann. Der Steuercode kann dann über eine einfache Schnittstelle in den Controller geladen werden. 



Es können verschiedene Konfigurationen aktiviert werden. Diese Konfigurationen werden im Programm festgelegt. 

| Konfiguration | digitale Ausgänge<br />(Dout) | digitale Eingänge<br />(Din) | analoge Ausgänge (PWM)<br />(Aout) | analoge Eingänge (0-5V)<br />(Ain) | Servos<br />(Srv) | RC Eingänge<br />RC | Tone |
| ------------- | ----------------------------- | ---------------------------- | ---------------------------------- | ---------------------------------- | ----------------- | ------------------- | ---- |
| 1             | 6                             | 6                            | 2                                  | 2                                  | 0                 | 0                   | 0    |
| 2             | 6                             | 6                            | 0                                  | 2                                  | 2                 | 2                   | 0    |
| 3             | 4                             | 4                            | 0                                  | 2                                  | 4                 | 1                   | 1    |
| 4             | 4                             | 4                            | 0                                  | 0                                  | 4                 | 4                   | 0    |
| 5             | 8                             | 2                            | 2                                  | 0                                  | 2                 | 2                   | 0    |



# RCPC Arduino Uno (Nano) Mapping

| Pin  | Funktion | K 1 | K 2 | K3 | K4 | K5 |
| ---- | -------- | ---- | ---- | -------- | ---- | ---- |
|  D0    | Rx (Prg) |      |     |          |      |      |
|  D1   | Tx (Prg) |      |      |          |      |      |
|  D2   |          | Din1 |Din1      | Din1 | Din1 | Din1  |
|  D3   | PWM | Din2 |Din2      | Din2 | Din2 | Aout1 |
|  D4   |          | Din3 |Din3      | Din3 | Din3 | Din2 |
|  D5   | PWM | Din4  |Din4      | Din4 | Din4 | Aout2 |
|  D6   | PWM | Dout1 |Dout1      | Dout1 | Dout1 | Dout1 |
|  D7   |          | Dout2 | Dout2 | Dout2 | Dout2 | Dout2 |
|  D8   |          | Dout3 |Dout3      | Dout3 | Dout3 | Dout3 |
| D9   | PWM | Aout1 | Srv1 | Srv1 | Srv1 | Srv1 |
| D10 | PWM | Aout2 | Srv2 | Srv2 | Srv2 | Srv2 |
| D11 | PWM | Dout4 | Dout4 | Dout4 | Dout4 | Dout4 |
| D12 | | Dout5 | Dout5 | Srv3 | Srv3 | Dout5 |
| D13 | LED | Dout6 | Dout6 | Srv4 | Srv4 | Dout6 |
| A0 | | Ain1 | Ain1 | Ain1 | RC3 | Dout7 |
| A1 | | Ain2 | Ain2 | Ain2 | RC4 | Dout8 |
| A2 | | Din5 | RC1 | RC1 | RC1 | RC1 |
| A3 | | Din6 | RC2 | Tone | RC2 | RC2 |
| A4 | SDA |  |  | | | |
| A5 | SCL |  |  | | | |

# Register/Speicher

All Register sind 8 Bit groß und können Werte von 0..255 auf nehmen. Negative Werte werden derzeit nicht unterstützt.

A, B sind die beiden Rechenregister.

C, D sind Hilfsregister

Z ist das Carry Register. >0 die letzte Operation erzeugte einen Überlauf.

X,Y sind Indexregister (für Schleifen)

Page: ist das Pageregister für Sprünge.

Weiterhin gibt es einen 16 Byte großer Stack. 

Der Programmspeicher ist min. 1kB groß.



# Mnemonics

## Control

| Mnemonic | Code   | short description      | Beschreibung                                                 |
| -------- | ------ | ---------------------- | ------------------------------------------------------------ |
| NOP      | 00     | No Operation           | es wird nichts gemacht, aber die Ausführung kostet trotzdem etwas Zeit und Prozessorzyklen |
| WAIT #x  | 2x     | wait <time>            | Verzögerung des Programms, 1⇐ x ⇐ 15, nach der vorher stehenden Tabelle . Wahlweise kann auch ein String verwendet werden. Mögliche Werte dazu: 1ms, 2ms, 5ms, 10,20,50ms, 100,200,500ms, 1,2,5, 10,20,50s, 60s. |
| RJMP #x  | 3x     | Jump back              | Sprung zurück um x Befehle, (bzw: zu :label  max 15 Befehle) |
| PAGE #x  | 8x     | page #x                | Setzen des Page Registers 0 ⇐ x < 16                         |
| JMP #x   | 9x     | jump to #x + 16 * page | Springe nach Adresse #x + (16 * page) 0 ⇐ x < 16 (max. 255 Zellen) |
| LOOPX #x | Ax     | loop x                 | Ist X>0: X=X-1; Springe nach #x + (16*page), sonst gehe zum nächsten Befehl |
| LOOPY #x | Bx     | loop y                 | Ist Y>0: Y=Y-1; Springe nach #x + (16*page), sonst gehe zum nächsten Befehl |
| SKIP0    | C0     | Skip if A = 0          | Überspringe nächsten Befehl wenn A = 0                       |
| AGTB     | C1     | Skip if A > B          | Überspringe nächsten Befehl wenn A > B                       |
| ALTB     | C2     | Skip if A < B          | Überspringe nächsten Befehl wenn A < B                       |
| AEQB     | C3     | Skip if A = B          | Überspringe nächsten Befehl wenn A = B                       |
| DEQ#y #x | C4..CB | Skip if Din.X = Y      | Überspringe nächsten Befehl wenn der Eingang #y 0 oder 1 ist. x=1..6 |
| CALL #x  | Dx     | Call #x + (16 * page)  |                                                              |
| RTR      | E0     | Return                 | Zurück springen nach einem Call-Befehl oder am Ende einer Subroutine |
| REST     | EF     | Restart program        |                                                              |
| PEND     | FF     | Program end            |                                                              |
|          |        |                        |                                                              |
|          |        |                        |                                                              |
| CASB #x  | E1..E6 | Call sub #x            | Aufruf der Subroutine x, 1 ⇐ x ⇐ 6                           |
| DFSB #x  | E8..ED | Define sub #x          | Definierung der Subroutine x, 1 ⇐ x ⇐ 6                      |

## Laden und Speichern

| Mnemonic | Code | short description         | Beschreibung                                                 |
| -------- | ---- | ------------------------- | ------------------------------------------------------------ |
| MOV A,#x | 4x   | A = #x                    | Register A wird mit dem festen Wert #x geladen               |
| SWAP     | 50   | Swap A & B                | Tauschen der Register A und B                                |
| MOV X,Y  | 5x   |                           | Kopiert ein Register in ein anderes. Entweder Quelle oder Ziel muss das Register A sein. mögliche Werte: A, B, C, D, X, Y, P<br />Berechnung: <br /><br />Operant op: A=0, B=1, C=2, D=3, X=4, Y=5, P=6 <br />A ist Quelle: x = op, A ist Ziel: x = 8 + op<br />Beispiel MOV A,B: A ist Ziel, B ist Quelle -> mn = 59<br />Beispiel MOV X,A: A ist Quelle, X ist Ziel -> mn = 54<br /> |
| PUSH     | 57   | push A on stack           | Das A Register wird auf dem Stack gelegt                     |
| POP      | 58   | Pop value from stack to A | Wert aus dem Stack in das A Register übertragen              |

## Mathematik

| Mnemonic | Code | short description | Carry | Beschreibung                                                 |
| -------- | ---- | ----------------- | ----- | ------------------------------------------------------------ |
| INC      | 71   | A = A + 1         | x     | Increment A                                                  |
| DEC      | 72   | A = A - 1         | x     | Decrement A                                                  |
| ADD      | 73   | A = A + B         | x     |                                                              |
| SUB      | 74   | A = A - B         | x     |                                                              |
| MUL      | 75   | A = A * B         | x     |                                                              |
| DIV      | 76   | A = A / B         | x     |                                                              |
| AND      | 77   | A = A and B       |       |                                                              |
| OR       | 78   | A = A or B        |       |                                                              |
| XOR      | 79   | A = A xor B       |       |                                                              |
| NOT      | 7A   | A = not A         |       |                                                              |
| MOD      | 7B   | A = A % B         |       |                                                              |
| BYTE     | 7C   | A = A + 16 * B    |       |                                                              |
| BSUBA    | 7D   | A = B - A         | x     |                                                              |
| SHR      | 7E   | A = A >> 1        | x     | Shift right (/2), das Carry Register enthält das rausgeschobene Bit |
| SHL      | 7F   | A = A << 1        | x     | Shift left (*2), das Carry Register enthält das rausgeschobene Bit |

## Ein/Ausgabe

| Mnemonic  | Code   | short description | Beschreibung                                                 |
| --------- | ------ | ----------------- | ------------------------------------------------------------ |
| LDA DIN   | 64     | A = Din           | Register A wird mit dem Wert vom digitalen Eingang geladen. Alle Bits. |
| LDA DINx  | 65..68 | A = Din.x         | Register A wird mit dem Wert vom digitalen Eingang #x geladen. x = 1..4 |
| LDA ADCx  | 69, 6A | A = ADC.x         | Register A wird mit dem Wert vom analogen Eingang #x geladen. x = 1..2 |
| LDA RCx   | 6B, 6C | A = RC.x          | Register A wird mit dem Wert vom RC Empfängereingang #x geladen. x = 1..2 |
| PORT #x   | 1x     | Dout = #x         | Direkte Ausgabe von Wert #x auf den Ausgängen. 0 ⇐ x ⇐ 15.   |
| STA DOUT  | 54     | Dout = A          | Ausgabe von A auf den Ausgängen                              |
| STA DOUTx | 55..58 | Dout.x = A        | Ausgabe von A.x auf den Ausgang x = 1..4                     |
| STA PWMx  | 59, 5A | PWM.x = A         | Der Wert aus dem A Register wird als PWM.x ausgegeben. x = 1,2 |
| STA SRVx  | 5B, 5C | Servo.x = A       | Der Wert aus dem A Register wird als Servo.x ausgegeben. x = 1,2 |

## Byte Befehle

| Mnemonic  | Code   | short description | Beschreibung                                                 |
| --------- | ------ | ----------------- | ------------------------------------------------------------ |
| BLDA ADCx | F0, F1 | A = ADC.x         | Bytewert auslesen eines Analogeinganges x = 1,2              |
| BLDA RCx  | F2, F3 | A = RC.x          | Bytewert auslesen eines Fernsteuerungskanals x = 1,2         |
| BSTA PWMx | F4, F5 | PWM.x = A         | Der Byte Wert aus dem A Register wird als PWM.x ausgegeben. x = 1,2 |
| BSTA SRVx | F6, F7 | Servo.x = A       | Der Byte Wert aus dem A Register wird als Servo.x ausgegeben. x = 1,2 |
| TONE      | F8     | Tone A            | Ausgabe eines Tones nach Midi 32 ⇐ A ⇐ 108                   |



