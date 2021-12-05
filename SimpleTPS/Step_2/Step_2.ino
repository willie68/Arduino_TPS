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

// Buttons
const byte SW_PRG = 8;  // (S2)
const byte SW_SEL = 11; // (S1)

// the defined Commands
const byte PORT = 0x10; // directly output to LEDs
const byte DELAY = 0x20; // wait just a little bit
const byte JUMP_BACK = 0x30; // jump back to address

// the actual address pointer of the program
word addr;

// this will be called only once after controller reset
void setup() {
  // set all TPS Outputs to output mode
  pinMode(Dout_0, OUTPUT);
  pinMode(Dout_1, OUTPUT);
  pinMode(Dout_2, OUTPUT);
  pinMode(Dout_3, OUTPUT);

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
  addr = addr - data;
}
