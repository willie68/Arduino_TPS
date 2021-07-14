".\3rd party\GoVersionSetter.exe" -i -e ino -f ./TPS/version.h -o "#ifndef TPS_VERSIN_H\r\n#define TPS_VERSION_H\r\n#define TPS_VERSION \"%%s\"\r\n#endif"

cd builder 
start build.cmd
cd ..
