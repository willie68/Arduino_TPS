".\3rd party\GoVersionSetter.exe" -i
".\3rd party\GoVersionSetter.exe" -e ino -f ./TPS/version.h -o "#ifndef TPS_VERSIN_H\r\n#define TPS_VERSION_H\r\n#define TPS_VERSION \"%%s\"\r\n#endif"
".\3rd party\GoVersionSetter.exe" -e txt -f ./builder/build_tps.sh -o "0,TPS_VERSION=\"%%s\""
rd /s /q dest_bck
xcopy /s /v /e dest\*  dest_bck\
rd /s /q dest

md dest
docker build ./ -t arduino-tps-builder 
docker run --name arduino-tps-builder arduino-tps-builder bash
docker cp arduino-tps-builder:/home/arduinocli/Arduino_TPS/dest/ ./
docker stop arduino-tps-builder
@echo off
echo if you like to look into the container than please break this script at this point. (CTRL-C)
echo after that simply start the container again with "docker run -it --name arduino-tps-builder arduino-tps-builder bash"
echo you will find the builded artefacts on this computer in ../dest and on the docker container in ../dest
pause

docker rm arduino-tps-builder