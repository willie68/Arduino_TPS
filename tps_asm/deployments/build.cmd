@echo off
go build -ldflags="-s -w" -o tpsasm-service.exe cmd/service/main.go