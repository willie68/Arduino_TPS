/*
  entering the programming mode
*/

#define BLINK_DELAY 500
#define SHOW_DELAY 1000
#define KEY_DELAY 250
#define ADDR_LOOP 50

const byte demoPrg[] = { 0x4F, 0x59, 0x1F, 0x29, 0x10, 0x29, 0x5A, 0x40,
                         0x59, 0x64, 0x54, 0x29, 0x4F, 0x59, 0x10, 0xCD,
                         0x11, 0x28, 0xCC, 0x18, 0x28, 0x4F, 0x59, 0x5A,
                         0x72, 0x26, 0xC0, 0x35, 0x80, 0x90, 0xFF
                       };


enum PROGRAMMING_MODE {ADDRESS, COMMAND, DATA};

PROGRAMMING_MODE prgMode;

void prgDemoPrg() {
  byte value = EEPROM.read(0);
  if (value == 0xFF) {
    value = EEPROM.read(1);
    if (value == 0xFF) {
      for (byte i = 0; i < sizeof(demoPrg); i++) {
        EEPROM.write(i, demoPrg[i]);
      }
    }
  }
}

void programMode() {
  // checking if advance programmer board connected?
  doPort(0x08);
  while (digitalRead(SW_PRG) == 0) {
    // waiting for PRG to release
  }
  blinkAll();
  prgMode = ADDRESS;
  addr = 0;
  do {
    blinkD1();
    // LoNibble Adresse anzeigen
    doAddr(addr);
    //delay(SHOW_DELAY);

    blinkD2();
    // HiNibble Adresse anzeigen
    data = (addr & 0xf0) >> 4;                                  //Adresse anzeigen
    doAddr(data);
    //delay(SHOW_DELAY);

    byte Eebyte = EEPROM.read(addr);
    data = Eebyte & 15;
    com = Eebyte >> 4;

    blinkD3();
    prgMode = COMMAND;
    doPort(com); //show command

    do {
      if (digitalRead(SW_SEL) == 0) {
        delay(KEY_DELAY);
        com += 1;
        com = com & 0x0F;
        doPort(com);
      }
    }
    while (digitalRead(SW_PRG) == 1);
    delay(DEBOUNCE);

    blinkD4();
    prgMode = DATA;
    doPort(data); //show data

    do {
      if (digitalRead(SW_SEL) == 0) {
        delay(KEY_DELAY);
        data += 1;
        data = data & 0x0F;
        doPort(data);
      }
    }
    while (digitalRead(SW_PRG) == 1); // S2 = 1
    delay(DEBOUNCE);

    byte newValue = (com << 4) + data;
    if (newValue != Eebyte) {
      EEPROM.write(addr, newValue); //           Writeeeprom Eebyte , Addr
      blinkAll();
    }
    addr += 1;
  }
  while (true);
}

void blinkAll() {
  blinkNull();
  doPort(0x0F);
  delay(BLINK_DELAY);
}

void blinkD1() {
  blinkNull();
  doPort(0x01);
  delay(BLINK_DELAY);
  blinkNull();
}

void blinkD2() {
  blinkNull();
  doPort(0x02);
  delay(BLINK_DELAY);
  blinkNull();
}

void blinkD3() {
  blinkNull();
  doPort(0x04);
  delay(BLINK_DELAY);
  blinkNull();
}

void blinkD4() {
  blinkNull();
  doPort(0x08);
  delay(BLINK_DELAY);
  blinkNull();
}

void blinkNull() {
  doPort(0x00);
  delay(BLINK_DELAY);
}

void doAddr(byte value) {
  for (byte i = ADDR_LOOP; i > 0; i--) {
    doPort(value);
    delay(19);
    doPort(0x0F);
    delay(1);
  }
}
