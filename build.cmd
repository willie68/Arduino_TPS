".\3rd party\GoVersionSetter.exe" -i -e ino -f ./TPS/version.h -o "#ifndef TPS_VERSIN_H\r\n#define TPS_VERSION_H\r\n#define TPS_VERSION \"%%s\"\r\n#endif"
".\3rd party\GoVersionSetter.exe" -e txt -f ./builder/build_tps.sh -o "0,TPS_VERSION=\"%%s\""

cd builder 
call build.cmd
cd ..
