package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"time"

	flag "github.com/spf13/pflag"
	"github.com/willie68/tps_cc/pkg/model"
	"github.com/willie68/tps_cc/pkg/tmpl"
	"github.com/willie68/tps_cc/pkg/tpsfile"
	"github.com/willie68/tps_cc/pkg/tpsgen"
)

var SrcFile string
var TemplateCmds model.TemplateCmds
var Debug bool
var BuildZip bool
var AutoCompile bool
var Board string

func init() {
	log.Println("init")
	rand.Seed(time.Now().UnixNano())

	flag.StringVarP(&SrcFile, "source", "s", "", "this is the path and filename to the source file")
	flag.BoolVarP(&Debug, "Debug", "d", false, "this will set the debug flag in the source file")
	flag.BoolVarP(&BuildZip, "Zip", "z", false, "building a zip file from the gerenrated sources")
	flag.BoolVarP(&AutoCompile, "Compile", "c", false, "compile the gerenrated sources")
	flag.StringVarP(&Board, "board", "b", "arduino_uno", "board to use in compile task")

	dat, err := tmpl.TemplateFS.ReadFile("files/template.json")
	if err != nil {
		log.Panicf("can't read template.json: %v", err)
	}

	err = json.Unmarshal(dat, &TemplateCmds)
	if err != nil {
		log.Panicf("can't unmarshall template.json: %v", err)
	}
}

func main() {
	log.Println("starting")

	flag.Parse()

	// reading tps file
	if SrcFile == "" {
		log.Panicf("source file cant't be empty.")
	}

	name := strings.TrimSuffix(SrcFile, ".tps")
	name = filepath.Base(name)

	path := generateRamdomPath(filepath.Dir(SrcFile))

	// generating structures
	commandSrc, err := tpsfile.ParseFile(SrcFile)
	if err != nil {
		log.Panicf("parsing error: %v", err)
	}

	flags := tpsfile.GetBuildFlags(commandSrc)

	// generating sources
	tpsgen := tpsgen.TPSGen{
		Name:         name,
		Path:         path,
		TemplateCmds: TemplateCmds,
		Debug:        Debug,
		Board:        Board,
		BuildFlags:   flags,
	}
	err = tpsgen.Generate(commandSrc)
	if err != nil {
		log.Panicf("generating error: %v", err)
	}

	// compiling sources

	if AutoCompile {
		// arduino-cli compile --clean -e -b arduino:avr:uno --output-dir /home/arduinocli/Arduino_TPS/dest/ ./ --build-property="build.extra_flags=-DTPS_ENHANCEMENT -DTPS_SERVO -DTPS_SERIAL_PRG -DTPS_TONE" >>log.$TPS_VERSION.log 2>&1
		// cp /home/arduinocli/Arduino_TPS/dest/TPS.ino.hex /home/arduinocli/Arduino_TPS/dest/TPS.$TPS_VERSION.UNO.ENHANCEMENT.SERVO.TONE.SERIAL_PRG.hex
		log.Printf("starting compilation for %s", Board)
		r, err := tpsgen.Compile()
		if err != nil {
			log.Panicf("compile error: %s,\r\n %v", r, err)
		}
		log.Printf("compilation finnished with: %s", r)
	}

	// building a zip
	if BuildZip {
		zip, err := tpsgen.ZipFiles()
		if err != nil {
			log.Panicf("generating zip error: %v", err)
		}
		log.Printf("building Zip File: %s", zip)
	}
}

func generateRamdomPath(base string) string {
	path := base
	exists := true
	for exists {
		path = base + "/" + randSeq(8)
		exists = false
		if _, err := os.Stat(path); err == nil {
			exists = true
		}
	}

	os.MkdirAll(path, os.ModePerm)
	return path
}

var letters = []rune("abcdefghijklmnopqrstuvwxyz0123456789")

func randSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
