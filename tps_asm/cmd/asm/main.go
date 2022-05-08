package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	flag "github.com/spf13/pflag"
	log "github.com/willie68/tps_asm/internal/logging"
)

type label struct {
	name       string
	prgcounter int
}

var (
	destination string
	tpsfile     string
	labels      map[string]label
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
	labels = make(map[string]label)
	//handle errors while opening
	if err != nil {
		log.Fatalf("Error when opening file: %s", err)
	}

	fileScanner := bufio.NewScanner(file)

	// read line by line
	for fileScanner.Scan() {
		line := fileScanner.Text()
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, ":") {
			parts := strings.Split(line, " ")
			labelName := strings.ToLower(parts[0])
			labels[labelName] = label{
				name:       labelName,
				prgcounter: 0,
			}
		}
		fmt.Println(line)
	}
	// handle first encountered error while reading
	if err := fileScanner.Err(); err != nil {
		log.Fatalf("Error while reading file: %s", err)
	}

	file.Close()
}
