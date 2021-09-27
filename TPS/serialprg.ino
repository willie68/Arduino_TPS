#ifdef TPS_SERIAL_PRG
#if defined(__AVR_ATtiny861__) || defined(__AVR_ATtiny84__) || defined(__AVR_ATmega328P__) || defined(ESP32)
#define BAUDRATE 9600
#else
#error "not implemetned yet: Maybe no serial port for programming available"
#endif

#include "version.h"

void initSerialPrg() {
  Serial.end();
  Serial.begin(BAUDRATE);
  Serial.println();
}

void sendHeader() {
#ifdef __AVR_ATtiny84__
  Serial.print("Tiny_TPS ");
#endif
#ifdef __AVR_ATmega328P__
  Serial.print("Arduino_TPS ");
#endif
  Serial.println(TPS_VERSION);
  Serial.print("max prg size:");
  Serial.print(STORESIZE, HEX);
  Serial.println();
  Serial.println("waiting for command:");
  Serial.println("w: write HEX file, r: read EPPROM, e: end");
  Serial.println("i: get inputs, o: write to output, a#: get analog in from #, p#: set pwm #, b: get button state");
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
      if ((myChar == 'W') || (myChar == 'w')) {
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
      if ((myChar == 'R') || (myChar == 'r')) {
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
      if ((myChar == 'I') || (myChar == 'i')) {
        value = digitalRead(Din_1) + (digitalRead(Din_2) << 1) + (digitalRead(Din_3) << 2) + (digitalRead(Din_4) << 3);
        Serial.print("i:0x");
        printHex4(value);
        Serial.println();
      }
      if ((myChar == 'O') || (myChar == 'o')) {
        c = getNextChar();
        value = hexToByte(c);
        doPort(value);
        Serial.print("o:0x");
        printHex4(value);
        Serial.println();
      }
      if ((myChar == 'A') || (myChar == 'a')) {
        c = getNextChar();
        if (c == '1') {
          Serial.print("a1:0x");
          value = getAnalog(ADC_0);
        } else {
          Serial.print("a2:0x");
          value = getAnalog(ADC_1);
        }
        printHex8(value);
        Serial.println();
      }
      if ((myChar == 'P') || (myChar == 'p')) {
        c = getNextChar();
        d = getNextChar();
        a = getNextChar();
        b = getNextChar();
        value = (hexToByte(a) << 4) + hexToByte(b);
        if (c == '1') {
          Serial.print("p1:0x");
          analogWrite(PWM_1, value);
        } else {
          Serial.print("p2:0x");
          analogWrite(PWM_2, value);
        }
        printHex8(value);
        Serial.println();
      }
      if ((myChar == 'B') || (myChar == 'b')) {
        value = digitalRead(SW_PRG) + (digitalRead(SW_SEL) << 1);
        Serial.print("b:0x");
        printHex4(value);
        Serial.println();
      }
      if ((myChar == 'E') ||(myChar == 'e')) {
        // end of programm
        endOfPrg = true;
      }
      if ((myChar == 'H') || (myChar == 'h')) {
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

void printHex4(int num) {
  char tmp[2];
  tmp[0] = nibbleToHex(num);
  tmp[1] = 0x00;
  Serial.print(tmp);
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
  if ( (c >= 'a') && (c <= 'f') ) {
    return (c - 'a') + 10;
  }
  return 0;
}

byte nibbleToHex (byte value) {
  byte c = value & 0x0F;
  if (c <= 9) {
    return c + '0';
  }
  return (c + 'A') - 10;
}
#endif
