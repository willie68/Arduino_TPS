# RCPC, Remote Controlled Programmable Controller

Die RCPC ist eine Steuerungvon Modellbauer für Modellbauer. Sie soll einfache Steuerungsaufgaben ermöglichen, in einer einfachen Steuersprache, auch am Teich oder der Rennstrecke. 

Die Sprache ist einfach und systematisch, reduziert auf das nötigste. 

Als Hardware werden je nach Bedarf verschiedene Controller verwendet. Die Firmware beinhaltet den notwendigen Übersetzer. 

Im Web gibt es eine entsprechende Seite, mit deren Hilfe man die Steuerprogramm in den Steuercode üebrsetzen kann. Der Steuercode kann dann über eine einfache Schnittstelle in den Controller geladen werden. 

Es gibt insgesamt 8 Eingänge und Ausgänge. Es können verschiedene Konfigurationen aktiviert werden.

| Konfiguration | digitale Ausgänge<br />(Dout) | digitale Eingänge<br />(Din) | analoge Ausgänge (PWM)<br />(Aout) | analoge Eingänge (0-5V)<br />(Ain) | Servos<br />(Srv) | RC Eingänge<br />RC | Tone |
| ------------- | ----------------------------- | ---------------------------- | ---------------------------------- | ---------------------------------- | ----------------- | ------------------- | ---- |
| 1             | 6                             | 6                            | 2                                  | 2                                  | 0                 | 0                   | 0    |
| 2             | 6                             | 6                            | 0                                  | 0                                  | 2                 | 2                   | 0    |



# RCPC Arduino Uno (Nano)

| Pin  | Funktion | K 1 | K 2 |      |      |
| ---- | -------- | ---- | ---- | -------- | ---- |
|  D0    | Rx (Prg) |      |     |          |      |
|  D1   | Tx (Prg) |      |      |          |      |
|  D2   |          | Din1 |Din1      |          |      |
|  D3   | PWM | Din2 |Din2      |          |      |
|  D4   |          | Din3 |Din3      |  |      |
|  D5   | PWM | Din4  |Din4      |          |      |
|  D6   | PWM | Dout1 |Dout1      |          |      |
|  D7   |          | Dout2 | Dout2 |          |      |
|  D8   |          | Dout3 |Dout3      |          |      |
| D9   | PWM | Aout1 | Srv1 |   |   |
| D10 | PWM | Aout2 | Srv2 | | |
| D11 | PWM | Dout4 | Dout4 | | |
| D12 | | Dout5 | Dout5 | | |
| D13 | LED | Dout6 | Dout6 | | |
| A0 | | Ain1 |  | | |
| A1 | | Ain2 |  | | |
| A2 | | Din5 |  | | |
| A3 | | Din6 |  | | |
| A4 | SDA |  |  | | |
| A5 | SCL |  |  | | |

