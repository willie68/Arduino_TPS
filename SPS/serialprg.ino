#ifdef __AVR_ATtiny84__
#define BAUDRATE 9600
#endif

#ifdef __AVR_ATmega328P__
#define BAUDRATE 9600
#endif

void serialPrg() {
  //int value1, value2, value3, value4;
  byte value;
  bool endOfPrg = false;
  bool endOfFile = false;
  byte data[32];
  char tmp[16];
  word readAddress;
  int crc, readcrc;
  byte count;
  byte type;

  addr = 0;
  Serial.end();
  Serial.begin(BAUDRATE);
  Serial.println("serPrgStart");
  Serial.println("waiting for command:");
  Serial.println("w: write HEX file, r: read EPPROM, e: end");
  while (!endOfPrg) {
    while (Serial.available() > 0) {
      // look for the next valid integer in the incoming serial stream:
      char myChar = Serial.read();
      if (myChar == 'w') {
        // hexfile is comming to programm
        Serial.println("waitin");
        addr = 0;
        do {
          for (byte i = 0; i < 8; i++) {
            data[i] = 0xFF;
          }
          do {
            c = getNextChar();
          } while (!(c == ':'));

          Serial.print(".");

          // read counter
          c = getNextChar();
          count = hexToByte(c) << 4;
          c = getNextChar();
          count += hexToByte(c);
          printHex8(count);

          crc = count;
          Serial.print(".");

          // address
          c = getNextChar();
          readAddress = hexToByte(c) << 12;
          c = getNextChar();
          readAddress += hexToByte(c) << 8;
          c = getNextChar();
          readAddress += hexToByte(c) << 4;
          c = getNextChar();
          readAddress += hexToByte(c);

          printHex16(readAddress);

          crc += readAddress >> 8;
          crc += readAddress & 0x00FF;
          Serial.print(".");

          // reading data type
          c = getNextChar();
          type = hexToByte(c) << 4;
          c = getNextChar();
          type += hexToByte(c);
          printHex8(type);

          crc += type;

          Serial.print(".");

          if (type == 0x01) {
            endOfFile = true;
          }

          // read data bytes
          for (byte x = 0; x < count; x++) {
            c = getNextChar();
            value = hexToByte(c) << 4;
            c = getNextChar();
            value += hexToByte(c);
            printHex8(value);

            Serial.print(".");
            data[x] = value;
            crc += value;
          }
          // read CRC
          c = getNextChar();
          readcrc = hexToByte(c) << 4;
          c = getNextChar();
          readcrc += hexToByte(c);
          printHex8(readcrc);

          crc += readcrc;
          // check CRC
          value = crc & 0x00FF;
          printHex8(value);

          if (value == 0) {
            Serial.print("ok");
            // adding value to EEPROM
            for (byte x = 0; x < count; x++) {
              EEPROM.write(readAddress + x, data[x]);
            }
          } else {
            Serial.println("CRC Error");
            endOfFile = true;
          }

          Serial.println();
        } while (!(endOfFile));
        Serial.println("endOfFile");
      }
      if (myChar == 'r') {
        // write eeprom as hexfile to receiver
        Serial.println("EEPROM data:");
        byte checksum = 0;
        for (int addr = 0; addr <= E2END; addr++) {
          value = EEPROM.read(addr);
          if ((addr % 16) == 0) {
            printCheckSum(checksum);
            checksum = 0;
            Serial.print(":10");
            checksum += 0x10;
            printHex16(addr);
            checksum += (addr >> 8);
            checksum += (addr & 0x00FF);
            Serial.print("00");
          }
          printHex8(value);
          checksum += value;
        }
        printCheckSum(checksum);
        // ending
        Serial.println(":00000001FF");
      }
      if (myChar == 'e') {
        // end of programm
        endOfPrg = true;
      }
    }
  }
  Serial.println("end");
  Serial.end();
  doReset();
}

char getNextChar() {
  while (!Serial.available()) {
  }
  return  Serial.read();
}
void printCheckSum(byte checksum) {
  printHex8(checksum);
  Serial.println();
}

void printHex8(int num) {
  char tmp[16];
  sprintf(tmp, "%.2X", num);
  Serial.print(tmp);
}

void printHex16(int num) {
  char tmp[16];
  sprintf(tmp, "%.4X", num);
  Serial.print(tmp);
}

byte hexToByte (char c) {
  if ( (c >= '0') && (c <= '9') ) {
    return c - '0';
  }
  if ( (c >= 'A') && (c <= 'F') ) {
    return (c - 'A') + 10;
  }
}
