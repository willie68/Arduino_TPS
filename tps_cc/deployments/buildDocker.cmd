@echo off
docker build -f ./build/package/Dockerfile -t mcs/tpscc-service:V1
docker run --name tpscc-service -p 9543:8443 -p 9080:8080 mcs/tpscc-service:V1