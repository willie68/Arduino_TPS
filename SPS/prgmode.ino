/*
  entering the programming mode
*/

#define BLINK_DELAY 250
#define SHOW_DELAY 500
#define KEY_DELAY 250

enum PROGRAMMING_MODE {ADDRESS, COMMAND, DATA};

PROGRAMMING_MODE prgMode;

void programMode() {
  // checking if advance programmer board connected?
#ifdef SPS_ENHANCEMENT
  if (digitalRead(SW_SEL) == 0) {
    advancePrg();
  }
  else {
#endif
    dbgOutLn("PrgMode");
    blinkAll();
    while (digitalRead(SW_PRG) == 0) {
      // waiting for PRG to release
    }
    prgMode = ADDRESS;
    addr = 0;
    do {
      blinkD1();
      dbgOut("Adr:");
      dbgOutLn(addr);
      // LoNibble Adresse anzeigen
      doPort(addr);
      delay(SHOW_DELAY);

      blinkD2();
      // HiNibble Adresse anzeigen
      data = (addr & 0xf0) >> 4;                                  //Adresse anzeigen
      doPort(data);
      delay(SHOW_DELAY);

      byte Eebyte = EEPROM.read(addr);
      data = Eebyte & 15;
      com = Eebyte >> 4;

      blinkD3();
      prgMode = COMMAND;
      dbgOutLn("com");
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
      dbgOutLn("dat");
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
#ifdef SPS_ENHANCEMENT
  }
#endif
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
