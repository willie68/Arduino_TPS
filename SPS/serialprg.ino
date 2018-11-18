

void serialPrg() {
  int value1, value2, value3, value4;
  byte value;
  bool endOfPrg = false;

  addr = 0;
  Serial.end();
  Serial.begin(57600);
  Serial.println("serPrgStart");
  while (!endOfPrg) {
    while (Serial.available() > 0) {

      // look for the next valid integer in the incoming serial stream:
      char myChar = Serial.read();
      if (myChar == 'd') {
        for (byte i = 0; i < 8; i++) {
          char c;
          while (!Serial.available()) {
          }
          do {
            c = Serial.read();
          } while (!(((c >= '0') && (c <= '9')) || ((c >= 'A') && (c <= 'F'))));
          value = hexToByte(c);

          while (!Serial.available()) {
          }
          do {
            c = Serial.read();
          } while (!(((c >= '0') && (c <= '9')) || ((c >= 'A') && (c <= 'F'))));

          value = value * 16 + hexToByte(c);
          EEPROM.write(addr + i, value);
        }
        Serial.print("w:");
        for (byte i = 0; i < 8; i++) {
          value = EEPROM.read(addr + i);;
          Serial.print(value, HEX);
          if (i < 7) {
            Serial.print(",");
          }
        }
        Serial.println();
        addr = addr + 8;
      }
      if (myChar == '*') {
        endOfPrg = true;
      }
    }
  }
  Serial.println("end");
  Serial.end();
  doReset();
}

byte hexToByte (char c) {
  if ( (c >= '0') && (c <= '9') ) {
    return c - '0';
  }
  if ( (c >= 'A') && (c <= 'F') ) {
    return (c - 'A') + 10;
  }
}
