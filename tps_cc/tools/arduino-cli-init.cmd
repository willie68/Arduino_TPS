echo off
echo Install AVR core
arduino-cli config init
arduino-cli config add board_manager.additional_urls https://dl.espressif.com/dl/package_esp32_index.json
arduino-cli config add board_manager.additional_urls http://drazzy.com/package_drazzy.com_index.json

arduino-cli core update-index

arduino-cli core download esp32:esp32
arduino-cli core download ATTinyCore:avr
arduino-cli core install arduino:avr
arduino-cli core install esp32:esp32
arduino-cli core install ATTinyCore:avr

echo install needed libraries
arduino-cli lib install switch
arduino-cli lib install Servo
arduino-cli lib install ESP32Servo
