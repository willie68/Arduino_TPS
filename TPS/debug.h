/*
  debug.h - Definition von 2 Debugfunktionen - Version 0.1
  Copyright (c) 2012 Wilfried Klaas.  All right reserved.

  This library is free software; you can redistribute it and/or
  modify it under the terms of the GNU Lesser General Public
  License as published by the Free Software Foundation; either
  version 2.1 of the License, or (at your option) any later version.

  This library is distributed in the hope that it will be useful,
  but WITHOUT ANY WARRANTY; without even the implied warranty of
  MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the GNU
  Lesser General Public License for more details.

  You should have received a copy of the GNU Lesser General Public
  License along with this library; if not, write to the Free Software
  Foundation, Inc., 51 Franklin St, Fifth Floor, Boston, MA  02110-1301  USA
*/

// Program im Debugmodus kompilieren, dann werden zus. Ausgaben auf die serielle Schnittstelle geschrieben.
// Zum aktivieren der Debug Funktion bitte den Define VOR dem #include "debug.h" in die Hauptdatei eintragen.
//#define debug
#include <Arduino.h>
#include <inttypes.h>
#ifndef debug_lib
#define debug_lib

#if defined(debug) && !defined(__AVR_ATtiny84__)
#define dbgOut(S) \
Serial.print(S); 
#define dbgOut2(S,P) \
Serial.print(S,P); 
#define dbgOutLn(S) \
Serial.println(S); 
#define dbgOutLn2(S,P) \
Serial.println(S,P); 
#define initDebug() \
  Serial.begin(9600); \
  Serial.flush(); \
  delay(100);
#else
#define dbgOut(S)
#define dbgOut2(S,P)
#define dbgOutLn(S)
#define dbgOutLn2(S,P)
#define initDebug()

#endif

/* 
void dbgOutHex(uint8_t data) {
  dbgOut(" 0x");
  if (data < 0x10) {
    dbgOut('0');
  }
  dbgOut2(data, HEX);
}
*/
#endif
