TPS_VERSION="0.13.28"
export TPS_VERSION

echo start building attiny tps
rm -f /home/arduinocli/Arduino_TPS/dest/*

echo TPS.$TPS_VERSION.TINY.ENHANCEMENT.SERVO.SERIAL_PRG >>log.$TPS_VERSION.log
arduino-cli compile --clean -e -b ATTinyCore:avr:attinyx4 --output-dir /home/arduinocli/Arduino_TPS/dest/ ./ --build-property="build.extra_flags=-DTPS_ENHANCEMENT -DTPS_SERVO -DTPS_SERIAL_PRG " >>log.$TPS_VERSION.log 2>&1
cp /home/arduinocli/Arduino_TPS/dest/TPS.ino.hex /home/arduinocli/Arduino_TPS/dest/TPS.$TPS_VERSION.TINY.ENHANCEMENT.SERVO.SERIAL_PRG.hex

echo TPS.$TPS_VERSION.TINY.ENHANCEMENT.TONE >>log.$TPS_VERSION.log
arduino-cli compile --clean -e -b ATTinyCore:avr:attinyx4 --output-dir /home/arduinocli/Arduino_TPS/dest/ ./ --build-property="build.extra_flags=-DTPS_ENHANCEMENT -DTPS_TONE" >>log.$TPS_VERSION.log 2>&1
cp /home/arduinocli/Arduino_TPS/dest/TPS.ino.hex /home/arduinocli/Arduino_TPS/dest/TPS.$TPS_VERSION.TINY.ENHANCEMENT.TONE.hex

echo TPS.$TPS_VERSION.TINY >>log.$TPS_VERSION.log
arduino-cli compile --clean -e -b ATTinyCore:avr:attinyx4 --output-dir /home/arduinocli/Arduino_TPS/dest/ ./ --build-property="build.extra_flags=" >>log.$TPS_VERSION.log 2>&1
cp /home/arduinocli/Arduino_TPS/dest/TPS.ino.hex /home/arduinocli/Arduino_TPS/dest/TPS.$TPS_VERSION.TINY.hex
rm -f /home/arduinocli/Arduino_TPS/dest/TPS.ino.*
mkdir -p /home/arduinocli/Arduino_TPS/dest/tiny
cp /home/arduinocli/Arduino_TPS/dest/*.hex /home/arduinocli/Arduino_TPS/dest/tiny/
rm -f /home/arduinocli/Arduino_TPS/dest/*.hex

echo start building arduino uno
rm -f /home/arduinocli/Arduino_TPS/dest/*
echo TPS.$TPS_VERSION.UNO.ENHANCEMENT.SERVO.TONE.SERIAL_PRG.RCRECEIVER.DISPLAY  >>log.$TPS_VERSION.log
arduino-cli compile --clean -e -b arduino:avr:uno --output-dir /home/arduinocli/Arduino_TPS/dest/ ./ --build-property="build.extra_flags=-DTPS_ENHANCEMENT -DTPS_SERVO -DTPS_SERIAL_PRG -DTPS_TONE -DTPS_RCRECEIVER -DTPS_USE_DISPLAY" >>log.$TPS_VERSION.log 2>&1
cp /home/arduinocli/Arduino_TPS/dest/TPS.ino.hex /home/arduinocli/Arduino_TPS/dest/TPS.$TPS_VERSION.UNO.ENHANCEMENT.SERVO.TONE.SERIAL_PRG.RCRECEIVER.DISPLAY.hex
rm -f /home/arduinocli/Arduino_TPS/dest/TPS.ino.*

echo TPS.$TPS_VERSION.UNO.ENHANCEMENT.SERVO.TONE.SERIAL_PRG.RCRECEIVER >>log.$TPS_VERSION.log
arduino-cli compile --clean -e -b arduino:avr:uno --output-dir /home/arduinocli/Arduino_TPS/dest/ ./ --build-property="build.extra_flags=-DTPS_ENHANCEMENT -DTPS_SERVO -DTPS_SERIAL_PRG -DTPS_TONE -DTPS_RCRECEIVER" >>log.$TPS_VERSION.log 2>&1
cp /home/arduinocli/Arduino_TPS/dest/TPS.ino.hex /home/arduinocli/Arduino_TPS/dest/TPS.$TPS_VERSION.UNO.ENHANCEMENT.SERVO.TONE.SERIAL_PRG.RCRECEIVER.hex
rm -f /home/arduinocli/Arduino_TPS/dest/TPS.ino.*

echo TPS.$TPS_VERSION.UNO.ENHANCEMENT.SERVO.TONE.SERIAL_PRG >>log.$TPS_VERSION.log
arduino-cli compile --clean -e -b arduino:avr:uno --output-dir /home/arduinocli/Arduino_TPS/dest/ ./ --build-property="build.extra_flags=-DTPS_ENHANCEMENT -DTPS_SERVO -DTPS_SERIAL_PRG -DTPS_TONE" >>log.$TPS_VERSION.log 2>&1
cp /home/arduinocli/Arduino_TPS/dest/TPS.ino.hex /home/arduinocli/Arduino_TPS/dest/TPS.$TPS_VERSION.UNO.ENHANCEMENT.SERVO.TONE.SERIAL_PRG.hex
rm -f /home/arduinocli/Arduino_TPS/dest/TPS.ino.*

echo TPS.$TPS_VERSION.UNO >>log.$TPS_VERSION.log
arduino-cli compile --clean -e -b arduino:avr:uno --output-dir /home/arduinocli/Arduino_TPS/dest/ ./ >>log.$TPS_VERSION.log 2>&1
cp /home/arduinocli/Arduino_TPS/dest/TPS.ino.hex /home/arduinocli/Arduino_TPS/dest/TPS.$TPS_VERSION.UNO.hex
rm -f /home/arduinocli/Arduino_TPS/dest/TPS.ino.*
mkdir -p /home/arduinocli/Arduino_TPS/dest/uno
cp /home/arduinocli/Arduino_TPS/dest/*.hex /home/arduinocli/Arduino_TPS/dest/uno/
rm -f /home/arduinocli/Arduino_TPS/dest/*.hex

echo start building arduino nano
rm -f /home/arduinocli/Arduino_TPS/dest/*

echo TPS.$TPS_VERSION.NANO.ENHANCEMENT.SERVO.TONE.SERIAL_PRG.RCRECEIVER.DISPLAY >>log.$TPS_VERSION.log
arduino-cli compile --clean -e -b arduino:avr:nano --output-dir /home/arduinocli/Arduino_TPS/dest/ ./ --build-property="build.extra_flags=-DTPS_ENHANCEMENT -DTPS_SERVO -DTPS_SERIAL_PRG -DTPS_TONE -DTPS_RCRECEIVER -DTPS_USE_DISPLAY"  >>log.$TPS_VERSION.log 2>&1
cp /home/arduinocli/Arduino_TPS/dest/TPS.ino.hex /home/arduinocli/Arduino_TPS/dest/TPS.$TPS_VERSION.NANO.ENHANCEMENT.SERVO.TONE.SERIAL_PRG.RCRECEIVER.DISPLAY.hex
rm -f /home/arduinocli/Arduino_TPS/dest/TPS.ino.*

echo TPS.$TPS_VERSION.NANO.ENHANCEMENT.SERVO.TONE.SERIAL_PRG.RCRECEIVER >>log.$TPS_VERSION.log
arduino-cli compile --clean -e -b arduino:avr:nano --output-dir /home/arduinocli/Arduino_TPS/dest/ ./ --build-property="build.extra_flags=-DTPS_ENHANCEMENT -DTPS_SERVO -DTPS_SERIAL_PRG -DTPS_TONE -DTPS_RCRECEIVER" >>log.$TPS_VERSION.log 2>&1
cp /home/arduinocli/Arduino_TPS/dest/TPS.ino.hex /home/arduinocli/Arduino_TPS/dest/TPS.$TPS_VERSION.NANO.ENHANCEMENT.SERVO.TONE.SERIAL_PRG.RCRECEIVER.hex
rm -f /home/arduinocli/Arduino_TPS/dest/TPS.ino.*

echo TPS.$TPS_VERSION.NANO.ENHANCEMENT.SERVO.TONE.SERIAL_PRG >>log.$TPS_VERSION.log
arduino-cli compile --clean -e -b arduino:avr:nano --output-dir /home/arduinocli/Arduino_TPS/dest/ ./ --build-property="build.extra_flags=-DTPS_ENHANCEMENT -DTPS_SERVO -DTPS_SERIAL_PRG -DTPS_TONE" >>log.$TPS_VERSION.log 2>&1
cp /home/arduinocli/Arduino_TPS/dest/TPS.ino.hex /home/arduinocli/Arduino_TPS/dest/TPS.$TPS_VERSION.NANO.ENHANCEMENT.SERVO.TONE.SERIAL_PRG.hex
rm -f /home/arduinocli/Arduino_TPS/dest/TPS.ino.*

echo TPS.$TPS_VERSION.NANO >>log.$TPS_VERSION.log
arduino-cli compile --clean -e -b arduino:avr:nano --output-dir /home/arduinocli/Arduino_TPS/dest/ ./ >>log.$TPS_VERSION.log 2>&1
cp /home/arduinocli/Arduino_TPS/dest/TPS.ino.hex /home/arduinocli/Arduino_TPS/dest/TPS.$TPS_VERSION.NANO.hex
rm -f /home/arduinocli/Arduino_TPS/dest/TPS.ino.*
mkdir -p /home/arduinocli/Arduino_TPS/dest/nano
cp /home/arduinocli/Arduino_TPS/dest/*.hex /home/arduinocli/Arduino_TPS/dest/nano/
rm -f /home/arduinocli/Arduino_TPS/dest/*.hex

echo start building arduino esp32 d1 mini
rm -f /home/arduinocli/Arduino_TPS/dest/*

echo TPS.$TPS_VERSION.ESP32.D1 >>log.$TPS_VERSION.log
arduino-cli compile --clean -e -b esp32:esp32:d1_mini32 --output-dir /home/arduinocli/Arduino_TPS/dest/ ./ --build-property="build.extra_flags=-DESP32" >>log.$TPS_VERSION.log 2>&1
cp /home/arduinocli/Arduino_TPS/dest/TPS.ino.bin /home/arduinocli/Arduino_TPS/dest/TPS.$TPS_VERSION.ESP32.D1.bin
rm -f /home/arduinocli/Arduino_TPS/dest/TPS.ino.bin

echo TPS.$TPS_VERSION.ESP32.D1.ENHANCEMENT >>log.$TPS_VERSION.log
arduino-cli compile --clean -e -b esp32:esp32:d1_mini32 --output-dir /home/arduinocli/Arduino_TPS/dest/ ./ --build-property="build.extra_flags=-DESP32 -DTPS_ENHANCEMENT" >>log.$TPS_VERSION.log 2>&1
cp /home/arduinocli/Arduino_TPS/dest/TPS.ino.bin /home/arduinocli/Arduino_TPS/dest/TPS.$TPS_VERSION.ESP32.D1.ENHANCEMENT.bin
rm -f /home/arduinocli/Arduino_TPS/dest/TPS.ino.bin

echo TPS.$TPS_VERSION.ESP32.D1.ENHANCEMENT.SERVO.TONE.SERIAL_PRG >>log.$TPS_VERSION.log
arduino-cli compile --clean -e -b esp32:esp32:d1_mini32 --output-dir /home/arduinocli/Arduino_TPS/dest/ ./ --build-property="build.extra_flags=-DESP32 -DTPS_ENHANCEMENT -DTPS_SERVO -DTPS_TONE -DTPS_SERIAL_PRG" >>log.$TPS_VERSION.log 2>&1
cp /home/arduinocli/Arduino_TPS/dest/TPS.ino.bin /home/arduinocli/Arduino_TPS/dest/TPS.$TPS_VERSION.ESP32.D1.ENHANCEMENT.SERVO.TONE.SERIAL_PRG.bin
rm -f /home/arduinocli/Arduino_TPS/dest/TPS.ino.bin

echo start building arduino esp32 dev module
echo TPS.$TPS_VERSION.ESP32.DEV >>log.$TPS_VERSION.log
arduino-cli compile --clean -e -b esp32:esp32:esp32 --output-dir /home/arduinocli/Arduino_TPS/dest/ ./ --build-property="build.extra_flags=-DESP32" >>log.$TPS_VERSION.log 2>&1
cp /home/arduinocli/Arduino_TPS/dest/TPS.ino.bin /home/arduinocli/Arduino_TPS/dest/TPS.$TPS_VERSION.ESP32.DEV.bin
rm -f /home/arduinocli/Arduino_TPS/dest/TPS.ino.bin

echo TPS.$TPS_VERSION.ESP32.DEV.ENHANCEMENT >>log.$TPS_VERSION.log
arduino-cli compile --clean -e -b esp32:esp32:esp32 --output-dir /home/arduinocli/Arduino_TPS/dest/ ./ --build-property="build.extra_flags=-DESP32 -DTPS_ENHANCEMENT" >>log.$TPS_VERSION.log 2>&1
cp /home/arduinocli/Arduino_TPS/dest/TPS.ino.bin /home/arduinocli/Arduino_TPS/dest/TPS.$TPS_VERSION.ESP32.DEV.ENHANCEMENT.bin
rm -f /home/arduinocli/Arduino_TPS/dest/TPS.ino.bin

echo TPS.$TPS_VERSION.ESP32.DEV.ENHANCEMENT.SERVO.TONE.SERIAL_PRG >>log.$TPS_VERSION.log
arduino-cli compile --clean -e -b esp32:esp32:esp32 --output-dir /home/arduinocli/Arduino_TPS/dest/ ./ --build-property="build.extra_flags=-DESP32 -DTPS_ENHANCEMENT -DTPS_SERVO -DTPS_TONE -DTPS_SERIAL_PRG" >>log.$TPS_VERSION.log 2>&1
cp /home/arduinocli/Arduino_TPS/dest/TPS.ino.bin /home/arduinocli/Arduino_TPS/dest/TPS.$TPS_VERSION.ESP32.DEV.ENHANCEMENT.SERVO.TONE.SERIAL_PRG.bin
rm -f /home/arduinocli/Arduino_TPS/dest/TPS.ino.bin

rm -f /home/arduinocli/Arduino_TPS/dest/TPS.ino.*
mkdir -p /home/arduinocli/Arduino_TPS/dest/esp32
cp /home/arduinocli/Arduino_TPS/dest/*.bin /home/arduinocli/Arduino_TPS/dest/esp32/
rm -f /home/arduinocli/Arduino_TPS/dest/TPS.*

cp log.$TPS_VERSION.log /home/arduinocli/Arduino_TPS/dest/log.$TPS_VERSION.log
