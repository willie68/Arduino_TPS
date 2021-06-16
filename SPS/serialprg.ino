#ifdef SPS_SERIAL_PRG
#if defined(__AVR_ATtiny861__) || defined(__AVR_ATtiny84__) || defined(__AVR_ATmega328P__) || defined(ESP32)
#define BAUDRATE 9600
#endif

void initSerialPrg() {
  Serial.end();
  Serial.begin(BAUDRATE);  
  Serial.println();
}

void sendHeader() {
#ifdef __AVR_ATtiny84__
  Serial.println("TinySPS");
#endif
#ifdef __AVR_ATmega328P__
  Serial.println("ArduinoSPS");
#endif
  Serial.print("max prg size:");
  Serial.print(STORESIZE, HEX);
  Serial.println();
  Serial.println("waiting for command:");
  Serial.println("w: write HEX file, r: read EPPROM, e: end");  
}

void serialPrg() {
  byte value;
  bool endOfPrg = false;
  bool endOfFile = false;
  byte data[32];
  word readAddress;
  int crc, readcrc;
  byte count;
  byte type;

  addr = 0;
  sendHeader();
  while (!endOfPrg) {
    while (Serial.available() > 0) {
      // look for the next valid integer in the incoming serial stream:
      char myChar = Serial.read();
      if (myChar == 'w') {
        // hexfile is comming to programm
        endOfFile = false;
        Serial.println("ready");
        addr = 0;
        do {
          for (byte i = 0; i < 8; i++) {
            data[i] = 0xFF;
          }
          do {
            c = getNextChar();
          } while (!(c == ':'));

#ifdef debug
          Serial.print(".");
#endif
          // read counter
          c = getNextChar();
          count = hexToByte(c) << 4;
          c = getNextChar();
          count += hexToByte(c);
#ifdef debug
          printHex8(count);
#endif

          crc = count;
#ifdef debug
          Serial.print(".");
#endif
          // address
          c = getNextChar();
          readAddress = hexToByte(c) << 12;
          c = getNextChar();
          readAddress += hexToByte(c) << 8;
          c = getNextChar();
          readAddress += hexToByte(c) << 4;
          c = getNextChar();
          readAddress += hexToByte(c);

#ifdef debug
          printHex16(readAddress);
#endif

          crc += readAddress >> 8;
          crc += readAddress & 0x00FF;
#ifdef debug
          Serial.print(".");
#endif

          // reading data type
          c = getNextChar();
          type = hexToByte(c) << 4;
          c = getNextChar();
          type += hexToByte(c);
#ifdef debug
          printHex8(type);
#endif

          crc += type;

#ifdef debug
          Serial.print(".");
#endif

          if (type == 0x01) {
            endOfFile = true;
          }

          // read data bytes
          for (byte x = 0; x < count; x++) {
            c = getNextChar();
            value = hexToByte(c) << 4;
            c = getNextChar();
            value += hexToByte(c);

#ifdef debug
            printHex8(value);
            Serial.print(".");
#endif

            data[x] = value;
            crc += value;
          }
          // read CRC
          c = getNextChar();
          readcrc = hexToByte(c) << 4;
          c = getNextChar();
          readcrc += hexToByte(c);

#ifdef debug
          printHex8(readcrc);
          Serial.print(".");
#endif
          

          crc += readcrc;
          // check CRC
          value = crc & 0x00FF;
#ifdef debug
          printHex8(value);
#endif

          if (value == 0) {
            Serial.print("ok");
            // adding value to EEPROM
            for (byte x = 0; x < count; x++) {
              writebyte(readAddress + x, data[x]);
            }
          } else {
            Serial.println(", CRC Error");
            endOfFile = true;
          }

          Serial.println();
        } while (!(endOfFile));
        Serial.println("endOfFile");
        store();
      }
      if (myChar == 'r') {
        // write eeprom as hexfile to receiver
        Serial.println("EEPROM data:");
        int checksum = 0;
        for (int addr = 0; addr <= STORESIZE; addr++) {
          value = readbyte(addr);
          if ((addr % 8) == 0) {
            printCheckSum(checksum);
            checksum = 0;
            Serial.print(":08");
            checksum += 0x08;
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
      if (myChar == 'h') {
        sendHeader();
      }
    }
  }
  Serial.println("end");
}

char getNextChar() {
  while (!Serial.available()) {
  }
  return  Serial.read();
}

void printCheckSum(int value) {
  int checksum = value & 0xFF;
  checksum = (checksum ^ 0xFF) + 1;
  printHex8(checksum);
  Serial.println();
}

void printHex8(int num) {
  char tmp[3];
  tmp[0] = nibbleToHex(num >> 4);
  tmp[1] = nibbleToHex(num);
  tmp[2] = 0x00;
  Serial.print(tmp);
}

void printHex16(int num) {
  char tmp[5];
  tmp[0] = nibbleToHex(num >> 12);
  tmp[1] = nibbleToHex(num >> 8);
  tmp[2] = nibbleToHex(num >> 4);
  tmp[3] = nibbleToHex(num);
  tmp[4] = 0x00;
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

byte nibbleToHex (byte value) {
  byte c = value & 0x0F;
  if ( (c >= 0) && (c <= 9) ) {
    return c + '0';
  }
  if ( (c >= 10) && (c <= 15) ) {
    return (c + 'A') - 10;
  }
}
#endif
