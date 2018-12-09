/*
  entering the programming mode
*/
void programMode() {
  // checking if advance programmer board connected?
#ifdef SPS_ENHANCEMENT
  if (digitalRead(SW_SEL) == 0) {
    advancePrg();
  }
  else {
#endif
    dbgOutLn("PrgMode");
    addr = 0;
    do {
      dbgOut("Adr:");
      dbgOutLn(addr);
      // LoNibble Adresse anzeigen
      doPort(addr);
      delay(300);
      lighting();
      // HiNibble Adresse anzeigen
      data = (addr & 0xf0) >> 4;                                  //Adresse anzeigen
      doPort(data);
      delay(300);
      lighting();

      byte Eebyte = EEPROM.read(addr);
      data = Eebyte & 15;
      com = Eebyte >> 4;
      dbgOutLn('commando eingeben');
      doPort(com); //Befehl anzeigen
      digitalWrite(PWM_1, HIGH);
      do {
      }
      while (digitalRead(SW_SEL) == 1); // S2 = 1
      delay(DEBOUNCE);

      prog = 1;                                            //Phase 1: Befehl anzeigen
      do {
        dbgOut("C:");
        dbgOut(com);
        dbgOut(", D:");
        dbgOut(data);
        dbgOut(", P:");
        dbgOutLn(prog);

        if (digitalRead(SW_PRG) == 0) {
          if (prog == 1) {
            prog = 2;
            com = 15;
          }
          if (prog == 2) {                                   //Phase 2: Befehl verändert
            com = com + 1;
            com = com & 15;
            doPort(com);
            digitalWrite(PWM_1, HIGH);
          }
          if (prog == 3) {                                   //Phase 3: Befehl unverändert, Daten ändern
            prog = 5;
            data = 15;
          }
          if (prog == 4) {                                   //Phase 4: Befehl und Daten geändert
            prog = 5;
            data = 15;
          }
          if (prog == 5) {                                   //Phase 5: Daten verändert
            data += 1;
            data = data & 15;
            doPort(data);
            digitalWrite(PWM_1, LOW);
          }
          delay(DEBOUNCE);
          do {
          }
          while (digitalRead(SW_PRG) == 1);
          delay(DEBOUNCE);
        }

        if (digitalRead(SW_SEL) == 0) {
          if (prog == 3) {
            prog = 7;          //nur angezeigt, nicht verändert
          }
          if (prog == 4) {
            doPort(255 - data);
            digitalWrite(PWM_1, LOW);
            prog = 6;
          }
          if (prog == 2) {
            doPort(data);             // Portd = Dat Or &HF0
            digitalWrite(PWM_1, LOW);
            prog = 4;
          }
          if (prog == 6) {                                    //nur Kommando wurde verändert
            data = data & 15;
            Eebyte = com * 16;
            Eebyte = Eebyte + data;
            EEPROM.write(addr, Eebyte); //           Writeeeprom Eebyte , Addr
            doPort(0x0F);
            delay(600);
            addr += 1;
            prog = 0;
          }
          if (prog == 5) {                                     //Daten wurden verändert
            data = data & 15;
            Eebyte = com * 16;
            Eebyte = Eebyte + data;
            EEPROM.write(addr, Eebyte); //           Writeeeprom Eebyte , Addr
            doPort(0xF0);
            delay(600);
            addr += 1;
            prog = 0;
          }
          if (prog == 7) {
            addr += 1;
            prog = 0;
          }
          delay(DEBOUNCE);
          do {
          }
          while (digitalRead(SW_SEL) == 0);
          delay(DEBOUNCE);
        }
      }
      while (prog != 0);
    }
    while (true);
#ifdef SPS_ENHANCEMENT
  }
#endif
}

void lighting() {
  doPort(0x0F);
  delay(200);
}
