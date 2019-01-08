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
    // light up all LEDs
    doPort(0x08);
    while (digitalRead(SW_PRG) == 0) {
      // waiting for PRG to release
    }
    blinkAll();
    prgMode = ADDRESS;
    addr = 0;
    do {
      blinkD1();
      dbgOut("Adr:");
      dbgOutLn(addr);
      // LoNibble Adresse anzeigen
      doDimPort(addr);

      blinkD2();
      // HiNibble Adresse anzeigen
      data = (addr & 0xf0) >> 4;                                  //Adresse anzeigen
      doDimPort(data);

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

void doDimPort(byte data) {
  long value = millis() + SHOW_DELAY;
  do {
    // on part
    digitalWrite(Dout_0, 1);
    digitalWrite(Dout_1, 1);
    digitalWrite(Dout_2, 1);
    digitalWrite(Dout_3, 1);
    delay(5);

    // off part
    if ((data & 0x01) == 0) {
      digitalWrite(Dout_0, 0);
    }
    if ((data & 0x02) == 0) {
      digitalWrite(Dout_1, 0);
    }
    if ((data & 0x04) == 0) {
      digitalWrite(Dout_2, 0);
    }
    if ((data & 0x08) == 0) {
      digitalWrite(Dout_3, 0);
    }
    delay(5);
  } while (value < millis());
}
