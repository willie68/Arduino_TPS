##### BUILDER #####

FROM golang:1.17-alpine3.14 as builder

## Task: Install build deps

# hadolint ignore=DL3018
RUN set -eux; \
    apk add --no-progress --quiet --no-cache --upgrade --virtual .build-deps \
        gcc \
        git \
        musl-dev

## Task: copy source files

COPY . /src
WORKDIR /src

## Task: fetch project deps

RUN go mod download

## Task: build project

ENV GOOS="linux"
ENV GOARCH="amd64"
ENV CGO_ENABLED="0"

RUN go build -ldflags="-s -w" -o tpscc-service cmd/service/main.go 

## Task: set permissions

RUN chmod 0755 /src/tpscc-service

## Task: runtime dependencies

# hadolint ignore=DL3018
RUN set -eux; \
    apk add --no-progress --quiet --no-cache --upgrade --virtual .run-deps \
        tzdata

# hadolint ignore=DL3018,SC2183,DL4006
RUN set -eu +x; \
    apk add --no-progress --quiet --no-cache --upgrade ncurses; \
    apk update --quiet; \
    printf '%30s\n' | tr ' ' -; \
    echo "RUNTIME DEPENDENCIES"; \
    PKGNAME=$(apk info --depends .run-deps \
        | sed '/^$/d;/depends/d' \
        | sort -u ); \
    printf '%s\n' "${PKGNAME}" \
        | while IFS= read -r pkg; do \
                apk info --quiet --description --no-network "${pkg}" \
                | sed -n '/description/p' \
                | sed -r "s/($(echo "${pkg}" | sed -r 's/\+/\\+/g'))-(.*)\s.*/\1=\2/"; \
                done \
        | tee -a /usr/share/rundeps; \
    printf '%30s\n' | tr ' ' - 


##### TARGET #####

FROM ubuntu:21.04

ARG RELEASE
ENV IMG_VERSION="${RELEASE}"

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


COPY --from=builder /src/tpscc-service /usr/local/bin/
COPY --from=builder /src/configs/service.yaml /config/
COPY --from=builder /src/configs/secret.yaml /config/
COPY --from=builder /usr/share/rundeps /usr/share/rundeps

ENTRYPOINT ["/usr/local/bin/tpscc-service"]
CMD ["--config","/config/service.yaml"]

EXPOSE 8080 8443

HEALTHCHECK --interval=30s --timeout=5s --retries=3 --start-period=10s \
  CMD wget -q -T 5 --spider http://localhost:8080/healthz

LABEL org.opencontainers.image.title="TPS-CC" \
      org.opencontainers.image.description="MCS TPS Cross Compiler" \
      org.opencontainers.image.version="${IMG_VERSION}" \
      org.opencontainers.image.source="https://github.com/willie68/Arduino_TPS" \
      org.opencontainers.image.vendor="MCS (www.rcarduino.de)" \
      org.opencontainers.image.authors="info@wk-music.de" \
      maintainer="MCS" \
      NAME="tpscc"

