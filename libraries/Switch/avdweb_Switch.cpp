/*
Switch.cpp
Copyright (C) 2012  Albert van Dalen http://www.avdweb.nl
This program is free software: you can redistribute it and/or modify it under the terms of the GNU General Public License
as published by the Free Software Foundation, either version 3 of the License, or (at your option) any later version.
This program is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty
of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public License at http://www.gnu.org/licenses .

AUTHOR: Albert van Dalen
WEBSITE: http://www.avdweb.nl/arduino/hardware-interfacing/simple-switch-debouncer.html

HISTORY:
1.0.0    20-04-2013 _debouncePeriod=50
1.0.1    22-05-2013 Added longPress, doubleClick
1.0.2    01-12-2015 added process(input)
1.0.3    15-01-2016 added deglitching
1.0.5    25-01-2017 file renamed to avdweb_Switch
1.1.0    28-07-2018 added callbacks (code by Sean Lanigan, added by Martin Laclaustra)
1.2.0-rc 28-07-2017 added singleclick. Reorganize, keeping variables for each event in one function

..........................................DEGLITCHING..............................

                        ________________   _
               on      |                | | |    _
                       |                | | |   | |
                       |                |_| |___| |__
 analog        off_____|_____________________________|____________________________

                        ________________   _     _
 input            _____|                |_| |___| |_______________________________

 poll            ^ ^ ^ ^ ^ ^ ^ ^ ^ ^ ^ ^ ^ ^ ^ ^ ^ ^ ^ ^ ^ ^ ^ ^ ^ ^ ^ ^ ^ ^ ^ ^ ^

 equal           0 1 1 0 1 1 1 1 1 1 1 1 0 0 0 1 0 0 1 1 1 1 1 1 1 1 1 1 1 1 1 1 1

 deglitchPeriod          <--------><--   <--     <-  <--------><--------><--------
                                    ___________________________
 deglitched       _________________|                           |__________________

 deglitchTime            ^         ^     ^       ^   ^         ^        ^

 ..........................................DEBOUNCING.............................

 debouncePeriod                    <-------------------------------->
                                    _________________________________
 debounced        _________________|                                 |____________
                                    _                                 _
 _switched        _________________| |_______________________________| |__________

 switchedTime                      ^                                 ^


**********************************************************************************
........................................DOUBLE CLICK..............................

                           __________         ______
 debounced        ________|          |_______|      |_____________________________

 poll            ^ ^ ^ ^ ^ ^ ^ ^ ^ ^ ^ ^ ^ ^ ^ ^ ^ ^ ^ ^ ^ ^ ^ ^ ^ ^ ^ ^ ^ ^ ^ ^ ^
                           _                  _
 pushed          _________| |________________| |__________________________________

 pushedTime               ^                  ^
                                      _              _
 released        ____________________| |____________| |___________________________

 releasedTime                        ^              ^

 doubleClickPeriod         <------------------------------------->
                                              _
 _doubleClick     ___________________________| |__________________________________


........................................LONG PRESS................................

                           ___________________________
 debounced        ________|                           |___________________________

 poll            ^ ^ ^ ^ ^ ^ ^ ^ ^ ^ ^ ^ ^ ^ ^ ^ ^ ^ ^ ^ ^ ^ ^ ^ ^ ^ ^ ^ ^ ^ ^ ^ ^

 longPressPeriod            <--------------->
                            _                           _
 _switched        _________| |_________________________| |________________________
                                              __________
 longPressDisable ___________________________|          |_________________________
                                              _
 _longPress       ___________________________| |__________________________________


........................................SINGLE CLICK..............................

                           __________                                 ______
 debounced        ________|          |_______________________________|      |_____

 poll            ^ ^ ^ ^ ^ ^ ^ ^ ^ ^ ^ ^ ^ ^ ^ ^ ^ ^ ^ ^ ^ ^ ^ ^ ^ ^ ^ ^ ^ ^ ^ ^ ^

 longPressPeriod           <--------------->
 doubleClickPeriod         <------------------------------------->
                            _         _                               _      _
 _switched        _________| |_______| |_____________________________| |____| |___
                                                                  _____
 singleClickDisable______________________________________________|     |__________
                                                                  _
 _singleClick     _______________________________________________| |______________

*/

#include "Arduino.h"
#include "avdweb_Switch.h"

