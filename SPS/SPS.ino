/*
  SPS System mit dem Arduino.
  Version 0.12
  27.01.2019
  - adding demo program,
  11.01.2018
  - some refactoring
  
  07.01.2018
  - programming: 1/2 duty cycle for 0 values in address display

  Version 0.11
  17.12.2018
  - adding Shift left and shift right to register A

  Version 0.10
  7.12.2018
  - new define for serial programming

  18.11.2018 WKLA
  - new standard programming mode
  I added a new programming mode for the default programming, because i thing the old one was a little bit clumsy.
  the new one has a nicer interface, as you now always know where you are.
  Starting with PRG pushed after Reset.
  as a result, all LEDs will shortly blink
  now you are in programming mode.
    the D1 LED will blink
    the higher nibble of the address will be shown
    the D2 LED will blink
    the lower nibble of the address will be shown
    the D3 LED will blink
    the command part (high nibble) will be shown
    with SEL you can step thru all commands
    PRG will save the command
    the D4 LED will blink
    the data part (low nibble) will be shown
    with SEL you can step thru all datas
    PRG will save the data
    if the new value has been changed, all LEDs will flash as the byte will be written to the EEPROM
    address will be increased and now it will start with blinking of the D1 LED

    To leave the programming simply push reset.

  Version 0.9
  18.11.2018 WKLA
  - BUGs entfernt. Release.
  10.11.2018 WKLA
  - Implementierung Tone Befehl

  Version 0.8
  06.11.2018 WKLA
  - Umstellung auf dbgOut
  - Display TM1637 Anbindung

  Version 0.7
  24.09.2012 WKLA
  - neue Berechnung A = B - A und Swap A,B...
  - Stack auf 16 Bytes berschränkt, wird zu oft gepusht, werden die alten Werte rausgeschoben.

  Basierd auf dem TPS System vom elektronik-labor.
  Erweiterungen:
  - es können bis zu 6 Unterroutinen definiert werden und diese direkt angesprungen werden.
  - neben return gibt's auch einen restart
  - 2 Servoausgänge für übliche RC Servos. (10° Auflösung in Nibble Modus, <1° Auflösung im Bytemodus)
  ACHTUNG: Servo und PWM Ausgänge sind nicht mischbar und können auch nicht gleichzeitig benutzt werden.
  - 2 RC Eingänge (16 Schritte auflösung im nibble Modus, Mitte 8, 255 Schritte im Byte Modus)
  - fkt. auch mit einem ATTiny84 (44 ist leider auf GRund der Programmgröße nicht mehr für den erweiterten Befehlssatz möglich)
  - call stack von bis zu 16 Unterfunktionen
  - neue Register e,f
*/

/*
 * Here are the defines used in this software to control special parts of the implementation
 * #define SPS_USE_DISPLAY: using a external TM1637 Display for displaying address and data at one time
 * #define SPS_RECEIVER: using a RC receiver input
 * #define SPS_ENHANCEMENT: all of the other enhancments
 * #define SPS_SERVO: using servo outputs
 * #define SPS_TONE: using a tone output
 * #define SPS_SERIAL_PRG: activates the serial programming feature
 */
// Program im Debugmodus kompilieren, dann werden zus. Ausgaben auf die serielle Schnittstelle geschrieben.
//#define debug

// defining different hardware platforms
#ifdef __AVR_ATmega328P__
//#define SPS_USE_DISPLAY
//#define SPS_RECEIVER
//#define SPS_ENHANCEMENT
//#define SPS_SERIAL_PRG
//#define SPS_SERVO
//#define SPS_TONE
#endif

#ifdef __AVR_ATtiny84__
#define SPS_ENHANCEMENT
#define SPS_SERIAL_PRG
#define SPS_SERVO
//#define SPS_TONE
#endif

#ifdef __AVR_ATtiny861__
#define SPS_RCRECEIVER
#define SPS_ENHANCEMENT
#define SPS_SERIAL_PRG
//#define SPS_SERVO
#define SPS_TONE
#endif

