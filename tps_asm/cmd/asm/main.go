package main

import (
	"bufio"
	"fmt"
	"os"

	flag "github.com/spf13/pflag"
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
	parseOne()
}

func parseOne() {
	file, err := os.Open(tpsfile)
	//handle errors while opening
	if err != nil {
		log.Fatalf("Error when opening file: %s", err)
	}

	fileScanner := bufio.NewScanner(file)

	// read line by line
	for fileScanner.Scan() {
		fmt.Println(fileScanner.Text())
	}
	// handle first encountered error while reading
	if err := fileScanner.Err(); err != nil {
		log.Fatalf("Error while reading file: %s", err)
	}

	file.Close()
}
