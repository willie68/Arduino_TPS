// defining the hardware connections

#ifdef __AVR_ATmega328P__
// Arduino Hardware
const byte Din_1 = 0;
const byte Din_2 = 1;
const byte Din_3 = 2;
const byte Din_4 = 3;

const byte Dout_1 = 4;
const byte Dout_2 = 5;
const byte Dout_3 = 6;
const byte Dout_4 = 7;

const byte ADC_0 = 0; //(15)
const byte ADC_1 = 1; //(16)
const byte PWM_1 = 9;
const byte PWM_2 = 10;

#ifdef TPS_TONE
const byte TONE_OUT = PWM_2;
#endif

#ifdef TPS_RCRECEIVER
const byte RC_0 = 18;
const byte RC_1 = 19;
#endif

#ifdef TPS_SERVO
const byte SERVO_1 = 9;
const byte SERVO_2 = 10;
#endif

const byte SW_PRG = 8;
const byte SW_SEL = 11;

#ifdef TPS_USE_DISPLAY
const byte DIGIT_DATA_IO = 12;
const byte DIGIT_CLOCK = 13;
#endif

#define getAnalog(pin) (analogRead(pin) >> 2)
#define initHardware()

#endif

#ifdef __AVR_ATtiny84__
// ATTiny84 Hardware
const byte Dout_1 = 6;
const byte Dout_2 = 5;
const byte Dout_3 = 4;
const byte Dout_4 = 1;

const byte Din_1 = 10;
const byte Din_2 = 9;
const byte Din_3 = 8;
const byte Din_4 = 7;
const byte ADC_0 = 0;
const byte ADC_1 = 1;
const byte PWM_1 = 2;
const byte PWM_2 = 3;

#ifdef TPS_TONE
const byte TONE_OUT = PWM_2;
#endif

#ifdef TPS_RCRECEIVER
const byte RC_0 = 10;
const byte RC_1 = 9;
#endif

#ifdef TPS_SERVO
const byte SERVO_1 = 2;
const byte SERVO_2 = 3;
#endif

const byte SW_PRG = 0;
const byte SW_SEL = 8;

#ifdef TPS_USE_DISPLAY
const byte DIGIT_DATA_IO = 4;
const byte DIGIT_CLOCK = 5;
#endif

#define getAnalog(pin) (analogRead(pin) >> 2)
#define initHardware()

#endif

#ifdef __AVR_ATtiny4313__
// ATTiny4313 Hardware
const byte Dout_1 = 0;
const byte Dout_2 = 1;
const byte Dout_3 = 2;
const byte Dout_4 = 3;

const byte Din_1 = 4;
const byte Din_2 = 5;
const byte Din_3 = 6;
const byte Din_4 = 7;
const byte ADC_0 = 13;
const byte ADC_1 = 14;
const byte PWM_1 = 11;
const byte PWM_2 = 12;

#ifdef TPS_TONE
const byte TONE_OUT = PWM_2;
#endif

#ifdef TPS_RCRECEIVER
const byte RC_0 = 15;
const byte RC_1 = 16;
#endif

#ifdef TPS_SERVO
const byte SERVO_1 = 11;
const byte SERVO_2 = 12;
#endif

const byte SW_PRG = 9;
const byte SW_SEL = 8;

#define getAnalog(pin) (analogRead(pin) >> 2)
#define initHardware()

#endif

#ifdef __AVR_ATtiny861__
// ATTiny4313 Hardware
const byte Dout_1 = 0;
const byte Dout_2 = 1;
const byte Dout_3 = 2;
const byte Dout_4 = 3;

const byte Din_1 = 4;
const byte Din_2 = 5;
const byte Din_3 = 6;
const byte Din_4 = 7;
const byte ADC_0 = 13;
const byte ADC_1 = 14;
const byte PWM_1 = 11;
const byte PWM_2 = 12;

#ifdef TPS_TONE
const byte TONE_OUT = PWM_2;
#endif

#ifdef TPS_RCRECEIVER
const byte RC_0 = 15;
const byte RC_1 = 16;
#endif

#ifdef TPS_SERVO
const byte SERVO_1 = 11;
const byte SERVO_2 = 12;
#endif

const byte SW_PRG = 9;
const byte SW_SEL = 8;

#define getAnalog(pin) (analogRead(pin) >> 2)
#define initHardware()

#endif

#ifdef _MICROBIT_V2_
// Microbit V2 Hardware
const byte Din_1 = 0;
const byte Din_2 = 1;
const byte Din_3 = 2;
const byte Din_4 = 3;

const byte Dout_1 = 4;
const byte Dout_2 = 5;
const byte Dout_3 = 6;
const byte Dout_4 = 7;

const byte ADC_0 = 0; //(15)
const byte ADC_1 = 1; //(16)
const byte PWM_1 = 9;
const byte PWM_2 = 10;

#ifdef TPS_TONE
const byte TONE_OUT = PWM_2;
#endif

#ifdef TPS_RCRECEIVER
const byte RC_0 = 18;
const byte RC_1 = 19;
#endif

#ifdef TPS_SERVO
const byte SERVO_1 = 9;
const byte SERVO_2 = 10;
#endif

const byte SW_PRG = 8;
const byte SW_SEL = 11;

#ifdef TPS_USE_DISPLAY
const byte DIGIT_DATA_IO = 12;
const byte DIGIT_CLOCK = 13;
#endif

#define getAnalog(pin) (analogRead(pin) >> 2)
#define initHardware()

#endif


#ifdef ESP32
#ifdef ARDUINO_ESP32_DEV
  #define LED_BUILTIN 2
#endif
#include "esp32.h"

// ESP32 Hardware
const byte Din_1 = 26;
const byte Din_2 = 18;
const byte Din_3 = 19;
const byte Din_4 = 23;

const byte Dout_1 = 22;
const byte Dout_2 = 21;
const byte Dout_3 = 17;
const byte Dout_4 = 16;

const byte ADC_0 = 36; 
const byte ADC_1 = 39; 
const byte PWM_1 = 27;
const byte PWM_2 = 25;

#ifdef TPS_TONE
const byte TONE_OUT = PWM_2;
#endif

#ifdef TPS_RCRECEIVER
const byte RC_0 = 34;
const byte RC_1 = 35;
#endif

#ifdef TPS_SERVO
const byte SERVO_1 = 27; //14
const byte SERVO_2 = 25; //32
#endif

const byte SW_PRG = 13;
const byte SW_SEL = 12;

#ifdef TPS_USE_DISPLAY
const byte DIGIT_DATA_IO = 32;
const byte DIGIT_CLOCK = 33;
#endif

#define getAnalog(pin) (analogRead(pin) >> 4)
#define initHardware() initESP32()
#endif
