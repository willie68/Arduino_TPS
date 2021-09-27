/* this is the eeprom abstraction for the microbit v2.
    the eeprom will be mapped to a simple file in the file system, later.
*/
#ifdef _MICROBIT_V2_

const int STORESIZE = 256;
byte program[STORESIZE];
bool loaded = false;

void load() {
  loaded = true;
}

byte readbyte(int addr) {
  if (!loaded) {
    load();
  }
  if ((addr >= 0) && (addr < STORESIZE)) {
    return program[addr];
  }
  return 0xFF;
}

void writebyte(int addr, byte value) {
  if ((addr >= 0) && (addr < STORESIZE)) {
    program[addr] = value;
  }
}

void store() {
}
#endif

#ifdef ESP32
#include <EEPROM.h>

const int STORESIZE = 1024;
byte program[STORESIZE];
bool loaded = false;

void load() {
  dbgOutLn("load prg from nvs");
  EEPROM.begin(STORESIZE);
  word readed = EEPROM.readBytes(0, program, STORESIZE);
  dbgOut("read:");
  dbgOut(readed);
  dbgOutLn(" bytes");
#ifdef debug1
  for(int i = 0; i < STORESIZE; i++) {
    if ((i % 16) == 0) {
      dbgOutLn();
      dbgOut2(i, HEX);
      dbgOut(": ");
    }
    dbgOut2(program[i], HEX);
    dbgOut(", ");
  }
  dbgOutLn();
#endif
  loaded = true;
}

byte readbyte(int addr) {
  if (!loaded) {
    load();
  }
  if ((addr >= 0) && (addr < STORESIZE)) {
    return program[addr];
  }
  return 0xFF;
}

void writebyte(int addr, byte value) {
  if ((addr >= 0) && (addr < STORESIZE)) {
    program[addr] = value;
  }
}

void store() {
  EEPROM.writeBytes(0, program, STORESIZE);
  EEPROM.commit();
}
#endif

#if defined(__AVR_ATmega328P__) || defined(__AVR_ATtiny84__) || defined(__AVR_ATtiny861__) || defined(__AVR_ATtiny4313__)

#include <EEPROM.h>
#include <avr/eeprom.h>

const int STORESIZE = E2END;

byte readbyte(int addr) {
  return EEPROM.read(addr);
}

void writebyte(int addr, byte value) {
  EEPROM.write(addr, value);
}

void store() {
}

#endif
