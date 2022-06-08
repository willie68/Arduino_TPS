package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	flag "github.com/spf13/pflag"
	"github.com/willie68/tps_asm/internal/asm"
	log "github.com/willie68/tps_asm/internal/logging"
)

var (
	destination  string
	tpsfile      string
	includes     string
	outputfile   string
	outputformat string
	fs           *flag.FlagSet
)

func init() {
	// variables for parameter override
	fs = flag.NewFlagSet("main", flag.ContinueOnError)
	fs.StringVarP(&tpsfile, "file", "f", "", "source file to compile")
	fs.StringVarP(&destination, "destination", "d", "HOLTEK", "destination hardware to use. HOLTEK, ATMEGA8, ARDUINOTPS, TINYTPS")
	fs.StringVarP(&includes, "includes", "i", "", "base folder for inclusion")
	fs.StringVarP(&outputfile, "output", "o", "", "output file")
	fs.StringVarP(&outputformat, "format", "t", "", "output format: BIN, IntelHEX, TPS")
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
		Hardware:     asm.ParseHardware(destination),
		Source:       src,
		Includes:     includes,
		Filename:     tpsfile,
		Outputformat: strings.ToLower(outputformat),
	}
	errs := tpsasm.Parse()
	if len(errs) > 0 {
		for _, err := range errs {
			log.Errorf("%v", err)
		}
	}

	/*
		jstr, err := json.MarshalIndent(tpsasm, "", "  ")
		if err != nil {
			panic(err)
		}
		log.Infof("parse: \r\n%s", jstr)
	*/
	if outputfile != "" {
		f, err := os.Create(outputfile)
		if err != nil {
			panic(err)
		}
		defer f.Close()
		err = asm.TpsFile(f, tpsasm)
		if err != nil {
			panic(err)
		}
	}
}
