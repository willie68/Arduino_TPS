// defining the hardware connections

#ifdef __AVR_ATmega328P__
// Arduino Hardware
const byte Din_0 = 0;
const byte Din_1 = 1;
const byte Din_2 = 2;
const byte Din_3 = 3;

const byte Dout_0 = 4;
const byte Dout_1 = 5;
const byte Dout_2 = 6;
const byte Dout_3 = 7;

const byte ADC_0 = 0; //(15)
const byte ADC_1 = 1; //(16)
const byte PWM_1 = 9;
const byte PWM_2 = 10;

#ifdef SPS_RCRECEIVER
const byte RC_0 = 18;
const byte RC_1 = 19;
#endif

#ifdef SPS_SERVO
const byte SERVO_1 = 9;
const byte SERVO_2 = 10;
#endif

const byte SW_PRG = 8;
const byte SW_SEL = 11;

#ifdef SPS_USE_DISPLAY
const byte DIGIT_DATA_IO = 12;
const byte DIGIT_CLOCK = 13;
#endif
#endif

#ifdef __AVR_ATtiny84__
// ATTiny84 Hardware
const byte Dout_0 = 6;
const byte Dout_1 = 5;
const byte Dout_2 = 4;
const byte Dout_3 = 1;

const byte Din_0 = 10;
const byte Din_1 = 9;
const byte Din_2 = 8;
const byte Din_3 = 7;
const byte ADC_0 = 0;
const byte ADC_1 = 1;
const byte PWM_1 = 2;
const byte PWM_2 = 3;
#ifdef SPS_RCRECEIVER
const byte RC_0 = 10;
const byte RC_1 = 9;
#endif

#ifdef SPS_SERVO
const byte SERVO_1 = 2;
const byte SERVO_2 = 3;
#endif

const byte SW_PRG = 0;
const byte SW_SEL = 8;

#ifdef SPS_USE_DISPLAY
const byte DIGIT_DATA_IO = 4;
const byte DIGIT_CLOCK = 5;
#endif
#endif

#ifdef __AVR_ATtiny4313__
// ATTiny4313 Hardware
const byte Dout_0 = 0;
const byte Dout_1 = 1;
const byte Dout_2 = 2;
const byte Dout_3 = 3;

const byte Din_0 = 4;
const byte Din_1 = 5;
const byte Din_2 = 6;
const byte Din_3 = 7;
const byte ADC_0 = 13;
const byte ADC_1 = 14;
const byte PWM_1 = 11;
const byte PWM_2 = 12;

#ifdef SPS_RCRECEIVER
const byte RC_0 = 15;
const byte RC_1 = 16;
#endif

#ifdef SPS_SERVO
const byte SERVO_1 = 11;
const byte SERVO_2 = 12;
#endif

const byte SW_PRG = 9;
const byte SW_SEL = 8;
#endif

#ifdef __AVR_ATtiny861__
// ATTiny4313 Hardware
const byte Dout_0 = 0;
const byte Dout_1 = 1;
const byte Dout_2 = 2;
const byte Dout_3 = 3;

const byte Din_0 = 4;
const byte Din_1 = 5;
const byte Din_2 = 6;
const byte Din_3 = 7;
const byte ADC_0 = 13;
const byte ADC_1 = 14;
const byte PWM_1 = 11;
const byte PWM_2 = 12;

#ifdef SPS_RCRECEIVER
const byte RC_0 = 15;
const byte RC_1 = 16;
#endif

#ifdef SPS_SERVO
const byte SERVO_1 = 11;
const byte SERVO_2 = 12;
#endif

const byte SW_PRG = 9;
const byte SW_SEL = 8;
#endif
