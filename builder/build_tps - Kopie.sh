TPS_VERSION="0.13.1"
export TPS_VERSION
echo start building attiny tps
rm -f /home/arduinocli/Arduino_TPS/dest/*
arduino-cli compile --clean -e -b ATTinyCore:avr:attinyx4 --output-dir /home/arduinocli/Arduino_TPS/dest/ ./ --build-property="build.extra_flags=-DTPS_ENHANCEMENT -DTPS_SERVO -DTPS_SERIAL_PRG "
cp /home/arduinocli/Arduino_TPS/dest/TPS.ino.hex /home/arduinocli/Arduino_TPS/dest/TPS.$TPS_VERSION.TINY.ENHANCEMENT.SERVO.SERIAL_PRG.hex
arduino-cli compile --clean -e -b ATTinyCore:avr:attinyx4 --output-dir /home/arduinocli/Arduino_TPS/dest/ ./ --build-property="build.extra_flags=-DTPS_ENHANCEMENT -DTPS_TONE"
cp /home/arduinocli/Arduino_TPS/dest/TPS.ino.hex /home/arduinocli/Arduino_TPS/dest/TPS.$TPS_VERSION.TINY.ENHANCEMENT.TONE.hex
arduino-cli compile --clean -e -b ATTinyCore:avr:attinyx4 --output-dir /home/arduinocli/Arduino_TPS/dest/ ./ --build-property="build.extra_flags="
cp /home/arduinocli/Arduino_TPS/dest/TPS.ino.hex /home/arduinocli/Arduino_TPS/dest/TPS.$TPS_VERSION.TINY.hex
rm -f /home/arduinocli/Arduino_TPS/dest/TPS.ino.*
mkdir -p /home/arduinocli/Arduino_TPS/dest/tiny
cp /home/arduinocli/Arduino_TPS/dest/*.hex /home/arduinocli/Arduino_TPS/dest/tiny/
rm -f /home/arduinocli/Arduino_TPS/dest/*.hex

echo start building arduino uno
rm -f /home/arduinocli/Arduino_TPS/dest/*
arduino-cli compile --clean -e -b arduino:avr:uno --output-dir /home/arduinocli/Arduino_TPS/dest/ ./ --build-property="build.extra_flags=-DTPS_ENHANCEMENT -DTPS_SERVO -DTPS_SERIAL_PRG -DTPS_TONE -DTPS_RCRECEIVER -DTPS_USE_DISPLAY"
cp /home/arduinocli/Arduino_TPS/dest/TPS.ino.hex /home/arduinocli/Arduino_TPS/dest/TPS.$TPS_VERSION.UNO.ENHANCEMENT.SERVO.TONE.SERIAL_PRG.RCRECEIVER.DISPLAY.hex
rm -f /home/arduinocli/Arduino_TPS/dest/TPS.ino.*
arduino-cli compile --clean -e -b arduino:avr:uno --output-dir /home/arduinocli/Arduino_TPS/dest/ ./ --build-property="build.extra_flags=-DTPS_ENHANCEMENT -DTPS_SERVO -DTPS_SERIAL_PRG -DTPS_TONE -DTPS_RCRECEIVER"
cp /home/arduinocli/Arduino_TPS/dest/TPS.ino.hex /home/arduinocli/Arduino_TPS/dest/TPS.$TPS_VERSION.UNO.ENHANCEMENT.SERVO.TONE.SERIAL_PRG.RCRECEIVER.hex
rm -f /home/arduinocli/Arduino_TPS/dest/TPS.ino.*
arduino-cli compile --clean -e -b arduino:avr:uno --output-dir /home/arduinocli/Arduino_TPS/dest/ ./ --build-property="build.extra_flags=-DTPS_ENHANCEMENT -DTPS_SERVO -DTPS_SERIAL_PRG -DTPS_TONE"
cp /home/arduinocli/Arduino_TPS/dest/TPS.ino.hex /home/arduinocli/Arduino_TPS/dest/TPS.$TPS_VERSION.UNO.ENHANCEMENT.SERVO.TONE.SERIAL_PRG.hex
rm -f /home/arduinocli/Arduino_TPS/dest/TPS.ino.*
arduino-cli compile --clean -e -b arduino:avr:uno --output-dir /home/arduinocli/Arduino_TPS/dest/ ./
cp /home/arduinocli/Arduino_TPS/dest/TPS.ino.hex /home/arduinocli/Arduino_TPS/dest/TPS.$TPS_VERSION.UNO.hex
rm -f /home/arduinocli/Arduino_TPS/dest/TPS.ino.*
mkdir -p /home/arduinocli/Arduino_TPS/dest/uno
cp /home/arduinocli/Arduino_TPS/dest/*.hex /home/arduinocli/Arduino_TPS/dest/uno/
rm -f /home/arduinocli/Arduino_TPS/dest/*.hex

echo start building arduino nano
rm -f /home/arduinocli/Arduino_TPS/dest/*
arduino-cli compile --clean -e -b arduino:avr:nano --output-dir /home/arduinocli/Arduino_TPS/dest/ ./ --build-property="build.extra_flags=-DTPS_ENHANCEMENT -DTPS_SERVO -DTPS_SERIAL_PRG -DTPS_TONE -DTPS_RCRECEIVER -DTPS_USE_DISPLAY"
cp /home/arduinocli/Arduino_TPS/dest/TPS.ino.hex /home/arduinocli/Arduino_TPS/dest/TPS.$TPS_VERSION.NANO.ENHANCEMENT.SERVO.TONE.SERIAL_PRG.RCRECEIVER.DISPLAY.hex
rm -f /home/arduinocli/Arduino_TPS/dest/TPS.ino.*
arduino-cli compile --clean -e -b arduino:avr:nano --output-dir /home/arduinocli/Arduino_TPS/dest/ ./ --build-property="build.extra_flags=-DTPS_ENHANCEMENT -DTPS_SERVO -DTPS_SERIAL_PRG -DTPS_TONE -DTPS_RCRECEIVER"
cp /home/arduinocli/Arduino_TPS/dest/TPS.ino.hex /home/arduinocli/Arduino_TPS/dest/TPS.$TPS_VERSION.NANO.ENHANCEMENT.SERVO.TONE.SERIAL_PRG.RCRECEIVER.hex
rm -f /home/arduinocli/Arduino_TPS/dest/TPS.ino.*
arduino-cli compile --clean -e -b arduino:avr:nano --output-dir /home/arduinocli/Arduino_TPS/dest/ ./ --build-property="build.extra_flags=-DTPS_ENHANCEMENT -DTPS_SERVO -DTPS_SERIAL_PRG -DTPS_TONE"
cp /home/arduinocli/Arduino_TPS/dest/TPS.ino.hex /home/arduinocli/Arduino_TPS/dest/TPS.$TPS_VERSION.NANO.ENHANCEMENT.SERVO.TONE.SERIAL_PRG.hex
rm -f /home/arduinocli/Arduino_TPS/dest/TPS.ino.*
arduino-cli compile --clean -e -b arduino:avr:nano --output-dir /home/arduinocli/Arduino_TPS/dest/ ./
cp /home/arduinocli/Arduino_TPS/dest/TPS.ino.hex /home/arduinocli/Arduino_TPS/dest/TPS.$TPS_VERSION.NANO.hex
rm -f /home/arduinocli/Arduino_TPS/dest/TPS.ino.*
mkdir -p /home/arduinocli/Arduino_TPS/dest/nano
cp /home/arduinocli/Arduino_TPS/dest/*.hex /home/arduinocli/Arduino_TPS/dest/nano/
rm -f /home/arduinocli/Arduino_TPS/dest/*.hex

echo start building arduino esp32 d1 mini
rm -f /home/arduinocli/Arduino_TPS/dest/*
arduino-cli compile --clean -e -b esp32:esp32:d1_mini32 --output-dir /home/arduinocli/Arduino_TPS/dest/ ./ --build-property="build.extra_flags=-DESP32"
cp /home/arduinocli/Arduino_TPS/dest/TPS.ino.bin /home/arduinocli/Arduino_TPS/dest/TPS.$TPS_VERSION.ESP32.D1.bin

arduino-cli compile --clean -e -b esp32:esp32:d1_mini32 --output-dir /home/arduinocli/Arduino_TPS/dest/ ./ --build-property="build.extra_flags=-DESP32 -DTPS_ENHANCEMENT"
cp /home/arduinocli/Arduino_TPS/dest/TPS.ino.bin /home/arduinocli/Arduino_TPS/dest/TPS.$TPS_VERSION.ESP32.D1.ENHANCEMENT.bin

arduino-cli compile --clean -e -b esp32:esp32:d1_mini32 --output-dir /home/arduinocli/Arduino_TPS/dest/ ./ --build-property="build.extra_flags=-DESP32 -DTPS_ENHANCEMENT -DTPS_SERVO -DTPS_TONE -DTPS_SERIAL_PRG"
cp /home/arduinocli/Arduino_TPS/dest/TPS.ino.hex /home/arduinocli/Arduino_TPS/dest/TPS.$TPS_VERSION.ESP32.D1.ENHANCEMENT.SERVO.TONE.SERIAL_PRG.bin
rm -f /home/arduinocli/Arduino_TPS/dest/TPS.ino.*
mkdir -p /home/arduinocli/Arduino_TPS/dest/esp32
cp /home/arduinocli/Arduino_TPS/dest/*.bin /home/arduinocli/Arduino_TPS/dest/esp32/
rm -f /home/arduinocli/Arduino_TPS/dest/TPS.*
