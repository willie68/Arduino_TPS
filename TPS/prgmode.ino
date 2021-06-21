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
/* this is the program 
Addr   BD   Befehl   Daten     Kommentar
0x00:  4F   0X00     XXXX      A=#:A=15   ,""
0x01:  59   0X0X     X00X      =A:PWM.1=A   ,""
0x02:  1F   000X     XXXX      Dout:Output 1111   ,""
0x03:  29   00X0     X00X      Delay:Delay 1s   ,""
0x04:  10   000X     0000      Dout:Output 0000   ,""
0x05:  29   00X0     X00X      Delay:Delay 1s   ,""
0x06:  5A   0X0X     X0X0      =A:PWM.2=A   ,""
0x07:  40   0X00     0000      A=#:A=0   ,""
0x08:  59   0X0X     X00X      =A:PWM.1=A   ,""
0x09:  64   0XX0     0X00      A=:A=Din   ,""
0x0A:  54   0X0X     0X00      =A:Dout=A   ,""
0x0B:  29   00X0     X00X      Delay:Delay 1s   ,""
0x0C:  4F   0X00     XXXX      A=#:A=15   ,""
0x0D:  59   0X0X     X00X      =A:PWM.1=A   ,""
0x0E:  10   000X     0000      Dout:Output 0000   ,""
0x0F:  CD   XX00     XX0X      Skip if:S_PRG=0   ,""
0x10:  11   000X     000X      Dout:Output 0001   ,""
0x11:  28   00X0     X000      Delay:Delay 500ms   ,""
0x12:  CC   XX00     XX00      Skip if:S_SEL=0   ,""
0x13:  18   000X     X000      Dout:Output 1000   ,""
0x14:  28   00X0     X000      Delay:Delay 500ms   ,""
0x15:  4F   0X00     XXXX      A=#:A=15   ,""
0x16:  59   0X0X     X00X      =A:PWM.1=A   ,""
0x17:  5A   0X0X     X0X0      =A:PWM.2=A   ,""
0x18:  72   0XXX     00X0      A=Calculation:A=A-1   ,""
0x19:  26   00X0     0XX0      Delay:Delay 100ms   ,""
0x1A:  C0   XX00     0000      Skip if:A=0   ,""
0x1B:  35   00XX     0X0X      Jump -:jump -5   ,""
0x1C:  80   X000     0000      Page:Page 0   ,""
0x1D:  90   X00X     0000      Jump:Jump 0   ,""
0x1E:  FF   XXXX     XXXX      Byte/Board:PrgEnd   ,""
*/

enum PROGRAMMING_MODE {ADDRESS, COMMAND, DATA};

PROGRAMMING_MODE prgMode;

void prgDemoPrg() {
  byte value = readbyte(0);
  if (value == 0xFF) {
    dbgOutLn("recreate demo program");
    value = readbyte(1);
    if (value == 0xFF) {
      for (byte i = 0; i < sizeof(demoPrg); i++) {
        writebyte(i, demoPrg[i]);
      }
    }
  }
}

void programMode() {
  // checking if advance programmer board connected?
#ifdef TPS_ENHANCEMENT
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
      doAddr(addr);
      //delay(SHOW_DELAY);

      blinkD2();
      // HiNibble Adresse anzeigen
      data = (addr & 0xf0) >> 4;                                  //Adresse anzeigen
      doAddr(data);
      //delay(SHOW_DELAY);

      byte Eebyte = readbyte(addr);
      data = Eebyte & 15;
      cmd = Eebyte >> 4;

      blinkD3();
      prgMode = COMMAND;
      dbgOutLn("cmd");
      doPort(cmd); //show command

      do {
        if (digitalRead(SW_SEL) == 0) {
          delay(KEY_DELAY);
          cmd += 1;
          cmd = cmd & 0x0F;
          doPort(cmd);
        }
      }
      while (digitalRead(SW_PRG) == 1);
      delay(DEBOUNCE);
      if (digitalRead(SW_SEL) == 0) {
        break;
      }
      
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
      if (digitalRead(SW_SEL) == 0) {
        break;
      }

      byte newValue = (cmd << 4) + data;
      if (newValue != Eebyte) {
        writebyte(addr, newValue); //           Writeeeprom Eebyte , Addr
        blinkAll();
      }
      addr += 1;
    }
    while (true);
#ifdef TPS_ENHANCEMENT
  }
#endif
  dbgOutLn("save program data");
  store();
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
