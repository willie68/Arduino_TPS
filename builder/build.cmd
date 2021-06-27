docker build ./ -t arduino-tps-builder
docker run --name arduino-tps-builder arduino-tps-builder bash
docker cp arduino-tps-builder:/home/arduinocli/Arduino_TPS/dest/ ../
docker stop arduino-tps-builder
