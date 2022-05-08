// Program im Debugmodus kompilieren, dann werden zus. Ausgaben auf die serielle Schnittstelle geschrieben.
#define debug

/*
   Here are the defines used in this software to control special parts of the implementation
   #define TPS_RCRECEIVER: using a RC receiver input
   #define TPS_ENHANCEMENT: all of the other enhancments
   #define TPS_SERVO: using servo outputs
   #define TPS_TONE: using a tone output

*/
#define TPS_RCRECEIVER 
#define TPS_ENHANCEMENT 
#define TPS_SERVO 


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

}

void loop() {
  // put your main code here, to run repeatedly:
  static void *array[] = { &&label_0, &&label_1, &&label_2, &&label_3, &&label_4, &&label_5, &&label_6, &&label_7, &&label_8, &&label_9, &&label_10, &&label_11, &&label_12, &&label_13, &&label_14, &&label_15, &&label_16, &&label_17, &&label_18, &&label_19, &&label_20, &&label_21, &&label_22, &&label_23, &&label_24, &&label_25, &&label_26};
label_0: // "Nullpunktwert ins B Register"
  a=8;
label_1: // ""
  b=a;
label_2: // ""
  a=0;
label_3: // ""
  a=b-a;
label_4: // ""
  b=a;
label_5: // "Reglerkanal leseen"
  #ifdef TPS_RCRECEIVER
  tmpValue = pulseIn(RC_0, HIGH, 100000);
  if (tmpValue < 1000) {
    tmpValue = 1000;
  }
  if (tmpValue > 2000) {
    tmpValue = 2000;
  }
  a = (tmpValue - 1000) / 4; //(Umrechnen auf 8 bit)
#endif;
label_6: // "sichern"
  e=a;
label_7: // "Rückwärts"
  if (a<b) { goto *array[9];}
label_8: // "Vorwärts"
  goto *array[(page * 16) + 11];
label_9: // "Klappenservo ausfahren"
  a=15;
label_10: // "Ist immer erfüllt"
  if (a<b) { goto *array[12];}
label_11: // "Vorwärts kein Klappenservo"
  a=0;
label_12: // ""
  #ifdef TPS_SERVO
  tmpValue = ((a & 0x0f) * 10) + 10;
  servo1.write(tmpValue);
  #endif;
label_13: // "RC Wert zurückholen"
  a=e;
label_14: // ""
  page=1;
label_15: // ""
  if (a<b) { goto *array[17];}
label_16: // "Bei Vorwärts ohne Berechnung direkt springen"
  goto *array[(page * 16) + 7];
label_17: // ""
  a=15;
label_18: // ""
  b=a;
label_19: // ""
  a=b-a;
label_20: // ""
  b=a;
label_21: // ""
  a=e;
label_22: // ""
  a=a^b;
label_23: // "Fahrt an den Fahrservo übergeben"
  analogWrite(PWM_2, a);;
label_24: // ""
  page=0;
label_25: // "wieder nach vorne springen"
  goto *array[(page * 16) + 0];
label_26: // ""
  #ifdef __AVR_ATmega328P__
  digitalWrite(LED_BUILTIN, 0);
#endif;
  ;
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

void push() {
  if (stackCnt < SAVE_CNT) {
    stack[stackCnt] = a;
    stackCnt += 1;
  } else {
    for (int i = 1; i < SAVE_CNT; i++) {
      stack[i - 1] = stack[i];
    }
    stack[SAVE_CNT-1] = a;
  }
}

void pop() {
  dbgOut("pop ");
  dbgOutLn(stackCnt);
  if (stackCnt > 0) {
    stackCnt -= 1;
    a = stack[stackCnt];
  } else {
    a = 0;
  }
}