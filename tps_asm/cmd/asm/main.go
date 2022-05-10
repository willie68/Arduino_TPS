package main

import (
	"encoding/json"

	flag "github.com/spf13/pflag"
	"github.com/willie68/tps_asm/internal/asm"
	log "github.com/willie68/tps_asm/internal/logging"
)

var (
	destination string
	tpsfile     string
)

func init() {
	// variables for parameter override
	log.Logger.Info("init tps asm")
	flag.StringVarP(&destination, "destination", "d", "HOLTEK", "destination hardware to use. HOLTEK, ATMEGA8, ARDUINOSPS, TINYSPS")
	flag.StringVarP(&tpsfile, "file", "f", "tps.tps", "source file to compile")
}

func main() {
	flag.Parse()
	log.Logger.Info("configuration")
	log.Logger.Infof("hardware: %s", destination)
	log.Logger.Infof("file    : %s", tpsfile)
	tpsasm := asm.Assembler{
		Hardware: asm.ParseHardware(destination),
		Source:   tpsfile,
	}
	errs := tpsasm.Parse()
	if errs != nil && len(errs) > 0 {
		for _, err := range errs {
			log.Errorf("%v", err)
		}
	}

	jstr, err := json.MarshalIndent(tpsasm, "", "  ")
	if err != nil {
		panic(err)
	}
	log.Infof("parse: \r\n%s", jstr)
}
