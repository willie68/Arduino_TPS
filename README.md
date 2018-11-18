# ArduinoSPS

A TPS Variant for the arduino with some major enhancments.

**Version 0.9**
  18.11.2018 WKLA
  * BUGs entfernt. Release.
  10.11.2018 WKLA
  * Implementierung Tone Befehl

**Version 0.8**
  06.11.2018 WKLA
  * Umstellung auf dbgOut
  * Display TM1637 Anbindung

**Version 0.7**
  24.09.2012 WKLA
  * neue Berechnung A = B - A und Swap A,B...
  * Stack auf 16 Bytes berschränkt, wird zu oft gepusht, werden die alten Werte rausgeschoben.

  Basierd auf dem TPS System vom elektronik-labor.
  Erweiterungen:
  * es können bis zu 6 Unterroutinen definiert werden und diese direkt angesprungen werden.
  * neben return gibt's auch einen restart
  * 2 Servoausgänge für übliche RC Servos. (10° Auflösung in Nibble Modus, <1° Auflösung im Bytemodus)
  ACHTUNG: Servo und PWM Ausgänge sind nicht mischbar und können auch nicht gleichzeitig benutzt werden.
  * 2 RC Eingänge (16 Schritte auflösung im nibble Modus, Mitte 8, 255 Schritte im Byte Modus)
  * fkt. auch mit einem ATTiny84 (44 ist leider auf GRund der Programmgröße nicht mehr für den erweiterten Befehlssatz möglich)
  * call stack von bis zu 16 Unterfunktionen
  * neue Register e,f
