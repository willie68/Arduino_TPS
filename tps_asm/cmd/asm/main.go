package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/marcinbor85/gohex"
	flag "github.com/spf13/pflag"
	"github.com/willie68/tps_asm/internal/asm"
	log "github.com/willie68/tps_asm/internal/logging"
)

var (
	destination string
	tpsfile     string
	includes    string
	outputfile  string
	fs          *flag.FlagSet
)

func init() {
	// variables for parameter override
	fs = flag.NewFlagSet("main", flag.ContinueOnError)
	fs.StringVarP(&tpsfile, "file", "f", "", "source file to compile")
	fs.StringVarP(&destination, "destination", "d", "HOLTEK", "destination hardware to use. HOLTEK, ATMEGA8, ARDUINOSPS, TINYSPS")
	fs.StringVarP(&includes, "includes", "i", "", "base folder for inclusion")
	fs.StringVarP(&outputfile, "output", "o", "", "output file")
	fs.SortFlags = false
}

func main() {
	fs.Parse(os.Args[1:])
	if tpsfile == "" {
		fmt.Fprintf(os.Stderr, "missing parameter for file. Usage: %s \r\n", filepath.Base(os.Args[0]))
		fs.PrintDefaults()
		os.Exit(-1)
	}
	log.Logger.Info("init tps asm")
	log.Logger.Info("configuration")
	log.Logger.Infof("hardware: %s", destination)
	log.Logger.Infof("file    : %s", tpsfile)
	log.Logger.Infof("includes: %s", includes)

	file, err := os.Open(tpsfile)
	//handle errors while opening
	if err != nil {
		log.Errorf("error reading file: %v", err)
		os.Exit(-1)
	}
	defer file.Close()

	fileScanner := bufio.NewScanner(file)
	src := make([]string, 0)
	for fileScanner.Scan() {
		src = append(src, fileScanner.Text())
		// handle first encountered error while reading
		if err := fileScanner.Err(); err != nil {
			log.Errorf("error reading file: %v", err)
			os.Exit(-1)
		}
	}

	tpsasm := asm.Assembler{
		Hardware: asm.ParseHardware(destination),
		Source:   src,
		Includes: includes,
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

	if outputfile != "" {
		err = os.WriteFile(outputfile, tpsasm.Binary, 0644)
		if err != nil {
			panic(err)
		}
	}
	mem := gohex.NewMemory()
	mem.AddBinary(0, tpsasm.Binary)

	mem.DumpIntelHex(os.Stdout, 8)
}
