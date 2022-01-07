docker build ./ -t mcs/tpscc-service:V1
docker run -d --restart always --name tpscc-service -p 9643:8443 -p 9680:8080 mcs/tpscc-service:V1