#ifdef __AVR_ATtiny4313__
#define SPS_RCRECEIVER
#endif

// libraries
#include <debug.h>
#include <makros.h>
#include <EEPROM.h>
#include <avr/eeprom.h>

#ifdef SPS_SERVO
#include <Servo.h>
#endif

#ifdef SPS_ENHANCEMENT
#include <avdweb_Switch.h>
#endif

#ifdef SPS_TONE
#include "notes.h"
#endif

#include "hardware.h"

// Commands
const byte PORT = 0x10;
const byte DELAY = 0x20;
const byte JUMP_BACK = 0x30;
const byte SET_A = 0x40;
const byte IS_A = 0x50;
const byte A_IS = 0x60;
const byte CALC = 0x70;
const byte PAGE = 0x80;
const byte JUMP = 0x90;
const byte C_COUNT = 0xA0;
const byte D_COUNT = 0xB0;
const byte SKIP_IF = 0xC0;
const byte CALL = 0xD0;
const byte CALL_SUB = 0xE0;
const byte CMD_BYTE = 0xF0;

// debouncing with 100ms
const byte DEBOUNCE = 100;

// sub routines
const byte subCnt = 7;
word subs[subCnt];

// the actual address of the program
word addr;
// page register
word page;
// defining register
byte a, b, c, d;
#ifdef SPS_ENHANCEMENT
byte e, f;
#endif

#ifdef SPS_ENHANCEMENT
const byte SAVE_CNT = 16;
#else
const byte SAVE_CNT = 1;
#endif

word saveaddr[SAVE_CNT];
byte saveCnt;

#ifdef SPS_ENHANCEMENT
byte stack[SAVE_CNT];
byte stackCnt;
#endif

unsigned long tmpValue;

#ifdef SPS_SERVO
Servo servo1;
Servo servo2;
#endif

byte prog = 0;
byte data = 0;
byte com = 0;

void setup() {
  pinMode(Dout_0, OUTPUT);
  pinMode(Dout_1, OUTPUT);
  pinMode(Dout_2, OUTPUT);
  pinMode(Dout_3, OUTPUT);

  pinMode(PWM_1, OUTPUT);
  pinMode(PWM_2, OUTPUT);

  pinMode(Din_0, INPUT_PULLUP);
  pinMode(Din_1, INPUT_PULLUP);
  pinMode(Din_2, INPUT_PULLUP);
  pinMode(Din_3, INPUT_PULLUP);

  pinMode(SW_PRG, INPUT_PULLUP);
  pinMode(SW_SEL, INPUT_PULLUP);

  digitalWrite(Dout_0, 1);
  delay(1000);
  digitalWrite(Dout_0, 0);
#ifdef SPS_USE_DISPLAY
  initDisplay();
#endif

  // Serielle Schnittstelle einstellen
#ifndef __AVR_ATtiny84__
  initDebug();
#endif

  prgDemoPrg();
  doReset();

  if (digitalRead(SW_PRG) == 0) {
    programMode();
  }
#ifdef SPS_ENHANCEMENT
  pinMode(LED_BUILTIN, OUTPUT);
#endif

#ifdef SPS_SERIAL_PRG
  if (digitalRead(SW_SEL) == 0) {
    serialPrg();
  }
#endif
}

void doReset() {
  dbgOutLn("Reset");
#ifdef SPS_SERVO
  servo1.detach();
  servo2.detach();
#endif

  for (int i = 0; i < subCnt; i++) {
    subs[i] = 0;
  }

  readProgram();

  addr = 0;
  page = 0;
  saveCnt = 0;
  a = 0;
  b = 0;
  c = 0;
  d = 0;
#ifdef SPS_ENHANCEMENT
  e = 0;
  f = 0;
#endif
}

