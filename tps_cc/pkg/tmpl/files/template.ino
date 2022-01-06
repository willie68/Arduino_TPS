// Program im Debugmodus kompilieren, dann werden zus. Ausgaben auf die serielle Schnittstelle geschrieben.
{{.debug}}

/*
   Here are the defines used in this software to control special parts of the implementation
   #define TPS_RCRECEIVER: using a RC receiver input
   #define TPS_ENHANCEMENT: all of the other enhancments
   #define TPS_SERVO: using servo outputs
   #define TPS_TONE: using a tone output

*/
{{.flags}}

// libraries
#include "debug.h"
#include "makros.h"

#ifdef ESP32
#include <ESP32Servo.h>
#ifdef TPS_TONE
#include <ESP32Tone.h>
#endif
#endif

#ifdef TPS_SERVO
#if defined(__AVR_ATmega328P__) || defined(__AVR_ATtiny84__) || defined(__AVR_ATtiny861__) || defined(__AVR_ATtiny4313__) || defined(__AVR_ATtiny2313A__)
#include <Servo.h>
#endif
#endif

#ifdef TPS_TONE
#include "notes.h"
#endif

#include "hardware.h"
// sub routines
const byte subCnt = 7;
word subs[subCnt];

// page register
word page;
// defining register
byte a, b, c, d, e, f;

#ifdef TPS_ENHANCEMENT
const byte SAVE_CNT = 16;
#else
const byte SAVE_CNT = 1;
#endif

word saveaddr[SAVE_CNT];
int saveCnt;

#ifdef TPS_ENHANCEMENT
byte stack[SAVE_CNT];
byte stackCnt;
#endif

unsigned long tmpValue;

#ifdef TPS_SERVO
Servo servo1;
Servo servo2;
#endif

void setup() {
  initDebug();
  // put your setup code here, to run once:
  pinMode(Dout_1, OUTPUT);
  pinMode(Dout_2, OUTPUT);
  pinMode(Dout_3, OUTPUT);
  pinMode(Dout_4, OUTPUT);

  pinMode(PWM_1, OUTPUT);
  pinMode(PWM_2, OUTPUT);

  pinMode(Din_1, INPUT_PULLUP);
  pinMode(Din_2, INPUT_PULLUP);
  pinMode(Din_3, INPUT_PULLUP);
  pinMode(Din_4, INPUT_PULLUP);

  pinMode(SW_PRG, INPUT_PULLUP);
  pinMode(SW_SEL, INPUT_PULLUP);

  initHardware();

  digitalWrite(Dout_1, 1);
  delay(1000);
  digitalWrite(Dout_1, 0);

  doReset();

#ifdef TPS_SERVO
  servo1.attach(SERVO_1);
  servo2.attach(SERVO_2);
#endif

#ifdef TPS_ENHANCEMENT
  pinMode(LED_BUILTIN, OUTPUT);
#endif

// todo setup servos
{{.setup}}
}

void loop() {
  // put your main code here, to run repeatedly:
{{.main}}
}

void doReset() {

  for (int i = 0; i < subCnt; i++) {
    subs[i] = 0;
  }

  page = 0;
  saveCnt = 0;
  a = 0;
  b = 0;
  c = 0;
  d = 0;
  e = 0;
  f = 0;
#ifdef TPS_ENHANCEMENT
  stackCnt = 0;
  for (int i = 0; i < SAVE_CNT; i++) {
    stack[i] = 0;
  }
#endif
}

/*
  output to port
*/
void doPort(byte data) {
  digitalWrite(Dout_1, (data & 0x01) > 0);
  digitalWrite(Dout_2, (data & 0x02) > 0);
  digitalWrite(Dout_3, (data & 0x04) > 0);
  digitalWrite(Dout_4, (data & 0x08) > 0);
}

/*
  delay in ms
*/
void doDelay(byte data) {
  switch (data) {
    case 0:
      delay(1);
      break;
    case 1:
      delay(2);
      break;
    case 2:
      delay(5);
      break;
    case 3:
      delay(10);
      break;
    case 4:
      delay(20);
      break;
    case 5:
      delay(50);
      break;
    case 6:
      delay(100);
      break;
    case 7:
      delay(200);
      break;
    case 8:
      delay(500);
      break;
    case 9:
      delay(1000);
      break;
    case 10:
      delay(2000);
      break;
    case 11:
      delay(5000);
      break;
    case 12:
      delay(10000);
      break;
    case 13:
      delay(20000);
      break;
    case 14:
      delay(30000);
      break;
    case 15:
      delay(60000);
      break;
    default:
      break;
  }
}