docker build ./ -t arduino-tps-builder
docker run --name arduino-tps-builder arduino-tps-builder bash
docker cp arduino-tps-builder:/home/arduinocli/Arduino_TPS/dest/ ../
docker stop arduino-tps-builder
@echo off
echo if you like to look into the container than please break this script at this point. (CTRL-C)
echo after that simply start the container again with "docker run -it --name arduino-tps-builder arduino-tps-builder bash"
echo you will find the builded artefacts on this computer in ../dest and on the docker container in ../dest
pause

docker rm arduino-tps-builder