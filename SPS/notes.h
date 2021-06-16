#define MIDI_START 36 // here starts our table, lowest is C2
#define MIDI_NOTES 73

#if defined(__AVR_ATmega328P__) || defined(__AVR_ATtiny84__) || defined(__AVR_ATtiny861__) || defined(__AVR_ATtiny4313__)
#include "avr/pgmspace.h"
// converting midi notes into the desired frequwnz for the tone function
// the frequenz was calculated with libreoffice calc with this  fm = (440 Hz)* 2^((m−69)/12) and than rounded.
const PROGMEM unsigned int midiNoteToFreq[MIDI_NOTES] =
{
  65, // C2
  69, 73, 78, 82, 87, 92, 98, 104, 110, 117, 123,
  131, // C3
  139, 147, 156, 165, 175, 185, 196, 208, 220, 233, 247,
  262, //C4
  277, 294, 311, 330, 349, 370, 392, 415, 440, 466, 494,
  523, //C5
  554, 587, 622, 659, 698, 740, 784, 831, 880, 932, 988,
  1047, //C6
  1109, 1175, 1245, 1319, 1397, 1480, 1568, 1661, 1760, 1865, 1976,
  2093, //C7
  2217, 2349, 2489, 2637, 2794, 2960, 3136, 3322, 3520, 3729, 3951,
  4186 //C8
};

#define getFrequency(n) (pgm_read_word(n - MIDI_START + midiNoteToFreq))
#endif

#ifdef ESP32
const unsigned int midiNoteToFreq[MIDI_NOTES] =
{
  65, // C2
  69, 73, 78, 82, 87, 92, 98, 104, 110, 117, 123,
  131, // C3
  139, 147, 156, 165, 175, 185, 196, 208, 220, 233, 247,
  262, //C4
  277, 294, 311, 330, 349, 370, 392, 415, 440, 466, 494,
  523, //C5
  554, 587, 622, 659, 698, 740, 784, 831, 880, 932, 988,
  1047, //C6
  1109, 1175, 1245, 1319, 1397, 1480, 1568, 1661, 1760, 1865, 1976,
  2093, //C7
  2217, 2349, 2489, 2637, 2794, 2960, 3136, 3322, 3520, 3729, 3951,
  4186 //C8
};
#define getFrequency(n) (midiNoteToFreq[n-MIDI_START])
#endif
