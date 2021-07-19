##### BUILDER #####

FROM ubuntu:21.04 as builder

## Task: Install build deps
ENV DEBIAN_FRONTEND=noninteractive
RUN apt-get update
RUN set -eux; \
    apt-get install --no-install-recommends --allow-change-held-packages -y\
    git \
    wget \ 
    python3 \
    python3-pip \
    python3-dev \
    python3-serial

RUN ln -s /usr/bin/python3 /usr/bin/python

WORKDIR /root
#RUN wget https://github.com/sgerrand/alpine-pkg-glibc/releases/download/2.30-r0/glibc-2.30-r0.apk
#RUN apk add --allow-untrusted glibc-2.30-r0.apk
RUN addgroup arduino && adduser arduinocli && adduser arduinocli arduino 

#USER arduinocli
WORKDIR /home/arduinocli
RUN mkdir -p /home/arduinocli/cli
WORKDIR /home/arduinocli/cli
RUN wget https://downloads.arduino.cc/arduino-cli/arduino-cli_latest_Linux_64bit.tar.gz
RUN tar zxvf arduino-cli_latest_Linux_64bit.tar.gz
RUN rm arduino-cli_latest_Linux_64bit.tar.gz
ENV PATH="/home/arduinocli/cli:${PATH}"

# Install AVR core
RUN arduino-cli config init
RUN arduino-cli config add board_manager.additional_urls https://dl.espressif.com/dl/package_esp32_index.json
RUN arduino-cli config add board_manager.additional_urls http://drazzy.com/package_drazzy.com_index.json

RUN arduino-cli core update-index

RUN arduino-cli core download esp32:esp32
RUN arduino-cli core download ATTinyCore:avr
RUN arduino-cli core install arduino:avr
RUN arduino-cli core install esp32:esp32
RUN arduino-cli core install ATTinyCore:avr

# install needed libraries
RUN arduino-cli lib install switch
RUN arduino-cli lib install Servo
RUN arduino-cli lib install ESP32Servo

## Task: get source files
WORKDIR /home/arduinocli
RUN git clone https://github.com/willie68/Arduino_TPS.git

## Task switch to right branch
WORKDIR /home/arduinocli/Arduino_TPS/
RUN git pull
RUN git checkout develop

COPY --chown=arduino:arduinocli ./build_tps.sh ./builder/build_tps.sh
RUN chmod u+x ./builder/build_tps.sh

COPY --chown=arduino:arduinocli ./git_pull.sh ./builder/git_pull.sh
RUN chmod u+x ./builder/git_pull.sh

## Task compile versions
WORKDIR /home/arduinocli/Arduino_TPS/TPS
RUN mkdir -p /home/arduinocli/Arduino_TPS/dest

RUN ../builder/git_pull.sh