/*
  getting all addresses of sub programms
*/
void readProgram() {
  dbgOutLn("Read program");
  word addr = 0;
  for ( addr = 0; addr <= E2END; addr++) {
    byte value = EEPROM.read(addr);

#ifdef debug
    dbgOut2(value, HEX);
    if (((addr + 1) % 16) == 0) {
      dbgOutLn();
    }
    else {
      dbgOut(",");
    }
#endif

    if (value == 0xFF) {
      // ende des Programms
      break;
    }
    byte cmd = (value & 0xF0);
    byte data = (value & 0x0F);

    dbgOut("(");
    dbgOut2(cmd, HEX);
    dbgOut2(data, HEX);
    dbgOut(")");

    if (cmd == CALL_SUB) {
      if (data >= 8) {
        data = data - 8;
        subs[data] = addr + 1;
      }
    }
#ifdef SPS_SERVO
    if ((cmd == IS_A) && (data == 0x0B)) {
      if (!servo1.attached()) {
        dbgOutLn("attach Srv1");
        servo1.attach(SERVO_1);
      }
    } else if ((cmd == CMD_BYTE) && (data == 0x06)) {
      if (!servo1.attached()) {
        dbgOutLn("attach Srv1");
        servo1.attach(SERVO_1);
      }
    } else if ((cmd == IS_A) && (data == 0x0C)) {
      if (!servo2.attached()) {
        dbgOutLn("attach Srv2");
        servo2.attach(SERVO_2);
      }
    } else if ((cmd == CMD_BYTE) && (data == 0x07)) {
      if (!servo2.attached()) {
        dbgOutLn("attach Srv2");
        servo2.attach(SERVO_2);
      }
    }
#endif
  }
  dbgOutLn();
}

/*
  main loop
*/
void loop() {
  byte value = EEPROM.read(addr);
  byte cmd = (value & 0xF0);
  byte data = (value & 0x0F);

  debugOutputRegister();

  addr = addr + 1;
  switch (cmd) {
    case PORT:
      doPort(data);
      break;
    case DELAY:
      doDelay(data);
      break;
    case JUMP_BACK:
      doJumpBack(data);
      break;
    case SET_A:
      doSetA(data);
      break;
    case A_IS:
      doAIs(data);
      break;
    case IS_A:
      doIsA(data);
      break;
    case CALC:
      doCalc(data);
      break;
    case PAGE:
      doPage(data);
      break;
    case JUMP:
      doJump(data);
      break;
    case C_COUNT:
      doCCount(data);
      break;
    case D_COUNT:
      doDCount(data);
      break;
    case SKIP_IF:
      doSkipIf(data);
      break;
    case CALL:
      doCall(data);
      break;
    case CALL_SUB:
      doCallSub(data);
      break;
    case CMD_BYTE:
      doByte(data);
      break;
    default:
      ;
  }
  if (addr > E2END) {
    doReset();
  }
}

void debugOutputRegister() {
  dbgOut2(addr, HEX); dbgOut(":"); dbgOut2(prog, HEX); dbgOut(",");
  dbgOut2(com, HEX); dbgOut(","); dbgOut2(data, HEX); dbgOut(",a:");
  dbgOut2(a, HEX); dbgOut(","); dbgOut2(b, HEX); dbgOut(",");
  dbgOut2(c, HEX); dbgOut(","); dbgOut2(d, HEX); dbgOut(",");
  dbgOut2(e, HEX); dbgOut(","); dbgOut2(f, HEX); dbgOutLn();
}

/*
  output to port
*/
void doPort(byte data) {
  digitalWrite(Dout_0, (data & 0x01) > 0);
  digitalWrite(Dout_1, (data & 0x02) > 0);
  digitalWrite(Dout_2, (data & 0x04) > 0);
  digitalWrite(Dout_3, (data & 0x08) > 0);
}

