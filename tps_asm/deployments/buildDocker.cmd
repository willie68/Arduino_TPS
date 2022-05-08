@echo off
docker build -f ./build/package/Dockerfile ./ -t mcs/tpsasm-service:V1
docker run --name tpsasm-service -p 9543:8443 -p 9080:8080 mcs/tpsasm-service:V1