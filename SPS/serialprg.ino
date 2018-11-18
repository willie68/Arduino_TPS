

void serialPrg() {
  int value1, value2, value3, value4;
  byte value;
  bool endOfPrg = false;

  addr = 0;
  Serial.begin(9600);
  Serial.println("serPrgStart");
  while (!endOfPrg) {
    while (Serial.available() > 0) {

      // look for the next valid integer in the incoming serial stream:
      char myChar = Serial.read();
      if (myChar == 'd') {
        for (byte i = 0; i < 8; i++) {
          value = readHexValue;
          myChar = readHexValue;
          value = value * 16 + hexToByte(myChar);
          EEPROM.write(addr, value);
          addr++;
        }
        Serial.println("w");
      }
      if (myChar == '*') {
        endOfPrg = true;
      }
    }
    Serial.println("end");
    Serial.end();
    doReset();
  }
}

byte readHexValue() {
  char c;
  do {
    c = Serial.read();
  } while (!((c >= '0') && (c <= '9')) || ((c >= 'A') && (c <= 'F')));
  return hexToByte(c);
}

byte hexToByte (char c) {
  if ( (c >= '0') && (c <= '9') ) {
    return c - '0';
  }
  if ( (c >= 'A') && (c <= 'F') ) {
    return (c - 'A') + 10;
  }
}