/*
  delay in ms
*/
void doDelay(byte data) {
  dbgOut("dly: ");
  dbgOutLn2(data, HEX);

  switch (data) {
    case 0:
      delay(1);
      break;
    case 1:
      delay(2);
      break;
    case 2:
      delay(5);
      break;
    case 3:
      delay(10);
      break;
    case 4:
      delay(20);
      break;
    case 5:
      delay(50);
      break;
    case 6:
      delay(100);
      break;
    case 7:
      delay(200);
      break;
    case 8:
      delay(500);
      break;
    case 9:
      delay(1000);
      break;
    case 10:
      delay(2000);
      break;
    case 11:
      delay(5000);
      break;
    case 12:
      delay(10000);
      break;
    case 13:
      delay(20000);
      break;
    case 14:
      delay(30000);
      break;
    case 15:
      delay(60000);
      break;
    default:
      break;
  }
}

/*
  jump relative back
*/
void doJumpBack(byte data) {
  addr = addr - data - 1;
}

/*
  a = data
*/
void doSetA(byte data) {
  a = data;
}

/*
  a = somthing;
*/
void doAIs(byte data) {
  switch (data) {
    case 1:
      a = b;
      break;
    case 2:
      a = c;
      break;
    case 3:
      a = d;
      break;
    case 4:
      a = digitalRead(Din_0) + (digitalRead(Din_1) << 1) + (digitalRead(Din_2) << 2) + (digitalRead(Din_3) << 3);
      break;
    case 5:
      a = digitalRead(Din_0);
      break;
    case 6:
      a = digitalRead(Din_1);
      break;
    case 7:
      a = digitalRead(Din_2);
      break;
    case 8:
      a = digitalRead(Din_3);
      break;
#ifndef __AVR_ATtiny4313__
    case 9:
      tmpValue = analogRead(ADC_0);
      a = tmpValue / 64; //(Umrechnen auf 4 bit)
      break;
    case 10:
      tmpValue = analogRead(ADC_1);
      a = tmpValue / 64; //(Umrechnen auf 4 bit)
      break;
#else
    case 9:
      a = digitalRead(ADC_0);
      break;
    case 10:
      a = digitalRead(ADC_1);
      break;
#endif
#ifdef SPS_RCRECEIVER
    case 11:
      tmpValue = pulseIn(RC_0, HIGH, 100000);
      if (tmpValue < 1000) {
        tmpValue = 1000;
      }
      if (tmpValue > 2000) {
        tmpValue = 2000;
      }
      a = (tmpValue - 1000) / 64; //(Umrechnen auf 4 bit)
      dbgOut("RC1:");
      dbgOut(tmpValue);
      dbgOut("=");
      dbgOutLn(a);
      break;
    case 12:
      tmpValue = pulseIn(RC_1, HIGH, 100000);
      if (tmpValue < 1000) {
        tmpValue = 1000;
      }
      if (tmpValue > 2000) {
        tmpValue = 2000;
      }
      a = (tmpValue - 1000) / 64; //(Umrechnen auf 4 bit)
      dbgOut("RC2:");
      dbgOut(tmpValue);
      dbgOut("=");
      dbgOutLn(a);
      break;
#endif
#ifdef SPS_ENHANCMENT
    case 13:
      a = e;
      break;
    case 14:
      a = f;
      break;
    case 15:
      if (stackCnt > 0) {
        stackCnt -= 1;
        a = stack[stackCnt];
      } else {
        a = 0;
      }
      break;
#endif
    default:
      break;
  }
}

