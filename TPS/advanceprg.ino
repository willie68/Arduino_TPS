#ifdef TPS_USE_DISPLAY
#include <TM1637Display.h>
#endif

#ifndef TPS_USE_DISPLAY
void advancePrg() {}
#endif

#ifdef TPS_USE_DISPLAY
TM1637Display display(DIGIT_CLOCK, DIGIT_DATA_IO);
byte displayValues[4];

const byte PRGS_ADDR = 0;
const byte PRGS_COM = 1;
const byte PRGS_DATA = 2;

Switch prgSwitch = Switch(SW_PRG);
Switch selSwitch = Switch(SW_SEL);

void advancePrg() {
  resetDisplay();
  showPrg();
  dbgOutLn("adv prg mode");
  while ((digitalRead(SW_PRG) == 0) || (digitalRead(SW_SEL) == 0));
  delay(DEBOUNCE);

  addr = 0;
  byte Eebyte = EEPROM.read(addr);
  boolean dirty = false;

  showAddress(addr, Eebyte);

  unsigned long saveTime = millis();

  byte prgState = PRGS_COM;
  showActualState(prgState);
  prgSwitch.poll();
  selSwitch.poll();

  // hier die Programmschleife
  do {
    prgSwitch.poll();
    selSwitch.poll();

    boolean lit = (((millis() - saveTime) % 500) > 250);
    display.setColon(lit);
    outputDisplay();
    if (prgSwitch.pushed()) {
      prgState += 1;
      prgState = prgState % 3;
      if (prgState == PRGS_ADDR) {
        if (dirty) {
          dbgOut("prg:");
          dbgOut2(addr, HEX);
          dbgOut(':');
          dbgOutLn2(Eebyte, HEX);

          displayValues[2] = 0x40;
          displayValues[3] = 0x40;
          outputDisplay();

          EEPROM.write(addr, Eebyte);
          delay(500);

          displayValues[2] = 0x00;
          displayValues[3] = 0x00;
          outputDisplay();
        }
        addr = (addr + 1) % E2END;
        Eebyte = EEPROM.read(addr);
        dirty = false;
        prgState = PRGS_COM;
        showAddress(addr, Eebyte);
      }
      showActualState(prgState);
    }

    if (selSwitch.pushed()) {
      switch (prgState) {
        case PRGS_COM:
          dbgOut('c');
          dirty = true;
          Eebyte = Eebyte + 0x10;
          showAddress(addr, Eebyte);
          break;
        case PRGS_DATA:
          dbgOut('d');
          dirty = true;
          byte data = Eebyte & 0x0F;
          data += 1;
          data = data & 0x0F;
          Eebyte = (Eebyte & 0xF0) + data;
          showAddress(addr, Eebyte);
          break;
      }
    }
  }
  while (true);
}

void showActualState(byte prgState) {
  switch (prgState) {
    case PRGS_ADDR:
      blinkAddr();
      delay(250);
      blinkAddr();
      break;
    case PRGS_COM:
      blinkCol(2);
      delay(250);
      blinkCol(2);
      break;
    case PRGS_DATA:
      blinkCol(3);
      delay(250);
      blinkCol(3);
      break;
  }
}

void blinkAddr() {
  byte saveValue1 = displayValues[0];
  byte saveValue2 = displayValues[1];
  displayValues[0] = 0;
  displayValues[1] = 0;
  outputDisplay();
  delay(250);
  displayValues[0] = saveValue1;
  displayValues[1] = saveValue2;
  outputDisplay();
}

void blinkCol(byte col) {
  byte saveValue1 = displayValues[col];
  displayValues[col] = 0;
  outputDisplay();
  delay(250);
  displayValues[col] = saveValue1;
  outputDisplay();
}

void showPrg() {
  displayValues[0] = 0b01110011;
  displayValues[1] = 0b01010000;
  displayValues[2] = 0b01101111;
  outputDisplay();
}

void showAddress(word addr, byte Eebyte) {
  displayValues[0] = display.encodeDigit(addr >> 4);
  displayValues[1] = display.encodeDigit(addr & 15);

  data = Eebyte & 15;
  com = Eebyte >> 4;

  displayValues[2] = display.encodeDigit(com);
  displayValues[3] = display.encodeDigit(data);
  outputDisplay();
}

// routines for the tm1637 display
void outputDisplay() {
  display.setSegments(displayValues);
}

void initDisplay() {
  display.setBrightness(0x0f);
  resetDisplay();
  for (byte i = 0; i < 4; i++) {
    displayValues[i] = 0x40;
    delay(150);
    outputDisplay();
  }
  delay(250);
  resetDisplay();
}

void resetDisplay() {
  for (byte i = 0; i < 4; i++) {
    displayValues[i] = 0;
  }
  display.setSegments(displayValues);
}
#endif