Switch::Switch(byte _pin, byte PinMode, bool polarity, int debouncePeriod, int longPressPeriod, int doubleClickPeriod, int deglitchPeriod):
pin(_pin), polarity(polarity), deglitchPeriod(deglitchPeriod), debouncePeriod(debouncePeriod), longPressPeriod(longPressPeriod), doubleClickPeriod(doubleClickPeriod)
{ pinMode(pin, PinMode);
  switchedTime = millis();
  debounced = digitalRead(pin);
  singleClickDisable = true;
}

bool Switch::poll()
{ input = digitalRead(pin);
  ms = millis();
  return process();
}

bool Switch::process()
{ deglitch();
  debounce();
  calcSingleClick();
  calcDoubleClick();
  calcLongPress();
  if(switched())
  { switchedTime = ms; //stores last times for future rounds
    if(pushed())
    { pushedTime = ms;
    } else { releasedTime = ms;
    }
  }
  triggerCallbacks();
  return _switched;
}

void inline Switch::deglitch()
{
  if(input == lastInput) equal = 1;
  else
  { equal = 0;
    deglitchTime = ms;
  }
  if(equal && ((ms - deglitchTime) > deglitchPeriod)) // max 50ms, disable deglitch: 0ms
  { deglitched = input;
    deglitchTime = ms;
  }
  lastInput = input;
}

void inline Switch::debounce()
{
  _switched = 0;
  if((deglitched != debounced) && ((ms - switchedTime) > debouncePeriod))
  { debounced = deglitched;
    _switched = 1;
  }
}

void inline Switch::calcSingleClick()
{ _singleClick = false;
  if(pushed())
  { if((ms - pushedTime) >= doubleClickPeriod)
    { singleClickDisable = false; //resets when pushed not in second click of doubleclick
    } else { singleClickDisable = true; //silence single click in second cl. doublecl.
    }
  }
  if(!singleClickDisable)
  { _singleClick = !switched() && !on() && ((releasedTime - pushedTime) <= longPressPeriod) && ((ms - pushedTime) >= doubleClickPeriod); // true just one time between polls
    singleClickDisable = _singleClick; // will be reset at next push
  }
}

void inline Switch::calcDoubleClick()
{ _doubleClick = pushed() && ((ms - pushedTime) < doubleClickPeriod);
}

void inline Switch::calcLongPress()
{ _longPress = false;
  if(released()) longPressDisable = false; //resets when released
  if(!longPressDisable)
  { _longPress = !switched() && on() && ((ms - pushedTime) > longPressPeriod); // true just one time between polls
    longPressDisable = _longPress; // will be reset at next release
  }
}

bool Switch::switched()
{ return _switched;
}

bool Switch::on()
{ return !(debounced^polarity);
}

bool Switch::pushed()
{ return _switched && !(debounced^polarity);
}

bool Switch::released()
{ return _switched && (debounced^polarity);
}

bool Switch::longPress()
{ return _longPress;
}

bool Switch::doubleClick()
{ return _doubleClick;
}

bool Switch::singleClick()
{ return _singleClick;
}

void Switch::triggerCallbacks()
{
  if(_pushedCallback && pushed())
  { _pushedCallback(_pushedCallbackParam);
  }
    else if(_releasedCallback && released())
  { _releasedCallback(_releasedCallbackParam);
  }

  if(_longPressCallback && longPress())
  { _longPressCallback(_longPressCallbackParam);
  }

  if(_doubleClickCallback && doubleClick())
  { _doubleClickCallback(_doubleClickCallbackParam);
  }

  if(_singleClickCallback && singleClick())
  { _singleClickCallback(_singleClickCallbackParam);
  }
}

void Switch::setPushedCallback(switchCallback_t cb, void* param)
{ /// Store the "pushed" callback function
  _pushedCallback = cb;
  _pushedCallbackParam = param;
}

void Switch::setReleasedCallback(switchCallback_t cb, void* param)
{ /// Store the "released" callback function
  _releasedCallback = cb;
  _releasedCallbackParam = param;
}

void Switch::setLongPressCallback(switchCallback_t cb, void* param)
{ /// Store the "long press" callback function
  _longPressCallback = cb;
  _longPressCallbackParam = param;
}

void Switch::setDoubleClickCallback(switchCallback_t cb, void* param)
{ /// Store the "double click" callback function
  _doubleClickCallback = cb;
  _doubleClickCallbackParam = param;
}

void Switch::setSingleClickCallback(switchCallback_t cb, void* param)
{ /// Store the "double click" callback function
  _singleClickCallback = cb;
  _singleClickCallbackParam = param;
}
