# Arduino_TPS

A TPS Variant for the Arduino and some other options with some major enhancements.

For more information and description see 
http://www.rcarduino.de/doku.php?id=en:arduino:arduinosps

And there is now a print book for this version at amazon 
https://www.amazon.com/dp/1731232535

No Hardware – No Problem. Get started on paper. The new book explains the 40 examples and more 
https://www.amazon.co.uk/gp/product/B08MN15NMQ/ref=dbs_a_def_rwt_bibl_vppi_i24

**---- 18.06.2021 -----**
Renaming the project from ArduinoSPS to Arduino_TPS

**Version 0.13.0**
  16.06.2021

  * feature: ESP32 implementation: https://wkla.no-ip.biz/ArduinoWiki/doku.php?id=arduino:arduinosps:esp32

**Version 0.12.4**
  11.06.2021

  * bug: pop not working
  * bug: a = not a not working

**Version 0.12.3**
  10.06.2021

  * adding auto programming feature for the SPS Emulator

**Version 0.12.2**
  07.06.2021

  * bug with servo in 4-bit mode, evaluate the full 8 bit.

**Version 0.12.1**
  03.09.2019

  * changing the variable names in debug mode
  * release

**Version 0.12**
  27.01.2019
  * Release

  * automatically adding demo program to new mcu

  11.01.2018

  * some refactoring

  07.01.2018
  * programming: 1/2 duty cycle for 0 values in address display

**Version 0.11**
  17.12.2018
  * adding Shift left and shift right to register A


**Version 0.10**
  9.12.2018
  * Release

  7.12.2018
  * new define for serial programming

  18.11.2018 WKLA
  * new standard programming mode
    I added a new programming mode for the default programming, because i thing the old one was a little bit clumsy.
    The new one has a nicer interface, as you now always know where you are.
    Starting with PRG pushed after Reset.
    as a result, all LEDs will shortly blink
    now you are in programming mode.
  * the D1 LED will blink
  * the higher nibble of the address will be shown
  * the D2 LED will blink
  * the lower nibble of the address will be shown
  * the D3 LED will blink
  * the command part (high nibble) will be shown
  * with SEL you can step thru all commands
  * PRG will save the command
  * the D4 LED will blink 
  * the data part (low nibble) will be shown
  * with SEL you can step thru all datas
  * PRG will save the data
  * if the new value has been changed, all LEDs will flash as the byte will be written to the EEPROM
  * address will be increased and now it will start with blinking of the D1 LED
  * 
  * To leave the programming simply push reset.

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
