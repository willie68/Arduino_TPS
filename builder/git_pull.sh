echo get source files
cd /home/arduinocli
git clone https://github.com/willie68/Arduino_TPS.git

echo Task switch to right branch
cd /home/arduinocli/Arduino_TPS/
git checkout develop

cp ../build_tps.sh ./builder/build_tps.sh
cd /home/arduinocli/Arduino_TPS/TPS
mkdir -p /home/arduinocli/Arduino_TPS/dest
../builder/build_tps.sh