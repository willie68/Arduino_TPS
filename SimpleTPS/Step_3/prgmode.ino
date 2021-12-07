/*
  entering the programming mode
*/

#define BLINK_DELAY 500
#define SHOW_DELAY 1000
#define KEY_DELAY 250
#define ADDR_LOOP 50
#define DEBOUNCE 50

// simple blink program
const byte demoPrg[] = { 
  0x11, // Dout=1
  0x29, // 1000ms
  0x18, // Dout=8
  0x29, // 1000ms
  0x34, // Addr = addr -4
  0xFF  // EoP End of program
};

enum PROGRAMMING_MODE {ADDRESS, COMMAND, DATA};

byte data = 0;
byte com = 0;

PROGRAMMING_MODE prgMode;

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