/*
  somthing = a;
*/
void doIsA(byte data) {
  switch (data) {
#ifdef SPS_ENHANCEMENT
    case 0:
      swap(a, b, byte);
      break;
#endif
    case 1:
      b = a;
      break;
    case 2:
      c = a;
      break;
    case 3:
      d = a;
      break;
    case 4:
      doPort(a);
      break;
    case 5:
      digitalWrite(Dout_0, (a & 0x01) > 0);
      break;
    case 6:
      digitalWrite(Dout_1, (a & 0x01) > 0);
      break;
    case 7:
      digitalWrite(Dout_2, (a & 0x01) > 0);
      break;
    case 8:
      digitalWrite(Dout_3, (a & 0x01) > 0);
      break;
    case 9:
      tmpValue = a * 16;
      dbgOut("PWM1:");
      dbgOutLn(tmpValue);
      analogWrite(PWM_1, tmpValue);
      break;
    case 10:
      tmpValue = a * 16;
      dbgOut("PWM2:");
      dbgOutLn(tmpValue);
      analogWrite(PWM_2, tmpValue);
      break;
#ifdef SPS_SERVO
    case 11:
      if (servo1.attached()) {
        tmpValue = (a * 10) + 10;
        dbgOut("Srv1:");
        dbgOutLn(tmpValue);
        servo1.write(tmpValue);
      }
      break;
    case 12:
      if (servo2.attached()) {
        tmpValue = (a * 10) + 10;
        dbgOut("Srv2:");
        dbgOutLn(tmpValue);
        servo2.write(tmpValue);
      }
      break;
#endif
#ifdef SPS_ENHANCEMENT
    case 13:
      e = a;
      break;
    case 14:
      f = a;
      break;
    case 15:
      if (stackCnt < SAVE_CNT) {
        stack[stackCnt] = a;
        stackCnt += 1;
      }
      else {
        for (int i = 1; i <= SAVE_CNT; i++) {
          stack[i - 1] = stack[i];
        }
        stack[stackCnt] = a;
      }
      break;
#endif
    default:
      break;
  }
}

/*
  calculations
*/
void doCalc(byte data) {
  switch (data) {
    case 1:
      a = a + 1;
      break;
    case 2:
      a = a - 1;
      break;
    case 3:
      a = a + b;
      break;
    case 4:
      a = a - b;
      break;
    case 5:
      a = a * b;
      break;
    case 6:
      a = a / b;
      break;
    case 7:
      a = a & b;
      break;
    case 8:
      a = a | b;
      break;
    case 9:
      a = a ^ b;
      break;
    case 10:
      a = !a;
      break;
#ifdef SPS_ENHANCEMENT
    case 11:
      a = a % b;
      break;
    case 12:
      a = a + 16 * b;
      break;
    case 13:
      a = b - a;
      break;
    case 14:
      a = a >> 1;
      break;
    case 15:
      a = a << 1;
      break;
#endif
    default:
      break;
  }
#ifndef SPS_ENHANCEMENT
  a = a & 15;
#endif
}

/*
  setting page
*/
void doPage(byte data) {
  page = data * 16;
}

/*
  jump absolute
*/
void doJump(byte data) {
#ifdef debug
  dbgOut("J");
  dbgOut2(page, HEX);
  dbgOutLn2(data, HEX);
#endif
  addr = page + data;
}

/*
  counting with c register
*/
void doCCount(byte data) {
  if (c > 0) {
    c -= 1;
    c = c & 0x0F;
    doJump(data);
  }
}

/*
  counting with d register
*/
void doDCount(byte data) {
  if (d > 0) {
    d -= 1;
    d = d & 0x0F;
    doJump(data);
  }
}

/*
  simple condition = true, skip next command
*/
void doSkipIf(byte data) {
  bool skip = false;
  switch (data) {
#ifdef SPS_ENHANCEMENT
    case 0:
      skip = (a == 0);
      break;
#endif
    case 1:
      skip = (a > b);
      break;
    case 2:
      skip = (a < b);
      break;
    case 3:
      skip = (a == b);
      break;
    case 4:
      skip = digitalRead(Din_0);
      break;
    case 5:
      skip = digitalRead(Din_1);
      break;
    case 6:
      skip = digitalRead(Din_2);
      break;
    case 7:
      skip = digitalRead(Din_3);
      break;
    case 8:
      skip = !digitalRead(Din_0);
      break;
    case 9:
      skip = !digitalRead(Din_1);
      break;
    case 10:
      skip = !digitalRead(Din_2);
      break;
    case 11:
      skip = !digitalRead(Din_3);
      break;
    case 12:
      skip = !digitalRead(SW_PRG);
      break;
    case 13:
      skip = !digitalRead(SW_SEL);
      break;
    case 14:
      skip = digitalRead(SW_PRG);
      break;
    case 15:
      skip = digitalRead(SW_SEL);
      break;
    default:
      break;
  }
  if (skip) {
    addr += 1;
  }
}

