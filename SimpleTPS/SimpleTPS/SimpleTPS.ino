/*
  Simple SPS System with Arduino Uno.
*/

// libraries
#include <EEPROM.h>
#include <avr/eeprom.h>

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

// the actual address of the program
word addr;
// page register
word page;
// defining register
byte a, b, c, d;

const byte SAVE_CNT = 1;

word saveaddr[SAVE_CNT];
byte saveCnt;

unsigned long tmpValue;

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

  prgDemoPrg();
  doReset();

  if (digitalRead(SW_PRG) == 0) {
    programMode();
  }
}

void doReset() {
  addr = 0;
  page = 0;
  saveCnt = 0;
  a = 0;
  b = 0;
  c = 0;
  d = 0;
}

/*
  main loop
*/
void loop() {
  // reading from address
  byte value = EEPROM.read(addr);
  // splitting into command
  byte cmd = (value & 0xF0);
  // and data
  byte data = (value & 0x0F);

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
    case IS_A:
      doIsA(data);
      break;
    case A_IS:
      doAIs(data);
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
    default:
      ;
  }
  if (addr > E2END) {
    doReset();
  }
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
    case 9:
      tmpValue = analogRead(ADC_0);
      a = tmpValue / 64; //(Umrechnen auf 4 bit)
      break;
    case 10:
      tmpValue = analogRead(ADC_1);
      a = tmpValue / 64; //(Umrechnen auf 4 bit)
      break;
    default:
      break;
  }
}

/*
  somthing = a;
*/
void doIsA(byte data) {
  switch (data) {
    case 0:
      tmpValue = b;
      b = a;
      a = tmpValue;
      break;
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
      analogWrite(PWM_1, tmpValue);
      break;
    case 10:
      tmpValue = a * 16;
      analogWrite(PWM_2, tmpValue);
      break;
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
      a = ~a;
      break;
    default:
      break;
  }
  a = a & 15;
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
    case 0:
      skip = (a == 0);
      break;
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
    return;
  }
}
