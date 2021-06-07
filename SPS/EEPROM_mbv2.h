/* this is the eeprom abstraction for the microbit v2. 
 *  the eeprom will be mapped to a simple file in the file system, later.
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
  if ((addr >=0) && (addr < STORESIZE)) {
    return program[addr];
  }
  return 0xFF;
}

void writebyte(int addr, byte value) { 
  if ((addr >=0) && (addr < STORESIZE)) {
    program[addr] = value;
  }
}

void store() {
}

#else

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
