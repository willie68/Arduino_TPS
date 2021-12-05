/*
  entering the programming mode
*/

// simple blink program
const byte demoPrg[] = { 
  0x11, // Dout=1
  0x29, // 1000ms
  0x18, // Dout=8
  0x29, // 1000ms
  0x34, // Addr = addr -4
  0xFF  // EoP End of program
};

// just putting this demo program into the eeprom
void prgDemoPrg() {
  // only need to program the demo program if there is nothing in the eeprom (value on address 0 is FF)
  byte value = EEPROM.read(0);
  if (value == 0xFF) {
    for (byte i = 0; i < sizeof(demoPrg); i++) {
      EEPROM.write(i, demoPrg[i]);
    }
  }
}

byte data = 0;
byte com = 0;

byte prgMode;

void programMode() {
  // waiting for PRG to release
  prgMode = 0;
  addr = 0;
  do {
    // show Low Nibble of Addresse
    doPort(addr & 0x0F);
    delay(300);
    doPort(0);
    delay(200);

    byte Eebyte = EEPROM.read(addr);
    data = Eebyte & 15;
    com = Eebyte >> 4;
    doPort(com); // show command
    do {
    }
    while (digitalRead(SW_PRG) == 0);
    delay(50);
    prgMode = 1; // Phase 1: show command

    do {
      if (digitalRead(SW_SEL) == 0) {
        if (prgMode == 1) {
          prgMode = 2;
          com = 15;
        }
        if (prgMode == 2) { // Phase 2: command changed
          com += 1;
          com = com & 0x0F;
          doPort(com);
        }
        if (prgMode == 3) { // Phase 3: command not changed, data changed
          prgMode = 5;
          data = 15;
        }
        if (prgMode == 4) { // phase 4: command and data changed
          prgMode = 5;
          data = 15;
        }
        if (prgMode == 5) { // phase 5 data changed
          data += 1;
          data = data & 0x0F;
          doPort(data);
        }
        delay(50);
        do {
        }
        while (digitalRead(SW_SEL) == 0);
        delay(50);
      }

      if (digitalRead(SW_PRG) == 0) {
        if (prgMode == 3) {
          prgMode = 7; // only shown not changed
        }
        if (prgMode == 1) {
          doPort(data);
          prgMode = 3;
        }
        if (prgMode == 4) {
          doPort(data);
          prgMode = 6;
        }
        if (prgMode == 2) {
          doPort(data);
          prgMode = 4;
        }
        if (prgMode == 6) { // only command changed
          data = data & 0x0F;
          Eebyte = (com << 4) + data;
          EEPROM.write(addr, Eebyte);
          doPort(0);
          delay(600);
          addr += 1;
          prgMode = 0;
        }
        if (prgMode == 5) { // data changed
          data = data & 0x0F;
          Eebyte = (com << 4) + data;
          EEPROM.write(addr, Eebyte);
          doPort(0);
          delay(600);
          addr += 1;
          prgMode = 0;
        }
        if (prgMode == 7) { // only command changed
          addr += 1;
          prgMode = 0;
        }
        delay(50);
        do {
        }
        while (digitalRead(SW_PRG) == 0);
        delay(50);
      }
    }
    while (digitalRead(SW_PRG) == 1);
  }
  while (true);
}