/*
  calling a subroutine
*/
void doCall(byte data) {
  saveaddr[saveCnt] = addr;
  saveCnt++;
  addr = page + data;
}

/*
  calling a subroutine, calling return and restart
*/
void doCallSub(byte data) {
  if (data == 0) {
    if (saveCnt < 0) {
      doReset();
      return;
    }
    saveCnt -= 1;
    addr = saveaddr[saveCnt];
    dbgOut("r:");
    dbgOutLn(addr);
    return;
  }
#ifdef SPS_ENHANCEMENT
  if (data <= 7) {
    // call subroutine number
    doCall(addr);
    addr = subs[data - 1];
    dbgOut("c:");
    dbgOutLn(addr);
    return;
  }
  if (data == 0x0f) {
    doReset();
  }
#endif
}

/*
  calling a byte methods
*/
void doByte(byte data) {
#ifdef SPS_ENHANCEMENT
  dbgOut("B ");
  switch (data) {
    case 0:
      tmpValue = analogRead(ADC_0);
      a = tmpValue >> 2; //(Umrechnen auf 8 bit)
      break;
    case 1:
      tmpValue = analogRead(ADC_1);
      a = tmpValue >> 2; //(Umrechnen auf 8 bit)
      break;
#ifdef SPS_RCRECEIVER
    case 2:
      tmpValue = pulseIn(RC_0, HIGH, 100000);
      if (tmpValue < 1000) {
        tmpValue = 1000;
      }
      if (tmpValue > 2000) {
        tmpValue = 2000;
      }
      a = (tmpValue - 1000) / 4; //(Umrechnen auf 4 bit)
      dbgOut("RC1:");
      dbgOut(tmpValue);
      dbgOut("=");
      dbgOutLn(a);
      break;
    case 3:
      tmpValue = pulseIn(RC_1, HIGH, 100000);
      if (tmpValue < 1000) {
        tmpValue = 1000;
      }
      if (tmpValue > 2000) {
        tmpValue = 2000;
      }
      a = (tmpValue - 1000) / 4; //(Umrechnen auf 4 bit)
      dbgOut("RC2:");
      dbgOut(tmpValue);
      dbgOut("=");
      dbgOutLn(a);
      break;
#endif
    case 4:
      tmpValue = a;
      dbgOut("PWM1:");
      dbgOutLn(a);
      analogWrite(PWM_1, a);
      break;
    case 5:
      tmpValue = a;
      dbgOut("PWM2:");
      dbgOutLn(a);
      analogWrite(PWM_2, a);
      break;
#ifdef SPS_SERVO
    case 6:
      if (servo1.attached()) {
        dbgOut("Srv1:");
        tmpValue = map(a, 0, 255, 0, 180);
        dbgOutLn(tmpValue);
        servo1.write(tmpValue);
      }
      break;
    case 7:
      if (servo2.attached()) {
        dbgOut("Srv2:");
        tmpValue = map(a, 0, 255, 0, 180);
        dbgOutLn(tmpValue);
        servo2.write(tmpValue);
      }
      break;
#endif
#ifdef SPS_TONE
    case 8:
      if (a == 0) {
        dbgOutLn("Tone off");
        noTone(PWM_2);
      } else {
        if (between(a, MIDI_START, MIDI_START + MIDI_NOTES)) {
          word frequenz = pgm_read_word(a - MIDI_START + midiNoteToFreq);
          dbgOut("Tone on: midi ");
          dbgOut2(a, DEC);
          dbgOut(", ");
          dbgOut2(frequenz, DEC);
          dbgOutLn("Hz");
          tone(PWM_2, frequenz);
        }
      }
      break;
#endif
#ifdef __AVR_ATmega328P__
    case 13:
      digitalWrite(LED_BUILTIN, 0);
      break;
    case 14:
      digitalWrite(LED_BUILTIN, 1);
      break;
#endif
  }
#endif
}
