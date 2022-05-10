@echo off
go build -ldflags="-s -w" -o tpsasm.exe cmd/asm/main.go
go build -ldflags="-s -w" -o tpsasm-service.exe cmd/service/main.go