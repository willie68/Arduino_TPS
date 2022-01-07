/*
  Simple SPS System with Arduino Uno.
*/

// libraries for access to eeprom
#include <EEPROM.h>
#include <avr/eeprom.h>

// defining the hardware connections
// Outputs
const byte Dout_0 = 4;
const byte Dout_1 = 5;
const byte Dout_2 = 6;
const byte Dout_3 = 7;

// Inputs
const byte Din_0 = 0;
const byte Din_1 = 1;
const byte Din_2 = 2;
const byte Din_3 = 3;

// Specials
const byte ADC_0 = 0; 
const byte ADC_1 = 1; 
const byte PWM_1 = 9;
const byte PWM_2 = 10;

// Buttons
const byte SW_PRG = 8;  // (S2)
const byte SW_SEL = 11; // (S1)

// the defined Commands
const byte PORT = 0x10; // directly output to LEDs
const byte DELAY = 0x20; // wait just a little bit
const byte JUMP_BACK = 0x30; // jump back to address
const byte SET_A = 0x40; // directly set a value into the A register
const byte IS_A = 0x50; // put the value of A into another 
const byte A_IS = 0x60; // put something into A
const byte CALC = 0x70; // calculations with A

// the registers
byte a, b, c, d;

// the actual address pointer of the program
word addr;

// this will be called only once after controller reset
void setup() {
  // set all TPS Outputs to output mode
  pinMode(Dout_0, OUTPUT);
  pinMode(Dout_1, OUTPUT);
  pinMode(Dout_2, OUTPUT);
  pinMode(Dout_3, OUTPUT);

  // analog output
  pinMode(PWM_1, OUTPUT);
  pinMode(PWM_2, OUTPUT);

  // digital input with pullup
  pinMode(Din_0, INPUT_PULLUP);
  pinMode(Din_1, INPUT_PULLUP);
  pinMode(Din_2, INPUT_PULLUP);
  pinMode(Din_3, INPUT_PULLUP);

  // program the demo program
  prgDemoPrg();

  // reset the program
  doReset();

// starting the programming mode on startup, if SW_PRG is pressed
  if (digitalRead(SW_PRG) == 0) {
    programMode();
  }
}

// reset all
void doReset() {
  // reset address pointer
  addr = 0;
  // all LED off
  doPort(0);
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

  // switch to the right method for the command
  switch (cmd) {
    // the port command
    case PORT:
      doPort(data);
      break;
    // the delay command
    case DELAY:
      doDelay(data);
      break;
    // the jump command
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
    // if there is an unkown command, reset the tps
    default:
      doReset();
      return;
  }
  // increment address pointer
  addr = addr + 1;
  // if address is at the end of the eeprom, simply reset the program
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
  somthing = a;
*/
void doIsA(byte data) {
  byte tmpValue;
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
  a = somthing;
*/
void doAIs(byte data) {
  word tmpValue;
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
    case 11:
      a = a % b;
      break;
    case 13:
      a = b - a ;
      break;
    case 14:
      a = a >> 1;
      break;
    case 15:
      a = a << 1;
      break;
    default:
      break;
  }
  a = a & 15;
}
